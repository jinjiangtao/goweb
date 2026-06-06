
&lt;template&gt;
  &lt;div class="bookings"&gt;
    &lt;div class="header"&gt;
      &lt;h2&gt;预订管理&lt;/h2&gt;
      &lt;div class="filters"&gt;
        &lt;el-date-picker v-model="filterDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" clearable /&gt;
        &lt;el-select v-model="filterRoom" placeholder="选择会议室" clearable style="width: 180px"&gt;
          &lt;el-option v-for="room in rooms" :key="room.id" :label="room.name" :value="room.id" /&gt;
        &lt;/el-select&gt;
        &lt;el-select v-model="filterStatus" placeholder="选择状态" clearable style="width: 140px"&gt;
          &lt;el-option label="全部" value="" /&gt;
          &lt;el-option label="已预订" value="1" /&gt;
          &lt;el-option label="已取消" value="2" /&gt;
        &lt;/el-select&gt;
        &lt;el-button type="primary" @click="fetchBookings"&gt;查询&lt;/el-button&gt;
      &lt;/div&gt;
    &lt;/div&gt;

    &lt;el-table :data="bookings" stripe style="width: 100%"&gt;
      &lt;el-table-column label="会议室"&gt;
        &lt;template #default="{ row }"&gt;{{ row.room?.name }}&lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column prop="name" label="预订人" width="120" /&gt;
      &lt;el-table-column prop="phone" label="手机号" width="130" /&gt;
      &lt;el-table-column prop="date" label="日期" width="120" /&gt;
      &lt;el-table-column label="时间段" width="140"&gt;
        &lt;template #default="{ row }"&gt;{{ row.start_time }} - {{ row.end_time }}&lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column prop="purpose" label="用途" show-overflow-tooltip /&gt;
      &lt;el-table-column label="状态" width="100"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-tag :type="row.status === 1 ? 'success' : 'danger'"&gt;
            {{ row.status === 1 ? '已预订' : '已取消' }}
          &lt;/el-tag&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column prop="cancel_reason" label="取消原因" show-overflow-tooltip /&gt;
      &lt;el-table-column label="操作" width="120"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-button v-if="row.status === 1" size="small" type="danger" @click="handleCancel(row)"&gt;取消&lt;/el-button&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
    &lt;/el-table&gt;

    &lt;div class="pagination"&gt;
      &lt;el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        @current-change="fetchBookings"
        layout="total, prev, pager, next"
      /&gt;
    &lt;/div&gt;

    &lt;el-dialog v-model="cancelDialogVisible" title="取消预订" width="500px"&gt;
      &lt;el-form :model="cancelForm" label-width="80px"&gt;
        &lt;el-form-item label="取消原因"&gt;
          &lt;el-input v-model="cancelForm.cancel_reason" type="textarea" :rows="4" placeholder="请输入取消原因" /&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
      &lt;template #footer&gt;
        &lt;el-button @click="cancelDialogVisible = false"&gt;取消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="confirmCancel" :loading="canceling"&gt;确定&lt;/el-button&gt;
      &lt;/template&gt;
    &lt;/el-dialog&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const authStore = useAuthStore()
const bookings = ref([])
const rooms = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const filterDate = ref('')
const filterRoom = ref('')
const filterStatus = ref('')
const cancelDialogVisible = ref(false)
const canceling = ref(false)
const cancelingBooking = ref(null)
const cancelForm = ref({ cancel_reason: '' })

const fetchBookings = async () =&gt; {
  const params = {
    page: page.value,
    page_size: pageSize.value
  }
  if (filterDate.value) params.date = filterDate.value
  if (filterRoom.value) params.room_id = filterRoom.value
  if (filterStatus.value) params.status = filterStatus.value
  const res = await authStore.api.get('/bookings', { params })
  bookings.value = res.data.list
  total.value = res.data.total
}

const fetchRooms = async () =&gt; {
  const res = await authStore.api.get('/rooms')
  rooms.value = res.data
}

const handleCancel = (booking) =&gt; {
  cancelingBooking.value = booking
  cancelForm.value.cancel_reason = ''
  cancelDialogVisible.value = true
}

const confirmCancel = async () =&gt; {
  if (!cancelForm.value.cancel_reason) {
    ElMessage.warning('请输入取消原因')
    return
  }
  canceling.value = true
  try {
    await authStore.api.put(`/bookings/${cancelingBooking.value.id}/cancel`, cancelForm.value)
    ElMessage.success('取消成功')
    cancelDialogVisible.value = false
    fetchBookings()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '取消失败')
  } finally {
    canceling.value = false
  }
}

onMounted(() =&gt; {
  fetchBookings()
  fetchRooms()
})
&lt;/script&gt;

&lt;style scoped&gt;
.bookings {
  background: #1e293b;
  border-radius: 12px;
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h2 {
  color: #e2e8f0;
  font-size: 20px;
}

.filters {
  display: flex;
  gap: 12px;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}
&lt;/style&gt;

