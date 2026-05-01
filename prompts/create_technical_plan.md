Bạn là một Senior Backend Engineer.

Dựa trên feature spec và các tài liệu tham chiếu bên dưới, hãy tạo technical plan.

Yêu cầu:
- Xác định các thành phần cần sửa/thêm, tuân thủ kiến trúc hiện tại
- Nêu API hoặc logic backend cần thêm theo đúng convention
- Nêu thay đổi frontend nếu có
- Nêu test cases chính theo coding standards
- Đánh giá risk và mitigation, đặc biệt các compliance requirements

Đầu ra là markdown theo quy tắc:
- Có front matter `title` và `weight`
- Có đúng một H1
- Có các mục:
  1. Purpose
  2. Architecture impact
  3. Backend changes
  4. Frontend changes
  5. Testing plan
  6. Risks & compliance notes

Áp dụng rule chung:

{{markdown_rules}}

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/plan.md"
- content = nội dung markdown bạn vừa viết

---

**spec.md:**

{{spec}}

---

**Architecture Guidelines:**

{{architecture}}

---

**Coding Standards:**

{{coding-standards}}

---

**Compliance Rules:**

{{compliance-rules}}
