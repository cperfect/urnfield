package ietf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/* Useful info

https://www.rfc-editor.org/rfc-index.html
https://www.ietf.org/download/rfc-index.txt
https://www.iana.org/assignments/params/params.xml#params-1
https://www.iana.org/assignments/xml-registry/xml-registry.xhtml
*/

var (
	testIetfUrnData = []string{
		"urn:ietf:rfc:4350", //https://www.ietf.org/rfc/rfc4350.html"
		"urn:ietf:fyi:20",   //https://tools.ietf.org/html/fyi20
		"urn:ietf:std:90",   //https://tools.ietf.org/html/std90
		"urn:ietf:bcp:13",   //https://tools.ietf.org/html/bcp13
		"urn:ietf:id:draft-www-opsawg-yang-vpn-service-pm-02", //https://tools.ietf.org/html/draft-www-opsawg-yang-vpn-service-pm-02
		"urn:ietf:params:xml:ns:allocationToken-1.0",          //https://www.iana.org/assignments/xml-registry/ns/allocationToken-1.0.txt
		"urn:ietf:any-string123",                              //anything else
	}
)

func TestSchema(t *testing.T) {
	for _, u := range testIetfUrnData {
		err := Schema.Validate(u)
		assert.NoError(t, err, "%s did not validate: %s", u, err)
	}
}
