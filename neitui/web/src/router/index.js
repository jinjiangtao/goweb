import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'JobList',
    component: () => import('@/views/JobList.vue')
  },
  {
    path: '/job/:id',
    name: 'JobDetail',
    component: () => import('@/views/JobDetail.vue')
  },
  {
    path: '/my-submissions',
    name: 'MySubmissions',
    component: () => import('@/views/MySubmissions.vue')
  },
  {
    path: '/favorites',
    name: 'Favorites',
    component: () => import('@/views/Favorites.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
