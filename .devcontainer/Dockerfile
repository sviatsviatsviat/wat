FROM mcr.microsoft.com/devcontainers/go:2.1.5-1.26-bookworm

ENV PATH="/go/bin:${PATH}"
# Named volume inherits ownership from this path on first use.
RUN mkdir -p /workspaces/wat && chown -R vscode:vscode /workspaces
