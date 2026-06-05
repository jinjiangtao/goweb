import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const menus = ref(JSON.parse(localStorage.getItem('menus') || '[]'))

  const login = async (username, password) => {
    const response = await axios.post('/api/admin/login', { username, password })
    token.value = response.data.token
    user.value = response.data.user
    menus.value = response.data.menus
    localStorage.setItem('token', token.value)
    localStorage.setItem('user', JSON.stringify(user.value))
    localStorage.setItem('menus', JSON.stringify(menus.value))
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
    return response.data
  }

  const logout = () => {
    token.value = ''
    user.value = null
    menus.value = []
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('menus')
    delete axios.defaults.headers.common['Authorization']
  }

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
    axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
  }

  const setMenus = (newMenus) => {
    menus.value = newMenus
    localStorage.setItem('menus', JSON.stringify(newMenus))
  }

  return {
    token,
    user,
    menus,
    login,
    logout,
    setToken,
    setMenus
  }
})
