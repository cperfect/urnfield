package urnfield

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Pattern for the regex to map and parse an URN
const Pattern = `urn:(?P<NID>[-A-Za-z0-9]+):(?P<NSS>[\/:._\-;A-Za-z0-9]+)(?P<QUERY>\?=[=*&_\-\w]+)?(?P<RESOLVERS>\?\+[=*&_\-\w]+)?(?P<FRAGMENT>#[_*\-/\(\)\w]+)?`

// the compile regex
var urnRegex = regexp.MustCompile(Pattern)

// Urn represents a parsed URN
// see https://tools.ietf.org/html/rfc8141
type Urn struct {
	//Nid is the Namespace ID (Nid)
	Nid string
	//NssSlashDelimiter indicates if this Urn uses / as delimiter in addition to :
	NssSlashDelimiter bool
	//Nss holds the Name Space Specific (Nss) strings in order
	Nss []string
	//Query holds the Query component if one exists
	Query map[string][]string
	//Resolves holds the Resolvers component if one exists
	Resolvers map[string][]string
	//Fragmeent holds the fragement componnent if one exsts
	Fragment string
}

// Parse an urn from a string returning the parsed Urn struct or an error
// If *any* separators are "/" then NSSSlashDelimiter will be true
func Parse(urn string) (Urn, error) {
	m := urnRegex.FindAllStringSubmatch(urn, -1)
	if m == nil {
		return Urn{}, errors.New("urn did not match pattern")
	} else if len(m) > 1 {
		return Urn{}, errors.New("too many groups in match")
	} else if len(m[0]) != 6 {
		return Urn{}, fmt.Errorf("not enough groups (%d) in match", len(m[0]))
	}
	mp := m[0]
	if mp[1] == "" {
		return Urn{}, errors.New("No NID")
	}

	u := Urn{
		Nid: mp[1],
	}
	if mp[2] == "" {
		return Urn{}, errors.New("No NSS")
	} else if strings.Contains(mp[2], ":") {
		u.Nss = strings.Split(mp[2], ":")
	} else if strings.Contains(mp[2], "/") {
		u.NssSlashDelimiter = true
		u.Nss = strings.Split(mp[2], "/")
	} else {
		u.Nss = []string{mp[2]}
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

// utility method to break a q or r component string to a map
func keyValuesToMap(kvStr string) map[string][]string {
	m := map[string][]string{}
	if kvStr == "" {
		return m
	}
	kvs := strings.Split(kvStr, "&")
	for _, kv := range kvs {
		kva := strings.Split(kv, "=")
		if len(kva) == 1 {
			if m[kva[0]] == nil {
				m[kva[0]] = []string{}
			}
		} else {
			m[kva[0]] = append(m[kva[0]], kva[1])
		}
	}
	return m
}

// IsWellFormed checks for Well Formedness
func (u *Urn) IsWellFormed() error {
	//TODO better error handling
	if u == nil {
		return errors.New("Urn object is null")
	}
	if u.Nid == "" {
		return errors.New("NID is required")
	}
	if u.Nss == nil || len(u.Nss) == 0 {
		return errors.New("NSS is required")
	}
	return nil
}

// ToString converts the Urn struct to a string
// will throw an error if not well formed
// this is a synonym of Format
func (u Urn) ToString() (string, error) {
	return u.Format()
}

// Format a urn string from a Urn struct
// will throw an error if not well formed
// If NSSSlashDelimiter is true then all delimiters will be "/"
func (u Urn) Format() (string, error) {
	err := (&u).IsWellFormed()
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	sb.WriteString("urn:")
	sb.WriteString(u.Nid)
	sb.WriteString(":")

	var nssDelim = ":"
	if u.NssSlashDelimiter {
		nssDelim = "/"
	}
	sb.WriteString(strings.Join(u.Nss, nssDelim))

	if u.Query != nil && len(u.Query) > 0 {
		sb.WriteString("?=")
		writeKeyValuesMap(&sb, u.Query)
	}
	if u.Resolvers != nil && len(u.Resolvers) > 0 {
		sb.WriteString("?+")
		writeKeyValuesMap(&sb, u.Resolvers)
	}
	if len(u.Fragment) > 0 {
		sb.WriteString("#")
		sb.WriteString(u.Fragment)
	}
	return sb.String(), nil
}

// ordering the params is not part of the spec but makes testing easier!
func writeKeyValuesMap(sb *strings.Builder, m map[string][]string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for count, k := range keys {
		vs, hasVal := m[k]
		if count > 0 {
			sb.WriteString("&")
		}
		if hasVal {
			if len(vs) == 0 {
				sb.WriteString(k)
			} else if len(vs) > 0 {
				for vcount, v := range vs {
					if vcount > 0 {
						sb.WriteString("&")
					}
					sb.WriteString(k)
					sb.WriteString("=")
					sb.WriteString(v)
				}
			}
		}
	}
}
