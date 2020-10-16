package urnfield

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type urnTestData struct {
	desc          string
	urn           string
	shouldSucceed bool
}

//test data
//see https://www.iana.org/assignments/urn-namespaces/urn-namespaces.xhtml
var testUrns = []urnTestData{
	urnTestData{
		"Good ISBN schema",
		"urn:isbn:0451450523",
		true,
	},
	urnTestData{
		"Good ISAN schema",
		"urn:isan:0000-0000-2CEA-0000-1-0000-0000-Y",
		true,
	},
	urnTestData{
		"Good ISSN Schema",
		"urn:ISSN:0167-6423",
		true,
	},
	urnTestData{
		"Good IETF Schema rfc",
		"urn:ietf:rfc:2648",
		true,
	},
	urnTestData{
		"Good mpeg Schema",
		"urn:mpeg:mpeg7:schema:2001",
		true,
	},
	urnTestData{
		"Good LDAP oid",
		"urn:oid:2.16.840",
		true,
	},
	urnTestData{
		"Good uuid schema",
		"urn:uuid:6e8bc430-9c3a-11d9-9669-0800200c9a66",
		true,
	},
	urnTestData{
		"Good nbn Schema",
		"urn:nbn:de:bvb:19-146642",
		true,
	},
	urnTestData{
		"Good lex Schema",
		"urn:lex:eu:council:directive:2010-03-09;2010-19-UE",
		true,
	},
	urnTestData{
		"Good lsid Schema",
		"urn:lsid:zoobank.org:pub:CDC8D258-8F57-41DC-B560-247E17D3DC8C",
		true,
	},
	urnTestData{
		"Good lex Schema with slash",
		"urn:lex:eu/council/directive/2010-03-09;2010-19-UE",
		true,
	},
	urnTestData{
		"Good IETF Schema with q-component",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble",
		true,
	},
	urnTestData{
		"Good IETF Schema with r-component & no val param",
		"urn:ietf:rfc:2648?+sparrow=african&niiii",
		true,
	},
	urnTestData{
		"Good IETF Schema with fragment",
		"urn:ietf:rfc:2648#some/fragment()",
		true,
	},
	urnTestData{
		"Good IETF Schema with query and fragment components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble#some/fragment()",
		true,
	},
	urnTestData{
		"Good IETF Schema with resolver and fragment components",
		"urn:ietf:rfc:2648?+sparrow=african&niiii#some/fragment()",
		true,
	},
	urnTestData{
		"Good IETF Schema with query, resolver and fragment components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble?+sparrow=african&niiii#some/fragment()",
		true,
	},
	urnTestData{
		"Good IETF Schema with query and resolver components",
		"urn:ietf:rfc:2648?=foo=bar&quux=wibble?+sparrow=african&niiii",
		true,
	},
	urnTestData{
		"Good IETF Schema with query and resolver components and slash",
		"urn:ietf:rfc/2648?=foo=bar&quux=wibble?+sparrow=african&niiii",
		true,
	},
	urnTestData{
		"Good IETF Schema with query and resolver components and globs",
		"urn:ietf:rfc:2648?=foo=bar&quux=*?+foo=bar&quux=*#some/*",
		true,
	},
}

func TestParseUrns(t *testing.T) {
	for _, d := range testUrns {
		fmt.Printf("Test Urn Data %v\n", d)
		u, err := Parse(d.urn)

		if d.shouldSucceed {
			assert.NoError(t, err, "Test case %s (%s) should not have produced error %s", d.desc, d.urn, err)
			//TODO test parsed values
			//assert.Equal(t, urn, u.ToString()) no worky as maps are not ordered DUH
		} else {
			assert.Error(t, err, "Test case %s (%s) should have produced an error", d.desc, d.urn)
		}

		fmt.Printf("%#v\n", u)
	}

}
