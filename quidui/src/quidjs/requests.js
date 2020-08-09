import axios from 'axios'
import quidException from "./exceptions"

export default class QuidRequests {
  refreshToken = null;
  #accessToken = null;

  constructor({ namespace, axiosConfig, timeouts = {
    accessToken: "20m",
    refreshToken: "24h"
  },
    accessTokenUri = null,
    verbose = false }) {
    if (typeof namespace !== 'string') {
      throw quidException({ error: 'Parameter namespace has to be set' });
    }
    if (typeof axiosConfig !== 'object') {
      throw quidException({ error: 'Parameter axiosConfig has to be set' });
    }
    this.namespace = namespace
    this.axiosConfig = axiosConfig;
    this.axios = axios.create(this.axiosConfig);
    this.timeouts = timeouts
    this.verbose = verbose
    this.accessTokenUri = accessTokenUri
  }

  async get(uri) {
    return await this._requestWithRetry(uri, "get")
  }

  async post(uri, payload) {
    return await this._requestWithRetry(uri, "post", payload)
  }

  async adminLogin(username, password) {
    let uri = "/admin_login";
    let payload = {
      namespace: "quid",
      username: username,
      password: password,
    }
    try {
      let response = await axios.post(uri, payload, this.axiosConfig);
      this.refreshToken = response.data.token;
    } catch (e) {
      if (e.response) {
        if (e.response.status === 401) {
          throw quidException({ error: null, unauthorized: true });
        }
      }
      throw quidException({ error: e });
    }
  }

  async _requestWithRetry(uri, method, payload, retry = 0) {
    if (this.verbose) {
      console.log(method + " request to " + uri)
    }
    await this.checkTokens();
    try {
      if (method === "get") {
        return await this.axios.get(uri, this.axiosConfig);
      } else {
        return await axios.post(uri, payload, this.axiosConfig);
      }
    } catch (e) {
      if (e.response) {
        if (e.response.status === 401) {
          if (this.verbose) {
            console.log("Access token has expired")
          }
          this.#accessToken = null;
          await this.checkTokens();
          if (retry > 2) {
            throw quidException({ error: "too many retries" });
          }
          retry++
          if (this.verbose) {
            console.log("Retrying", method, "request to", uri, ", retry", retry)
          }
          return await this._requestWithRetry(uri, method, payload, retry)
        } else {
          throw quidException({ error: e });
        }
      } else {
        throw quidException({ error: e });
      }
    }
  }

  async checkTokens() {
    if (this.refreshToken === null) {
      if (this.verbose) {
        console.log("Tokens check: no refresh token")
      }
      throw quidException({ error: 'No refresh token found', hasToLogin: true });
    }
    if (this.#accessToken === null) {
      if (this.verbose) {
        console.log("Tokens check: no access token")
      }
      let { token, error, statusCode } = await this._getAccessToken();
      if (error !== null) {
        if (statusCode === 401) {
          if (this.verbose) {
            console.log("The refresh token has expired")
          }
          throw quidException({ error: 'The refresh token has expired', hasToLogin: true });
        } else {
          throw quidException({ error: error });
        }
      }
      this.#accessToken = token;
      this.axiosConfig.headers.Authorization = "Bearer " + this.#accessToken
      this.axios = axios.create(this.axiosConfig);
    }
  }

  async _getAccessToken() {
    try {
      let payload = {
        namespace: this.namespace,
        refresh_token: this.refreshToken,
      }
      let url = "/token/access/" + this.timeouts.accessToken
      if (this.accessTokenUri !== null) {
        url = this.accessTokenUri
      }
      if (this.verbose) {
        console.log("Getting an access token from", url, payload)
      }
      let response = await axios.post(url, payload, this.axiosConfig);
      return { token: response.data.token, error: null, statusCode: response.status };
    } catch (e) {
      if (e.response !== undefined) {
        return { token: null, error: e.response.data.error, statusCode: e.response.status };
      }
      return { token: null, error: e, statusCode: null }
    }
  }
}