---
title: User Management Discovery
weight: 10
---

# User Management Discovery

## Purpose

Phân tích yêu cầu xây dựng tính năng quản lý người dùng trong ứng dụng `react-mui`, xác định mục tiêu, phạm vi và các câu hỏi cần làm rõ trước khi thiết kế kỹ thuật.

## Goal

- Cung cấp giao diện cho admin quản lý danh sách người dùng trong hệ thống.
- Cho phép tạo mới, chỉnh sửa thông tin, khoá/mở khoá và xoá tài khoản người dùng.
- Hỗ trợ phân quyền theo vai trò với tối đa 50 role khác nhau.
- Tích hợp với tầng xác thực đã có trong tính năng `login-tool`.

## Assumptions

- `Assumption:` Chỉ người dùng có role `admin` mới có quyền truy cập module quản lý người dùng; protected route dựa trên token/role từ `login-tool`.
- `Assumption:` Dữ liệu người dùng được lưu phía backend — không dùng `localStorage` như `pdf-notes`.
- `Assumption:` Backend đã cung cấp RESTful API cho các thao tác CRUD người dùng.
- `Assumption:` Mỗi người dùng có ít nhất các trường: `id`, `email`, `name`, `role`, `status` (`active` / `inactive`), `createdAt`.
- `Assumption:` Giao diện xây dựng bằng React + MUI theo chuẩn hiện tại của dự án (`react-mui`).
- `Assumption:` Sử dụng `@tanstack/react-query` (đang dùng cho `pdf-notes`) để quản lý server state.
- `Assumption:` Danh sách người dùng đủ nhỏ để phân trang client-side (toàn bộ dữ liệu tải một lần).
- `Decision:` Khi tạo user mới, hệ thống gửi email mời — admin không đặt mật khẩu thay.
- `Decision:` Xoá user không kéo theo dữ liệu liên quan (ghi chú PDF, v.v. giữ nguyên).
- `Decision:` Phân trang client-side.
- `Decision:` Hệ thống có 50 role — cần dropdown có search khi chọn role.

## Scope

- Trang danh sách người dùng (`/users`) với bảng MUI:
  - Hiển thị: tên, email, vai trò, trạng thái, ngày tạo.
  - Phân trang client-side.
  - Tìm kiếm theo email / tên.
  - Lọc theo role (dropdown searchable, 50 options) và status.
- Tạo người dùng mới qua dialog/form (email, name, role) — gửi email mời tự động.
- Chỉnh sửa thông tin người dùng (name, role, status).
- Khoá / mở khoá tài khoản (toggle `status`).
- Xoá người dùng với dialog xác nhận (không cascade dữ liệu liên quan).
- Phân quyền route: chỉ `admin` truy cập được `/users`.
- Hiển thị toast / snackbar thông báo kết quả thao tác.

## Out of scope

- Reset mật khẩu người dùng từ giao diện admin.
- Quản lý quyền chi tiết theo từng resource (fine-grained RBAC).
- Lịch sử hoạt động / audit log của người dùng.
- Import / export danh sách người dùng qua CSV/JSON.
- Đăng nhập thay mặt người dùng (impersonation).
- Xác thực hai yếu tố (2FA).
- Backend API implementation.
- Quản lý nhóm / tổ chức (multi-tenant).
- Bulk actions (khoá nhiều user cùng lúc).

## Open questions

1. `Open question:` Bulk actions (khoá / xoá nhiều user cùng lúc) có cần trong phạm vi này không?
2. `Open question:` Có cần trang chi tiết riêng (`/users/:id`) hay chỉ dialog inline là đủ?
