---
title: User Management Spec
weight: 20
---

# User Management Spec

## Purpose

Đặc tả chi tiết kỹ thuật cho tính năng quản lý người dùng, đủ để engineer implement mà không cần hỏi thêm.

## Architectural Note

> **Xung đột với kiến trúc hiện tại:** `architecture.md` và `compliance-rules.md` ghi rõ _"không có backend — toàn bộ dữ liệu lưu localStorage"_. Tính năng này **phá vỡ rule đó** vì yêu cầu backend API (đã xác nhận có sẵn) và xử lý PII (email, name). Các điểm lệch cần chú ý:
>
> - Không dùng `localStorage` cho user data — dữ liệu lấy từ API.
> - Token xác thực **không lưu `localStorage`** (compliance rule). Cơ chế lưu token cần được quyết định riêng (HTTP-only cookie khuyến nghị).
> - Cần cân nhắc scrub PII trước khi log lỗi.

## Data Model

### Zod Schema — `src/features/users/schema/user.schema.ts`

```ts
import { z } from 'zod';

export const UserStatusSchema = z.enum(['active', 'inactive']);

export const UserSchema = z.object({
  id: z.string(),
  email: z.string().email(),
  name: z.string(),
  roleId: z.string(),
  roleName: z.string(),
  status: UserStatusSchema,
  createdAt: z.string(),
});

export const UsersSchema = z.array(UserSchema);

export const RoleSchema = z.object({
  id: z.string(),
  name: z.string(),
});

export const RolesSchema = z.array(RoleSchema);

export const CreateUserSchema = z.object({
  email: z.string().email('Email không hợp lệ'),
  name: z.string().min(1, 'Tên không được để trống'),
  roleId: z.string().min(1, 'Vai trò không được để trống'),
});

export const UpdateUserSchema = z.object({
  name: z.string().min(1, 'Tên không được để trống'),
  roleId: z.string().min(1, 'Vai trò không được để trống'),
});

export type User = z.infer<typeof UserSchema>;
export type Role = z.infer<typeof RoleSchema>;
export type CreateUserInput = z.infer<typeof CreateUserSchema>;
export type UpdateUserInput = z.infer<typeof UpdateUserSchema>;
```

## API Contract

Base URL và auth header do caller cung cấp qua HTTP client wrapper (ngoài scope spec này).

| Method | Path | Body | Response |
|---|---|---|---|
| `GET` | `/api/users` | — | `User[]` |
| `GET` | `/api/roles` | — | `Role[]` |
| `POST` | `/api/users` | `CreateUserInput` | `User` (gửi email mời tự động) |
| `PUT` | `/api/users/:id` | `UpdateUserInput` | `User` |
| `PATCH` | `/api/users/:id/status` | `{ status: 'active' \| 'inactive' }` | `User` |
| `DELETE` | `/api/users/:id` | — | `204 No Content` |

## File Structure

```
src/
  features/
    users/
      schema/
        user.schema.ts
      repositories/
        user.repository.ts
      hooks/
        useUsersQuery.ts
        useRolesQuery.ts
        useCreateUserMutation.ts
        useUpdateUserMutation.ts
        useToggleUserStatusMutation.ts
        useDeleteUserMutation.ts
      facade/
        user-list.facade.ts
        user-list.facade.types.ts
      components/
        UserListTable.tsx
        UserFormDialog.tsx
        DeleteUserDialog.tsx
  pages/
    user-list.tsx
```

## Repository — `user.repository.ts`

```ts
export const userRepository = {
  getAll(): Promise<User[]>,
  getRoles(): Promise<Role[]>,
  create(input: CreateUserInput): Promise<User>,
  update(id: string, input: UpdateUserInput): Promise<User>,
  toggleStatus(id: string, status: 'active' | 'inactive'): Promise<User>,
  delete(id: string): Promise<void>,
};
```

Parse response với Zod. Nếu parse thất bại throw `Error`.

## Hooks

| File | Query key | Gọi |
|---|---|---|
| `useUsersQuery.ts` | `['users']` | `userRepository.getAll()` |
| `useRolesQuery.ts` | `['roles']` | `userRepository.getRoles()` |
| `useCreateUserMutation.ts` | invalidate `['users']` on success | `userRepository.create(input)` |
| `useUpdateUserMutation.ts` | invalidate `['users']` on success | `userRepository.update(id, input)` |
| `useToggleUserStatusMutation.ts` | invalidate `['users']` on success | `userRepository.toggleStatus(id, status)` |
| `useDeleteUserMutation.ts` | invalidate `['users']` on success | `userRepository.delete(id)` |

## Facade — `user-list.facade.ts`

```ts
interface UserListFacade {
  users: User[];        // filtered + paginated slice
  roles: Role[];
  totalCount: number;
  isLoading: boolean;
  error: string | null;

  search: string;
  filterRoleId: string | null;
  filterStatus: 'active' | 'inactive' | null;
  page: number;
  pageSize: number;

  createDialogOpen: boolean;
  editTarget: User | null;
  deleteTarget: User | null;
  dialogError: string | null;
  isSubmitting: boolean;

  setSearch(v: string): void;
  setFilterRoleId(v: string | null): void;
  setFilterStatus(v: 'active' | 'inactive' | null): void;
  setPage(n: number): void;
  openCreateDialog(): void;
  openEditDialog(user: User): void;
  openDeleteDialog(user: User): void;
  closeDialogs(): void;
  submitCreate(input: CreateUserInput): Promise<void>;
  submitUpdate(input: UpdateUserInput): Promise<void>;
  toggleStatus(user: User): Promise<void>;
  submitDelete(): Promise<void>;
}
```

Filtering pipeline (client-side):
1. Match `search` trên `name` hoặc `email` (case-insensitive).
2. Lọc `filterRoleId` nếu không null.
3. Lọc `filterStatus` nếu không null.
4. Slice theo `page` × `pageSize` (default `pageSize = 10`).

## Components

### `UserListTable.tsx`

- MUI `Table`: columns **Tên**, **Email**, **Vai trò**, **Trạng thái**, **Ngày tạo**, **Hành động**.
- Row actions: edit (`EditOutlined`), toggle status (`LockOutlined` / `LockOpenOutlined`), xoá (`DeleteOutline`).
- Toggle status gọi `toggleStatus(user)` trực tiếp — không cần dialog xác nhận.
- Empty state: "Không tìm thấy người dùng nào."
- `TablePagination` (rows per page: 10 / 25 / 50).
- Loading: `CircularProgress` centered.

### Toolbar (inline trong `user-list.tsx`)

- `TextField` search (debounce 300 ms).
- `Autocomplete` lọc role (searchable, 50 options).
- `Select` lọc status: Tất cả / Đang hoạt động / Đã khoá.
- Button **"Thêm người dùng"**.

### `UserFormDialog.tsx`

Props: `open`, `mode: 'create' | 'edit'`, `initialValues?`, `roles`, `onSubmit`, `onClose`, `isSubmitting`, `error`.

Fields:
- **Email** — disabled khi `mode === 'edit'`.
- **Tên** — `TextField`.
- **Vai trò** — `Autocomplete` searchable.

Footer: **Huỷ** + **Lưu**. Khi `mode === 'create'`: note _"Email mời sẽ được gửi tới địa chỉ trên."_

### `DeleteUserDialog.tsx`

- Title: "Xoá người dùng?", body: _"Bạn sắp xoá **{name}** ({email}). Hành động này không thể hoàn tác."_
- Nút **Xoá** color `error`. Hiển thị `dialogError` nếu có.

## Page — `src/pages/user-list.tsx`

- Route: `/users`, protected (admin only → redirect `/login`).
- Snackbar thành công sau mỗi action.

## Error Handling

| Tình huống | Xử lý |
|---|---|
| `401` | Redirect `/login` |
| `403` | Snackbar "Bạn không có quyền thực hiện thao tác này." |
| `4xx` khác | Message từ response body trong `dialogError` |
| `5xx` | "Lỗi hệ thống, vui lòng thử lại." |
| Zod parse fail | Throw `Error`, TanStack Query bắt qua `isError` |
