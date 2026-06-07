import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import MyBookings from '../views/MyBookings.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/my-bookings', name: 'MyBookings', component: MyBookings }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
