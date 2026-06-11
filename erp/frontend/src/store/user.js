import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, getCurrentUser } from '@/api/auth'
import { getMenusByRole } from '@/api/menu'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  const menus = ref(JSON.parse(localStorage.getItem('menus') || '[]'))

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  const setMenus = (newMenus) => {
    menus.value = newMenus
    localStorage.setItem('menus', JSON.stringify(newMenus))
  }

  const login = async (loginData) => {
    const res = await loginApi(loginData)
    setToken(res.token)
    setUserInfo(res.user)
    setMenus(res.menus || [])
    return res
  }

  const fetchUserInfo = async () => {
    const res = await getCurrentUser()
    setUserInfo(res)
    return res
  }

  const logout = () => {
    token.value = ''
    userInfo.value = {}
    menus.value = []
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('menus')
  }

  return {
    token,
    userInfo,
    menus,
    setToken,
    setUserInfo,
    setMenus,
    login,
    fetchUserInfo,
    logout
  }
})
