
&lt;template&gt;
  &lt;div class="my-bookings"&gt;
    &lt;div class="header"&gt;
      &lt;el-button text @click="$router.back()" class="back-btn"&gt;
        &lt;el-icon&gt;&lt;ArrowLeft /&gt;&lt;/el-icon&gt;
      &lt;/el-button&gt;
      &lt;h1&gt;我的预订&lt;/h1&gt;
      &lt;div style="width: 32px;"&gt;&lt;/div&gt;
    &lt;/div&gt;

    &lt;div class="content"&gt;
      &lt;div class="booking-list" v-loading="loading"&gt;
        &lt;el-empty v-if="!loading &amp;&amp; bookings.length === 0" description="暂无预订记录" /&gt;

        &lt;div class="booking-card" v-for="booking in bookings" :key="booking.id"&gt;
          &lt;div class="booking-header"&gt;
            &lt;div class="room-name"&gt;{{ booking.room?.name || '会议室' }}&lt;/div&gt;
            &lt;el-tag :type="booking.status === 1 ? 'success' : 'info'" size="small"&gt;
              {{ booking.status === 1 ? '已预订' : '已取消' }}
            &lt;/el-tag&gt;
          &lt;/div&gt;
          &lt;div class="booking-info"&gt;
            &lt;div class="info-item"&gt;
              &lt;el-icon&gt;&lt;Calendar /&gt;&lt;/el-icon&gt;
              &lt;span&gt;{{ booking.date }}&lt;/span&gt;
            &lt;/div&gt;
            &lt;div class="info-item"&gt;
              &lt;el-icon&gt;&lt;Clock /&gt;&lt;/el-icon&gt;
              &lt;span&gt;{{ booking.start_time }} - {{ booking.end_time }}&lt;/span&gt;
            &lt;/div&gt;
            &lt;div class="info-item" v-if="booking.purpose"&gt;
              &lt;el-icon&gt;&lt;Document /&gt;&lt;/el-icon&gt;
              &lt;span&gt;{{ booking.purpose }}&lt;/span&gt;
            &lt;/div&gt;
          &lt;/div&gt;
          &lt;div class="booking-actions" v-if="booking.status === 1"&gt;
            &lt;el-button type="danger" size="small" @click="cancelBooking(booking)"&gt;取消预订&lt;/el-button&gt;
          &lt;/div&gt;
        &lt;/div&gt;
      &lt;/div&gt;
    &lt;/div&gt;

    &lt;div class="footer-btn"&gt;
      &lt;el-button type="primary" class="back-home-btn" @click="$router.push('/')"&gt;
        &lt;el-icon&gt;&lt;House /&gt;&lt;/el-icon&gt;
        返回预订
      &lt;/el-button&gt;
    &lt;/div&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Calendar, Clock, Document, House } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const bookings = ref([])

onMounted(() =&gt; {
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
&lt;/script&gt;

&lt;style scoped&gt;
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
&lt;/style&gt;

