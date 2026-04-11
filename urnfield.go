package urnfield

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Pattern is the regex used to parse a complete URN string into its components.
// It is anchored with ^ and $ so it matches the full input string only; strings
// with surrounding content will not match.
const Pattern = `^urn:(?P<NID>[-A-Za-z0-9]+):(?P<NSS>[\/:._\-;A-Za-z0-9]+)(?P<QUERY>\?=[=*&_\-\w]+)?(?P<RESOLVERS>\?\+[=*&_\-\w]+)?(?P<FRAGMENT>#[_*\-/\(\)\w]+)?$`

var urnRegex = regexp.MustCompile(Pattern)

// Urn represents a parsed URN
// see https://tools.ietf.org/html/rfc8141
type Urn struct {
	// Nid is the Namespace Identifier.
	Nid string
	// NssSlashDelimiter indicates if this URN uses "/" as the NSS delimiter instead of ":".
	NssSlashDelimiter bool
	// Nss holds the Namespace-Specific String elements in order.
	Nss []string
	// Query holds the query component ("?=") if one exists.
	Query map[string][]string
	// Resolvers holds the resolvers component ("?+") if one exists.
	Resolvers map[string][]string
	// Fragment holds the fragment component ("#") if one exists.
	Fragment string
}

// Parse parses a complete URN string and returns the parsed Urn struct or an error.
// The input must be a standalone URN — strings with surrounding content will not match.
// If *any* NSS separators are "/" then NssSlashDelimiter will be true.
func Parse(urn string) (Urn, error) {
	// FindStringSubmatch returns nil if the anchored pattern does not match the
	// full input, so no further length checks for multiple matches are needed.
	mp := urnRegex.FindStringSubmatch(urn)
	if mp == nil {
		return Urn{}, errors.New("urn did not match pattern")
	} else if len(mp) != 6 {
		return Urn{}, fmt.Errorf("not enough groups (%d) in match", len(mp))
	}
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
	for kv := range strings.SplitSeq(kvStr, "&") {
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

// IsWellFormed reports whether u is well-formed per RFC 8141, returning a
// descriptive error if not. Returns an error if u is nil.
func (u *Urn) IsWellFormed() error {
	if u == nil {
		return errors.New("Urn object is nil")
	}
	if u.Nid == "" {
		return errors.New("NID is required")
	}
	if len(u.Nss) == 0 {
		return errors.New("NSS is required")
	}
	return nil
}

// ToString converts the Urn to its string representation.
// It is a synonym for Format; prefer Format for consistency.
func (u Urn) ToString() (string, error) {
	return u.Format()
}

// Format formats the Urn as a URN string per RFC 8141.
// Returns an error if the Urn is not well-formed.
// If NssSlashDelimiter is true, all NSS delimiters will be "/" instead of ":".
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

	if len(u.Query) > 0 {
		sb.WriteString("?=")
		writeKeyValuesMap(&sb, u.Query)
	}
	if len(u.Resolvers) > 0 {
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
		vs := m[k]
		if count > 0 {
			sb.WriteString("&")
		}
		if len(vs) == 0 {
			// key with no value — write bare key (e.g. "?+niii")
			sb.WriteString(k)
		} else {
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
