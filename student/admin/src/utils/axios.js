import axios from 'axios'

const instance = axios.create({
  timeout: 10000
})

const token = localStorage.getItem('token')
if (token) {
  instance.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

export const setAuthToken = (token) => {
  if (token) {
    instance.defaults.headers.common['Authorization'] = `Bearer ${token}`
    localStorage.setItem('token', token)
  } else {
    delete instance.defaults.headers.common['Authorization']
    localStorage.removeItem('token')
  }
}

export default instance