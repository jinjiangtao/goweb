<template>
  <div class="my-bookings">
    <div class="header">
      <el-button text @click="$router.back()" class="back-btn">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>我的预订</h1>
      <div style="width: 32px;"></div>
    </div>

    <div class="content">
      <div class="booking-list" v-loading="loading">
        <el-empty v-if="!loading && bookings.length === 0" description="暂无预订记录" />

        <div class="booking-card" v-for="booking in bookings" :key="booking.id">
          <div class="booking-header">
            <div class="room-name">{{ booking.room?.name || '会议室' }}</div>
            <el-tag :type="booking.status === 1 ? 'success' : 'info'" size="small">
              {{ booking.status === 1 ? '已预订' : '已取消' }}
            </el-tag>
          </div>
          <div class="booking-info">
            <div class="info-item">
              <el-icon><Calendar /></el-icon>
              <span>{{ booking.date }}</span>
            </div>
            <div class="info-item">
              <el-icon><Clock /></el-icon>
              <span>{{ booking.start_time }} - {{ booking.end_time }}</span>
            </div>
            <div class="info-item" v-if="booking.purpose">
              <el-icon><Document /></el-icon>
              <span>{{ booking.purpose }}</span>
            </div>
          </div>
          <div class="booking-actions" v-if="booking.status === 1">
            <el-button type="danger" size="small" @click="cancelBooking(booking)">取消预订</el-button>
          </div>
        </div>
      </div>
    </div>

    <div class="footer-btn">
      <el-button type="primary" class="back-home-btn" @click="$router.push('/')">
        <el-icon><House /></el-icon>
        返回预订
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Calendar, Clock, Document, House } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const bookings = ref([])

onMounted(() => {
  fetchBookings()
})

async function fetchBookings() {
  loading.value = true
  try {
    const res = await api.get('/bookings/my')
    bookings.value = res.data
  } catch (error) {
    ElMessage.error('获取预订记录失败')
  } finally {
    loading.value = false
  }
}

async function cancelBooking(booking) {
  try {
    await ElMessageBox.confirm(
      '确定要取消这个预订吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return
  }

  try {
    await api.delete(`/bookings/${booking.id}`)
    ElMessage.success('取消成功')
    fetchBookings()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '取消失败')
  }
}
</script>

<style scoped>
.my-bookings {
  min-height: 100vh;
  padding-bottom: 80px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
  color: white;
}

.back-btn {
  color: white;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header h1 {
  font-size: 18px;
  font-weight: 600;
}

.content {
  padding: 0 0 20px 0;
}

.booking-list {
  padding: 16px 20px;
}

.booking-card {
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 12px;
}

.booking-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.room-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
}

.booking-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-secondary);
  font-size: 14px;
}

.booking-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.footer-btn {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 430px;
  padding: 16px 20px;
  background: white;
  border-top: 1px solid var(--border-color);
}

.back-home-btn {
  width: 100%;
  border-radius: 6px;
  background-color: var(--primary-color);
}

.back-home-btn:hover {
  background-color: var(--primary-light) !important;
}
</style>
