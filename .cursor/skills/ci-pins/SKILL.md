---
name: ci-pins
description: >-
  Keep GitHub Actions, Go tools, and Dev Container features exactly pinned and
  synchronized across the repository.
---

# CI and tool pins

All workflow `uses:` values must use exact release tags with patch versions.
Tool versions in workflow commands, Dev Container features, and installation
documentation must also be exact. Do not use `@latest`, branches, or floating
major/minor tags.

When bumping one dependency:

1. locate every reference with `rg`;
2. update all references to that dependency in the same change;
3. preserve unrelated pins unless compatibility requires another bump;
4. run the workflow-tag audit;
5. run the normal repository verification when config affects builds.

Examples of synchronized references:

- golangci-lint: workflow, `.devcontainer/devcontainer.json`,
  `CONTRIBUTING.md`, and `AGENTS.md`;
- govulncheck: every workflow command using `govulncheck@...`;
- an Action: every `uses:` entry for that Action.

Audit floating major tags:

```bash
rg 'uses:\s*[\w./-]+@v\d+$' .github/workflows
```

The audit must produce no matches. Also search manually for `@latest`,
`@main`, and `@master`.
