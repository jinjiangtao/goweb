import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000
})

request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    ElMessage.error(error.response?.data?.error || '请求失败')
    return Promise.reject(error)
  }
)

export const noteApi = {
  list: (params = {}) => request.get('/notes', { params }),
  get: (id) => request.get(`/notes/${id}`),
  create: (data) => request.post('/notes', data),
  update: (id, data) => request.put(`/notes/${id}`, data),
  remove: (id, permanent = false) => request.delete(`/notes/${id}`, { params: { permanent } }),
  verify: (id, password) => request.post(`/notes/${id}/verify`, { password }),
  restore: (id) => request.post(`/notes/${id}/restore`),
  togglePin: (id) => request.post(`/notes/${id}/pin`),
  toggleArchive: (id) => request.post(`/notes/${id}/archive`),
  emptyTrash: () => request.delete('/trash/empty'),
  stats: () => request.get('/stats')
}

export const tagApi = {
  list: (params = {}) => request.get('/tags', { params }),
  get: (id) => request.get(`/tags/${id}`),
  create: (data) => request.post('/tags', data),
  update: (id, data) => request.put(`/tags/${id}`, data),
  remove: (id) => request.delete(`/tags/${id}`),
  getCount: (id) => request.get(`/tags/${id}/count`),
  reorder: (tagIds) => request.post('/tags/reorder', { tagIds })
}

export const tagGroupApi = {
  list: () => request.get('/tag-groups'),
  create: (data) => request.post('/tag-groups', data),
  update: (id, data) => request.put(`/tag-groups/${id}`, data),
  remove: (id) => request.delete(`/tag-groups/${id}`)
}

export default request
