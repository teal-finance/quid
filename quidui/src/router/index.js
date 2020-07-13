import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    component: () => import('../views/Home.vue'),
  },
  {
    path: '/users',
    component: () => import('../views/Users.vue'),
  },
  {
    path: '/groups',
    component: () => import('../views/Groups.vue'),
  },
  {
    path: '/tokens',
    component: () => import('../views/Tokens.vue'),
  },
  {
    path: '/namespaces',
    component: () => import('../views/Namespaces.vue'),
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
