
&lt;template&gt;
  &lt;div&gt;
    &lt;div style="margin-bottom: 20px"&gt;
      &lt;el-button type="primary" @click="openDialog()"&gt;添加奖品&lt;/el-button&gt;
    &lt;/div&gt;
    &lt;el-table :data="prizes" border&gt;
      &lt;el-table-column prop="id" label="ID" width="80" /&gt;
      &lt;el-table-column prop="name" label="奖品名称" /&gt;
      &lt;el-table-column prop="probability" label="中奖概率(%)" width="120" /&gt;
      &lt;el-table-column label="库存" width="150"&gt;
        &lt;template #default="{ row }"&gt;
          {{ row.stockUsed }} / {{ row.stock }}
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column prop="description" label="描述" show-overflow-tooltip /&gt;
      &lt;el-table-column label="图片" width="100"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-image v-if="row.imageUrl" :src="row.imageUrl" style="width: 60px; height: 60px" fit="cover" /&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column label="状态" width="100"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-switch :model-value="row.enabled" @change="togglePrize(row)" /&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column label="操作" width="180" fixed="right"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-button type="primary" size="small" @click="openDialog(row)"&gt;编辑&lt;/el-button&gt;
          &lt;el-button type="danger" size="small" @click="deletePrize(row)"&gt;删除&lt;/el-button&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
    &lt;/el-table&gt;

    &lt;el-dialog v-model="dialogVisible" :title="isEdit ? '编辑奖品' : '添加奖品'" width="500px"&gt;
      &lt;el-form :model="prizeForm" label-width="100px"&gt;
        &lt;el-form-item label="奖品名称"&gt;
          &lt;el-input v-model="prizeForm.name" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="中奖概率(%)"&gt;
          &lt;el-input-number v-model="prizeForm.probability" :min="0" :max="100" :precision="2" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="库存"&gt;
          &lt;el-input-number v-model="prizeForm.stock" :min="0" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="描述"&gt;
          &lt;el-input v-model="prizeForm.description" type="textarea" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="图片链接"&gt;
          &lt;el-input v-model="prizeForm.imageUrl" /&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
      &lt;template #footer&gt;
        &lt;el-button @click="dialogVisible = false"&gt;取消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="savePrize"&gt;保存&lt;/el-button&gt;
      &lt;/template&gt;
    &lt;/el-dialog&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
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

const fetchPrizes = async () =&gt; {
  const res = await api.get('/admin/prizes')
  prizes.value = res.data
}

const openDialog = (row = null) =&gt; {
  if (row) {
    isEdit.value = true
    prizeForm.value = { ...row }
  } else {
    isEdit.value = false
    prizeForm.value = { id: null, name: '', probability: 0, stock: 0, description: '', imageUrl: '' }
  }
  dialogVisible.value = true
}

const savePrize = async () =&gt; {
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

const deletePrize = async (row) =&gt; {
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

const togglePrize = async (row) =&gt; {
  await api.put(`/admin/prizes/${row.id}/toggle`)
  fetchPrizes()
}

onMounted(() =&gt; {
  fetchPrizes()
})
&lt;/script&gt;

