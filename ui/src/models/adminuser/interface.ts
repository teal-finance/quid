interface AdminUserContract {
  "id": number;
  "user_id": number;
  "namespace_id": number;
  "username": string
}

interface AdminUserTable {
  id: number;
  userName: string;
  userId: number;
}

export { AdminUserContract, AdminUserTable }