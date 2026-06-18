<template>
  <div class="todo-page">
    <div class="page-header">
      <h2>待我审批</h2>
      <div class="header-actions">
        <el-select
          v-model="filterStatus"
          placeholder="状态筛选"
          style="width: 150px; margin-right: 15px;"
          clearable
          @change="loadApprovals"
        >
          <el-option label="审批中" value="pending" />
          <el-option label="已通过" value="approved" />
          <el-option label="已驳回" value="rejected" />
        </el-select>
      </div>
    </div>

    <el-table :data="approvalList" border stripe v-loading="loading">
      <el-table-column prop="request_no" label="申请编号" width="200" />
      <el-table-column prop="title" label="申请标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)">{{ getTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="applicant.name" label="申请人" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="current_node_id" label="当前节点" width="120" />
      <el-table-column prop="created_at" label="申请时间" width="180">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="goToDetail(row.id)">审批</el-button>
          <el-button
            v-if="row.status === 'pending'"
            size="small"
            type="success"
            @click="handleQuickApprove(row)"
          >同意</el-button>
          <el-button
            v-if="row.status === 'pending'"
            size="small"
            type="danger"
            @click="handleQuickReject(row)"
          >驳回</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="loadApprovals"
      @current-change="loadApprovals"
      style="margin-top: 20px; justify-content: flex-end;"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMyApprovals, approveApproval, rejectApproval } from '@/api/approval'

const router = useRouter()

const loading = ref(false)
const approvalList = ref([])
const filterStatus = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadApprovals()
})

const loadApprovals = async () => {
  loading.value = true
  try {
    const res = await getMyApprovals({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: filterStatus.value
    })
    approvalList.value = res.data.list
    pagination.total = res.data.total
  } catch (e) {
    ElMessage.error('加载审批列表失败')
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

const getStatusType = (status) => {
  switch (status) {
    case 'approved':
      return 'success'
    case 'rejected':
      return 'danger'
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
    case 'pending':
      return '审批中'
    default:
      return status
  }
}

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const goToDetail = (id) => {
  router.push(`/approval/detail/${id}`)
}

const handleQuickApprove = async (row) => {
  try {
    const { value: comment } = await ElMessageBox.prompt('请输入审批意见（可选）', '同意申请', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '请输入审批意见',
      type: 'success'
    })
    await approveApproval(row.id, { comment: comment || '同意' })
    ElMessage.success('审批通过')
    loadApprovals()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '审批失败')
    }
  }
}

const handleQuickReject = async (row) => {
  try {
    const { value: comment } = await ElMessageBox.prompt('请输入驳回原因', '驳回申请', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '请输入驳回原因',
      inputValidator: (value) => {
        if (!value || value.trim() === '') {
          return '请输入驳回原因'
        }
        return true
      },
      type: 'warning'
    })
    await rejectApproval(row.id, { comment })
    ElMessage.success('已驳回')
    loadApprovals()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '驳回失败')
    }
  }
}
</script>

<style scoped lang="scss">
.todo-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 20px;
      color: #303133;
    }
  }
}
</style>
