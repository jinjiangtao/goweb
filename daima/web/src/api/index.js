import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000
})

request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      if (res.code === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
      }
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  error => {
    if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error('网络错误')
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  register: data => request.post('/register', data),
  login: data => request.post('/login', data),
  getMe: () => request.get('/user/me')
}

export const snippetApi = {
  list: params => request.get('/snippets', { params }),
  search: params => request.get('/snippets/search', { params }),
  get: id => request.get(`/snippets/${id}`),
  create: data => request.post('/snippets', data),
  update: (id, data) => request.put(`/snippets/${id}`, data),
  delete: id => request.delete(`/snippets/${id}`),
  my: params => request.get('/my/snippets', { params }),
  myFavorites: params => request.get('/my/favorites', { params }),
  like: id => request.post(`/snippets/${id}/like`),
  favorite: id => request.post(`/snippets/${id}/favorite`),
  fork: id => request.post(`/snippets/${id}/fork`),
  diff: data => request.post('/diff', data),
  getLanguages: () => request.get('/languages')
}

export const commentApi = {
  list: snippetId => request.get(`/snippets/${snippetId}/comments`),
  create: (snippetId, data) => request.post(`/snippets/${snippetId}/comments`, data),
  delete: commentId => request.delete(`/comments/${commentId}`)
}

export default request
