import { requests } from "@/api";
import NamespaceContract from "./contract";

export default class Namespace {
  id: number;
  name: string;
  maxTokenTtl: string;
  maxRefreshTokenTtl: string;
  publicEndpointEnabled: boolean;

  constructor({ id, name, max_token_ttl, max_refresh_token_ttl, public_endpoint_enabled }: {
    id: number,
    name: string,
    max_token_ttl: string,
    max_refresh_token_ttl: string,
    public_endpoint_enabled: boolean
  }) {
    this.id = id;
    this.name = name;
    this.maxTokenTtl = max_token_ttl;
    this.maxRefreshTokenTtl = max_refresh_token_ttl;
    this.publicEndpointEnabled = public_endpoint_enabled;
  }

  static async fetchAll(): Promise<Set<Namespace>> {
    const url = "/admin/namespaces/all";
    const ns = new Set<Namespace>([]);
    try {
      const resp = await requests.get<Array<NamespaceContract>>(url);
      resp.forEach((row) => {
        console.log(row)
        ns.add(new Namespace(row))
      });
    } catch (e) {
      console.log("Err", e);
      throw e;
    }
    return ns;
  }
}