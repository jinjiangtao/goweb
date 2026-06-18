import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layout/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'template',
        name: 'Template',
        component: () => import('@/views/Template.vue'),
        meta: { title: '流程模板' }
      },
      {
        path: 'template/designer',
        name: 'TemplateDesigner',
        component: () => import('@/views/TemplateDesigner.vue'),
        meta: { title: '流程设计器' }
      },
      {
        path: 'template/designer/:id',
        name: 'TemplateDesignerEdit',
        component: () => import('@/views/TemplateDesigner.vue'),
        meta: { title: '编辑流程' }
      },
      {
        path: 'approval',
        name: 'Approval',
        component: () => import('@/views/Approval.vue'),
        meta: { title: '我的申请' }
      },
      {
        path: 'approval/todo',
        name: 'ApprovalTodo',
        component: () => import('@/views/ApprovalTodo.vue'),
        meta: { title: '待我审批' }
      },
      {
        path: 'approval/submit/:templateId',
        name: 'ApprovalSubmit',
        component: () => import('@/views/ApprovalSubmit.vue'),
        meta: { title: '提交审批' }
      },
      {
        path: 'approval/detail/:id',
        name: 'ApprovalDetail',
        component: () => import('@/views/ApprovalDetail.vue'),
        meta: { title: '审批详情' }
      },
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/User.vue'),
        meta: { title: '用户管理' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
