---
title: Book Download Discovery
weight: 10
---

# Book Download Discovery

## Purpose

Phân tích yêu cầu xây dựng tính năng tải sách (download book) trong ứng dụng `react-mui`, xác định mục tiêu, phạm vi và các câu hỏi cần làm rõ trước khi thiết kế kỹ thuật.

## Goal

- Cho phép người dùng tải file sách (PDF hoặc các định dạng khác) về máy trực tiếp từ giao diện danh sách sách hoặc trang chi tiết sách.
- Ghi nhận lượt tải để hỗ trợ thống kê (nếu có yêu cầu).
- Kiểm soát quyền tải sách theo role người dùng (nếu cần).

## Assumptions

- `Assumption:` Module `books` đã tồn tại trong dự án (dựa trên commit `feat(books): implement book deletion functionality`); tính năng download là bổ sung vào module này.
- `Assumption:` File sách được lưu trên server/cloud storage (S3, GCS, hoặc tương đương) và backend cung cấp URL hoặc endpoint để tải.
- `Assumption:` Định dạng file chủ yếu là PDF; có thể có EPUB hoặc các định dạng khác.
- `Assumption:` Backend trả về signed URL (hoặc stream file trực tiếp) để tránh expose URL tĩnh công khai.
- `Assumption:` Giao diện xây dựng bằng React + MUI nhất quán với các module hiện có.
- `Assumption:` Sử dụng `@tanstack/react-query` để gọi API lấy download URL.
- `Assumption:` Người dùng đã đăng nhập mới được phép tải sách; download không hỗ trợ public/anonymous.
- `Decision:` Tải file bằng cách tạo thẻ `<a download>` với URL nhận từ backend — không stream qua frontend proxy.

## Scope

- Nút **Download** trên:
  - Mỗi hàng trong bảng danh sách sách (`/books`).
  - Trang/dialog chi tiết sách (nếu tồn tại).
- Gọi API backend để lấy download URL (có thể là signed URL với thời hạn ngắn).
- Kích hoạt tải file tự động sau khi nhận URL (dùng `<a download>` hoặc `window.open`).
- Hiển thị trạng thái loading trên nút khi đang lấy URL.
- Hiển thị toast/snackbar thông báo lỗi nếu lấy URL thất bại.
- Phân quyền: chỉ user đã đăng nhập mới thấy và dùng được nút Download.

## Out of scope

- Tải nhiều sách cùng lúc (bulk download / zip).
- Xem trước sách (PDF viewer) trong trình duyệt.
- Upload sách mới hoặc cập nhật file sách.
- Quản lý phiên bản file sách (versioning).
- Thống kê / báo cáo lượt tải chi tiết theo user.
- Giới hạn số lượt tải mỗi user (download quota).
- DRM (Digital Rights Management) hoặc watermarking file.
- Tải sách ở chế độ offline (service worker / cache).
- Backend API implementation.

## Open questions

1. `Open question:` Backend trả về **signed URL** (redirect) hay **stream file** trực tiếp? Cách này ảnh hưởng đến cách frontend xử lý download.
2. `Open question:` Có cần kiểm soát quyền tải theo role không (ví dụ: chỉ `premium` hoặc `admin` mới tải được)?
3. `Open question:` Tên file khi tải về được xác định như thế nào — theo `title` của sách hay tên file gốc trên storage?
4. `Open question:` Có ghi nhận lượt tải (analytics) không? Nếu có, ghi ở frontend hay backend tự log khi cấp URL?
5. `Open question:` Sách có thể có nhiều định dạng (PDF, EPUB)? Nếu có, nút Download cần cho phép chọn định dạng không?
