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
// TODO: update to https://tools.ietf.org/html/rfc6924
var Schema = &urnfield.Schema{
	Description: "Schema for https://tools.ietf.org/html/rfc2648",
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

				//draft: string
				{
					Description: "draft-nss id: string",
					ElementValidator: urnfield.EqualsNssElementValidatorFunc(
						"id",
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
