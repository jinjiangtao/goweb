import { defineStore } from 'pinia'
import { login as apiLogin, register as apiRegister, getUserInfo, getUserStats } from '@/utils/api'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    stats: {
      album_count: 0,
      photo_count: 0
    }
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.user?.nickname || state.user?.username || ''
  },

  actions: {
    async login(username, password) {
      const res = await apiLogin({ username, password })
      this.token = res.token
      this.user = res.user
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      return res
    },

    async register(username, password, nickname) {
      const res = await apiRegister({ username, password, nickname })
      this.token = res.token
      this.user = res.user
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      return res
    },

    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },

    async fetchUserInfo() {
      const res = await getUserInfo()
      this.user = res
      localStorage.setItem('user', JSON.stringify(res))
    },

    async fetchStats() {
      const res = await getUserStats()
      this.stats = res
      return res
    }
  }
})
