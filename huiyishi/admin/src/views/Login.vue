<template>
  <div class="login-container">
    <div class="login-box">
      <h1>会议室预订系统</h1>
      <h2>管理后台</h2>
      <el-form :model="form" @submit.prevent="handleLogin">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码" prefix-icon="Lock" size="large" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="handleLogin">登录</el-button>
        </el-form-item>
      </el-form>
      <div class="hint">默认账号: admin / 123456</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: ''
})
const loading = ref(false)

const handleLogin = async () => {
  console.log('开始登录...')
  console.log('用户名:', form.value.username)
  console.log('密码:', form.value.password)
  
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    console.log('发送请求到: http://localhost:8080/api/admin/login')
    const res = await axios.post('http://localhost:8080/api/admin/login', form.value)
    console.log('响应:', res)
    
    authStore.setAuth(res.data.token, res.data.admin)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('登录错误:', error)
    ElMessage.error(error.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, transparent 70%);
  animation: rotate 30s linear infinite;
}

@keyframes rotate {
  from { transform: translate(-25%, -25%) rotate(0deg); }
  to { transform: translate(-25%, -25%) rotate(360deg); }
}

.login-box {
  background: rgba(30, 41, 59, 0.9);
  backdrop-filter: blur(10px);
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 400px;
  border: 1px solid rgba(71, 85, 105, 0.5);
  position: relative;
  z-index: 1;
}

.login-box h1 {
  text-align: center;
  font-size: 24px;
  margin-bottom: 8px;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-box h2 {
  text-align: center;
  font-size: 16px;
  color: #94a3b8;
  margin-bottom: 32px;
  font-weight: normal;
}

.hint {
  text-align: center;
  color: #64748b;
  font-size: 12px;
  margin-top: 16px;
}
</style>
