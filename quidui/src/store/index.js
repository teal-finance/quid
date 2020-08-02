import Vue from 'vue'
import Vuex from 'vuex'
import conf from '@/conf'

Vue.use(Vuex)

const store = new Vuex.Store({
    state: {
        isAuthenticated: false,
        username: null,
        action: null,
        showActionBar: false,
        refreshAction: null,
    },
    mutations: {
        action(state, name) {
            state.action = name;
            state.showActionBar = true;
        },
        endAction(state) {
            state.action = null;
            state.showActionBar = false;
        },
        refreshAction(state, name) {
            state.action = name;
            state.showActionBar = true;
        },
        endRefreshAction(state) {
            state.action = null;
            state.showActionBar = false;
        },
        authenticate(state, { username, token = null }) {
            state.isAuthenticated = true;
            state.username = username
            if (token && !conf.isProduction) {
                console.log("Storing refresh token", token)
                localStorage.setItem("refreshToken", token)
                localStorage.setItem("username", username)
            }
        },
        unauthenticate(state) {
            state.isAuthenticated = false;
            state.username = null;
            if (!conf.isProduction) {
                localStorage.removeItem("refreshToken")
                localStorage.removeItem("username")
            }
        },
    },
    getters: {
        showActionBar: (state) => {
            return state.showActionBar;
        },
    }
});

export default store;
