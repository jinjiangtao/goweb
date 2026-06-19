<template>
  <div class="property-config">
    <div class="prop-group">
      <div class="prop-title">基础属性</div>

      <el-form :model="localField" label-position="top" size="default">
        <el-form-item label="组件标签 (Label)">
          <el-input v-model="localField.label" placeholder="请输入标签名称" />
        </el-form-item>

        <el-form-item label="组件ID">
          <el-input v-model="localField.id" disabled />
        </el-form-item>

        <el-form-item v-if="hasPlaceholder" label="占位提示">
          <el-input v-model="localField.placeholder" placeholder="请输入占位文字" />
        </el-form-item>

        <el-form-item v-if="hasDefaultValue && !isOptionsType" label="默认值">
          <el-input v-if="isTextType" v-model="localField.defaultValue" placeholder="默认值" />
          <el-input-number v-else-if="field.type === 'number'" v-model="localField.defaultValue" controls-position="right" style="width: 100%" />
          <el-switch v-else-if="field.type === 'switch'" v-model="localField.defaultValue" />
          <el-rate v-else-if="field.type === 'rate'" v-model="localField.defaultValue" />
          <el-slider v-else-if="field.type === 'slider'" v-model="localField.defaultValue" :show-input="true" />
        </el-form-item>

        <el-form-item v-if="!isLayoutType" label="是否必填">
          <el-switch v-model="localField.required" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="hasOptions" class="prop-group">
      <div class="prop-title">选项配置</div>
      <div class="options-editor">
        <div
          v-for="(opt, idx) in localField.options"
          :key="idx"
          class="option-row"
        >
          <el-input
            v-model="opt.label"
            placeholder="标签"
            size="small"
            style="flex: 1"
          />
          <el-input
            v-model="opt.value"
            placeholder="值"
            size="small"
            style="flex: 1; margin-left: 8px"
          />
          <el-button
            size="small"
            :icon="Delete"
            text
            type="danger"
            style="margin-left: 4px"
            @click="removeOption(idx)"
          />
        </div>
        <el-button size="small" :icon="Plus" style="width: 100%; margin-top: 8px" @click="addOption">
          添加选项
        </el-button>
      </div>
    </div>

    <div v-if="field.type === 'input' || field.type === 'textarea'" class="prop-group">
      <div class="prop-title">文本属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="最大长度">
          <el-input-number v-model="localField.maxLength" :min="0" :max="10000" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="field.type === 'textarea'" label="行数">
          <el-input-number v-model="localField.rows" :min="1" :max="20" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="field.type === 'input'" label="输入类型">
          <el-select v-model="localField.inputType" style="width: 100%">
            <el-option label="普通文本" value="text" />
            <el-option label="密码" value="password" />
            <el-option label="URL" value="url" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'number'" class="prop-group">
      <div class="prop-title">数字属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="最小值">
          <el-input-number v-model="localField.min" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="最大值">
          <el-input-number v-model="localField.max" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="步长">
          <el-input-number v-model="localField.step" :min="0.1" controls-position="right" style="width: 100%" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'date'" class="prop-group">
      <div class="prop-title">日期属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="日期类型">
          <el-select v-model="localField.dateType" style="width: 100%">
            <el-option label="日期" value="date" />
            <el-option label="日期时间" value="datetime" />
            <el-option label="周" value="week" />
            <el-option label="月" value="month" />
            <el-option label="年" value="year" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否范围选择">
          <el-switch v-model="localField.range" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'upload'" class="prop-group">
      <div class="prop-title">上传属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="允许多文件">
          <el-switch v-model="localField.multiple" />
        </el-form-item>
        <el-form-item label="最大大小 (MB)">
          <el-input-number v-model="localField.maxSize" :min="1" :max="500" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="允许类型">
          <el-input v-model="localField.accept" placeholder="如: .jpg,.png,.pdf 留空表示不限" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'rate'" class="prop-group">
      <div class="prop-title">评分属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="最大分值">
          <el-input-number v-model="localField.max" :min="1" :max="10" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="允许半星">
          <el-switch v-model="localField.allowHalf" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'switch'" class="prop-group">
      <div class="prop-title">开关属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="开启文字">
          <el-input v-model="localField.activeText" />
        </el-form-item>
        <el-form-item label="关闭文字">
          <el-input v-model="localField.inactiveText" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'slider'" class="prop-group">
      <div class="prop-title">滑块属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="最小值">
          <el-input-number v-model="localField.min" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="最大值">
          <el-input-number v-model="localField.max" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="步长">
          <el-input-number v-model="localField.step" :min="1" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="显示输入框">
          <el-switch v-model="localField.showInput" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'divider'" class="prop-group">
      <div class="prop-title">分割线属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="中间文字">
          <el-input v-model="localField.content" placeholder="可选，留空则无文字" />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="field.type === 'description'" class="prop-group">
      <div class="prop-title">说明文字属性</div>
      <el-form label-position="top" size="default">
        <el-form-item label="说明内容">
          <el-input v-model="localField.content" type="textarea" :rows="4" placeholder="请输入说明内容" />
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, watch } from 'vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import { useFormStore } from '@/stores/form'

const props = defineProps({
  field: {
    type: Object,
    required: true
  }
})

const formStore = useFormStore()
const localField = reactive({ ...props.field })

watch(() => props.field, (newVal) => {
  Object.keys(localField).forEach(key => delete localField[key])
  Object.assign(localField, newVal)
}, { deep: true })

watch(localField, (newVal) => {
  formStore.updateField(props.field.id, { ...newVal })
}, { deep: true })

const hasPlaceholder = computed(() => {
  return ['input', 'textarea', 'number', 'email', 'phone', 'select', 'date', 'time'].includes(props.field.type)
})

const hasDefaultValue = computed(() => {
  return !['divider', 'description', 'upload'].includes(props.field.type)
})

const isTextType = computed(() => {
  return ['input', 'textarea', 'email', 'phone'].includes(props.field.type)
})

const isOptionsType = computed(() => {
  return ['radio', 'checkbox', 'select'].includes(props.field.type)
})

const hasOptions = computed(() => isOptionsType.value)

const isLayoutType = computed(() => {
  return ['divider', 'description'].includes(props.field.type)
})

function addOption() {
  const idx = localField.options.length + 1
  localField.options.push({ label: `选项${idx}`, value: `opt${idx}_${Math.random().toString(36).substr(2, 5)}` })
}

function removeOption(idx) {
  if (localField.options.length <= 1) return
  localField.options.splice(idx, 1)
}
</script>

<style scoped>
.property-config {
  padding: 4px;
}

.prop-group {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f2f5;
}

.prop-group:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.prop-title {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 14px;
  padding-left: 8px;
  border-left: 3px solid #409EFF;
}

.options-editor {
  background: #fafbfc;
  border-radius: 6px;
  padding: 12px;
}

.option-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.option-row:last-of-type {
  margin-bottom: 0;
}
</style>
