import Cookies from 'js-cookie';
import { ResponseError } from './errors';
import { UseApiParams } from './interfaces';

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
  let _credentials: RequestCredentials | null = params?.credentials ?? "include";
  // state
  let _csrfToken: string | null = null;

  const _hasCsrfCookie = (): boolean => {
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

  /** Set a csrf token to use with request headers */
  const setCsrfToken = (token: string) => {
    _csrfToken = token;
  }

  /** Get the csrf token from a cookie and set it to use with request headers */
  const setCsrfTokenFromCookie = (verbose = false): boolean => {
    if (_hasCsrfCookie()) {
      if (verbose) {
        console.log("User logged in with csrf cookie, setting api token", _csrfFromCookie);
      }
      _csrfToken = _csrfFromCookie();
      return true
    } else {
      if (verbose) {
        console.log("User does not have csrf cookie")
      }
    }
    return false
  }

  /** Post request */
  const post = async <T>(
    uri: string,
    payload: Array<any> | Record<string, any> | FormData,
    multipart: boolean = false,
    verbose = false
  ): Promise<T> => {
    const opts = _postHeader(payload, "post", multipart);
    let url = _serverUrl + uri;
    if (verbose) {
      console.log("POST", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    if (!response.ok) {
      throw new ResponseError(response, "request");
    }
    if (response.headers.get("content-type") != "application/json") {
      // no data
      return await response.text() as unknown as T;
    }
    let _data: T
    try {
      _data = (await response.json()) as unknown as T
    } catch (e) {
      throw new ResponseError(response, "json", `${e}`);
    }
    return _data;
  }

  /** Patch request */
  const patch = async <T>(uri: string, payload: Array<any> | Record<string, any>, verbose = false) => {
    const opts = _postHeader(payload, "patch");
    let url = _serverUrl + uri;
    if (verbose) {
      console.log("PATCH", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    if (!response.ok) {
      throw new ResponseError(response, "request");
    }
    if (response.headers.get("content-type") != "application/json") {
      // no data
      return await response.text() as unknown as T;
    }
    let _data: T
    try {
      _data = (await response.json()) as T
    } catch (e) {
      throw new ResponseError(response, "json", `${e}`);
    }
    return _data;
  }

  /** Put request */
  const put = async <T>(uri: string, payload: Array<any> | Record<string, any>, verbose = false) => {
    let url = _serverUrl + uri;
    const opts = _postHeader(payload, "put");
    if (verbose) {
      console.log("PUT", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    if (!response.ok) {
      throw new ResponseError(response, "request");
    }
    if (response.headers.get("content-type") != "application/json") {
      // no data
      return await response.text() as unknown as T;
    }
    let _data: T
    try {
      _data = (await response.json()) as T
    } catch (e) {
      throw new ResponseError(response, "json", `${e}`);
    }
    return _data;
  }

  /** Get request */
  const get = async <T>(uri: string, verbose = false): Promise<T> => {
    let url = _serverUrl + uri;
    const opts = _getHeader("get");
    if (verbose) {
      console.log("GET", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    if (!response.ok) {
      throw new ResponseError(response, "request");
    }
    if (response.headers.get("content-type") != "application/json") {
      // no data
      return await response.text() as unknown as T;
    }
    let _data: T
    try {
      _data = (await response.json()) as T
    } catch (e) {
      throw new ResponseError(response, "json", `${e}`);
    }
    return _data;
  }

  /** Delete request */
  const del = async (uri: string, verbose = false): Promise<void> => {
    const url = _serverUrl + uri;
    const opts = _getHeader("delete");
    if (verbose) {
      console.log("DELETE", url);
      console.log(JSON.stringify(opts, null, "  "));
    }
    const response = await fetch(url, opts);
    if (!response.ok) {
      throw new ResponseError(response, "request");
    }
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
    if (_csrfToken !== null) {
      h.headers = { "Content-Type": "application/json" }
      h.headers[_csrfHeaderKey] = _csrfToken;
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
    }
    if (_credentials !== null) {
      r.credentials = _credentials
    }
    if (_csrfToken !== null) {
      if (multipart) {
        r.headers = {}
        r.headers[_csrfHeaderKey] = _csrfToken;
      } else {
        r.headers = { "Content-Type": "application/json" }
        r.headers[_csrfHeaderKey] = _csrfToken;
      }
    }
    return r;
  }

  return {
    setCsrfToken,
    setCsrfTokenFromCookie,
    get,
    post,
    put,
    patch,
    del,
  }
}

export { useApi }
