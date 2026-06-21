<template>
  <div class="board-container" @paste="handlePaste">
    <header class="board-header">
      <div class="header-left">
        <el-button text @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <el-input 
          v-if="editingName" 
          v-model="boardName" 
          size="small" 
          class="name-input"
          @blur="saveBoardName"
          @keyup.enter="saveBoardName"
          ref="nameInput"
        />
        <h2 v-else class="board-title" @click="startEditName">{{ boardName || '未命名白板' }}</h2>
      </div>
      
      <div class="header-center">
        <span class="connection-status" :class="{ connected: wsConnected }">
          <span class="dot"></span>
          {{ wsConnected ? '已连接' : '连接中...' }}
        </span>
      </div>
      
      <div class="header-right">
        <el-button size="small" @click="saveBoard">
          <el-icon><DocumentAdd /></el-icon>
          保存
        </el-button>
        <el-button size="small" @click="showHistory = true">
          <el-icon><Clock /></el-icon>
          历史
        </el-button>
        <el-button size="small" @click="showLayers = !showLayers">
          <el-icon><Collection /></el-icon>
          图层
        </el-button>
        <el-button size="small" @click="toggleFullscreen">
          <el-icon><FullScreen /></el-icon>
        </el-button>
      </div>
    </header>

    <div class="board-content">
      <Toolbar @clear="showClearDialog = true" />
      
      <div 
        class="canvas-wrapper"
        ref="canvasWrapper"
        @wheel="handleWheel"
      >
        <canvas 
          ref="canvasRef" 
          class="whiteboard-canvas"
          @mousedown="startDrawing"
          @mousemove="draw"
          @mouseup="stopDrawing"
          @mouseleave="stopDrawing"
          @touchstart="startDrawing"
          @touchmove="draw"
          @touchend="stopDrawing"
        ></canvas>
        
        <div 
          v-for="user in onlineUsers" 
          :key="user.id"
          class="user-cursor"
          v-if="user.id !== currentUser.id && user.cursor"
          :style="{
            left: user.cursor.x * zoom + pan.x + 'px',
            top: user.cursor.y * zoom + pan.y + 'px',
            borderColor: user.color
          }"
        >
          <div class="cursor-tip" :style="{ background: user.color }">{{ user.name }}</div>
        </div>
        
        <div 
          v-if="textInputVisible"
          class="text-input-wrapper"
          :style="{
            left: textInputPosition.x * zoom + pan.x + 'px',
            top: textInputPosition.y * zoom + pan.y + 'px'
          }"
        >
          <el-input
            v-model="textInputValue"
            size="small"
            placeholder="输入文字..."
            @blur="confirmText"
            @keyup.enter="confirmText"
            ref="textInput"
            class="text-input"
          />
        </div>
      </div>
      
      <UserList />
      
      <LayerPanel v-if="showLayers" @close="showLayers = false" />
      
      <HistoryPanel v-if="showHistory" @close="showHistory = false" />
    </div>
    
    <div class="zoom-controls">
      <el-button size="small" circle @click="zoomOut">
        <el-icon><Minus /></el-icon>
      </el-button>
      <span class="zoom-value">{{ Math.round(zoom * 100) }}%</span>
      <el-button size="small" circle @click="zoomIn">
        <el-icon><Plus /></el-icon>
      </el-button>
      <el-button size="small" @click="resetView">重置</el-button>
    </div>
    
    <el-dialog v-model="showClearDialog" title="清空白板" width="400px">
      <p>确定要清空当前白板的所有内容吗？此操作不可撤销。</p>
      <template #footer>
        <el-button @click="showClearDialog = false">取消</el-button>
        <el-button type="danger" @click="confirmClear">清空</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useBoardStore } from '@/store/board'
import { useCanvas } from '@/composables/useCanvas'
import { useWebSocket } from '@/composables/useWebSocket'
import { boardApi } from '@/utils/request'
import { debounce } from '@/utils/common'
import Toolbar from '@/components/Toolbar.vue'
import UserList from '@/components/UserList.vue'
import LayerPanel from '@/components/LayerPanel.vue'
import HistoryPanel from '@/components/HistoryPanel.vue'

const route = useRoute()
const router = useRouter()
const store = useBoardStore()

const canvasRef = ref(null)
const canvasWrapper = ref(null)
const nameInput = ref(null)
const textInput = ref(null)

const boardName = ref('')
const editingName = ref(false)
const showLayers = ref(false)
const showHistory = ref(false)
const showClearDialog = ref(false)
const isFullscreen = ref(false)

const {
  initCanvas,
  render,
  startDrawing,
  draw,
  stopDrawing,
  addText,
  addImage,
  textInputPosition,
  textInputVisible,
  textInputValue,
  handleWheel
} = useCanvas(canvasRef)

const {
  connected: wsConnected,
  connect,
  disconnect,
  sendClear
} = useWebSocket()

const currentUser = computed(() => store.currentUser)
const onlineUsers = computed(() => store.onlineUsers)
const zoom = computed(() => store.zoom)
const pan = computed(() => store.pan)

async function loadBoard() {
  const boardId = route.params.id
  store.setBoardId(boardId)
  
  try {
    const [boardRes, elementsRes] = await Promise.all([
      boardApi.get(boardId),
      boardApi.getElements(boardId)
    ])
    
    boardName.value = boardRes.data.name
    store.setBoardName(boardRes.data.name)
    store.setBackground(boardRes.data.background)
    store.setElements(elementsRes.data || [])
    
    connect(boardId, currentUser.value.name, currentUser.value.color)
    
    nextTick(() => {
      initCanvas()
    })
  } catch (e) {
    ElMessage.error('加载白板失败')
    router.push('/')
  }
}

function goBack() {
  router.push('/')
}

function startEditName() {
  editingName.value = true
  nextTick(() => {
    nameInput.value?.focus()
    nameInput.value?.select()
  })
}

async function saveBoardName() {
  editingName.value = false
  if (boardName.value.trim() && boardName.value !== store.boardName) {
    try {
      await boardApi.update(store.boardId, { name: boardName.value })
      store.setBoardName(boardName.value)
      ElMessage.success('已保存')
    } catch (e) {
      ElMessage.error('保存失败')
      boardName.value = store.boardName
    }
  } else {
    boardName.value = store.boardName
  }
}

const saveBoard = debounce(async () => {
  try {
    await boardApi.saveElements(store.boardId, store.elements)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败')
  }
}, 500)

function confirmText() {
  if (textInputValue.value.trim() && textInputPosition.value) {
    addText(textInputValue.value, textInputPosition.value)
  }
  textInputVisible.value = false
  textInputValue.value = ''
}

function handlePaste(e) {
  const items = e.clipboardData?.items
  if (!items) return
  
  for (const item of items) {
    if (item.type.startsWith('image/')) {
      e.preventDefault()
      const file = item.getAsFile()
      if (file) {
        const reader = new FileReader()
        reader.onload = (event) => {
          const img = new Image()
          img.onload = () => {
            const maxWidth = 400
            const maxHeight = 400
            let width = img.width
            let height = img.height
            
            if (width > maxWidth) {
              height = (maxWidth / width) * height
              width = maxWidth
            }
            if (height > maxHeight) {
              width = (maxHeight / height) * width
              height = maxHeight
            }
            
            const pos = {
              x: (canvasWrapper.value?.clientWidth || 800) / 2 / store.zoom - store.pan.x - width / 2,
              y: (canvasWrapper.value?.clientHeight || 600) / 2 / store.zoom - store.pan.y - height / 2
            }
            
            addImage(event.target?.result, pos, width, height)
            ElMessage.success('图片已粘贴')
          }
          img.src = event.target?.result
        }
        reader.readAsDataURL(file)
      }
      break
    }
  }
}

function zoomIn() {
  store.setZoom(store.zoom * 1.2)
  render()
}

function zoomOut() {
  store.setZoom(store.zoom / 1.2)
  render()
}

function resetView() {
  store.setZoom(1)
  store.setPan(0, 0)
  render()
}

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

function confirmClear() {
  store.clearElements()
  sendClear()
  render()
  showClearDialog.value = false
  ElMessage.success('白板已清空')
}

function handleKeydown(e) {
  if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
    return
  }
  
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    saveBoard()
  } else if ((e.ctrlKey || e.metaKey) && e.key === 'z') {
    e.preventDefault()
    if (store.undo()) {
      render()
    }
  } else if ((e.ctrlKey || e.metaKey) && e.key === 'y') {
    e.preventDefault()
    if (store.redo()) {
      render()
    }
  } else if (e.key === 'Delete' || e.key === 'Backspace') {
    if (store.selectedLayerId) {
      store.removeElement(store.selectedLayerId)
      store.selectLayer(null)
      render()
    }
  } else if (e.key === 'Escape') {
    store.setTool('pen')
  }
}

onMounted(() => {
  loadBoard()
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  disconnect()
  window.removeEventListener('keydown', handleKeydown)
})

watch(() => store.currentTool, (tool) => {
  if (canvasRef.value) {
    const cursors = {
      pen: 'crosshair',
      eraser: 'crosshair',
      line: 'crosshair',
      rect: 'crosshair',
      circle: 'crosshair',
      ellipse: 'crosshair',
      triangle: 'crosshair',
      arrow: 'crosshair',
      text: 'text',
      pan: 'grab',
      select: 'default'
    }
    canvasRef.value.style.cursor = cursors[tool] || 'default'
  }
})

defineExpose({
  showClearDialog
})
</script>

<style lang="scss" scoped>
.board-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  overflow: hidden;
}

.board-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  z-index: 100;
  
  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .name-input {
      width: 200px;
    }
    
    .board-title {
      font-size: 18px;
      margin: 0;
      color: #303133;
      cursor: pointer;
      padding: 4px 8px;
      border-radius: 4px;
      transition: background 0.2s;
      
      &:hover {
        background: #f5f7fa;
      }
    }
  }
  
  .header-center {
    .connection-status {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: #909399;
      
      .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #f56c6c;
        animation: pulse 2s infinite;
      }
      
      &.connected {
        color: #67c23a;
        
        .dot {
          background: #67c23a;
          animation: none;
        }
      }
    }
  }
  
  .header-right {
    display: flex;
    gap: 8px;
  }
}

.board-content {
  flex: 1;
  display: flex;
  position: relative;
  overflow: hidden;
}

.canvas-wrapper {
  flex: 1;
  position: relative;
  overflow: hidden;
  background: #e9e9e9;
  background-image: 
    linear-gradient(45deg, #dcdcdc 25%, transparent 25%),
    linear-gradient(-45deg, #dcdcdc 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #dcdcdc 75%),
    linear-gradient(-45deg, transparent 75%, #dcdcdc 75%);
  background-size: 20px 20px;
  background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
}

.whiteboard-canvas {
  position: absolute;
  top: 0;
  left: 0;
  touch-action: none;
}

.user-cursor {
  position: absolute;
  width: 12px;
  height: 12px;
  border: 2px solid;
  border-radius: 50%;
  pointer-events: none;
  z-index: 1000;
  transform: translate(-50%, -50%);
  transition: left 0.05s linear, top 0.05s linear;
  
  .cursor-tip {
    position: absolute;
    top: 14px;
    left: 8px;
    padding: 2px 6px;
    border-radius: 4px;
    color: #fff;
    font-size: 11px;
    white-space: nowrap;
    opacity: 0.9;
  }
}

.text-input-wrapper {
  position: absolute;
  z-index: 1000;
  transform: translate(-4px, -4px);
  
  .text-input {
    min-width: 150px;
  }
}

.zoom-controls {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  padding: 8px 12px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
  
  .zoom-value {
    font-size: 13px;
    color: #606266;
    min-width: 50px;
    text-align: center;
  }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
