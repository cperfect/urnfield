# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`urnfield` is a Go library for parsing, validating, and formatting URNs (Uniform Resource Names) per RFC 8141. It supports pluggable schema validators for namespace-specific NSS validation.

## Commands

```bash
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

The library has two core layers:

### 1. URN Parsing (`urnfield.go`)
Parses URN strings into `Urn` structs using regex, then formats them back to strings. The `Urn` struct holds:
- `Nid` — Namespace ID
- `Nss []string` — Namespace-Specific String, split on `:` or `/` delimiters (tracked via `NssSlashDelimiter` flag)
- `Query map[string][]string` — `?=` component
- `Resolvers map[string][]string` — `?+` component
- `Fragment string` — `#` component

Key functions: `Parse()` (string → Urn), `Format()` / `ToString()` (Urn → string). Query and Resolvers keys are sorted on output for deterministic testing.

### 2. Schema Validation (`schema.go`)
Schemas define validation rules for a specific URN namespace. Validation is a recursive chain: each `NssElementValidator` processes one NSS element and returns the remainder plus the next validator. This allows composing complex hierarchical validations.

Built-in validator factories:
- `RegexNssElementValidatorFunc` — regex match
- `EqualsNssElementValidatorFunc` — exact match
- `SimpleOrNssElementValidatorFunc` — one of a list of strings
- `ComplexOrNssElementValidatorFunc` — one of a list of validators
- `GlobNssElementValidatorFunc` — glob pattern match (uses `github.com/gobwas/glob`)

### 3. Example (`examples/ietf/`)
A real-world implementation of the IETF URN namespace (RFC 2648). Read this to understand how to compose the validator factories for hierarchical, branching NSS structures.

## Key Patterns

- **Validator chain**: `NssElementValidator` functions receive remaining NSS elements and return the next validator to use — this is the core extensibility mechanism.
- The public API is the exported symbols in `urnfield.go` (parsing/formatting) and `schema.go` (validation). There is no separate `api.go`.
- Tests use data-driven structs comparing expected vs. actual `Urn` values.

## Directives
- Ignore THOUGHTS.md unless explicitly told otherwise