<template>
  <div class="home-page">
    <header class="header">
      <div class="header-content">
        <div class="logo">
          <el-icon :size="28" color="#409eff"><Document /></el-icon>
          <span class="logo-text">简历生成器</span>
        </div>
        <el-button type="primary" @click="goToTemplates">
          <el-icon><Plus /></el-icon>
          新建简历
        </el-button>
      </div>
    </header>

    <main class="main-content">
      <div class="container">
        <h2 class="page-title">我的简历</h2>

        <div v-if="loading" class="loading-wrapper">
          <el-icon :size="40" class="loading-icon"><Loading /></el-icon>
          <p>加载中...</p>
        </div>

        <div v-else-if="resumeList.length === 0" class="empty-wrapper">
          <el-empty description="还没有简历，创建你的第一份简历吧">
            <el-button type="primary" @click="goToTemplates">
              <el-icon><Plus /></el-icon>
              创建第一份简历
            </el-button>
          </el-empty>
        </div>

        <div v-else class="resume-grid">
          <el-card
            v-for="resume in resumeList"
            :key="resume.id"
            class="resume-card"
            shadow="hover"
          >
            <template #header>
              <div class="card-header">
                <span class="resume-title">{{ resume.title }}</span>
                <el-tag size="small" :type="getTemplateTagType(resume.template_id)">
                  {{ resume.template_name || '模板' + resume.template_id }}
                </el-tag>
              </div>
            </template>

            <div class="card-body">
              <div class="resume-info">
                <div class="info-item">
                  <el-icon :size="14"><Clock /></el-icon>
                  <span class="info-label">更新时间：</span>
                  <span class="info-value">{{ formatDate(resume.updated_at) }}</span>
                </div>
                <div class="info-item">
                  <el-icon :size="14"><Calendar /></el-icon>
                  <span class="info-label">创建时间：</span>
                  <span class="info-value">{{ formatDate(resume.created_at) }}</span>
                </div>
              </div>

              <div class="card-actions">
                <el-button type="primary" size="small" @click="editResume(resume.id)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button type="success" size="small" @click="handleExportPDF(resume)">
                  <el-icon><Download /></el-icon>
                  导出PDF
                </el-button>
                <el-button type="danger" size="small" @click="handleDelete(resume)">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  Plus,
  Clock,
  Calendar,
  Edit,
  Download,
  Delete,
  Loading
} from '@element-plus/icons-vue'
import { useResumeStore } from '@/stores/resume'
import { getResumeList, deleteResume, exportPDF } from '@/api/resume'

const router = useRouter()
const resumeStore = useResumeStore()

const loading = ref(false)
const resumeList = ref([])

const fetchResumeList = async () => {
  loading.value = true
  try {
    const data = await getResumeList(resumeStore.userIdentity)
    resumeList.value = data || []
    resumeStore.setResumeList(resumeList.value)
  } catch (error) {
    ElMessage.error('获取简历列表失败')
  } finally {
    loading.value = false
  }
}

const goToTemplates = () => {
  router.push('/templates')
}

const editResume = (id) => {
  router.push(`/editor/${id}`)
}

const handleDelete = async (resume) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除简历「${resume.title}」吗？删除后无法恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteResume(resume.id)
    ElMessage.success('删除成功')
    fetchResumeList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleExportPDF = async (resume) => {
  try {
    const blob = await exportPDF(resume.id)
    const url = window.URL.createObjectURL(new Blob([blob]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `${resume.title}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出PDF失败')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getTemplateTagType = (templateId) => {
  const types = ['primary', 'success', 'warning', 'info']
  return types[(templateId - 1) % types.length] || 'info'
}

onMounted(() => {
  fetchResumeList()
})
</script>

<style lang="scss" scoped>
.home-page {
  min-height: 100%;
  background-color: #f5f7fa;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 16px 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 8px;

    .logo-text {
      font-size: 20px;
      font-weight: 600;
      color: #303133;
    }
  }
}

.main-content {
  padding: 32px 24px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 24px;
  }
}

.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  color: #909399;

  .loading-icon {
    animation: rotate 1s linear infinite;
    margin-bottom: 16px;
  }

  @keyframes rotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
}

.empty-wrapper {
  padding: 80px 0;
}

.resume-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.resume-card {
  transition: transform 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-2px);
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .resume-title {
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      max-width: 200px;
    }
  }

  .card-body {
    .resume-info {
      margin-bottom: 16px;

      .info-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #606266;
        margin-bottom: 8px;

        &:last-child {
          margin-bottom: 0;
        }

        .info-label {
          color: #909399;
        }

        .info-value {
          color: #606266;
        }
      }
    }

    .card-actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
  }
}
</style>
