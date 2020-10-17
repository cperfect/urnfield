package urnfield

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type urnTestData struct {
	Urn
	desc          string
	urnString     string
	shouldSucceed bool
}

//test data
//see https://www.iana.org/assignments/urn-namespaces/urn-namespaces.xhtml
var testUrns = []urnTestData{
	urnTestData{
		Urn{
			NID: "isbn",
			NSS: []string{"0451450523"},
		},
		"Good ISBN schema",
		"urn:isbn:0451450523",
		true,
	},
	urnTestData{
		Urn{
			NID: "isan",
			NSS: []string{"0000-0000-2CEA-0000-1-0000-0000-Y"},
		},
		"Good ISAN schema",
		"urn:isan:0000-0000-2CEA-0000-1-0000-0000-Y",
		true,
	},
	urnTestData{
		Urn{
			NID: "ISSN",
			NSS: []string{"0167-6423"},
		},
		"Good ISSN Schema",
		"urn:ISSN:0167-6423",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
		},
		"Good IETF Schema rfc",
		"urn:ietf:rfc:2648",
		true,
	},
	urnTestData{
		Urn{
			NID: "mpeg",
			NSS: []string{"mpeg7", "schema", "2001"},
		},
		"Good mpeg Schema",
		"urn:mpeg:mpeg7:schema:2001",
		true,
	},
	urnTestData{
		Urn{
			NID: "oid",
			NSS: []string{"2.16.840"},
		},
		"Good LDAP oid",
		"urn:oid:2.16.840",
		true,
	},
	urnTestData{
		Urn{
			NID: "uuid",
			NSS: []string{"6e8bc430-9c3a-11d9-9669-0800200c9a66"},
		},
		"Good uuid schema",
		"urn:uuid:6e8bc430-9c3a-11d9-9669-0800200c9a66",
		true,
	},
	urnTestData{
		Urn{
			NID: "nbn",
			NSS: []string{"de", "bvb", "19-146642"},
		},
		"Good nbn Schema",
		"urn:nbn:de:bvb:19-146642",
		true,
	},
	urnTestData{
		Urn{
			NID: "lex",
			NSS: []string{"eu", "council", "directive", "2010-03-09;2010-19-UE"},
		},
		"Good lex Schema",
		"urn:lex:eu:council:directive:2010-03-09;2010-19-UE",
		true,
	},
	urnTestData{
		Urn{
			NID: "lsid",
			NSS: []string{"zoobank.org", "pub", "CDC8D258-8F57-41DC-B560-247E17D3DC8C"},
		},
		"Good lsid Schema",
		"urn:lsid:zoobank.org:pub:CDC8D258-8F57-41DC-B560-247E17D3DC8C",
		true,
	},
	urnTestData{
		Urn{
			NID:               "lex",
			NSS:               []string{"eu", "council", "directive", "2010-03-09;2010-19-UE"},
			NSSSlashDelimiter: true,
		},
		"Good lex Schema with slash",
		"urn:lex:eu/council/directive/2010-03-09;2010-19-UE",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
			Query: map[string]string{
				"foo":  "bar",
				"quux": "wibble",
			},
		},
		"Good IETF Schema with q-component",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
			Resolvers: map[string]string{
				"sparrow": "african",
				"niii":    "",
			},
		},
		"Good IETF Schema with resolvers component & no val param",
		"urn:ietf:rfc:2648?+sparrow=african&niii",
		true,
	},
	urnTestData{
		Urn{
			NID:      "ietf",
			NSS:      []string{"rfc", "2648"},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with fragment",
		"urn:ietf:rfc:2648#some/fragment()",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
			Query: map[string]string{
				"foo":  "bar",
				"quux": "wibble",
			},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with query and fragment components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble#some/fragment()",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
			Resolvers: map[string]string{
				"sparrow": "african",
				"niii":    "",
			},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with resolver and fragment components",
		"urn:ietf:rfc:2648?+sparrow=african&niii#some/fragment()",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
			Query: map[string]string{
				"foo":  "bar",
				"quux": "wibble",
			},
			Resolvers: map[string]string{
				"sparrow": "african",
				"niii":    "",
			},
			Fragment: "some/fragment()",
		},
		"Good IETF Schema with query, resolver and fragment components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble?+sparrow=african&niii#some/fragment()",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
			Query: map[string]string{
				"foo":  "bar",
				"quux": "wibble",
			},
			Resolvers: map[string]string{
				"sparrow": "african",
				"niii":    "",
			},
		},
		"Good IETF Schema with query and resolver components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble?+sparrow=african&niii",
		true,
	},
	urnTestData{
		Urn{
			NID:               "ietf",
			NSS:               []string{"rfc", "2648"},
			NSSSlashDelimiter: true,
			Query: map[string]string{
				"foo":  "bar",
				"quux": "wibble",
			},
			Resolvers: map[string]string{
				"sparrow": "african",
				"niii":    "",
			},
		},
		"Good IETF Schema with query and resolver components and slash",
		"urn:ietf:rfc/2648?=foo=bar&quux=wibble?+sparrow=african&niii",
		true,
	},
	urnTestData{
		Urn{
			NID: "ietf",
			NSS: []string{"rfc", "2648"},
			Query: map[string]string{
				"foo":  "bar",
				"quux": "*",
			},
			Resolvers: map[string]string{
				"sparrow": "african",
				"niii":    "*",
			},
			Fragment: "some/*",
		},
		"Good IETF Schema with query, resolvers and globs",
		"urn:ietf:rfc:2648?=foo=bar&quux=*?+sparrow=african&niii=*#some/*",
		true,
	},
}

func TestParseUrns(t *testing.T) {
	for _, d := range testUrns {
		fmt.Printf("Test Urn Data %v\n", d)
		u, err := Parse(d.urnString)

		if d.shouldSucceed {
			assert.NoError(t, err, "Parse Test case %s (%s) should not have produced error %s", d.desc, d.urnString, err)

			assert.Equal(t, d.Urn, u, "Parse Test case %s (%s) result not equal", d.desc, d.urnString)
		} else {
			assert.Error(t, err, "Test case %s (%s) should have produced an error", d.desc, d.urnString)
		}

		fmt.Printf("%#v\n", u)
	}

}
