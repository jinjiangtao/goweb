import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

export const boardApi = {
  list: () => request.get('/boards'),
  get: (id) => request.get(`/boards/${id}`),
  create: (data) => request.post('/boards', data),
  update: (id, data) => request.put(`/boards/${id}`, data),
  delete: (id) => request.delete(`/boards/${id}`),
  saveElements: (id, elements) => request.post(`/boards/${id}/elements`, { elements }),
  getElements: (id) => request.get(`/boards/${id}/elements`),
  clearElements: (id) => request.post(`/boards/${id}/clear`),
  getHistory: (id) => request.get(`/boards/${id}/history`),
  restoreHistory: (id, historyId) => request.post(`/boards/${id}/history/${historyId}/restore`)
}

export default request
