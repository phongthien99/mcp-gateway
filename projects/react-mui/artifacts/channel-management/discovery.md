---
title: Channel Management Discovery
weight: 10
---

# Channel Management Discovery

## Purpose

Phân tích yêu cầu xây dựng tính năng quản lý channel trong ứng dụng `react-mui`, xác định mục tiêu, phạm vi và các câu hỏi cần làm rõ trước khi thiết kế kỹ thuật.

## Goal

- Cung cấp giao diện cho admin quản lý danh sách channel trong hệ thống.
- Cho phép tạo mới, chỉnh sửa thông tin, kích hoạt/vô hiệu hoá và xoá channel.
- Hỗ trợ gán thành viên (users) vào channel và phân quyền theo vai trò trong channel.
- Tích hợp với module `user-management` đã có để chọn người dùng khi thêm vào channel.

## Assumptions

- `Assumption:` Chỉ người dùng có role `admin` mới có quyền quản lý toàn bộ channel; người dùng thường chỉ xem và tham gia channel được cấp phép.
- `Assumption:` Mỗi channel có ít nhất các trường: `id`, `name`, `description`, `status` (`active` / `inactive`), `createdAt`, `memberCount`.
- `Assumption:` Backend đã cung cấp RESTful API cho các thao tác CRUD channel và quản lý thành viên.
- `Assumption:` Dữ liệu channel lưu phía backend — không dùng `localStorage`.
- `Assumption:` Giao diện xây dựng bằng React + MUI theo chuẩn hiện tại của dự án.
- `Assumption:` Sử dụng `@tanstack/react-query` để quản lý server state, nhất quán với `user-management` và `pdf-notes`.
- `Assumption:` Danh sách channel đủ nhỏ để phân trang client-side trong giai đoạn đầu.
- `Decision:` Xoá channel không tự động xoá dữ liệu liên quan (tin nhắn, file, ghi chú) — dữ liệu được giữ nguyên hoặc archive.
- `Decision:` Mỗi channel có thể có nhiều admin-channel (owner) bên cạnh các member thông thường.

## Scope

- Trang danh sách channel (`/channels`) với bảng MUI:
  - Hiển thị: tên, mô tả, số thành viên, trạng thái, ngày tạo.
  - Phân trang client-side.
  - Tìm kiếm theo tên channel.
  - Lọc theo trạng thái (`active` / `inactive`).
- Tạo channel mới qua dialog/form (`name`, `description`, `status`).
- Chỉnh sửa thông tin channel (`name`, `description`, `status`).
- Kích hoạt / vô hiệu hoá channel (toggle `status`).
- Xoá channel với dialog xác nhận (không cascade dữ liệu liên quan).
- Quản lý thành viên channel:
  - Xem danh sách thành viên trong channel.
  - Thêm user vào channel (chọn từ danh sách user — tích hợp `user-management`).
  - Xoá thành viên khỏi channel.
  - Gán role thành viên trong channel (`owner` / `member`).
- Phân quyền route: chỉ `admin` truy cập được `/channels`.
- Hiển thị toast / snackbar thông báo kết quả thao tác.

## Out of scope

- Giao diện chat / nhắn tin trong channel.
- Tìm kiếm nội dung bên trong channel (tin nhắn, file).
- Thông báo real-time (WebSocket / SSE) khi có thay đổi thành viên.
- Import / export danh sách channel qua CSV/JSON.
- Phân quyền chi tiết theo từng tài nguyên trong channel (fine-grained RBAC).
- Lịch sử hoạt động / audit log của channel.
- Channel công khai tự đăng ký (self-join public channel) cho end-user.
- Bulk actions (kích hoạt / xoá nhiều channel cùng lúc).
- Backend API implementation.
- Multi-tenant / workspace isolation.

## Open questions

1. `Open question:` Channel có phân loại (`public` / `private`) không? Nếu có, người dùng thường có thể tự request tham gia channel public không?
2. `Open question:` Có cần trang chi tiết riêng (`/channels/:id`) với tab Members, Settings hay chỉ dùng dialog inline là đủ?
3. `Open question:` Số lượng thành viên tối đa mỗi channel là bao nhiêu? Cần phân trang server-side cho danh sách thành viên không?
4. `Open question:` Khi xoá channel, dữ liệu liên quan (tin nhắn, file, ghi chú PDF) xử lý như thế nào — archive hay giữ nguyên nhưng orphan?
5. `Open question:` Có cần bulk actions (thêm/xoá nhiều thành viên cùng lúc) trong phạm vi này không?
