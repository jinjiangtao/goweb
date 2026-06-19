import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/notes',
    children: [
      {
        path: 'notes',
        name: 'Notes',
        component: () => import('@/views/Notes.vue')
      },
      {
        path: 'notes/:id',
        name: 'NoteDetail',
        component: () => import('@/views/NoteDetail.vue')
      },
      {
        path: 'editor',
        name: 'Editor',
        component: () => import('@/views/Editor.vue')
      },
      {
        path: 'editor/:id',
        name: 'EditorEdit',
        component: () => import('@/views/Editor.vue')
      },
      {
        path: 'tags',
        name: 'Tags',
        component: () => import('@/views/Tags.vue')
      },
      {
        path: 'archived',
        name: 'Archived',
        component: () => import('@/views/Notes.vue'),
        props: { status: 'archived' }
      },
      {
        path: 'trash',
        name: 'Trash',
        component: () => import('@/views/Notes.vue'),
        props: { status: 'deleted' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
