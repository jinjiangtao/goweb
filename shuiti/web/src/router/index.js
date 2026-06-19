import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/practice',
    name: 'Practice',
    component: () => import('@/views/Practice.vue')
  },
  {
    path: '/exam/:mode',
    name: 'Exam',
    component: () => import('@/views/Exam.vue'),
    props: true
  },
  {
    path: '/result/:recordId',
    name: 'Result',
    component: () => import('@/views/Result.vue'),
    props: true
  },
  {
    path: '/wrong',
    name: 'Wrong',
    component: () => import('@/views/Wrong.vue')
  },
  {
    path: '/stats',
    name: 'Stats',
    component: () => import('@/views/Stats.vue')
  },
  {
    path: '/records',
    name: 'Records',
    component: () => import('@/views/Records.vue')
  },
  {
    path: '/record/:id',
    name: 'RecordDetail',
    component: () => import('@/views/RecordDetail.vue'),
    props: true
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/Admin.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
