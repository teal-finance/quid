/** The composable parameters */
interface UseApiParams {
  serverUrl?: string;
  csrfCookieName?: string;
  csrfHeaderKey?: string;
  credentials?: RequestCredentials | null;
  mode?: RequestMode;
}

/** The response error contructor params */
interface ResponseErrorParams {
  status: number;
  statusText: string;
  content?: Record<string, any> | Array<any>;
  text?: string;
  errMsg?: string;
}

export { UseApiParams, ResponseErrorParams }