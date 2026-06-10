
&lt;template&gt;
  &lt;div class="login-container"&gt;
    &lt;el-card class="login-card"&gt;
      &lt;h2&gt;抽奖转盘后台管理&lt;/h2&gt;
      &lt;el-form :model="loginForm" label-width="80px"&gt;
        &lt;el-form-item label="账号"&gt;
          &lt;el-input v-model="loginForm.username" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="密码"&gt;
          &lt;el-input v-model="loginForm.password" type="password" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item&gt;
          &lt;el-button type="primary" @click="handleLogin" style="width: 100%"&gt;登录&lt;/el-button&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
    &lt;/el-card&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loginForm = ref({
  username: '',
  password: ''
})

const handleLogin = async () =&gt; {
  try {
    const res = await axios.post('http://localhost:8080/api/login', loginForm.value)
    localStorage.setItem('token', res.data.token)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (err) {
    ElMessage.error('登录失败')
  }
}
&lt;/script&gt;

&lt;style scoped&gt;
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: #f0f2f5;
}
.login-card {
  width: 400px;
}
.login-card h2 {
  text-align: center;
  margin-bottom: 20px;
}
&lt;/style&gt;
