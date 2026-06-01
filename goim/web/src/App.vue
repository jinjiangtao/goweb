<script setup lang="ts">
import { onMounted, watch, ref, computed, onUnmounted } from 'vue'
import { useChatStore } from '@/stores/chat'
import LoginPage from '@/components/LoginPage.vue'
import FriendList from '@/components/FriendList.vue'
import ChatWindow from '@/components/ChatWindow.vue'

const chatStore = useChatStore()
const wsInitialized = ref(false)
let refreshInterval: ReturnType<typeof setInterval> | null = null

const shouldInitWs = computed(() => {
  return chatStore.isLoggedIn && chatStore.userId && !wsInitialized.value
})

watch(shouldInitWs, async (should) => {
  if (should && chatStore.userId) {
    wsInitialized.value = true
    console.log('App: initializing WebSocket for user', chatStore.userId)
    
    const { useWebSocket } = await import('@/composables/useWebSocket')
    const { messages: wsMessages, connect: connectWs, isConnected: wsConnected, sendMessage: wsSend } = useWebSocket(chatStore.userId)
    
    chatStore.setWsSendMessage(wsSend)
    
    watch(() => wsMessages.value, (newMessages) => {
      if (newMessages && newMessages.length > 0) {
        const latestMsg = newMessages[newMessages.length - 1]
        chatStore.addMessage(latestMsg)
        wsMessages.value = []
      }
    }, { deep: true })
    
    watch(wsConnected, (connected) => {
      chatStore.setWsConnected(connected)
      if (connected) {
        chatStore.refreshFriendsOnlineStatus()
        chatStore.loadOnlineUsers()
      }
    })
    
    setTimeout(() => {
      if (!wsConnected.value) {
        console.log('App: triggering WebSocket connection')
        connectWs()
      }
    }, 300)
    
    refreshInterval = setInterval(() => {
      chatStore.refreshFriendsOnlineStatus()
      chatStore.loadOnlineUsers()
    }, 5000)
  }
}, { immediate: true })

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