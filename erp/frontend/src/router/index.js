import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Layout from '../layout/index.vue'
import Dashboard from '../views/Dashboard.vue'
import UserManagement from '../views/UserManagement.vue'
import RoleManagement from '../views/RoleManagement.vue'
import MenuManagement from '../views/MenuManagement.vue'
import ProductManagement from '../views/ProductManagement.vue'
import CustomerManagement from '../views/CustomerManagement.vue'
import SupplierManagement from '../views/SupplierManagement.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: Dashboard,
        meta: { title: '仪表盘', icon: 'House' }
      },
      {
        path: 'system/user',
        name: 'Users',
        component: UserManagement,
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: 'system/role',
        name: 'Roles',
        component: RoleManagement,
        meta: { title: '角色管理', icon: 'UserFilled' }
      },
      {
        path: 'system/menu',
        name: 'Menus',
        component: MenuManagement,
        meta: { title: '菜单管理', icon: 'Menu' }
      },
      {
        path: 'product',
        name: 'Products',
        component: ProductManagement,
        meta: { title: '产品管理', icon: 'Goods' }
      },
      {
        path: 'customer',
        name: 'Customers',
        component: CustomerManagement,
        meta: { title: '客户管理', icon: 'User' }
      },
      {
        path: 'supplier',
        name: 'Suppliers',
        component: SupplierManagement,
        meta: { title: '供应商管理', icon: 'OfficeBuilding' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 简单的路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path === '/login') {
    if (token) {
      next('/')
    } else {
      next()
    }
  } else {
    if (token) {
      next()
    } else {
      next('/login')
    }
  }
})

export default router
