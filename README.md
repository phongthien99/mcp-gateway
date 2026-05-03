# MCP Workbench

An MCP (Model Context Protocol) server that gives AI assistants (Claude, Cursor, etc.) structured tools for managing project context, generating artifacts, and running multi-step development workflows.

## Overview

MCP Workbench acts as a persistent workspace layer between your AI assistant and your codebase. Instead of re-explaining your architecture every session, you store it once in a project context directory. The AI reads it via MCP tools, generates specs and task breakdowns as artifacts, and follows reusable workflow templates — turning ad-hoc prompts into a repeatable engineering process.

```
AI Assistant (Claude/Cursor)
        │  MCP protocol (SSE or stdio)
        ▼
  mcp-workbench server
  ├── Tools:     read/write/list artifacts, get/set context docs
  ├── Resources: artifact:// URIs for direct file access
  └── Prompts:   workflow-driven prompt templates
        │
        ▼
  data/
  ├── projects/{id}/context/   ← architecture, coding standards, compliance rules
  ├── projects/{id}/artifacts/ ← generated specs, plans, task lists
  ├── workflows/               ← YAML workflow definitions
  └── prompts/                 ← Markdown prompt templates
```

## Features

- **Context management** — store and retrieve project docs (architecture, standards, compliance) that persist across sessions
- **Artifact system** — AI writes structured outputs (specs, technical plans, tasks) to versioned files
- **Workflow engine** — chain prompt templates into multi-step flows (discovery → spec → plan → tasks)
- **SSE + stdio transport** — works as a local server or piped stdio process
- **Hugo documentation site** — bundled doc viewer served alongside the MCP server
- **File API** — REST endpoint for reading/writing markdown and YAML files from the browser

## Tech Stack

| Layer | Technology |
|-------|-----------|
| MCP server | Go 1.26, [mcp-go](https://github.com/mark3labs/mcp-go) |
| Dependency injection | uber/fx |
| Logging | uber/zap |
| Documentation | Hugo + hugo-book theme |
| Container | Docker, Docker Compose |

## Getting Started

### Prerequisites

- Docker and Docker Compose, **or**
- Go 1.21+ and Hugo extended

### Docker (recommended)

```bash
cp .env.example .env
docker compose up -d
```

| Service | URL |
|---------|-----|
| MCP server (SSE) | `http://localhost:8100/sse` |
| File API | `http://localhost:8110` |
| Documentation site | `http://localhost:1314` |

### Local development

```bash
# Start the MCP server with live reload
make mcp-run

# Start the Hugo documentation site
make hugo-book-run
```

Build a standalone binary:

```bash
make mcp-build
# Output: ./bin/mcp-workbench
```

## Configuration

Copy `.env.example` to `.env` and adjust as needed:

```env
MCP_TRANSPORT=sse          # sse | stdio
MCP_PORT=8099
MCP_NAME=mcp-workbench
MCP_VERSION=1.0.0

API_PORT=8110              # File API port
DASHBOARD_PORT=1313        # Hugo dev server
HUGO_BOOK_PORT=1314        # Standalone Hugo Book
FILEBROWSER_PORT=8080      # Markdown/YAML editor
```

## Connecting Your AI Assistant

### Claude Desktop / Claude Code

Add to your MCP config:

```json
{
  "mcpServers": {
    "mcp-workbench": {
      "url": "http://localhost:8100/sse"
    }
  }
}
```

For stdio transport:

```json
{
  "mcpServers": {
    "mcp-workbench": {
      "command": "./bin/mcp-workbench",
      "env": { "MCP_TRANSPORT": "stdio" }
    }
  }
}
```

## MCP Tools

| Tool | Description |
|------|-------------|
| `init_project` | Initialize a project context directory |
| `get_context` | Read a context document (architecture, standards, etc.) |
| `set_context` | Write or update a context document |
| `list_context` | List all context documents for a project |
| `read_artifact` | Read a generated artifact file |
| `write_artifact` | Write an artifact (spec, plan, task list) |
| `list_artifacts` | List artifacts for a project/feature |

## Workflows

Workflows are YAML files in `data/workflows/` that define ordered prompt steps. Each step references a Markdown prompt template and can read previous artifacts as input.

Built-in workflows:

| Workflow | Steps |
|----------|-------|
| `init-project` | Initialize project context from description |
| `dev-workflow` | Discovery → spec → technical plan → tasks |
| `gen-kanban` | Generate Kanban board from spec artifact |

See [docs/workflow-guide.md](docs/workflow-guide.md) for how to write custom workflows and prompt templates.

## Project Structure

```
mcp-gateway/
├── apps/
│   ├── mcp-server/          # Go MCP server
│   │   ├── main.go
│   │   └── src/
│   │       ├── module/      # App wiring (fx)
│   │       ├── tools/       # MCP tool handlers
│   │       ├── resources/   # MCP resource templates
│   │       ├── prompts/     # Workflow prompt registry
│   │       ├── workflow/    # Workflow runner
│   │       ├── api/         # File API server
│   │       └── config/      # Config loader
│   └── hugo-book-site/      # Hugo documentation
├── data/
│   ├── projects/            # Per-project context & artifacts
│   ├── workflows/           # YAML workflow definitions
│   ├── prompts/             # Markdown prompt templates
│   └── runs/                # Dry-run execution logs
├── docs/                    # Design docs and guides
├── docker-compose.yml
├── docker-compose.dev.yml
├── Dockerfile
├── Makefile
└── .env.example
```

## Documentation

- [Workflow Guide](docs/workflow-guide.md) — writing prompts and workflow YAML
- [User Stories](docs/user-stories-summary.md) — use cases and implementation status
- [Architecture Blog](docs/blog.md) — design decisions and rationale

## Makefile Targets

```bash
make up              # docker compose up (production)
make down            # docker compose down
make mcp-build       # build Go binary → ./bin/mcp-workbench
make mcp-run         # run with hot reload (dev)
make hugo-book-run   # Hugo dev server on :1314
make hugo-book-docker # run Hugo Book via Docker Compose
```

### Docker Image

```bash
# Build image (default: phongthien/mcp-workbench:latest)
make docker-build

# Build with a custom tag
make docker-build IMAGE=myrepo/mcp-workbench TAG=1.0.0

# Push to Docker Hub
make docker-push

# Build + push in one step
make docker-release

# Or use Docker directly
docker build -t phongthien/mcp-workbench:latest .
docker push phongthien/mcp-workbench:latest
```

## License

MIT
