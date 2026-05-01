---
title: Artifacts
weight: 40
---

# Artifacts

Workflow artifacts generated or read by the MCP artifact tools.

## Source of truth

Artifact files live under:

```txt
artifacts/{project}/{feature}/{file}
```

MCP reads these files with:

```txt
read_artifact(path="{project}/{feature}/{file}")
write_artifact(path="{project}/{feature}/{file}", content="...")
list_artifacts(workflow="{project}")
```

For markdown artifacts, the MCP resource URI format is:

```txt
resource://artifact/{project}/{feature}/{name}
```

## Writing rules

Use the shared [Markdown Rules]({{< relref "/docs/markdown-rules" >}}) for markdown artifact files.

When `write_artifact` creates a new folder, it also creates missing Hugo `_index.md` files so artifact folders appear in the documentation tree.

Artifact documents should capture workflow output:

- Discovery notes.
- Feature specs.
- Technical plans.
- Implementation summaries.
- Verification results.
- Open questions for the next workflow step.

Do not put long-term project rules in artifacts. Put durable rules in `context/`.

Agents should choose the artifact file extension based on the content:

| Content | Recommended file |
| --- | --- |
| Discovery, specs, plans, summaries | `.md` |
| Workflow/config snippets | `.yaml` or `.yml` |
| Structured schemas or API examples | `.json` |
| Plain logs or raw output | `.txt` |
| Diagrams or screenshots | image extension such as `.png` |
