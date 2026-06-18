import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)

  const setAuth = (userData, userToken) => {
    user.value = userData
    token.value = userToken
    localStorage.setItem('token', userToken)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  const logout = () => {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  const checkAuth = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser && token.value) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (e) {
        console.error(e)
      }
    }
  }

  const refreshUser = async () => {
    try {
      const res = await authApi.getMe()
      if (res.data) {
        user.value = res.data
        localStorage.setItem('user', JSON.stringify(res.data))
      }
    } catch (e) {
      console.error(e)
    }
  }

  return {
    user,
    token,
    isLoggedIn,
    setAuth,
    logout,
    checkAuth,
    refreshUser
  }
})
