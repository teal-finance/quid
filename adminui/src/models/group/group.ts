import { requests } from "@/api";
import { GroupContract } from "./contract";
import { GroupTable } from "./interface";

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
    const row = Object.create(this);
    row.actions = [];
    return row as GroupTable;
  }

  // *************************
  //    static methods
  // *************************

  static async fetchAll(nsid: number): Promise<Array<GroupTable>> {
    const url = "/admin/groups/nsall";
    const ns = new Array<GroupTable>();
    try {
      const payload = { namespace_id: nsid }
      const resp = await requests.post<Array<GroupContract>>(url, payload);
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

  static async delete(id: number) {
    await requests.post("/admin/groups/delete", {
      id: id,
    });
  }
}