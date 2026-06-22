import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@/pages/DashboardPage.vue'),
    meta: { title: '概览' },
  },
  {
    path: '/generate',
    name: 'generate',
    component: () => import('@/pages/GeneratePage.vue'),
    meta: { title: '生成中心' },
  },
  {
    path: '/parse',
    name: 'parse',
    component: () => import('@/pages/ParsePage.vue'),
    meta: { title: '解析中心' },
  },
  {
    path: '/templates',
    name: 'templates',
    component: () => import('@/pages/TemplatesPage.vue'),
    meta: { title: '模板库' },
  },
  {
    path: '/archive',
    name: 'archive',
    component: () => import('@/pages/ArchivePage.vue'),
    meta: { title: '归档管理' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  document.title = `${to.meta.title || 'QRForge'} · QRForge`
})

export default router
