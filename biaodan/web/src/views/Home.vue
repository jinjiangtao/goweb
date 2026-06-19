<template>
  <div class="home-container">
    <header class="header">
      <div class="header-inner">
        <div class="logo">
          <el-icon :size="28" color="#409EFF"><Tickets /></el-icon>
          <h1>低代码表单设计平台</h1>
        </div>
        <div class="header-actions">
          <el-input
            v-model="keyword"
            placeholder="搜索表单名称..."
            clearable
            :prefix-icon="Search"
            style="width: 260px; margin-right: 16px"
            @keyup.enter="loadList"
            @clear="loadList"
          />
          <el-button type="primary" :icon="Plus" @click="goDesigner">
            新建表单
          </el-button>
        </div>
      </div>
    </header>

    <main class="main-content">
      <el-tabs v-model="activeTab" class="template-tabs">
        <el-tab-pane label="全部表单" name="all">
          <div v-loading="loading" class="template-grid">
            <template v-if="templates.length">
              <div
                v-for="item in templates"
                :key="item.id"
                class="template-card"
                @click="viewData(item.id)"
              >
                <div class="card-header">
                  <div class="card-title" :title="item.name">{{ item.name }}</div>
                  <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, item)">
                    <el-button :icon="More" link type="primary" @click.stop />
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="edit">
                          <el-icon><Edit /></el-icon>编辑表单
                        </el-dropdown-item>
                        <el-dropdown-item command="fill">
                          <el-icon><EditPen /></el-icon>填写表单
                        </el-dropdown-item>
                        <el-dropdown-item command="data">
                          <el-icon><DataLine /></el-icon>查看数据
                        </el-dropdown-item>
                        <el-dropdown-item command="link" divided>
                          <el-icon><Link /></el-icon>复制链接
                        </el-dropdown-item>
                        <el-dropdown-item command="delete" divided>
                          <el-icon><Delete /></el-icon>删除表单
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
                <p class="card-desc" :title="item.description">{{ item.description || '暂无描述' }}</p>
                <div class="card-footer">
                  <div class="stat">
                    <el-icon color="#909399"><Document /></el-icon>
                    <span>{{ item.submissions }} 条数据</span>
                  </div>
                  <span class="card-time">{{ formatTime(item.createdAt) }}</span>
                </div>
              </div>
            </template>
            <el-empty v-else description="暂无表单，点击右上角新建吧~" />
          </div>
        </el-tab-pane>
      </el-tabs>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Search, More, Edit, EditPen, DataLine, Link, Delete, Tickets, Document
} from '@element-plus/icons-vue'
import { listTemplates, deleteTemplate } from '@/api'

const router = useRouter()
const loading = ref(false)
const keyword = ref('')
const activeTab = ref('all')
const templates = ref([])

async function loadList() {
  loading.value = true
  try {
    const res = await listTemplates({ keyword: keyword.value })
    templates.value = res.data || []
  } finally {
    loading.value = false
  }
}

function goDesigner() {
  router.push('/designer')
}

function viewData(id) {
  router.push(`/data/${id}`)
}

function handleCommand(cmd, item) {
  switch (cmd) {
    case 'edit':
      router.push(`/designer/${item.id}`)
      break
    case 'fill':
      window.open(`/fill/${item.id}`, '_blank')
      break
    case 'data':
      router.push(`/data/${item.id}`)
      break
    case 'link':
      const url = `${window.location.origin}/fill/${item.id}`
      navigator.clipboard.writeText(url)
      ElMessage.success('链接已复制到剪贴板')
      break
    case 'delete':
      ElMessageBox.confirm('确定要删除此表单吗？所有数据也将被删除。', '删除确认', {
        type: 'warning'
      }).then(async () => {
        await deleteTemplate(item.id)
        ElMessage.success('删除成功')
        loadList()
      }).catch(() => {})
      break
  }
}

function formatTime(t) {
  if (!t) return ''
  return t.substring(0, 10)
}

onMounted(loadList)
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo h1 {
  font-size: 22px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
}

.main-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
}

.template-tabs {
  background: transparent;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  padding: 8px 0;
}

.template-card {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
}

.template-card:hover {
  border-color: #409EFF;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.12);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding-right: 8px;
}

.card-desc {
  color: #909399;
  font-size: 13px;
  line-height: 1.5;
  margin: 0 0 16px;
  flex: 1;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f2f5;
}

.stat {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #909399;
  font-size: 12px;
}

.card-time {
  color: #c0c4cc;
  font-size: 12px;
}
</style>
