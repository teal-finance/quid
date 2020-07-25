import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    component: () => import('../views/ViewHome.vue'),
  },
  {
    path: '/users',
    component: () => import('../views/ViewUsers.vue'),
  },
  {
    path: '/groups',
    component: () => import('../views/ViewGroups.vue'),
  },
  {
    path: '/tokens',
    component: () => import('../views/ViewTokens.vue'),
  },
  {
    path: '/namespaces',
    component: () => import('../views/ViewNamespaces.vue'),
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.afterEach(() => {
  store.commit("endAction")
})

export default router
