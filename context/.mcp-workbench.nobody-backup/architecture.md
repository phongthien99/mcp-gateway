---
title: "Architecture"
weight: 10
---

# Architecture

## Stack

- **Backend**: Go 1.26, `uber/fx` (DI), `mark3labs/mcp-go` (MCP protocol), `gopkg.in/yaml.v3`
- **Dashboard**: Hugo + hugo-book theme (hot-reload dev server)
- **Transport**: SSE (default, port 8099) hoặc stdio (set `MCP_TRANSPORT=stdio`)
- **Build**: Go workspace (`go.work`) — monorepo với một module duy nhất tại `apps/mcp-server`

## Layer Structure

```
main.go
  └── module.NewApp()          ← khởi tạo fx.App, wire tất cả handler
        ├── tools.*            ← MCP Tools: hành động Claude có thể gọi
        ├── resources.*        ← MCP Resources: dữ liệu Claude đọc qua URI
        └── prompts.*          ← MCP Prompts: template prompt động từ YAML
```

Mỗi handler là một struct với constructor `NewXxx()` và method `Register(*mcpserver.MCPServer)`. Đăng ký vào fx qua `mcp.AsHandler(NewXxx)`.

## Workflow Engine

Hệ thống hoạt động hoàn toàn dựa trên file — **không cần sửa code Go** khi thêm workflow mới:

```
workflows/*.yaml      ← định nghĩa steps, reads/writes, context refs
prompts/*.md          ← nội dung prompt với {{placeholder}} syntax
artifacts/{project}/{feature}/  ← output Claude sinh ra (per step)
context/global/       ← tài liệu tham chiếu dùng chung
context/{project_id}/ ← override cho project cụ thể (ưu tiên hơn global)
```

Khi server khởi động, `WorkflowPrompts.Register()` scan `workflows/*.yaml` và tự đăng ký một MCP prompt cho mỗi step.

## API Conventions

- MCP Tools trả về `*mcp.CallToolResult` — dùng `mcp.NewToolResultText()` cho success, `mcp.NewToolResultError()` cho lỗi (không return Go error).
- MCP Resources URI scheme: `resource://artifact/{project}/{feature}/{name}`
- Placeholder trong prompt template: `{{tên_biến}}` — replace bằng string, không dùng Go template engine.

## References

- `docs/workflow-guide.md` — hướng dẫn đầy đủ cách viết workflow và prompt
- `apps/mcp-server/src/` — toàn bộ source Go
