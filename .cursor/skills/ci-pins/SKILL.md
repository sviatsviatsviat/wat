---
name: ci-pins
description: >-
  Pin exact GitHub Actions versions in workflows and audit for floating tags.
  Use when editing .github/workflows, bumping actions, or reviewing CI config.
---

# CI action version pinning

Every `uses:` reference in `.github/workflows/` must pin an **exact release tag** with a patch version (e.g. `actions/checkout@v4.3.1`), not a floating major or minor (`@v4`, `@v6`).

Tool versions inside `with:` blocks must also be exact (e.g. `version: v2.12.2`, `go run ...@v1.6.0`).

## Audit command

Run from repo root after editing workflows:

```bash
rg 'uses:\s*[\w./-]+@v\d+$' .github/workflows
```

Exit code 0 with **no matches** means no floating major pins remain. Any match is a violation.

Also scan manually for `@main`, `@master`, and `@latest`.

## When bumping an action

1. Find the latest patch for the intended major line on the action's releases page.
2. Update **every** reference to that action in the same change (workflow, README, AGENTS.md if mentioned).
3. Re-run the audit command.
4. Unrelated action pins do not need to move unless compatibility requires it (see AGENTS.md Dependabot note).

## Examples

```yaml
# Good
uses: actions/checkout@v4.3.1
uses: actions/setup-go@v5.6.0
uses: golangci/golangci-lint-action@v9.3.0

# Bad — floating major
uses: actions/checkout@v4
uses: actions/setup-go@v5
```

## Definition of done (workflow edits)

- [ ] Audit command reports no floating `@vN` pins
- [ ] `go-version-file` and other existing `with:` config preserved unless intentionally changed
- [ ] README / AGENTS.md updated if they cite the bumped action or tool version
