# Technical Plan: MCP Workbench React UI

## Stack

| Layer | Choice | Lý do |
|---|---|---|
| Frontend | React 18 + TypeScript | SPA, type-safe |
| UI Library | shadcn/ui + Tailwind CSS | Đẹp, không cần design từ đầu |
| Routing | React Router v6 | Standard SPA routing |
| Data fetching | TanStack Query (React Query) | Caching, loading state, refetch |
| State | React Query + local state | Đủ dùng, không cần Redux |
| Build tool | Vite | Nhanh, config đơn giản |
| Backend API | Go — REST over HTTP | Thêm endpoints vào Go server hiện tại |

---

## Cấu trúc thư mục

```
ui/                          # React SPA root
├── src/
│   ├── api/                 # API client functions
│   │   ├── artifacts.ts
│   │   ├── workflows.ts
│   │   └── prompts.ts
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── Artifacts.tsx
│   │   ├── Workflows.tsx
│   │   └── Prompts.tsx
│   ├── components/
│   │   ├── Layout.tsx
│   │   └── ...shared components
│   ├── App.tsx
│   └── main.tsx
├── package.json
├── vite.config.ts
└── tsconfig.json
```

---

## Backend API cần thêm (Go)

### Artifacts
| Method | Path | Mô tả |
|---|---|---|
| GET | `/api/artifacts` | List tất cả artifacts |
| GET | `/api/artifacts/*path` | Đọc nội dung artifact |
| PUT | `/api/artifacts/*path` | Tạo / cập nhật artifact |
| DELETE | `/api/artifacts/*path` | Xóa artifact |

### Workflows
| Method | Path | Mô tả |
|---|---|---|
| GET | `/api/workflows` | List các workflow YAML |
| GET | `/api/workflows/:name` | Đọc định nghĩa workflow |
| POST | `/api/workflows/:name/run` | Trigger workflow với input |
| GET | `/api/workflows/:name/runs` | List các lần chạy gần đây |
| GET | `/api/runs/:id` | Xem status + log của một run |

### Prompts
| Method | Path | Mô tả |
|---|---|---|
| GET | `/api/prompts` | List prompt files |
| GET | `/api/prompts/:name` | Đọc nội dung prompt |
| PUT | `/api/prompts/:name` | Cập nhật prompt |

### System
| Method | Path | Mô tả |
|---|---|---|
| GET | `/api/health` | Uptime, version, stats |

---

## Routing (React Router)

```
/                    → Dashboard
/artifacts           → Artifacts list
/artifacts/*path     → Artifact detail/edit
/workflows           → Workflows list
/workflows/:name     → Workflow detail + trigger
/runs/:id            → Run log viewer
/prompts             → Prompts list
/prompts/:name       → Prompt viewer/editor
```

---

## Phân chia task theo ưu tiên

### Phase 1 — Artifacts
- [ ] Setup Vite + React + shadcn/ui
- [ ] Layout component (sidebar nav)
- [ ] Go: `/api/artifacts` endpoints
- [ ] Artifacts list page
- [ ] Artifact detail / editor (markdown preview)
- [ ] Create / update / delete artifact

### Phase 2 — Workflows
- [ ] Go: `/api/workflows` + `/api/runs` endpoints
- [ ] Workflows list page
- [ ] Workflow detail + trigger form (input params)
- [ ] Run log viewer (polling mỗi 2s)

### Phase 3 — Dashboard
- [ ] Go: `/api/health` endpoint
- [ ] Dashboard page (stats cards + recent runs table)

### Phase 4 — Prompts
- [ ] Go: `/api/prompts` endpoints
- [ ] Prompts list + editor page

---

## CORS & Dev Setup

- Vite dev server chạy ở `http://localhost:5173`
- Go server chạy ở `http://localhost:8080` (hoặc port hiện tại)
- Go cần thêm CORS middleware cho `localhost:5173`
- Vite config proxy `/api` → Go server (tránh CORS khi dev)
