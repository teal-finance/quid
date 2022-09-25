import { AlgoType } from "@/interface";

interface NamespaceContract {
  id: number;
  name: string;
  alg: AlgoType;
  max_access_ttl: string;
  max_refresh_ttl: string;
  public_endpoint_enabled: boolean;
}

export default NamespaceContract
