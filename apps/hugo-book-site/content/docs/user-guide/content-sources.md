---
title: Content Sources
weight: 20
---

# Content Sources

The Hugo app can be used for authored documentation under `apps/hugo-book-site/content`.

The Docker Compose dashboard mount also keeps the existing repository data available to the server:

- `context/` for project context and standards.
- `artifacts/` for generated workflow outputs.
- `workflows/` for workflow definitions.
- `prompts/` for prompt templates.

## Ownership

Use [File Ownership]({{< relref "/docs/user-guide/file-ownership" >}}) to distinguish:

- System app files that implement the gateway and dashboard.
- Shared editable configuration that defines prompts, workflows, and global defaults.
- User workspace files that can change during normal project work.
