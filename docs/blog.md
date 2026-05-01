# MCP Workbench: Cầu nối Spec-to-Dev mà không tạo file trong repo

---

## 1. Đặt vấn đề

Khi dùng AI assistant (Claude Code, Cursor…) để phát triển tính năng, một quy trình điển hình trông như thế này:

1. Analyst viết **discovery** — phân tích yêu cầu
2. Designer viết **spec** — đặc tả tính năng
3. Engineer viết **technical plan** — thiết kế kỹ thuật
4. Tech Lead breakdown **tasks** — phân công việc

Vấn đề: những tài liệu này rất hữu ích trong quá trình làm việc, nhưng **không thuộc về repo của sản phẩm**. Chúng là output của quá trình tư duy, không phải source code. Nếu để trong repo:

- Phải thêm vào `.gitignore` từng file một
- Dễ bị commit nhầm gây nhiễu lịch sử git
- Mỗi developer có phiên làm việc riêng → xung đột path
- Mỗi AI agent cần biết phải lưu vào đâu, với format gì — không có chuẩn chung

Nếu xóa đi sau khi dùng → mất context, không thể trace lại quyết định thiết kế.

**Câu hỏi cốt lõi:** Làm sao để AI có thể thực hiện toàn bộ workflow spec→develop, sinh ra và đọc các tài liệu trung gian, mà **không tạo một file nào trong repo sản phẩm**?

---

## 2. Giải pháp: MCP Workbench

Ý tưởng: chạy một **MCP server riêng biệt** — gọi là *gateway* — đứng ngoài tất cả các repo sản phẩm. Claude Code kết nối vào gateway này qua giao thức MCP. Mọi artifact (discovery, spec, plan, tasks) được lưu **trong thư mục của gateway**, hoàn toàn tách khỏi repo đang phát triển.

Gateway cung cấp ba loại capability qua MCP:

| Loại | Tool / Prompt / Resource | Vai trò |
|---|---|---|
| **Tools** | `write_artifact`, `read_artifact`, `list_artifacts` | Claude ghi/đọc planning docs ngoài repo |
| **Tools** | `init_project`, `set_context`, `get_context`, `list_context` | Quản lý context (architecture, standards) |
| **Prompts** | `discover`, `spec`, `plan`, `tasks` | Template prompt đã inject sẵn context và artifact từ bước trước |
| **Resources** | `artifact://{project}/{feature}/{name}` | Đọc artifact qua URI |

Workflow chạy theo thứ tự step: mỗi step nhận một MCP prompt đã được gateway render đầy đủ → Claude điền output → gọi `write_artifact` → step tiếp theo đọc artifact vừa được tạo.

---

## 3. Thực hiện

### 3.1 C1 — System Context

```mermaid
C4Context
    title C1: System Context — MCP Workbench

    Person(developer, "Developer", "Software engineer dùng Claude Code để implement feature trong repo sản phẩm")

    System(gateway, "MCP Workbench", "Cầu nối spec-to-dev: quản lý workflow, artifact và context cho AI assistant — hoàn toàn bên ngoài repo sản phẩm")

    System_Ext(claudeCode, "Claude Code", "AI coding assistant. Kết nối MCP Workbench qua giao thức MCP để chạy workflow spec-to-dev")

    System_Ext(productRepo, "Product Repository", "Codebase đang phát triển. Không chứa bất kỳ planning doc hay AI-generated artifact nào")

    System_Ext(contextConfig, "Context Config", "Tài liệu kiến trúc, coding standards, compliance rules do team viết sẵn")

    Rel(developer, claudeCode, "Giao việc, review output", "IDE / terminal")
    Rel(claudeCode, gateway, "Gọi tools, prompts, resources", "MCP protocol (stdio / HTTP+SSE)")
    Rel(developer, gateway, "Cấu hình project context", "CLI / file editor")
    Rel(claudeCode, productRepo, "Đọc / viết source code", "File system")
    Rel(gateway, contextConfig, "Đọc context docs khi build prompt", "File system / Git")

    UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
```

**Điểm then chốt của C1:** Không có đường quan hệ nào từ gateway sang Product Repository. Gateway và repo sản phẩm hoàn toàn cô lập — đây chính là điều đảm bảo zero-intrusion.

---

### 3.2 C2 — Container

```mermaid
C4Container
    title C2: Container — MCP Workbench

    Person(developer, "Developer")
    System_Ext(claudeCode, "Claude Code", "AI coding assistant")

    Container_Boundary(gateway, "MCP Workbench") {
        Container(mcpServer, "MCP Server", "Go · uber/fx · mcp-go", "Entry point nhận MCP JSON-RPC. Wire tất cả handlers qua DI. Hỗ trợ stdio (local) và HTTP+SSE (team)")

        Container(artifactTools, "Artifact Tools", "Go", "Tool: read_artifact / write_artifact / list_artifacts. I/O cho planning docs bên ngoài repo")

        Container(contextTools, "Context Tools", "Go", "Tool: init_project / set_context / get_context / list_context. Fallback: project-scoped → global")

        Container(promptRegistry, "Prompt Registry", "Go", "Scan workflows/*.yaml khi startup. Đăng ký một MCP prompt per step. Inject artifact + context vào template trước khi trả về Claude")

        Container(resourceServer, "Resource Server", "Go", "Expose artifact qua MCP Resource URI: artifact://{project}/{feature}/{name}")

        Container(workflowRunner, "Workflow Runner", "Go · CLI", "Dry-run workflow locally. Simulate step-by-step, ghi placeholder. Dùng để validate YAML trước khi dùng với Claude")

        ContainerDb(artifactStore, "Artifact Store", "File System — artifacts/", "Lưu output của từng step: discovery.md, spec.md, plan.md, tasks.md")

        ContainerDb(contextStore, "Context Store", "File System — context/", "Architecture guidelines, coding standards, compliance rules. Project-scoped override global")

        ContainerDb(workflowDefs, "Workflow Definitions", "YAML — workflows/", "Khai báo workflow steps, thứ tự reads/writes, extra_args, context docs cần dùng")

        ContainerDb(promptTemplates, "Prompt Templates", "Markdown — prompts/", "Nội dung prompt với placeholders: {{project_id}}, {{artifact_name}}, {{context_name}}")
    }

    Rel(claudeCode, mcpServer, "MCP JSON-RPC", "stdio (v1) / HTTP+SSE (v2)")
    Rel(developer, workflowRunner, "Dry-run validate", "CLI")
    Rel(developer, contextStore, "Viết / cập nhật context docs", "File editor")

    Rel(mcpServer, artifactTools, "Route tool calls")
    Rel(mcpServer, contextTools, "Route tool calls")
    Rel(mcpServer, promptRegistry, "Route prompt requests")
    Rel(mcpServer, resourceServer, "Route resource reads")

    Rel(promptRegistry, workflowDefs, "Đọc khi startup để đăng ký prompts")
    Rel(promptRegistry, promptTemplates, "Render template khi Claude gọi prompt")
    Rel(promptRegistry, artifactStore, "Đọc artifact → inject vào prompt")
    Rel(promptRegistry, contextStore, "Đọc context docs → inject vào prompt")

    Rel(artifactTools, artifactStore, "Đọc / ghi artifacts")
    Rel(contextTools, contextStore, "Đọc / ghi context docs")
    Rel(resourceServer, artifactStore, "Đọc artifacts")

    Rel(workflowRunner, workflowDefs, "Đọc workflow YAML")
    Rel(workflowRunner, promptTemplates, "Render prompts (dry-run)")
    Rel(workflowRunner, contextStore, "Đọc context để validate")
    Rel(workflowRunner, artifactStore, "Ghi placeholder artifacts")
```

**Điểm then chốt của C2:** `Prompt Registry` là container trung tâm — nó kết nối tất cả luồng dữ liệu: đọc YAML định nghĩa workflow, load template markdown, inject artifact từ step trước và context docs của project, rồi trả về prompt đã hoàn chỉnh cho Claude. `MCP Server` chỉ là thin router không có logic nghiệp vụ.

Thêm workflow mới **không cần sửa code Go** — chỉ tạo `workflows/ten-workflow.yaml` và `prompts/ten-step.md`, restart server.

---

### 3.3 Luồng dữ liệu một workflow hoàn chỉnh

```
Developer mở Claude Code trong repo sản phẩm
        │
        ▼
Gọi MCP prompt "discover"
  (project_id=my-app, feature_id=export-csv, request="...")
        │
        ▼
Prompt Registry:
  - Đọc workflows/export-task-csv.yaml → lấy step "discover"
  - Render prompts/discover_requirement.md với {{request}}
        │
        ▼
Claude nhận prompt đầy đủ → sinh nội dung discovery
→ gọi write_artifact("my-app/export-csv/discovery", ...)
        │  lưu vào gateway/artifacts/my-app/export-csv/discovery.md
        ▼
Gọi MCP prompt "spec"
        │
        ▼
Prompt Registry:
  - Đọc artifact discovery.md → inject vào {{discovery}}
  - Render prompts/create_feature_spec.md
        │
        ▼
Claude sinh spec.md → write_artifact(...)
        │
        ▼
Gọi MCP prompt "plan"
        │
        ▼
Prompt Registry:
  - Inject spec.md + architecture.md + coding-standards.md
        │
        ▼
Claude sinh plan.md → write_artifact(...)
        │
        ▼
Gọi MCP prompt "tasks" → Claude sinh tasks.md
```

Tất cả file lưu trong `gateway/artifacts/` — **repo sản phẩm không bị chạm đến**.

---

## 4. Kết luận

MCP Workbench giải quyết một vấn đề thực tế và hay bị bỏ qua: **tách bạch lifecycle của tài liệu kỹ thuật khỏi lifecycle của source code**. Tài liệu spec không phải code, không nên sống trong repo code.

Bằng cách dùng MCP như giao thức chuẩn, gateway trở nên hoàn toàn trong suốt với Claude Code — developer không cần biết gateway đang chạy ở đâu, chỉ cần gọi prompt và tool như bình thường. Mọi artifact được quản lý ngoài repo, nhưng vẫn có thể đọc lại bất cứ lúc nào qua MCP resource URI.

Điểm mạnh của thiết kế:

- **Zero-intrusion** vào repo sản phẩm — không một file nào bị tạo thêm
- **Config-driven** — thêm workflow mới chỉ bằng YAML + Markdown, không sửa code Go
- **Portable** — cùng một gateway phục vụ nhiều project, nhiều feature độc lập với nhau

