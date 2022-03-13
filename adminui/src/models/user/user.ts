import { requests } from "@/api";
import UserContract from "./contract";
import { UserTable } from "./interface";

export default class User {
  id: number;
  name: string;

  constructor(data: UserContract) {
    this.id = data.id;
    this.name = data.username;
  }

  toTableRow(): UserTable {
    const row = Object.create(this);
    row.actions = [];
    return row as UserTable;
  }

  static async fetchAll(nsid: number): Promise<Array<UserTable>> {
    const url = "/admin/users/nsall";
    const data = new Array<UserTable>();
    try {
      const payload = { namespace_id: nsid }
      const resp = await requests.post<Array<UserContract>>(url, payload);
      resp.forEach((row) => data.push(new User(row).toTableRow()));
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return data;
  }

  static async delete(id: number) {
    await requests.post("/admin/users/delete", {
      id: id,
    });
  }
}