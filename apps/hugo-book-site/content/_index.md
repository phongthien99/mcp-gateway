---
title: MCP Workbench
type: docs
---

# MCP Workbench

This is the standalone Hugo Book documentation app for the MCP Workbench workspace.

Use the left navigation to browse project docs, workflows, context, and artifact notes.

## Local development

```sh
cd apps/hugo-book-site
hugo server --bind=0.0.0.0 --port=1313
```

The Docker Compose setup mounts this app into the MCP dashboard service and serves it from the configured dashboard port.
