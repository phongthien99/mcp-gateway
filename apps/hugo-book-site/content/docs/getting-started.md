---
title: Getting Started
weight: 10
---

# Getting Started

## What is MCP Workbench?

**MCP Workbench** is a local workspace that exposes your project context, artifacts, and files through a single [MCP](https://modelcontextprotocol.io) endpoint — compatible with any MCP-capable AI client (Claude, Cursor, Copilot, and others). It provides:

- **An MCP server** — exposes your project context, artifacts, and files as MCP tools and resources.
- **A documentation dashboard** — a Hugo-powered site that renders your project docs, specs, and notes in a browser.
- **An in-page file editor** — a lightweight file API so your AI client (or you) can read and write content directly from the dashboard.

The goal is a tight feedback loop: the AI client reads context from the workbench, generates or updates artifacts, and you review them in the live dashboard — all without leaving your local environment.

---

## Features

- **Domain-specific prompts** — define prompt templates as `.md` files tailored to any domain (software engineering, data science, marketing, etc.).
- **Custom workflows** — declare multi-step processes in YAML; each step links to a prompt template, reads artifacts from previous steps, and writes new ones.
- **Auto-registration** — the MCP server automatically registers all workflows and prompts at startup with no code changes required.
- **Context injection** — each step can reference project context documents (architecture, coding standards, compliance rules) so the AI client has full background.
- **Artifact pipeline** — output of each step (discovery, spec, plan, tasks) is persisted and passed as input to the next step.

---