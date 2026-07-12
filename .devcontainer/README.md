# Dev Container

Linux development environment for **wat** with Go 1.26, golangci-lint, ripgrep, and GitHub CLI (`gh`). The repo lives on a named Docker volume at `/workspaces/wat`. Only `.devcontainer/` is bind-mounted from the host so **Rebuild Container** picks up config edits from either side.

## First-time setup

1. Open this repo locally in Cursor (Windows path — used for `.devcontainer/` bind mount and rebuild).
2. Run **Dev Containers: Rebuild and Reopen in Container** from the Command Palette.
3. After the container starts, clone the repo inside it:

```bash
git clone https://github.com/sviatsviatsviat/wat.git /workspaces/wat
cd /workspaces/wat
go mod download
```

3. Open `/workspaces/wat` as the workspace folder if Cursor did not do so automatically.

## Troubleshooting

### "No remote authority found"

This means Cursor is still in a **local** window — the container connection did not complete. Close the window, reopen the repo locally, then run **Dev Containers: Rebuild and Reopen in Container** again. Do not use attach/reopen commands from a window that never connected.

Check **Output → Dev Containers** for build or server-install errors.

## Daily workflow

- **Update**: `git pull` from `/workspaces/wat`.
- **Verify**:

```bash
go vet ./... && go test ./... && golangci-lint run ./...
```

## Dev Container config

Edit `.devcontainer/` on the host or at `/workspaces/wat/.devcontainer/` in the container — they are the same files (host bind mount). Rebuild reads this folder; the rest of the repo stays on volume `wat-workspace`.

## Volume persistence

The repo persists in the Docker volume `wat-workspace` across container rebuilds. Removing it is destructive:

```bash
docker volume rm wat-workspace
```
