# Compliance Rules — CLS Frontend Commercial

> Version: 1.0 | Ratified: 2026-04-24 | Last Amended: 2026-04-24

---

## Linting is a Hard Gate

All of the following checks **must pass** before a merge request is accepted. CI blocks merges on failure.

| Check | Command | Scope |
|---|---|---|
| ESLint | `yarn lint` | All `src/**/*.{js,jsx,ts,tsx}` |
| Prettier | `yarn fm:check` | All `src/**/*.{js,jsx,ts,tsx}` |
| TypeScript | `tsc` (run as part of `yarn build`) | Whole project, strict mode |

Run `yarn fix:all` locally to auto-fix ESLint and Prettier violations before pushing.

`lefthook` enforces ESLint + Prettier on every `git commit` via `lint-staged`. Never skip hooks (`--no-verify`).

---

## Protected Files — Peer Review Required

Changes to the following files or directories **require at least one peer review approval** before merging, regardless of ticket size:

| File / Directory | Reason |
|---|---|
| `src/auth/` (entire directory) | Auth context is critical. A bug here affects every user session. |
| `src/routes/paths.ts` | Central path registry. A rename silently breaks navigation across the app. |
| `src/lib/axios.ts` | Shared HTTP client, token refresh logic, and endpoint registry. |
| `src/theme/` (entire directory) | Theme changes have visual regressions across all pages. |

---

## Breaking Change Policy

The following are considered **breaking changes** and must be flagged explicitly in the MR description:

- Renaming or removing a key in `paths.ts`.
- Changing the `endpoints` object structure in `src/lib/axios.ts`.
- Modifying the MUI theme token names or palette keys in `src/theme/`.
- Removing or renaming a prop from a globally shared component in `src/components/`.
- Changing the shape of a DTO that is consumed by multiple modules.

Breaking changes require a migration plan documented in the MR description.

---

## Vendor Chunk Splitting Rules

The `vite.config.ts` `manualChunks` function defines the chunk manifest. **Do not move packages between chunks or add new heavy packages without updating this manifest.**

| Chunk name | Packages |
|---|---|
| `charts-vendor` | `react-apexcharts`, `apexcharts` |
| `pdf-vendor` | `react-pdf`, `pdfjs-dist` |
| `mui-x-vendor` | `@mui/x-data-grid`, `@mui/x-date-pickers`, `@mui/x-tree-view` |
| `mui-icons-vendor` | `@mui/icons-material` |
| `mui-core-vendor` | `@mui/material`, `@mui/lab`, `@emotion/*` |
| `motion-vendor` | `framer-motion` |
| `react-vendor` | `react`, `react-dom`, `react-router`, `scheduler` |
| `app-runtime-vendor` | `i18next`, `dayjs`, `date-fns`, `@tanstack/react-query` |

The chunk size warning limit is **700 kB**. Adding a new dependency that causes a chunk to exceed this must be justified and the chunk split must be updated accordingly.

---

## No Test Suite

There is currently **no automated test suite**. This places additional responsibility on code structure:

- Business logic **must** live in hooks and facades — not in components.
- Components must be pure renderers so that logic can be reviewed and reasoned about in isolation.
- Code reviewers must manually verify logic correctness during MR review.
- This rule exists to ensure that when tests are introduced in the future, the logic is already in testable, isolated units.

---

## File Length Cap

**No single file may exceed 300 lines.** If a file reaches this limit before merging, it must be split into focused, single-responsibility units. This is enforced during code review — an oversized file is grounds to request changes.

---

## Standardised Error / Loading UI

- All loading states must use components from `src/components/loading-screen/` or equivalent shared components.
- All empty / error states must use `src/components/empty-content/` or `src/components/search-not-found/`.
- Ad-hoc inline `"Loading..."` strings, custom spinners, or one-off error `<div>`s are **not permitted**.
- Every async data boundary must render a consistent UX to the user.

---

## Forbidden Patterns

| Pattern | Why forbidden |
|---|---|
| `any` in TypeScript | Defeats type safety; use `unknown` + type guard |
| Copying React Query data into `useState` | Creates stale state and double source of truth |
| Direct `axios` import in components or views | Bypasses the service layer; violates separation of concerns |
| Cross-module imports (`modules/A` importing from `modules/B`) | Breaks feature isolation; creates hidden coupling |
| Raw CSS files or `style={{}}` objects for static styles | Bypasses MUI/Tailwind system; causes theme inconsistencies |
| Nested ternaries in JSX (more than 1 level) | Unreadable; extract to a named variable or helper |
| Components with more than 8 props | Sign of missing decomposition; split or compose |
| Props drilling beyond 2 levels | Extract a subcomponent or lift state |
| `--no-verify` on git commit | Bypasses lint-staged gate |

---

## Dependency Policy

- **UI components**: MUI v7 only. Do not introduce a second component library.
- **Date handling**: `dayjs` or `date-fns` (both are present). Do not add `moment`.
- **Utility functions**: prefer `es-toolkit` over `lodash` for new code; `lodash` is present but should not be extended.
- New dependencies must be discussed and approved before adding to `package.json`.
- Dev dependencies must not be imported in `src/` production code.

---

## Security Baseline

- All user-generated HTML rendered via `react-markdown` or `dangerouslySetInnerHTML` must be sanitised with `dompurify` first.
- No secrets, API keys, or credentials may be committed to the repository. Use Vite env variables (`VITE_*`) and `.env.local` (git-ignored).
- The `X-App-Type: commercial` header is injected by the axios interceptor — do not override or remove it.
- reCAPTCHA (`react-google-recaptcha`) must remain active on all public-facing forms.
