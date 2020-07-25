import axios from 'axios'
import store from './store'
import vue from './main'

function apiError(e) {
  console.log("API ERROR:", e);
  if (e.response === undefined || e.response === null) {
    vue.$bvToast.toast(
      `${e}`,
      {
        title: "Error",
        variant: "danger"
      }
    );
  }
  else {
    if (e.response.status !== 200) {
      if (e.response.status === 401) {
        store.commit("unauthenticate");
        return
      }
      if (e.response.status === 404) {
        vue.$bvToast.toast(
          `Not found`,
          {
            title: "Error",
            variant: "warning"
          }
        );
        return
      } else {
        vue.$bvToast.toast(
          `${e.response.status} ${e.response.data.error}`,
          {
            title: "Error",
            variant: "danger",
            noAutoHide: true
          }
        );
      }
    }
  }
}

const api = {
  get: async function (uri) {
    try {
      let response = await axios.get(uri, vue.$axiosConfig);
      return { response: response, error: null };
    } catch (e) {
      apiError(e)
      if (e.response !== undefined) {
        if (e.response.status !== 404) {
          return { response: null, error: e }
        }
      }
    }
  },
  post:
    async function (uri, payload) {
      try {
        let response = await axios.post(uri, payload, vue.$axiosConfig);
        return { response: response, error: null };
      } catch (e) {
        apiError(e)
        if (e.response !== undefined) {
          if (e.response.status !== 404) {
            return { response: null, error: e }
          }
        }
      }
    }
}

export default api;