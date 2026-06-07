import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

api.interceptors.request.use(config => {
  console.log('API 请求:', config)
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    console.error('API 响应错误:', error)
    return Promise.reject(error)
  }
)

export default api
