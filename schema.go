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

//Schema defines a valid urn in a specific Namespace
//Note that schema does not validate Q, R, F components
type Schema struct {
	Description string
	Nid         string
	NssSchema   *NssSchema
}

//Validate checks the urn string against the schema
func (s *Schema) Validate(urn string) error {
	u, err := Parse(urn)
	if err != nil {
		return err
	}
	return s.ValidateUrn(u)
}

//ValidateUrn validates a Urn object
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

//NssElementValidator defines the type of funcs for validating NSS elements
//returns hasNext == true and a non-nil pointer to the schema for the next
//elements if more are expected else returns hasNext == false and a nil schema,
//or returns fals, nil and a non-nil error is something has gone wrong
type NssElementValidator func(nss []string) (nssRemainder []string, next *NssSchema, err error)

//GlobNssElementValidatorFunc returns a NssElementValidator using glob matching https://github.com/gobwas/glob
//this could be useful for matching groups or classes of resources
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

//RegexNssElementValidatorFunc returns a NssElementValidator func based on the given regex pattern
func RegexNssElementValidatorFunc(pattern *regexp.Regexp, next *NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if !pattern.Match([]byte(nss[0])) {
			return nil, nil, fmt.Errorf("Bad value for element: value %s should match %s", nss[0], pattern.String())
		}
		return nss[1:], next, nil
	}
}

//EqualsNssElementValidatorFunc matches nss exactly
func EqualsNssElementValidatorFunc(nssEquals string, next *NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if nssEquals != nss[0] {
			return nil, nil, fmt.Errorf("Bad value for element: value %s should equal %s", nss[0], nssEquals)
		}
		return nss[1:], next, nil
	}
}

//SimpleOrNssElementValidatorFunc matches a set of distinct alternatives via string equals
func SimpleOrNssElementValidatorFunc(alternatives map[string]*NssSchema) NssElementValidator {
	return func(nss []string) ([]string, *NssSchema, error) {
		if next, ok := alternatives[nss[0]]; ok {
			return nss[1:], next, nil
		}
		return []string{}, nil, fmt.Errorf("No matching alternative for nss element %s in %v", nss[0], alternatives)
	}
}

//ComplexOrNssElementValidatorFunc Matches a set of alternatives via
//any other Validator fund
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

//NssSchema defines a valid Nss set for a specific scheme
type NssSchema struct {
	Description      string
	ElementValidator NssElementValidator
}

//recursively validate the nss elements
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
