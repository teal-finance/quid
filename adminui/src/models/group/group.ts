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

  static async fetchAll(): Promise<Array<GroupTable>> {
    const url = "/admin/groups/all";
    const ns = new Array<GroupTable>();
    try {
      const resp = await requests.get<Array<GroupContract>>(url);
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

}