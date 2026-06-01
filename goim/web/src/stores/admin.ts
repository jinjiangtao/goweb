import { defineStore } from 'pinia'
import { ref } from 'vue'
import { adminLogin, getAdminUsers, getAdminMessages, type AdminUser, type AdminPagination, type AdminMessage } from '@/api'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const isLoggedIn = ref(!!token.value)
  const users = ref<AdminUser[]>([])
  const messages = ref<AdminMessage[]>([])
  const pagination = ref<AdminPagination>({ page: 1, page_size: 10, total: 0, total_page: 0 })
  const messagesPagination = ref<AdminPagination>({ page: 1, page_size: 20, total: 0, total_page: 0 })
  const loading = ref(false)
  const messagesLoading = ref(false)

  async function login(username: string, password: string) {
    const result = await adminLogin(username, password)
    token.value = result.token
    localStorage.setItem('admin_token', result.token)
    isLoggedIn.value = true
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('admin_token')
    isLoggedIn.value = false
    users.value = []
    messages.value = []
  }

  async function fetchUsers(page = 1, pageSize = 10) {
    loading.value = true
    try {
      const result = await getAdminUsers(page, pageSize, token.value)
      users.value = result.users
      pagination.value = result.pagination
    } finally {
      loading.value = false
    }
  }

  async function fetchMessages(page = 1, pageSize = 20) {
    messagesLoading.value = true
    try {
      const result = await getAdminMessages(page, pageSize, token.value)
      messages.value = result.messages
      messagesPagination.value = result.pagination
    } finally {
      messagesLoading.value = false
    }
  }

  return {
    token,
    isLoggedIn,
    users,
    messages,
    pagination,
    messagesPagination,
    loading,
    messagesLoading,
    login,
    logout,
    fetchUsers,
    fetchMessages,
  }
})
