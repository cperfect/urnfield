package urnfield

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

//TODO - prevent use of nss delims in validator patterns - maybe allow in glob only?
//TODO - fix functor namings?

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

// GlobNssElementValidatorFunc returns an NssElementValidator that matches the remaining NSS
// elements (joined with ":") against the given glob pattern. Useful for matching groups or
// classes of resources. See https://github.com/gobwas/glob for pattern syntax.
func GlobNssElementValidatorFunc(glob glob.Glob) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if !glob.Match(strings.Join(nss, ":")) {
			//if next != nil {
			return nil, nil, fmt.Errorf("Bad value for element: value %s should match glob pattern", nss)
			//}
		}
		return []string{}, nil, nil
	}
}

// RegexNssElementValidatorFunc returns an NssElementValidator that matches the current NSS
// element against the given compiled regex pattern.
func RegexNssElementValidatorFunc(pattern *regexp.Regexp, next *NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if !pattern.Match([]byte(nss[0])) {
			return nil, nil, fmt.Errorf("Bad value for element: value %s should match %s", nss[0], pattern.String())
		}
		return nss[1:], next, nil
	}
}

// EqualsNssElementValidatorFunc returns an NssElementValidator that requires the current
// NSS element to equal nssEquals exactly.
func EqualsNssElementValidatorFunc(nssEquals string, next *NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if nssEquals != nss[0] {
			return nil, nil, fmt.Errorf("Bad value for element: value %s should equal %s", nss[0], nssEquals)
		}
		return nss[1:], next, nil
	}
}

// SimpleOrNssElementValidatorFunc returns an NssElementValidator that matches the current NSS
// element against a map of allowed string values, each mapping to the next NssSchema to use.
func SimpleOrNssElementValidatorFunc(alternatives map[string]*NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if next, ok := alternatives[nss[0]]; ok {
			return nss[1:], next, nil
		}
		return []string{}, nil, fmt.Errorf("No matching alternative for nss element %s in %v", nss[0], alternatives)
	}
}

// ComplexOrNssElementValidatorFunc returns an NssElementValidator that tries each of the
// provided NssSchemas in order, returning the result of the first one that succeeds.
func ComplexOrNssElementValidatorFunc(alternatives []*NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		for _, schema := range alternatives {
			nss, next, err := schema.ElementValidator(nss)
			if err == nil { //TODO counter intuitive - find a better way
				return nss, next, err
			}
		}
		return nil, nil, fmt.Errorf("No matching alternative for nss element %s in %v", nss[0], alternatives)
	}
}

// NssSchema defines the validation rules for a set of NSS elements.
type NssSchema struct {
	Description      string
	ElementValidator NssElementValidator
}

// recursively validate the nss elements
func (ns *NssSchema) validate(nss []string) error {
	if nss == nil {
		return fmt.Errorf("No value for nss element %s)", ns.Description)
	}
	nss, next, err := ns.ElementValidator(nss)
	if err != nil {
		return fmt.Errorf("Invalid value for element %s: %s", ns.Description, err)
	} else if next == nil {
		if nss != nil && len(nss) > 0 {
			return errors.New("Too many nss elements")
		}
		return nil
	}

	if nss == nil || len(nss) < 1 {
		return errors.New("Not enough nss elements")
	}
	return next.validate(nss)
}
