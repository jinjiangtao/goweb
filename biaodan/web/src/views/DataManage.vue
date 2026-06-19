<template>
  <div class="data-container">
    <header class="data-header">
      <div class="header-inner">
        <div class="header-left">
          <el-button :icon="ArrowLeft" link @click="$router.push('/')">返回列表</el-button>
          <div class="header-info">
            <h1 class="template-name">{{ template?.name || '数据管理' }}</h1>
            <p class="template-desc" v-if="template?.description">{{ template.description }}</p>
          </div>
        </div>
        <div class="header-right">
          <el-button :icon="Edit" @click="goEdit">编辑表单</el-button>
          <el-button :icon="Link" @click="copyFillLink">复制填写链接</el-button>
          <el-button type="primary" :icon="Download" :loading="exporting" @click="exportData">
            导出CSV
          </el-button>
        </div>
      </div>
    </header>

    <div class="stats-section" v-if="stats">
      <div class="stat-cards">
        <div class="stat-card total">
          <div class="stat-icon"><el-icon :size="28"><DataLine /></el-icon></div>
          <div class="stat-content">
            <div class="stat-num">{{ stats.total }}</div>
            <div class="stat-label">总提交数</div>
          </div>
        </div>
        <div class="stat-card today">
          <div class="stat-icon"><el-icon :size="28"><Sunny /></el-icon></div>
          <div class="stat-content">
            <div class="stat-num">{{ stats.todayCount }}</div>
            <div class="stat-label">今日新增</div>
          </div>
        </div>
        <div class="stat-card chart-card">
          <div class="chart-title">近7天提交趋势</div>
          <div class="chart-bars">
            <div
              v-for="(day, idx) in stats.last7Days"
              :key="idx"
              class="chart-bar-item"
            >
              <div class="bar-wrap">
                <div
                  class="bar"
                  :style="{ height: getBarHeight(day.count) + '%' }"
                ></div>
              </div>
              <span class="bar-label">{{ formatDay(day.date) }}</span>
              <span class="bar-count">{{ day.count }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <main class="data-main">
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索数据..."
          :prefix-icon="Search"
          clearable
          style="width: 260px"
          @keyup.enter="loadData"
          @clear="loadData"
        />
        <div class="toolbar-right">
          <el-button :icon="Refresh" @click="loadData">刷新</el-button>
        </div>
      </div>

      <el-table
        v-loading="loading"
        :data="records"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column type="index" label="#" width="56" align="center" />
        <el-table-column
          v-for="field in fields"
          :key="field.id"
          :label="field.label"
          min-width="140"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <component :is="'div'" v-html="formatCellValue(row.data[field.id], field.type)" />
          </template>
        </el-table-column>
        <el-table-column label="提交时间" width="170" align="center">
          <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="IP" width="140" align="center">
          <template #default="{ row }">{{ row.ip || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="danger" link @click="deleteRecord(row.id)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Edit, Link, Download, Search, Refresh,
  DataLine, Sunny
} from '@element-plus/icons-vue'
import {
  getTemplate, listSubmissions, deleteSubmission,
  exportSubmissions, getTemplateStats
} from '@/api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const exporting = ref(false)
const template = ref(null)
const records = ref([])
const fields = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const stats = ref(null)

async function loadAll() {
  await Promise.all([loadTemplate(), loadStats(), loadData()])
}

async function loadTemplate() {
  const id = route.params.id
  try {
    const res = await getTemplate(id)
    template.value = res.data
  } catch (e) {}
}

async function loadStats() {
  const id = route.params.id
  try {
    const res = await getTemplateStats(id)
    stats.value = res.data
  } catch (e) {}
}

async function loadData() {
  const id = route.params.id
  loading.value = true
  try {
    const res = await listSubmissions({
      templateId: id,
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value
    })
    records.value = res.data.records || []
    fields.value = res.data.fields || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

function goEdit() {
  router.push(`/designer/${route.params.id}`)
}

function copyFillLink() {
  const url = `${window.location.origin}/fill/${route.params.id}`
  navigator.clipboard.writeText(url)
  ElMessage.success('填写链接已复制到剪贴板')
}

async function deleteRecord(id) {
  ElMessageBox.confirm('确定删除此条数据吗？', '删除确认', {
    type: 'warning'
  }).then(async () => {
    await deleteSubmission(id)
    ElMessage.success('删除成功')
    loadData()
    loadStats()
  }).catch(() => {})
}

async function exportData() {
  exporting.value = true
  try {
    const blob = await exportSubmissions({ templateId: route.params.id })
    const url = URL.createObjectURL(new Blob([blob], { type: 'text/csv;charset=utf-8' }))
    const link = document.createElement('a')
    link.href = url
    link.download = `${template.value?.name || '表单数据'}_${Date.now()}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } finally {
    exporting.value = false
  }
}

function formatCellValue(value, type) {
  if (value === undefined || value === null || value === '') return '-'
  if (type === 'checkbox' && Array.isArray(value)) {
    return value.map(v => {
      if (typeof v === 'object' && v.label) return v.label
      return String(v)
    }).join('、') || '-'
  }
  if (type === 'upload') {
    if (Array.isArray(value)) {
      return value.map(f => `<a href="${f.url}" target="_blank">${f.name || '附件'}</a>`).join('<br>')
    } else if (typeof value === 'object' && value.url) {
      return `<a href="${value.url}" target="_blank">${value.name || '附件'}</a>`
    }
  }
  if (typeof value === 'object' && value.url) {
    return `<a href="${value.url}" target="_blank">${value.name || '附件'}</a>`
  }
  if (Array.isArray(value)) {
    return value.map(v => String(v)).join('、')
  }
  return String(value)
}

function formatTime(t) {
  if (!t) return ''
  return t.replace('T', ' ').substring(0, 19)
}

function formatDay(date) {
  if (!date) return ''
  return date.substring(5)
}

const maxCount = computed(() => {
  if (!stats.value?.last7Days) return 0
  return Math.max(...stats.value.last7Days.map(d => d.count), 1)
})

function getBarHeight(count) {
  if (!count) return 2
  return Math.max(4, (count / maxCount.value) * 100)
}

onMounted(loadAll)
</script>

<style scoped>
.data-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.data-header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}

.header-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  min-width: 0;
}

.header-info {
  min-width: 0;
}

.template-name {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 4px;
  color: #303133;
}

.template-desc {
  margin: 0;
  font-size: 13px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 400px;
}

.header-right {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.stats-section {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px 24px 0;
}

.stat-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
}

.stat-card.total .stat-icon {
  background: #ecf5ff;
  color: #409EFF;
}

.stat-card.today .stat-icon {
  background: #fdf6ec;
  color: #E6A23C;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.chart-card {
  flex-direction: column;
  align-items: stretch;
  gap: 12px;
}

.chart-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.chart-bars {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 8px;
  height: 80px;
}

.chart-bar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  min-width: 0;
}

.bar-wrap {
  width: 100%;
  height: 50px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar {
  width: 60%;
  background: linear-gradient(180deg, #409EFF 0%, #66b1ff 100%);
  border-radius: 3px 3px 0 0;
  min-height: 4px;
  transition: height 0.3s;
}

.bar-label {
  font-size: 11px;
  color: #909399;
}

.bar-count {
  font-size: 10px;
  color: #c0c4cc;
}

.data-main {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px 24px 32px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  background: #fff;
  padding: 12px 16px;
  border-radius: 8px;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 960px) {
  .stat-cards {
    grid-template-columns: 1fr;
  }
}
</style>
