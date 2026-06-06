
import { defineStore } from 'pinia'
import { ref } from 'vue'

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

  return { token, admin, setAuth, logout }
})
