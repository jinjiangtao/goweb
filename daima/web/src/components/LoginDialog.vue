<template>
  <el-dialog v-model="visible" title="用户登录" width="420px" :close-on-click-modal="false">
    <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" placeholder="请输入用户名" size="large" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input v-model="form.password" type="password" placeholder="请输入密码" size="large" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false" size="large">取消</el-button>
      <el-button type="primary" :loading="loading" size="large" @click="handleLogin">登录</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { authApi } from '../api'
import { useUserStore } from '../store/user'

const props = defineProps({
  visible: Boolean
})
const emit = defineEmits(['update:visible', 'success'])
const userStore = useUserStore()

const visible = computed({
  get: () => props.visible,
  set: val => emit('update:visible', val)
})

const formRef = ref(null)
const loading = ref(false)
const form = reactive({
  username: '',
  password: ''
})
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

watch(() => props.visible, (v) => {
  if (v) {
    formRef.value?.resetFields()
  }
})

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  try {
    loading.value = true
    const res = await authApi.login(form)
    userStore.setAuth(res.data.user, res.data.token)
    ElMessage.success('登录成功')
    emit('success')
    visible.value = false
  } finally {
    loading.value = false
  }
}
</script>
