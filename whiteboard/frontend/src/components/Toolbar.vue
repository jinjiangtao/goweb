<template>
  <div class="toolbar">
    <div class="tool-group">
      <el-tooltip content="画笔" placement="right">
        <el-button 
          :type="currentTool === 'pen' ? 'primary' : 'default'" 
          circle
          @click="setTool('pen')"
        >
          <el-icon><EditPen /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="橡皮擦" placement="right">
        <el-button 
          :type="currentTool === 'eraser' ? 'primary' : 'default'" 
          circle
          @click="setTool('eraser')"
        >
          <el-icon><MagicStick /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="文字" placement="right">
        <el-button 
          :type="currentTool === 'text' ? 'primary' : 'default'" 
          circle
          @click="setTool('text')"
        >
          <el-icon><Edit /></el-icon>
        </el-button>
      </el-tooltip>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <el-tooltip content="直线" placement="right">
        <el-button 
          :type="currentTool === 'line' ? 'primary' : 'default'" 
          circle
          @click="setTool('line')"
        >
          <el-icon><Minus /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="箭头" placement="right">
        <el-button 
          :type="currentTool === 'arrow' ? 'primary' : 'default'" 
          circle
          @click="setTool('arrow')"
        >
          <el-icon><Right /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="矩形" placement="right">
        <el-button 
          :type="currentTool === 'rect' ? 'primary' : 'default'" 
          circle
          @click="setTool('rect')"
        >
          <el-icon><Grid /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="圆形" placement="right">
        <el-button 
          :type="currentTool === 'circle' ? 'primary' : 'default'" 
          circle
          @click="setTool('circle')"
        >
          <el-icon><CircleCheck /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="椭圆" placement="right">
        <el-button 
          :type="currentTool === 'ellipse' ? 'primary' : 'default'" 
          circle
          @click="setTool('ellipse')"
        >
          <el-icon><Odometer /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="三角形" placement="right">
        <el-button 
          :type="currentTool === 'triangle' ? 'primary' : 'default'" 
          circle
          @click="setTool('triangle')"
        >
          <el-icon><CaretTop /></el-icon>
        </el-button>
      </el-tooltip>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <el-tooltip content="抓手工具" placement="right">
        <el-button 
          :type="currentTool === 'pan' ? 'primary' : 'default'" 
          circle
          @click="setTool('pan')"
        >
          <el-icon><Rank /></el-icon>
        </el-button>
      </el-tooltip>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <el-tooltip content="撤销" placement="right">
        <el-button 
          circle
          @click="handleUndo"
          :disabled="!canUndo"
        >
          <el-icon><RefreshLeft /></el-icon>
        </el-button>
      </el-tooltip>
      
      <el-tooltip content="重做" placement="right">
        <el-button 
          circle
          @click="handleRedo"
          :disabled="!canRedo"
        >
          <el-icon><RefreshRight /></el-icon>
        </el-button>
      </el-tooltip>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <div class="color-picker-wrapper">
        <el-tooltip content="线条颜色" placement="right">
          <el-color-picker 
            v-model="strokeColor" 
            size="small"
            @change="updateStyle"
          />
        </el-tooltip>
      </div>
      
      <div class="color-picker-wrapper">
        <el-tooltip content="填充颜色" placement="right">
          <el-color-picker 
            v-model="fillColor" 
            size="small"
            @change="updateStyle"
          />
        </el-tooltip>
      </div>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <el-tooltip content="线条粗细" placement="right">
        <div class="slider-wrapper">
          <el-icon><Aim /></el-icon>
          <el-slider 
            v-model="lineWidth" 
            :min="1" 
            :max="50" 
            :step="1"
            class="width-slider"
            @change="updateStyle"
          />
          <span class="slider-value">{{ lineWidth }}px</span>
        </div>
      </el-tooltip>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <el-tooltip content="透明度" placement="right">
        <div class="slider-wrapper">
          <el-icon><View /></el-icon>
          <el-slider 
            v-model="opacity" 
            :min="0.1" 
            :max="1" 
            :step="0.1"
            class="width-slider"
            @change="updateStyle"
          />
          <span class="slider-value">{{ Math.round(opacity * 100) }}%</span>
        </div>
      </el-tooltip>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <el-tooltip content="背景颜色" placement="right">
        <el-color-picker 
          v-model="backgroundColor" 
          size="small"
          @change="updateBackground"
        />
      </el-tooltip>
    </div>
    
    <div class="divider"></div>
    
    <div class="tool-group">
      <el-tooltip content="粘贴图片 (Ctrl+V)" placement="right">
        <el-button 
          circle
          @click="triggerFileInput"
        >
          <el-icon><Picture /></el-icon>
        </el-button>
      </el-tooltip>
      <input 
        type="file" 
        ref="fileInput" 
        accept="image/*" 
        style="display: none"
        @change="handleFileSelect"
      />
      
      <el-tooltip content="清空画布" placement="right">
        <el-button 
          type="danger"
          circle
          @click="handleClear"
        >
          <el-icon><Delete /></el-icon>
        </el-button>
      </el-tooltip>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useBoardStore } from '@/store/board'
import { useWebSocket } from '@/composables/useWebSocket'
import { ElMessage } from 'element-plus'
import { boardApi } from '@/utils/request'

const store = useBoardStore()
const ws = useWebSocket()
const fileInput = ref(null)

const strokeColor = ref(store.currentStyle.color)
const fillColor = ref(store.currentStyle.fill || 'transparent')
const lineWidth = ref(store.currentStyle.width)
const opacity = ref(store.currentStyle.opacity)
const backgroundColor = ref(store.background)

const currentTool = computed(() => store.currentTool)
const canUndo = computed(() => store.elements.length > 0)
const canRedo = computed(() => false)

const emit = defineEmits(['clear'])

function setTool(tool) {
  store.setTool(tool)
}

function updateStyle() {
  store.setStyle({
    color: strokeColor.value,
    width: lineWidth.value,
    opacity: opacity.value,
    fill: fillColor.value,
    fillOpacity: 0.5
  })
}

function updateBackground(color) {
  store.setBackground(color)
  boardApi.update(store.boardId, { background: color })
}

function handleUndo() {
  if (store.undo()) {
    ws.sendUndo()
  }
}

function handleRedo() {
  if (store.redo()) {
    ws.sendRedo()
  }
}

function handleClear() {
  emit('clear')
}

function triggerFileInput() {
  fileInput.value?.click()
}

function handleFileSelect(e) {
  const input = e.target
  const file = input.files?.[0]
  if (file && file.type.startsWith('image/')) {
    const reader = new FileReader()
    reader.onload = (event) => {
      const img = new Image()
      img.onload = () => {
        const event = new ClipboardEvent('paste', {
          clipboardData: new DataTransfer()
        })
        event.clipboardData?.items.add(file)
        document.dispatchEvent(event)
        
        ElMessage.success('图片已添加')
      }
      img.src = event.target?.result
    }
    reader.readAsDataURL(file)
  }
  input.value = ''
}

watch(() => store.currentStyle, (style) => {
  strokeColor.value = style.color
  lineWidth.value = style.width
  opacity.value = style.opacity
  if (style.fill) fillColor.value = style.fill
}, { deep: true })

watch(() => store.background, (bg) => {
  backgroundColor.value = bg
})
</script>

<style lang="scss" scoped>
.toolbar {
  width: 64px;
  background: #fff;
  border-right: 1px solid #e4e7ed;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  overflow-y: auto;
  z-index: 50;
  
  .tool-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: center;
  }
  
  .divider {
    width: 40px;
    height: 1px;
    background: #e4e7ed;
    margin: 4px 0;
  }
  
  .color-picker-wrapper {
    :deep(.el-color-picker__trigger) {
      width: 32px;
      height: 32px;
    }
  }
  
  .slider-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: 4px;
    
    .width-slider {
      width: 40px;
      
      :deep(.el-slider__runway) {
        margin: 0;
      }
    }
    
    .slider-value {
      font-size: 10px;
      color: #909399;
    }
  }
}

.toolbar::-webkit-scrollbar {
  width: 4px;
}

.toolbar::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 2px;
}
</style>
