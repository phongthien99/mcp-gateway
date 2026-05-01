---
title: Context
weight: 30
---

# MCP Context

Project and global context documents exposed by the MCP context tools.

## Source of truth

Context files live under:

```txt
context/global/{name}.md
context/{project_id}/{name}.md
```

MCP reads these files with:

```txt
get_context(project_id="{project_id}", name="{name}")
list_context(project_id="{project_id}")
```

## Writing rules

Use the shared [Markdown Rules]({{< relref "/docs/markdown-rules" >}}) for all context files.

Context documents should capture stable knowledge:

- Architecture and module boundaries.
- Coding standards.
- Compliance and security rules.
- Project-specific constraints.
- Decisions that future work must preserve.

Do not use context files for temporary notes, task logs, or generated workflow output. Put those in `artifacts/`.
