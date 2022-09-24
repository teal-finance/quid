/** Type of a response error */
type ResponseErrorType = "request" | "json";

/** The composable parameters */
interface UseApiParams {
  serverUrl?: string;
  csrfCookieName?: string;
  csrfHeaderKey?: string;
  credentials?: RequestCredentials | null;
  mode?: RequestMode;
}

export { ResponseErrorType, UseApiParams }