import { createRouter, createWebHistory } from 'vue-router'
import Login from './views/Login.vue'
import Layout from './views/Layout.vue'
import Prizes from './views/Prizes.vue'
import Records from './views/Records.vue'
import Stats from './views/Stats.vue'
import AddressManage from './views/AddressManage.vue'

const routes = [
  { path: '/login', component: Login },
  {
    path: '/',
    component: Layout,
    redirect: '/stats',
    children: [
      { path: 'prizes', component: Prizes },
      { path: 'records', component: Records },
      { path: 'stats', component: Stats },
      { path: 'address', component: AddressManage }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
