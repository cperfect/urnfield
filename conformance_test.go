// Conformance runner: drives the library against the language-neutral test
// vectors published by the urnfield specification (vendored as a git submodule
// under submodules/urnfield-spec). These fixtures are the executable contract
// for the spec version pinned by that submodule; passing them is what "conformant"
// means. See submodules/urnfield-spec/conformance/README.md for the fixture format.
//
// This is an external test package (urnfield_test): it exercises only the public
// API, exactly as a downstream consumer would.
package urnfield_test

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/cperfect/urnfield/v2"
	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
)

const conformanceDir = "submodules/urnfield-spec/conformance"

// model mirrors the parsed-URN shape used in the fixtures (parse.yaml `expected`
// and format.yaml `input`). It maps to urnfield.Urn.
type model struct {
	Nid               string              `yaml:"nid"`
	Nss               []string            `yaml:"nss"`
	NssSlashDelimiter bool                `yaml:"nssSlashDelimiter"`
	Query             map[string][]string `yaml:"query"`
	Resolvers         map[string][]string `yaml:"resolvers"`
	Fragment          string              `yaml:"fragment"`
}

func (m model) toUrn() urnfield.Urn {
	return urnfield.Urn{
		Nid:               m.Nid,
		Nss:               m.Nss,
		NssSlashDelimiter: m.NssSlashDelimiter,
		Query:             m.Query,
		Resolvers:         m.Resolvers,
		Fragment:          m.Fragment,
	}
}

// loadFixture reads and YAML-decodes one fixture file into dst.
func loadFixture(t *testing.T, name string, dst any) {
	t.Helper()
	path := filepath.Join(conformanceDir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (did you run `git submodule update --init --recursive`?)", path, err)
	}
	if err := yaml.Unmarshal(data, dst); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}

// normalizeMap converts a nil map to empty and nil value lists to empty so that
// the library's representation (which uses empty, non-nil slices for bare keys)
// compares equal to fixtures regardless of nil/empty distinctions.
func normalizeMap(m map[string][]string) map[string][]string {
	out := map[string][]string{}
	for k, v := range m {
		if v == nil {
			v = []string{}
		}
		out[k] = v
	}
	return out
}

func modelsEqual(got urnfield.Urn, want model) bool {
	if got.Nid != want.Nid || got.NssSlashDelimiter != want.NssSlashDelimiter || got.Fragment != want.Fragment {
		return false
	}
	gotNss, wantNss := got.Nss, want.Nss
	if len(gotNss) == 0 {
		gotNss = nil
	}
	if len(wantNss) == 0 {
		wantNss = nil
	}
	if !reflect.DeepEqual(gotNss, wantNss) {
		return false
	}
	return reflect.DeepEqual(normalizeMap(got.Query), normalizeMap(want.Query)) &&
		reflect.DeepEqual(normalizeMap(got.Resolvers), normalizeMap(want.Resolvers))
}

func TestConformanceParse(t *testing.T) {
	var fixture struct {
		Version string `yaml:"version"`
		Cases   []struct {
			Desc     string `yaml:"desc"`
			Input    string `yaml:"input"`
			Ok       bool   `yaml:"ok"`
			Expected *model `yaml:"expected"`
		} `yaml:"cases"`
	}
	loadFixture(t, "parse.yaml", &fixture)

	for _, c := range fixture.Cases {
		t.Run(c.Desc, func(t *testing.T) {
			got, err := urnfield.Parse(c.Input)
			if !c.Ok {
				if err == nil {
					t.Fatalf("Parse(%q) = %+v, want parse error", c.Input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) returned error %v, want success", c.Input, err)
			}
			if c.Expected == nil {
				t.Fatalf("ok case %q has no expected model", c.Desc)
			}
			if !modelsEqual(got, *c.Expected) {
				t.Fatalf("Parse(%q)\n got: %+v\nwant: %+v", c.Input, got, *c.Expected)
			}
		})
	}
}

func TestConformanceFormat(t *testing.T) {
	var fixture struct {
		Version string `yaml:"version"`
		Cases   []struct {
			Desc     string `yaml:"desc"`
			Input    model  `yaml:"input"`
			Ok       bool   `yaml:"ok"`
			Expected string `yaml:"expected"`
		} `yaml:"cases"`
	}
	loadFixture(t, "format.yaml", &fixture)

	for _, c := range fixture.Cases {
		t.Run(c.Desc, func(t *testing.T) {
			got, err := c.Input.toUrn().Format()
			if !c.Ok {
				if err == nil {
					t.Fatalf("Format(%+v) = %q, want error", c.Input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Format(%+v) returned error %v, want success", c.Input, err)
			}
			if got != c.Expected {
				t.Fatalf("Format(%+v) = %q, want %q", c.Input, got, c.Expected)
			}
		})
	}
}

// matcher mirrors the declarative element-matcher model from SPECIFICATION.md §9.
// `alternatives` is polymorphic (a map for oneOfStrings, a sequence for
// oneOfSubschemas), so it is captured as a raw node and decoded per type.
type matcher struct {
	Type         string    `yaml:"type"`
	Value        string    `yaml:"value"`
	Pattern      string    `yaml:"pattern"`
	Separators   []string  `yaml:"separators"`
	Alternatives yaml.Node `yaml:"alternatives"`
	Next         *matcher  `yaml:"next"`
}

// buildSchema translates a matcher tree into the library's NssSchema, using the
// public validator factories.
func buildSchema(t *testing.T, m *matcher) *urnfield.NssSchema {
	if m == nil {
		return nil
	}
	ns := &urnfield.NssSchema{Description: m.Type}
	switch m.Type {
	case "exact":
		ns.ElementValidator = urnfield.EqualsNssElementValidatorFunc(m.Value, buildSchema(t, m.Next))
	case "regex":
		ns.ElementValidator = urnfield.RegexNssElementValidatorFunc(regexp.MustCompile(m.Pattern), buildSchema(t, m.Next))
	case "oneOfStrings":
		alts := map[string]*urnfield.NssSchema{}
		content := m.Alternatives.Content
		for i := 0; i+1 < len(content); i += 2 {
			alts[content[i].Value] = nodeToSchema(t, content[i+1])
		}
		ns.ElementValidator = urnfield.SimpleOrNssElementValidatorFunc(alts)
	case "oneOfSubschemas":
		var subs []*urnfield.NssSchema
		for _, c := range m.Alternatives.Content {
			subs = append(subs, nodeToSchema(t, c))
		}
		ns.ElementValidator = urnfield.ComplexOrNssElementValidatorFunc(subs)
	case "glob":
		ns.ElementValidator = urnfield.GlobNssElementValidatorFunc(glob.MustCompile(m.Pattern, separatorsToRunes(m.Separators)...))
	default:
		t.Fatalf("unknown matcher type %q", m.Type)
	}
	return ns
}

// nodeToSchema decodes a YAML node (a matcher or an explicit null terminal) into
// an NssSchema.
func nodeToSchema(t *testing.T, n *yaml.Node) *urnfield.NssSchema {
	if n == nil || n.Tag == "!!null" {
		return nil
	}
	var m matcher
	if err := n.Decode(&m); err != nil {
		t.Fatalf("decode matcher node: %v", err)
	}
	return buildSchema(t, &m)
}

func separatorsToRunes(seps []string) []rune {
	var out []rune
	for _, s := range seps {
		for _, r := range s {
			out = append(out, r)
		}
	}
	return out
}

func TestConformanceValidate(t *testing.T) {
	var fixture struct {
		Version string `yaml:"version"`
		Schemas map[string]struct {
			Nid string   `yaml:"nid"`
			Nss *matcher `yaml:"nss"`
		} `yaml:"schemas"`
		Cases []struct {
			Desc   string `yaml:"desc"`
			Schema string `yaml:"schema"`
			Input  string `yaml:"input"`
			Ok     bool   `yaml:"ok"`
		} `yaml:"cases"`
	}
	loadFixture(t, "validate.yaml", &fixture)

	// Build each declared schema once.
	schemas := map[string]*urnfield.Schema{}
	for name, sd := range fixture.Schemas {
		schemas[name] = &urnfield.Schema{
			Nid:       sd.Nid,
			NssSchema: buildSchema(t, sd.Nss),
		}
	}

	for _, c := range fixture.Cases {
		t.Run(c.Desc, func(t *testing.T) {
			s, ok := schemas[c.Schema]
			if !ok {
				t.Fatalf("case references unknown schema %q", c.Schema)
			}
			err := s.Validate(c.Input)
			if c.Ok && err != nil {
				t.Fatalf("Validate(%q) against %q returned error %v, want success", c.Input, c.Schema, err)
			}
			if !c.Ok && err == nil {
				t.Fatalf("Validate(%q) against %q succeeded, want failure", c.Input, c.Schema)
			}
		})
	}
}

func TestConformanceEquals(t *testing.T) {
	var fixture struct {
		Version string `yaml:"version"`
		Cases   []struct {
			Desc       string `yaml:"desc"`
			A          string `yaml:"a"`
			B          string `yaml:"b"`
			Equivalent bool   `yaml:"equivalent"`
		} `yaml:"cases"`
	}
	loadFixture(t, "equals.yaml", &fixture)

	for _, c := range fixture.Cases {
		t.Run(c.Desc, func(t *testing.T) {
			a, err := urnfield.Parse(c.A)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.A, err)
			}
			b, err := urnfield.Parse(c.B)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.B, err)
			}
			if got := a.Equivalent(b); got != c.Equivalent {
				t.Fatalf("Equivalent(%q, %q) = %v, want %v", c.A, c.B, got, c.Equivalent)
			}
		})
	}
}
