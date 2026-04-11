---
name: fix-review
description: Works through a code review findings file issue by issue, following a structured protocol. Presents each finding, recommends a fix, applies it on instruction, updates tests and comments, updates the findings file, and offers a commit after each resolution.
argument-hint: path/to/code-review-findings.md
allowed-tools: Read Glob Edit Write Bash(git add *) Bash(git commit *) Bash(git log *) Bash(git status *) Bash(git diff *) Bash(go test *) Bash(go build *) Bash(go vet *) Bash(gofmt *)
---

## Findings File

Determine which findings file to work from:

1. If `$ARGUMENTS` is set, use that path as the findings file.
2. Otherwise, if a code review findings file is already in context (the user has it open or recently referenced it), use that.
3. Otherwise, use `glob` to find all `code-review-findings_*.md` files in the repository root, sort by filename descending, and use the most recent one.

Read the findings file before proceeding.

## Protocol

Work through the unresolved findings **in priority order** (top of the prioritised list first). Skip any finding already marked as resolved (~~strikethrough~~ or ✓).

For **each finding**, follow these steps exactly:

### Step 1 — Present

Present the finding to the user:
- Restate the description and why it matters
- Show the relevant code
- Think carefully about the best fix, considering correctness, idiom, and any knock-on effects

### Step 2 — Recommend

Propose a concrete fix, then offer exactly three options:

> **How would you like to proceed?**
> - **a)** Apply the fix
> - **b)** Think again (I will provide more input)
> - **c)** Add comments and move on (add appropriate `TODO` or `FIXME` to the relevant code)

Wait for the user's response before doing anything.

### Step 3 — Act

**If a):** Apply the fix:
- Make the code change
- Update or add tests to cover the fixed behaviour
- Update any affected comments or GoDoc
- Add inline comments where the reason for the fix is not obvious from the code alone
- Run `go test ./...` to confirm everything passes

**If b):** Incorporate the user's input, revise the recommendation, and re-present Step 2.

**If c):** Add a `TODO` or `FIXME` comment at the relevant location explaining the issue and what needs to be done.

### Step 4 — Update findings file

Mark the finding as resolved in the findings file:
- Strikethrough the heading: `## ~~Finding N — ...~~ ✓ Resolved`
- Add a **Resolution:** line summarising what was changed
- Update the prioritised table entry to match

### Step 5 — Offer commit

Ask the user if they would like a commit for the changes. If yes, create one following conventional commits style, including all changed files (code, tests, findings file). Always add:
```
Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

Then move on to the next unresolved finding.

## Notes

- Never apply a fix without the user choosing option a)
- Never skip the offer to commit after each resolved finding
- If all findings are resolved, say so clearly and stop
