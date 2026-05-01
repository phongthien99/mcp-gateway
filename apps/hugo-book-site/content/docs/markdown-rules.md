---
title: Markdown Rules
weight: 25
---

# Markdown Rules

Use these rules for documentation files stored under `context/` and `artifacts/`. These files are read by both Hugo and MCP tools, so they should be easy for humans to scan and easy for LLMs to parse.

Agents should choose the file type that matches the artifact. Use `.md` for documentation by default, but use `.yaml`, `.json`, `.txt`, images, or other formats when the output is not a markdown document.

## File naming

- Use lowercase kebab-case file names: `technical-plan.md`, `coding-standards.md`.
- Keep names stable because MCP references files by path or logical name.
- Use `.md` for authored documentation.
- Use the correct extension for non-document artifacts, for example `workflow.yaml`, `schema.json`, or `diagram.png`.
- Avoid spaces, Vietnamese accents, and special characters in file names.

## Folder rules

When a folder contains multiple markdown files, add an `_index.md` file for that folder.

Example:

```txt
artifacts/project/nay/
  _index.md
  discovery.md
  spec.md
  technical-plan.md
  tasks.md
```

The `_index.md` file should explain what the folder contains:

```md
---
title: Nay
weight: 10
bookCollapseSection: false
---

# Nay

## Purpose

Collect workflow artifacts for the `nay` feature.

## Contents

- `discovery.md`: Requirement discovery.
- `spec.md`: Feature specification.
- `technical-plan.md`: Implementation plan.
- `tasks.md`: Task breakdown.
```

Folder rules:

- Use lowercase kebab-case folder names.
- Add `_index.md` when the folder has child docs that should appear as a section in Hugo.
- Give each child markdown file a `weight` so Hugo navigation order is predictable.
- Keep nesting shallow. Prefer `project/feature/file.md` over deeply nested folders.
- Put generated support files next to the document that references them.
- If a folder mixes docs and non-doc artifacts, describe those non-doc files in `_index.md`.
- Do not rely on folder names alone to explain meaning; write a short Purpose section.

## Required structure

Every markdown file should include:

```md
---
title: Short Human Title
weight: 10
---

# Short Human Title

## Purpose

What this document is for.

## Scope

What this document covers and does not cover.

## Content

The main body.

## References

- Link or path to related files.
```

For MCP-generated artifact files that are not meant to appear directly as Hugo pages, front matter is optional. If a markdown file is mounted into Hugo content, add front matter.

## Front matter

Use TOML-style YAML front matter at the top of Hugo-visible files:

```yaml
---
title: Technical Plan
weight: 20
---
```

Recommended fields:

- `title`: Display name in Hugo navigation.
- `weight`: Sort order inside a section.
- `bookCollapseSection`: Use on `_index.md` files when a section has many children.

Keep front matter small. Do not put operational data there if MCP tools need to read it.

## Heading rules

- Use exactly one `# H1` per file.
- Use `##` for main sections.
- Use `###` only when a section is long enough to need grouping.
- Do not skip levels, for example do not jump from `##` to `####`.
- Keep headings descriptive: `API Contracts`, not `Notes`.

## Writing rules

- Start each document with the purpose and scope.
- Prefer short paragraphs and bullet lists.
- Put decisions in explicit language: `Decision: ...`.
- Put assumptions in explicit language: `Assumption: ...`.
- Put unresolved work in explicit language: `Open question: ...`.
- Include exact file paths, tool names, command names, and IDs in backticks.
- Avoid vague words like "some", "maybe", "stuff", and "etc." when the missing detail matters.

## Code and commands

Use fenced code blocks with a language tag:

````md
```sh
docker compose up hugo-book-site
```

```go
func main() {}
```
````

Use inline code for short identifiers like `project_id`, `context/global/architecture.md`, and `resource://artifact/project/nay/technical-plan`.

## Links and paths

- Link related docs when they are stable.
- Use repository-relative paths for local files.
- For MCP resource references, write the full URI.

Example:

```md
Related artifact: `resource://artifact/project/nay/technical-plan`
Context file: `context/global/architecture.md`
```

## Tables

Use tables only for compact comparison data. Keep them small enough to read on mobile.

```md
| Field | Meaning |
| --- | --- |
| `project_id` | Directory under `context/` |
| `name` | Markdown file name without `.md` |
```

## Context documents

Context files should describe stable rules and project knowledge. They should not be daily logs.

Use this structure:

```md
---
title: Architecture
weight: 10
---

# Architecture

## Purpose

## Current System

## Constraints

## Decisions

## References
```

Expected context paths:

```txt
context/global/architecture.md
context/global/coding-standards.md
context/global/compliance-rules.md
context/{project_id}/architecture.md
context/{project_id}/coding-standards.md
context/{project_id}/compliance-rules.md
```

## Artifact documents

Markdown artifact files should describe workflow outputs and decisions for a specific project or feature.

Not every artifact has to be markdown:

```txt
artifacts/{project}/{feature}/discovery.md
artifacts/{project}/{feature}/technical-plan.md
artifacts/{project}/{feature}/workflow.yaml
artifacts/{project}/{feature}/api-schema.json
artifacts/{project}/{feature}/diagram.png
```

Use this structure:

```md
---
title: Technical Plan
weight: 20
---

# Technical Plan

## Purpose

## Inputs

## Findings

## Decisions

## Implementation Plan

## Verification

## Open Questions
```

Recommended markdown artifact paths:

```txt
artifacts/{project}/{feature}/discovery.md
artifacts/{project}/{feature}/technical-plan.md
artifacts/{project}/{feature}/spec.md
```

MCP markdown resource URI:

```txt
resource://artifact/{project}/{feature}/{name}
```

The `resource://artifact/...` template maps to `.md` files. For non-markdown artifacts, use `read_artifact` with the exact relative path.

## Quality checklist

Before committing a markdown file:

- The file has a clear title.
- The first sections explain purpose and scope.
- Headings are ordered and not skipped.
- Paths and commands are wrapped in backticks.
- Decisions, assumptions, and open questions are explicit.
- The file can be read without external context.
- The file path matches the MCP mapping.
