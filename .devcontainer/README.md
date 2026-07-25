# Dev Container

The Dev Container is the authoritative Linux environment for wat's complete
verification suite. It includes:

- Go 1.26 from `mcr.microsoft.com/devcontainers/go:2-1.26-bookworm`;
- golangci-lint 2.12.2;
- ripgrep 15.1.0;
- GitHub CLI 2.96.0;
- Node.js 22.23.1 and the Claude Code CLI (feature `1.0.5`, latest CLI at
  rebuild time; auto-update disabled via `DISABLE_AUTOUPDATER`).

Claude Code auth and settings persist in the named volume `wat-claude-config`
mounted at `/home/vscode/.claude`. The Dockerfile seeds that path as
`vscode`-owned so the volume is writable on first create; `postStartCommand`
also `chown`s it in case an older root-owned volume already exists.

If transcript or settings writes still fail with `EACCES`, fix ownership or
recreate the volume:

```bash
sudo chown -R vscode:vscode /home/vscode/.claude
# or, from the host after stopping the container:
docker volume rm wat-claude-config
```

The repository lives in the named Docker volume `wat-workspace` at
`/workspaces/wat`. Only the local `.devcontainer/` directory is bind-mounted,
so configuration edits are visible during a rebuild without putting the Go
workspace on a Windows bind mount.

## First-time setup

1. Open the local checkout in an editor with Dev Containers support.
2. Run **Dev Containers: Rebuild and Reopen in Container**.
3. Clone the repository into the named volume:

```bash
git clone https://github.com/sviatsviatsviat/wat.git /workspaces/wat
cd /workspaces/wat
go mod download
```

4. Open `/workspaces/wat` as the workspace folder if it was not selected
   automatically.

The `postStartCommand` runs `go mod download` on later starts when `go.mod` is
present.

## Daily workflow

Run commands from `/workspaces/wat`:

```bash
git pull
go vet ./...
go test ./...
golangci-lint run ./...
go build ./cmd/wat
```

The Linux environment matters because some CLI tests execute Unix utilities
such as `echo` and `sleep`.

## Updating the environment

Edit `.devcontainer/Dockerfile` or `.devcontainer/devcontainer.json` in the
host checkout, then rebuild the container. The bind mount makes the same files
visible at `/workspaces/wat/.devcontainer/`.

Keep image, Feature, and tool versions exact. When bumping a tool, update every
reference to that same version in workflows and contributor/agent
documentation.

## Troubleshooting

### No remote authority found

The editor is still in a local window or the container connection did not
complete. Reopen the local checkout and run **Rebuild and Reopen in Container**
again. Check the editor's Dev Containers output for build or server-install
errors.

### Workspace is empty after rebuild

The named volume may not contain a clone yet. Follow the first-time clone
steps. Container rebuilds preserve `wat-workspace`; they do not populate it
from the host checkout.

## Volume lifecycle

The repository persists across container rebuilds. Removing the Docker volume
deletes that clone and all unpushed work inside it:

```bash
docker volume rm wat-workspace
```

Only remove the volume when that destructive reset is intentional.
