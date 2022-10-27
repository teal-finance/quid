interface AdminUserContract {
  "id": number;
  "usr_id": number;
  "ns_id": number;
  "name": string
}

interface AdminUserTable {
  id: number;
  name: string;
  usrId: number;
}

export { AdminUserContract, AdminUserTable }
