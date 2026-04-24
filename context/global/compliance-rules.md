# Compliance Rules

## Data export
- Mọi file export phải có audit log: ai export, khi nào, bao nhiêu records
- Không export dữ liệu nhạy cảm (password, token, PII) ra file
- File export phải được xóa khỏi server sau khi user tải về (hoặc không lưu server)

## API security
- Rate limit: tối đa 100 request/phút per user
- Endpoint export phải yêu cầu authentication
- Log tất cả request đến endpoint export

## PII (Personally Identifiable Information)
- Không log email, phone, tên đầy đủ
- Mask data trong response nếu user không có quyền xem đầy đủ
