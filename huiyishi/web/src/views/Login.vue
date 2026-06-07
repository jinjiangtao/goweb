<template>
  <div class="login-container">
    <div class="login-card">
      <h1>会议室预订系统</h1>
      <div class="tabs">
        <button 
          :class="['tab', { active: !isRegister }]" 
          @click="isRegister = false"
        >
          登录
        </button>
        <button 
          :class="['tab', { active: isRegister }]" 
          @click="isRegister = true"
        >
          注册
        </button>
      </div>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>用户名</label>
          <input 
            v-model="form.username" 
            type="text" 
            placeholder="请输入用户名" 
            required
          />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input 
            v-model="form.password" 
            type="password" 
            placeholder="请输入密码" 
            required
          />
        </div>
        <div v-if="isRegister" class="form-group">
          <label>确认密码</label>
          <input 
            v-model="form.confirmPassword" 
            type="password" 
            placeholder="请再次输入密码" 
            required
          />
        </div>
        <div v-if="isRegister" class="form-group">
          <label>真实姓名</label>
          <input 
            v-model="form.realName" 
            type="text" 
            placeholder="请输入真实姓名" 
            required
          />
        </div>
        <div v-if="isRegister" class="form-group">
          <label>手机号</label>
          <input 
            v-model="form.phone" 
            type="tel" 
            placeholder="请输入手机号" 
            maxlength="11"
            required
          />
        </div>
        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? '请稍候...' : (isRegister ? '注册' : '登录') }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
      <p v-if="!isRegister" class="hint">
        没有账号？<a @click="isRegister = true">立即注册</a>
      </p>
      <p v-if="isRegister" class="hint">
        已有账号？<a @click="isRegister = false">立即登录</a>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api, { setAuthToken } from '../api'

const router = useRouter()
const isRegister = ref(false)
const loading = ref(false)
const error = ref('')

const form = ref({
  username: '',
  password: '',
  confirmPassword: '',
  realName: '',
  phone: ''
})

const handleSubmit = async () => {
  error.value = ''
  loading.value = true
  
  try {
    if (isRegister.value) {
      if (form.value.password !== form.value.confirmPassword) {
        error.value = '两次密码输入不一致'
        return
      }
      if (form.value.password.length < 6) {
        error.value = '密码长度至少6位'
        return
      }
      if (form.value.phone.length !== 11) {
        error.value = '请输入正确的手机号'
        return
      }
      
      await api.post('/register', {
        username: form.value.username,
        password: form.value.password,
        real_name: form.value.realName,
        phone: form.value.phone
      })
      
      alert('注册成功，请登录')
      isRegister.value = false
    } else {
      const res = await api.post('/login', {
        username: form.value.username,
        password: form.value.password
      })
      
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      setAuthToken(res.data.token)
      
      router.push('/')
    }
  } catch (err) {
    error.value = err.response?.data?.error || '请求失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0,0,0,0.2);
  width: 100%;
  max-width: 400px;
}

.login-card h1 {
  text-align: center;
  color: #333;
  margin-bottom: 30px;
  font-size: 24px;
}

.tabs {
  display: flex;
  margin-bottom: 30px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

.tab {
  flex: 1;
  padding: 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.3s;
}

.tab.active {
  background: #667eea;
  color: white;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #666;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  box-sizing: border-box;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.submit-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.submit-btn:hover {
  transform: translateY(-2px);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error {
  color: #e74c3c;
  text-align: center;
  margin-top: 16px;
}

.hint {
  text-align: center;
  margin-top: 20px;
  color: #666;
}

.hint a {
  color: #667eea;
  cursor: pointer;
  text-decoration: underline;
}
</style>
