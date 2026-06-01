<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { ElMessage, ElTable, ElTableColumn, ElPagination, ElButton, ElTag } from 'element-plus'

const router = useRouter()
const adminStore = useAdminStore()

const currentPage = ref(1)
const pageSize = ref(10)
let refreshInterval: ReturnType<typeof setInterval> | null = null

async function handlePageChange(page: number) {
  currentPage.value = page
  await adminStore.fetchUsers(page, pageSize.value)
}

async function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  await adminStore.fetchUsers(1, size)
}

async function handleRefresh() {
  await adminStore.fetchUsers(currentPage.value, pageSize.value)
  ElMessage.success('刷新成功')
}

function handleLogout() {
  adminStore.logout()
  router.push('/admin')
}

watch(() => adminStore.isLoggedIn, (loggedIn) => {
  if (!loggedIn) {
    router.push('/admin')
  }
})

onMounted(() => {
  if (!adminStore.isLoggedIn) {
    router.push('/admin')
    return
  }
  adminStore.fetchUsers(currentPage.value, pageSize.value)
  
  // 每5秒刷新一次在线状态
  refreshInterval = setInterval(() => {
    adminStore.fetchUsers(currentPage.value, pageSize.value)
  }, 5000)
})

// 组件卸载时清除定时器
import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<template>
  <div class="admin-users-container">
    <div class="admin-header">
      <h1>GoIM 管理后台 - 用户管理</h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleRefresh" :loading="adminStore.loading">
          刷新
        </el-button>
        <el-button @click="handleLogout">
          退出
        </el-button>
      </div>
    </div>
    
    <div class="admin-content">
      <el-table :data="adminStore.users" v-loading="adminStore.loading" style="width: 100%">
        <el-table-column prop="id" label="用户ID" width="200" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="nickname" label="昵称" width="150" />
        <el-table-column prop="avatar" label="头像" width="100">
          <template #default="{ row }">
            <div v-if="row.avatar" class="avatar-preview">
              <img :src="row.avatar" :alt="row.nickname" />
            </div>
            <div v-else class="avatar-placeholder">
              {{ row.nickname.charAt(0).toUpperCase() }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="online" label="在线状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.online ? 'success' : 'info'">
              {{ row.online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" />
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="adminStore.pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-users-container {
  min-height: 100vh;
  background: #f0f2f5;
}

.admin-header {
  background: white;
  padding: 20px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.admin-header h1 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.admin-content {
  padding: 20px 40px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.avatar-preview {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
}
</style>
