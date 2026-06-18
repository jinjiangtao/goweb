<template>
  <div class="dynamic-form">
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item
        v-for="field in fields"
        :key="field.name"
        :label="field.label"
        :prop="field.name"
      >
        <el-input
          v-if="field.type === 'input'"
          v-model="formData[field.name]"
          :placeholder="field.placeholder || '请输入' + field.label"
          clearable
        />

        <el-input
          v-else-if="field.type === 'textarea'"
          v-model="formData[field.name]"
          type="textarea"
          :rows="4"
          :placeholder="field.placeholder || '请输入' + field.label"
          clearable
        />

        <el-input-number
          v-else-if="field.type === 'number'"
          v-model="formData[field.name]"
          :min="field.validation?.min || 0"
          :max="field.validation?.max"
          :step="field.validation?.step || 1"
          :precision="field.validation?.precision"
          style="width: 100%"
        />

        <el-select
          v-else-if="field.type === 'select'"
          v-model="formData[field.name]"
          :placeholder="field.placeholder || '请选择' + field.label"
          clearable
          style="width: 100%"
        >
          <el-option
            v-for="option in field.options"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>

        <el-date-picker
          v-else-if="field.type === 'date'"
          v-model="formData[field.name]"
          type="date"
          :placeholder="field.placeholder || '请选择' + field.label"
          value-format="YYYY-MM-DD"
          style="width: 100%"
        />

        <el-date-picker
          v-else-if="field.type === 'datetime'"
          v-model="formData[field.name]"
          type="datetime"
          :placeholder="field.placeholder || '请选择' + field.label"
          value-format="YYYY-MM-DD HH:mm:ss"
          style="width: 100%"
        />

        <el-switch
          v-else-if="field.type === 'switch'"
          v-model="formData[field.name]"
        />

        <el-radio-group
          v-else-if="field.type === 'radio'"
          v-model="formData[field.name]"
        >
          <el-radio
            v-for="option in field.options"
            :key="option.value"
            :label="option.value"
          >{{ option.label }}</el-radio>
        </el-radio-group>

        <el-checkbox-group
          v-else-if="field.type === 'checkbox'"
          v-model="formData[field.name]"
        >
          <el-checkbox
            v-for="option in field.options"
            :key="option.value"
            :label="option.value"
          >{{ option.label }}</el-checkbox>
        </el-checkbox-group>

        <el-input
          v-else
          v-model="formData[field.name]"
          :placeholder="field.placeholder || '请输入' + field.label"
          clearable
        />
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  fields: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Object,
    default: () => ({})
  },
  readonly: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'submit'])

const formRef = ref(null)

const formData = ref({ ...props.modelValue })

watch(
  () => props.modelValue,
  (newVal) => {
    formData.value = { ...newVal }
  },
  { deep: true }
)

watch(
  formData,
  (newVal) => {
    emit('update:modelValue', { ...newVal })
  },
  { deep: true }
)

const rules = computed(() => {
  const result = {}
  props.fields.forEach(field => {
    if (field.required) {
      result[field.name] = [
        { required: true, message: `请${field.type === 'select' ? '选择' : '输入'}${field.label}`, trigger: ['blur', 'change'] }
      ]
    }
    if (field.validation) {
      if (!result[field.name]) {
        result[field.name] = []
      }
      if (field.validation.min !== undefined && field.type !== 'number') {
        result[field.name].push({
          min: field.validation.min,
          message: `最少${field.validation.min}个字符`,
          trigger: 'blur'
        })
      }
      if (field.validation.max !== undefined && field.type !== 'number') {
        result[field.name].push({
          max: field.validation.max,
          message: `最多${field.validation.max}个字符`,
          trigger: 'blur'
        })
      }
      if (field.validation.pattern) {
        result[field.name].push({
          pattern: new RegExp(field.validation.pattern),
          message: field.validation.message || '格式不正确',
          trigger: 'blur'
        })
      }
    }
  })
  return result
})

const validate = () => {
  return formRef.value?.validate()
}

const resetFields = () => {
  formRef.value?.resetFields()
}

defineExpose({
  validate,
  resetFields,
  formData
})
</script>

<style scoped lang="scss">
.dynamic-form {
  :deep(.el-form-item) {
    margin-bottom: 22px;
  }
}
</style>
