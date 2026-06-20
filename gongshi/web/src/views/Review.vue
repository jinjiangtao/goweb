<template>
  <div class="page-container">
    <div class="card">
      <div class="card-header">
        <div class="card-title">
          <el-icon style="vertical-align: middle; margin-right: 6px; color: #409eff"><Check /></el-icon>
          工时审核
          <el-tag type="warning" style="margin-left: 10px" v-if="pendingCount > 0">待审核 {{ pendingCount }} 条</el-tag>
        </div>
        <div class="flex-row" style="gap: 8px">
          <el-button type="success" @click="handleBatchApprove" :disabled="selectedIds.length === 0" :loading="batchLoading">
            <el-icon><Select /></el-icon>批量通过 ({{ selectedIds.length }})
          </el-button>
        </div>
      </div>

      <el-form :inline="true" :model="filterForm" style="margin-bottom: 16px">
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="filterForm.dates"
            type="daterange"
            value-format="YYYY-MM-DD"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          />
        </el-form-item>
        <el-form-item label="员工">
          <el-select v-model="filterForm.user_id" placeholder="全部员工" clearable style="width: 160px" filterable>
            <el-option v-for="u in users" :key="u.id" :label="u.name" :value="u.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="项目">
          <el-select v-model="filterForm.project_id" placeholder="全部项目" clearable style="width: 180px">
            <el-option v-for="p in projects" :key="p.id" :label="p.code + ' - ' + p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadRecords">
            <el-icon><Search /></el-icon>查询
          </el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="records" border @selection-change="handleSelectionChange" v-loading="loading">
        <el-table-column type="selection" width="55" :selectable="(row) => row.status === 'pending'" />
        <el-table-column prop="user.name" label="员工" width="100" />
        <el-table-column prop="work_date" label="日期" width="120" sortable />
        <el-table-column prop="project.name" label="项目" min-width="160">
          <template #default="{ row }">
            <el-tag size="small">{{ row.project?.code }}</el-tag> {{ row.project?.name }}
          </template>
        </el-table-column>
        <el-table-column prop="hours" label="工时" width="90" align="center">
          <template #default="{ row }">
            <span style="color: #409eff; font-weight: 600">{{ row.hours }}h</span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="工作内容" min-width="240" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'pending'" type="warning" size="small">待审核</el-tag>
            <el-tag v-else-if="row.status === 'approved'" type="success" size="small">已通过</el-tag>
            <el-tag v-else type="danger" size="small">已驳回</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="rejectReason" label="驳回原因" width="140" show-overflow-tooltip />
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">查看</el-button>
            <el-button v-if="row.status === 'pending'" type="success" link size="small" @click="handleApprove(row)">通过</el-button>
            <el-button v-if="row.status === 'pending'" type="danger" link size="small" @click="handleReject(row)">驳回</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="detailVisible" title="工时详情" width="560px">
      <el-descriptions :column="2" border v-if="currentRecord">
        <el-descriptions-item label="填报人">{{ currentRecord.user?.name }}</el-descriptions-item>
        <el-descriptions-item label="工作日期">{{ currentRecord.work_date }}</el-descriptions-item>
        <el-descriptions-item label="所属项目">{{ currentRecord.project?.code }} - {{ currentRecord.project?.name }}</el-descriptions-item>
        <el-descriptions-item label="工时">
          <span style="color: #409eff; font-weight: 600">{{ currentRecord.hours }} 小时</span>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag v-if="currentRecord.status === 'pending'" type="warning">待审核</el-tag>
          <el-tag v-else-if="currentRecord.status === 'approved'" type="success">已通过</el-tag>
          <el-tag v-else type="danger">已驳回</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="填报时间">{{ currentRecord.created_at }}</el-descriptions-item>
        <el-descriptions-item label="工作内容" :span="2">{{ currentRecord.content }}</el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.reject_reason" label="驳回原因" :span="2">
          <span style="color: #f56c6c">{{ currentRecord.reject_reason }}</span>
        </el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.reviewer" label="审核人">{{ currentRecord.reviewer?.name }}</el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.reviewed_at" label="审核时间">{{ currentRecord.reviewed_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer v-if="currentRecord && currentRecord.status === 'pending'">
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="success" @click="handleApprove(currentRecord)">通过</el-button>
        <el-button type="danger" @click="handleReject(currentRecord)">驳回</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="rejectVisible" title="驳回原因" width="440px">
      <el-form :model="rejectForm" :rules="rejectRules" ref="rejectFormRef" label-width="80px">
        <el-form-item label="驳回原因" prop="reason">
          <el-input v-model="rejectForm.reason" type="textarea" :rows="4" placeholder="请输入驳回原因..." maxlength="200" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmReject" :loading="rejectLoading">确认驳回</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { getProjects } from '../api/project'
import { getUsers } from '../api/auth'
import { getWorkRecords, approveWorkRecord, rejectWorkRecord, batchApprove } from '../api/workrecord'

const projects = ref([])
const users = ref([])
const records = ref([])
const loading = ref(false)
const batchLoading = ref(false)
const rejectLoading = ref(false)
const detailVisible = ref(false)
const rejectVisible = ref(false)
const currentRecord = ref(null)
const selectedIds = ref([])
const rejectFormRef = ref()

const filterForm = reactive({
  dates: [dayjs().subtract(1, 'month').format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')],
  user_id: null,
  project_id: null,
  status: ''
})

const rejectForm = reactive({
  reason: ''
})

const rejectRules = {
  reason: [{ required: true, message: '请输入驳回原因', trigger: 'blur', min: 2 }]
}

const pendingCount = computed(() => records.value.filter(r => r.status === 'pending').length)

const resetFilter = () => {
  filterForm.dates = [dayjs().subtract(1, 'month').format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')]
  filterForm.user_id = null
  filterForm.project_id = null
  filterForm.status = ''
  loadRecords()
}

const loadRecords = async () => {
  loading.value = true
  selectedIds.value = []
  try {
    const params = {
      start_date: filterForm.dates?.[0] || '',
      end_date: filterForm.dates?.[1] || ''
    }
    if (filterForm.user_id) params.user_id = filterForm.user_id
    if (filterForm.project_id) params.project_id = filterForm.project_id
    if (filterForm.status) params.status = filterForm.status
    records.value = await getWorkRecords(params)
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (selection) => {
  selectedIds.value = selection.filter(r => r.status === 'pending').map(r => r.id)
}

const viewDetail = (row) => {
  currentRecord.value = row
  detailVisible.value = true
}

const handleApprove = async (row) => {
  ElMessageBox.confirm(`确认通过 ${row.user?.name} 在 ${row.work_date} 的工时吗？`, '确认', { type: 'info' })
    .then(async () => {
      await approveWorkRecord(row.id)
      ElMessage.success('审核通过')
      detailVisible.value = false
      loadRecords()
    }).catch(() => {})
}

const handleReject = (row) => {
  currentRecord.value = row
  rejectForm.reason = ''
  rejectVisible.value = true
}

const confirmReject = async () => {
  await rejectFormRef.value.validate()
  rejectLoading.value = true
  try {
    await rejectWorkRecord(currentRecord.value.id, { reject_reason: rejectForm.reason })
    ElMessage.success('已驳回')
    rejectVisible.value = false
    detailVisible.value = false
    loadRecords()
  } finally {
    rejectLoading.value = false
  }
}

const handleBatchApprove = async () => {
  if (selectedIds.value.length === 0) return
  ElMessageBox.confirm(`确认批量通过选中的 ${selectedIds.value.length} 条工时记录吗？`, '确认', { type: 'warning' })
    .then(async () => {
      batchLoading.value = true
      try {
        await batchApprove({ ids: selectedIds.value })
        ElMessage.success('批量审核成功')
        loadRecords()
      } finally {
        batchLoading.value = false
      }
    }).catch(() => {})
}

onMounted(() => {
  getProjects().then(res => projects.value = res)
  getUsers().then(res => users.value = res)
  loadRecords()
})
</script>
