import Vue from 'vue'
import Vuex from 'vuex'
import axios from "axios";
import Conf from "@/conf";

Vue.use(Vuex)

const store = new Vuex.Store({
    state: {
        isAuthenticated: false,
        username: null,
        key: null,
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
        authenticate(state, payload) {
            state.isAuthenticated = true;
            state.username = payload.username;
            state.key = payload.key;
            const axiosConf = {
                baseURL: Conf.quidUrl,
                timeout: 5000,
                headers: { Authorization: "Bearer " + payload.key }
            }
            const ax = axios.create(axiosConf);
            /*ax.interceptors.request.use(request => {
                console.log("Starting Request", request);
                return request;
            });*/
            Vue.prototype.$axiosConfig = axiosConf
            Vue.prototype.$axios = ax;
        },
        unauthenticate(state) {
            state.isAuthenticated = false;
            state.key = null;
            state.username = null;
            const axiosConf = {
                baseURL: Conf.quidUrl,
                timeout: 5000,
            }
            const ax = axios.create(axiosConf);
            Vue.prototype.$axiosConfig = axiosConf
            Vue.prototype.$axios = ax;
        },
    },
    getters: {
        showActionBar: (state) => {
            return state.showActionBar;
        }
    }
});

export default store;
