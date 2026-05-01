---
title: "Compliance Rules"
weight: 30
---

# Compliance Rules

## Scope

Tool nội bộ, chạy local, không expose ra internet. Không có yêu cầu compliance bên ngoài.

## Rules

- **Không lưu credentials** vào `artifacts/` hoặc `context/` — các thư mục này được mount vào container và có thể được commit vào git.
- **Không expose MCP port ra ngoài localhost** trừ khi có chủ ý (kiểm tra `ports` trong `docker-compose.yml`).
- Artifact và context files được phục vụ qua Hugo dashboard mà không có auth — không đặt thông tin nhạy cảm vào đây.
- Không có yêu cầu audit log hay data retention.
