<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'

onMounted(() => {
  const authStore = useAuthStore()
  if (authStore.token) {
    import('axios').then(({ default: axios }) => {
      axios.defaults.headers.common['Authorization'] = `Bearer ${authStore.token}`
    })
  }
})
</script>
