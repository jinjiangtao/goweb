import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi } from '@/api/user'
import request from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)

  function setUserInfo(info) {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function clearUser() {
    userInfo.value = {}
    token.value = ''
    localStorage.removeItem('userInfo')
    localStorage.removeItem('token')
  }

  async function login(loginForm) {
    const res = await loginApi(loginForm)
    setToken(res.data.token)
    setUserInfo(res.data.user)
    return res
  }

  function logout() {
    clearUser()
  }

  async function fetchUserInfo() {
    const res = await request({
      url: '/user/profile',
      method: 'get'
    })
    setUserInfo(res.data)
    return res
  }

  return {
    userInfo,
    token,
    isLoggedIn,
    setUserInfo,
    setToken,
    clearUser,
    login,
    logout,
    fetchUserInfo
  }
})
