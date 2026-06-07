
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

// 初始化时设置token
function initAuth() {
  const token = localStorage.getItem('token')
  if (token) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`
  }
}

initAuth()

// 更新token的方法
export function setAuthToken(token) {
  if (token) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`
  } else {
    delete api.defaults.headers.common['Authorization']
  }
}

api.interceptors.request.use(config =&gt; {
  console.log('API 请求:', config)
  return config
})

api.interceptors.response.use(
  response =&gt; response,
  error =&gt; {
    console.error('API 响应错误:', error)
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      setAuthToken(null)
      window.location.href = '#/login'
    }
    return Promise.reject(error)
  }
)

export default api

