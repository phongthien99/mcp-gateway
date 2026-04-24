# Architecture — CLS Frontend Commercial

> Version: 1.0 | Ratified: 2026-04-24 | Last Amended: 2026-04-24

---

## Folder Map (`src/`)

```
src/
├── assets/           # Static assets: icons, images, illustrations, logo, seed data
├── auth/             # Authentication subsystem (context, guards, hooks, services, views)
├── components/       # Global shared UI components (reusable across all modules)
├── hooks/            # Global shared React hooks
├── layouts/          # Page shell layouts (dashboard, auth-split, simple, high-level)
├── lib/              # Low-level infrastructure: axios instance, endpoints registry
├── locales/          # i18n resources and language utilities
├── modules/          # Feature modules — primary business logic lives here
├── pages/            # Route entry points — thin wrappers that import module views
├── routes/           # Route definitions, lazy loaders, route hooks
├── sections/         # Non-module UI sections (error, payment-return, blank)
├── theme/            # MUI theme configuration and settings overrides
├── types/            # Global TypeScript type definitions
└── utils/            # Pure utility functions (formatters, mappers, helpers)
```

---

## Module Anatomy (`src/modules/<feature>/`)

Each feature module is a self-contained vertical slice. All layers are optional — use only what the feature needs — but the responsibilities below are strictly scoped.

| Layer | Folder | Responsibility |
|---|---|---|
| **View** | `view/` | Container component. The **only** place allowed to call the facade hook. Passes data down to `components/` as props. Owns page-level layout and routing logic. One view per module entry point. |
| **Facade** | `facade/` | Composes multiple hooks into a single interface for the View. Coordinates cross-cutting concerns (e.g., trigger refetch after mutation). No UI, no direct API calls. |
| **Hooks** | `hooks/` | Wraps React Query `useQuery` / `useMutation`. Owns query keys, caching strategy, and cache invalidation. Calls repository/service functions only. |
| **Services / Repositories** | `services/` | Calls axios endpoints, maps raw API response → typed DTO. No business logic, no React dependencies. |
| **Components** | `components/` | Pure UI for this module. Receives all data via props. No API calls, no service/repository imports, no facade imports. |
| **Schema** | `schema/` | Zod schemas for form validation. Source of truth for input shape within this module. |
| **Types / DTO** | `types.ts` or `dto/` | TypeScript types for request payloads and API responses. Shared across layers within the module only. |
| **Utils** | `utils/` | Module-scoped pure utility functions. No React, no side effects. |
| **Constants** | `constants/` | Module-scoped constants and static lookup data. |

### Data Flow (strict, one-direction)

```
Page (src/pages/)
  └── View (view/)
        └── Facade (facade/)
              ├── Hook A (hooks/)  →  Service/Repository (services/)  →  axios (src/lib/axios.ts)
              └── Hook B (hooks/)  →  Service/Repository (services/)  →  axios (src/lib/axios.ts)
        └── Components (components/)   ← receives props only, never imports upward
```

---

## Tech Stack

| Layer | Decision |
|---|---|
| **Runtime** | React 19, TypeScript 5.8 (strict) |
| **Build** | Vite 6 + `@vitejs/plugin-react-swc` |
| **Routing** | React Router v7 (file-per-page under `src/pages/`, lazy-loaded) |
| **Server State** | TanStack React Query v5 — owns all remote data, caching, and invalidation |
| **Client State** | Zustand v5 — for UI-only global state that has no server representation |
| **Auth** | JWT (access + refresh tokens); `src/lib/axios.ts` handles silent refresh and force-logout on 401; `src/auth/` owns all auth context and guards |
| **HTTP** | Axios v1 via shared `axiosInstance` in `src/lib/axios.ts`; all endpoints registered in the `endpoints` const in that file |
| **UI** | MUI v7 (`@mui/material`, `@mui/lab`, `@mui/x-data-grid`, `@mui/x-date-pickers`, `@mui/x-tree-view`) — sole UI component library |
| **Styling** | Emotion (MUI default) for component styles; Tailwind CSS v4 for layout utilities only; no raw CSS files |
| **Animation** | Framer Motion v12 |
| **Forms** | React Hook Form v7 + `@hookform/resolvers` + Zod v3 |
| **i18n** | i18next v25 + `react-i18next` + `i18next-browser-languagedetector` |
| **Charts** | ApexCharts via `react-apexcharts` |
| **PDF** | `react-pdf` + `pdfjs-dist` |
| **Icons** | `@iconify/react` (primary) + `lucide-react` + `@mui/icons-material` |
| **Notifications / Toasts** | `sonner` (toast) + `sweetalert2` (confirm dialogs) |
| **Package Manager** | Yarn 1 (`yarn.lock` is source of truth) |
| **Node** | ≥ 20 |

---

## Infrastructure Details

### Axios (`src/lib/axios.ts`)
- Single `axiosInstance` with `baseURL = CONFIG.serverUrl`.
- Request interceptor: attaches `Authorization: Bearer <access>` and `X-App-Type: commercial`.
- Response interceptor: handles 401 (silent token refresh with queue), 503 (maintenance event), and force-logout on refresh failure.
- All API endpoints are declared in the `endpoints` const in the same file — no magic strings elsewhere.

### Auth (`src/auth/`)
- JWT stored in `localStorage` via `setSession` / `getAccess` / `getRefresh`.
- Auth context provided at app root; protected routes use guards in `src/auth/guard/`.
- **Critical path** — no edits without peer review.

### Routing (`src/routes/`)
- Pages under `src/pages/` are thin wrappers importing module views.
- All routes are lazy-loaded via React Router's `lazy()`.
- Centralised path constants in `src/routes/paths.ts` — changes require peer review.

### Theme (`src/theme/`)
- MUI theme configured in `src/theme/core/`; settings overrides in `src/theme/with-settings/`.
- Changes to the theme require peer review.

### Vendor Chunk Splitting (`vite.config.ts`)
See the **Governance** document for the chunk manifest.
