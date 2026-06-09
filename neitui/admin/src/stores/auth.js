import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthenticated = () => !!token.value

  const login = async (username, password) => {
    const res = await api.post('/user/login', { username, password })
    token.value = res.token
    user.value = res.user
    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(res.user))
  }

  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  const isEmployee = () => user.value?.role === 'employee'
  const isHR = () => ['hr', 'admin'].includes(user.value?.role)
  const isAdmin = () => user.value?.role === 'admin'

  return { token, user, isAuthenticated, login, logout, isEmployee, isHR, isAdmin }
})
