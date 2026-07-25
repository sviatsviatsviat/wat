# Architecture

`wat` separates native hook protocols, portable policy, process dispatch, and
CLI orchestration. The central design rule is that protocol-specific behavior
stays in its agent SDK while the portable SDK depends on those SDKs and adapts
only their common behavior.

## Dependency direction

```text
cmd/wat
  -> sdk/claude, sdk/copilot, sdk/cursor
  -> sdk/run

generated .wat hook project
  -> sdk/agnostic and/or native SDKs
  -> sdk/run

sdk/agnostic
  -> sdk/claude, sdk/copilot, sdk/cursor
  -> sdk/run

sdk/claude, sdk/copilot, sdk/cursor
  -> sdk/run
  -> internal/hookkit

sdk/run
  -> internal/hookkit
```

The module has no third-party runtime dependencies. The per-agent SDKs do not
import `sdk/agnostic`; reversing that edge would create both an import cycle and
an ownership problem.

## Package responsibilities

### `cmd/wat`

The command root owns argument parsing, help text, and process exit codes.
Substantive behavior lives in focused internal packages:

| Package | Responsibility |
|---|---|
| `internal/project` | Upward `.wat/` discovery and `WAT_PROJECT_DIR` |
| `internal/initproj` | Scaffold and `go mod tidy` |
| `internal/buildcache` | Content-addressed bootstrap builds |
| `internal/hookmanifest` | Load `run.Inspect` data from the built hook |
| `internal/installcfg` | Reconcile wat-managed native config entries |
| `internal/hostconfig/*` | Native configuration schemas |
| `internal/hookrun` | Live stdin/stdout/stderr execution |
| `internal/hooktest` | Fixture validation, execution, and reporting |
| `internal/doctor` | Independent diagnostic checks and fixes |
| `internal/dialect`, `paths`, `hookconfig` | Small shared CLI concepts |

Command files in `cmd/wat` should remain thin. New behavior belongs in an
internal package with injected filesystem, environment, process, and I/O
dependencies so it can be tested without mutating the developer's machine.

### `sdk/run`

`run` is the public process boundary. SDK registrars implement `run.Hooks` and
contribute native handlers to a sealed registry.

`Serve` performs one dispatch cycle:

```text
contribute registrations
  -> read stdin
  -> detect dialect
  -> peek hook_event_name
  -> find handlers
  -> decode once
  -> invoke and merge in order
  -> encode once
  -> exit
```

`Inspect` follows the same contribution path without invoking handlers and
returns the native dialect/event manifest used by install, test, and doctor.
Registration metadata must therefore never be maintained in a separate CLI
table.

### Per-agent SDKs

`sdk/claude`, `sdk/copilot`, and `sdk/cursor` own their complete native wire
contracts:

- event-name and dialect constants;
- typed event payloads and embedded native envelopes;
- typed result builders and sealed outputs;
- output merge, stop, encode, and exit behavior;
- payload detection and codec registration;
- a fluent `UseHooks` registrar.

Public package roots are facades. Most public types are aliases to vertical
event packages under:

```text
sdk/<agent>/internal/hooks/<domain>/<event>/
```

A result-producing event normally has:

| File | Role |
|---|---|
| `event.go` | Native input payload and `EventName` |
| `results.go` | Hook-scoped builder interface and implementation |
| `output.go` | Sealed output plus fluent `With*` methods |
| `bind.go` | Typed registration into the shared dialect handler bag |
| `event_test.go` | Decode, builder, encode, merge, and edge behavior |

Observe-only events omit result/output code. Shared native concepts belong in
that SDK's `internal/event`, `internal/runtime`, or `internal/tools`, then are
aliased deliberately from the package root.

Not every decoded event must be registerable. The public `UseHooks` methods are
the supported handler surface; exported decode-only event types can still be
used by internal adapters or future registrations.

### `sdk/agnostic`

The portable SDK is an adapter layer, not a fourth wire protocol.

`sdk/agnostic/internal/model` defines the normalized leaf events and portable
result interfaces. Per-agent adapter packages:

```text
sdk/agnostic/internal/{claude,copilot,cursor}
```

map native input to a normalized event, wrap native result builders behind the
portable interface, and return native outputs. A root `On*` method appends all
supported native registrations to one deferred registrar.

Canonical tool names and typed portable inputs live in
`sdk/agnostic/tools`. Native names are always preserved on `ToolCall.Native`;
normalization must not discard protocol information.

Portable APIs expose only behavior that all three agents can represent
honestly. Agent-specific fields such as Claude session environment or Cursor
prompt blocking remain in the native SDKs.

### `internal/hookkit`

`hookkit` contains module-private protocol machinery shared by SDKs:

- codec and decoder registration;
- the internal `Event`, `Output`, and handler abstractions;
- handler queues and dialect bags;
- JSON, environment, shell, timestamp, merge, and tool-name helpers.

It is not an author-facing SDK. Public packages should expose typed concepts,
not leak general-purpose hookkit helpers merely for convenience.

## Core implementation patterns

### Deferred fluent registration

`UseHooks` does not perform runtime work. It queues typed registration
functions. `Contribute` installs them into a registry later. This allows the
same authored values to drive both `run.Serve` and `run.Inspect`.

Multiple registrars for one dialect merge. Handler order is the order in which
`Hooks` values contribute, followed by method registration order inside each
value.

### Decode once, encode once

The native codec peeks the event name before decoding. If no handler is
registered, the process exits successfully without decoding that event. When
handlers exist, one typed event is shared across them and their outputs are
folded before a single encode.

This keeps wire policy in one place and avoids invalid combinations produced by
concatenating independent JSON responses.

### Hook-scoped result builders

Handlers do not construct outputs with struct literals. Each result-producing
registration supplies a builder whose methods match that event's legal
responses. Fluent `With*` methods add advanced fields to a builder-produced
value.

Return `nil` for no opinion. Builders may also expose `Noop` on native APIs,
but portable handlers use `nil` because it maps consistently.

### Preserve native information at normalization boundaries

Normalized events carry an `Envelope` with the dialect and native event name.
Tool calls carry both canonical and native names plus raw typed input. When a
native protocol has no portable representation, skip the portable mapping or
document the deliberate downgrade; do not fabricate parity.

## How to extend the code

### Add or change a native event

1. Confirm the native payload and output contract.
2. Implement or update the event vertical slice under the owning SDK.
3. Register its decoder in the SDK codec.
4. Alias the intended public types/constants from the package root.
5. Add a typed fluent method to `hooks.go` when the event is author-facing.
6. Test input decoding, result building, output JSON/exit code, merge behavior,
   and nil handler registration.
7. Update godoc, the native API table, protocol matrix, and changelog.

### Add or change a portable event

1. Establish behavior genuinely supported by all three agents.
2. Define the normalized event and narrow result interfaces in
   `internal/model`.
3. Implement one adapter per agent, including input mapping and result wrapping.
4. Fan those adapters out from the root `On*` method.
5. Test the native registration expansion with `run.Inspect` and mapping/output
   behavior for every agent.
6. Update the SDK guide, protocol matrix, godoc, and changelog.

### Add or change a CLI command

1. Keep flag parsing and help in `cmd/wat`.
2. Put behavior in a focused `cmd/wat/internal` package with explicit
   dependencies.
3. Preserve stream conventions: help to stdout, invalid usage to stderr.
4. Add unit tests for injected behavior and end-to-end CLI tests for the public
   contract.
5. Update usage documentation, help text, and changelog together.

## Architectural invariants

- Native SDKs never import `sdk/agnostic`.
- Native wire encoding stays inside the owning SDK.
- CLI installation derives events from `run.Inspect`; it does not guess.
- `hooks.go` registration has no external side effects.
- Exported identifiers have godoc and author-facing examples where useful.
- Merge helpers do not mutate caller-owned maps or slices.
- Public behavior changes include tests and user documentation in the same
  change.
