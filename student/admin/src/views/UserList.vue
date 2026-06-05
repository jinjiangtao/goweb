<template>
  <div class="user-list">
    <div class="search-bar">
      <el-input v-model="searchForm.username" placeholder="按用户名搜索" class="search-input" />
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button type="success" @click="openAddModal">添加用户</el-button>
    </div>

    <el-table :data="tableData" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="角色" width="120">
        <template #default="scope">
          <el-tag :type="scope.row.role === 'super_admin' ? 'danger' : 'primary'">
            {{ scope.row.role === 'super_admin' ? '超级管理员' : '普通管理员' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
            {{ scope.row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_login_time" label="最后登录时间" width="180">
        <template #default="scope">
          {{ scope.row.last_login_time ? formatDate(scope.row.last_login_time) : '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button type="text" @click="openEditModal(scope.row)">编辑</el-button>
          <el-button type="text" @click="openResetModal(scope.row)">重置密码</el-button>
          <el-button v-if="scope.row.id !== currentUserId" type="text" danger @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="pagination.page"
      :page-sizes="[10, 20, 50, 100]"
      :page-size="pagination.pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="pagination.total"
      class="pagination"
    />

    <el-dialog :title="addForm.id ? '编辑用户' : '添加用户'" :visible.sync="addModalVisible" width="400px">
      <el-form :model="addForm" ref="addFormRef" label-width="80px">
        <el-form-item v-if="!addForm.id" label="用户名" prop="username">
          <el-input v-model="addForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item v-if="!addForm.id" label="密码" prop="password">
          <el-input v-model="addForm.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="addForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="addForm.role">
            <el-option label="超级管理员" value="super_admin" />
            <el-option label="普通管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="addForm.id" label="状态" prop="status">
          <el-select v-model="addForm.status">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addModalVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog title="重置密码" :visible.sync="resetModalVisible" width="400px">
      <el-form :model="resetForm" ref="resetFormRef" label-width="80px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="resetForm.password" type="password" placeholder="请输入新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetModalVisible = false">取消</el-button>
        <el-button type="primary" @click="handleResetPassword">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const currentUserId = authStore.user?.id

const searchForm = reactive({
  username: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const tableData = ref([])
const addModalVisible = ref(false)
const resetModalVisible = ref(false)
const addFormRef = ref()
const resetFormRef = ref()

const addForm = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  role: 'admin',
  status: 1
})

const resetForm = reactive({
  password: ''
})

const resetUserId = ref(0)

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const loadData = async () => {
  try {
    const response = await axios.get('/api/admin/users', {
      params: {
        username: searchForm.username,
        page: pagination.page,
        pageSize: pagination.pageSize
      }
    })
    tableData.value = response.data.data
    pagination.total = response.data.total
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  loadData()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  loadData()
}

const openAddModal = () => {
  addForm.id = 0
  addForm.username = ''
  addForm.password = ''
  addForm.nickname = ''
  addForm.role = 'admin'
  addForm.status = 1
  addModalVisible.value = true
}

const openEditModal = (row) => {
  addForm.id = row.id
  addForm.username = row.username
  addForm.password = ''
  addForm.nickname = row.nickname
  addForm.role = row.role
  addForm.status = row.status
  addModalVisible.value = true
}

const openResetModal = (row) => {
  resetUserId.value = row.id
  resetForm.password = ''
  resetModalVisible.value = true
}

const handleSave = async () => {
  if (!addForm.id && (!addForm.username || !addForm.password)) {
    ElMessage.warning('请填写用户名和密码')
    return
  }
  if (!addForm.nickname) {
    ElMessage.warning('请填写昵称')
    return
  }

  try {
    if (addForm.id) {
      await axios.put(`/api/admin/users/${addForm.id}`, {
        nickname: addForm.nickname,
        role: addForm.role,
        status: addForm.status
      })
      ElMessage.success('更新成功')
    } else {
      await axios.post('/api/admin/users', {
        username: addForm.username,
        password: addForm.password,
        nickname: addForm.nickname,
        role: addForm.role
      })
      ElMessage.success('添加成功')
    }
    addModalVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

const handleResetPassword = async () => {
  if (!resetForm.password) {
    ElMessage.warning('请输入新密码')
    return
  }

  try {
    await axios.post(`/api/admin/users/${resetUserId.value}/reset-password`, {
      password: resetForm.password
    })
    ElMessage.success('重置密码成功')
    resetModalVisible.value = false
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '重置密码失败')
  }
}

const handleDelete = async (row) => {
  try {
    await axios.delete(`/api/admin/users/${row.id}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.user-list {
  padding: 20px;
}
.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}
.search-input {
  width: 200px;
}
.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>