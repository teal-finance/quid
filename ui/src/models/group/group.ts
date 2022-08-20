import { api } from "@/api";
import { GroupContract } from "./contract";
import { GroupTable } from "./interface";
import { user } from "@/state";

export default class Group {
  id: number;
  name: string;

  constructor({ id, name }: { id: number, name: string }) {
    this.id = id;
    this.name = name;
  }

  // *************************
  //   factory constructors
  // *************************

  static fromContract(data: GroupContract): Group {
    return new Group({ id: data.id, name: data.name })
  }

  // *************************
  //         methods
  // *************************

  toTableRow(): GroupTable {
    const row: GroupTable = {
      id: this.id,
      name: this.name,
      actions: [],
    };
    return row;
  }

  // *************************
  //    static methods
  // *************************

  static async fetchAll(nsid: number): Promise<Array<GroupTable>> {
    const url = user.adminUrl + "/groups/nsall";
    const ns = new Array<GroupTable>();
    try {
      const payload = { namespace_id: nsid }
      const resp = await api.post<Array<GroupContract>>(url, payload, false, true);
      console.log("RESP", resp)
      resp.forEach((row) => {
        //console.log(row)
        ns.push(new Group(row).toTableRow())
      });
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return ns;
  }

  static async fetchUserGroups(uid: number) {
    const url = user.adminUrl + "/users/groups";
    const data = new Array<GroupTable>();
    try {
      const payload = { id: uid, namespace_id: user.namespace.value.id }
      const resp = await api.post<{ groups: Array<GroupContract> }>(url, payload);
      //console.log("RESP", JSON.stringify(resp.groups, null, "  "))
      if (resp.groups.length > 0) {
        resp.groups.forEach((row) => {
          data.push(new Group(row).toTableRow())
        });
      }
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return data;
  }

  static async addUserToGroup(uid: number, gid: number) {
    const url = user.adminUrl + "/groups/add_user";
    try {
      const payload = {
        user_id: uid,
        group_id: gid,
        namespace_id: user.namespace.value.id
      }
      await api.post(url, payload);
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
  }

  static async removeUserFromGroup(uid: number, gid: number) {
    const url = user.adminUrl + "/groups/remove_user";
    try {
      const payload = {
        user_id: uid,
        group_id: gid,
        namespace_id: user.namespace.value.id
      }
      await api.post(url, payload);
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
  }

  static async delete(id: number) {
    await api.post(user.adminUrl + "/groups/delete", {
      id: id,
      namespace_id: user.namespace.value.id
    });
  }
}