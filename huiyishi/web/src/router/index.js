
import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import MyBookings from '../views/MyBookings.vue'
import Login from '../views/Login.vue'

const routes = [
  { path: '/login', name: 'Login', component: Login },
  { path: '/', name: 'Home', component: Home, meta: { requiresAuth: true } },
  { path: '/my-bookings', name: 'MyBookings', component: MyBookings, meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to, from, next) =&gt; {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth &amp;&amp; !token) {
    next('/login')
  } else if (to.path === '/login' &amp;&amp; token) {
    next('/')
  } else {
    next()
  }
})

export default router

