import { createRouter, createWebHashHistory } from 'vue-router'
import Home from './views/Home.vue'
import MyPrizes from './views/MyPrizes.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/my-prizes', component: MyPrizes }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
