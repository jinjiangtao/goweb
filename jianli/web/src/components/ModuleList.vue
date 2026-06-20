<template>
  <div class="module-list">
    <div class="module-list-header">
      <span class="title">模块列表</span>
    </div>

    <draggable
      v-model="localModules"
      item-key="id"
      handle=".drag-handle"
      class="module-drag-list"
      @end="onDragEnd"
    >
      <template #item="{ element }">
        <div
          class="module-item"
          :class="{
            active: element.id === selectedModuleId,
            hidden: !element.visible
          }"
          @click="selectModule(element.id)"
        >
          <span class="drag-handle">
            <el-icon><Sort /></el-icon>
          </span>
          <span class="module-name">{{ element.name }}</span>
          <div class="module-actions">
            <el-switch
              :model-value="element.visible"
              size="small"
              @click.stop
              @change="toggleVisibility(element.id)"
            />
            <el-button
              v-if="element.removable"
              type="danger"
              link
              size="small"
              @click.stop="removeModule(element.id)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
      </template>
    </draggable>

    <div class="add-module-section">
      <el-dropdown trigger="click" @command="handleAddModule">
        <el-button type="primary" class="add-btn">
          <el-icon><Plus /></el-icon>
          添加模块
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="mod in availableModules"
              :key="mod.type"
              :command="mod.type"
            >
              {{ mod.name }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import draggable from 'vuedraggable'
import { useResumeStore } from '@/stores/resume'
import { Sort, Delete, Plus } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'

const store = useResumeStore()
const localModules = ref([])

watch(
  () => store.modules,
  (newVal) => {
    localModules.value = [...newVal]
  },
  { immediate: true, deep: true }
)

const selectedModuleId = computed(() => store.selectedModuleId)

const availableModules = computed(() => store.getAvailableModules())

function selectModule(id) {
  store.selectModule(id)
}

function toggleVisibility(id) {
  store.toggleModuleVisibility(id)
}

function removeModule(id) {
  ElMessageBox.confirm('确定要删除该模块吗？删除后数据将无法恢复。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    store.removeModule(id)
  }).catch(() => {})
}

function handleAddModule(type) {
  store.addModule(type)
  ElMessage.success(`已添加${availableModules.value.find(m => m.type === type)?.name || '模块'}`)
}

function onDragEnd() {
  store.reorderModules([...localModules.value])
}
</script>

<style lang="scss" scoped>
.module-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-right: 1px solid #e4e7ed;

  .module-list-header {
    padding: 16px;
    border-bottom: 1px solid #e4e7ed;
    font-weight: 600;
    font-size: 14px;
  }

  .module-drag-list {
    flex: 1;
    padding: 8px;
    overflow-y: auto;
  }

  .module-item {
    display: flex;
    align-items: center;
    padding: 10px 12px;
    margin-bottom: 4px;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #f5f7fa;
    }

    &.active {
      background: #ecf5ff;
      color: #409eff;
    }

    &.hidden {
      opacity: 0.5;
    }

    .drag-handle {
      margin-right: 8px;
      cursor: grab;
      color: #c0c4cc;

      &:active {
        cursor: grabbing;
      }
    }

    .module-name {
      flex: 1;
      font-size: 14px;
    }

    .module-actions {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }

  .add-module-section {
    padding: 12px;
    border-top: 1px solid #e4e7ed;

    .add-btn {
      width: 100%;
    }
  }
}
</style>
