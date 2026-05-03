---
title: User Management Tasks
weight: 40
---

# User Management Tasks

## Purpose

Danh sách task implementation cho tính năng quản lý người dùng, có thể giao cho developer theo thứ tự dependency.

## Tasks

| Task ID | Title | Description | Depends on | Output |
|---|---|---|---|---|
| U-01 | Schema layer | Tạo `src/features/users/schema/user.schema.ts` với đầy đủ Zod schemas: `UserStatusSchema`, `UserSchema`, `UsersSchema`, `RoleSchema`, `RolesSchema`, `CreateUserSchema`, `UpdateUserSchema`. Export inferred types. Kiểm tra parse một object hợp lệ và một object thiếu field để xác nhận schema đúng. | — | `user.schema.ts` build thành công, TypeScript strict pass |
| U-02 | Repository layer | Tạo `src/features/users/repositories/user.repository.ts`. Implement 6 methods dùng `fetch` + `credentials: 'include'`. Base URL từ `import.meta.env.VITE_API_BASE_URL`. Parse response bằng Zod; nếu `!res.ok` đọc body JSON lấy message rồi throw. Không import React. | U-01 | `user.repository.ts` build thành công, lint pass |
| U-03 | Query hooks | Tạo `useUsersQuery.ts` (queryKey `['users']`) và `useRolesQuery.ts` (queryKey `['roles']`) gọi repository tương ứng. | U-02 | 2 hook files, build + lint pass |
| U-04 | Mutation hooks | Tạo 4 mutation hooks: `useCreateUserMutation.ts`, `useUpdateUserMutation.ts`, `useToggleUserStatusMutation.ts`, `useDeleteUserMutation.ts`. Mỗi hook invalidate `['users']` on success. | U-02 | 4 hook files, build + lint pass |
| U-05 | Facade types | Tạo `src/features/users/facade/user-list.facade.types.ts` export interface `UserListFacade` theo spec. | U-01 | `user-list.facade.types.ts`, build pass |
| U-06 | Facade implementation | Tạo `src/features/users/facade/user-list.facade.ts` export `useUserListFacade()`. Implement filter pipeline (search → roleId → status → paginate). `isSubmitting` = OR các mutation pending. `dialogError` reset khi mở dialog mới. `toggleStatus` đảo `status` rồi gọi mutation. | U-03, U-04, U-05 | `user-list.facade.ts`, build + lint pass |
| U-07 | UserListTable component | Tạo `src/features/users/components/UserListTable.tsx`. MUI `Table` với 6 columns. `Chip` cho status (color `success`/`default`). Row actions: edit, toggle status (icon `LockOutlined`/`LockOpenOutlined`), xoá. `TablePagination` rowsPerPageOptions `[10, 25, 50]`. Empty state + loading `CircularProgress`. | U-05 | `UserListTable.tsx`, build + lint pass |
| U-08 | UserFormDialog component | Tạo `src/features/users/components/UserFormDialog.tsx`. Controlled form với `useState`. Validate bằng `CreateUserSchema`/`UpdateUserSchema` `.safeParse` khi submit. `Autocomplete` role searchable (50 options). Email disabled khi `mode === 'edit'`. Note "Email mời sẽ được gửi" hiển thị khi `mode === 'create'`. | U-01 | `UserFormDialog.tsx`, build + lint pass |
| U-09 | DeleteUserDialog component | Tạo `src/features/users/components/DeleteUserDialog.tsx`. Hiển thị name + email của `target`. Nút **Xoá** color `error`, loading khi `isSubmitting`. Hiển thị `error` nếu có. | U-01 | `DeleteUserDialog.tsx`, build + lint pass |
| U-10 | Page + toolbar | Tạo `src/pages/user-list.tsx`. Gọi `useUserListFacade()`. Toolbar inline: `TextField` search debounce 300ms, `Autocomplete` role, `Select` status, Button "Thêm người dùng". Render `UserListTable` + `UserFormDialog` + `DeleteUserDialog`. `Snackbar`+`Alert` success/error. Auth guard: `localStorage.getItem('react-mui.session')` → redirect `/login` nếu null; ghi `// TODO(login-tool): replace with real auth check`. ⚠️ **Compliance:** không log email/name trong catch blocks. | U-06, U-07, U-08, U-09 | `user-list.tsx`, build + lint pass |
| U-11 | Routing (app-shell) | Sửa `src/app-shell.tsx`: thêm `isUsersUrl()` check `pathname === '/users'`; thêm nhánh render `UserListPage` **trước** `isBookUrl()` để tránh conflict. Import `UserListPage` từ `src/pages/user-list.tsx`. | U-10 | `app-shell.tsx` updated, route `/users` hoạt động |
| U-12 | Schema unit tests | Viết tests cho `user.schema.ts`: parse valid `User` object ✓; parse `User` thiếu `email` ✗; parse `CreateUserInput` email sai format ✗; parse `CreateUserInput` name rỗng ✗; `UserStatusSchema` reject giá trị ngoài enum ✗. Coverage schema layer ≥ 80%. | U-01 | Test file `user.schema.test.ts`, tất cả pass |
| U-13 | Facade unit tests | Viết tests cho filter pipeline trong facade (mock hooks): search match name/email case-insensitive; filterRoleId null = không lọc; filterStatus null = không lọc; pagination đúng slice; `submitCreate` success → `createDialogOpen = false`; `submitCreate` error → `dialogError` có message. Coverage facade ≥ 80%. | U-06 | Test file `user-list.facade.test.ts`, tất cả pass |
| U-14 | Component smoke tests | Viết tests cho 3 components: `UserFormDialog` create mode — email enabled, note hiển thị; `UserFormDialog` edit mode — email disabled; `UserFormDialog` submit email rỗng — validation error hiển thị; `DeleteUserDialog` render đúng name+email; `UserListTable` empty state render "Không tìm thấy người dùng nào.". Coverage components ≥ 80%. | U-07, U-08, U-09 | Test files cho 3 components, tất cả pass |
| U-15 | E2E smoke (manual) | Kiểm tra thủ công trên browser: (1) truy cập `/users` khi chưa auth → redirect `/login`; (2) render danh sách user từ API; (3) tạo user → dialog đóng, snackbar "Đã tạo người dùng"; (4) toggle status → chip đổi màu; (5) xoá user → row biến mất; (6) search/filter hoạt động client-side. | U-11 | Checklist manual test pass, không có console error |
