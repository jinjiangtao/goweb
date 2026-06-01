<script setup lang="ts">
import { ref } from 'vue'
import { useChatStore } from '@/stores/chat'

const chatStore = useChatStore()

const isLogin = ref(true)
const username = ref('')
const password = ref('')
const nickname = ref('')
const errorMessage = ref('')

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
      await chatStore.register(username.value, password.value, nickname.value)
      await chatStore.login(username.value, password.value)
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