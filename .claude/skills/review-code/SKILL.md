---
name: review-code
description: Performs a code review. On main, reviews the full codebase. On a feature branch, reviews the diff against main. Accepts an optional PR number to review a specific PR. On main all the source code files, including the tests.
argument-hint: pr-number
allowed-tools: Read Grep Glob Bash(git branch *) Bash(git diff *) Bash(git ls-files *) Bash(gh pr diff *) Bash(gh pr review *) Bash(gofmt *) Bash(go vet *) Bash(staticcheck *)
---

Read [CONTRIBUTING.md](../../../CONTRIBUTING.md), then review the code defined in **Scope** below against every guideline it defines.

Look for any potential bugs or issues, but also look for places where the code could be improved to better follow the guidelines, best practices, performance or readability, even if it is not strictly in violation of them.

## Scope

Determine the scope using the following logic:

1. If `$ARGUMENTS` is set, run `gh pr diff $ARGUMENTS` and review the PR diff.
2. Otherwise, run `git branch --show-current` to get the current branch.
   - If the branch is not `main`, run `git diff main...HEAD` and review the branch diff.
   - If the branch is `main`, run `git ls-files` filtered to `*.go` files and review the full codebase.

## Static Analysis

Run the following tools and collect their output. Include any findings in the review, filtered to files within the scope where the scope is a diff or PR.

1. **Format:** `gofmt -l .` — list files with formatting issues.
2. **Vet:** `go vet ./...` — report suspicious constructs.
3. **Staticcheck:** `staticcheck ./...` — report style, correctness, and performance issues (includes hints). If `staticcheck` is not installed, skip this step and note it was unavailable.

## Instructions

- If the scope is a diff, review only the changed lines and their immediate context — including test files.
- If the scope is a file list, read each file before reviewing.
- Work through each scope item one by one, starting with the non-test files, prioritising those with fewest dependencies first.
- Do not flag issues outside the scope.

For each issue found, report:
- **Description** — what the issue is and why it matters
- **Classification** — classify finding as: Issue (bug or violation of guidelines, security problem or major performance issue), Improvement (not a violation but could be better), or Suggestion (minor improvement or best practice).
- **File and line** (where applicable). Where a file is applicable include the path relative to the repository root, and line numbers if the issue is localised.
- **Guideline** — cite the section from CONTRIBUTING.md
- **Recommendation** — what to change and why

End with a one-paragraph overall assessment and a prioritised list of findings (most important first). If there are no issues, say so clearly.

Offer to write the findings to a file `code-review-findings_<datetimestamp>.md` if the user would like. Number each finding sequentially starting from 1.
