---
title: Add MCP to Claude
weight: 15
---

# Add MCP to Claude

This workbench runs an MCP server over SSE.

## Configure the project

Copy the example env file and adjust values as needed:

```sh
cp .env.example .env
```

Key variables in `.env`:

| Variable | Default | Description |
|---|---|---|
| `MCP_PORT` | `8100` | Host port mapped to the MCP SSE endpoint |
| `API_PORT` | `8110` | Host port for the file API |
| `HUGO_PORT` | `1314` | Host port for the Hugo dashboard |
| `MCP_NAME` | `mcp-workbench` | Server name shown to LLM clients |
| `MCP_VERSION` | `1.0.0` | Server version shown to LLM clients |

Most defaults work out of the box. Only change a port if something on your machine already uses it.

## Start the gateway

Run the stack:

```sh
docker compose up
```

By default, the MCP endpoint is:

```txt
http://localhost:8100/sse
```

To scope all tool calls to a specific project, append the `?project=` query parameter:

```txt
http://localhost:8100/sse?project=<project-id>
```

Replace `<project-id>` with the directory name of your project under `projects/` (e.g. `react-mui`). When this parameter is set, the server automatically uses that project's context and artifacts without requiring `project_id` to be passed in every tool call.

The host port comes from:

```yaml
ports:
  - "${MCP_PORT:-8100}:8099"
```

Inside Docker, the MCP server listens on `8099`. From Claude on your host machine, use `8100` unless `MCP_PORT` is overridden.

## Add to Claude Code

Add the workbench as a local project MCP server:

```sh
claude mcp add --transport sse --scope local mcp-workbench "http://localhost:8100/sse?project=<project-id>"
```

Check that Claude sees it:

```sh
claude mcp list
claude mcp get mcp-workbench
```

Inside Claude Code, check connection status:

```txt
/mcp
```

## Project-shared config

To share the MCP config with the repo, use project scope:

```sh
claude mcp add --transport sse --scope project mcp-workbench "http://localhost:8100/sse?project=<project-id>"
```

This writes a `.mcp.json` file in the project root. Claude Code will ask for approval before using project-scoped MCP servers from that file.

Equivalent `.mcp.json`:

```json
{
  "mcpServers": {
    "mcp-workbench": {
      "type": "sse",
      "url": "http://localhost:8100/sse?project=<project-id>"
    }
  }
}
```

## Remove

Remove the server when needed:

```sh
claude mcp remove mcp-workbench
```

## Notes

Claude Code currently recommends HTTP transport for remote MCP servers when available. This gateway is configured for SSE, so use `--transport sse` unless the server is changed to expose an HTTP MCP endpoint.

Reference: [Claude Code MCP documentation](https://docs.anthropic.com/en/docs/claude-code/mcp).
