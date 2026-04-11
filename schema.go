package urnfield

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

// Schema defines a valid URN in a specific namespace.
// Note that Schema does not validate the query, resolvers, or fragment components.
type Schema struct {
	Description string
	Nid         string
	NssSchema   *NssSchema
}

// Validate parses urn and checks it against the schema.
func (s *Schema) Validate(urn string) error {
	u, err := Parse(urn)
	if err != nil {
		return err
	}
	return s.ValidateUrn(u)
}

// ValidateUrn validates a parsed Urn against the schema.
func (s *Schema) ValidateUrn(u Urn) error {
	if s.Nid != u.Nid {
		return fmt.Errorf("NID must be %s: got %s", s.Nid, u.Nid)
	}
	err := s.NssSchema.validate(u.Nss)
	if err != nil {
		return err
	}

	return nil
}

// NssElementValidator is a function that validates one or more NSS elements.
// It returns the remaining unprocessed elements and the NssSchema to use for
// the next element, or a non-nil error if validation fails. When no further
// elements are expected, next is nil and nssRemainder should be empty.
type NssElementValidator func(nss []string) (nssRemainder []string, next *NssSchema, err error)

// GlobNssElementValidatorFunc returns an NssElementValidator that matches all
// remaining NSS elements (joined with ":") against the given glob pattern.
// This validator always terminates the chain — it consumes every remaining
// element in one match, so no next NssSchema is accepted. This is intentional:
// glob patterns are used for opaque or arbitrarily deep sub-namespaces (e.g.
// the IETF "params" sub-namespace) where per-element structure is not defined.
// For single-element pattern matching use RegexNssElementValidatorFunc instead.
// See https://github.com/gobwas/glob for pattern syntax.
func GlobNssElementValidatorFunc(glob glob.Glob) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if len(nss) == 0 {
			return nil, nil, fmt.Errorf("expected element matching glob pattern but got none")
		}
		if !glob.Match(strings.Join(nss, ":")) {
			return nil, nil, fmt.Errorf("bad value for element: value %s should match glob pattern", nss)
		}
		return []string{}, nil, nil
	}
}

// RegexNssElementValidatorFunc returns an NssElementValidator that matches the current NSS
// element against the given compiled regex pattern. Patterns operate on individual
// pre-split elements — do not include ":" or "/" in the pattern, as these delimiters
// are consumed by Parse before validation.
func RegexNssElementValidatorFunc(pattern *regexp.Regexp, next *NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if len(nss) == 0 {
			return nil, nil, fmt.Errorf("expected element matching %s but got none", pattern.String())
		}
		if !pattern.MatchString(nss[0]) {
			return nil, nil, fmt.Errorf("bad value for element: value %s should match %s", nss[0], pattern.String())
		}
		return nss[1:], next, nil
	}
}

// EqualsNssElementValidatorFunc returns an NssElementValidator that requires the current
// NSS element to equal nssEquals exactly. The value operates on a single pre-split
// element — do not include ":" or "/" in nssEquals.
func EqualsNssElementValidatorFunc(nssEquals string, next *NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if len(nss) == 0 {
			return nil, nil, fmt.Errorf("expected element equal to %s but got none", nssEquals)
		}
		if nssEquals != nss[0] {
			return nil, nil, fmt.Errorf("bad value for element: value %s should equal %s", nss[0], nssEquals)
		}
		return nss[1:], next, nil
	}
}

// SimpleOrNssElementValidatorFunc returns an NssElementValidator that matches the current NSS
// element against a map of allowed string values, each mapping to the next NssSchema to use.
// Use this when each alternative is a fixed string (e.g. "rfc", "fyi", "std"). For alternatives
// that require their own sub-validation logic, use ComplexOrNssElementValidatorFunc instead.
func SimpleOrNssElementValidatorFunc(alternatives map[string]*NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if len(nss) == 0 {
			return nil, nil, fmt.Errorf("expected one of %v but got none", alternatives)
		}
		if next, ok := alternatives[nss[0]]; ok {
			return nss[1:], next, nil
		}
		return []string{}, nil, fmt.Errorf("no matching alternative for nss element %s in %v", nss[0], alternatives)
	}
}

// ComplexOrNssElementValidatorFunc returns an NssElementValidator that tries each of the
// provided NssSchemas in order, returning the result of the first one that succeeds.
// Use this when alternatives need their own validation logic beyond a fixed string match
// (e.g. "rfc" followed by digits, vs "params" followed by an opaque glob). For simple
// fixed-string branching use SimpleOrNssElementValidatorFunc instead.
func ComplexOrNssElementValidatorFunc(alternatives []*NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if len(nss) == 0 {
			return nil, nil, fmt.Errorf("expected one of %v but got none", alternatives)
		}
		for _, schema := range alternatives {
			nss, next, err := schema.ElementValidator(nss)
			if err != nil {
				continue
			}
			return nss, next, nil
		}
		return nil, nil, fmt.Errorf("no matching alternative for nss element %s in %v", nss[0], alternatives)
	}
}

// NssSchema defines the validation rules for a set of NSS elements.
type NssSchema struct {
	Description      string
	ElementValidator NssElementValidator
}

// recursively validate the nss elements
func (ns *NssSchema) validate(nss []string) error {
	// len(nil) == 0, so this covers both nil and empty slices
	if len(nss) == 0 {
		return fmt.Errorf("no value for nss element %s", ns.Description)
	}
	nss, next, err := ns.ElementValidator(nss)
	if err != nil {
		return fmt.Errorf("invalid value for element %s: %s", ns.Description, err)
	} else if next == nil {
		if len(nss) > 0 {
			return errors.New("too many nss elements")
		}
		return nil
	}

	if len(nss) < 1 {
		return errors.New("not enough nss elements")
	}
	return next.validate(nss)
}
