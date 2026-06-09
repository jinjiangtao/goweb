<template>
  <div>
    <el-card>
      <template #header>
        <span>投递管理</span>
      </template>
      <div style="margin-bottom: 20px">
        <el-select v-model="statusFilter" placeholder="按状态筛选" clearable @change="loadSubmissions" style="width: 150px; margin-right: 10px">
          <el-option label="初筛中" value="screening" />
          <el-option label="已面试" value="interviewing" />
          <el-option label="已发offer" value="offer" />
          <el-option label="已入职" value="hired" />
          <el-option label="已淘汰" value="rejected" />
        </el-select>
        <el-select v-model="jobFilter" placeholder="按职位筛选" clearable @change="loadSubmissions" style="width: 200px; margin-right: 10px">
          <el-option v-for="job in jobs" :key="job.id" :label="job.title" :value="job.id" />
        </el-select>
        <el-input v-model="searchText" placeholder="搜索姓名/手机号" clearable @keyup.enter="loadSubmissions" style="width: 200px; margin-right: 10px" />
        <el-button type="primary" @click="loadSubmissions">搜索</el-button>
      </div>
      <el-table :data="submissions" stripe>
        <el-table-column prop="name" label="投递人" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="jobTitle" label="投递职位" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
            <el-tag v-if="row.converted" type="success" size="small" style="margin-left: 5px">已转内推</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resumePath" label="简历">
          <template #default="{ row }">
            <el-link v-if="row.resumePath" type="primary" :href="row.resumePath" target="_blank">下载</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="投递时间">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="openUpdateDialog(row)" :disabled="row.converted">更新状态</el-button>
            <el-button type="success" link @click="openConvertDialog(row)" :disabled="row.converted">转为内推</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @current-change="loadSubmissions"
        layout="total, prev, pager, next"
        style="margin-top: 20px; justify-content: center"
      />
    </el-card>

    <el-dialog v-model="updateDialogVisible" title="更新状态" width="500px">
      <el-form :model="updateForm" label-width="100px">
        <el-form-item label="状态">
          <el-select v-model="updateForm.status" style="width: 100%">
            <el-option label="初筛中" value="screening" />
            <el-option label="已面试" value="interviewing" />
            <el-option label="已发offer" value="offer" />
            <el-option label="已入职" value="hired" />
            <el-option label="已淘汰" value="rejected" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="updateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="loading">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="convertDialogVisible" title="转为内推" width="500px">
      <el-form :model="convertForm" label-width="100px">
        <el-form-item label="推荐人">
          <el-select v-model="convertForm.employeeId" placeholder="请选择推荐人" style="width: 100%">
            <el-option v-for="emp in employees" :key="emp.id" :label="emp.real_name" :value="emp.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="convertDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleConvert" :loading="loading">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const submissions = ref([])
const jobs = ref([])
const employees = ref([])
const statusFilter = ref('')
const jobFilter = ref(null)
const searchText = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const updateDialogVisible = ref(false)
const convertDialogVisible = ref(false)
const loading = ref(false)
const updateForm = ref({
  id: null,
  status: ''
})
const convertForm = ref({
  id: null,
  employeeId: null
})

const loadSubmissions = async () => {
  const params = { page: page.value, page_size: pageSize.value }
  if (statusFilter.value) params.status = statusFilter.value
  if (jobFilter.value) params.job_id = jobFilter.value
  if (searchText.value) params.search = searchText.value
  const res = await api.get('/admin/submissions', { params })
  submissions.value = res.data
  total.value = res.total
}

const loadJobs = async () => {
  const res = await api.get('/admin/jobs')
  jobs.value = res
}

const loadEmployees = async () => {
  const res = await api.get('/admin/users')
  employees.value = res.filter(u => u.role === 'employee' && u.status === 'enabled')
}

const openUpdateDialog = (row) => {
  updateForm.value = {
    id: row.id,
    status: row.status
  }
  updateDialogVisible.value = true
}

const handleUpdate = async () => {
  if (!updateForm.value.status) {
    ElMessage.warning('请选择状态')
    return
  }
  loading.value = true
  try {
    await api.put(`/admin/submissions/${updateForm.value.id}/status`, updateForm.value)
    ElMessage.success('更新成功')
    updateDialogVisible.value = false
    loadSubmissions()
  } catch (e) {
    ElMessage.error('更新失败')
  } finally {
    loading.value = false
  }
}

const openConvertDialog = (row) => {
  convertForm.value = {
    id: row.id,
    employeeId: null
  }
  convertDialogVisible.value = true
}

const handleConvert = async () => {
  if (!convertForm.value.employeeId) {
    ElMessage.warning('请选择推荐人')
    return
  }
  loading.value = true
  try {
    await api.post(`/admin/submissions/${convertForm.value.id}/convert`, { employee_id: convertForm.value.employeeId })
    ElMessage.success('转为内推成功')
    convertDialogVisible.value = false
    loadSubmissions()
  } catch (e) {
    ElMessage.error('操作失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const map = {
    screening: 'info',
    interviewing: 'warning',
    offer: 'primary',
    hired: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

const getStatusLabel = (status) => {
  const map = {
    screening: '初筛中',
    interviewing: '已面试',
    offer: '已发offer',
    hired: '已入职',
    rejected: '已淘汰'
  }
  return map[status] || status
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  loadSubmissions()
  loadJobs()
  loadEmployees()
})
</script>

<style scoped>
</style>
