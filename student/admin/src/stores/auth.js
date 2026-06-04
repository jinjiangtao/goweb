import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')

  const login = async (username, password) => {
    const response = await axios.post('/api/admin/login', { username, password })
    token.value = response.data.token
    localStorage.setItem('token', token.value)
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
    return response.data
  }

  const logout = () => {
    token.value = ''
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
  }

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
    axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
  }

  return {
    token,
    login,
    logout,
    setToken
  }
})
