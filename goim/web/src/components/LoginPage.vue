<script setup lang="ts">
import { ref } from 'vue'
import { useChatStore } from '@/stores/chat'
import * as api from '@/api'

const chatStore = useChatStore()

const isLogin = ref(true)
const username = ref('')
const password = ref('')
const nickname = ref('')
const errorMessage = ref('')
const avatarPreview = ref('')
const avatarFile = ref<File | null>(null)
const avatarInputRef = ref<HTMLInputElement | null>(null)

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

async function handleSubmit() {
  errorMessage.value = ''
  
  if (!username.value || !password.value) {
    errorMessage.value = '请填写用户名和密码'
    return
  }

  if (!isLogin.value && !nickname.value) {
    errorMessage.value = '请填写昵称'
    return
  }

  try {
    if (isLogin.value) {
      await chatStore.login(username.value, password.value)
    } else {
      let avatarUrl = ''
      // 如果有头像文件，先上传
      if (avatarFile.value) {
        // 先注册用户获取ID
        const user = await api.register({ 
          username: username.value, 
          password: password.value, 
          nickname: nickname.value 
        })
        // 然后登录获取token
        await api.login({ username: username.value, password: password.value })
        // 再上传头像
        const avatarResult = await api.uploadAvatar(user.id, avatarFile.value)
        avatarUrl = avatarResult.avatar
        // 更新用户信息
        await api.updateProfile(user.id, nickname.value, avatarUrl)
        // 重新登录加载完整用户信息
        await chatStore.login(username.value, password.value)
      } else {
        // 没有头像，直接注册登录
        await chatStore.register(username.value, password.value, nickname.value)
        await chatStore.login(username.value, password.value)
      }
    }
  } catch (error) {
    errorMessage.value = (error as Error).message
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-form">
      <h2>{{ isLogin ? '登录' : '注册' }}</h2>
      
      <div v-if="!isLogin" class="avatar-upload">
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
      </div>
      
      <div class="form-group">
        <label>用户名</label>
        <input v-model="username" type="text" placeholder="请输入用户名" />
      </div>

      <div class="form-group">
        <label>密码</label>
        <input v-model="password" type="password" placeholder="请输入密码" />
      </div>

      <div v-if="!isLogin" class="form-group">
        <label>昵称</label>
        <input v-model="nickname" type="text" placeholder="请输入昵称" />
      </div>

      <button class="btn btn-primary" @click="handleSubmit">
        {{ isLogin ? '登录' : '注册' }}
      </button>

      <button class="btn btn-secondary" @click="isLogin = !isLogin">
        {{ isLogin ? '切换到注册' : '切换到登录' }}
      </button>

      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
    </div>
  </div>
</template>

<style scoped>
.avatar-upload {
  margin-bottom: 20px;
  text-align: center;
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
  border: 2px dashed #667eea;
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

.avatar-upload:hover .avatar-placeholder {
  border-color: #764ba2;
  background-color: #f8f8ff;
}
</style>