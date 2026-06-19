import { defineStore } from 'pinia'
import { ref } from 'vue'

const STORAGE_KEY = 'shuiti_exam_progress'

export const useExamStore = defineStore('exam', () => {
  const progress = ref(JSON.parse(localStorage.getItem(STORAGE_KEY) || 'null'))

  const saveProgress = (data) => {
    progress.value = data
    localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
  }

  const clearProgress = () => {
    progress.value = null
    localStorage.removeItem(STORAGE_KEY)
  }

  return {
    progress,
    saveProgress,
    clearProgress
  }
})
