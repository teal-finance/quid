import { QuidParams, QuidLoginParams } from "./types";
import QuidError from "./errors";

export default class QuidRequests {
  public refreshToken: string | null = null;
  private accessToken: string | null = null;
  private quidUri: string;
  private serverUri: string;
  private namespace: string;
  private timeouts: Record<string, string>;
  private verbose: boolean;
  private headers: HeadersInit;
  private credentials: string | null;
  private accessTokenUri: string | null;

  public constructor({ quidUri, serverUri, namespace, timeouts = {
    accessToken: "20m",
    refreshToken: "24h"
  },
    credentials = "include",
    verbose = false,
    accessTokenUri = null }: QuidParams) {
    this.quidUri = quidUri;
    this.serverUri = serverUri;
    this.namespace = namespace;
    this.timeouts = timeouts;
    this.credentials = credentials;
    this.verbose = verbose;
    this.headers = {
      'Content-Type': 'application/json',
    } as HeadersInit;
    this.accessTokenUri = accessTokenUri;
    if (verbose) {
      console.log("Initializing QuidRequests", this.quidUri);
    }
  }

  async get<T>(url: string): Promise<T> {
    return this._request<T>(url, "get");
  }

  async post<T>(url: string): Promise<T> {
    return this._request<T>(url, "post");
  }

  async login(username: string, password: string) {
    await this.getRefreshToken({ username: username, password: password } as QuidLoginParams);
    await this.checkTokens();
  }

  async getRefreshToken({ username, password, refreshTokenTtl = "24h" }: QuidLoginParams) {
    const uri = this.quidUri + "/token/refresh/" + refreshTokenTtl;
    const payload = {
      namespace: this.namespace,
      username: username,
      password: password,
    }
    try {
      const opts: RequestInit = {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(payload),
      };
      const response = await fetch(uri, opts);
      if (!response.ok) {
        console.log("RESP NOT OK", response);
        throw new Error(response.statusText)
      }
      const t = await response.json();
      if (this.verbose) {
        console.log("Setting refresh token")
      }
      this.refreshToken = t.token;
    } catch (e) {
      throw new Error(`${e}`);
    }
  }

  async checkTokens(): Promise<void> {
    if (this.refreshToken === null) {
      if (this.verbose) {
        console.log("Tokens check: no refresh token")
      }
      throw new QuidError('No refresh token found', true);
    }
    if (this.accessToken === null) {
      if (this.verbose) {
        console.log("Tokens check: no access token")
      }
      const status = await this._getAccessToken();
      if (status === 401) {
        if (this.verbose) {
          console.log("The refresh token has expired")
        }
        throw new QuidError('The refresh token has expired');
      }
      this.headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': "Bearer " + this.accessToken
      } as HeadersInit;
    }
  }

  private async _request<T>(url: string, method: string, payload: Record<string, any> | Array<any> | null = null, retry = 0): Promise<T> { // eslint-disable-line
    if (this.verbose) {
      console.log(method + " request to " + url)
    }
    await this.checkTokens();
    let opts: RequestInit;
    if (method === "post") {
      //console.log("GET", this.#accessToken, uri);
      opts = {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(payload),
      };
      if (this.credentials !== null) {
        opts.credentials = this.credentials as RequestCredentials;
      }
    } else {
      opts = {
        method: 'GET',
        headers: this.headers,
      };
      if (this.credentials !== null) {
        opts.credentials = this.credentials as RequestCredentials;
      }
    }
    //console.log("FETCH", this.serverUri + url);
    //console.log(JSON.stringify(opts, null, "  "))
    const response = await fetch(this.serverUri + url, opts);
    if (!response.ok) {
      if (response.status === 401) {
        this.accessToken = null;
        this.checkTokens();
        retry++;
        if (retry > 2) {
          throw new Error("Too many retries")
        }
        if (this.verbose) {
          console.log("Request retry", retry)
        }
        return this._request<T>(url, method, payload, retry);
      }
      console.log("RESP NOT OK", response);
      throw new Error(response.statusText)
    }
    return await response.json() as T;
  }

  private async _getAccessToken(): Promise<number> {
    const payload = {
      namespace: this.namespace,
      refresh_token: this.refreshToken, // eslint-disable-line
    }
    let url = this.quidUri + "/token/access/" + this.timeouts.accessToken;
    if (this.accessTokenUri !== null) {
      url = this.accessTokenUri
    }
    if (this.verbose) {
      console.log("Getting an access token from", url, payload)
    }
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    }
    const opts: RequestInit = {
      method: 'POST',
      headers,
      body: JSON.stringify(payload),
    };
    const response = await fetch(url, opts);
    if (!response.ok) {
      return response.status;
    }
    const data = await response.json();
    this.accessToken = data.token;
    return response.status;
  }
}