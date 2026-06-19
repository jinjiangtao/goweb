import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '表单模板列表' }
  },
  {
    path: '/designer/:id?',
    name: 'Designer',
    component: () => import('@/views/Designer.vue'),
    meta: { title: '表单设计器' }
  },
  {
    path: '/fill/:id',
    name: 'Fill',
    component: () => import('@/views/FillForm.vue'),
    meta: { title: '填写表单' }
  },
  {
    path: '/data/:id',
    name: 'Data',
    component: () => import('@/views/DataManage.vue'),
    meta: { title: '数据管理' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = (to.meta.title || '低代码表单') + ' - 低代码表单设计平台'
  next()
})

export default router
