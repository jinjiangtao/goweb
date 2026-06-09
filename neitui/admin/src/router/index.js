import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/jobs',
    children: [
      {
        path: 'jobs',
        name: 'Jobs',
        component: () => import('@/views/employee/Jobs.vue'),
        meta: { role: 'employee' }
      },
      {
        path: 'referrals',
        name: 'Referrals',
        component: () => import('@/views/employee/Referrals.vue'),
        meta: { role: 'employee' }
      },
      {
        path: 'admin/referrals',
        name: 'AdminReferrals',
        component: () => import('@/views/hr/Referrals.vue'),
        meta: { role: 'hr' }
      },
      {
        path: 'admin/jobs',
        name: 'AdminJobs',
        component: () => import('@/views/hr/Jobs.vue'),
        meta: { role: 'hr' }
      },
      {
        path: 'admin/stats',
        name: 'Stats',
        component: () => import('@/views/hr/Stats.vue'),
        meta: { role: 'hr' }
      },
      {
        path: 'admin/users',
        name: 'Users',
        component: () => import('@/views/admin/Users.vue'),
        meta: { role: 'admin' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.path === '/login') {
    if (authStore.isAuthenticated()) {
      next('/')
    } else {
      next()
    }
    return
  }

  if (!authStore.isAuthenticated()) {
    next('/login')
    return
  }

  if (to.meta.role === 'hr' && !authStore.isHR()) {
    next('/')
    return
  }

  if (to.meta.role === 'admin' && !authStore.isAdmin()) {
    next('/')
    return
  }

  next()
})

export default router
