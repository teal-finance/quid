import { requests } from "@/api";
import OrgContract from "./contract";
import { OrgTable } from "./interface";

export default class Org {
  id: number;
  name: string;

  constructor(data: OrgContract) {
    this.id = data.id;
    this.name = data.name;
  }

  toTableRow(): OrgTable {
    const row = Object.create(this);
    row.actions = [];
    return row as OrgTable;
  }

  static async fetchAll(): Promise<Array<OrgTable>> {
    const url = "/admin/orgs/all";
    const data = new Array<OrgTable>();
    try {
      const resp = await requests.get<Array<OrgContract>>(url);
      resp.forEach((row) => data.push(new Org(row).toTableRow()));
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return data;
  }

  static async delete(id: number) {
    await requests.post("/admin/orgs/delete", {
      id: id,
    });
  }
}