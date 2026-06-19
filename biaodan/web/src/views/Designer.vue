<template>
  <div class="designer-container">
    <header class="designer-header">
      <div class="header-left">
        <el-button :icon="ArrowLeft" link @click="$router.push('/')">
          返回
        </el-button>
        <el-input
          v-model="formStore.name"
          class="form-name-input"
          placeholder="请输入表单名称"
          size="large"
        />
      </div>
      <div class="header-right">
        <el-button :icon="View" @click="showPreview = true">预览</el-button>
        <el-button :icon="Share" @click="copyShareLink" v-if="formStore.templateId">分享</el-button>
        <el-button :icon="Refresh" @click="resetConfirm">重置</el-button>
        <el-button type="primary" :icon="Check" @click="saveTemplate">保存</el-button>
      </div>
    </header>

    <div class="designer-main">
      <aside class="sidebar left-sidebar">
        <div class="sidebar-title">组件库</div>
        <div class="component-palette">
          <div v-for="group in componentTypes" :key="group.group" class="component-group">
            <div class="group-title">{{ group.group }}</div>
            <div class="group-items">
              <div
                v-for="comp in group.items"
                :key="comp.type"
                class="component-item"
                draggable="true"
                @dragstart="onDragStart($event, comp.type)"
                @click="addComponent(comp.type)"
              >
                <el-icon><component :is="comp.icon" /></el-icon>
                <span>{{ comp.label }}</span>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <section class="canvas-area">
        <div
          class="canvas-wrapper"
          @dragover.prevent="onDragOver"
          @drop="onDrop"
          @click.self="deselectField"
        >
          <div class="form-header-info">
            <input
              v-model="formStore.name"
              class="form-title-input"
              placeholder="请输入表单标题"
            />
            <textarea
              v-model="formStore.description"
              class="form-desc-input"
              placeholder="请输入表单描述（可选）"
              rows="2"
            />
          </div>

          <div v-if="formStore.fields.length === 0" class="empty-canvas">
            <el-empty description="从左侧拖拽组件到这里，或点击组件添加" />
          </div>

          <draggable-container
            v-model="formStore.fields"
            item-key="id"
            class="fields-container"
            :animation="200"
            ghost-class="ghost"
          >
            <template #item="{ element, index }">
              <div
                class="field-wrapper"
                :class="{ selected: formStore.selectedFieldId === element.id }"
                @click.stop="formStore.selectField(element.id)"
                @dragover.prevent
                @drop.stop="onFieldDrop($event, index)"
              >
                <div class="field-actions">
                  <el-button size="small" :icon="Top" text @click.stop="moveUp(index)" :disabled="index === 0" />
                  <el-button size="small" :icon="Bottom" text @click.stop="moveDown(index)" :disabled="index === formStore.fields.length - 1" />
                  <el-button size="small" :icon="CopyDocument" text type="primary" @click.stop="formStore.duplicateField(element.id)" />
                  <el-button size="small" :icon="Delete" text type="danger" @click.stop="removeField(element.id)" />
                </div>
                <form-renderer :field="element" :preview="true" />
              </div>
            </template>
          </draggable-container>
        </div>
      </section>

      <aside class="sidebar right-sidebar">
        <div class="sidebar-title">属性配置</div>
        <div v-if="formStore.selectedField" class="property-panel">
          <property-config :field="formStore.selectedField" />
        </div>
        <el-empty v-else description="请选择一个组件" :image-size="80" />
      </aside>
    </div>

    <el-dialog v-model="showPreview" title="表单预览" width="720px" top="4vh">
      <div class="preview-dialog">
        <div class="preview-form">
          <h2 class="preview-title">{{ formStore.name || '表单标题' }}</h2>
          <p v-if="formStore.description" class="preview-desc">{{ formStore.description }}</p>
          <div class="preview-fields">
            <form-renderer
              v-for="field in formStore.fields"
              :key="field.id"
              :field="field"
              v-model="previewData"
            />
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="showPreview = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, View, Share, Refresh, Check,
  Top, Bottom, CopyDocument, Delete
} from '@element-plus/icons-vue'
import { useFormStore, componentTypes } from '@/stores/form'
import { getTemplate, createTemplate, updateTemplate } from '@/api'
import FormRenderer from '@/components/FormRenderer.vue'
import PropertyConfig from '@/components/PropertyConfig.vue'
import DraggableContainer from '@/components/DraggableContainer.vue'

const route = useRoute()
const router = useRouter()
const formStore = useFormStore()
const showPreview = ref(false)
const previewData = reactive({})
const draggedType = ref('')
const dropIndex = ref(-1)

function addComponent(type) {
  formStore.addField(type)
}

function removeField(id) {
  ElMessageBox.confirm('确定删除此组件吗？', '删除确认', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
    formStore.removeField(id)
    ElMessage.success('已删除')
  }).catch(() => {})
}

function moveUp(index) {
  formStore.moveField(index, index - 1)
}

function moveDown(index) {
  formStore.moveField(index, index + 1)
}

function deselectField() {
  formStore.selectField(null)
}

function onDragStart(e, type) {
  draggedType.value = type
  e.dataTransfer.effectAllowed = 'copy'
}

function onDragOver(e) {
  e.dataTransfer.dropEffect = 'copy'
}

function onDrop(e) {
  if (draggedType.value) {
    formStore.addField(draggedType.value)
    draggedType.value = ''
  }
}

function onFieldDrop(e, index) {
  if (draggedType.value) {
    formStore.addField(draggedType.value, index)
    draggedType.value = ''
  }
}

async function saveTemplate() {
  if (!formStore.name.trim()) {
    ElMessage.warning('请输入表单名称')
    return
  }
  if (formStore.fields.length === 0) {
    ElMessage.warning('请至少添加一个表单组件')
    return
  }

  const data = {
    name: formStore.name,
    description: formStore.description,
    schema: formStore.getSchema()
  }

  try {
    let res
    if (formStore.templateId) {
      res = await updateTemplate(formStore.templateId, data)
      ElMessage.success('更新成功')
    } else {
      res = await createTemplate(data)
      formStore.templateId = res.data.id
      ElMessage.success('创建成功')
    }
  } catch (e) {}
}

function resetConfirm() {
  ElMessageBox.confirm('确定要重置所有内容吗？此操作不可撤销。', '重置确认', {
    type: 'warning'
  }).then(() => {
    const id = formStore.templateId
    formStore.reset()
    if (id) {
      loadTemplate(id)
    }
    ElMessage.success('已重置')
  }).catch(() => {})
}

async function loadTemplate(id) {
  try {
    const res = await getTemplate(id)
    formStore.loadTemplate(res.data)
  } catch (e) {}
}

function copyShareLink() {
  const url = `${window.location.origin}/fill/${formStore.templateId}`
  navigator.clipboard.writeText(url)
  ElMessage.success('分享链接已复制')
}

onMounted(() => {
  const id = route.params.id
  if (id) {
    loadTemplate(id)
  } else {
    formStore.reset()
    formStore.name = '新建表单'
  }
})
</script>

<style scoped>
.designer-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f0f2f5;
}

.designer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.form-name-input {
  width: 360px;
}

.header-right {
  display: flex;
  gap: 8px;
}

.designer-main {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 260px;
  background: #fff;
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
}

.sidebar.left-sidebar {
  border-right: 1px solid #e4e7ed;
}

.sidebar.right-sidebar {
  width: 320px;
  border-right: none;
  border-left: 1px solid #e4e7ed;
}

.sidebar-title {
  padding: 16px;
  font-weight: 600;
  font-size: 14px;
  border-bottom: 1px solid #f0f2f5;
  background: #fafbfc;
}

.component-palette {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.component-group {
  margin-bottom: 16px;
}

.group-title {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
  padding: 0 4px;
}

.group-items {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.component-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 14px 8px;
  background: #fafbfc;
  border: 1px dashed #dcdfe6;
  border-radius: 6px;
  cursor: grab;
  transition: all 0.2s;
  color: #606266;
  font-size: 12px;
}

.component-item:hover {
  background: #ecf5ff;
  border-color: #409EFF;
  color: #409EFF;
}

.component-item:active {
  cursor: grabbing;
}

.canvas-area {
  flex: 1;
  overflow: auto;
  padding: 24px;
}

.canvas-wrapper {
  max-width: 720px;
  margin: 0 auto;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  padding: 32px;
  min-height: calc(100vh - 180px);
}

.form-header-info {
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f2f5;
}

.form-title-input {
  width: 100%;
  font-size: 24px;
  font-weight: 600;
  border: none;
  outline: none;
  background: transparent;
  margin-bottom: 12px;
  color: #303133;
}

.form-title-input::placeholder {
  color: #c0c4cc;
}

.form-desc-input {
  width: 100%;
  font-size: 14px;
  color: #606266;
  border: none;
  outline: none;
  background: transparent;
  resize: none;
  font-family: inherit;
  line-height: 1.6;
}

.form-desc-input::placeholder {
  color: #c0c4cc;
}

.empty-canvas {
  padding: 60px 0;
  border: 2px dashed #dcdfe6;
  border-radius: 8px;
  margin-top: 20px;
}

.fields-container {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.field-wrapper {
  position: relative;
  padding: 12px 14px;
  border: 2px solid transparent;
  border-radius: 6px;
  transition: all 0.15s;
  background: #fff;
}

.field-wrapper:hover {
  background: #fafbfc;
  border-color: #d9ecff;
}

.field-wrapper.selected {
  border-color: #409EFF;
  background: #f5f9ff;
}

.field-actions {
  position: absolute;
  top: 4px;
  right: 4px;
  display: none;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  z-index: 10;
}

.field-wrapper:hover .field-actions,
.field-wrapper.selected .field-actions {
  display: flex;
}

.property-panel {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.preview-dialog {
  max-height: 70vh;
  overflow-y: auto;
}

.preview-form {
  padding: 12px;
}

.preview-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 12px;
}

.preview-desc {
  color: #606266;
  margin: 0 0 24px;
  line-height: 1.6;
}

.preview-fields {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ghost {
  opacity: 0.5;
  background: #f0f7ff;
}
</style>
