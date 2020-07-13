import Vue from 'vue'
import router from './router'
import axios from "axios";
import VueCookies from 'vue-cookies'

import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from './App.vue'
import store from './store'
import Conf from './conf';
import Api from "./api";

Vue.config.productionTip = false

Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)
Vue.use(VueCookies)
Vue.$cookies.config('1d', '', '', false, 'Strict')

const axiosConfig = {
  baseURL: Conf.quidUrl,
  timeout: 5000
};
const ax = axios.create(axiosConfig);

Vue.prototype.$axiosConfig = axiosConfig
Vue.prototype.$axios = ax
Vue.prototype.$api = Api

const vue = new Vue({
  router,
  store: store,
  render: h => h(App)
}).$mount('#app');

export default vue;
