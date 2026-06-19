<template>
  <div class="note-detail" v-loading="loading">
    <template v-if="note">
      <div class="detail-header flex items-center justify-between mb-md">
        <div class="flex items-center gap-sm">
          <el-button link @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h1 class="detail-title">{{ note.title || '无标题笔记' }}</h1>
        </div>
        <div class="flex gap-sm">
          <el-button @click="editNote">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-button :type="note.isPinned ? 'warning' : ''" @click="togglePin">
            <el-icon><component :is="note.isPinned ? 'Top' : 'Bottom'" /></el-icon>
            {{ note.isPinned ? '取消置顶' : '置顶' }}
          </el-button>
          <el-button @click="toggleArchive">
            <el-icon><Folder /></el-icon>
            {{ note.isArchived ? '取消归档' : '归档' }}
          </el-button>
          <el-button type="danger" @click="deleteNote">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>

      <div class="detail-meta flex items-center gap-md mb-md">
        <el-tag v-if="note.isEncrypted" type="warning" effect="dark">
          <el-icon><Lock /></el-icon>
          已加密
        </el-tag>
        <el-tag v-if="note.isPinned" type="warning">
          <el-icon><Top /></el-icon>
          置顶
        </el-tag>
        <el-tag v-if="note.isArchived" type="info">已归档</el-tag>
        <span class="text-muted text-small">创建于: {{ formatDate(note.createdAt) }}</span>
        <span class="text-muted text-small">更新于: {{ formatDate(note.updatedAt) }}</span>
      </div>

      <div class="detail-tags flex gap-sm mb-md" style="flex-wrap: wrap">
        <el-tag
          v-for="tag in note.tags"
          :key="tag.id"
          :color="tag.color"
          effect="light"
        >
          {{ tag.name }}
        </el-tag>
        <span v-if="!note.tags?.length" class="text-muted text-small">无标签</span>
      </div>

      <div class="detail-content">
        <div v-if="displayContent" class="content-html" v-html="displayContent" />
        <div v-else class="text-muted">暂无内容</div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { noteApi } from '@/api'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const note = ref(null)
const displayContent = ref('')

const formatDate = (date) => {
  if (!date) return ''
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const loadNote = async () => {
  loading.value = true
  try {
    note.value = await noteApi.get(route.params.id)
    if (note.value.isEncrypted) {
      if (route.query.decrypted === '1') {
        showPasswordDialog()
      } else {
        displayContent.value = '<p class="text-muted"><el-icon><Lock /></el-icon> 此笔记已加密，请先输入密码查看</p>'
      }
    } else {
      displayContent.value = note.value.content
    }
  } finally {
    loading.value = false
  }
}

const showPasswordDialog = async () => {
  try {
    const { value: password } = await ElMessageBox.prompt('请输入密码', '查看加密笔记', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'password',
      inputValidator: (v) => !!v || '请输入密码'
    })
    const result = await noteApi.verify(note.value.id, password)
    displayContent.value = result.content
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('密码错误')
    }
  }
}

const goBack = () => {
  router.back()
}

const editNote = () => {
  router.push(`/editor/${note.value.id}`)
}

const togglePin = async () => {
  await noteApi.togglePin(note.value.id)
  note.value.isPinned = !note.value.isPinned
  ElMessage.success(note.value.isPinned ? '已置顶' : '已取消置顶')
}

const toggleArchive = async () => {
  await noteApi.toggleArchive(note.value.id)
  note.value.isArchived = !note.value.isArchived
  ElMessage.success(note.value.isArchived ? '已归档' : '已取消归档')
}

const deleteNote = async () => {
  try {
    await ElMessageBox.confirm('确定删除到回收站？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await noteApi.remove(note.value.id)
    ElMessage.success('已删除')
    router.back()
  } catch (e) {}
}

watch(() => route.params.id, loadNote)

onMounted(loadNote)
</script>

<style lang="scss" scoped>
.note-detail {
  background: #fff;
  border-radius: 8px;
  padding: 32px;
  min-height: 100%;
}

.detail-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.detail-meta {
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.detail-tags {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-content {
  padding-top: 24px;
}

.content-html {
  :deep(img) {
    max-width: 100%;
    height: auto;
    border-radius: 4px;
  }

  :deep(h1), :deep(h2), :deep(h3) {
    margin-top: 24px;
    margin-bottom: 16px;
  }

  :deep(p) {
    line-height: 1.8;
    margin-bottom: 16px;
  }

  :deep(ul), :deep(ol) {
    padding-left: 24px;
    margin-bottom: 16px;
  }

  :deep(blockquote) {
    border-left: 4px solid #409EFF;
    padding-left: 16px;
    color: #606266;
    margin: 16px 0;
  }

  :deep(code) {
    background: #f5f7fa;
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Consolas', monospace;
  }

  :deep(pre) {
    background: #282c34;
    color: #abb2bf;
    padding: 16px;
    border-radius: 8px;
    overflow-x: auto;

    code {
      background: transparent;
      color: inherit;
      padding: 0;
    }
  }

  :deep(table) {
    border-collapse: collapse;
    width: 100%;
    margin-bottom: 16px;

    th, td {
      border: 1px solid #dcdfe6;
      padding: 8px 12px;
    }

    th {
      background: #f5f7fa;
    }
  }
}
</style>
