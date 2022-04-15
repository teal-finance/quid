import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: () => import('./views/HomeView.vue')
  },
  {
    path: '/namespaces',
    name: 'Namespaces',
    component: () => import('./views/NamespaceView.vue')
  },
  {
    path: '/admins',
    name: 'Admins',
    component: () => import('./views/AdminsView.vue')
  },
  {
    path: '/group',
    name: 'Groups',
    component: () => import('./views/GroupView.vue')
  },
  {
    path: '/user',
    name: 'Users',
    component: () => import('./views/UserView.vue')
  },
  {
    path: '/org',
    name: 'Orgs',
    component: () => import('./views/OrgView.vue')
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
