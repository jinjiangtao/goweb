<script setup lang="ts">
import { onMounted, watch, ref, computed, onUnmounted } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useWebSocket } from '@/composables/useWebSocket'
import LoginPage from '@/components/LoginPage.vue'
import FriendList from '@/components/FriendList.vue'
import ChatWindow from '@/components/ChatWindow.vue'
import type { WSMessage } from '@/types'

const chatStore = useChatStore()
const { sendMessage, connect, onMessage, disconnect, connected } = useWebSocket()

const wsInitialized = ref(false)
let refreshInterval: ReturnType<typeof setInterval> | null = null

const shouldInitWs = computed(() => {
  return chatStore.isLoggedIn && chatStore.userId && !wsInitialized.value
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
    wsInitialized.value = false
  }
})

onMounted(() => {
  if (chatStore.isLoggedIn && !chatStore.currentUser) {
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
  <LoginPage v-if="!chatStore.isLoggedIn" />
  
  <div v-else class="chat-container">
    <FriendList />
    <ChatWindow />
  </div>
</template>
