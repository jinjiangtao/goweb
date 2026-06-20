<template>
  <div class="editor-page">
    <header class="editor-header">
      <div class="header-left">
        <el-button type="primary" text @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </div>
      <div class="header-center">
        <el-input
          v-model="resumeTitle"
          class="title-input"
          placeholder="简历标题"
          @blur="handleTitleChange"
        />
      </div>
      <div class="header-right">
        <el-button @click="showVersionHistory">
          <el-icon><Clock /></el-icon>
          版本历史
        </el-button>
        <el-button type="primary" :loading="isSaving" @click="handleSave">
          <el-icon><Check /></el-icon>
          {{ isSaving ? '保存中...' : '保存' }}
        </el-button>
        <el-button type="success" @click="handleExportPDF">
          <el-icon><Download /></el-icon>
          导出PDF
        </el-button>
      </div>
    </header>

    <div class="editor-body">
      <aside class="sidebar-left">
        <ModuleList />
        <StylePanel />
      </aside>

      <main class="content-area">
        <ContentEditor />
      </main>

      <aside class="sidebar-right">
        <ResumePreview />
      </aside>
    </div>

    <VersionPanel
      v-model="versionPanelVisible"
      :resume-id="store.resumeId"
      @restored="handleVersionRestored"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useResumeStore } from '@/stores/resume'
import { exportPDF } from '@/api/resume'
import { useDebounceFn } from '@vueuse/core'
import ModuleList from '@/components/ModuleList.vue'
import ContentEditor from '@/components/ContentEditor.vue'
import StylePanel from '@/components/StylePanel.vue'
import ResumePreview from '@/components/ResumePreview.vue'
import VersionPanel from '@/components/VersionPanel.vue'
import { ArrowLeft, Clock, Check, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const store = useResumeStore()

const resumeTitle = ref('未命名简历')
const versionPanelVisible = ref(false)

const isSaving = computed(() => store.isSaving)

watch(
  () => store.currentResume,
  (resume) => {
    if (resume?.title) {
      resumeTitle.value = resume.title
    }
  },
  { immediate: true, deep: true }
)

const debouncedSave = useDebounceFn(async () => {
  if (store.resumeId && store.resumeId !== 'new') {
    await store.saveResume()
  }
}, 3000)

const stopWatching = watch(
  () => [store.content, store.styleConfig, store.modules],
  () => {
    if (store.resumeId && store.resumeId !== 'new') {
      debouncedSave()
    }
  },
  { deep: true }
)

onMounted(async () => {
  const id = route.params.id
  await store.loadResume(id)
})

onBeforeUnmount(() => {
  stopWatching && stopWatching()
  debouncedSave.flush()
})

function goBack() {
  router.push('/')
}

function handleTitleChange() {
  if (store.currentResume) {
    store.currentResume.title = resumeTitle.value
  }
  if (store.resumeId && store.resumeId !== 'new') {
    store.saveResume()
  }
}

async function handleSave() {
  if (store.resumeId === 'new') {
    const newId = await store.createNewResume(resumeTitle.value)
    if (newId) {
      router.replace(`/editor/${newId}`)
      ElMessage.success('创建成功')
    }
  } else {
    await store.saveResume()
    ElMessage.success('保存成功')
  }
}

async function handleExportPDF() {
  if (!store.resumeId || store.resumeId === 'new') {
    ElMessage.warning('请先保存简历')
    return
  }
  try {
    const modulesWithOrder = store.modules.map((m, idx) => ({
      ...m,
      order: idx
    }))
    const fullContent = {
      ...store.content,
      modules: modulesWithOrder
    }
    const res = await exportPDF(fullContent, store.styleConfig, `${resumeTitle.value || '简历'}.pdf`)
    const blob = new Blob([res], { type: 'application/pdf' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${resumeTitle.value || '简历'}.pdf`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (e) {
    ElMessage.error('导出失败')
    console.error(e)
  }
}

function showVersionHistory() {
  if (!store.resumeId || store.resumeId === 'new') {
    ElMessage.warning('请先保存简历')
    return
  }
  versionPanelVisible.value = true
}

async function handleVersionRestored() {
  if (store.resumeId && store.resumeId !== 'new') {
    await store.loadResume(store.resumeId)
  }
}
</script>

<style lang="scss" scoped>
.editor-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;

  .editor-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    background: #fff;
    border-bottom: 1px solid #e4e7ed;
    height: 60px;
    flex-shrink: 0;

    .header-left {
      min-width: 120px;
    }

    .header-center {
      flex: 1;
      display: flex;
      justify-content: center;

      .title-input {
        width: 400px;
        text-align: center;
        font-size: 16px;
        font-weight: 600;

        :deep(.el-input__inner) {
          text-align: center;
        }
      }
    }

    .header-right {
      display: flex;
      gap: 12px;
      min-width: 350px;
      justify-content: flex-end;
    }
  }

  .editor-body {
    display: flex;
    flex: 1;
    overflow: hidden;

    .sidebar-left {
      width: 260px;
      display: flex;
      flex-direction: column;
      flex-shrink: 0;
      overflow: hidden;
    }

    .content-area {
      flex: 1;
      overflow: hidden;
      min-width: 0;
    }

    .sidebar-right {
      width: 420px;
      flex-shrink: 0;
      overflow: hidden;
    }
  }

}
</style>
