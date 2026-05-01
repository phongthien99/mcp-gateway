---
title: Getting Started
weight: 10
---

# Getting Started

Run the full stack:

```sh
docker compose up
```

Run the MCP server with the embedded Hugo dashboard:

```sh
make mcp-run
```

By default, the MCP SSE endpoint is available on port `8099` and the dashboard is available on port `1313`.

When running through Docker Compose, Claude should connect to the host MCP endpoint:

```txt
http://localhost:8100/sse
```

See [Add MCP to Claude]({{< relref "/docs/add-to-claude" >}}) for the Claude Code setup command.

The file API used by the in-page editor is available at:

```txt
http://localhost:8110
```

## File Browser

Docker Compose also starts File Browser for editing workspace Markdown/YAML files:

```txt
http://localhost:8080
```

It is bound to `127.0.0.1` and runs without login for local editing.

Mounted editable folders:

- `prompts/`
- `workflows/`
- `context/`
- `artifacts/`
- `runs/`
- `docs/`
- `apps/hugo-book-site/content/` as `hugo-content/`

Use File Browser for workspace docs and shared configuration. Use the IDE for application code changes under `apps/mcp-server/`.
