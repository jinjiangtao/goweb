<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useWebSocket } from '@/composables/useWebSocket'

const chatStore = useChatStore()
const messageInput = ref('')
const messagesContainer = ref<HTMLElement | null>(null)
const showEmojiPicker = ref(false)
const { sendMessage } = useWebSocket()

// 常用表情
const emojis = [
  '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂',
  '🙂', '🙃', '😉', '😊', '😇', '🥰', '😍', '🤩',
  '😘', '😗', '😚', '😋', '😛', '😜', '🤪', '😝',
  '🤑', '🤗', '🤭', '🤫', '🤔', '🤐', '🤨', '😐',
  '😑', '😶', '😏', '😒', '🙄', '😬', '🤥', '😌',
  '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢',
  '🤮', '🤧', '🥵', '🥶', '🥴', '😵', '🤯', '🤠',
  '🥳', '😎', '🤓', '🧐', '👍', '👎', '👏', '🙌'
]

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

function insertEmoji(emoji: string) {
  messageInput.value += emoji
  showEmojiPicker.value = false
}

async function handleImageUpload(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0] && chatStore.userId && chatStore.currentUser) {
    const currentUser = chatStore.currentUser
    try {
      // 上传图片（复用头像上传接口，或者需要创建专门的消息图片上传接口）
      // 这里暂时直接发送图片URL，实际项目中应该有专门的消息图片上传接口
      const reader = new FileReader()
      reader.onload = async (e) => {
        const imageUrl = e.target?.result as string
        const wsMsg = await chatStore.sendMessage(imageUrl, 1)
        if (wsMsg) {
          const localMessage = {
            id: `temp-${Date.now()}`,
            sender_id: currentUser.id,
            receiver_id: wsMsg.to,
            receiver_type: wsMsg.to_type,
            content: imageUrl,
            type: 1,
            status: 0,
            created_at: new Date().toISOString()
          }
          chatStore.messages.push(localMessage)
          await sendMessage(wsMsg)
        }
      }
      reader.readAsDataURL(target.files[0])
    } catch (error) {
      console.error('上传图片失败:', error)
    }
  }
}

async function handleSend() {
  if (!messageInput.value.trim() || !chatStore.currentUser) return
  const currentUser = chatStore.currentUser

  const wsMsg = await chatStore.sendMessage(messageInput.value.trim())
  if (wsMsg) {
    console.log('Sending message via WebSocket:', wsMsg)
    const localMessage = {
      id: `temp-${Date.now()}`,
      sender_id: currentUser.id,
      receiver_id: wsMsg.to,
      receiver_type: wsMsg.to_type,
      content: wsMsg.content,
      type: wsMsg.msg_type,
      status: 0,
      created_at: new Date().toISOString()
    }
    chatStore.messages.push(localMessage)
    
    await sendMessage(wsMsg)
    messageInput.value = ''
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

function openImage(url: string) {
  window.open(url, '_blank')
}
</script>

<template>
  <div class="chat-window">
    <template v-if="chatStore.currentFriend || chatStore.currentGroup || chatStore.currentOwner">
      <div class="chat-header">
        <div class="avatar">
          <img 
            v-if="chatStore.currentFriend?.avatar || chatStore.currentGroup?.avatar || chatStore.currentOwner?.avatar" 
            :src="chatStore.currentFriend?.avatar || chatStore.currentGroup?.avatar || chatStore.currentOwner?.avatar" 
            alt="头像"
          />
          <span v-else>
            {{ getInitials(chatStore.currentFriend?.nickname || chatStore.currentGroup?.name || chatStore.currentOwner?.name || '') }}
          </span>
        </div>
        <div>
          <div class="chat-title">
            <h4>{{ chatStore.currentFriend?.nickname || chatStore.currentGroup?.name || chatStore.currentOwner?.name }}</h4>
            <span v-if="chatStore.currentOwner" class="owner-badge">官方群主</span>
          </div>
          <p v-if="chatStore.currentFriend">
            {{ chatStore.currentFriend.online ? '在线' : '离线' }}
          </p>
          <p v-if="chatStore.currentOwner">
            {{ chatStore.currentOwner.description }}
          </p>
        </div>
      </div>
      
      <div ref="messagesContainer" class="chat-messages">
        <div 
          v-for="msg in chatStore.messages" 
          :key="msg.id"
          :class="['message', { self: isSelf(msg), other: !isSelf(msg) }]"
        >
          <div class="message-avatar">
            <img 
              v-if="chatStore.getAvatar(msg.sender_id)" 
              :src="chatStore.getAvatar(msg.sender_id)" 
              alt="头像"
            />
            <span v-else>{{ getInitials(chatStore.getNickname(msg.sender_id)) }}</span>
          </div>
          <div class="message-content">
            <div class="message-sender">{{ chatStore.getNickname(msg.sender_id) }}</div>
            <div class="message-bubble">
              <template v-if="msg.type === 1">
                <img :src="msg.content" class="message-image" @click="openImage(msg.content)" />
              </template>
              <template v-else>
                {{ msg.content }}
              </template>
            </div>
            <div class="message-time">{{ formatTime(msg.created_at) }}</div>
          </div>
        </div>
      </div>
      
      <div class="chat-input-area">
        <div class="emoji-container">
          <button class="emoji-btn" @click="showEmojiPicker = !showEmojiPicker">
            😊
          </button>
          <div v-if="showEmojiPicker" class="emoji-picker">
            <div 
              v-for="emoji in emojis" 
              :key="emoji"
              class="emoji-item"
              @click="insertEmoji(emoji)"
            >
              {{ emoji }}
            </div>
          </div>
        </div>
        <label class="upload-btn">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
            <polyline points="17 8 12 3 7 8" />
            <line x1="12" y1="3" x2="12" y2="15" />
          </svg>
          <input 
            type="file" 
            accept="image/*" 
            @change="handleImageUpload"
            style="display: none"
          />
        </label>
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

<style scoped>
.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.message {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.message.self {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  flex-shrink: 0;
  overflow: hidden;
}

.message-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.message-image {
  max-width: 300px;
  max-height: 300px;
  border-radius: 8px;
  cursor: pointer;
}

.emoji-container {
  position: relative;
}

.emoji-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  padding: 8px;
}

.emoji-btn:hover {
  background: #f0f0f0;
  border-radius: 50%;
}

.emoji-picker {
  position: absolute;
  bottom: 100%;
  left: 0;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  padding: 12px;
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 4px;
  z-index: 100;
}

.emoji-item {
  padding: 8px;
  font-size: 20px;
  cursor: pointer;
  text-align: center;
  border-radius: 6px;
}

.emoji-item:hover {
  background: #f0f0f0;
}

.upload-btn {
  background: none;
  border: none;
  color: #666;
  font-size: 24px;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-btn:hover {
  background: #f0f0f0;
  border-radius: 50%;
}

.chat-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.owner-badge {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  color: white;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}
</style>
