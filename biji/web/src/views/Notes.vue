<template>
  <div class="notes-page">
    <div class="page-header flex items-center justify-between mb-md">
      <h2 class="page-title">
        {{ pageTitle }}
        <span class="text-muted text-small">({{ notes.length }})</span>
      </h2>
      <div v-if="status === 'deleted'" class="flex gap-sm">
        <el-button type="danger" @click="handleEmptyTrash">清空回收站</el-button>
      </div>
    </div>

    <div v-loading="loading" class="notes-grid">
      <el-empty v-if="!loading && notes.length === 0" :description="emptyText" />
      <div
        v-for="note in notes"
        :key="note.id"
        class="note-card"
        :class="{ pinned: note.isPinned }"
        @click="handleNoteClick(note)"
      >
        <div class="note-header flex items-center justify-between">
          <h3 class="note-title ellipsis">{{ note.title || '无标题笔记' }}</h3>
          <div class="note-actions flex gap-sm">
            <el-tooltip content="置顶">
              <el-icon
                class="action-icon"
                :class="{ active: note.isPinned }"
                @click.stop="togglePin(note)"
              >
                <component :is="note.isPinned ? 'Top' : 'Bottom'" />
              </el-icon>
            </el-tooltip>
            <el-tooltip v-if="status !== 'deleted'" content="归档">
              <el-icon
                class="action-icon"
                :class="{ active: note.isArchived }"
                @click.stop="toggleArchive(note)"
              >
                <Folder />
              </el-icon>
            </el-tooltip>
            <el-tooltip content="加密">
              <el-icon
                class="action-icon"
                :class="{ active: note.isEncrypted }"
                @click.stop="handleEncryptClick(note)"
              >
                <Lock />
              </el-icon>
            </el-tooltip>
            <el-dropdown trigger="click" @command="(cmd) => handleMoreAction(cmd, note)" @click.stop>
              <el-icon class="action-icon"><MoreFilled /></el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="status === 'deleted'" command="restore">恢复</el-dropdown-item>
                  <el-dropdown-item command="edit">编辑</el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    {{ status === 'deleted' ? '永久删除' : '删除' }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>

        <div class="note-preview ellipsis-2">
          <template v-if="note.isEncrypted">
            <span class="text-muted"><el-icon><Lock /></el-icon> 已加密，点击查看</span>
          </template>
          <template v-else>
            {{ note.plainText || '暂无内容' }}
          </template>
        </div>

        <div class="note-footer flex items-center justify-between">
          <div class="note-tags flex gap-sm" style="flex-wrap: wrap">
            <el-tag
              v-for="tag in note.tags?.slice(0, 3)"
              :key="tag.id"
              :color="tag.color"
              effect="plain"
              size="small"
            >
              {{ tag.name }}
            </el-tag>
            <span v-if="note.tags?.length > 3" class="text-muted text-small">
              +{{ note.tags.length - 3 }}
            </span>
          </div>
          <span class="text-muted text-small">{{ formatDate(note.updatedAt) }}</span>
        </div>

        <div v-if="note.isPinned" class="pin-indicator">
          <el-icon :size="14"><Top /></el-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { noteApi } from '@/api'

const props = defineProps({
  status: {
    type: String,
    default: 'normal'
  },
  searchKeyword: String,
  selectedTagIds: {
    type: Array,
    default: () => []
  },
  dateRange: {
    type: Array,
    default: null
  }
})

const emit = defineEmits(['filters-changed'])

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const notes = ref([])

const pageTitle = computed(() => {
  switch (props.status) {
    case 'archived': return '归档笔记'
    case 'deleted': return '回收站'
    default: return '全部笔记'
  }
})

const emptyText = computed(() => {
  switch (props.status) {
    case 'archived': return '暂无归档笔记'
    case 'deleted': return '回收站为空'
    default: return '暂无笔记，点击右上角新建笔记'
  }
})

const formatDate = (date) => {
  if (!date) return ''
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const loadNotes = async () => {
  loading.value = true
  try {
    const params = {
      status: props.status
    }
    if (props.searchKeyword) {
      params.keyword = props.searchKeyword
    }
    if (props.selectedTagIds && props.selectedTagIds.length) {
      params.tagIds = props.selectedTagIds
    }
    if (props.dateRange && props.dateRange.length === 2) {
      params.startDate = dayjs(props.dateRange[0]).format('YYYY-MM-DD')
      params.endDate = dayjs(props.dateRange[1]).format('YYYY-MM-DD')
    }
    notes.value = await noteApi.list(params)
  } finally {
    loading.value = false
  }
}

const handleNoteClick = (note) => {
  if (note.isEncrypted) {
    showPasswordDialog(note)
  } else {
    router.push(`/notes/${note.id}`)
  }
}

const showPasswordDialog = async (note) => {
  try {
    const { value: password } = await ElMessageBox.prompt('请输入密码', '查看加密笔记', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'password',
      inputValidator: (v) => !!v || '请输入密码'
    })
    const result = await noteApi.verify(note.id, password)
    router.push({
      path: `/notes/${note.id}`,
      query: { decrypted: '1' }
    })
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('密码错误')
    }
  }
}

const handleEncryptClick = async (note) => {
  if (note.isEncrypted) {
    try {
      await ElMessageBox.confirm('确定要取消加密吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      await noteApi.update(note.id, { isEncrypted: false, password: '' })
      ElMessage.success('已取消加密')
      loadNotes()
    } catch (e) {}
  } else {
    try {
      const { value: password } = await ElMessageBox.prompt('设置笔记密码', '加密笔记', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password',
        inputValidator: (v) => (v && v.length >= 4) || '密码至少4位'
      })
      await noteApi.update(note.id, { isEncrypted: true, password })
      ElMessage.success('加密成功')
      loadNotes()
    } catch (e) {}
  }
}

const togglePin = async (note) => {
  await noteApi.togglePin(note.id)
  ElMessage.success(note.isPinned ? '已取消置顶' : '已置顶')
  loadNotes()
}

const toggleArchive = async (note) => {
  await noteApi.toggleArchive(note.id)
  ElMessage.success(note.isArchived ? '已取消归档' : '已归档')
  loadNotes()
}

const handleMoreAction = async (command, note) => {
  switch (command) {
    case 'edit':
      router.push(`/editor/${note.id}`)
      break
    case 'restore':
      await noteApi.restore(note.id)
      ElMessage.success('已恢复')
      loadNotes()
      break
    case 'delete':
      try {
        const permanent = props.status === 'deleted'
        await ElMessageBox.confirm(
          permanent ? '确定永久删除？此操作不可恢复！' : '确定删除到回收站？',
          '提示',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        await noteApi.remove(note.id, permanent)
        ElMessage.success('已删除')
        loadNotes()
      } catch (e) {}
      break
  }
}

const handleEmptyTrash = async () => {
  try {
    await ElMessageBox.confirm('确定清空回收站？所有笔记将永久删除！', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await noteApi.emptyTrash()
    ElMessage.success('回收站已清空')
    loadNotes()
  } catch (e) {}
}

watch([
  () => props.status,
  () => props.searchKeyword,
  () => props.selectedTagIds,
  () => props.dateRange
], loadNotes, { deep: true })

onMounted(loadNotes)
</script>

<style lang="scss" scoped>
.notes-page {
  height: 100%;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.notes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.note-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid #ebeef5;
  position: relative;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
    border-color: #409EFF;
  }

  &.pinned {
    border-color: #E6A23C;
    background: linear-gradient(135deg, #fff 0%, #fdf6ec 100%);
  }
}

.note-header {
  margin-bottom: 8px;
}

.note-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  flex: 1;
  margin-right: 8px;
}

.note-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.note-card:hover .note-actions {
  opacity: 1;
}

.action-icon {
  font-size: 16px;
  color: #909399;
  padding: 4px;
  border-radius: 4px;

  &:hover {
    background: #f0f0f0;
    color: #409EFF;
  }

  &.active {
    color: #409EFF;
  }
}

.note-preview {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 12px;
  min-height: 44px;
}

.note-footer {
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
}

.note-tags {
  flex: 1;
  margin-right: 8px;
  min-height: 20px;
}

.pin-indicator {
  position: absolute;
  top: 8px;
  right: 8px;
  color: #E6A23C;
}
</style>
