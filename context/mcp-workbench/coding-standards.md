---
title: "Coding Standards"
weight: 20
---

# Coding Standards

## Language / Framework

- **Go**: tuân theo `gofmt` + `golangci-lint`. Chạy lint trước khi commit.
- Tên package: lowercase, không underscore (`tools`, `resources`, `prompts`, `workflow`).
- Constructor pattern: `NewXxx() *Xxx` — không dùng bare struct literal bên ngoài package.
- Handler mới phải implement `Register(*mcpserver.MCPServer)` và đăng ký vào `module/app.go` qua `mcp.AsHandler(NewXxx)`.
- Không dùng `init()` function.

## Workflow / Prompt Files

- Mỗi step trong YAML phải có `writes` nếu sinh ra artifact; `reads` phải khớp với `writes` của step trước.
- Placeholder trong `.md` phải khớp chính xác tên trong `reads`, `context`, `extra_args` của YAML tương ứng.
- Path trong `write_artifact` luôn có dạng `"{{project_id}}/{{feature_id}}/{tên_artifact}"`.

## Git

- Branching: **Gitflow** — `main`, `develop`, `feature/*`, `release/*`, `hotfix/*`.
- Commit message: tiếng Anh, imperative mood, prefix theo loại: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`.
- PR merge vào `develop`; chỉ `release/*` và `hotfix/*` merge vào `main`.

## Testing

- Không bắt buộc unit test. Ưu tiên test thủ công qua MCP Inspector hoặc Claude client.
