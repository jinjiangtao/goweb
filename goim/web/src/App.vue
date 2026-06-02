<script setup lang="ts">
import { onMounted, watch, ref, computed, onUnmounted } from 'vue'
import { useChatStore } from './stores/chat'
import { useWebSocket } from './composables/useWebSocket'
import { useRoute } from 'vue-router'
import LoginPage from './components/LoginPage.vue'
import FriendList from './components/FriendList.vue'
import ChatWindow from './components/ChatWindow.vue'
import type { WSMessage } from './types'

const route = useRoute()
const chatStore = useChatStore()
const { connect, onMessage, disconnect, removeMessageHandler } = useWebSocket()

const wsInitialized = ref(false)
let refreshInterval: ReturnType<typeof setInterval> | null = null

const isAdminPage = computed(() => {
  return route.path.startsWith('/admin')
})

const shouldInitWs = computed(() => {
  return !isAdminPage.value && chatStore.isLoggedIn && chatStore.userId && !wsInitialized.value
})

const handleWsMessage = (msg: WSMessage) => {
  chatStore.addMessage(msg)
}

watch(shouldInitWs, (should) => {
  if (should && chatStore.userId) {
    wsInitialized.value = true
    console.log('App: initializing WebSocket for user', chatStore.userId)
    
    onMessage(handleWsMessage)
    
    setTimeout(() => {
      connect(chatStore.userId)
      
      refreshInterval = setInterval(() => {
        chatStore.refreshFriendsOnlineStatus()
        chatStore.loadOnlineUsers()
      }, 5000)
    }, 200)
  }
}, { immediate: true })

watch(() => chatStore.isLoggedIn, (loggedIn) => {
  if (!loggedIn) {
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
    disconnect()
    removeMessageHandler(handleWsMessage)
    wsInitialized.value = false
  }
})

onMounted(() => {
  if (!isAdminPage.value && chatStore.isLoggedIn && !chatStore.currentUser) {
    const payload = JSON.parse(atob(chatStore.token.split('.')[1]))
    chatStore.currentUser = {
      id: payload.user_id,
      username: payload.username,
      nickname: payload.username,
      avatar: ''
    }
    chatStore.loadFriends()
    chatStore.loadGroups()
    chatStore.loadOnlineUsers()
  }
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<template>
  <div id="app-root">
    <!-- 只在管理页面时显示 router-view -->
    <template v-if="isAdminPage">
      <router-view />
    </template>
    
    <!-- 聊天界面 -->
    <template v-else>
      <!-- 未登录时显示登录页面 -->
      <LoginPage v-if="!chatStore.isLoggedIn" />
      
      <!-- 已登录时显示聊天界面 -->
      <div v-else class="chat-container">
        <FriendList />
        <ChatWindow />
      </div>
    </template>
  </div>
</template>

<style scoped>
#app-root {
  width: 100%;
  height: 100vh;
}
</style>
