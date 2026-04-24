# Coding Standards — CLS Frontend Commercial

> Version: 1.0 | Ratified: 2026-04-24 | Last Amended: 2026-04-24

---

## Core Principles

1. **Feature Modularity** — Each feature is a self-contained module under `src/modules/`. Cross-module imports are forbidden.
2. **Separation of Concerns** — Strict layered flow: `View → Facade → Hook → Service → Axios`. No layer may skip or reverse this chain.
3. **Server State is Sacred** — All server data is owned by React Query. Never copy server data into `useState`.
4. **Type Safety Everywhere** — TypeScript strict mode, no `any`, Zod schemas for all forms, typed DTOs for all API responses.
5. **UI Consistency** — MUI v7 is the sole UI library. Tailwind for layout utilities only. No raw CSS.

---

## Layer Rules

### View (`view/`)
- Calls **only the facade hook** for data and actions.
- Passes everything to child components as props.
- Owns page-level layout structure and route-level logic.
- One view component per module entry point.

### Facade (`facade/`)
- Composes multiple hooks into one interface consumed by the View.
- Coordinates cross-cutting concerns: trigger refetch after mutation, combine loading states, etc.
- **No JSX, no direct axios/service calls.**

### Hooks (`hooks/`)
- Wraps `useQuery` / `useMutation` from React Query.
- Owns query keys, stale time, cache invalidation.
- Calls service functions only — never calls axios directly.

### Services / Repositories (`services/`)
- Calls `axiosInstance` endpoints and maps response → typed DTO.
- Pure async functions — no React, no hooks, no side effects beyond the HTTP call.

### Components (`components/`)
- Pure UI renderers. Accept data via props.
- **Forbidden imports**: `services/`, `hooks/`, `facade/`, `axios`, React Query.
- If a component needs data fetching, that is a sign the data should be lifted to the View via the Facade.

### Schema (`schema/`)
- Zod schemas only. No business logic.
- Each form has its own schema file. Schema is the single source of truth for input shape.

### Types / DTO (`types.ts` or `dto/`)
- Plain TypeScript interfaces and types.
- Scoped to the module. Global types live in `src/types/`.

---

## TypeScript

- **Strict mode is on** — `strict: true` in `tsconfig.json`.
- **No `any`** — use `unknown` + type guard if the shape is genuinely unknown.
- All API response shapes must have a corresponding DTO type.
- All form inputs must be typed from their Zod schema via `z.infer<typeof schema>`.
- Prefer `type` over `interface` for object shapes; use `interface` only when declaration merging is needed.

---

## Component Design

### Single Responsibility
Each component does one thing. A component that renders a list **and** handles selection **and** manages a dialog must be broken up into focused sub-components.

### File Length Cap
**No single file may exceed 300 lines.** If a file grows past this limit, split it into focused, single-responsibility units before merging.

### Props Limit
A component with more than **8 props** is a signal it should be split or composed differently. Props drilling beyond 2 levels is a signal to extract a subcomponent or move state up.

### No Complex Inline Logic in JSX
Ternaries beyond 1 level, array transforms, and conditional chains must be extracted into named variables or helper functions **before** the `return` statement. JSX should read like a template, not an algorithm.

```tsx
// Bad
return (
  <div>
    {items.filter(i => i.active).map(i => (
      <Item key={i.id} label={i.type === 'A' ? 'Alpha' : i.type === 'B' ? 'Beta' : 'Other'} />
    ))}
  </div>
);

// Good
const activeItems = items.filter(i => i.active);
const labelFor = (type: string) => ({ A: 'Alpha', B: 'Beta' }[type] ?? 'Other');

return (
  <div>
    {activeItems.map(i => <Item key={i.id} label={labelFor(i.type)} />)}
  </div>
);
```

### Prefer Composition Over Configuration
Build behavior by composing small components and hooks, not by adding more props and flags to an existing one. Adding a boolean prop to change rendering behavior is usually a sign the component should be split.

---

## State Management

| Data type | Where it lives |
|---|---|
| Server/async data | React Query (`useQuery` / `useMutation`) |
| Global UI state (theme, sidebar open, etc.) | Zustand store |
| Local ephemeral UI state (open/closed, hover) | `useState` inside the component |
| Form state | React Hook Form |

**Never** copy server data into `useState`. Use `select` or derived variables from the query result instead.

---

## Forms

- Every form uses **React Hook Form** + **Zod** schema resolver (`@hookform/resolvers/zod`).
- Schema lives in the module's `schema/` folder.
- Input types are inferred from schema: `type FormValues = z.infer<typeof mySchema>`.
- Use shared form components from `src/components/hook-form/`.

---

## Error and Loading States

- Loading and error states **must** use the shared components from `src/components/` (e.g., `LoadingScreen`, `EmptyContent`).
- Ad-hoc spinners, inline `"Loading..."` strings, or custom error `<div>`s are not permitted.
- Every async boundary must have a consistent UX.

---

## Styling

- **MUI `sx` prop or `styled()`** for component-level styles.
- **Tailwind utility classes** for layout spacing (`flex`, `gap-*`, `p-*`, `w-full`, etc.) when `sx` would be verbose.
- **No raw CSS files**, no `style={{}}` inline objects for anything other than dynamic values that cannot be expressed in MUI/Tailwind.
- Do not override MUI component internals via CSS class selectors.

---

## Imports

- Use the `src/` path alias — never use relative paths that traverse more than 1 level (`../../`).
- Import order is enforced by `eslint-plugin-perfectionist` and `eslint-plugin-import`. Run `yarn lint:fix` before committing.
- No circular imports. Modules must not import from each other.

---

## Comments

Write no comments by default. Add a comment only when the **why** is non-obvious: a hidden constraint, a subtle invariant, a known framework bug workaround. Never describe what the code does — well-named identifiers already do that.

---

## Naming Conventions

| Thing | Convention | Example |
|---|---|---|
| React components | PascalCase | `SearchResultCard` |
| Hooks | `use` + PascalCase | `useSearchResults` |
| Facade hooks | `use` + PascalCase + `Facade` | `useSearchFacade` |
| Service functions | camelCase verb phrase | `fetchSearchResults` |
| Zod schemas | camelCase + `Schema` | `searchFilterSchema` |
| DTO types | PascalCase + `DTO` or `Response` | `SearchResultDTO` |
| Files | kebab-case | `search-result-card.tsx` |
| Folders | kebab-case | `search-layout/` |

---

## Git & Branch Hygiene

- Branch format: `feature/TICKET-ID-short-description` or `fix/TICKET-ID-short-description`.
- Commit messages must reference the ticket: `CLEXIX-123: short imperative description`.
- All commits go through `lefthook` pre-commit (ESLint + Prettier via `lint-staged`).
