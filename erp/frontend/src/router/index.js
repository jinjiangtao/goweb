import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () =&gt; import('@/views/Login.vue')
  },
  {
    path: '/',
    component: () =&gt; import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () =&gt; import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'House' }
      },
      {
        path: 'system/user',
        name: 'Users',
        component: () =&gt; import('@/views/UserManagement.vue'),
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: 'system/role',
        name: 'Roles',
        component: () =&gt; import('@/views/RoleManagement.vue'),
        meta: { title: '角色管理', icon: 'UserFilled' }
      },
      {
        path: 'system/menu',
        name: 'Menus',
        component: () =&gt; import('@/views/MenuManagement.vue'),
        meta: { title: '菜单管理', icon: 'Menu' }
      },
      {
        path: 'product',
        name: 'Products',
        component: () =&gt; import('@/views/ProductManagement.vue'),
        meta: { title: '产品管理', icon: 'Goods' }
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
