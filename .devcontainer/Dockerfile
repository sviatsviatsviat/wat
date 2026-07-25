FROM mcr.microsoft.com/devcontainers/go:2-1.26-bookworm

ENV PATH="/go/bin:${PATH}"
# Named volumes inherit ownership from these paths on first use.
RUN mkdir -p /workspaces/wat /home/vscode/.claude \
	&& chown -R vscode:vscode /workspaces /home/vscode/.claude
