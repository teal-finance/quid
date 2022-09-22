import { api } from "@/api";
import UserContract from "./contract";
import { UserTable } from "./interface";
import { user } from "@/state";

export default class User {
  id: number;
  name: string;

  constructor(user: UserContract) {
    this.id = user.id;
    this.name = user.name;
  }

  toTableRow(): UserTable {
    const row = Object.create(this);
    row.actions = [];
    return row as UserTable;
  }

  static async fetchAll(nsid: number): Promise<Array<UserTable>> {
    const url = user.adminUrl + "/users/nsall";
    const data = new Array<UserTable>();
    try {
      const payload = { ns_id: nsid }
      const resp = await api.post<Array<UserContract>>(url, payload);
      resp.forEach((row) => data.push(new User(row).toTableRow()));
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return data;
  }

  static async delete(id: number) {
    const url = user.adminUrl + "/users/delete";
    await api.post(url, {
      id: id,
      ns_id: user.namespace.value.id
    });
  }
}
