import { AlgoType } from "@/interface";

interface NamespaceTable {
  id: number;
  name: string;
  algo: AlgoType;
  maxTokenTtl: string;
  maxRefreshTokenTtl: string;
  publicEndpointEnabled: boolean;
  actions: Array<string>;
}

export default NamespaceTable