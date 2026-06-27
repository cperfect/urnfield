# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`urnfield` is a Go library for parsing, validating, and formatting URNs (Uniform Resource Names) per RFC 8141. It supports pluggable schema validators for namespace-specific NSS validation.

It is the **Go reference implementation** of the technology-agnostic [`urnfield` specification](https://github.com/cperfect/urnfield-spec), vendored here as a git submodule under `submodules/urnfield-spec` (pinned to a spec version tag). The spec and its conformance suite are the source of truth for behaviour — follow the spec-first workflow (see [CONTRIBUTING.md](./CONTRIBUTING.md) "Spec-First Development"): behaviour changes go through the spec repo and its conformance vectors first; non-behavioural changes (refactors, docs, perf, examples) can be made locally.

## Commands

```bash
# Initialise the spec submodule (needed for the conformance tests)
git submodule update --init --recursive

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test package
go test ./examples/ietf

# Run a single test
go test -run TestFunctionName ./...

# Tidy dependencies
go mod tidy
```

## Contributing

Follow all guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md), except any sections that explicitly exclude agents.

## Architecture

The library has two core layers, plus a conformance layer:

### 1. URN Parsing (`urnfield.go`)
Parses URN strings into `Urn` structs using regex, then formats them back to strings. The `Urn` struct holds:
- `Nid` — Namespace ID
- `Nss []string` — Namespace-Specific String, split on `:` or `/` delimiters (tracked via `NssSlashDelimiter` flag)
- `Query map[string][]string` — `?=` component
- `Resolvers map[string][]string` — `?+` component
- `Fragment string` — `#` component

Key functions: `Parse()` (string → Urn), `Format()` / `ToString()` (Urn → string), and `Urn.Equivalent()` (spec §8 identity: NID case-insensitive, NSS case-sensitive and element-wise; delimiter/query/resolvers/fragment ignored). Query and Resolvers keys are sorted on output for deterministic testing.

### 2. Schema Validation (`schema.go`)
Schemas define validation rules for a specific URN namespace. Validation is a recursive chain: each `NssElementValidator` processes one NSS element and returns the remainder plus the next validator. This allows composing complex hierarchical validations.

Built-in validator factories:
- `RegexNssElementValidatorFunc` — regex match
- `EqualsNssElementValidatorFunc` — exact match
- `SimpleOrNssElementValidatorFunc` — one of a list of strings
- `ComplexOrNssElementValidatorFunc` — one of a list of validators
- `GlobNssElementValidatorFunc` — glob pattern match (uses `github.com/gobwas/glob`); matches against the **raw NSS tail** rejoined with the URN's active delimiter (`:` or `/`), per spec §10.

### 3. Example (`examples/ietf/`)
A real-world implementation of the IETF URN namespace (RFC 2648). Read this to understand how to compose the validator factories for hierarchical, branching NSS structures. It is a separate Go module that `replace`s the library with the local source.

### 4. Conformance (`conformance_test.go`)
An external-package (`urnfield_test`) runner that loads the spec's language-neutral YAML vectors from `submodules/urnfield-spec/conformance` (`parse`, `format`, `validate`, `equals`) and drives the public API against them. The `validate` harness translates the spec's declarative matcher model (SPECIFICATION.md §9) into `NssSchema` trees via the public validator factories. Passing this suite is what "conformant with spec version X.Y.Z" means.

## Key Patterns

- **Validator chain**: `NssElementValidator` functions receive remaining NSS elements plus the URN's active delimiter and return the next validator to use — this is the core extensibility mechanism. The delimiter is used only by glob matchers.
- The public API is the exported symbols in `urnfield.go` (parsing/formatting) and `schema.go` (validation). There is no separate `api.go`.
- Tests use data-driven structs comparing expected vs. actual `Urn` values.

## Directives
- Ignore THOUGHTS.md unless explicitly told otherwise