# Coding Standards

## Go
- Tên hàm, struct: PascalCase (exported), camelCase (unexported)
- Error phải được wrap: `fmt.Errorf("context: %w", err)`
- Không dùng `panic` trong production code
- Test file: `_test.go`, dùng `testify/assert`
- Mỗi package có một file `doc.go` mô tả package

## Frontend (React + TypeScript)
- Component: PascalCase, file `.tsx`
- Hook: `use` prefix, file `.ts`
- Không dùng `any` type
- Props phải có interface định nghĩa rõ ràng

## Git
- Branch: `feature/`, `fix/`, `chore/`
- Commit message: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`
- PR phải có description và test plan

## Testing
- Unit test coverage tối thiểu 80%
- Integration test cho mọi API endpoint
- Không mock database trong integration test
