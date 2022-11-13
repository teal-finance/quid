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
      const payload = { ns_id: nsid }
      const resp = await api.post<Array<GroupContract>>(url, payload, false, true);
      resp.data.forEach((row) => {
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
      const payload = { id: uid, ns_id: user.namespace.value.id }
      const resp = await api.post<{ groups: Array<GroupContract> | null }>(url, payload);
      //console.log("RESP", JSON.stringify(resp.groups, null, "  "))
      if (resp.data.groups) {
        if (resp.data.groups.length > 0) {
          resp.data.groups.forEach((row) => {
            data.push(new Group(row).toTableRow())
          });
        }
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
        usr_id: uid,
        grp_id: gid,
        ns_id: user.namespace.value.id
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
        usr_id: uid,
        grp_id: gid,
        ns_id: user.namespace.value.id
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
      ns_id: user.namespace.value.id
    });
  }
}
