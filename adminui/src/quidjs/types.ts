interface QuidParams {
  quidUri: string;
  serverUri: string;
  namespace: string;
  timeouts: Record<string, string>;
  credentials?: string | null;
  verbose: boolean;
  accessTokenUri?: string | null;
}

interface QuidLoginParams {
  username: string;
  password: string;
  refreshTokenTtl: string;
}

export { QuidParams, QuidLoginParams };