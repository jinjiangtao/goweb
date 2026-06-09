<template>
  <div>
    <el-card>
      <template #header>
        <span>职位列表</span>
      </template>
      <el-table :data="jobs" stripe>
        <el-table-column prop="title" label="职位名称" />
        <el-table-column prop="requirement" label="职位要求" show-overflow-tooltip />
        <el-table-column prop="salary_range" label="薪资范围" />
        <el-table-column prop="location" label="工作地点" />
        <el-table-column prop="created_by_name" label="发布人" />
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
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const jobs = ref([])

const loadJobs = async () => {
  const res = await api.get('/admin/jobs')
  jobs.value = res
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  loadJobs()
})
</script>

<style scoped>
</style>
