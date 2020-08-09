import vue from './main'
import QuidRequests from '@/quidjs/requests'
import Conf from "@/conf";

var requests = new QuidRequests({
  namespace: "quid",
  timeouts: {
    accessToken: "5m",
    refreshToken: "24h"
  },
  accessTokenUri: "/admin_token/access/",
  axiosConfig: {
    baseURL: Conf.quidUrl,
    timeout: 5000,
    withCredentials: process.env.NODE === 'production',
    headers: { 'Content-Type': 'application/json' }
  },
  verbose: Conf.isProduction
})

function apiError(e) {
  console.log("API ERROR:", e);
  if (e.response === undefined || e.response === null) {
    if (typeof e === 'object') {
      if ("error" in e) {
        vue.$notify.error(`${e.error}`)
        return
      }
    }
    vue.$notify.error(`${e}`)
  }
  else {
    if (e.response.status !== 200) {
      if (e.response.status === 404) {
        vue.$notify.warning({
          title: "Not found",
          content: e.response.uri
        })
        return
      } else {
        vue.$notify.error(`${e.response.status} ${e.response.data.error}`)
      }
    }
  }
}

const api = {
  requests: requests,
  adminLogin: async function (username, password) {
    try {
      await requests.adminLogin(username, password)
    } catch (e) {
      if (e.unauthorized) {
        vue.$notify.warning({
          title: "Login failed",
          content: "Authentication refused"
        });
        return "unauthorized"
      } else {
        return e
      }
    }
    return null
  },
  get: async function (uri) {
    try {
      let response = await this.requests.get(uri);
      return { response: response, error: null }
    } catch (e) {
      if (e.hasToLogin) {
        console.log("has to login")
        vue.$store.commit("unauthenticate")
        return
      }
      apiError(e)
      if (e.response !== undefined) {
        if (e.response.status !== 404) {
          return { response: null, error: e }
        }
      }
      return e;
    }
  },
  post:
    async function (uri, payload) {
      try {
        let response = await this.requests.post(uri, payload)
        return { response: response, error: null }
      } catch (e) {
        if (e.hasToLogin) {
          console.log("has to login")
          vue.$store.commit("unauthenticate")
          return
        }
        apiError(e)
        if (e.response !== undefined) {
          if (e.response.status !== 404) {
            return { response: null, error: e }
          }
        }
        return e
      }
    }
}

export default api;