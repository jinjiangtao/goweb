import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('user_token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('user_info') || 'null'))

  const isLogin = computed(() => !!token.value)

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('user_token', newToken)
  }

  function setUserInfo(info) {
    userInfo.value = info
    localStorage.setItem('user_info', JSON.stringify(info))
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('user_token')
    localStorage.removeItem('user_info')
  }

  return {
    token,
    userInfo,
    isLogin,
    setToken,
    setUserInfo,
    logout,
  }
})
