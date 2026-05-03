---
title: Compliance Rules
weight: 30
---

# Compliance Rules

## Purpose

Ràng buộc về bảo mật, data privacy, và logging để đảm bảo ứng dụng không rò rỉ dữ liệu người dùng.

## Scope

Áp dụng cho mọi code đọc/ghi dữ liệu, gọi API bên ngoài, hoặc render nội dung người dùng.

## Data Privacy

- Decision: **không có backend**. Toàn bộ dữ liệu (sách, ghi chú, nét vẽ) lưu trong `localStorage` của trình duyệt và không bao giờ rời khỏi thiết bị.
- Không gửi nội dung PDF, ghi chú, hay văn bản đã chọn lên bất kỳ server bên thứ ba nào.
- Translation: sử dụng Chrome Translator API (on-device). Assumption: API không gửi text ra ngoài; nếu API unavailable thì fail gracefully — không fallback sang external API mà không có sự đồng ý rõ ràng của người dùng.
- Không thu thập analytics, không nhúng tracking pixel, không dùng cookies.

## localStorage Security

- Mọi key phải có prefix `react-mui.` để tránh xung đột với extension/tab khác.
- Dữ liệu đọc từ `localStorage` phải được validate bằng Zod trước khi dùng; nếu parse thất bại trả về mảng rỗng / safe default — không throw uncaught error.
- Không lưu thông tin nhạy cảm (credentials, token, PII) vào `localStorage`.

## XSS Prevention

- Không dùng `dangerouslySetInnerHTML` trừ khi content đã được sanitize.
- `react-markdown` render ghi chú người dùng — không cho phép raw HTML trong markdown input.
- URL PDF từ người dùng (`pdfUrl`) chỉ được truyền vào `<Document file={...}>` của react-pdf — không dùng làm `href` hay `src` trong thẻ `<a>`/`<img>` mà không validate.

## Content Security

- Worker của pdfjs load từ `unpkg.com` — đây là CDN đã được cố định. Không thay đổi sang URL tùy ý.
- Không eval, không `new Function()`, không dynamic `import()` với string từ user input.

## Logging

- Không log dữ liệu người dùng (nội dung ghi chú, văn bản PDF) ra console trong production build.
- `catch {}` blocks không log errors có chứa content người dùng.
- Open question: chưa có error reporting (Sentry, etc.) — nếu thêm vào phải scrub PII trước khi gửi.

## Dependency Policy

- Không thêm dependency mới mà không kiểm tra license (project license: MIT).
- Dependency mới phải có bundle size impact được cân nhắc (project là SPA, no SSR).
