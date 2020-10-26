package urnfield

import (
	"errors"
	"fmt"
	"regexp"
)

//Schema defines a valid urn in a specific Namespace
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
type NssElementValidator func(nsse string) error

//RegexNssElementValidatorFunc returns a NssElementValidator func based on the given regex pattern
func RegexNssElementValidatorFunc(pattern *regexp.Regexp) NssElementValidator {
	return func(nsse string) error {
		if !pattern.Match([]byte(nsse)) {
			return fmt.Errorf("Bad value for element: value %s should match %s", nsse, pattern.String())
		}
		return nil
	}
}

//NssSchema defines a valid Nss set for a specific scheme
type NssSchema struct {
	Description      string
	Next             *NssSchema
	ElementValidator NssElementValidator
}

func (ns *NssSchema) validate(nss []string) error {
	if nss == nil {
		return fmt.Errorf("No value for nss element %s)", ns.Description)
	}
	err := ns.ElementValidator(nss[0])
	if err != nil {
		return fmt.Errorf("Invalid value for element %s: %s", ns.Description, err)
	}
	nss = nss[1:]
	if ns.Next == nil {
		if len(nss) == 0 {
			return nil
		}
		return fmt.Errorf("Too many nss element(s) %s", nss)
	} else if len(nss) == 0 {
		return errors.New("Too few nss elements")
	}
	return ns.Next.validate(nss)

}
