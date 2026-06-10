
&lt;template&gt;
  &lt;div&gt;
    &lt;div style="margin-bottom: 20px; display: flex; gap: 10px"&gt;
      &lt;el-input v-model="query.name" placeholder="姓名" style="width: 200px" clearable @clear="fetchRecords" /&gt;
      &lt;el-input v-model="query.phone" placeholder="手机号" style="width: 200px" clearable @clear="fetchRecords" /&gt;
      &lt;el-select v-model="query.prizeName" placeholder="奖品" style="width: 200px" clearable @clear="fetchRecords"&gt;
        &lt;el-option v-for="prize in prizeNames" :key="prize" :label="prize" :value="prize" /&gt;
      &lt;/el-select&gt;
      &lt;el-select v-model="query.isWin" placeholder="是否中奖" style="width: 150px" clearable @clear="fetchRecords"&gt;
        &lt;el-option label="是" value="true" /&gt;
        &lt;el-option label="否" value="false" /&gt;
      &lt;/el-select&gt;
      &lt;el-button type="primary" @click="fetchRecords"&gt;搜索&lt;/el-button&gt;
      &lt;el-button type="success" @click="exportRecords"&gt;导出Excel&lt;/el-button&gt;
    &lt;/div&gt;

    &lt;el-table :data="records" border&gt;
      &lt;el-table-column prop="id" label="ID" width="80" /&gt;
      &lt;el-table-column prop="name" label="姓名" width="120" /&gt;
      &lt;el-table-column prop="phone" label="手机号" width="150" /&gt;
      &lt;el-table-column prop="prizeName" label="奖品" /&gt;
      &lt;el-table-column label="是否中奖" width="100"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-tag :type="row.isWin ? 'success' : 'info'"&gt;{{ row.isWin ? '是' : '否' }}&lt;/el-tag&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column prop="status" label="状态" width="120"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-tag :type="row.status === '已领取' ? 'success' : 'warning'"&gt;{{ row.status }}&lt;/el-tag&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column prop="createdAt" label="时间" width="180" /&gt;
      &lt;el-table-column label="操作" width="120" fixed="right"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-button v-if="row.isWin &amp;&amp; row.status === '待领取'" type="primary" size="small" @click="claimRecord(row)"&gt;标记领取&lt;/el-button&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
    &lt;/el-table&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
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

const prizeNames = computed(() =&gt; {
  const names = new Set()
  records.value.forEach(r =&gt; {
    if (r.prizeName &amp;&amp; r.prizeName !== '未中奖') {
      names.add(r.prizeName)
    }
  })
  return Array.from(names)
})

const fetchRecords = async () =&gt; {
  const params = {}
  if (query.value.name) params.name = query.value.name
  if (query.value.phone) params.phone = query.value.phone
  if (query.value.prizeName) params.prizeName = query.value.prizeName
  if (query.value.isWin) params.isWin = query.value.isWin
  const res = await api.get('/records', { params })
  records.value = res.data
}

const claimRecord = async (row) =&gt; {
  await api.put(`/records/${row.id}/claim`)
  ElMessage.success('已标记领取')
  fetchRecords()
}

const exportRecords = () =&gt; {
  window.open(`http://localhost:8080/api/records/export?token=${localStorage.getItem('token')}`)
}

onMounted(() =&gt; {
  fetchRecords()
})
&lt;/script&gt;
