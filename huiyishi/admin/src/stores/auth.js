
import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/admin',
  timeout: 10000
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  console.log('请求:', config)
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    console.error('响应错误:', error)
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const admin = ref(JSON.parse(localStorage.getItem('admin') || 'null'))

  const setAuth = (newToken, newAdmin) => {
    token.value = newToken
    admin.value = newAdmin
    localStorage.setItem('token', newToken)
    localStorage.setItem('admin', JSON.stringify(newAdmin))
  }

  const logout = () => {
    token.value = ''
    admin.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('admin')
  }

  return { token, admin, setAuth, logout, api }
})
