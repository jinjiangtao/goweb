<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的职位</span>
          <el-button type="primary" @click="openDialog">发布新职位</el-button>
        </div>
      </template>
      <el-table :data="jobs" stripe>
        <el-table-column prop="title" label="职位名称" />
        <el-table-column prop="requirement" label="职位要求" show-overflow-tooltip />
        <el-table-column prop="salary_range" label="薪资范围" />
        <el-table-column prop="location" label="工作地点" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '发布中' : '已关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="toggleStatus(row)">
              {{ row.status === 'active' ? '关闭' : '开启' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="发布职位" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="职位名称">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="职位要求">
          <el-input v-model="form.requirement" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="薪资范围">
          <el-input v-model="form.salary_range" placeholder="如：15-25K" />
        </el-form-item>
        <el-form-item label="工作地点">
          <el-input v-model="form.location" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="loading">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const jobs = ref([])
const dialogVisible = ref(false)
const loading = ref(false)
const form = ref({
  title: '',
  requirement: '',
  salary_range: '',
  location: ''
})

const loadJobs = async () => {
  const res = await api.get('/jobs/my')
  jobs.value = res
}

const openDialog = () => {
  form.value = { title: '', requirement: '', salary_range: '', location: '' }
  dialogVisible.value = true
}

const handleCreate = async () => {
  if (!form.value.title || !form.value.requirement) {
    ElMessage.warning('请填写必填项')
    return
  }
  loading.value = true
  try {
    await api.post('/jobs', form.value)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadJobs()
  } catch (e) {
    ElMessage.error('创建失败')
  } finally {
    loading.value = false
  }
}

const toggleStatus = async (row) => {
  const newStatus = row.status === 'active' ? 'closed' : 'active'
  await api.put(`/jobs/${row.id}/status`, { status: newStatus })
  ElMessage.success('更新成功')
  loadJobs()
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  loadJobs()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
