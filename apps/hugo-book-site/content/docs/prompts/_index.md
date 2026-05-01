---
title: Prompts
weight: 20
---

# Prompts

Prompt templates exposed by the MCP prompt registry.

## Source of truth

Prompt files live under:

```txt
prompts/{name}.md
```

MCP loads these files through workflow prompt definitions and exposes them as reusable prompt steps.

## Writing rules

Use the shared [Markdown Rules]({{< relref "/docs/markdown-rules" >}}) for all prompt templates.

Prompt files should be stable templates:

- Define the role and task clearly.
- State required inputs and expected output.
- Keep formatting rules explicit.
- Avoid project-specific one-off details unless the prompt is intentionally scoped.

Put generated results in `artifacts/`, and put long-term project facts in `context/`.
