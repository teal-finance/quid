import { api } from "@/api";
import { notify } from "@/state";
import { AdminUserContract, AdminUserTable } from "./interface";

export default class AdminUser {
  id: number;
  userName: string;
  usrId: number;

  constructor(row: AdminUserContract) {
    this.id = row.id;
    this.usrId = row.usr_id;
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
      const payload = { ns_id: nsid }
      try {
        const resp = await api.post<Array<AdminUserContract>>(url, payload);
        resp.forEach((row) => data.push(new AdminUser(row).toTableRow()));
      } catch (e) {
        console.log("QERR", JSON.stringify(e, null, "  "))
      }
    } catch (e) {
      console.log("Err", e);

      throw e;
    }
    return data;
  }

  static async searchNonAdmins(nsid: number, username: string): Promise<Array<AdminUser>> {
    const url = "/admin/nsadmin/nonadmins";
    const data = new Array<AdminUser>();
    try {
      const payload = { ns_id: nsid, username: username }
      const resp = await api.post<{ users: Array<AdminUserContract> }>(url, payload);
      resp.users.forEach((row) => data.push(new AdminUser(row)));
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return data;
  }

  static async fetchAdd(nsId: number, usrIds: Array<number>) {
    try {
      await api.post("/admin/nsadmin/add", {
        ns_id: nsId,
        usr_ids: usrIds,
      });
    } catch (e) {
      console.log(e)
      notify.error("Error adding admin users")
    }
  }

  static async delete(uid: number, nsid: number) {
    console.log("Delete", {
      ns_id: nsid,
      usr_id: uid,
    })
    try {
      await api.post("/admin/nsadmin/delete", {
        ns_id: nsid,
        usr_id: uid,
      });
    } catch (e) {
      console.log(e)
      notify.error("Error deleting admin users")
    }
  }
}
