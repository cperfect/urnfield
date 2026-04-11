package ietf

import (
	"testing"

	"github.com/cperfect/urnfield"
	"github.com/stretchr/testify/assert"
)

/* Useful info

https://www.rfc-editor.org/rfc-index.html
https://www.ietf.org/download/rfc-index.txt
https://www.iana.org/assignments/params/params.xml#params-1
https://www.iana.org/assignments/xml-registry/xml-registry.xhtml
*/

type urnTestData struct {
	urnfield.Urn
	desc          string
	urnString     string
	shouldSucceed bool
}

type testSchemaData struct {
	urnTestData
	*urnfield.Schema
}

var (
	testIetfUrnData = []string{
		"urn:ietf:rfc:4350", //https://www.ietf.org/rfc/rfc4350.html"
		"urn:ietf:fyi:20",   //https://tools.ietf.org/html/fyi20
		"urn:ietf:std:90",   //https://tools.ietf.org/html/std90
		"urn:ietf:bcp:13",   //https://tools.ietf.org/html/bcp13
		"urn:ietf:id:draft-www-opsawg-yang-vpn-service-pm-02", //https://tools.ietf.org/html/draft-www-opsawg-yang-vpn-service-pm-02
		"urn:ietf:mtg:ietf55",                                 //RFC 6924 mtg sub-namespace
		"urn:ietf:params:xml:ns:allocationToken-1.0",          //https://www.iana.org/assignments/xml-registry/ns/allocationToken-1.0.txt
		"urn:ietf:any-string123",                              //anything else
	}
)

var testUrnSchemas = []testSchemaData{
	//"urn:ietf:rfc:4350", //https://www.ietf.org/rfc/rfc4350.html"
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"rfc", "4350"},
			},
			"Good rfc",
			"urn:ietf:rfc:4350",
			true,
		},
		Schema,
	},
	//"urn:ietf:fyi:20",   //https://tools.ietf.org/html/fyi20
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"fyi", "20"},
			},
			"Good fyi",
			"urn:ietf:fyi:20",
			true,
		},
		Schema,
	},
	//"urn:ietf:std:90",   //https://tools.ietf.org/html/std90
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"std", "90"},
			},
			"Good std",
			"urn:ietf:std:90",
			true,
		},
		Schema,
	},
	//"urn:ietf:bcp:13",   //https://tools.ietf.org/html/bcp13
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"bcp", "13"},
			},
			"Good bcp",
			"urn:ietf:bcp:13",
			true,
		},
		Schema,
	},
	//"urn:ietf:id:draft-www-opsawg-yang-vpn-service-pm-02", //https://tools.ietf.org/html/draft-www-opsawg-yang-vpn-service-pm-02
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"id", "draft-www-opsawg-yang-vpn-service-pm-02"},
			},
			"Good internet draft",
			"urn:ietf:id:draft-www-opsawg-yang-vpn-service-pm-02",
			true,
		},
		Schema,
	},
	//"urn:ietf:mtg:ietf55",
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"mtg", "ietf55"},
			},
			"Good mtg (meeting)",
			"urn:ietf:mtg:ietf55",
			true,
		},
		Schema,
	},
	//"urn:ietf:params:xml:ns:allocationToken-1.0",          //https://www.iana.org/assignments/xml-registry/ns/allocationToken-1.0.txt
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"params", "xml", "ns", "allocationToken-1.0"},
			},
			"Good params (xml ns)",
			"urn:ietf:params:xml:ns:allocationToken-1.0",
			true,
		},
		Schema,
	},
	//"urn:ietf:any-string123",
	{
		urnTestData{
			urnfield.Urn{
				Nid: "ietf",
				Nss: []string{"any-string123"},
			},
			"Good any string",
			"urn:ietf:any-string123",
			true,
		},
		Schema,
	},
}

func TestSchema(t *testing.T) {
	for _, u := range testIetfUrnData {
		err := Schema.Validate(u)
		assert.NoError(t, err, "%s did not validate: %s", u, err)
	}
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
