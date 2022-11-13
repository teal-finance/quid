/** The composable parameters */
interface UseApiParams {
  serverUrl?: string;
  csrfCookieName?: string;
  csrfHeaderKey?: string;
  credentials?: RequestCredentials | null;
  mode?: RequestMode;
}

/** The standard api response with typed data */
interface ApiResponse<T = Record<string, any> | Array<any>> {
  ok: boolean;
  url: string;
  headers: Record<string, string>;
  status: number;
  statusText: string;
  data: T;
  text: string;
}

type OnResponseHook = <T>(res: ApiResponse<T>) => Promise<ApiResponse<T>>;

export { UseApiParams, ApiResponse, OnResponseHook }