interface NamespaceContract {
  id: number;
  name: string;
  max_token_ttl: string;
  max_refresh_token_ttl: string;
  public_endpoint_enabled: boolean;
}

export default NamespaceContract