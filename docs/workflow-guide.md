# Hướng dẫn viết Prompt và Workflow

## Tổng quan

Hệ thống hoạt động theo nguyên tắc: **YAML định nghĩa luồng, Markdown chứa nội dung**.

```
workflows/my-workflow.yaml   ← định nghĩa các bước, thứ tự, artifact nào đọc/ghi
prompts/my-step.md           ← nội dung prompt gửi cho Claude ở bước đó
artifacts/{project}/{feature}/ ← nơi lưu kết quả từng bước
```

Thêm workflow hoặc step mới **không cần sửa code Go**, chỉ cần tạo file YAML và `.md`.

---

## 1. Viết Workflow (YAML)

### Cấu trúc file

```yaml
# workflows/ten-workflow.yaml

id: ten-workflow          # định danh duy nhất, dùng làm prefix prompt name
name: Tên hiển thị       # mô tả ngắn

steps:
  - id: buoc-1            # tên prompt MCP sẽ được đăng ký (phải unique)
    prompt_file: file.md  # file prompt tương ứng trong prompts/
    description: "..."    # mô tả hiển thị trong MCP client
    extra_args:           # tham số bổ sung ngoài project_id và feature_id
      - name: ten_arg
        description: Mô tả tham số
        required: true
    reads: []             # tên artifact cần đọc (phải đã được sinh bởi step trước)
    writes: ten-artifact  # tên artifact step này sẽ ghi ra
```

### Quy tắc quan trọng

**`reads` và `writes` phải khớp nhau giữa các step:**

```yaml
steps:
  - id: discover
    reads: []
    writes: discovery       # ← step này tạo ra "discovery"

  - id: spec
    reads: [discovery]      # ← step này đọc "discovery" (phải đã tồn tại)
    writes: spec

  - id: tasks
    reads: [spec, plan]     # ← đọc nhiều artifact cùng lúc
    writes: tasks
```

**`extra_args` chỉ dùng cho tham số người dùng nhập vào.** Các artifact trong `reads` được load tự động — không cần khai báo thành `extra_args`.

**`id` của step chính là tên prompt MCP.** Claude Code sẽ thấy prompt với tên này.

### Ví dụ đầy đủ

```yaml
id: bug-fix-workflow
name: Bug Fix Analysis

steps:
  - id: analyze_bug
    prompt_file: analyze_bug.md
    description: "Analyst: phân tích bug report → sinh analysis.md"
    extra_args:
      - name: bug_report
        description: Nội dung bug report
        required: true
      - name: severity
        description: "Mức độ: low / medium / high"
        required: false
    reads: []
    writes: analysis

  - id: propose_fix
    prompt_file: propose_fix.md
    description: "Engineer: đề xuất hướng fix → sinh fix_plan.md"
    reads: [analysis]
    writes: fix_plan

  - id: write_tests
    prompt_file: write_tests.md
    description: "QA: viết test cases → sinh tests.md"
    reads: [analysis, fix_plan]
    writes: tests
```

---

## 2. Viết Prompt (Markdown)

### Cấu trúc file

```markdown
<!-- prompts/ten-step.md -->

Vai trò và nhiệm vụ của Claude ở step này.

Hướng dẫn chi tiết:
- Điểm 1
- Điểm 2

Format đầu ra mong muốn.

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/ten-artifact"
- content = nội dung bạn vừa viết

---

<!-- Phần động: các artifact đọc vào và tham số người dùng -->

**Tên artifact:**

{{ten_artifact}}
```

### Placeholder có sẵn

Các placeholder luôn có mặt, không cần khai báo thêm:

| Placeholder | Giá trị | Ví dụ |
|---|---|---|
| `{{project_id}}` | ID project | `my-app` |
| `{{feature_id}}` | Tên feature | `export-task-csv` |

Placeholder cho **artifact** — khai báo trong `reads`:

| Khai báo trong YAML | Placeholder trong `.md` |
|---|---|
| `reads: [discovery]` | `{{discovery}}` |
| `reads: [spec, plan]` | `{{spec}}` và `{{plan}}` |

Placeholder cho **extra_args** — khai báo trong `extra_args`:

| Khai báo trong YAML | Placeholder trong `.md` |
|---|---|
| `extra_args: [{name: request}]` | `{{request}}` |
| `extra_args: [{name: severity}]` | `{{severity}}` |

### Ví dụ prompt cho step đọc một artifact

```markdown
<!-- prompts/propose_fix.md -->

Bạn là một Senior Engineer.

Dựa trên bug analysis bên dưới, hãy đề xuất hướng fix cụ thể.

Yêu cầu:
- Xác định root cause
- Đề xuất ít nhất 2 hướng fix với trade-off
- Chọn hướng recommend và giải thích lý do

Đầu ra markdown với các mục:
1. Root cause
2. Fix options
3. Recommended approach

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/fix_plan"
- content = nội dung bạn vừa viết

---

**analysis.md:**

{{analysis}}
```

### Ví dụ prompt cho step đọc nhiều artifact

```markdown
<!-- prompts/write_tests.md -->

Bạn là một QA Engineer.

Dựa trên analysis và fix plan bên dưới, hãy viết test cases.

...

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/tests"
- content = nội dung bạn vừa viết

---

**analysis.md:**

{{analysis}}

---

**fix_plan.md:**

{{fix_plan}}
```

### Ví dụ prompt cho step có extra_args

```markdown
<!-- prompts/analyze_bug.md -->

Bạn là một Senior Analyst.

Phân tích bug report bên dưới với severity level: **{{severity}}**

...

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/analysis"
- content = nội dung bạn vừa viết

---

**Bug report:**

{{bug_report}}
```

---

## 3. Artifact

Artifact được lưu tự động tại:

```
artifacts/
  {project_id}/
    {feature_id}/
      discovery.md
      spec.md
      plan.md
      tasks.md
```

**Đọc artifact qua MCP Resource:**

```
resource://artifact/{project_id}/{feature_id}/{name}
```

Ví dụ:
```
resource://artifact/my-app/export-task-csv/spec
```

---

## 4. Thêm workflow mới — checklist

```
[ ] Tạo workflows/ten-workflow.yaml
      - Đặt id duy nhất
      - Định nghĩa steps theo thứ tự
      - Đảm bảo reads/writes khớp nhau giữa các step

[ ] Tạo prompts/{ten-step}.md cho mỗi step
      - Dùng đúng placeholder tương ứng với reads và extra_args
      - Luôn có hướng dẫn gọi write_artifact ở cuối
      - Ghi đúng path: "{{project_id}}/{{feature_id}}/{writes}"

[ ] Restart MCP server
      - Prompts mới tự động xuất hiện trong Claude
```

---

## 5. Lỗi thường gặp

| Lỗi | Nguyên nhân | Cách fix |
|---|---|---|
| `artifact "X" not found` | Step trước chưa chạy hoặc `writes` sai tên | Chạy đúng thứ tự step, kiểm tra `writes` trong YAML |
| `prompt file "X" not found` | `prompt_file` trỏ sai tên file | Kiểm tra tên file trong `prompts/` |
| Placeholder không được replace | Sai tên placeholder | Tên trong `.md` phải khớp `reads`/`extra_args` trong YAML |
| Claude ghi sai đường dẫn | `write_artifact path` trong `.md` thiếu `{{project_id}}` | Path phải là `"{{project_id}}/{{feature_id}}/{writes}"` |
