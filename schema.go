package urnfield

import (
	"errors"
	"fmt"
	"regexp"
)

type Schema struct {
	Description string
	Nid         string
	NssSchema   *NssSchema
	// QComponentSchema *MapSchema
	// RComponentSchema *MapSchema
	// FComponentSchema *FComponentSchema
}

func (s *Schema) Validate(urn string) error {
	u, err := Parse(urn)
	if err != nil {
		return err
	}
	return s.ValidateUrn(u)
}

func (s *Schema) ValidateUrn(u Urn) error {
	if s.Nid != u.Nid {
		return fmt.Errorf("NID must be %s: got %s", s.Nid, u.Nid)
	}
	err := s.NssSchema.validate(u.Nss)
	if err != nil {
		return err
	}

	// if s.QComponentSchema != nil {
	// 	err = s.QComponentSchema.validate(u.Query)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else if u.Query != nil {
	// 	return errors.New("No schema for query component")
	// }

	// if s.RComponentSchema != nil {
	// 	err = s.RComponentSchema.validate(u.Resolvers)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else if u.Resolvers != nil {
	// 	return errors.New("No schema for resolver component")
	// }

	// if s.FComponentSchema != nil {
	// 	err = s.FComponentSchema.validate(u.Fragment)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else if s.FComponentSchema != nil {
	// 	return errors.New("No schema for resolver component")
	// }
	return nil
}

// type MapSchema struct {
// 	Required        bool
// 	KeyValueSchemas []*KeyValuesSchema
// }

// func (ms *MapSchema) validate(m map[string][]string) error {
// 	if m == nil && (ms.KeyValueSchemas != nil || len(ms.KeyValueSchemas) == 0) {
// 		return errors.New("Expected parameters")
// 	} else if ms.KeyValueSchemas == nil && (m != nil || len(m) == 0) {
// 		return errors.New("Did not expect parameters")
// 	}
// 	for _, kvs := range ms.KeyValueSchemas {
// 		vals, ok := m[kvs.Key]

// 		if !ok && kvs.Required {
// 			return fmt.Errorf("Key %s is required", kvs.Key)
// 		} else if len(vals) < kvs.MinVals {
// 			return fmt.Errorf("Not enough values for key %s: got %d < %d", kvs.Key, len(vals), kvs.MinVals)
// 		} else if kvs.MaxVals != nil && len(vals) > *kvs.MaxVals {
// 			return fmt.Errorf("Too many values for key %s: got %d > %d", kvs.Key, len(vals), kvs.MaxVals)
// 		}
// 		if kvs.Pattern != nil {
// 			for _, v := range vals {
// 				if !kvs.Pattern.Match([]byte(v)) {
// 					return fmt.Errorf("Invalid value for key %s: %s doesn't match %s", kvs.Key, v, kvs.Pattern.String())
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// type KeyValuesSchema struct {
// 	Description string
// 	Key         string
// 	Required    bool
// 	MinVals     int
// 	MaxVals     *int
// 	Pattern     *regexp.Regexp
// }

// func (kvs KeyValuesSchema) validate(key string, vals []string) error {
// 	if key != kvs.Key {
// 		return fmt.Errorf("Key %s is not expected: %s", kvs.Key, key)
// 	}
// 	if len(vals) < kvs.MinVals {
// 		return fmt.Errorf("Not enough values for key %s: %d", kvs.Key, len(vals))
// 	} else if kvs.MaxVals != nil && len(vals) > *kvs.MaxVals {
// 		return fmt.Errorf("Too many values for key %s: %d", kvs.Key, len(vals))
// 	}
// 	for _, v := range vals {
// 		if !kvs.Pattern.Match([]byte(v)) {
// 			return fmt.Errorf("Invalid value for key %s: got %s expected %s ", kvs.Key, v, kvs.Pattern.String())
// 		}
// 	}
// 	return nil
// }

// type FComponentSchema struct {
// 	Description string
// 	Required    bool
// 	Pattern     *regexp.Regexp
// }

// func (fc *FComponentSchema) validate(f string) error {
// 	if fc.Required && len(f) < 1 {
// 		return fmt.Errorf("Fragment %s is required but missing", fc.Description)
// 	}
// 	if !fc.Pattern.Match([]byte(f)) {
// 		return fmt.Errorf("Bad value for fragment %s: value %s should match %s", fc.Description, f, fc.Pattern.String())
// 	}
// 	return nil
// }

type NssElementValidator func(nsse string) error

func RegexNssElementValidatorFunc(pattern *regexp.Regexp) NssElementValidator {
	return func(nsse string) error {
		if !pattern.Match([]byte(nsse)) {
			return fmt.Errorf("Bad value for element: value %s should match %s", nsse, pattern.String())
		}
		return nil
	}
}

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
