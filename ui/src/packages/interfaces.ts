type ResponseErrorType = "request" | "json";

interface UseApiParams {
  serverUrl?: string;
  csrfCookieName?: string;
  csrfHeaderKey?: string;
  credentials?: RequestCredentials | null;
  mode?: RequestMode;
}

export { ResponseErrorType, UseApiParams }