# Markdown output rules

When writing a documentation artifact, use a `.md` path.

Required markdown structure:

```md
---
title: Short Human Title
weight: 10
---

# Short Human Title

## Purpose

State why this document exists.

## Scope

State what this document covers and does not cover.
```

Rules:

- Use exactly one `# H1`.
- Use `##` for main sections.
- Do not skip heading levels.
- Use kebab-case file names, for example `technical-plan.md`.
- When a folder contains multiple markdown documents, add `_index.md` for that folder.
- `_index.md` should include front matter, one H1, Purpose, and a Contents list.
- Give sibling markdown files `weight` values so navigation order is predictable.
- Write decisions as `Decision: ...`.
- Write assumptions as `Assumption: ...`.
- Write unresolved items as `Open question: ...`.
- Wrap file paths, commands, tool names, and IDs in backticks.
- Use fenced code blocks with a language tag for commands, code, YAML, or JSON.
- If the output is not documentation, choose the correct extension instead of `.md`, for example `.yaml`, `.json`, `.txt`, or an image extension.
