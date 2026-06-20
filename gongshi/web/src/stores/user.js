import { defineStore } from 'pinia'
import { login } from '../api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('gongshi_user') || 'null')
  }),
  getters: {
    isAdmin: (state) => state.user?.role === 'admin',
    isLoggedIn: (state) => !!state.user
  },
  actions: {
    async doLogin(username, password) {
      const res = await login({ username, password })
      this.user = res.user
      localStorage.setItem('gongshi_user', JSON.stringify(res.user))
      localStorage.setItem('gongshi_token', res.token)
      return res.user
    },
    logout() {
      this.user = null
      localStorage.removeItem('gongshi_user')
      localStorage.removeItem('gongshi_token')
    }
  }
})
