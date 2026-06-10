<template>
  <div>
    <div style="margin-bottom: 20px">
      <el-button type="primary" @click="openDialog()">添加奖品</el-button>
    </div>
    <el-table :data="prizes" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="奖品名称" />
      <el-table-column prop="probability" label="中奖概率(%)" width="120" />
      <el-table-column label="库存" width="150">
        <template #default="{ row }">
          {{ row.stockUsed }} / {{ row.stock }}
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column label="图片" width="100">
        <template #default="{ row }">
          <el-image v-if="row.imageUrl" :src="row.imageUrl" style="width: 60px; height: 60px" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch :model-value="row.enabled" @change="togglePrize(row)" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="openDialog(row)">编辑</el-button>
          <el-button type="danger" size="small" @click="deletePrize(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑奖品' : '添加奖品'" width="500px">
      <el-form :model="prizeForm" label-width="100px">
        <el-form-item label="奖品名称">
          <el-input v-model="prizeForm.name" />
        </el-form-item>
        <el-form-item label="中奖概率(%)">
          <el-input-number v-model="prizeForm.probability" :min="0" :max="100" :precision="2" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="prizeForm.stock" :min="0" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="prizeForm.description" type="textarea" />
        </el-form-item>
        <el-form-item label="图片链接">
          <el-input v-model="prizeForm.imageUrl" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="savePrize">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const prizes = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const prizeForm = ref({
  id: null,
  name: '',
  probability: 0,
  stock: 0,
  description: '',
  imageUrl: ''
})

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})

const fetchPrizes = async () => {
  const res = await api.get('/admin/prizes')
  prizes.value = res.data
}

const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    prizeForm.value = { ...row }
  } else {
    isEdit.value = false
    prizeForm.value = { id: null, name: '', probability: 0, stock: 0, description: '', imageUrl: '' }
  }
  dialogVisible.value = true
}

const savePrize = async () => {
  try {
    if (isEdit.value) {
      await api.put(`/admin/prizes/${prizeForm.value.id}`, prizeForm.value)
    } else {
      await api.post('/admin/prizes', prizeForm.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchPrizes()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '保存失败')
  }
}

const deletePrize = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除吗?', '提示')
    await api.delete(`/admin/prizes/${row.id}`)
    ElMessage.success('删除成功')
    fetchPrizes()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const togglePrize = async (row) => {
  await api.put(`/admin/prizes/${row.id}/toggle`)
  fetchPrizes()
}

onMounted(() => {
  fetchPrizes()
})
</script>
