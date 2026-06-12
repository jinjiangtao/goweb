import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi } from '../api/auth.js'

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
    console.log('设置菜单数据:', newMenus) // 调试日志
    menus.value = newMenus || []
    localStorage.setItem('menus', JSON.stringify(newMenus || []))
  }

  const login = async (loginData) => {
    const res = await loginApi(loginData)
    console.log('登录完整响应:', res) // 调试日志
    setToken(res.data.token)
    setUserInfo(res.data.user)
    setMenus(res.data.menus)
    return res
  }

  const fetchUserInfo = async () => {
    // 暂时不需要实现
    return
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
