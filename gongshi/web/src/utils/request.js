import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000
})

request.interceptors.request.use(config => {
  const user = JSON.parse(localStorage.getItem('gongshi_user') || '{}')
  if (user && user.id) {
    config.headers['X-User-ID'] = user.id
    config.headers['X-User-Role'] = user.role
  }
  return config
})

request.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response) {
      const msg = error.response.data?.error || '请求失败'
      ElMessage.error(msg)
    } else {
      ElMessage.error('网络连接失败')
    }
    return Promise.reject(error)
  }
)

export default request
