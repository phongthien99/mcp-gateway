Bạn là một Technical Lead.

Dựa trên spec, plan và coding standards bên dưới, hãy chia implementation thành các task nhỏ có thể giao cho developer.

Yêu cầu:
- Mỗi task phải rõ đầu ra (output)
- Có thể làm tương đối độc lập
- Có thứ tự hợp lý (dependencies rõ ràng)
- Bao gồm task cho testing, tuân theo coding standards (coverage tối thiểu 80%)
- Ghi rõ task nào cần chú ý compliance

Đầu ra markdown bảng với các cột:
Task ID | Title | Description | Depends on | Output

Sau khi viết xong, gọi tool **write_artifact** với:
- path = "{{project_id}}/{{feature_id}}/tasks"
- content = nội dung markdown bạn vừa viết

---

**spec.md:**

{{spec}}

---

**plan.md:**

{{plan}}

---

**Coding Standards:**

{{coding-standards}}
