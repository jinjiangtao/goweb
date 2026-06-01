<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useWebSocket } from '@/composables/useWebSocket'
import LoginPage from '@/components/LoginPage.vue'
import FriendList from '@/components/FriendList.vue'
import ChatWindow from '@/components/ChatWindow.vue'

const chatStore = useChatStore()
const { messages: wsMessages, connect: connectWs, isConnected: wsConnected } = useWebSocket(chatStore.userId)
const hasInitializedWs = ref(false)

watch(() => wsMessages.value, (newMessages) => {
  if (newMessages && newMessages.length > 0) {
    const latestMsg = newMessages[newMessages.length - 1]
    chatStore.addMessage(latestMsg)
    wsMessages.value = []
  }
}, { deep: true })

watch(() => chatStore.isLoggedIn, (isLoggedIn) => {
  if (isLoggedIn && !hasInitializedWs.value) {
    hasInitializedWs.value = true
    setTimeout(() => {
      connectWs()
    }, 200)
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
</script>

<template>
  <LoginPage v-if="!chatStore.isLoggedIn" />
  
  <div v-else class="chat-container">
    <FriendList />
    <ChatWindow />
  </div>
</template>