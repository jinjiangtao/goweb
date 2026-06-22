import axios from 'axios'

const client = axios.create({
  baseURL: '/api',
  timeout: 30000,
  responseType: 'json',
})

client.interceptors.response.use(
  (res) => res,
  (err) => {
    console.error('API error:', err?.response?.data?.message || err.message)
    return Promise.reject(err)
  },
)

export default client
