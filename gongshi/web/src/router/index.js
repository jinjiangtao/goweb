import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    redirect: '/timesheet',
    children: [
      {
        path: 'timesheet',
        name: 'Timesheet',
        component: () => import('../views/Timesheet.vue'),
        meta: { title: '工时填报' }
      },
      {
        path: 'history',
        name: 'History',
        component: () => import('../views/History.vue'),
        meta: { title: '历史记录' }
      },
      {
        path: 'review',
        name: 'Review',
        component: () => import('../views/Review.vue'),
        meta: { title: '工时审核', admin: true }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('../views/Statistics.vue'),
        meta: { title: '数据统计' }
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('../views/Projects.vue'),
        meta: { title: '项目管理', admin: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const user = JSON.parse(localStorage.getItem('gongshi_user') || 'null')
  if (to.path === '/login') {
    next()
    return
  }
  if (!user) {
    next('/login')
    return
  }
  if (to.meta.admin && user.role !== 'admin') {
    next('/timesheet')
    return
  }
  next()
})

export default router
