Library for using URN fields in go structs or func params
====================================================

WARNING WIP - use at your own risk

See https://tools.ietf.org/html/rfc8141

## How URNs work (RFC 8141)

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

