<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="openAddDialog">添加用户</el-button>
        </div>
      </template>
      <el-table :data="users" stripe>
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="real_name" label="姓名" />
        <el-table-column prop="role" label="角色">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : row.role === 'hr' ? 'warning' : 'primary'">
              {{ row.role === 'admin' ? '管理员' : row.role === 'hr' ? 'HR' : '员工' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === 'enabled' ? 'success' : 'danger'">
              {{ row.status === 'enabled' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="toggleStatus(row)">
              {{ row.status === 'enabled' ? '禁用' : '启用' }}
            </el-button>
            <el-button type="warning" link @click="resetPassword(row)">重置密码</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="addDialogVisible" title="添加用户" width="500px">
      <el-form :model="addForm" label-width="100px">
        <el-form-item label="用户名">
          <el-input v-model="addForm.username" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="addForm.real_name" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="addForm.role" style="width: 100%">
            <el-option label="员工" value="employee" />
            <el-option label="HR" value="hr" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAdd" :loading="loading">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const users = ref([])
const addDialogVisible = ref(false)
const loading = ref(false)
const addForm = ref({
  username: '',
  real_name: '',
  role: 'employee'
})

const loadUsers = async () => {
  const res = await api.get('/admin/users')
  users.value = res
}

const openAddDialog = () => {
  addForm.value = { username: '', real_name: '', role: 'employee' }
  addDialogVisible.value = true
}

const handleAdd = async () => {
  if (!addForm.value.username || !addForm.value.real_name) {
    ElMessage.warning('请填写必填项')
    return
  }
  loading.value = true
  try {
    await api.post('/admin/users', addForm.value)
    ElMessage.success('添加成功，初始密码：123456')
    addDialogVisible.value = false
    loadUsers()
  } catch (e) {
    ElMessage.error('添加失败')
  } finally {
    loading.value = false
  }
}

const toggleStatus = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要${row.status === 'enabled' ? '禁用' : '启用'}该用户吗？`, '提示')
    const newStatus = row.status === 'enabled' ? 'disabled' : 'enabled'
    await api.put(`/admin/users/${row.id}/status`, { status: newStatus })
    ElMessage.success('更新成功')
    loadUsers()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

const resetPassword = async (row) => {
  try {
    await ElMessageBox.confirm('确定要重置密码吗？重置后密码为：123456', '提示')
    await api.put(`/admin/users/${row.id}/reset-password`)
    ElMessage.success('密码已重置为 123456')
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
