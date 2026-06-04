<template>
  <div class="login-container">
    <div class="login-box">
      <h2>学生报名管理系统</h2>
      <el-form :model="form" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading" style="width: 100%">登录</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { ElMessage } from 'element-plus';
const router = useRouter();
const authStore = useAuthStore();
const formRef = ref();
const loading = ref(false);
const form = reactive({
 username: '',
 password: ''
});
const handleLogin = async () => {
 if (!form.username || !form.password) {
 ElMessage.warning('请填写用户名和密码');
 return;
 }
 loading.value = true;
 try {
 await authStore.login(form.username, form.password);
 ElMessage.success('登录成功');
 router.push('/');
 }
 catch (error) {
 ElMessage.error(error.response?.data?.error || '登录失败');
 }
 finally {
 loading.value = false;
 }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-box {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  width: 400px;
}
.login-box h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}
</style>
