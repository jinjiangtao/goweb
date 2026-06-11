import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'Odometer' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/UserManagement.vue'),
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: 'roles',
        name: 'Roles',
        component: () => import('@/views/RoleManagement.vue'),
        meta: { title: '角色管理', icon: 'UserFilled' }
      },
      {
        path: 'menus',
        name: 'Menus',
        component: () => import('@/views/MenuManagement.vue'),
        meta: { title: '菜单管理', icon: 'Menu' }
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/ProductManagement.vue'),
        meta: { title: '产品管理', icon: 'Box' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.path === '/login') {
    if (userStore.token) {
      next('/')
    } else {
      next()
    }
  } else {
    if (userStore.token) {
      next()
    } else {
      next('/login')
    }
  }
})

export default router
