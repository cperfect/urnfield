package urnfield

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//PatternString for the regex to map and parse an URN
const PatternString = `urn:(?P<NID>[-A-Za-z0-9]+):(?P<NSS>[\/:._\-;A-Za-z0-9]+)(?P<QUERY>\?=[=*&_\-\w]+)?(?P<RESOLVERS>\?\+[=*&_\-\w]+)?(?P<FRAGMENT>#[_*\-/\(\)\w]+)?`

//the compile regex
var pattern = regexp.MustCompile(PatternString)

//Urn represents a parsed URN
//see https://tools.ietf.org/html/rfc8141
type Urn struct {
	//NID is the Namespace ID
	NID string
	//NSSSlashDelimiter indicates if this Urn uses / as delimiter in addition to :
	NSSSlashDelimiter bool
	//NSS holds the Name Space Specific strings in order
	NSS []string
	//Query holds the Query component if one exists
	Query map[string]string
	//Resolves holds the Resolvers component if one exists
	Resolvers map[string]string
	//Fragmeent holds the fragement componnent if one exsts
	Fragment string
}

//Parse an urn from a string returning the parsed Urn struct or an error
//If *any* separators are "/" then NSSSlashDelimiter will be true
func Parse(urn string) (Urn, error) {
	m := pattern.FindAllStringSubmatch(urn, -1)
	if m == nil {
		return Urn{}, errors.New("urn did not match pattern")
	} else if len(m) > 1 {
		return Urn{}, errors.New("too many urns in match")
	} else if len(m[0]) != 6 {
		return Urn{}, fmt.Errorf("not enough groups (%d) in match", len(m[0]))
	}
	mp := m[0]
	if mp[1] == "" {
		return Urn{}, errors.New("No NID")
	}

	u := Urn{
		NID: mp[1],
	}
	if mp[2] == "" {
		return Urn{}, errors.New("No NSS")
	} else if strings.Contains(mp[2], ":") {
		u.NSS = strings.Split(mp[2], ":")
	} else if strings.Contains(mp[2], "/") {
		u.NSSSlashDelimiter = true
		u.NSS = strings.Split(mp[2], "/")
	} else {
		u.NSS = []string{mp[2]}
	}

	if mp[3] != "" {
		u.Query = keyValuesToMap(strings.TrimPrefix(mp[3], "?="))
	}
	if mp[4] != "" {
		u.Resolvers = keyValuesToMap(strings.TrimPrefix(mp[4], "?+"))
	}
	if mp[5] != "" {
		u.Fragment = strings.TrimPrefix(mp[5], "#")
	}
	return u, nil
}

//utility method to break a q or r component string to a map
func keyValuesToMap(kvStr string) map[string]string {
	m := map[string]string{}
	if kvStr == "" {
		return m
	}
	kvs := strings.Split(kvStr, "&")
	for _, kv := range kvs {
		kva := strings.Split(kv, "=")
		if len(kva) == 1 {
			m[kva[0]] = ""
		} else {
			m[kva[0]] = kva[1]
		}
	}
	return m
}

//ToString converts the Urn struct to a string
//this is a synonym of Format
func (u Urn) ToString() string {
	return u.Format()
}

//Format a urn string from a Urn struct
//If NSSSlashDelimiter is true then all delimiters will be "/"
func (u Urn) Format() string {
	sb := strings.Builder{}
	sb.WriteString("urn:")
	sb.WriteString(u.NID)
	sb.WriteString(":")

	var nssDelim = ":"
	if u.NSSSlashDelimiter {
		nssDelim = "/"
	}
	sb.WriteString(strings.Join(u.NSS, nssDelim))

	if u.Query != nil && len(u.Query) > 0 {
		sb.WriteString("?=")
		var count = 0
		for k, v := range u.Query {
			if count > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(k)
			if len(v) > 0 {
				sb.WriteString("=")
				sb.WriteString(v)
			}
			count++
		}
	}
	if u.Resolvers != nil && len(u.Resolvers) > 0 {
		sb.WriteString("?+")
		var count = 0
		for k, v := range u.Resolvers {
			if count > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(k)
			if len(v) > 0 {
				sb.WriteString("=")
				sb.WriteString(v)
			}
			count++
		}
	}
	if len(u.Fragment) > 0 {
		sb.WriteString("#")
		sb.WriteString(u.Fragment)
	}
	return sb.String()
}
