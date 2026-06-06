<template>
  <div class="bookings">
    <div class="header">
      <h2>预订管理</h2>
      <div class="filters">
        <el-date-picker v-model="filterDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" clearable />
        <el-select v-model="filterRoom" placeholder="选择会议室" clearable style="width: 180px">
          <el-option v-for="room in rooms" :key="room.id" :label="room.name" :value="room.id" />
        </el-select>
        <el-select v-model="filterStatus" placeholder="选择状态" clearable style="width: 140px">
          <el-option label="全部" value="" />
          <el-option label="已预订" value="1" />
          <el-option label="已取消" value="2" />
        </el-select>
        <el-button type="primary" @click="fetchBookings">查询</el-button>
      </div>
    </div>

    <el-table :data="bookings" stripe style="width: 100%">
      <el-table-column label="会议室">
        <template #default="{ row }">{{ row.room?.name }}</template>
      </el-table-column>
      <el-table-column prop="name" label="预订人" width="120" />
      <el-table-column prop="phone" label="手机号" width="130" />
      <el-table-column prop="date" label="日期" width="120" />
      <el-table-column label="时间段" width="140">
        <template #default="{ row }">{{ row.start_time }} - {{ row.end_time }}</template>
      </el-table-column>
      <el-table-column prop="purpose" label="用途" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '已预订' : '已取消' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="cancel_reason" label="取消原因" show-overflow-tooltip />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button v-if="row.status === 1" size="small" type="danger" @click="handleCancel(row)">取消</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        @current-change="fetchBookings"
        layout="total, prev, pager, next"
      />
    </div>

    <el-dialog v-model="cancelDialogVisible" title="取消预订" width="500px">
      <el-form :model="cancelForm" label-width="80px">
        <el-form-item label="取消原因">
          <el-input v-model="cancelForm.cancel_reason" type="textarea" :rows="4" placeholder="请输入取消原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="cancelDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCancel" :loading="canceling">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { ElMessage } from 'element-plus'

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

const fetchBookings = async () => {
  const params = {
    page: page.value,
    page_size: pageSize.value
  }
  if (filterDate.value) params.date = filterDate.value
  if (filterRoom.value) params.room_id = filterRoom.value
  if (filterStatus.value) params.status = filterStatus.value
  try {
    const res = await authStore.api.get('/bookings', { params })
    bookings.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  }
}

const fetchRooms = async () => {
  try {
    const res = await authStore.api.get('/rooms')
    rooms.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const handleCancel = (booking) => {
  cancelingBooking.value = booking
  cancelForm.value.cancel_reason = ''
  cancelDialogVisible.value = true
}

const confirmCancel = async () => {
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

onMounted(() => {
  fetchBookings()
  fetchRooms()
})
</script>

<style scoped>
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
</style>
