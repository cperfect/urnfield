package urnfield

import (
	"errors"
	"fmt"
	"regexp"
)

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
type NssElementValidator func(nsse string) (hasNext bool, next *NssSchema, err error)

//RegexNssElementValidatorFunc returns a NssElementValidator func based on the given regex pattern
func RegexNssElementValidatorFunc(pattern *regexp.Regexp, next *NssSchema) NssElementValidator {
	return func(nsse string) (bool, *NssSchema, error) {
		if !pattern.Match([]byte(nsse)) {
			return false, nil, fmt.Errorf("Bad value for element: value %s should match %s", nsse, pattern.String())
		}
		return next != nil, next, nil
	}
}

func EqualsNssElementValidatorFunc(nssEquals string, next *NssSchema) NssElementValidator {
	return func(nsse string) (bool, *NssSchema, error) {
		if nssEquals != nsse {
			return false, nil, fmt.Errorf("Bad value for element: value %s should equal %s", nsse, nssEquals)
		}
		return next != nil, next, nil
	}
}

func SimpleOrNssElementValidatorFunc(alternatives map[string]*NssSchema) NssElementValidator {
	return func(nsse string) (bool, *NssSchema, error) {
		if next, ok := alternatives[nsse]; ok {
			return (next != nil), next, nil
		}
		return false, nil, fmt.Errorf("No matching alternative for nss element %s in %v", nsse, alternatives)
	}
}

//NssSchema defines a valid Nss set for a specific scheme
type NssSchema struct {
	Description      string
	ElementValidator NssElementValidator
}

func (ns *NssSchema) validate(nss []string) error {
	if nss == nil {
		return fmt.Errorf("No value for nss element %s)", ns.Description)
	}
	hasNext, next, err := ns.ElementValidator(nss[0])
	if err != nil {
		return fmt.Errorf("Invalid value for element %s: %s", ns.Description, err)
	}
	nss = nss[1:]
	if len(nss) == 0 {
		if hasNext {
			return errors.New("Too few nss elements")
		}
		return nil
	} else if !hasNext && len(nss) > 0 {
		return fmt.Errorf("Too many nss element(s) %s", nss)
	}
	return next.validate(nss)

}
