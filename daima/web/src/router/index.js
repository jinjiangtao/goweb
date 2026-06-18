import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/snippet/:id',
    name: 'snippet-detail',
    component: () => import('../views/SnippetDetail.vue')
  },
  {
    path: '/create',
    name: 'snippet-create',
    component: () => import('../views/SnippetCreate.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/edit/:id',
    name: 'snippet-edit',
    component: () => import('../views/SnippetCreate.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/my',
    name: 'my-snippets',
    component: () => import('../views/MySnippets.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/favorites',
    name: 'my-favorites',
    component: () => import('../views/MyFavorites.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/search',
    name: 'search',
    component: () => import('../views/Search.vue')
  },
  {
    path: '/diff',
    name: 'diff',
    component: () => import('../views/Diff.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    ElMessage.warning('请先登录')
    next('/')
  } else {
    next()
  }
})

export default router
