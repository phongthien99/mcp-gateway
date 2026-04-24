# Architecture Guidelines

## Stack
- Backend: Go (Gin framework)
- Frontend: React + TypeScript
- Database: PostgreSQL
- Cache: Redis

## Layer structure
```
handler → service → repository → database
```
- Handler: HTTP request/response only, no business logic
- Service: business logic, orchestration
- Repository: database queries only, no business logic

## API conventions
- RESTful, JSON response
- Error format: `{"error": "message", "code": "ERROR_CODE"}`
- Pagination: `?page=1&limit=20`, response includes `total`, `page`, `limit`
- Authentication: Bearer token via `Authorization` header

## File download endpoints
- Content-Type phải set đúng MIME type
- Luôn set `Content-Disposition: attachment; filename="..."`
- Với CSV: thêm UTF-8 BOM (`\xEF\xBB\xBF`) để Excel Windows đọc đúng

## Database
- Không dùng raw SQL, dùng query builder hoặc ORM
- Transaction cho các thao tác ghi nhiều bảng
- Index trên các cột filter/sort thường dùng
