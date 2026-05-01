---
title: File Ownership
weight: 25
---

# File Ownership

Repository content is split into three groups:

- **System app files**: code and dashboard files that implement the MCP workbench.
- **Shared editable configuration**: prompts, workflows, and global context that users can change, but changes affect every project or run that uses them.
- **User workspace files**: project-specific context and generated outputs that users and agents can change during normal work.

## System app files

These files implement the MCP workbench and dashboard. Change them when developing the app itself.

| Path | Purpose |
| --- | --- |
| `apps/hugo-book-site/content/` | Authored dashboard documentation. |
| `apps/mcp-server/` | MCP server implementation. |
| `docker-compose.yml` | Local service wiring and mounted content sources. |

## Shared editable configuration

These files are user-editable shared defaults. Change them when you want to update the baseline behavior for multiple projects or runs.

| Path | Purpose |
| --- | --- |
| `prompts/*.md` | Reusable prompt templates used by workflow steps. |
| `workflows/*.yaml` | Workflow definitions and step wiring. |
| `context/global/*.md` | Default global standards and constraints. |

## User workspace files

These files are expected to change per project, feature, or run.

| Path | Purpose |
| --- | --- |
| `context/{project_id}/*.md` | Project-specific architecture, standards, compliance rules, and local notes. |
| `artifacts/{project}/{feature}/` | Generated discovery notes, specs, plans, summaries, and verification output. |
| `runs/*.yaml` | Run-specific workflow configuration. |

## Editing rule

For normal project-specific work, update `context/{project_id}/`, `artifacts/`, and `runs/`.

Update `prompts/`, `workflows/`, and `context/global/` when changing shared defaults.

Update `apps/` or Docker files when changing the gateway implementation or dashboard.
