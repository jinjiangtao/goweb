
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
  console.log('API 请求:', config)
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    console.error('API 响应错误:', error)
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
