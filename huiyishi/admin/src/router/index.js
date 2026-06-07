
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Login from '../views/Login.vue'
import Layout from '../views/Layout.vue'
import Rooms from '../views/Rooms.vue'
import Bookings from '../views/Bookings.vue'
import Stats from '../views/Stats.vue'
import Users from '../views/Users.vue'

const routes = [
  { path: '/login', component: Login },
  {
    path: '/',
    component: Layout,
    redirect: '/stats',
    children: [
      { path: 'stats', component: Stats },
      { path: 'rooms', component: Rooms },
      { path: 'bookings', component: Bookings },
      { path: 'users', component: Users }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) =&gt; {
  const authStore = useAuthStore()
  if (to.path !== '/login' &amp;&amp; !authStore.token) {
    next('/login')
  } else if (to.path === '/login' &amp;&amp; authStore.token) {
    next('/')
  } else {
    next()
  }
})

export default router
