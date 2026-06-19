import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userId = ref(localStorage.getItem('shuiti_user_id') || '')

  const setUserId = (id) => {
    userId.value = id
    localStorage.setItem('shuiti_user_id', id)
  }

  return {
    userId,
    setUserId
  }
})
