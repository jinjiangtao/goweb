<template>
  <div class="editor-page">
    <div class="editor-header flex items-center justify-between">
      <div class="flex items-center gap-sm">
        <el-button link @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <el-input
          v-model="title"
          placeholder="请输入标题..."
          class="title-input"
          size="large"
        />
      </div>
      <div class="flex items-center gap-sm">
        <span v-if="saveStatus" class="save-status" :class="saveStatusClass">
          {{ saveStatusText }}
        </span>
        <el-switch
          v-model="isEncrypted"
          active-text="加密"
          inactive-text="不加密"
          style="margin-right: 16px"
        />
        <el-button @click="saveNote" :loading="saving">
          <el-icon><Check /></el-icon>
          保存
        </el-button>
        <el-button type="primary" @click="publishNote" :loading="saving">
          <el-icon><Promotion /></el-icon>
          保存并返回
        </el-button>
      </div>
    </div>

    <div class="editor-toolbar">
      <div class="flex items-center gap-md">
        <span class="text-muted">标签:</span>
        <el-select
          v-model="selectedTagIds"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="选择或创建标签"
          style="width: 400px"
          @change="markUnsaved"
        >
          <el-option
            v-for="tag in flatTags"
            :key="tag.id"
            :label="tag.name"
            :value="tag.id"
          >
            <div class="flex items-center gap-sm">
              <span
                class="tag-color-dot"
                :style="{ backgroundColor: tag.color }"
              />
              <span>{{ tag.name }}</span>
            </div>
          </el-option>
        </el-select>
        <el-color-picker v-model="newTagColor" />
        <el-button size="small" @click="handleCreateTag">
          <el-icon><Plus /></el-icon>
          新建标签
        </el-button>
      </div>
    </div>

    <div v-if="isEncrypted && !noteId" class="password-section p-md">
      <el-form inline>
        <el-form-item label="加密密码" required>
          <el-input
            v-model="password"
            type="password"
            placeholder="请输入密码（至少4位）"
            show-password
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item label="确认密码" required>
          <el-input
            v-model="confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            show-password
            style="width: 240px"
          />
        </el-form-item>
      </el-form>
    </div>

    <div class="editor-container">
      <div style="border: 1px solid #ccc; z-index: 100">
        <Toolbar
          style="border-bottom: 1px solid #ccc"
          :editor="editorRef"
          :default-config="toolbarConfig"
          mode="default"
        />
        <Editor
          v-model="content"
          style="height: 500px; overflow-y: hidden"
          :default-config="editorConfig"
          mode="default"
          @onCreated="handleEditorCreated"
          @onChange="handleEditorChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, shallowRef, onBeforeUnmount, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import '@wangeditor/editor/dist/css/style.css'
import { noteApi, tagApi } from '@/api'

const router = useRouter()
const route = useRoute()

const noteId = computed(() => route.params.id)

const title = ref('')
const content = ref('')
const selectedTagIds = ref([])
const isEncrypted = ref(false)
const password = ref('')
const confirmPassword = ref('')
const newTagColor = ref('#409EFF')
const allTags = ref([])
const saving = ref(false)
const saveStatus = ref('')
const hasUnsavedChanges = ref(false)

const editorRef = shallowRef()

const saveStatusClass = computed(() => {
  if (saveStatus.value === 'saving') return 'status-saving'
  if (saveStatus.value === 'saved') return 'status-saved'
  if (saveStatus.value === 'unsaved') return 'status-unsaved'
  return ''
})

const saveStatusText = computed(() => {
  if (saveStatus.value === 'saving') return '保存中...'
  if (saveStatus.value === 'saved') return '已自动保存'
  if (saveStatus.value === 'unsaved') return '有未保存的更改'
  return ''
})

const flatTags = computed(() => {
  const result = []
  const flatten = (tags) => {
    tags.forEach(tag => {
      result.push(tag)
      if (tag.children && tag.children.length) {
        flatten(tag.children)
      }
    })
  }
  flatten(allTags.value)
  return result
})

const toolbarConfig = {
  excludeKeys: ['group-video']
}

const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {}
}

let autoSaveTimer = null
let saveStatusTimer = null

const handleEditorCreated = (editor) => {
  editorRef.value = editor
}

const handleEditorChange = () => {
  markUnsaved()
}

const markUnsaved = () => {
  hasUnsavedChanges.value = true
  saveStatus.value = 'unsaved'
  scheduleAutoSave()
}

const scheduleAutoSave = () => {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => {
    if (noteId.value) {
      doSave(false)
    }
  }, 2000)
}

const loadNote = async () => {
  if (!noteId.value) return

  const note = await noteApi.get(noteId.value)
  title.value = note.title
  selectedTagIds.value = (note.tags || []).map(t => t.id)
  isEncrypted.value = note.isEncrypted

  if (note.isEncrypted) {
    try {
      const { value: pwd } = await ElMessageBox.prompt('请输入密码解密笔记', '编辑加密笔记', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password',
        inputValidator: (v) => !!v || '请输入密码'
      })
      const result = await noteApi.verify(noteId.value, pwd)
      content.value = result.content
    } catch (e) {
      if (e === 'cancel') {
        router.back()
      }
    }
  } else {
    content.value = note.content || ''
  }
}

const loadTags = async () => {
  allTags.value = await tagApi.list()
}

const handleCreateTag = async () => {
  const { value: name } = await ElMessageBox.prompt('输入标签名称', '新建标签', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputValidator: (v) => !!v?.trim() || '请输入标签名称'
  })

  const tag = await tagApi.create({
    name: name.trim(),
    color: newTagColor.value
  })
  allTags.value.push(tag)
  selectedTagIds.value.push(tag.id)
  ElMessage.success('标签创建成功')
}

const validateForm = () => {
  if (!title.value.trim() && !content.value.trim()) {
    ElMessage.warning('标题和内容不能都为空')
    return false
  }
  if (isEncrypted.value && !noteId.value) {
    if (password.value.length < 4) {
      ElMessage.warning('密码至少4位')
      return false
    }
    if (password.value !== confirmPassword.value) {
      ElMessage.warning('两次输入的密码不一致')
      return false
    }
  }
  return true
}

const doSave = async (showMsg = true) => {
  if (!validateForm()) return false

  saving.value = true
  saveStatus.value = 'saving'

  try {
    const data = {
      title: title.value.trim(),
      content: content.value,
      tagIds: selectedTagIds.value,
      isEncrypted: isEncrypted.value,
      password: isEncrypted.value ? (password.value || undefined) : undefined
    }

    let savedNote
    if (noteId.value) {
      savedNote = await noteApi.update(noteId.value, data)
    } else {
      savedNote = await noteApi.create(data)
    }

    hasUnsavedChanges.value = false
    saveStatus.value = 'saved'

    if (saveStatusTimer) clearTimeout(saveStatusTimer)
    saveStatusTimer = setTimeout(() => {
      saveStatus.value = ''
    }, 3000)

    if (showMsg) {
      ElMessage.success('保存成功')
    }

    return savedNote
  } catch (e) {
    saveStatus.value = 'unsaved'
    return null
  } finally {
    saving.value = false
  }
}

const saveNote = () => {
  doSave(true)
}

const publishNote = async () => {
  const saved = await doSave(true)
  if (saved) {
    router.push(`/notes/${saved.id}`)
  }
}

const goBack = () => {
  if (hasUnsavedChanges.value) {
    ElMessageBox.confirm('有未保存的更改，确定要离开吗？', '提示', {
      confirmButtonText: '离开',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      router.back()
    }).catch(() => {})
  } else {
    router.back()
  }
}

watch(isEncrypted, (val) => {
  if (!val) {
    password.value = ''
    confirmPassword.value = ''
  }
  markUnsaved()
})

watch([title, selectedTagIds], () => {
  markUnsaved()
}, { deep: true })

onMounted(() => {
  loadNote()
  loadTags()

  window.addEventListener('beforeunload', (e) => {
    if (hasUnsavedChanges.value) {
      e.preventDefault()
      e.returnValue = ''
    }
  })
})

onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()

  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  if (saveStatusTimer) clearTimeout(saveStatusTimer)
})
</script>

<style lang="scss" scoped>
.editor-page {
  background: #fff;
  border-radius: 8px;
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

.editor-header {
  padding: 16px 24px;
  border-bottom: 1px solid #ebeef5;
}

.title-input {
  width: 400px;

  :deep(.el-input__wrapper) {
    box-shadow: none;
  }

  :deep(.el-input__inner) {
    font-size: 20px;
    font-weight: 600;
  }
}

.save-status {
  font-size: 12px;
  margin-right: 16px;

  &.status-saving {
    color: #E6A23C;
  }

  &.status-saved {
    color: #67C23A;
  }

  &.status-unsaved {
    color: #F56C6C;
  }
}

.editor-toolbar {
  padding: 12px 24px;
  background: #fafafa;
  border-bottom: 1px solid #ebeef5;
}

.tag-color-dot {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.password-section {
  background: #fdf6ec;
  border-bottom: 1px solid #faecd8;
}

.editor-container {
  flex: 1;
  padding: 24px;
  overflow: auto;
}
</style>
