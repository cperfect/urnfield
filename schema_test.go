package urnfield

import (
	"regexp"
	"testing"

	"github.com/gobwas/glob"
	"github.com/stretchr/testify/assert"
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
			Description: "First element",
			ElementValidator: RegexNssElementValidatorFunc(regexp.MustCompile(`^\d{1,10}$`),
				&NssSchema{
					Description:      "Second element",
					ElementValidator: RegexNssElementValidatorFunc(regexp.MustCompile(`^\w{1,6}$`), nil),
				}),
		},
	}

	xTestOrSchema = &Schema{
		Description: "Test Or schema",
		Nid:         "x-test-or",
		NssSchema: &NssSchema{
			Description: "First element alternatives",
			ElementValidator: SimpleOrNssElementValidatorFunc(
				map[string]*NssSchema{
					"foo": {
						Description:      "foo Second element",
						ElementValidator: RegexNssElementValidatorFunc(regexp.MustCompile(`^\w{1,6}$`), nil),
					},
					"bar": nil,
				},
			),
		},
	}

	xTestGlobSchema = &Schema{
		Description: "Test Or schema",
		Nid:         "x-test-glob",
		NssSchema: &NssSchema{
			Description: "First element",
			ElementValidator: EqualsNssElementValidatorFunc("zow",
				&NssSchema{
					Description:      "Second element",
					ElementValidator: GlobNssElementValidatorFunc(glob.MustCompile("**.peowww", ':')),
				}),
		},
	}

	xTestComplexOrSchema = &Schema{
		Description: "Test Complex Or schema",
		Nid:         "x-test-or-complex",
		NssSchema: &NssSchema{
			Description: "First element alternatives",
			ElementValidator: ComplexOrNssElementValidatorFunc(
				[]*NssSchema{
					{
						Description: "Foo first element",
						ElementValidator: EqualsNssElementValidatorFunc("foo", &NssSchema{
							Description:      "Foo Second element",
							ElementValidator: RegexNssElementValidatorFunc(regexp.MustCompile(`^\w{1,6}$`), nil),
						})},
					{
						Description:      "Foo first element",
						ElementValidator: EqualsNssElementValidatorFunc("bar", nil)},
				},
			),
		},
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
	{
		urnTestData{
			Urn{
				Nid: "x-test-or",
				Nss: []string{"bibble", "Addasd"},
			},
			"Bad x-test-or - bad first element",
			"urn:x-test-or:bibble:Addasd",
			false,
		},
		xTestOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or",
				Nss: []string{"foo", "Addasd"},
			},
			"Good x-test-or - foo",
			"urn:x-test-or:foo:Addasd",
			true,
		},
		xTestOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or",
				Nss: []string{"foo", "1234567"},
			},
			"Bad x-test-or - foo bad next element",
			"urn:x-test-or:foo:1234567",
			false,
		},
		xTestOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or",
				Nss: []string{"bar"},
			},
			"Good x-test-or - bar",
			"urn:x-test-or:bar",
			true,
		},
		xTestOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or",
				Nss: []string{"bar", "bad"},
			},
			"Bad x-test-or - bar - too many elements",
			"urn:x-test-or:bar:bad",
			false,
		},
		xTestOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-glob",
				Nss: []string{"zow", "zpoe"},
			},
			"Bad x-test-glob - bad second element",
			"urn:x-test-glob:zow:zpoe",
			false,
		},
		xTestGlobSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-glob",
				Nss: []string{"zow", "zsadasdasd.peowww"},
			},
			"Good x-test-glob",
			"urn:x-test-glob:zow:zsadasdasd.peowww",
			true,
		},
		xTestGlobSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-glob",
				Nss: []string{"zow", "zsadasdasd:dsaddas:asdadsads:asdsdds.peowww"},
			},
			"Good x-test-glob-multinns",
			"urn:x-test-glob:zow:zsadasdasd:dsaddas:asdadsads:asdsdds.peowww",
			true,
		},
		xTestGlobSchema,
	},

	{
		urnTestData{
			Urn{
				Nid: "x-test-or-complex",
				Nss: []string{"bibble", "Addasd"},
			},
			"Bad x-test-or-complex - bad first element",
			"urn:x-test-or-complex:bibble:Addasd",
			false,
		},
		xTestComplexOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or-complex",
				Nss: []string{"foo", "Addasd"},
			},
			"Good x-test-or-complex - foo",
			"urn:x-test-or-complex:foo:Addasd",
			true,
		},
		xTestComplexOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or-complex",
				Nss: []string{"foo", "1234567"},
			},
			"Bad x-test-or-complex - foo bad next element",
			"urn:x-test-or-complex:foo:1234567",
			false,
		},
		xTestComplexOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or-complex",
				Nss: []string{"bar"},
			},
			"Good x-test-or-complex - bar",
			"urn:x-test-or-complex:bar",
			true,
		},
		xTestComplexOrSchema,
	},
	{
		urnTestData{
			Urn{
				Nid: "x-test-or-complex",
				Nss: []string{"bar", "bad"},
			},
			"Bad x-test-or-complex - bar - too many elements",
			"urn:x-test-or-complex:bar:bad",
			false,
		},
		xTestComplexOrSchema,
	},
}

func TestSchemas(t *testing.T) {
	for _, sd := range testUrnSchemas {
		err := sd.Schema.Validate(sd.urnString)
		if sd.shouldSucceed {
			assert.NoError(t, err, "Test Schemas %s (%s) should not have produced error: %s", sd.desc, sd.urnString, err)
		} else {
			assert.Error(t, err, "Test Schemas %s (%s) should have produced an error", sd.desc, sd.urnString)
		}
	}
}
