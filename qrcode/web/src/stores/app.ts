import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Stats } from '@/types'
import { getStats } from '@/api'

export const useAppStore = defineStore('app', () => {
  const stats = ref<Stats | null>(null)
  const loading = ref(false)

  async function refreshStats() {
    loading.value = true
    try {
      stats.value = await getStats()
    } finally {
      loading.value = false
    }
  }

  return { stats, loading, refreshStats }
})
