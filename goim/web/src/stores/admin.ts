import { defineStore } from 'pinia'
import { ref } from 'vue'
import { adminLogin, getAdminUsers, type AdminUser, type AdminPagination } from '@/api'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const isLoggedIn = ref(!!token.value)
  const users = ref<AdminUser[]>([])
  const pagination = ref<AdminPagination>({ page: 1, page_size: 10, total: 0, total_page: 0 })
  const loading = ref(false)

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

  return {
    token,
    isLoggedIn,
    users,
    pagination,
    loading,
    login,
    logout,
    fetchUsers,
  }
})
