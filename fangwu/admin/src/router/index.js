import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    redirect: '/houses',
    children: [
      {
        path: 'houses',
        name: 'Houses',
        component: () => import('../views/House.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.name !== 'Login' && !userStore.token) {
    next({ name: 'Login' })
  } else if (to.name === 'Login' && userStore.token) {
    next({ name: 'Houses' })
  } else {
    next()
  }
})

export default router
