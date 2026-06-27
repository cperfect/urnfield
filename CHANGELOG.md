# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2026-06-27
### ⚠ Breaking Changes

- Migrate module path to v2 — the module import path is now github.com/cperfect/urnfield/v2; consumers must update their imports. The package name is unchanged (urnfield).
### Bug Fixes

- **security, CI**: Potential fix for code scanning alert no. 2: Workflow does not contain permissions
- Conform to spec 1.0.0 equivalence and glob raw-tail rules
### Continuous Integration

- Check out spec submodule for conformance vectors
- Check out spec submodule in prepare-release workflow
### Documentation

- Document spec-first workflow and conformance suite
### Features

- **spec**: Add spec v1.0.0 as submodule
### Maintenance

- **devcontainer**: Remove second claude feature
- **deps**: Bump actions/checkout from 6 to 7 in the all-actions group
- **ai**: Grants Claude Go development task permissions
- Drop manual CHANGELOG entry
### Tests

- Add spec conformance runner (red)
### Build

- **cliff**: Improve changelog rendering
- **[BREAKING]** Migrate module path to v2
  the module import path is now github.com/cperfect/urnfield/v2; consumers must update their imports. The package name is unchanged (urnfield).
- **cliff**: Collect breaking changes into a top section
## [1.0.0] - 2026-04-22
### Bug Fixes

- **fmt**: Configures Claude for gofmt; cleans regex patterns
- **review issues**: Simplifies nil checks and string splitting
- **schema**: Guard against empty nss slice in validate and all validators
- **parse**: Anchor Pattern regex to reject embedded URNs
- **schema**: Document and finalise GlobNssElementValidatorFunc semantics
- **schema**: Use MatchString instead of Match([]byte(...))
- **urnfield**: Remove dead hasVal branch in writeKeyValuesMap
- Lowercase error strings per Go style guide
- **schema**: Clarify ComplexOrNssElementValidatorFunc success path
- **fmt**: Aligns comments in schema test
### Documentation

- **urnfield**: Fix GoDoc comments to follow function-name-first convention
- **ietf**: Fix misleading //draft comment on id sub-namespace validator
- **schema**: Resolve design TODOs with documentation
- **claude**: Remove stale api.go reference
### Features

- **AI**: Adds Claude AI code review capabilities
- **CI**: Adds Go code validation workflow
- **AI**: Enhances code review with Go static analysis
- **CI**: Adds `go vet` to CI workflow
- **AI**: Add fix-review skill for structured code review resolution
- **ietf**: Add missing mtg sub-namespace and update to RFC 6924
- **security**: Adds Dependabot configuration
- **CI**: Introduces automated release workflow
### Maintenance

- Add devcontainer and update to go 1.26
- Update deps
- **CONTRIBUTING**: Adds contributing guidelines
- **README**: Update
- Rename main files to match package name
- **CONTRIBUTING**: Adds guidelines for AI tool usage in contributions
- **AI**: Init claude etc.
- Removes api.go placeholder
- **format**: And GoDoc
- Introduces a thoughts file and AI directive
- **CONTRIBUTING**: Enhances developer contributing guide
- **README**: Adds comprehensive usage examples to README
- **README**: Removes WIP warning from README
- **comments**: Adds comprehensive package-level documentation
- **CI**: Enables static code analysis in CI
- **structure**: Exclude examples from lib module
- Removes development notes file
- **CI**: Configures Dependabot to group dependencies
- **docs**: Adds release process documentation
- **deps**: Bump the all-actions group with 2 updates
- **release**: Prepare v1.0.0
### Tests

- **ietf**: Exercise testUrnSchemas via TestSchemas
- **urnfield**: Add missing Parse failure cases and remove TODO comments

