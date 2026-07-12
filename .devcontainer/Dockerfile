FROM mcr.microsoft.com/devcontainers/go:1.26-bookworm

# Pin golangci-lint to match CI and AGENTS.md.
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

# ripgrep for ci-pins audit (rg 'uses:...').
RUN apt-get update && apt-get install -y --no-install-recommends ripgrep \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/go/bin:${PATH}"

# Named volume inherits ownership from this path on first use.
RUN mkdir -p /workspaces/wat && chown -R vscode:vscode /workspaces
