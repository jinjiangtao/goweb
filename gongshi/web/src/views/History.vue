<template>
  <div class="page-container">
    <div class="card">
      <div class="card-header">
        <div class="card-title">
          <el-icon style="vertical-align: middle; margin-right: 6px; color: #409eff"><Document /></el-icon>
          历史填报记录
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

      <el-table :data="records" border v-loading="loading">
        <el-table-column prop="work_date" label="日期" width="120" sortable>
          <template #default="{ row }">
            <span style="font-weight: 600">{{ row.work_date }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="project.name" label="项目名称" min-width="160">
          <template #default="{ row }">
            <el-tag size="small">{{ row.project?.code }}</el-tag> {{ row.project?.name }}
          </template>
        </el-table-column>
        <el-table-column prop="hours" label="工时" width="90" align="center" sortable>
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
        <el-table-column prop="rejectReason" label="驳回原因" width="160" show-overflow-tooltip />
        <el-table-column prop="reviewer.name" label="审核人" width="100" />
        <el-table-column prop="reviewed_at" label="审核时间" width="160" />
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">详情</el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)" :disabled="row.status === 'approved'">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)" :disabled="row.status === 'approved'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="detailVisible" title="工时详情" width="560px">
      <el-descriptions :column="2" border v-if="currentRecord">
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
        <el-descriptions-item label="填报人">{{ currentRecord.user?.name }}</el-descriptions-item>
        <el-descriptions-item label="填报时间">{{ currentRecord.created_at }}</el-descriptions-item>
        <el-descriptions-item label="工作内容" :span="2">{{ currentRecord.content }}</el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.reject_reason" label="驳回原因" :span="2">
          <span style="color: #f56c6c">{{ currentRecord.reject_reason }}</span>
        </el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.reviewer" label="审核人">{{ currentRecord.reviewer?.name }}</el-descriptions-item>
        <el-descriptions-item v-if="currentRecord.reviewed_at" label="审核时间">{{ currentRecord.reviewed_at }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <el-dialog v-model="editVisible" title="编辑工时" width="500px" destroy-on-close>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="90px">
        <el-form-item label="工作日期">
          <el-date-picker v-model="form.work_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="所属项目" prop="project_id">
          <el-select v-model="form.project_id" style="width: 100%" filterable>
            <el-option v-for="p in projects" :key="p.id" :label="p.code + ' - ' + p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="工作时长" prop="hours">
          <el-input-number v-model="form.hours" :min="0.5" :max="24" :step="0.5" controls-position="right" style="width: 100%" />
        </el-form-item>
        <el-form-item label="工作内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="4" maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleEditSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { getProjects } from '../api/project'
import { getWorkRecords, updateWorkRecord, deleteWorkRecord } from '../api/workrecord'

const projects = ref([])
const records = ref([])
const loading = ref(false)
const detailVisible = ref(false)
const editVisible = ref(false)
const currentRecord = ref(null)
const submitting = ref(false)
const formRef = ref()

const filterForm = reactive({
  dates: [dayjs().subtract(3, 'month').format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')],
  project_id: null,
  status: ''
})

const form = reactive({
  project_id: null,
  work_date: '',
  hours: 8,
  content: ''
})

const rules = {
  project_id: [{ required: true, message: '请选择项目', trigger: 'change' }],
  hours: [{ required: true, message: '请填写工时', trigger: 'blur' }],
  content: [{ required: true, message: '请填写工作内容', trigger: 'blur', min: 2 }]
}

const resetFilter = () => {
  filterForm.dates = [dayjs().subtract(3, 'month').format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')]
  filterForm.project_id = null
  filterForm.status = ''
  loadRecords()
}

const loadRecords = async () => {
  loading.value = true
  try {
    const params = {
      start_date: filterForm.dates?.[0] || '',
      end_date: filterForm.dates?.[1] || ''
    }
    if (filterForm.project_id) params.project_id = filterForm.project_id
    if (filterForm.status) params.status = filterForm.status
    records.value = await getWorkRecords(params)
  } finally {
    loading.value = false
  }
}

const viewDetail = (row) => {
  currentRecord.value = row
  detailVisible.value = true
}

const handleEdit = (row) => {
  Object.assign(form, {
    project_id: row.project_id,
    work_date: row.work_date,
    hours: row.hours,
    content: row.content
  })
  currentRecord.value = row
  editVisible.value = true
}

const handleEditSubmit = async () => {
  await formRef.value.validate()
  submitting.value = true
  try {
    await updateWorkRecord(currentRecord.value.id, { ...form })
    ElMessage.success('修改成功')
    editVisible.value = false
    loadRecords()
  } finally {
    submitting.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除该工时记录吗？`, '提示', { type: 'warning' })
    .then(async () => {
      await deleteWorkRecord(row.id)
      ElMessage.success('删除成功')
      loadRecords()
    }).catch(() => {})
}

onMounted(() => {
  getProjects().then(res => projects.value = res)
  loadRecords()
})
</script>
