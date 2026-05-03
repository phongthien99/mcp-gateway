---
title: Coding Standards
weight: 20
---

# Coding Standards

## Purpose

Quy tắc đặt tên, cấu trúc file, error handling và git workflow để đảm bảo code nhất quán.

## Scope

Áp dụng cho toàn bộ code trong `src/`. Không áp dụng cho config files (`vite.config.ts`, `eslint.config.js`).

## Naming Conventions

| Artifact | Convention | Example |
|---|---|---|
| React component file | PascalCase `.tsx` | `BookList.tsx` |
| Hook / util / repository file | kebab-case `.ts` | `use-is-mobile.ts`, `book.repository.ts` |
| Schema file | kebab-case `.schema.ts` | `book.schema.ts` |
| React component | PascalCase | `function BookList(...)` |
| Query hook | `use<Entity>Query` | `useBooksQuery` |
| Mutation hook | `use<Entity><Action>Mutation` | `useDeleteBookMutation` |
| Facade hook | `use<Feature>Facade` | `useBooksFacade` |
| Repository object | `<entity>Repository` (camelCase) | `bookRepository` |
| localStorage key | `react-mui.<entity>` | `react-mui.books` |

## TypeScript Rules

- `strict: true` — no implicit `any`, no non-null assertions without justification.
- Types come from `z.infer<typeof Schema>` — never duplicate a Zod schema as a manual interface.
- Use `type` imports (`import type { Foo }`) for type-only imports.
- Prefer `unknown` over `any` at boundaries; narrow with Zod `.safeParse`.

## Zod Usage

```ts
// Define schema first
export const BookSchema = z.object({ ... });
// Infer type — do NOT write a manual interface
export type Book = z.infer<typeof BookSchema>;
```

- Parse external / stored data with `.safeParse`; on failure return a safe default or throw.
- Surface `ZodError` messages via `err.issues[0]?.message`.

## React & Hooks

- Facade hooks return a plain typed object — not a tuple.
- `useCallback` / `useMemo` only when the dependency is passed as a prop or used in another hook (not premature optimisation).
- Never call a hook conditionally.
- Colocate local state in the smallest component that needs it; lift only when necessary.

## Comments

- Default: no comments.
- Add a comment only when the WHY is non-obvious (hidden constraint, workaround, subtle invariant).
- No JSDoc on every function; no `// What this does` comments.

## Error Handling

- Validate at storage boundaries (read from `localStorage`) and user-input boundaries only.
- Trust TanStack Query's `onError` / `isError` for mutation errors; surface to UI via facade `dialogError` field.
- Silent catch (`catch {}` or `catch { return [] }`) is acceptable only for localStorage read failures where a safe default exists.

## Git Workflow

- Branch off `master`; PR back to `master`.
- Commit message format: `type(scope): short description` — types: `feat`, `fix`, `refactor`, `docs`, `chore`.
- One logical change per commit.
- No `--no-verify`.

## Linting

```bash
npm run lint   # eslint with @typescript-eslint, react-hooks, react-refresh
```

Max warnings: 0. Fix lint errors before committing.
