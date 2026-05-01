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
