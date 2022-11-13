import Cookies from 'js-cookie';
import { ApiResponse, OnResponseHook, UseApiParams } from './interfaces';

/** The main api composable */
const useApi = (params: UseApiParams = {
  serverUrl: "",
  csrfCookieName: "csrftoken",
  csrfHeaderKey: "X-CSRFToken",
  credentials: "include",
  mode: "cors",
}) => {
  // options
  let _serverUrl = params.serverUrl ?? "";
  let _csrfCookieName = params?.csrfCookieName ?? "csrftoken";
  let _csrfHeaderKey = params?.csrfHeaderKey ?? "X-CSRFToken";
  let _mode = params?.mode ?? "cors";
  let _credentials: RequestCredentials | null = params.credentials ?? "include";
  // state
  let csrfToken: string | null = null;
  // hooks
  let _onResponse: OnResponseHook;

  const hasCsrfCookie = (): boolean => {
    const cookie = Cookies.get(_csrfCookieName);
    if (cookie) {
      return true
    }
    return false
  }

  const _csrfFromCookie = (): string => {
    const c = Cookies.get(_csrfCookieName);
    if (!c) {
      throw ("Csrf cookie not found")
    }
    return c
  }

  /** Set the on response hook */
  const onResponse = (hook: OnResponseHook) => {
    _onResponse = hook;
  }

  /** Set a csrf token to use with request headers */
  const setCsrfToken = (token: string) => {
    csrfToken = token;
  }

  /** Get the csrf token from a cookie and set it to use with request headers */
  const setCsrfTokenFromCookie = (verbose = false): boolean => {
    if (hasCsrfCookie()) {
      csrfToken = _csrfFromCookie();
      if (verbose) {
        console.log("User logged in with csrf cookie, setting api token", csrfToken);
      }
      return true
    } else {
      if (verbose) {
        console.log("User does not have csrf cookie")
      }
    }
    return false
  }

  const _processResponse = async <T>(response: Response): Promise<ApiResponse<T>> => {
    const head: Record<string, string> = {};
    response.headers.forEach((v, k) => head[k] = v);
    let apiResp: ApiResponse<T> = {
      ok: response.ok,
      url: response.url,
      headers: head,
      status: response.status,
      statusText: response.statusText,
      data: {} as unknown as T,
      text: "",
    }
    if (!response.headers.get("Content-Type")?.startsWith("application/json")) {
      const txt = await response.text();
      apiResp.text = txt;
    } else {
      try {
        apiResp.data = (await response.json()) as T
      } catch (e) {
        throw new Error(`Json parsing error: ${e}`);
      }
    }
    if (_onResponse) {
      apiResp = await _onResponse(apiResp)
    }
    return apiResp
  }

  /** Post request */
  const post = async <T>(
    uri: string,
    payload: Array<any> | Record<string, any> | FormData,
    multipart: boolean = false,
    verbose = false
  ): Promise<ApiResponse<T>> => {
    const opts = _postHeader(payload, "post", multipart);
    let url = _serverUrl + uri;
    if (verbose) {
      console.log("POST", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    return await _processResponse<T>(response)
  }

  /** Patch request */
  const patch = async <T>(
    uri: string, payload: Array<any> | Record<string, any>, verbose = false
  ): Promise<ApiResponse<T>> => {
    const opts = _postHeader(payload, "patch");
    let url = _serverUrl + uri;
    if (verbose) {
      console.log("PATCH", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    return await _processResponse<T>(response)
  }

  /** Put request */
  const put = async <T>(
    uri: string, payload: Array<any> | Record<string, any>, verbose = false
  ): Promise<ApiResponse<T>> => {
    let url = _serverUrl + uri;
    const opts = _postHeader(payload, "put");
    if (verbose) {
      console.log("PUT", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    return await _processResponse<T>(response)
  }

  /** Get request */
  const get = async <T>(uri: string, verbose = false): Promise<ApiResponse<T>> => {
    let url = _serverUrl + uri;
    const opts = _getHeader("get");
    if (verbose) {
      console.log("GET", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    return await _processResponse<T>(response)
  }

  /** Delete request */
  const del = async <T>(uri: string, verbose = false): Promise<ApiResponse<T>> => {
    const url = _serverUrl + uri;
    const opts = _getHeader("delete");
    if (verbose) {
      console.log("DELETE", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    return await _processResponse<T>(response)
  }

  const _getHeader = (method: string = "get"): RequestInit => {
    const h = {
      method: method,
      headers: { "Content-Type": "application/json" },
      mode: _mode,
    } as RequestInit;
    if (_credentials !== null) {
      h.credentials = _credentials
    }
    if (csrfToken !== null) {
      h.headers = { "Content-Type": "application/json" }
      h.headers[_csrfHeaderKey] = csrfToken;
    }
    return h;
  }

  const _postHeader = (payload: Array<any> | Record<string, any> | FormData, method = "post", multipart: boolean = false): RequestInit => {
    const pl = multipart ? payload as FormData : JSON.stringify(payload);
    const r: RequestInit = {
      method: method,
      mode: _mode,
      body: pl
    };
    if (!multipart) {
      r.headers = { "Content-Type": "application/json" }
    } else {
      r.headers = { "Content-Type": "multipart/form-data" }
    }
    if (_credentials !== null) {
      r.credentials = _credentials
    }
    if (csrfToken !== null) {
      r.headers[_csrfHeaderKey] = csrfToken;
    }
    return r;
  }

  return {
    csrfToken,
    hasCsrfCookie,
    setCsrfToken,
    setCsrfTokenFromCookie,
    onResponse,
    get,
    post,
    put,
    patch,
    del,
  }
}

export { useApi }
