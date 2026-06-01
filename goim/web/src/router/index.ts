import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Chat',
    component: () => import('../App.vue'),
  },
  {
    path: '/admin',
    name: 'AdminLogin',
    component: () => import('../views/AdminLogin.vue'),
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: () => import('../views/AdminUsers.vue'),
  },
  {
    path: '/admin/messages',
    name: 'AdminMessages',
    component: () => import('../views/AdminMessages.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
