import { requests } from "@/api";
import User from "@/models/user/user";
import UserContract from "@/models/user/contract";
import { UserTable } from "@/models/user/interface";

export default class AdminUser {

  toTableRow(): UserTable {
    const row = Object.create(this);
    row.actions = [];
    return row as UserTable;
  }

  static async fetchAll(nsid: number): Promise<Array<UserTable>> {
    const url = "/admin/nsadmin/nsall";
    const data = new Array<UserTable>();
    try {
      const payload = { namespace_id: nsid }
      const resp = await requests.post<Array<UserContract>>(url, payload);
      console.log("RESP", JSON.stringify(resp, null, "  "))
      resp.forEach((row) => data.push(new User(row).toTableRow()));
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return data;
  }
}
