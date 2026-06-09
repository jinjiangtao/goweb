<template>
  <div>
    <el-row :gutter="20" style="margin-bottom: 20px">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num">{{ stats.total }}</div>
            <div class="stat-label">总推荐数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num">{{ stats.status?.hired || 0 }}</div>
            <div class="stat-label">已入职</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num">{{ stats.status?.interviewing || 0 }}</div>
            <div class="stat-label">面试中</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num">{{ stats.status?.offer || 0 }}</div>
            <div class="stat-label">Offer</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的推荐</span>
          <el-button type="primary" @click="openDialog">推荐候选人</el-button>
        </div>
      </template>
      <div style="margin-bottom: 20px">
        <el-select v-model="statusFilter" placeholder="按状态筛选" clearable @change="loadReferrals" style="width: 150px; margin-right: 10px">
          <el-option label="初筛中" value="screening" />
          <el-option label="面试中" value="interviewing" />
          <el-option label="发Offer" value="offer" />
          <el-option label="已入职" value="hired" />
          <el-option label="淘汰" value="rejected" />
        </el-select>
      </div>
      <el-table :data="referrals" stripe>
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
        <el-table-column prop="hr_remark" label="HR备注" show-overflow-tooltip />
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
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button 
              type="success" 
              link 
              @click="openEvaluationDialog(row)"
              v-if="row.evaluation_score"
            >查看评价</el-button>
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

    <el-dialog v-model="dialogVisible" title="推荐候选人" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="选择职位">
          <el-select v-model="form.job_id" placeholder="请选择职位" style="width: 100%">
            <el-option v-for="job in jobs" :key="job.id" :label="job.title" :value="job.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="候选人姓名">
          <el-input v-model="form.candidate_name" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.candidate_phone" />
        </el-form-item>
        <el-form-item label="简历">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            :on-change="handleFileChange"
            accept=".pdf,.doc,.docx"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">支持pdf/doc/docx格式</div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="loading">推荐</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="evaluationDialogVisible" title="查看评价" width="500px">
      <el-descriptions :column="1" border v-if="currentEvaluation">
        <el-descriptions-item label="评分">
          <el-rate v-model="currentEvaluation.evaluation_score" disabled show-score />
        </el-descriptions-item>
        <el-descriptions-item label="优点">
          {{ currentEvaluation.evaluation_pros || '暂无' }}
        </el-descriptions-item>
        <el-descriptions-item label="不足">
          {{ currentEvaluation.evaluation_cons || '暂无' }}
        </el-descriptions-item>
        <el-descriptions-item label="评价时间">
          {{ currentEvaluation.evaluation_time ? formatDate(currentEvaluation.evaluation_time) : '暂无' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="evaluationDialogVisible = false">关闭</el-button>
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
const stats = ref({ total: 0, status: {} })
const dialogVisible = ref(false)
const evaluationDialogVisible = ref(false)
const loading = ref(false)
const statusFilter = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const currentEvaluation = ref(null)
const form = ref({
  job_id: null,
  candidate_name: '',
  candidate_phone: ''
})
let selectedFile = null

const loadReferrals = async () => {
  const params = { page: page.value, page_size: pageSize.value }
  if (statusFilter.value) params.status = statusFilter.value
  const res = await api.get('/referrals/my', { params })
  referrals.value = res.data
  total.value = res.total
}

const loadStats = async () => {
  const res = await api.get('/referrals/my/stats')
  stats.value = res
}

const loadJobs = async () => {
  const res = await api.get('/jobs/my')
  jobs.value = res.filter(j => j.status === 'active')
}

const openDialog = () => {
  form.value = { job_id: null, candidate_name: '', candidate_phone: '' }
  selectedFile = null
  dialogVisible.value = true
}

const openEvaluationDialog = async (row) => {
  try {
    const res = await api.get(`/referrals/${row.id}/evaluation`)
    currentEvaluation.value = res
    evaluationDialogVisible.value = true
  } catch (e) {
    ElMessage.error('获取评价失败')
  }
}

const handleFileChange = (file) => {
  selectedFile = file.raw
}

const handleCreate = async () => {
  if (!form.value.job_id || !form.value.candidate_name || !form.value.candidate_phone) {
    ElMessage.warning('请填写必填项')
    return
  }
  
  loading.value = true
  try {
    const formData = new FormData()
    formData.append('job_id', form.value.job_id)
    formData.append('candidate_name', form.value.candidate_name)
    formData.append('candidate_phone', form.value.candidate_phone)
    if (selectedFile) formData.append('resume', selectedFile)
    
    await api.post('/referrals', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    ElMessage.success('推荐成功')
    dialogVisible.value = false
    loadReferrals()
    loadStats()
  } catch (e) {
    ElMessage.error('推荐失败')
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
  loadStats()
  loadJobs()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.stat-item {
  text-align: center;
}
.stat-num {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
}
.stat-label {
  color: #999;
  margin-top: 5px;
}
</style>
