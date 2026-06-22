<template>
  <Layout>
    <div class="shares-page">
      <div class="page-header">
        <h3>分享管理</h3>
      </div>

      <el-table v-loading="loading" :data="shares" stripe>
        <el-table-column prop="album.name" label="相册名称" min-width="180">
          <template #default="{ row }">
            <div class="album-name-cell">
              <el-icon color="#409EFF"><Folder /></el-icon>
              {{ row.album?.name || '未知' }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="token" label="分享令牌" min-width="200">
          <template #default="{ row }">
            <span class="token-text">{{ row.token }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="访问次数" width="100" />
        <el-table-column prop="max_views" label="最大访问" width="100">
          <template #default="{ row }">
            {{ row.max_views || '不限' }}
          </template>
        </el-table-column>
        <el-table-column prop="is_active" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
              {{ row.is_active ? '有效' : '失效' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="password" label="密码保护" width="100">
          <template #default="{ row }">
            <el-icon v-if="row.password" color="#e6a23c"><Lock /></el-icon>
            <span v-else class="muted">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="expires_at" label="过期时间" width="180">
          <template #default="{ row }">
            {{ row.expires_at ? formatDate(row.expires_at) : '永久有效' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="copyLink(row)">
              <el-icon><CopyDocument /></el-icon>
              复制链接
            </el-button>
            <el-button link type="primary" @click="viewVisits(row)">
              <el-icon><View /></el-icon>
              查看记录
            </el-button>
            <el-button link type="danger" @click="deleteShare(row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="shares.length === 0 && !loading" class="empty-state">
        <el-empty description="暂无分享链接" />
      </div>
    </div>

    <el-dialog v-model="visitDialogVisible" title="访客记录" width="800px">
      <el-table :data="visitRecords" stripe>
        <el-table-column prop="ip_address" label="IP地址" width="150" />
        <el-table-column prop="user_agent" label="浏览器信息">
          <template #default="{ row }">
            <span class="ua-text text-ellipsis" :title="row.user_agent">{{ row.user_agent }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="visited_at" label="访问时间" width="200">
          <template #default="{ row }">
            {{ formatDate(row.visited_at) }}
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <div class="dialog-footer">
          <span class="visit-count">共 {{ visitRecords.length }} 条记录</span>
          <el-button @click="visitDialogVisible = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '@/components/Layout.vue'
import { getShareLinks, deleteShareLink, getVisitRecords } from '@/utils/api'

const shares = ref([])
const loading = ref(false)
const visitDialogVisible = ref(false)
const visitRecords = ref([])

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const loadShares = async () => {
  loading.value = true
  try {
    const res = await getShareLinks()
    shares.value = res.shares
  } catch {
  } finally {
    loading.value = false
  }
}

const copyLink = async (row) => {
  const link = `${window.location.origin}/share/${row.token}`
  try {
    await navigator.clipboard.writeText(link)
    ElMessage.success('链接已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

const viewVisits = async (row) => {
  try {
    const res = await getVisitRecords(row.id)
    visitRecords.value = res.records
    visitDialogVisible.value = true
  } catch {}
}

const deleteShare = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个分享链接吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteShareLink(row.id)
    ElMessage.success('删除成功')
    loadShares()
  } catch {}
}

onMounted(() => {
  loadShares()
})
</script>

<style scoped>
.shares-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.album-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.token-text {
  font-family: monospace;
  font-size: 12px;
  color: #606266;
}

.muted {
  color: #c0c4cc;
}

.ua-text {
  display: block;
  max-width: 400px;
  font-size: 12px;
  color: #909399;
}

.empty-state {
  padding: 60px 0;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.visit-count {
  font-size: 14px;
  color: #909399;
}

.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
