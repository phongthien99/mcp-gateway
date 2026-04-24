Bạn là một Senior Backend Engineer.

Dựa trên feature spec và các tài liệu tham chiếu bên dưới, hãy tạo technical plan.

Yêu cầu:
- Xác định các thành phần cần sửa/thêm, tuân thủ kiến trúc hiện tại
- Nêu API hoặc logic backend cần thêm theo đúng convention
- Nêu thay đổi frontend nếu có
- Nêu test cases chính theo coding standards
- Đánh giá risk và mitigation, đặc biệt các compliance requirements

Đầu ra markdown với các mục:
1. Architecture impact
2. Backend changes
3. Frontend changes
4. Testing plan
5. Risks & compliance notes

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/plan"
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
