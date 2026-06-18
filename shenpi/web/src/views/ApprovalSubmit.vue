<template>
  <div class="submit-page">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>提交审批申请</h2>
    </div>

    <div v-if="template" class="submit-content">
      <el-card class="template-info">
        <template #header>
          <div class="card-header">
            <span>模板信息</span>
            <el-tag :type="getTypeTagType(template.type)">{{ getTypeText(template.type) }}</el-tag>
          </div>
        </template>
        <div class="info-item">
          <span class="label">模板名称：</span>
          <span class="value">{{ template.name }}</span>
        </div>
        <div class="info-item">
          <span class="label">模板描述：</span>
          <span class="value">{{ template.description || '暂无描述' }}</span>
        </div>
        <div class="info-item">
          <span class="label">创建人：</span>
          <span class="value">{{ template.creator?.name }}</span>
        </div>
      </el-card>

      <el-card class="form-card">
        <template #header>
          <span>申请内容</span>
        </template>
        <el-form :model="formData" label-width="120px">
          <el-form-item label="申请标题" required>
            <el-input
              v-model="formData.title"
              placeholder="请输入申请标题"
              clearable
            />
          </el-form-item>
        </el-form>
        <DynamicForm
          ref="dynamicFormRef"
          :fields="formFields"
          v-model="formData.fields"
        />
        <div class="form-actions">
          <el-button @click="goBack">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">提交申请</el-button>
        </div>
      </el-card>
    </div>

    <el-empty v-else description="加载中..." v-loading="loading" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getTemplateDetail, getTemplateFields } from '@/api/template'
import { submitApproval } from '@/api/approval'
import DynamicForm from '@/components/DynamicForm.vue'

const route = useRoute()
const router = useRouter()

const templateId = route.params.templateId
const loading = ref(false)
const submitting = ref(false)
const template = ref(null)
const formFields = ref([])
const dynamicFormRef = ref(null)

const formData = reactive({
  title: '',
  fields: {}
})

onMounted(() => {
  loadTemplate()
})

const loadTemplate = async () => {
  loading.value = true
  try {
    const [templateRes, fieldsRes] = await Promise.all([
      getTemplateDetail(templateId),
      getTemplateFields(templateId)
    ])
    template.value = templateRes.data
    formFields.value = fieldsRes.data

    formFields.value.forEach(field => {
      if (field.default !== undefined) {
        formData.fields[field.name] = field.default
      }
    })
  } catch (e) {
    ElMessage.error('加载模板信息失败')
  } finally {
    loading.value = false
  }
}

const getTypeTagType = (type) => {
  switch (type) {
    case 'leave':
      return 'primary'
    case 'reimburse':
      return 'success'
    default:
      return 'info'
  }
}

const getTypeText = (type) => {
  switch (type) {
    case 'leave':
      return '请假审批'
    case 'reimburse':
      return '报销审批'
    default:
      return '通用申请'
  }
}

const handleSubmit = async () => {
  if (!formData.title.trim()) {
    ElMessage.warning('请输入申请标题')
    return
  }

  try {
    await dynamicFormRef.value.validate()
  } catch (e) {
    ElMessage.warning('请完善表单信息')
    return
  }

  submitting.value = true
  try {
    const fullFormData = {
      title: formData.title,
      ...formData.fields
    }

    await submitApproval({
      template_id: templateId,
      title: formData.title,
      form_data: fullFormData
    })

    ElMessage.success('申请提交成功')
    router.push('/approval')
  } catch (e) {
    ElMessage.error(e.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

const goBack = () => {
  router.back()
}
</script>

<style scoped lang="scss">
.submit-page {
  .page-header {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 20px;
      color: #303133;
    }
  }

  .submit-content {
    max-width: 800px;
    margin: 0 auto;

    .template-info {
      margin-bottom: 20px;

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .info-item {
        margin-bottom: 10px;
        font-size: 14px;

        .label {
          color: #909399;
        }

        .value {
          color: #303133;
        }
      }
    }

    .form-card {
      .form-actions {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 20px;
        padding-top: 20px;
        border-top: 1px solid #ebeef5;
      }
    }
  }
}
</style>
