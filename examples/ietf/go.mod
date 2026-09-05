module github.com/cperfect/urnfield/examples/ietf

go 1.26

require (
	github.com/cperfect/urnfield/v2 v2.0.0
	github.com/gobwas/glob v1.0.0
	github.com/stretchr/testify v1.12.1
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

// replace points at the local library source during development.
// Remove and pin to a tagged release version once the library is published.
replace github.com/cperfect/urnfield/v2 => ../..
