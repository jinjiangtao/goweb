<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { ElMessage, ElTable, ElTableColumn, ElPagination, ElButton, ElTag, ElAvatar } from 'element-plus'

const router = useRouter()
const adminStore = useAdminStore()

const currentPage = ref(1)
const pageSize = ref(20)

async function handlePageChange(page: number) {
  currentPage.value = page
  await adminStore.fetchMessages(page, pageSize.value)
}

async function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  await adminStore.fetchMessages(1, size)
}

async function handleRefresh() {
  await adminStore.fetchMessages(currentPage.value, pageSize.value)
  ElMessage.success('刷新成功')
}

function handleLogout() {
  adminStore.logout()
  router.push('/admin')
}

function goToUsers() {
  router.push('/admin/users')
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
  adminStore.fetchMessages(currentPage.value, pageSize.value)
})
</script>

<template>
  <div class="admin-messages-container">
    <div class="admin-header">
      <h1>GoIM 管理后台 - 聊天记录</h1>
      <div class="header-actions">
        <el-button @click="goToUsers">
          用户管理
        </el-button>
        <el-button type="primary" @click="handleRefresh" :loading="adminStore.messagesLoading">
          刷新
        </el-button>
        <el-button @click="handleLogout">
          退出
        </el-button>
      </div>
    </div>
    
    <div class="admin-content">
      <el-table :data="adminStore.messages" v-loading="adminStore.messagesLoading" style="width: 100%">
        <el-table-column prop="id" label="消息ID" width="180" />
        <el-table-column label="发送者" width="160">
          <template #default="{ row }">
            <div class="sender-info">
              <el-avatar :size="32" :src="row.sender?.avatar">
                {{ row.sender?.nickname?.charAt(0).toUpperCase() }}
              </el-avatar>
              <span class="sender-name">{{ row.sender?.nickname }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="接收者" width="160">
          <template #default="{ row }">
            <div class="receiver-info">
              <el-avatar v-if="row.receiver?.type === 'user'" :size="32" :src="row.receiver?.avatar">
                {{ row.receiver?.nickname?.charAt(0).toUpperCase() }}
              </el-avatar>
              <el-avatar v-else :size="32">
                {{ row.receiver?.name?.charAt(0).toUpperCase() }}
              </el-avatar>
              <span class="receiver-name">
                {{ row.receiver?.nickname || row.receiver?.name }}
                <el-tag v-if="row.receiver?.type === 'group'" size="small" type="info" style="margin-left: 8px">群组</el-tag>
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="消息内容" min-width="200">
          <template #default="{ row }">
            <template v-if="row.type === 1">
              <div class="image-preview">
                <img :src="row.content" alt="图片" />
              </div>
            </template>
            <template v-else>
              {{ row.content }}
            </template>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="消息类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 1 ? 'warning' : 'primary'">
              {{ row.type === 1 ? '图片' : '文字' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发送时间" width="180" />
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100]"
          :total="adminStore.messagesPagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-messages-container {
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

.sender-info,
.receiver-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sender-name,
.receiver-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.image-preview {
  max-width: 120px;
  max-height: 80px;
  overflow: hidden;
  border-radius: 4px;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
