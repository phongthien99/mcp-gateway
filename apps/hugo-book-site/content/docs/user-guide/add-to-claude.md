---
title: Add MCP to Claude
weight: 15
---

# Add MCP to Claude

This workbench runs an MCP server over SSE.

## Start the gateway

Run the stack:

```sh
docker compose up
```

By default, the MCP endpoint is:

```txt
http://localhost:8100/sse
```

The host port comes from:

```yaml
ports:
  - "${MCP_PORT:-8100}:8099"
```

Inside Docker, the MCP server listens on `8099`. From Claude on your host machine, use `8100` unless `MCP_PORT` is overridden.

## Add to Claude Code

Add the workbench as a local project MCP server:

```sh
claude mcp add --transport sse --scope local mcp-workbench http://localhost:8100/sse
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
claude mcp add --transport sse --scope project mcp-workbench http://localhost:8100/sse
```

This writes a `.mcp.json` file in the project root. Claude Code will ask for approval before using project-scoped MCP servers from that file.

Equivalent `.mcp.json`:

```json
{
  "mcpServers": {
    "mcp-workbench": {
      "type": "sse",
      "url": "http://localhost:8100/sse"
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
