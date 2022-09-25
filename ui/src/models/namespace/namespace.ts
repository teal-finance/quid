import { api } from "@/api";
import NamespaceContract from "./contract";
import NamespaceTable from "./interface";
import Group from "@/models/group";
import { GroupContract } from "../group/contract";
import { AlgoType } from "@/interface";

export default class Namespace {
  id: number;
  name: string;
  algo: AlgoType = "HS256";
  maxTokenTtl: string;
  maxRefreshTokenTtl: string;
  publicEndpointEnabled: boolean;

  constructor({ id, name, alg, max_access_ttl, max_refresh_ttl, public_endpoint_enabled }: {
    id: number,
    name: string,
    alg: AlgoType,
    max_access_ttl: string,
    max_refresh_ttl: string,
    public_endpoint_enabled: boolean
  }) {
    this.id = id;
    this.name = name;
    this.algo = alg;
    this.maxTokenTtl = max_access_ttl;
    this.maxRefreshTokenTtl = max_refresh_ttl;
    this.publicEndpointEnabled = public_endpoint_enabled;
  }

  // *************************
  //   factory constructors
  // *************************

  static fromNamespaceTable(nst: NamespaceTable): Namespace {
    return new Namespace({
      id: nst.id,
      name: nst.name,
      alg: nst.algo,
      max_access_ttl: nst.maxTokenTtl,
      max_refresh_ttl: nst.maxRefreshTokenTtl,
      public_endpoint_enabled: nst.publicEndpointEnabled,
    })
  }

  static empty(): Namespace {
    return new Namespace({ id: 0, name: "default", alg: "HS256", max_access_ttl: "10", max_refresh_ttl: "10m", public_endpoint_enabled: false })
  }

  // *************************
  //         methods
  // *************************

  toTableRow(): NamespaceTable {
    const row: NamespaceTable = {
      id: this.id,
      name: this.name,
      algo: this.algo,
      maxTokenTtl: this.maxTokenTtl,
      maxRefreshTokenTtl: this.maxRefreshTokenTtl,
      publicEndpointEnabled: this.publicEndpointEnabled,
      actions: []
    };
    return row
  }

  // *************************
  //    static methods
  // *************************

  static async saveMaxAccessTokenTtl(id: number, ttl: string) {
    await api.post("/admin/namespaces/max-ttl", {
      id: id,
      max_ttl: ttl,
    });
  }

  static async saveMaxRefreshTokenTtl(id: number, ttl: string) {
    await api.post("/admin/namespaces/max-refresh-ttl", {
      id: id,
      refresh_max_ttl: ttl,
    });
  }

  static async delete(id: number) {
    await api.post("/admin/namespaces/delete", {
      id: id,
    });
  }

  static async getKey(id: number): Promise<string> {
    const data = await api.post<{ key: string }>("/admin/namespaces/key", {
      id: id,
    });
    return data.key
  }

  static async togglePublicEndpoint(id: number, enabled: boolean): Promise<void> {
    await api.post("/admin/namespaces/endpoint", {
      id: id,
      enable: enabled,
    });
  }

  static async fetchAll(): Promise<Array<NamespaceTable>> {
    const url = "/admin/namespaces/all";
    const ns = new Array<NamespaceTable>();
    try {
      const resp = await api.get<Array<NamespaceContract>>(url, true);
      resp.forEach((row) => {
        console.log(row)
        ns.push(new Namespace(row).toTableRow())
      });
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return ns;
  }

  static async fetchRowInfo(id: number): Promise<{ numUsers: number, groups: Array<Group> }> {
    const res: { numUsers: number, groups: Array<Group> } = { numUsers: 0, groups: [] };
    const data = await api.post<{ num_users: number, groups: Array<GroupContract> | null }>("/admin/namespaces/info", {
      id: id,
    });
    res.numUsers = data.num_users;
    //console.log("GG", data)
    if (!data.groups) {
      return res
    }
    for (const groupdata of data.groups) {
      res.groups.push(Group.fromContract(groupdata))
    }
    return res;
  }
}
