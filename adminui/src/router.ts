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
    path: '/org',
    name: 'Orgs',
    component: () => import('./views/OrgView.vue')
  },
  {
    path: "/settings",
    component: () => import("./views/SettingsView.vue"),
    meta: {
      title: "Settings"
    }
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
