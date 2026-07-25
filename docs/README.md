# Documentation

The documentation is organized by audience:

| Document | Audience and scope |
|---|---|
| [Using wat](usage.md) | Hook authors using the CLI and operating a `.wat/` project |
| [SDK API](sdk.md) | Go developers writing portable or native hook handlers |
| [Architecture](architecture.md) | Maintainers extending packages, events, adapters, or commands |
| [Agent protocols](agent-formats.md) | Maintainers working on codecs, mappings, and tool normalization |
| [Contributing](../CONTRIBUTING.md) | Anyone preparing a change or pull request |
| [Changelog](../CHANGELOG.md) | Users evaluating the next release's visible behavior |

The root [README](../README.md) is the short entry point. Package-level details
also belong in godoc so they are available to library consumers:

```bash
go doc github.com/sviatsviatsviat/wat/sdk/agnostic
go doc github.com/sviatsviatsviat/wat/sdk/claude
go doc github.com/sviatsviatsviat/wat/sdk/copilot
go doc github.com/sviatsviatsviat/wat/sdk/cursor
go doc github.com/sviatsviatsviat/wat/sdk/run
```

## Documentation ownership

Update documentation in the same change as behavior:

| Change | Required documentation |
|---|---|
| CLI command, flag, discovery, cache, or exit behavior | `usage.md`, command help, tests, and changelog |
| Portable SDK API or mapping | `sdk.md`, godoc, protocol matrix when relevant, tests, and changelog |
| Native event or result API | `sdk.md`, godoc, protocol matrix, tests, and changelog |
| Package boundary or implementation pattern | `architecture.md` and contributor/agent instructions when necessary |
| Internal refactor with no user-visible effect | Architecture docs only if the documented design changed; no changelog entry |

Documentation should describe implemented behavior. Local design notes may live
under the gitignored `plan/` directory, but they are not a substitute for these
committed references.
