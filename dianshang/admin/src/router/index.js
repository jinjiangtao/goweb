import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('../layout/index.vue'),
    redirect: '/admin/users',
    children: [
      {
        path: '/admin/users',
        name: 'UserManage',
        component: () => import('../views/UserManage.vue')
      },
      {
        path: '/admin/menus',
        name: 'MenuManage',
        component: () => import('../views/MenuManage.vue')
      },
      {
        path: '/admin/roles',
        name: 'RoleManage',
        component: () => import('../views/RoleManage.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path === '/login') {
    next()
  } else {
    if (token) {
      next()
    } else {
      next('/login')
    }
  }
})

export default router