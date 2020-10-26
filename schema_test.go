package urnfield

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

type testSchemaData struct {
	urnTestData
	*Schema
}

var (
	xTestSchema = &Schema{
		Description: "Test schema",
		Nid:         "x-test",
		NssSchema: &NssSchema{
			Description:      "First element",
			ElementValidator: RegexNssElementValidatorFunc(regexp.MustCompile(`\d{1,10}$`)),
			Next: &NssSchema{
				Description:      "Second element",
				ElementValidator: RegexNssElementValidatorFunc(regexp.MustCompile(`\w{1,6}$`)),
			},
		},
		// QComponentSchema: &MapSchema{
		// 	KeyValueSchemas: {
		// 		&KeyValuesSchema{
		// 			Description: "",
		// 			Key:         "foo",
		// 			Pattern:     regexp.MustCompile("bar"),
		// 		},
		// 		&KeyValuesSchema{
		// 			Description: "",
		// 			Key:         "quux",
		// 			Pattern:     regexp.MustCompile("wibble|ptuey"),
		// 		},
		// 	},
		// },
		//RComponentSchema: &MapSchema{},
	}
)

var testUrnSchemas = []testSchemaData{
	{
		urnTestData{
			Urn{
				Nid: "x-test",
				Nss: []string{"0451450523", "Addasd"},
			},
			"Good x-test simple",
			"urn:x-test:0451450523:Addasd",
			true,
		},
		xTestSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test",
				Nss: []string{"0451450523", "Addasd"},
			},
			"Good x-test with q and r components",
			"urn:x-test:0451450523:Addasd?=foo=bar&quux=wibble?+niii&sparrow=african",
			true,
		},
		xTestSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test",
				Nss: []string{"0sddssd", "Addasd"},
			},
			"Bad x-test - first element wrong",
			"urn:x-test:0sddssd:Addasd",
			false,
		},
		xTestSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test",
				Nss: []string{"0451450523"},
			},
			"Bad x-test - missing 2nd nss element",
			"urn:x-test:0451450523",
			false,
		},
		xTestSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test",
				Nss: []string{"0sddssd", "Addasd", "another"},
			},
			"Bad x-test - extra nss element",
			"urn:x-test:0sddssd:Addasd:another",
			false,
		},
		xTestSchema,
	},
}

func TestSchemas(t *testing.T) {
	for _, sd := range testUrnSchemas {
		err := sd.Schema.Validate(sd.urnString)
		if sd.shouldSucceed {
			require.NoError(t, err, "Test Schemas %s (%s) should not have produced error %s", sd.desc, sd.urnString, err)
		} else {
			require.Error(t, err, "Test Schemas %s (%s) should have produced an error", sd.desc, sd.urnString)
		}
	}
}
