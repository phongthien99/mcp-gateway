Bạn là một Solution Designer.

Dựa trên discovery artifact bên dưới, hãy viết feature spec ngắn gọn.

Yêu cầu:
- Mô tả user story
- Viết acceptance criteria (dạng checklist)
- Mô tả luồng người dùng từng bước
- Ghi rõ các ràng buộc kỹ thuật ở mức cao

Đầu ra là markdown theo quy tắc:
- Có front matter `title` và `weight`
- Có đúng một H1
- Có các mục:
  1. Purpose
  2. Feature summary
  3. User story
  4. Acceptance criteria
  5. User flow
  6. Technical notes

Áp dụng rule chung:

{{markdown_rules}}

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/spec.md"
- content = nội dung markdown bạn vừa viết

---

**discovery.md:**

{{discovery}}
