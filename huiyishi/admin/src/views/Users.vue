<template>
  <div class="users">
    <div class="header">
      <h2>用户管理</h2>
      <el-button type="primary" @click="openCreateDialog">添加用户</el-button>
    </div>

    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索用户名、姓名或手机号"
        clearable
        style="width: 300px"
        @keyup.enter="fetchUsers"
      >
        <template #append>
          <el-button @click="fetchUsers">搜索</el-button>
        </template>
      </el-input>
    </div>

    <el-table :data="users" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="real_name" label="真实姓名" width="120" />
      <el-table-column prop="phone" label="手机号" width="130" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" fixed="right" width="280">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button type="warning" size="small" @click="toggleStatus(row)">
            {{ row.status === 1 ? '禁用' : '启用' }}
          </el-button>
          <el-button type="info" size="small" @click="resetPassword(row)">重置密码</el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="fetchUsers"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '添加用户'"
      width="450px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="form.real_name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" maxlength="11" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const users = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')

const form = reactive({
  id: null,
  username: '',
  password: '',
  real_name: '',
  phone: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, message: '密码至少6位', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }, { pattern: /^\d{11}$/, message: '请输入11位手机号', trigger: 'blur' }]
}

onMounted(() => {
  fetchUsers()
})

async function fetchUsers() {
  loading.value = true
  try {
    const res = await api.get('/users', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        search: searchKeyword.value
      }
    })
    users.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  isEdit.value = false
  Object.assign(form, { id: null, username: '', password: '', real_name: '', phone: '' })
  dialogVisible.value = true
}

function openEditDialog(row) {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    username: row.username,
    real_name: row.real_name,
    phone: row.phone
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await api.put(`/users/${form.id}`, {
        real_name: form.real_name,
        phone: form.phone
      })
      ElMessage.success('编辑成功')
    } else {
      await api.post('/users', form)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchUsers()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await ElMessageBox.confirm(
      `确定要${row.status === 1 ? '禁用' : '启用'}该用户吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return
  }

  try {
    await api.put(`/users/${row.id}/toggle-status`)
    ElMessage.success('操作成功')
    fetchUsers()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

async function resetPassword(row) {
  try {
    await ElMessageBox.confirm(
      '确定要重置密码吗？重置后密码为 123456',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return
  }

  try {
    await api.put(`/users/${row.id}/reset-password`)
    ElMessage.success('密码已重置为 123456')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      '确定要删除该用户吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return
  }

  try {
    await api.delete(`/users/${row.id}`)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}
</script>

<style scoped>
.users {
  background: #1e293b;
  border-radius: 8px;
  padding: 24px;
  min-height: 100%;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  color: #e2e8f0;
  margin: 0;
}

.search-bar {
  margin-bottom: 20px;
}

:deep(.el-table) {
  background: #1e293b;
  color: #e2e8f0;
}

:deep(.el-table th) {
  background: #334155;
  color: #cbd5e1;
}

:deep(.el-table tr) {
  background: #1e293b;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background: #27344b;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
