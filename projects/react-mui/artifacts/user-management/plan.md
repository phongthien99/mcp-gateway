---
title: User Management Technical Plan
weight: 30
---

# User Management Technical Plan

## Purpose

Đặc tả kỹ thuật triển khai, thứ tự implement các file, và các risk/compliance cần xử lý cho tính năng quản lý người dùng trong `react-mui`.

## Architecture Impact

Tính năng này là **ngoại lệ đầu tiên** phá vỡ kiến trúc localStorage-only của dự án:

| Thay đổi | Trước | Sau |
|---|---|---|
| Data source | `localStorage` | Backend REST API |
| Auth guard | Không có | Check session → redirect `/login` |
| PII handling | Không có | Email, name gửi qua network |
| Routing | `/`, `/books/:slug` | Thêm `/users` |

`Decision:` Repository layer của feature `users` gọi `fetch` thay vì đọc `localStorage`. Layering rules (schema → repository → hooks → facade → component) **giữ nguyên**.

`Decision:` Token xác thực dùng HTTP-only cookie — không lưu `localStorage` (tuân thủ compliance). Repository không tự gắn `Authorization` header; trình duyệt tự đính cookie vào mọi request cùng origin.

`Decision:` API base URL đọc từ `import.meta.env.VITE_API_BASE_URL`. Nếu biến không được set, fallback về chuỗi rỗng (relative URL, phù hợp khi frontend và API cùng origin).

`Assumption:` Hook file names follow codebase convention (camelCase, e.g. `useUsersQuery.ts`) — không phải kebab-case như coding-standards ghi, vì codebase hiện tại dùng camelCase cho hooks.

## Backend Changes

Backend API đã có sẵn (xác nhận từ discovery). Không có thay đổi backend trong scope này.

Các endpoint được dùng:

```
GET    /api/users
GET    /api/roles
POST   /api/users
PUT    /api/users/:id
PATCH  /api/users/:id/status   { status: 'active' | 'inactive' }
DELETE /api/users/:id
```

## Frontend Changes

### Thứ tự implement (theo dependency)

**Bước 1 — Schema** (không dependency)

- `src/features/users/schema/user.schema.ts`
  - Export: `UserSchema`, `UsersSchema`, `RoleSchema`, `RolesSchema`, `CreateUserSchema`, `UpdateUserSchema`, `UserStatusSchema`
  - Export types: `User`, `Role`, `CreateUserInput`, `UpdateUserInput`

**Bước 2 — Repository** (phụ thuộc schema)

- `src/features/users/repositories/user.repository.ts`
  - Plain object `userRepository` với 6 methods: `getAll`, `getRoles`, `create`, `update`, `toggleStatus`, `delete`
  - Mỗi method: `fetch` → kiểm tra `res.ok` → parse với Zod → throw nếu lỗi
  - Xử lý HTTP error: đọc `res.json()` để lấy message, throw `Error` với message đó
  - Không import React

**Bước 3 — Hooks** (phụ thuộc repository)

Tất cả 6 hooks có thể viết song song sau bước 2:

| File | Loại | invalidate |
|---|---|---|
| `useUsersQuery.ts` | `useQuery` | — |
| `useRolesQuery.ts` | `useQuery` | — |
| `useCreateUserMutation.ts` | `useMutation` | `['users']` |
| `useUpdateUserMutation.ts` | `useMutation` | `['users']` |
| `useToggleUserStatusMutation.ts` | `useMutation` | `['users']` |
| `useDeleteUserMutation.ts` | `useMutation` | `['users']` |

**Bước 4 — Facade types** (phụ thuộc schema)

- `src/features/users/facade/user-list.facade.types.ts`
  - Export interface `UserListFacade` (xem spec)

**Bước 5 — Facade** (phụ thuộc hooks + facade types)

- `src/features/users/facade/user-list.facade.ts`
  - Export `useUserListFacade(): UserListFacade`
  - Client-side filter pipeline: search → roleId → status → paginate
  - `isSubmitting` = OR của tất cả mutation `isPending`
  - `dialogError` reset về `null` khi mở dialog mới
  - `toggleStatus` gọi mutation với `status` đảo ngược, bắt lỗi vào snackbar (không dùng dialog)

**Bước 6 — Components** (phụ thuộc facade types + schema)

Ba components viết song song:

- `src/features/users/components/UserListTable.tsx`
  - Props: slice từ `UserListFacade` (users, roles, isLoading, page, pageSize, totalCount, onEdit, onToggleStatus, onDelete, onPageChange, onRowsPerPageChange)
  - Columns: Tên / Email / Vai trò / Trạng thái (`Chip`) / Ngày tạo / Hành động
  - `TablePagination` rowsPerPageOptions `[10, 25, 50]`

- `src/features/users/components/UserFormDialog.tsx`
  - Props: `open`, `mode`, `initialValues?`, `roles`, `onSubmit`, `onClose`, `isSubmitting`, `error`
  - Controlled form với `useState` nội bộ; validate bằng `CreateUserSchema` / `UpdateUserSchema` `.safeParse` khi submit
  - `Autocomplete` role: `options={roles}`, `getOptionLabel={r => r.name}`, `isOptionEqualToValue={(o,v) => o.id === v.id}`
  - Email field `disabled` khi `mode === 'edit'`

- `src/features/users/components/DeleteUserDialog.tsx`
  - Props: `open`, `target: User | null`, `onConfirm`, `onClose`, `isSubmitting`, `error`

**Bước 7 — Page** (phụ thuộc facade + components)

- `src/pages/user-list.tsx`
  - Gọi `useUserListFacade()`
  - Toolbar inline: `TextField` search + debounce 300ms (`setTimeout` + `clearTimeout`), `Autocomplete` role, `Select` status, Button "Thêm người dùng"
  - `Snackbar` + `Alert` cho success/error messages ngoài dialog
  - Auth guard: đọc `localStorage.getItem('react-mui.session')` (placeholder cho đến khi `login-tool` implement) — nếu null redirect `/login`

**Bước 8 — Routing** (phụ thuộc page)

- Sửa `src/app-shell.tsx`:
  - Thêm helper `isUsersUrl()` kiểm tra `pathname === '/users'`
  - Thêm nhánh render `UserListPage` trước `BookList`
  - Thêm link "Quản lý người dùng" vào navigation nếu có

### File thay đổi tổng hợp

```
src/features/users/          ← mới hoàn toàn (13 files)
src/pages/user-list.tsx      ← mới
src/app-shell.tsx            ← sửa (thêm /users route)
```

## Testing Plan

Dự án chưa có test runner setup — đây là test cases cần cover khi test được thêm vào:

### Repository

- `getAll()`: parse valid response → trả `User[]`; parse invalid → throw
- `create()`: `POST` đúng body; response 4xx → throw với message từ body
- `toggleStatus()`: gọi `PATCH /api/users/:id/status` với đúng payload
- `delete()`: response 204 → resolve void; response 4xx → throw

### Facade

- Filter pipeline: search match trên name/email, không match bỏ qua; filterRoleId/filterStatus null = bỏ qua bước lọc
- Pagination: page 0 trả 10 items đầu; page 1 trả 10 items tiếp theo
- `submitCreate` success → `createDialogOpen = false`, `dialogError = null`
- `submitCreate` error → `dialogError` có message, dialog vẫn mở
- `toggleStatus` không mở dialog; mutation pending block double-click

### Components

- `UserFormDialog` create mode: email enabled, note "Email mời" hiển thị
- `UserFormDialog` edit mode: email disabled, note ẩn
- `UserFormDialog` submit khi email rỗng: lỗi validation hiển thị
- `DeleteUserDialog`: hiển thị đúng name + email của target
- `UserListTable` empty state: render "Không tìm thấy người dùng nào."

## Risks & Compliance Notes

### Risk 1 — Auth guard placeholder

`src/pages/user-list.tsx` dùng `localStorage.getItem('react-mui.session')` làm guard tạm. Khi `login-tool` hoàn thành, cần thay bằng session hook thực. Nếu bỏ quên: bất kỳ user nào cũng truy cập được `/users`.

**Mitigation:** Ghi `// TODO(login-tool): replace with real auth check` ngay tại dòng guard.

### Risk 2 — PII qua network (compliance violation)

Compliance rule hiện tại cấm gửi dữ liệu ra ngoài thiết bị. Feature này gửi `email` và `name` lên server.

**Mitigation:** Đây là ngoại lệ được chấp nhận có chủ ý (user đã xác nhận API có sẵn). Ghi `Decision:` vào spec (đã có). Không log `email`/`name` trong catch blocks.

### Risk 3 — Cookie không tự động gửi cross-origin

Nếu frontend và API khác origin, HTTP-only cookie sẽ bị chặn bởi `SameSite` policy.

**Mitigation:** Repository dùng `credentials: 'include'` trong tất cả `fetch` calls. Confirm CORS config phía backend cho phép origin của frontend.

### Risk 4 — Zod v4 API khác v3

Dự án dùng Zod v4 (`zod` ≥ 4.x). Một số API thay đổi so v3 (ví dụ `.email()` không nhận options object, dùng `.email()` rồi `.superRefine()` nếu cần custom message).

**Mitigation:** Test schema parse ngay sau khi viết xong `user.schema.ts` trước khi viết repository.

### Risk 5 — Routing conflict với `app-shell.tsx`

`app-shell.tsx` hiện dùng `window.location.pathname` thủ công. Thêm `/users` sai vị trí có thể khiến BookList render lên trên.

**Mitigation:** Thêm check `/users` **trước** check `/books/` trong render logic. Order: `isUsersUrl()` → `isBookUrl()` → default BookList.
