<template>
  <div class="template-page">
    <header class="header">
      <div class="header-content">
        <div class="back-btn" @click="goBack">
          <el-icon :size="20"><ArrowLeft /></el-icon>
          <span>返回</span>
        </div>
        <h1 class="page-title">选择简历模板</h1>
        <div style="width: 80px;"></div>
      </div>
    </header>

    <main class="main-content">
      <div class="container">
        <p class="subtitle">选择一个模板，开始创建你的简历</p>

        <div v-if="loading" class="loading-wrapper">
          <el-icon :size="40" class="loading-icon"><Loading /></el-icon>
          <p>加载模板中...</p>
        </div>

        <div v-else class="template-grid">
          <div
            v-for="template in templates"
            :key="template.id"
            class="template-card"
            @click="handleSelectTemplate(template)"
          >
            <div class="template-preview">
              <div
                class="preview-content"
                :style="getPreviewStyle(template)"
              >
                <div class="preview-header" :style="{ backgroundColor: getHeaderColor(template) }">
                  <div class="preview-avatar" :style="{ backgroundColor: getAvatarColor(template) }"></div>
                  <div class="preview-title-bar">
                    <div class="preview-name" :style="{ backgroundColor: getNameColor(template) }"></div>
                    <div class="preview-subtitle" :style="{ backgroundColor: getSubtitleColor(template) }"></div>
                  </div>
                </div>
                <div class="preview-body">
                  <div class="preview-section">
                    <div class="preview-section-title" :style="{ backgroundColor: getSectionColor(template) }"></div>
                    <div class="preview-line" v-for="i in 3" :key="i"></div>
                  </div>
                  <div class="preview-section">
                    <div class="preview-section-title" :style="{ backgroundColor: getSectionColor(template) }"></div>
                    <div class="preview-line" v-for="i in 2" :key="i"></div>
                  </div>
                </div>
              </div>
              <div class="template-overlay">
                <el-button type="primary" size="large">
                  <el-icon><Check /></el-icon>
                  选择此模板
                </el-button>
              </div>
            </div>

            <div class="template-info">
              <h3 class="template-name">{{ template.name }}</h3>
              <p class="template-desc">{{ template.description }}</p>
              <el-tag size="small" type="info">{{ template.category || '通用' }}</el-tag>
            </div>
          </div>
        </div>
      </div>
    </main>

    <el-dialog
      v-model="dialogVisible"
      title="输入简历标题"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" label-position="top">
        <el-form-item label="简历标题">
          <el-input
            v-model="form.title"
            placeholder="请输入简历标题，如：前端开发工程师简历"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreateResume">
          创建简历
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Loading,
  Check
} from '@element-plus/icons-vue'
import { useResumeStore } from '@/stores/resume'
import { getTemplates, createResume } from '@/api/resume'

const router = useRouter()
const resumeStore = useResumeStore()

const loading = ref(false)
const templates = ref([])
const dialogVisible = ref(false)
const creating = ref(false)
const selectedTemplate = ref(null)

const form = reactive({
  title: ''
})

const colorSchemes = {
  1: {
    header: '#409eff',
    avatar: '#ffffff',
    name: '#ffffff',
    subtitle: 'rgba(255,255,255,0.8)',
    section: '#409eff'
  },
  2: {
    header: '#67c23a',
    avatar: '#ffffff',
    name: '#ffffff',
    subtitle: 'rgba(255,255,255,0.8)',
    section: '#67c23a'
  },
  3: {
    header: '#e6a23c',
    avatar: '#ffffff',
    name: '#ffffff',
    subtitle: 'rgba(255,255,255,0.8)',
    section: '#e6a23c'
  }
}

const fetchTemplates = async () => {
  loading.value = true
  try {
    const data = await getTemplates()
    templates.value = data.data || data || []
    resumeStore.setTemplates(templates.value)
  } catch (error) {
    ElMessage.error('获取模板列表失败')
  } finally {
    loading.value = false
  }
}

const getPreviewStyle = () => ({})

const getHeaderColor = (template) => {
  const scheme = colorSchemes[template.id] || colorSchemes[1]
  return scheme.header
}

const getAvatarColor = (template) => {
  const scheme = colorSchemes[template.id] || colorSchemes[1]
  return scheme.avatar
}

const getNameColor = (template) => {
  const scheme = colorSchemes[template.id] || colorSchemes[1]
  return scheme.name
}

const getSubtitleColor = (template) => {
  const scheme = colorSchemes[template.id] || colorSchemes[1]
  return scheme.subtitle
}

const getSectionColor = (template) => {
  const scheme = colorSchemes[template.id] || colorSchemes[1]
  return scheme.section
}

const goBack = () => {
  router.push('/')
}

const handleSelectTemplate = (template) => {
  selectedTemplate.value = template
  form.title = ''
  dialogVisible.value = true
}

const handleCreateResume = async () => {
  if (!form.title.trim()) {
    ElMessage.warning('请输入简历标题')
    return
  }

  creating.value = true
  try {
    const data = await createResume({
      user_identity: resumeStore.userIdentity,
      template_id: selectedTemplate.value.id,
      title: form.title.trim()
    })

    ElMessage.success('简历创建成功')
    dialogVisible.value = false

    const resumeId = data?.id
    if (resumeId) {
      router.push(`/editor/${resumeId}`)
    } else {
      router.push('/')
    }
  } catch (error) {
    ElMessage.error('创建简历失败')
  } finally {
    creating.value = false
  }
}

onMounted(() => {
  fetchTemplates()
})
</script>

<style lang="scss" scoped>
.template-page {
  min-height: 100%;
  background-color: #f5f7fa;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 16px 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .back-btn {
      display: flex;
      align-items: center;
      gap: 4px;
      color: #606266;
      cursor: pointer;
      transition: color 0.2s;

      &:hover {
        color: #409eff;
      }
    }

    .page-title {
      font-size: 20px;
      font-weight: 600;
      color: #303133;
      margin: 0;
    }
  }
}

.main-content {
  padding: 32px 24px;

  .container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .subtitle {
    text-align: center;
    color: #909399;
    margin-bottom: 32px;
    font-size: 14px;
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

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.template-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-4px);

    .template-preview {
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);

      .template-overlay {
        opacity: 1;
      }
    }
  }

  .template-preview {
    position: relative;
    background-color: #fff;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    aspect-ratio: 210 / 297;

    .preview-content {
      width: 100%;
      height: 100%;
      display: flex;
      flex-direction: column;

      .preview-header {
        padding: 16px;
        display: flex;
        align-items: center;
        gap: 12px;

        .preview-avatar {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          flex-shrink: 0;
        }

        .preview-title-bar {
          flex: 1;
          display: flex;
          flex-direction: column;
          gap: 6px;

          .preview-name {
            height: 14px;
            width: 60%;
            border-radius: 2px;
          }

          .preview-subtitle {
            height: 10px;
            width: 40%;
            border-radius: 2px;
          }
        }
      }

      .preview-body {
        flex: 1;
        padding: 16px;
        background-color: #fff;

        .preview-section {
          margin-bottom: 16px;

          .preview-section-title {
            height: 12px;
            width: 30%;
            border-radius: 2px;
            margin-bottom: 8px;
          }

          .preview-line {
            height: 8px;
            background-color: #f0f0f0;
            border-radius: 2px;
            margin-bottom: 6px;

            &:nth-child(2) {
              width: 90%;
            }
            &:nth-child(3) {
              width: 80%;
            }
            &:nth-child(4) {
              width: 85%;
            }
          }
        }
      }
    }

    .template-overlay {
      position: absolute;
      inset: 0;
      background-color: rgba(0, 0, 0, 0.5);
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: opacity 0.2s;
    }
  }

  .template-info {
    padding: 16px 8px;

    .template-name {
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 8px 0;
    }

    .template-desc {
      font-size: 13px;
      color: #606266;
      margin: 0 0 12px 0;
      line-height: 1.5;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
  }
}
</style>
