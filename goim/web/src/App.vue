<script setup lang="ts">
import { onMounted } from 'vue'
import { useChatStore } from '@/stores/chat'
import LoginPage from '@/components/LoginPage.vue'
import FriendList from '@/components/FriendList.vue'
import ChatWindow from '@/components/ChatWindow.vue'

const chatStore = useChatStore()

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