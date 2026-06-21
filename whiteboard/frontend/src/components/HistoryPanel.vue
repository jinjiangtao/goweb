<template>
  <div class="panel-mask" @click="$emit('close')">
    <div class="panel" @click.stop>
      <div class="panel-header">
        <h3 class="title">
          <el-icon color="#409eff"><Clock /></el-icon>
          历史记录
        </h3>
        <el-button text circle @click="$emit('close')">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>
      
      <div class="panel-content">
        <div 
          v-for="record in history" 
          :key="record.id"
          class="history-item"
        >
          <div class="history-icon">
            <el-icon>
              <component :is="getIcon(record.action)" />
            </el-icon>
          </div>
          <div class="history-info">
            <span class="history-action">{{ getActionName(record.action) }}</span>
            <span class="history-user">
              <span class="avatar" :style="{ background: getUserColor(record.userId) }">
                {{ record.userName.charAt(0) }}
              </span>
              {{ record.userName }}
            </span>
            <span class="history-time">{{ formatTime(record.timestamp) }}</span>
          </div>
          <div class="history-action-btn">
            <el-button 
              size="small" 
              type="primary" 
              @click="restoreHistory(record.id)"
            >
              恢复
            </el-button>
          </div>
        </div>
        
        <el-empty 
          v-if="history.length === 0" 
          description="暂无历史记录" 
          :image-size="60"
        >
          <template #image>
            <el-icon :size="48" color="#dcdfe6"><Clock /></el-icon>
          </template>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useBoardStore } from '@/store/board'
import { boardApi } from '@/utils/request'
import { formatTime as formatTimeUtil, generateColor } from '@/utils/common'
import { ElMessage, ElMessageBox } from 'element-plus'

const store = useBoardStore()
const emit = defineEmits(['close'])

const history = ref([])
const userColors = ref({})

async function loadHistory() {
  try {
    const res = await boardApi.getHistory(store.boardId)
    history.value = res.data || []
  } catch (e) {
    ElMessage.error('加载历史记录失败')
  }
}

async function restoreHistory(historyId) {
  try {
    await ElMessageBox.confirm('确定要恢复到此历史版本吗？当前内容将被覆盖。', '提示', {
      confirmButtonText: '恢复',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await boardApi.restoreHistory(store.boardId, historyId)
    ElMessage.success('恢复成功')
    emit('close')
  } catch {
    // 用户取消
  }
}

function getIcon(action) {
  const icons = {
    save: 'DocumentCheck',
    clear: 'Delete',
    create: 'Plus',
    update: 'Edit'
  }
  return icons[action] || 'Clock'
}

function getActionName(action) {
  const names = {
    save: '保存白板',
    clear: '清空白板',
    create: '创建白板',
    update: '更新白板'
  }
  return names[action] || action
}

function getUserColor(userId) {
  if (!userColors.value[userId]) {
    userColors.value[userId] = generateColor()
  }
  return userColors.value[userId]
}

function formatTime(timestamp) {
  return formatTimeUtil(timestamp)
}

onMounted(() => {
  loadHistory()
})
</script>

<style lang="scss" scoped>
.panel-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
}

.panel {
  width: 500px;
  max-height: 80vh;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  
  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid #e4e7ed;
    
    .title {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 16px;
      margin: 0;
      color: #303133;
    }
  }
  
  .panel-content {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }
  
  .history-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    border-radius: 6px;
    border: 1px solid #e4e7ed;
    margin-bottom: 8px;
    transition: all 0.2s;
    
    &:hover {
      background: #f5f7fa;
    }
    
    .history-icon {
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #ecf5ff;
      border-radius: 50%;
      color: #409eff;
      flex-shrink: 0;
    }
    
    .history-info {
      flex: 1;
      min-width: 0;
      
      .history-action {
        display: block;
        font-size: 14px;
        color: #303133;
        font-weight: 500;
        margin-bottom: 4px;
      }
      
      .history-user {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        color: #606266;
        margin-bottom: 2px;
        
        .avatar {
          width: 18px;
          height: 18px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
          font-size: 10px;
        }
      }
      
      .history-time {
        font-size: 11px;
        color: #909399;
      }
    }
  }
}

.panel-content::-webkit-scrollbar {
  width: 6px;
}

.panel-content::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}
</style>
