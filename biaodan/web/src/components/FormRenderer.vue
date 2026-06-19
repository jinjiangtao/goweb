<template>
  <div class="field-item" :class="`field-${field.type}`">
    <template v-if="field.type === 'divider'">
      <div class="divider-field">
        <span v-if="field.content" class="divider-text">{{ field.content }}</span>
        <hr />
      </div>
    </template>

    <template v-else-if="field.type === 'description'">
      <div class="description-field">
        <i class="desc-icon"><el-icon :size="16"><InfoFilled /></el-icon></i>
        <p>{{ field.content || '说明文字' }}</p>
      </div>
    </template>

    <template v-else>
      <label v-if="field.label" class="field-label" :for="field.id">
        <span class="required" v-if="field.required">*</span>
        {{ field.label }}
      </label>

      <div class="field-control">
        <template v-if="field.type === 'input' || field.type === 'email' || field.type === 'phone'">
          <el-input
            :id="field.id"
            v-model="model[field.id]"
            :type="field.type === 'email' ? 'email' : (field.type === 'phone' ? 'tel' : (field.inputType || 'text'))"
            :placeholder="field.placeholder || '请输入' + field.label"
            :maxlength="field.maxLength"
            :show-word-limit="field.maxLength > 0 && !preview"
            clearable
            :disabled="preview"
          />
        </template>

        <template v-else-if="field.type === 'textarea'">
          <el-input
            :id="field.id"
            v-model="model[field.id]"
            type="textarea"
            :rows="field.rows || 4"
            :placeholder="field.placeholder || '请输入' + field.label"
            :maxlength="field.maxLength"
            :show-word-limit="field.maxLength > 0 && !preview"
            :disabled="preview"
          />
        </template>

        <template v-else-if="field.type === 'number'">
          <el-input-number
            :id="field.id"
            v-model="model[field.id]"
            :min="field.min"
            :max="field.max"
            :step="field.step || 1"
            :placeholder="field.placeholder || '请输入数字'"
            controls-position="right"
            style="width: 100%"
            :disabled="preview"
          />
        </template>

        <template v-else-if="field.type === 'radio'">
          <el-radio-group
            :id="field.id"
            v-model="model[field.id]"
            :disabled="preview"
          >
            <el-radio
              v-for="opt in field.options"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </el-radio>
          </el-radio-group>
        </template>

        <template v-else-if="field.type === 'checkbox'">
          <el-checkbox-group
            :id="field.id"
            v-model="model[field.id]"
            :disabled="preview"
          >
            <el-checkbox
              v-for="opt in field.options"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </el-checkbox>
          </el-checkbox-group>
        </template>

        <template v-else-if="field.type === 'select'">
          <el-select
            :id="field.id"
            v-model="model[field.id]"
            :placeholder="field.placeholder || '请选择' + field.label"
            clearable
            style="width: 100%"
            :disabled="preview"
          >
            <el-option
              v-for="opt in field.options"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
        </template>

        <template v-else-if="field.type === 'date'">
          <template v-if="field.range">
            <el-date-picker
              :id="field.id"
              v-model="model[field.id]"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 100%"
              :disabled="preview"
              value-format="YYYY-MM-DD"
            />
          </template>
          <template v-else>
            <el-date-picker
              :id="field.id"
              v-model="model[field.id]"
              :type="field.dateType || 'date'"
              :placeholder="field.placeholder || '请选择日期'"
              style="width: 100%"
              :disabled="preview"
              value-format="YYYY-MM-DD"
            />
          </template>
        </template>

        <template v-else-if="field.type === 'time'">
          <el-time-picker
            :id="field.id"
            v-model="model[field.id]"
            :placeholder="field.placeholder || '请选择时间'"
            style="width: 100%"
            :disabled="preview"
            format="HH:mm"
            value-format="HH:mm"
          />
        </template>

        <template v-else-if="field.type === 'upload'">
          <el-upload
            :id="field.id"
            v-model:file-list="fileListMap[field.id]"
            :action="'/api/upload'"
            :multiple="field.multiple"
            :limit="field.multiple ? 10 : 1"
            :accept="field.accept"
            :before-upload="(file) => beforeUpload(file, field)"
            :on-success="(res, file) => onUploadSuccess(res, file, field)"
            :on-remove="(file) => onUploadRemove(file, field)"
            :disabled="preview"
            :auto-upload="true"
            list-type="text"
          >
            <el-button :icon="UploadFilled" :disabled="preview">
              点击上传
            </el-button>
            <template #tip>
              <div class="upload-tip">
                {{ field.maxSize ? `单文件不超过 ${field.maxSize}MB` : '' }}
                {{ field.accept ? `支持格式: ${field.accept}` : '' }}
              </div>
            </template>
          </el-upload>
        </template>

        <template v-else-if="field.type === 'rate'">
          <el-rate
            :id="field.id"
            v-model="model[field.id]"
            :max="field.max || 5"
            :allow-half="field.allowHalf"
            :disabled="preview"
          />
        </template>

        <template v-else-if="field.type === 'switch'">
          <el-switch
            :id="field.id"
            v-model="model[field.id]"
            :active-text="field.activeText || '是'"
            :inactive-text="field.inactiveText || '否'"
            :disabled="preview"
          />
        </template>

        <template v-else-if="field.type === 'slider'">
          <el-slider
            :id="field.id"
            v-model="model[field.id]"
            :min="field.min || 0"
            :max="field.max || 100"
            :step="field.step || 1"
            :show-input="field.showInput"
            :disabled="preview"
          />
        </template>

        <template v-else>
          <el-input
            v-model="model[field.id]"
            placeholder="未知组件类型"
            disabled
          />
        </template>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled, InfoFilled } from '@element-plus/icons-vue'

const props = defineProps({
  field: {
    type: Object,
    required: true
  },
  modelValue: {
    type: Object,
    default: () => ({})
  },
  preview: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const model = props.modelValue
const fileListMap = ref({})

function initDefaultValue() {
  if (props.field.type === 'upload') {
    if (!fileListMap.value[props.field.id]) {
      fileListMap.value[props.field.id] = []
    }
  }
  if (model[props.field.id] === undefined || model[props.field.id] === null) {
    if (props.field.type === 'checkbox') {
      model[props.field.id] = props.field.defaultValue || []
    } else if (props.field.type === 'switch') {
      model[props.field.id] = props.field.defaultValue || false
    } else if (props.field.type === 'number') {
      model[props.field.id] = props.field.defaultValue !== '' ? props.field.defaultValue : undefined
    } else if (props.field.type === 'rate' || props.field.type === 'slider') {
      model[props.field.id] = props.field.defaultValue || 0
    } else {
      model[props.field.id] = props.field.defaultValue || ''
    }
  }
  emit('update:modelValue', model)
}

function beforeUpload(file, field) {
  const maxSize = (field.maxSize || 50) * 1024 * 1024
  if (file.size > maxSize) {
    ElMessage.error(`文件大小不能超过 ${field.maxSize || 50}MB`)
    return false
  }
  return true
}

function onUploadSuccess(res, file, field) {
  if (res && res.data) {
    const url = res.data.url
    if (field.multiple) {
      if (!Array.isArray(model[field.id])) {
        model[field.id] = []
      }
      model[field.id].push({
        url,
        name: file.name,
        size: file.size
      })
    } else {
      model[field.id] = {
        url,
        name: file.name,
        size: file.size
      }
    }
    ElMessage.success('上传成功')
  }
}

function onUploadRemove(file, field) {
  if (field.multiple && Array.isArray(model[field.id])) {
    const idx = model[field.id].findIndex(f => f.name === file.name)
    if (idx > -1) model[field.id].splice(idx, 1)
  } else {
    model[field.id] = ''
  }
}

onMounted(initDefaultValue)
watch(() => props.field.id, initDefaultValue)
</script>

<style scoped>
.field-item {
  margin-bottom: 4px;
  padding: 8px 0;
}

.field-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  line-height: 1.4;
}

.required {
  color: #f56c6c;
  margin-right: 4px;
  font-weight: bold;
}

.field-control {
  width: 100%;
}

.divider-field {
  padding: 12px 0;
  position: relative;
  text-align: center;
}

.divider-field hr {
  border: none;
  border-top: 1px solid #ebeef5;
  margin: 0;
}

.divider-text {
  display: inline-block;
  padding: 0 16px;
  background: #fff;
  color: #909399;
  font-size: 13px;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.field-wrapper.selected .divider-text {
  background: #f5f9ff;
}

.description-field {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 12px 16px;
  background: #f4f4f5;
  border-radius: 6px;
  border-left: 3px solid #909399;
}

.desc-icon {
  color: #909399;
  flex-shrink: 0;
  margin-top: 2px;
}

.description-field p {
  margin: 0;
  color: #606266;
  font-size: 13px;
  line-height: 1.6;
}

.upload-tip {
  color: #909399;
  font-size: 12px;
  margin-top: 6px;
}
</style>
