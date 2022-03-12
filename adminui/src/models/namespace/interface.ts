interface NamespaceTable {
  id: number;
  name: string;
  maxTokenTtl: string;
  maxRefreshTokenTtl: string;
  publicEndpointEnabled: boolean;
  actions: Array<string>;
}

export default NamespaceTable