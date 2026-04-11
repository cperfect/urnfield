package ietf

import (
	"regexp"

	"github.com/cperfect/urnfield"
	"github.com/gobwas/glob"
)

// 1*DIGIT
var oneOrMoreDigitsNssSchema = &urnfield.NssSchema{
	Description:      "1*DIGIT",
	ElementValidator: urnfield.RegexNssElementValidatorFunc(regexp.MustCompile(`^\d+$`), nil),
}

// string = 1*(DIGIT / ALPHA / "-")
var stringNssSchema = &urnfield.NssSchema{
	Description:      "string",
	ElementValidator: urnfield.RegexNssElementValidatorFunc(regexp.MustCompile(`^[0-9a-zA-Z-]+$`), nil),
}

// Schema is the urnfield schema for the IETF URN namespace as defined in
// https://tools.ietf.org/html/rfc2648, updated by https://tools.ietf.org/html/rfc3553.
// RFC 6924 (https://tools.ietf.org/html/rfc6924) formalises the IANA registry for
// sub-namespaces but does not change their formats; future additions require IETF Review.
var Schema = &urnfield.Schema{
	Description: "Schema for the IETF URN namespace (RFC 2648, RFC 3553, RFC 6924)",
	Nid:         "ietf",
	NssSchema: &urnfield.NssSchema{
		Description: "First element alternatives",
		ElementValidator: urnfield.ComplexOrNssElementValidatorFunc(
			[]*urnfield.NssSchema{
				//rfc: 1*DIGIT
				{
					Description: "rfc-nss rfc: 1*DIGIT",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"rfc",
						oneOrMoreDigitsNssSchema,
					),
				},

				//fyi: 1*DIGIT
				{
					Description: "fyi-nss fyi: 1*DIGIT",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"fyi",
						oneOrMoreDigitsNssSchema,
					),
				},

				//std: 1*DIGIT
				{
					Description: "std-nss std: 1*DIGIT",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"std",
						oneOrMoreDigitsNssSchema,
					),
				},

				//bcp: 1*DIGIT
				{
					Description: "bcp-nss bcp: 1*DIGIT",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"bcp",
						oneOrMoreDigitsNssSchema,
					),
				},

				// id: string — internet drafts (e.g. urn:ietf:id:draft-foo-bar-01)
				{
					Description: "draft-nss id: string",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"id",
						stringNssSchema,
					),
				},

				// mtg: string — IETF meeting documents (e.g. urn:ietf:mtg:ietf55)
				{
					Description: "mtg-nss mtg: string",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"mtg",
						stringNssSchema,
					),
				},

				//params: "The namespace is primarily opaque"
				{
					Description: "params-nss params: The namespace is primarily opaque",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"params",
						&urnfield.NssSchema{
							Description: `
							The IANA, as operator of the
      registry, may take suggestions for names to assign but they
      reserve the right to assign whatever name they desire, within
      guidelines set by the IESG.  The colon character (":") is used to
      denote a very limited concept of hierarchy.  If a colon is present
      then the items on both sides of it are valid names.  In general,
      if a name has a colon then the item on the left hand side
      represents a class of those items that would contain other items
      of that class.
							`,
							ElementValidator: urnfield.GlobNssElementValidatorFunc(
								glob.MustCompile("*"),
							),
						},
					),
				},

				//string (other)
				stringNssSchema,
			},
		),
	},
}
