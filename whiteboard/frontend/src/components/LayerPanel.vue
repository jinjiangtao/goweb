<template>
  <div class="panel-mask" @click="$emit('close')">
    <div class="panel" @click.stop>
      <div class="panel-header">
        <h3 class="title">
          <el-icon color="#409eff"><Collection /></el-icon>
          图层管理
        </h3>
        <el-button text circle @click="$emit('close')">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>
      
      <div class="panel-content">
        <div 
          v-for="(element, index) in reversedElements" 
          :key="element.id"
          class="layer-item"
          :class="{ selected: selectedLayerId === element.id }"
          @click="selectLayer(element.id)"
        >
          <div class="layer-icon">
            <el-icon>
              <component :is="getIcon(element.type)" />
            </el-icon>
          </div>
          <div class="layer-info">
            <span class="layer-name">{{ getTypeName(element.type) }}</span>
            <span class="layer-user">用户: {{ getUserName(element.userId) }}</span>
          </div>
          <div class="layer-actions">
            <el-button 
              size="small" 
              text 
              :disabled="index === reversedElements.length - 1"
              @click.stop="moveLayerUp(element.id)"
            >
              <el-icon><Top /></el-icon>
            </el-button>
            <el-button 
              size="small" 
              text 
              :disabled="index === 0"
              @click.stop="moveLayerDown(element.id)"
            >
              <el-icon><Bottom /></el-icon>
            </el-button>
            <el-button 
              size="small" 
              text 
              type="danger"
              @click.stop="deleteLayer(element.id)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        
        <el-empty 
          v-if="elements.length === 0" 
          description="暂无图层" 
          :image-size="60"
        >
          <template #image>
            <el-icon :size="48" color="#dcdfe6"><Collection /></el-icon>
          </template>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useBoardStore } from '@/store/board'
import { useCanvas } from '@/composables/useCanvas'

const store = useBoardStore()
const emit = defineEmits(['close'])

const elements = computed(() => store.elements)
const reversedElements = computed(() => [...store.elements].reverse())
const selectedLayerId = computed(() => store.selectedLayerId)

function selectLayer(id) {
  store.selectLayer(id === selectedLayerId.value ? null : id)
}

function moveLayerUp(id) {
  store.moveLayerUp(id)
}

function moveLayerDown(id) {
  store.moveLayerDown(id)
}

function deleteLayer(id) {
  store.removeElement(id)
  store.selectLayer(null)
}

function getIcon(type) {
  const icons = {
    pen: 'EditPen',
    line: 'Minus',
    rect: 'Grid',
    circle: 'CircleCheck',
    ellipse: 'Odometer',
    triangle: 'CaretTop',
    arrow: 'Right',
    text: 'Edit',
    image: 'Picture'
  }
  return icons[type] || 'EditPen'
}

function getTypeName(type) {
  const names = {
    pen: '画笔',
    line: '直线',
    rect: '矩形',
    circle: '圆形',
    ellipse: '椭圆',
    triangle: '三角形',
    arrow: '箭头',
    text: '文字',
    image: '图片',
    eraser: '橡皮擦'
  }
  return names[type] || type
}

function getUserName(userId) {
  const user = store.users.find(u => u.id === userId)
  return user?.name || '未知用户'
}
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
  width: 400px;
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
  
  .layer-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    border-radius: 6px;
    border: 1px solid transparent;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      background: #f5f7fa;
    }
    
    &.selected {
      background: #ecf5ff;
      border-color: #409eff;
    }
    
    .layer-icon {
      width: 32px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f5f7fa;
      border-radius: 4px;
      color: #909399;
    }
    
    .layer-info {
      flex: 1;
      min-width: 0;
      
      .layer-name {
        display: block;
        font-size: 13px;
        color: #303133;
        font-weight: 500;
      }
      
      .layer-user {
        font-size: 11px;
        color: #909399;
      }
    }
    
    .layer-actions {
      display: flex;
      gap: 4px;
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
