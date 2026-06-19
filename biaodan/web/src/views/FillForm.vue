<template>
  <div class="fill-container">
    <div v-loading="loading" class="fill-inner">
      <template v-if="template">
        <div class="fill-card">
          <header class="fill-header">
            <h1 class="fill-title">{{ template.name }}</h1>
            <p v-if="template.description" class="fill-desc">{{ template.description }}</p>
          </header>

          <el-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            label-position="top"
            class="fill-form"
          >
            <form-renderer
              v-for="field in template.schema.fields"
              :key="field.id"
              :field="field"
              v-model="formData"
            />

            <div class="form-actions">
              <el-button
                type="primary"
                size="large"
                :loading="submitting"
                style="width: 160px"
                @click="submitForm"
              >
                提交表单
              </el-button>
            </div>
          </el-form>
        </div>

        <div v-if="success" class="success-overlay">
          <div class="success-card">
            <el-icon :size="64" color="#67C23A"><CircleCheckFilled /></el-icon>
            <h2>提交成功！</h2>
            <p>感谢您的填写，我们已收到您的信息。</p>
            <el-button type="primary" @click="resetForm">再填一份</el-button>
          </div>
        </div>
      </template>

      <el-empty v-else-if="!loading" description="表单不存在或已被删除" />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CircleCheckFilled } from '@element-plus/icons-vue'
import { getTemplate, submitForm as submitFormApi } from '@/api'
import FormRenderer from '@/components/FormRenderer.vue'

const route = useRoute()
const loading = ref(false)
const submitting = ref(false)
const template = ref(null)
const formData = reactive({})
const formRef = ref(null)
const success = ref(false)

const formRules = computed(() => {
  const rules = {}
  if (!template.value) return rules

  template.value.schema.fields.forEach(field => {
    if (field.required) {
      rules[field.id] = [{
        required: true,
        validator: (rule, value, callback) => {
          const v = formData[field.id]
          if (v === undefined || v === null || v === '') {
            callback(new Error(field.label + '不能为空'))
            return
          }
          if (Array.isArray(v) && v.length === 0) {
            callback(new Error(field.label + '不能为空'))
            return
          }
          if (field.type === 'email' && typeof v === 'string' && v) {
            if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)) {
              callback(new Error('邮箱格式不正确'))
              return
            }
          }
          if (field.type === 'phone' && typeof v === 'string' && v) {
            if (!/^1[3-9]\d{9}$/.test(v)) {
              callback(new Error('手机号格式不正确'))
              return
            }
          }
          callback()
        },
        trigger: ['blur', 'change']
      }]
    }
  })

  return rules
})

async function loadTemplate() {
  const id = route.params.id
  if (!id) return
  loading.value = true
  try {
    const res = await getTemplate(id)
    template.value = res.data
    initFormData()
  } finally {
    loading.value = false
  }
}

function initFormData() {
  if (!template.value) return
  template.value.schema.fields.forEach(field => {
    if (field.type === 'checkbox') {
      formData[field.id] = []
    } else if (field.type === 'switch') {
      formData[field.id] = false
    } else if (field.type === 'rate' || field.type === 'slider') {
      formData[field.id] = 0
    } else {
      formData[field.id] = ''
    }
  })
}

function clearFormData() {
  Object.keys(formData).forEach(key => delete formData[key])
  initFormData()
}

async function submitForm() {
  try {
    await formRef.value.validate()
  } catch (e) {
    ElMessage.warning('请完整填写表单')
    return
  }

  submitting.value = true
  try {
    await submitFormApi({
      templateId: template.value.id,
      data: { ...formData }
    })
    success.value = true
  } finally {
    submitting.value = false
  }
}

function resetForm() {
  success.value = false
  clearFormData()
  formRef.value?.clearValidate()
}

onMounted(loadTemplate)
</script>

<style scoped>
.fill-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px 16px;
}

.fill-inner {
  max-width: 680px;
  margin: 0 auto;
  position: relative;
}

.fill-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.fill-header {
  padding: 36px 36px 20px;
  background: linear-gradient(135deg, #409EFF 0%, #337ecc 100%);
  color: #fff;
}

.fill-title {
  font-size: 26px;
  font-weight: 600;
  margin: 0 0 8px;
}

.fill-desc {
  margin: 0;
  font-size: 14px;
  opacity: 0.9;
  line-height: 1.6;
}

.fill-form {
  padding: 28px 36px 36px;
}

.fill-form :deep(.field-item) {
  margin-bottom: 8px;
  padding: 12px 0;
}

.fill-form :deep(.field-label) {
  font-size: 15px;
  margin-bottom: 10px;
}

.form-actions {
  text-align: center;
  padding-top: 24px;
  margin-top: 8px;
  border-top: 1px solid #f0f2f5;
}

.success-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.success-card {
  background: #fff;
  border-radius: 12px;
  padding: 48px 40px;
  text-align: center;
  max-width: 400px;
  width: 90%;
}

.success-card h2 {
  margin: 20px 0 8px;
  color: #303133;
  font-size: 22px;
}

.success-card p {
  color: #909399;
  margin: 0 0 28px;
  line-height: 1.6;
}

@media (max-width: 640px) {
  .fill-container {
    padding: 0;
    background: #f5f7fa;
  }
  .fill-card {
    border-radius: 0;
    min-height: 100vh;
  }
  .fill-header {
    padding: 28px 20px 16px;
  }
  .fill-title {
    font-size: 22px;
  }
  .fill-form {
    padding: 20px;
  }
}
</style>
