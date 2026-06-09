<template>
  <div>
    <el-card>
      <template #header>
        <span>内推列表</span>
      </template>
      <div style="margin-bottom: 20px">
        <el-select v-model="statusFilter" placeholder="按状态筛选" clearable @change="loadReferrals" style="width: 150px; margin-right: 10px">
          <el-option label="初筛中" value="screening" />
          <el-option label="面试中" value="interviewing" />
          <el-option label="发Offer" value="offer" />
          <el-option label="已入职" value="hired" />
          <el-option label="淘汰" value="rejected" />
        </el-select>
        <el-select v-model="jobFilter" placeholder="按职位筛选" clearable @change="loadReferrals" style="width: 200px; margin-right: 10px">
          <el-option v-for="job in jobs" :key="job.id" :label="job.title" :value="job.id" />
        </el-select>
        <el-input v-model="searchText" placeholder="搜索姓名/手机号" clearable @keyup.enter="loadReferrals" style="width: 200px; margin-right: 10px" />
        <el-button type="primary" @click="loadReferrals">搜索</el-button>
      </div>
      <el-table :data="referrals" stripe>
        <el-table-column prop="employee_name" label="推荐人" />
        <el-table-column prop="candidate_name" label="候选人姓名" />
        <el-table-column prop="candidate_phone" label="手机号" />
        <el-table-column prop="job_title" label="推荐职位" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="评分" width="100">
          <template #default="{ row }">
            <span v-if="row.evaluation_score">
              <el-rate v-model="row.evaluation_score" disabled show-score />
            </span>
            <span v-else style="color: #999">未评分</span>
          </template>
        </el-table-column>
        <el-table-column label="评价摘要" width="200">
          <template #default="{ row }">
            <span v-if="row.evaluation_pros" style="color: #67c23a">
              {{ row.evaluation_pros.length > 30 ? row.evaluation_pros.substring(0, 30) + '...' : row.evaluation_pros }}
            </span>
            <span v-else style="color: #999">暂无评价</span>
          </template>
        </el-table-column>
        <el-table-column prop="resume_path" label="简历">
          <template #default="{ row }">
            <el-link v-if="row.resume_path" type="primary" :href="row.resume_path" target="_blank">查看</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="推荐时间">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="openUpdateDialog(row)">更新</el-button>
            <el-button type="success" link @click="openEvaluationDialog(row)">评价</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @current-change="loadReferrals"
        layout="total, prev, pager, next"
        style="margin-top: 20px; justify-content: center"
      />
    </el-card>

    <el-dialog v-model="updateDialogVisible" title="更新状态" width="500px">
      <el-form :model="updateForm" label-width="100px">
        <el-form-item label="状态">
          <el-select v-model="updateForm.status" style="width: 100%">
            <el-option label="初筛中" value="screening" />
            <el-option label="面试中" value="interviewing" />
            <el-option label="发Offer" value="offer" />
            <el-option label="已入职" value="hired" />
            <el-option label="淘汰" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="淘汰原因" v-if="updateForm.status === 'rejected'">
          <el-input v-model="updateForm.reject_reason" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="HR备注">
          <el-input v-model="updateForm.hr_remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="updateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="loading">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="evaluationDialogVisible" title="评价候选人" width="500px">
      <el-form :model="evaluationForm" label-width="100px">
        <el-form-item label="评分">
          <el-rate v-model="evaluationForm.score" />
        </el-form-item>
        <el-form-item label="优点">
          <el-input v-model="evaluationForm.pros" type="textarea" :rows="3" placeholder="请描述候选人的优点" />
        </el-form-item>
        <el-form-item label="不足">
          <el-input v-model="evaluationForm.cons" type="textarea" :rows="3" placeholder="请描述候选人的不足" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="evaluationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleEvaluation" :loading="loading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const referrals = ref([])
const jobs = ref([])
const statusFilter = ref('')
const jobFilter = ref(null)
const searchText = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const updateDialogVisible = ref(false)
const evaluationDialogVisible = ref(false)
const loading = ref(false)
const updateForm = ref({
  id: null,
  status: '',
  reject_reason: '',
  hr_remark: ''
})
const evaluationForm = ref({
  id: null,
  pros: '',
  cons: '',
  score: 3
})

const loadReferrals = async () => {
  const params = { page: page.value, page_size: pageSize.value }
  if (statusFilter.value) params.status = statusFilter.value
  if (jobFilter.value) params.job_id = jobFilter.value
  if (searchText.value) params.search = searchText.value
  const res = await api.get('/admin/referrals', { params })
  referrals.value = res.data
  total.value = res.total
}

const loadJobs = async () => {
  const res = await api.get('/admin/jobs')
  jobs.value = res
}

const openUpdateDialog = (row) => {
  updateForm.value = {
    id: row.id,
    status: row.status,
    reject_reason: row.reject_reason || '',
    hr_remark: row.hr_remark || ''
  }
  updateDialogVisible.value = true
}

const openEvaluationDialog = (row) => {
  evaluationForm.value = {
    id: row.id,
    pros: row.evaluation_pros || '',
    cons: row.evaluation_cons || '',
    score: row.evaluation_score || 3
  }
  evaluationDialogVisible.value = true
}

const handleUpdate = async () => {
  if (!updateForm.value.status) {
    ElMessage.warning('请选择状态')
    return
  }
  loading.value = true
  try {
    await api.put(`/admin/referrals/${updateForm.value.id}/status`, updateForm.value)
    ElMessage.success('更新成功')
    updateDialogVisible.value = false
    loadReferrals()
  } catch (e) {
    ElMessage.error('更新失败')
  } finally {
    loading.value = false
  }
}

const handleEvaluation = async () => {
  if (!evaluationForm.value.score) {
    ElMessage.warning('请选择评分')
    return
  }
  loading.value = true
  try {
    await api.put(`/admin/referrals/${evaluationForm.value.id}/evaluation`, evaluationForm.value)
    ElMessage.success('评价成功')
    evaluationDialogVisible.value = false
    loadReferrals()
  } catch (e) {
    ElMessage.error('评价失败')
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
    interviewing: '面试中',
    offer: '发Offer',
    hired: '已入职',
    rejected: '淘汰'
  }
  return map[status] || status
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  loadReferrals()
  loadJobs()
})
</script>

<style scoped>
</style>
