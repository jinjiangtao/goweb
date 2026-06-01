<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useWebSocket } from '@/composables/useWebSocket'

const chatStore = useChatStore()
const messageInput = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

const { sendMessage: sendWsMessage } = useWebSocket(chatStore.userId)

watch(() => chatStore.messages?.length || 0, async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
})

function getInitials(name: string): string {
  return name.charAt(0).toUpperCase()
}

function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

function isSelf(msg: { sender_id: string }): boolean {
  return msg.sender_id === chatStore.currentUser?.id
}

async function handleSend() {
  if (!messageInput.value.trim() || !chatStore.currentUser) return

  const wsMsg = await chatStore.sendMessage(messageInput.value.trim())
  if (wsMsg) {
    console.log('Sending message via WebSocket:', wsMsg)
    const localMessage = {
      id: `temp-${Date.now()}`,
      sender_id: chatStore.currentUser.id,
      receiver_id: wsMsg.to,
      receiver_type: wsMsg.to_type,
      content: wsMsg.content,
      type: wsMsg.msg_type,
      status: 0,
      created_at: new Date().toISOString()
    }
    chatStore.messages.push(localMessage)
    await sendWsMessage(wsMsg)
    messageInput.value = ''
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}
</script>

<template>
  <div class="chat-window">
    <template v-if="chatStore.currentFriend || chatStore.currentGroup">
      <div class="chat-header">
        <div class="avatar">
          {{ getInitials(chatStore.currentFriend?.nickname || chatStore.currentGroup?.name || '') }}
        </div>
        <div>
          <h4>{{ chatStore.currentFriend?.nickname || chatStore.currentGroup?.name }}</h4>
          <p v-if="chatStore.currentFriend">
            {{ chatStore.currentFriend.online ? '在线' : '离线' }}
          </p>
        </div>
      </div>
      
      <div ref="messagesContainer" class="chat-messages">
        <div 
          v-for="msg in chatStore.messages" 
          :key="msg.id"
          :class="['message', { self: isSelf(msg), other: !isSelf(msg) }]"
        >
          <div class="message-bubble">
            <template v-if="msg.type === 1">
              <img :src="msg.content" class="message-image" />
            </template>
            <template v-else>
              {{ msg.content }}
            </template>
          </div>
          <div class="message-time">{{ formatTime(msg.created_at) }}</div>
        </div>
      </div>
      
      <div class="chat-input-area">
        <button class="upload-btn">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
            <polyline points="17 8 12 3 7 8" />
            <line x1="12" y1="3" x2="12" y2="15" />
          </svg>
        </button>
        <input 
          v-model="messageInput"
          type="text" 
          placeholder="输入消息..."
          @keydown="handleKeydown"
        />
        <button class="send-btn" @click="handleSend">发送</button>
      </div>
    </template>
    
    <div v-else class="empty-chat">
      选择一个好友或群组开始聊天
    </div>
  </div>
</template>