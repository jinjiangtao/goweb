import { defineStore } from 'pinia'
import { ref } from 'vue'

const safeParse = (str, defaultValue) => {
  if (!str || str === 'undefined') {
    return defaultValue
  }
  try {
    return JSON.parse(str)
  } catch {
    return defaultValue
  }
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(safeParse(localStorage.getItem('user'), null))
  const menus = ref(safeParse(localStorage.getItem('menus'), []))

  function setToken(val) {
    token.value = val
    localStorage.setItem('token', val)
  }

  function setUser(val) {
    user.value = val
    localStorage.setItem('user', JSON.stringify(val))
  }

  function setMenus(val) {
    menus.value = val
    localStorage.setItem('menus', JSON.stringify(val))
  }

  function logout() {
    token.value = ''
    user.value = null
    menus.value = []
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('menus')
  }

  return { token, user, menus, setToken, setUser, setMenus, logout }
})