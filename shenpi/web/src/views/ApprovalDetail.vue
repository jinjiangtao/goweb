<template>
  <div class="detail-page">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>审批详情</h2>
    </div>

    <div v-if="detail" class="detail-content">
      <el-card class="base-info">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-tag :type="getStatusType(detail.request.status)">{{ getStatusText(detail.request.status) }}</el-tag>
          </div>
        </template>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">申请编号：</span>
            <span class="value">{{ detail.request.request_no }}</span>
          </div>
          <div class="info-item">
            <span class="label">申请标题：</span>
            <span class="value">{{ detail.request.title }}</span>
          </div>
          <div class="info-item">
            <span class="label">申请类型：</span>
            <span class="value">{{ getTypeText(detail.request.type) }}</span>
          </div>
          <div class="info-item">
            <span class="label">申请人：</span>
            <span class="value">{{ detail.request.applicant?.name }}</span>
          </div>
          <div class="info-item">
            <span class="label">申请时间：</span>
            <span class="value">{{ formatTime(detail.request.created_at) }}</span>
          </div>
          <div class="info-item">
            <span class="label">当前节点：</span>
            <span class="value">{{ detail.request.current_node_id }}</span>
          </div>
        </div>
      </el-card>

      <el-card class="form-data">
        <template #header>
          <span>申请内容</span>
        </template>
        <div class="form-data-list">
          <div v-for="(value, key) in detail.form_data" :key="key" class="form-item">
            <span class="label">{{ getFieldLabel(key) }}：</span>
            <span class="value">{{ formatValue(value) }}</span>
          </div>
        </div>
      </el-card>

      <ApprovalProgress
        :approval-id="route.params.id"
        :show-actions="canOperate"
        @approved="onApproved"
        @rejected="onRejected"
        @revoked="onRevoked"
      />
    </div>

    <el-empty v-else description="加载中..." v-loading="loading" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getApprovalDetail } from '@/api/approval'
import { getTemplateFields } from '@/api/template'
import { useUserStore } from '@/store/user'
import ApprovalProgress from '@/components/ApprovalProgress.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const detail = ref(null)
const formFields = ref([])

const canOperate = computed(() => {
  if (!detail.value) return false
  return detail.value.request.status === 'pending'
})

onMounted(() => {
  loadDetail()
})

const loadDetail = async () => {
  loading.value = true
  try {
    const res = await getApprovalDetail(route.params.id)
    detail.value = res.data

    if (res.data.request.template_id) {
      const fieldsRes = await getTemplateFields(res.data.request.template_id)
      formFields.value = fieldsRes.data
    }
  } catch (e) {
    ElMessage.error('加载详情失败')
  } finally {
    loading.value = false
  }
}

const getFieldLabel = (name) => {
  const field = formFields.value.find(f => f.name === name)
  return field ? field.label : name
}

const formatValue = (value) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }
  return String(value)
}

const getStatusType = (status) => {
  switch (status) {
    case 'approved':
      return 'success'
    case 'rejected':
      return 'danger'
    case 'revoked':
      return 'info'
    default:
      return 'warning'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'approved':
      return '已通过'
    case 'rejected':
      return '已驳回'
    case 'revoked':
      return '已撤回'
    case 'pending':
      return '审批中'
    default:
      return status
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

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const onApproved = () => {
  loadDetail()
}

const onRejected = () => {
  loadDetail()
}

const onRevoked = () => {
  loadDetail()
}

const goBack = () => {
  router.back()
}
</script>

<style scoped lang="scss">
.detail-page {
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

  .detail-content {
    max-width: 1000px;
    margin: 0 auto;

    .base-info {
      margin-bottom: 20px;

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .info-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 15px;

        .info-item {
          font-size: 14px;

          .label {
            color: #909399;
          }

          .value {
            color: #303133;
          }
        }
      }
    }

    .form-data {
      margin-bottom: 20px;

      .form-data-list {
        .form-item {
          display: flex;
          margin-bottom: 12px;
          font-size: 14px;

          .label {
            width: 120px;
            color: #909399;
            flex-shrink: 0;
          }

          .value {
            color: #303133;
            word-break: break-all;
          }
        }
      }
    }
  }
}
</style>
