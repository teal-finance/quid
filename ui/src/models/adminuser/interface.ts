interface AdminUserContract {
  "id": number;
  "usr_id": number;
  "ns_id": number;
  "username": string
}

interface AdminUserTable {
  id: number;
  userName: string;
  usrId: number;
}

export { AdminUserContract, AdminUserTable }
