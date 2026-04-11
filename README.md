# `urnfield` is a go library for using URN fields in structs or function params

WARNING WIP - use at your own risk

## How URNs work (RFC 8141)
See https://tools.ietf.org/html/rfc8141

A URN (Uniform Resource Name) is a persistent, location-independent identifier. The full syntax is:

```
urn:<NID>:<NSS>[?+<resolvers>][?=<query>][#<fragment>]
```

- **`urn:`** — fixed scheme prefix
- **NID** (Namespace Identifier) — names the namespace, e.g. `isbn`, `ietf`, `payments`. Case-insensitive, 2–32 chars.
- **NSS** (Namespace-Specific String) — the actual identifier within that namespace. Structure is defined by the namespace; elements are commonly delimited by `:` or `/`.
- **Resolvers** (`?+`) — optional hints about how to locate or access the resource (e.g. a service endpoint). Not part of the identity.
- **Query** (`?=`) — optional key/value pairs for passing parameters. Also not part of the identity.
- **Fragment** (`#`) — optional fragment, same semantics as in URLs.

Two URNs are considered equivalent if their NID and NSS match (case-insensitively for the NID, case-sensitively for the NSS by default). The resolvers, query, and fragment components do not affect identity.

### Examples

```
urn:isbn:978-0-13-110362-7
urn:ietf:rfc:2648
urn:payments:account:au:123456:78901234?+resolve=https://api.example.com?=currency=AUD
```

---
## Usage

> The intent is that once a URN string is created and set it is immutable (except possibly the resolvers component).

## Examples

### Parsing and formatting

```go
// Parse a URN string into a Urn struct
u, err := urnfield.Parse("urn:ietf:rfc:2648")
if err != nil {
    log.Fatal(err)
}

fmt.Println(u.Nid)  // "ietf"
fmt.Println(u.Nss)  // ["rfc", "2648"]

// Format it back to a string
s, err := u.Format()
if err != nil {
    log.Fatal(err)
}
fmt.Println(s) // "urn:ietf:rfc:2648"
```

### Defining a schema

Schemas validate the structure of a namespace's NSS using a chain of element validators. The IETF namespace ([RFC 2648](https://tools.ietf.org/html/rfc2648)) accepts several sub-namespaces (`rfc`, `fyi`, `std`, `bcp`, `id`, `params`) plus any other string:

```go
var oneOrMoreDigits = &urnfield.NssSchema{
    Description:      "1*DIGIT",
    ElementValidator: urnfield.RegexNssElementValidatorFunc(regexp.MustCompile(`^\d+$`), nil),
}

var IETFSchema = &urnfield.Schema{
    Description: "IETF URN namespace (RFC 2648)",
    Nid:         "ietf",
    NssSchema: &urnfield.NssSchema{
        Description: "sub-namespace",
        ElementValidator: urnfield.ComplexOrNssElementValidatorFunc(
            []*urnfield.NssSchema{
                // rfc: 1*DIGIT  e.g. urn:ietf:rfc:2648
                {ElementValidator: urnfield.EqualsNssElementValidatorFunc("rfc", oneOrMoreDigits)},
                // fyi: 1*DIGIT  e.g. urn:ietf:fyi:20
                {ElementValidator: urnfield.EqualsNssElementValidatorFunc("fyi", oneOrMoreDigits)},
                // params: *    e.g. urn:ietf:params:xml:ns:allocationToken-1.0
                {
                    ElementValidator: urnfield.EqualsNssElementValidatorFunc("params",
                        &urnfield.NssSchema{
                            ElementValidator: urnfield.GlobNssElementValidatorFunc(glob.MustCompile("*")),
                        }),
                },
            },
        ),
    },
}
```

### Validating a URN against a schema

```go
err := IETFSchema.Validate("urn:ietf:rfc:2648")   // nil
err  = IETFSchema.Validate("urn:ietf:rfc:abc")    // error: "abc" is not digits
err  = IETFSchema.Validate("urn:isbn:123")         // error: NID mismatch

// Validate a pre-parsed Urn directly
u, _ := urnfield.Parse("urn:ietf:params:xml:ns:allocationToken-1.0")
err   = IETFSchema.ValidateUrn(u)                  // nil
```

A full working implementation of the IETF schema is available in [`examples/ietf/`](examples/ietf/).

## Use cases

**Referencing resources by identity without fetching them:**

```go
payment := Payment{
  amount:      Currency.New("AUD", 1000),
  //assume we have a payments urn schema
  fromAccount: urn.MustParse("urn:payments:account:banka:au:123456:78901234"),
  toAccount: fromAccount: urn.MustParse("urn:payments:account:bankb:uk:98-99-00:945-234B"),
}
```

**Concise claims in auth tokens** 
```json
{
  "iss": "urn:payments:processor:acme",
  "sub": "urn:payments:user:acme:u-4f92a1",
  "aud": "urn:payments:processor:banka",
  "exp": 1744329600,
  "nbf": 1744243200,
  "iat": 1744243200,
  "jti": "76588473-b530-4e2b-8693-992f55a6c5b1",
  "permissions": [
    "urn:payments:account:banka:au:123456:78901234:read",
    "urn:payments:account:banka:au:123456:78901234:transfer",
    "urn:payments:account:bankb:uk:98-99-00:945-234B:read"
  ]
}
```

All seven registered claims from RFC 7519 are shown: `iss` (issuer), `sub` (subject), `aud` (audience), `exp` (expiry), `nbf` (not before), `iat` (issued at), and `jti` (JWT ID). URNs work naturally for the identity claims — they're globally unique, self-describing, and carry no location coupling.

Each scope URN precisely identifies the resource and the permitted action, without needing a separate schema document to interpret it.

