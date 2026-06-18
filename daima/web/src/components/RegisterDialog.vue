<template>
  <el-dialog v-model="visible" title="用户注册" width="420px" :close-on-click-modal="false">
    <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" placeholder="3-50个字符" size="large" />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" placeholder="选填" size="large" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input v-model="form.password" type="password" placeholder="至少6个字符" size="large" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false" size="large">取消</el-button>
      <el-button type="primary" :loading="loading" size="large" @click="handleRegister">注册</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { authApi } from '../api'
import { useUserStore } from '../store/user'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible', 'success'])
const userStore = useUserStore()

const visible = computed({
  get: () => props.visible,
  set: val => emit('update:visible', val)
})

const formRef = ref(null)
const loading = ref(false)
const form = reactive({ username: '', email: '', password: '' })
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '长度在 6 到 50 个字符', trigger: 'blur' }
  ]
}

watch(() => props.visible, (v) => v && formRef.value?.resetFields())

const handleRegister = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  try {
    loading.value = true
    const res = await authApi.register(form)
    userStore.setAuth(res.data.user, res.data.token)
    ElMessage.success('注册成功')
    emit('success')
    visible.value = false
  } finally {
    loading.value = false
  }
}
</script>
