<template>
  <div>
    <div style="margin-bottom: 20px; display: flex; gap: 10px">
      <el-input v-model="query.name" placeholder="姓名" style="width: 200px" clearable @clear="fetchRecords" />
      <el-input v-model="query.phone" placeholder="手机号" style="width: 200px" clearable @clear="fetchRecords" />
      <el-select v-model="query.prizeName" placeholder="奖品" style="width: 200px" clearable @clear="fetchRecords">
        <el-option v-for="prize in prizeNames" :key="prize" :label="prize" :value="prize" />
      </el-select>
      <el-select v-model="query.isWin" placeholder="是否中奖" style="width: 150px" clearable @clear="fetchRecords">
        <el-option label="是" value="true" />
        <el-option label="否" value="false" />
      </el-select>
      <el-button type="primary" @click="fetchRecords">搜索</el-button>
      <el-button type="success" @click="exportRecords">导出Excel</el-button>
    </div>

    <el-table :data="records" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column prop="phone" label="手机号" width="150" />
      <el-table-column prop="prizeName" label="奖品" />
      <el-table-column label="是否中奖" width="100">
        <template #default="{ row }">
          <el-tag :type="row.isWin ? 'success' : 'info'">{{ row.isWin ? '是' : '否' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="row.status === '已领取' ? 'success' : 'warning'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="时间" width="180" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.isWin && row.status === '待领取'" type="primary" size="small" @click="claimRecord(row)">标记领取</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const records = ref([])
const query = ref({
  name: '',
  phone: '',
  prizeName: '',
  isWin: ''
})

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})

const prizeNames = computed(() => {
  const names = new Set()
  records.value.forEach(r => {
    if (r.prizeName && r.prizeName !== '未中奖') {
      names.add(r.prizeName)
    }
  })
  return Array.from(names)
})

const fetchRecords = async () => {
  const params = {}
  if (query.value.name) params.name = query.value.name
  if (query.value.phone) params.phone = query.value.phone
  if (query.value.prizeName) params.prizeName = query.value.prizeName
  if (query.value.isWin) params.isWin = query.value.isWin
  const res = await api.get('/admin/records', { params })
  records.value = res.data
}

const claimRecord = async (row) => {
  await api.put(`/admin/records/${row.id}/claim`)
  ElMessage.success('已标记领取')
  fetchRecords()
}

const exportRecords = () => {
  window.open(`http://localhost:8080/api/admin/records/export?token=${localStorage.getItem('token')}`)
}

onMounted(() => {
  fetchRecords()
})
</script>
