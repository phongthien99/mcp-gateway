---
title: Architecture
weight: 10
---

# Architecture

## Purpose

Mô tả stack công nghệ, cấu trúc tầng, và các quy ước thiết kế để AI engineer có đủ context khi sinh code mới.

## Scope

Bao gồm: frontend stack, directory layout, data layer, state management, PDF engine, translation. Không bao gồm CI/CD, backend, hay deployment.

## Tech Stack

| Layer | Thư viện / Công cụ |
|---|---|
| UI Framework | React 18, MUI v7 (`@mui/material`), Emotion |
| Language | TypeScript 5 (strict) |
| Build | Vite 6 |
| Server State | TanStack Query v5 (`@tanstack/react-query`) |
| Validation | Zod v4 |
| PDF Render | `react-pdf` v10 + `pdfjs-dist` (worker via unpkg CDN) |
| PDF Drawing | `react-konva` + `konva` v9 |
| PDF Engine | `@embedpdf/*` plugin suite (scroll, zoom, annotation, search, …) |
| Translation | Chrome Translator API (on-device, `chrome-translator.repository.ts`) |
| Markdown render | `react-markdown` + `remark-gfm` |

## Directory Layout

```
src/
  features/
    books/
      schema/          # Zod schemas + inferred types (source of truth)
      dto/             # Data-transfer shapes (plain objects)
      repositories/    # Storage adapters (localStorage, async interface)
      facade/          # Aggregating hooks that compose queries + local state
      hooks/           # use<Name>Query / use<Name>Mutation wrappers
      components/      # Feature-scoped React components (PascalCase)
    pdf-notes/         # Same sub-folder pattern
    pdf-translation/   # Same sub-folder pattern
  components/          # Shared UI components
  hooks/               # Shared React hooks
  pages/               # Route-level components
  icons.tsx
  app-shell.tsx
  application.tsx      # PDF viewer page (large stateful component)
  main.tsx
```

## Layering Rules

1. **Schema layer** — Zod schemas define the canonical shape; TypeScript types are `z.infer<>` only.
2. **Repository layer** — Plain async objects (not classes). Interact with `localStorage` (key prefix `react-mui.*`). Never import React.
3. **Hooks layer** — `use<Name>Query` wraps `useQuery`; `use<Name>Mutation` wraps `useMutation`. Each hook calls exactly one repository method.
4. **Facade layer** — `use<Name>Facade` composes multiple hooks + `useState` into a single typed object consumed by the page/component.
5. **Component layer** — Receives facade output as props or calls the facade directly. Uses MUI primitives exclusively for UI.

## Data Storage

Decision: all data is stored in `localStorage`; there is no backend API.

- Books: key `react-mui.books` (JSON array, validated with `BooksSchema`)
- PDF notes: keyed per PDF URL
- PDF drawings: key prefix `react-mui.pdf-drawings.<pdfUrl>`

## State Management

- **Remote/async state** → TanStack Query (queryKey conventions: `[entity, ...params]`)
- **Local UI state** → `useState` / `useReducer` in facades or components
- No global state store (Redux, Zustand, etc.)

## PDF Rendering

- `react-pdf` renders pages via `pdfjs-dist`; worker configured once in `application.tsx`.
- Pages are lazy-loaded via `IntersectionObserver` (`rootMargin: '600px 0px'`).
- Drawing overlay: `PdfDrawingLayer` wraps `<Page>` with a Konva `<Stage>`.
- Zoom: baked into `width` prop only — never pass both `width` and `scale` to `<Page>`.
