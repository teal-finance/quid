import Vue from 'vue'
import router from './router'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from './App.vue'
import store from './store'
import Conf from './conf';
import api from "./api";
import notify from "./notify";

Vue.config.productionTip = false

Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)

const axiosConfig = {
  baseURL: Conf.quidUrl,
  timeout: 5000,
  withCredentials: false,
};

Vue.prototype.$axiosConfig = axiosConfig
Vue.prototype.$api = api

const vue = new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app');

Vue.prototype.$notify = notify(vue)

export default vue;
