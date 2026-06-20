<template>
  <div class="style-panel">
    <el-collapse v-model="activeNames">
      <el-collapse-item title="样式调整" name="style">
        <el-form label-width="100px" class="style-form">
          <el-form-item label="主色调">
            <el-color-picker
              v-model="primaryColor"
              @change="(val) => updateStyle('primary_color', val)"
            />
          </el-form-item>
          <el-form-item :label="`字体大小: ${fontSize}px`">
            <el-slider
              v-model="fontSize"
              :min="10"
              :max="16"
              :step="1"
              @change="(val) => updateStyle('font_size', val)"
            />
          </el-form-item>
          <el-form-item :label="`行高: ${lineHeight}`">
            <el-slider
              v-model="lineHeight"
              :min="1.2"
              :max="2.0"
              :step="0.1"
              @change="(val) => updateStyle('line_height', val)"
            />
          </el-form-item>
          <el-form-item :label="`页面边距: ${pageMargin}px`">
            <el-slider
              v-model="pageMargin"
              :min="10"
              :max="40"
              :step="2"
              @change="(val) => updateStyle('page_margin', val)"
            />
          </el-form-item>
        </el-form>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useResumeStore } from '@/stores/resume'

const store = useResumeStore()
const activeNames = ref(['style'])

const primaryColor = computed({
  get: () => store.styleConfig.primaryColor,
  set: (val) => store.updateStyle('primaryColor', val)
})

const accentColor = computed({
  get: () => store.styleConfig.accentColor,
  set: (val) => store.updateStyle('accentColor', val)
})

const fontSize = computed({
  get: () => store.styleConfig.fontSize,
  set: (val) => store.updateStyle('fontSize', val)
})

const lineHeight = computed({
  get: () => store.styleConfig.lineSpacing,
  set: (val) => store.updateStyle('lineSpacing', val)
})

const marginTop = computed({
  get: () => store.styleConfig.marginTop,
  set: (val) => store.updateStyle('marginTop', val)
})

const marginBottom = computed({
  get: () => store.styleConfig.marginBottom,
  set: (val) => store.updateStyle('marginBottom', val)
})

const marginLeft = computed({
  get: () => store.styleConfig.marginLeft,
  set: (val) => store.updateStyle('marginLeft', val)
})

const marginRight = computed({
  get: () => store.styleConfig.marginRight,
  set: (val) => store.updateStyle('marginRight', val)
})

function updateStyle(key, value) {
  store.updateStyle(key, value)
}
</script>

<style lang="scss" scoped>
.style-panel {
  border-top: 1px solid #e4e7ed;

  .style-form {
    padding: 12px 0;

    .el-form-item {
      margin-bottom: 20px;
    }
  }

  :deep(.el-collapse-item__header) {
    font-weight: 600;
  }
}
</style>
