import axios from 'axios'
import store from './store'
import vue from './main'

function apiError(e) {
  console.log("API ERR", e);
  if (e.response !== undefined) {
    if (e.response.status !== 200) {
      if (e.response.status === 401) {
        store.commit("unauthenticate");
        return
      }
      vue.$bvToast.toast(
        `${e.response.status} ${e.response.data.error}`,
        {
          title: "Error",
          variant: "danger"
        }
      );
    }
  } else {
    console.log("API ERROR:", e);
  }
}

const Api = {
  get: async function (uri) {
    try {
      let response = await axios.get(uri, vue.$axiosConfig);
      return { response: response, error: null };
    } catch (e) {
      if (e.response !== undefined) {
        return { response: null, error: e }
      }
      apiError(e)
    }
  },
  post:
    async function (uri, payload) {
      try {
        let response = await axios.post(uri, payload, vue.$axiosConfig);
        return { response: response, error: null };
      } catch (e) {
        if (e.response !== undefined) {
          return { response: null, error: e }
        }
        apiError(e)
      }
    }
}

export default Api;