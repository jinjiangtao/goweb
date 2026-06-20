<template>
  <el-drawer
    v-model="visible"
    title="版本历史"
    direction="rtl"
    size="420px"
    :before-close="handleClose"
  >
    <div class="version-panel">
      <div v-if="loading" class="loading-wrapper">
        <el-icon class="is-loading"><Loading /></el-icon>
        <span class="loading-text">加载中...</span>
      </div>

      <el-empty v-else-if="versions.length === 0" description="暂无版本历史" />

      <el-timeline v-else>
        <el-timeline-item
          v-for="version in sortedVersions"
          :key="version.id"
          :timestamp="formatTime(version.created_at)"
          placement="top"
          :type="version.is_current ? 'primary' : ''"
        >
          <div
            class="version-item"
            :class="{ 'version-current': version.is_current }"
          >
            <div class="version-header">
              <div class="version-title">
                <span class="version-number">v{{ version.version_number }}</span>
                <el-tag
                  v-if="version.is_current"
                  type="primary"
                  size="small"
                  effect="dark"
                >
                  当前
                </el-tag>
              </div>
              <el-button
                v-if="!version.is_current"
                type="primary"
                link
                size="small"
                @click="handleRestore(version)"
              >
                恢复到此版本
              </el-button>
            </div>
            <div class="version-snapshot">{{ version.snapshot_name }}</div>
            <div v-if="version.description" class="version-desc">
              {{ version.description }}
            </div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { useResumeStore } from '@/stores/resume'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  resumeId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'restored'])

const store = useResumeStore()
const visible = ref(props.modelValue)
const loading = ref(false)
const versions = ref([])

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val
    if (val && props.resumeId) {
      loadVersions()
    }
  }
)

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const sortedVersions = computed(() => {
  return [...versions.value].sort((a, b) => b.version_number - a.version_number)
})

async function loadVersions() {
  if (!props.resumeId) return
  loading.value = true
  try {
    const data = await store.loadVersions(props.resumeId)
    versions.value = data || []
  } catch (e) {
    ElMessage.error('获取版本历史失败')
  } finally {
    loading.value = false
  }
}

async function handleRestore(version) {
  try {
    await ElMessageBox.confirm(
      `确定要恢复到版本 v${version.version_number} 吗？\n当前未保存的修改将会丢失。`,
      '恢复版本',
      {
        confirmButtonText: '确定恢复',
        cancelButtonText: '取消',
        type: 'warning',
        distinguishCancelAndClose: true
      }
    )
    loading.value = true
    await store.restoreVersion(version.id)
    ElMessage.success('恢复成功')
    emit('restored')
    await loadVersions()
  } catch (e) {
    if (e !== 'cancel' && e.action !== 'cancel') {
      ElMessage.error('恢复失败')
    }
  } finally {
    loading.value = false
  }
}

function handleClose(done) {
  done()
}

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style lang="scss" scoped>
.version-panel {
  padding: 8px 0;

  .loading-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 0;
    color: #909399;

    .is-loading {
      font-size: 32px;
      margin-bottom: 12px;
    }

    .loading-text {
      font-size: 14px;
    }
  }

  .version-item {
    padding: 12px 16px;
    border: 1px solid #e4e7ed;
    border-radius: 8px;
    transition: all 0.2s;

    &:hover {
      border-color: #c0c4cc;
    }

    &.version-current {
      background: #ecf5ff;
      border-color: #409eff;
    }

    .version-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;

      .version-title {
        display: flex;
        align-items: center;
        gap: 8px;

        .version-number {
          font-weight: 600;
          font-size: 15px;
          color: #303133;
        }
      }
    }

    .version-snapshot {
      font-size: 14px;
      color: #303133;
      margin-bottom: 4px;
      line-height: 1.5;
    }

    .version-desc {
      font-size: 13px;
      color: #909399;
      line-height: 1.5;
    }
  }
}

:deep(.el-timeline-item__timestamp) {
  color: #909399;
  font-size: 12px;
}

:deep(.el-drawer__body) {
  padding: 16px 20px;
}
</style>
