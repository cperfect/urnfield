package urnfield

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type urnTestData struct {
	Urn
	desc          string
	urnString     string
	shouldSucceed bool
}

//test data
//see https://www.iana.org/assignments/urn-namespaces/urn-namespaces.xhtml
//note that I am ordering the param strings to allow for testing
//TODO find a better way to doing this
//TODO failure tests!
var testUrns = []urnTestData{
	{
		Urn{
			Nid: "isbn",
			Nss: []string{"0451450523"},
		},
		"Good ISBN schema",
		"urn:isbn:0451450523",
		true,
	},
	{
		Urn{
			Nid: "isbn",
			Nss: []string{}, //empty!
		},
		"Bad ISBN schema - empt NSS",
		"urn:isbn",
		false,
	},
	{
		Urn{
			Nid: "isan",
			Nss: []string{"0000-0000-2CEA-0000-1-0000-0000-Y"},
		},
		"Good ISAN schema",
		"urn:isan:0000-0000-2CEA-0000-1-0000-0000-Y",
		true,
	},
	{
		Urn{
			Nid: "ISSN",
			Nss: []string{"0167-6423"},
		},
		"Good ISSN Schema",
		"urn:ISSN:0167-6423",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
		},
		"Good IETF Schema rfc",
		"urn:ietf:rfc:2648",
		true,
	},
	{
		Urn{
			Nid: "mpeg",
			Nss: []string{"mpeg7", "schema", "2001"},
		},
		"Good mpeg Schema",
		"urn:mpeg:mpeg7:schema:2001",
		true,
	},
	{
		Urn{
			Nid: "oid",
			Nss: []string{"2.16.840"},
		},
		"Good LDAP oid",
		"urn:oid:2.16.840",
		true,
	},
	{
		Urn{
			Nid: "uuid",
			Nss: []string{"6e8bc430-9c3a-11d9-9669-0800200c9a66"},
		},
		"Good uuid schema",
		"urn:uuid:6e8bc430-9c3a-11d9-9669-0800200c9a66",
		true,
	},
	{
		Urn{
			Nid: "nbn",
			Nss: []string{"de", "bvb", "19-146642"},
		},
		"Good nbn Schema",
		"urn:nbn:de:bvb:19-146642",
		true,
	},
	{
		Urn{
			Nid: "lex",
			Nss: []string{"eu", "council", "directive", "2010-03-09;2010-19-UE"},
		},
		"Good lex Schema",
		"urn:lex:eu:council:directive:2010-03-09;2010-19-UE",
		true,
	},
	{
		Urn{
			Nid: "lsid",
			Nss: []string{"zoobank.org", "pub", "CDC8D258-8F57-41DC-B560-247E17D3DC8C"},
		},
		"Good lsid Schema",
		"urn:lsid:zoobank.org:pub:CDC8D258-8F57-41DC-B560-247E17D3DC8C",
		true,
	},
	{
		Urn{
			Nid:               "lex",
			Nss:               []string{"eu", "council", "directive", "2010-03-09;2010-19-UE"},
			NssSlashDelimiter: true,
		},
		"Good lex Schema with slash",
		"urn:lex:eu/council/directive/2010-03-09;2010-19-UE",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Query: map[string][]string{
				"foo":  {"bar"},
				"quux": {"wibble"},
			},
		},
		"Good IETF Schema with q-component",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Resolvers: map[string][]string{
				"sparrow": {"african"},
				"niii":    {},
			},
		},
		"Good IETF Schema with resolvers component & no val param",
		"urn:ietf:rfc:2648?+niii&sparrow=african",
		true,
	},
	{
		Urn{
			Nid:      "ietf",
			Nss:      []string{"rfc", "2648"},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with fragment",
		"urn:ietf:rfc:2648#some/fragment()",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Query: map[string][]string{
				"foo":  {"bar"},
				"quux": {"wibble"},
			},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with query and fragment components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble#some/fragment()",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Resolvers: map[string][]string{
				"sparrow": {"african"},
				"niii":    {},
			},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with resolver and fragment components",
		"urn:ietf:rfc:2648?+niii&sparrow=african#some/fragment()",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Query: map[string][]string{
				"foo":  {"bar"},
				"quux": {"wibble"},
			},
			Resolvers: map[string][]string{
				"sparrow": {"african"},
				"niii":    {},
			},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with query, resolver and fragment components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble?+niii&sparrow=african#some/fragment()",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Query: map[string][]string{
				"foo":  {"bar"},
				"quux": {"wibble"},
			},
			Resolvers: map[string][]string{
				"sparrow": {"african"},
				"niii":    {},
			},
		},
		"Good IETF Schema with query and resolver components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble?+niii&sparrow=african",
		true,
	},
	{
		Urn{
			Nid:               "ietf",
			Nss:               []string{"rfc", "2648"},
			NssSlashDelimiter: true,
			Query: map[string][]string{
				"foo":  {"bar"},
				"quux": {"wibble"},
			},
			Resolvers: map[string][]string{
				"sparrow": {"african"},
				"niii":    {},
			},
		},
		"Good IETF Schema with query and resolver components and slash",
		"urn:ietf:rfc/2648?=foo=bar&quux=wibble?+niii&sparrow=african",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Query: map[string][]string{
				"foo":  {"bar"},
				"quux": {"*"},
			},
			Resolvers: map[string][]string{
				"sparrow": {"african"},
				"niii":    {"*"},
			},
			Fragment: "some/*",
		},
		"Good IETF Schema with query, resolvers and globs",
		"urn:ietf:rfc:2648?=foo=bar&quux=*?+niii=*&sparrow=african#some/*",
		true,
	},
	{
		Urn{
			Nid: "ietf",
			Nss: []string{"rfc", "2648"},
			Query: map[string][]string{
				"foo":  {"bar", "zoing"},
				"quux": {"*"},
			},
			Resolvers: map[string][]string{
				"sparrow": {"european", "african"},
				"niii":    {"*"},
			},
			Fragment: "some/*",
		},
		"Good IETF Schema with query, resolvers and globs",
		"urn:ietf:rfc:2648?=foo=bar&foo=zoing&quux=*?+niii=*&sparrow=european&sparrow=african#some/*",
		true,
	},
}

func TestParseUrns(t *testing.T) {
	for _, d := range testUrns {
		fmt.Printf("Test Urn Data %v\n", d)
		u, err := Parse(d.urnString)

		if d.shouldSucceed {
			require.NoError(t, err, "Parse Test case %s (%s) should not have produced error %s", d.desc, d.urnString, err)

			assert.Equal(t, d.Urn, u, "Parse Test case %s (%s) result not equal", d.desc, d.urnString)
		} else {
			assert.Error(t, err, "Test case %s (%s) should have produced an error", d.desc, d.urnString)
		}

		fmt.Printf("%#v\n", u)
	}
}

func TestFormatUrns(t *testing.T) {
	for _, d := range testUrns {
		fmt.Printf("Test Urn Data %v\n", d)
		s, err := d.Urn.Format()
		if d.shouldSucceed {
			require.NoError(t, err, "Format should not have errored")
			assert.Equal(t, d.urnString, s, "Format Test case %s (%v) result not equal", d.desc, d.urnString)
		} else {
			assert.Error(t, err, "Test Format case %s (%s) should have produced an error", d.desc, d.urnString)
		}

		fmt.Printf("%#v\n", s)
	}
}
