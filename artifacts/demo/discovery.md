# Discovery: Thêm Giao Diện (UI) cho MCP Workbench

## 1. Goal

Xây dựng một **React SPA** cho MCP Workbench, cho phép người dùng tương tác với hệ thống (Artifacts, Workflows, Dashboard, Prompts) qua trình duyệt — không cần đăng nhập.

---

## 2. Assumptions

- Project hiện tại là một **MCP Workbench** — backend Go xử lý tools, workflows, artifacts, và prompts.
- Frontend là **React SPA** (single-page app), chạy tách biệt, gọi API từ Go backend.
- **Không có authentication** — UI public trong môi trường local/internal.
- Go backend sẽ cần expose thêm **REST API** để phục vụ frontend.
- Thứ tự ưu tiên phát triển: **Artifacts → Workflows → Dashboard → Prompts**.

---

## 3. Scope

### Artifacts
- Xem danh sách artifacts (list)
- Đọc nội dung artifact (read)
- Tạo / cập nhật artifact (write)
- Xóa artifact

### Workflows
- Xem danh sách workflows (từ thư mục `workflows/`)
- Xem chi tiết định nghĩa workflow (YAML)
- Trigger/chạy một workflow với input params
- Xem trạng thái và output/log của workflow đang chạy (có thể polling hoặc SSE)

### Dashboard
- Tổng quan: số lượng artifacts, workflows, prompts
- Trạng thái gateway (uptime, version)
- Các workflow gần đây đã chạy

### Prompts
- Xem danh sách prompt templates (từ thư mục `prompts/`)
- Xem nội dung từng prompt
- Chỉnh sửa / lưu prompt

---

## 4. Out of Scope

- Authentication / Authorization
- Mobile app
- Multi-tenant
- CI/CD pipeline cho frontend
- Tích hợp external services qua UI

---

## 5. Open Questions (đã giải quyết)

| Câu hỏi | Quyết định |
|---|---|
| Frontend framework | React SPA |
| Cách serve | SPA riêng (dev: port khác, prod: có thể serve static từ Go) |
| Authentication | Không cần |
| Ưu tiên tính năng | Artifacts → Workflows → Dashboard → Prompts |
| Real-time workflow log | TBD — polling trước, nâng lên SSE nếu cần |
