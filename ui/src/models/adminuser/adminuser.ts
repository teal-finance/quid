import { requests } from "@/api";
import { notify } from "@/state";
import { AdminUserContract, AdminUserTable } from "./interface";

export default class AdminUser {
  id: number;
  userName: string;
  userId: number;

  constructor(row: AdminUserContract) {
    this.id = row.id;
    this.userId = row.user_id;
    this.userName = row.username;
  }

  toTableRow(): AdminUserTable {
    const row = Object.create(this);
    row.actions = [];
    return row as AdminUserTable;
  }

  static async fetchAll(nsid: number): Promise<Array<AdminUserTable>> {
    const url = "/admin/nsadmin/nsall";
    const data = new Array<AdminUserTable>();
    try {
      const payload = { namespace_id: nsid }
      const resp = await requests.post<Array<AdminUserContract>>(url, payload);
      console.log("RESP", JSON.stringify(resp, null, "  "))
      resp.forEach((row) => data.push(new AdminUser(row).toTableRow()));
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return data;
  }

  static async fetchAdd(nsId: number, userIds: Array<number>) {
    try {
      await requests.post("/admin/nsadmin/add", {
        namespace_id: nsId,
        user_ids: userIds,
      });
    } catch (e) {
      console.log(e)
      notify.error("Error adding admin users")
    }
  }

  static async delete(uid: number, nsid: number) {
    console.log("Delete", {
      namespace_id: nsid,
      user_id: uid,
    })
    try {
      await requests.post("/admin/nsadmin/delete", {
        namespace_id: nsid,
        user_id: uid,
      });
    } catch (e) {
      console.log(e)
      notify.error("Error deleting admin users")
    }
  }
}