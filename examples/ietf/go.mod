module github.com/cperfect/urnfield/examples/ietf

go 1.26

require (
	github.com/cperfect/urnfield/v2 v2.0.0
	github.com/gobwas/glob v0.2.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace points at the local library source during development.
// Remove and pin to a tagged release version once the library is published.
replace github.com/cperfect/urnfield/v2 => ../..
