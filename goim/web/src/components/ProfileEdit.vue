<script setup lang="ts">
import { ref, watch } from 'vue'
import { useChatStore } from '@/stores/chat'
import * as api from '@/api'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits(['close'])

const chatStore = useChatStore()
const nickname = ref('')
const avatarPreview = ref('')
const avatarFile = ref<File | null>(null)
const loading = ref(false)
const errorMessage = ref('')
const avatarInputRef = ref<HTMLInputElement | null>(null)

watch(() => props.visible, (newVal) => {
  if (newVal && chatStore.currentUser) {
    nickname.value = chatStore.currentUser.nickname
    avatarPreview.value = chatStore.currentUser.avatar || ''
  }
})

function handleAvatarClick() {
  avatarInputRef.value?.click()
}

function handleAvatarChange(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    avatarFile.value = target.files[0]
    const reader = new FileReader()
    reader.onload = (e) => {
      avatarPreview.value = e.target?.result as string
    }
    reader.readAsDataURL(avatarFile.value)
  }
}

async function handleSave() {
  if (!nickname.value.trim()) {
    errorMessage.value = '昵称不能为空'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    let avatarUrl = avatarPreview.value
    if (avatarFile.value && chatStore.userId) {
      const result = await api.uploadAvatar(chatStore.userId, avatarFile.value)
      avatarUrl = result.avatar
    }

    await chatStore.updateProfile(nickname.value, avatarUrl)
    emit('close')
  } catch (error) {
    errorMessage.value = (error as Error).message
  } finally {
    loading.value = false
  }
}

function handleClose() {
  emit('close')
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleClose">
    <div class="modal-content">
      <div class="modal-header">
        <h3>编辑个人资料</h3>
        <button class="close-btn" @click="handleClose">×</button>
      </div>

      <div class="modal-body">
        <div class="avatar-edit">
          <label>头像</label>
          <div class="avatar-preview" v-if="avatarPreview" @click="handleAvatarClick">
            <img :src="avatarPreview" alt="头像预览" />
          </div>
          <div v-else class="avatar-placeholder" @click="handleAvatarClick">
            点击选择头像
          </div>
          <input 
            ref="avatarInputRef"
            type="file" 
            accept="image/*" 
            @change="handleAvatarChange" 
            class="avatar-input"
          />
          <label class="change-avatar-btn" @click="handleAvatarClick">更换头像</label>
        </div>

        <div class="form-group">
          <label>昵称</label>
          <input v-model="nickname" type="text" placeholder="请输入昵称" />
        </div>

        <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      </div>

      <div class="modal-footer">
        <button class="btn btn-secondary" @click="handleClose" :disabled="loading">
          取消
        </button>
        <button class="btn btn-primary" @click="handleSave" :disabled="loading">
          {{ loading ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 400px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  padding: 20px;
  border-top: 1px solid #eee;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.avatar-edit {
  text-align: center;
  margin-bottom: 20px;
}

.avatar-edit label {
  display: block;
  margin-bottom: 10px;
  font-weight: bold;
}

.avatar-preview,
.avatar-placeholder {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  margin: 10px auto;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
  border: 2px solid #667eea;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  background-color: #f0f0f0;
  color: #999;
  font-size: 14px;
}

.avatar-input {
  display: none;
}

.change-avatar-btn {
  display: inline-block;
  background: #667eea;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  margin-top: 10px;
}

.change-avatar-btn:hover {
  background: #764ba2;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 16px;
}

.error-message {
  color: #e74c3c;
  margin-top: 10px;
}
</style>
