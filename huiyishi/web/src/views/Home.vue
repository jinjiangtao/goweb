<template>
  <div class="home">
    <div class="banner">
      <div class="banner-top">
        <h1>会议室预订</h1>
        <el-button text class="logout-btn" @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          退出
        </el-button>
      </div>
      <p>高效协作，从订会议室开始</p>
      <p class="user-info">欢迎你，{{ user?.real_name }}</p>
    </div>

    <div class="content">
      <div class="section-title">
        <el-icon><OfficeBuilding /></el-icon>
        <span>选择会议室</span>
      </div>

      <div class="room-list" v-loading="loading">
        <div class="room-card" v-for="room in rooms" :key="room.id">
          <div class="room-header">
            <h3>{{ room.name }}</h3>
            <div class="capacity">
              <el-icon><User /></el-icon>
              <span>可容纳 {{ room.capacity }} 人</span>
            </div>
          </div>
          <div class="devices">
            <div class="device-item" v-for="device in parseDevices(room.devices)" :key="device">
              <el-icon :size="16">
                <component :is="getDeviceIcon(device)" />
              </el-icon>
              <span>{{ device }}</span>
            </div>
          </div>
          <el-button type="primary" class="book-btn" @click="openBookingDrawer(room)">立即预订</el-button>
        </div>
      </div>
    </div>

    <div class="footer-btn">
      <el-button type="primary" class="my-bookings-btn" @click="$router.push('/my-bookings')">
        <el-icon><Document /></el-icon>
        我的预订
      </el-button>
    </div>

    <el-drawer
      v-model="drawerVisible"
      title="预订会议室"
      direction="btt"
      size="80%"
    >
      <div class="booking-form">
        <el-form :model="form" ref="formRef" :rules="rules" label-width="80px">
          <el-alert
            :title="'预订会议室：' + selectedRoom?.name"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px;"
          />
          <el-form-item label="日期" prop="date">
            <el-date-picker
              v-model="form.date"
              type="date"
              placeholder="选择日期"
              :disabled-date="disabledDate"
              value-format="YYYY-MM-DD"
              style="width: 100%;"
            />
          </el-form-item>
          <el-form-item label="开始时间" prop="startTime">
            <el-select v-model="form.startTime" placeholder="选择开始时间" style="width: 100%;">
              <el-option
                v-for="time in timeOptions"
                :key="time"
                :label="time"
                :value="time"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="结束时间" prop="endTime">
            <el-select v-model="form.endTime" placeholder="选择结束时间" style="width: 100%;">
              <el-option
                v-for="time in filteredEndTimes"
                :key="time"
                :label="time"
                :value="time"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="用途" prop="purpose">
            <el-input
              v-model="form.purpose"
              type="textarea"
              :rows="3"
              placeholder="请输入预订用途（选填）"
            />
          </el-form-item>
        </el-form>

        <div class="form-buttons">
          <el-button @click="drawerVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="submitBooking">提交预订</el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { OfficeBuilding, User, Document, Monitor, VideoCamera, DataBoard, Phone, SwitchButton } from '@element-plus/icons-vue'
import api, { setAuthToken } from '../api'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const drawerVisible = ref(false)
const rooms = ref([])
const selectedRoom = ref(null)
const formRef = ref(null)
const user = ref(null)

const form = reactive({
  date: '',
  startTime: '',
  endTime: '',
  purpose: ''
})

const rules = {
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  startTime: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  endTime: [{ required: true, message: '请选择结束时间', trigger: 'change' }]
}

const timeOptions = ref([])

onMounted(() => {
  generateTimeOptions()
  fetchRooms()
  getUserInfo()
})

function getUserInfo() {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
  }
}

function generateTimeOptions() {
  const options = []
  for (let hour = 9; hour <= 18; hour++) {
    options.push(`${hour.toString().padStart(2, '0')}:00`)
    if (hour < 18) {
      options.push(`${hour.toString().padStart(2, '0')}:30`)
    }
  }
  timeOptions.value = options
}

const filteredEndTimes = computed(() => {
  if (!form.startTime) return timeOptions.value
  const startIndex = timeOptions.value.indexOf(form.startTime)
  return timeOptions.value.slice(startIndex + 1)
})

function disabledDate(time) {
  return time.getTime() < Date.now() - 8.64e7
}

function parseDevices(devicesStr) {
  if (!devicesStr) return []
  try {
    const parsed = JSON.parse(devicesStr)
    if (Array.isArray(parsed)) {
      return parsed.map(d => {
        if (typeof d === 'string') return d
        if (d.name) return d.name
        if (d.label) return d.label
        return String(d)
      })
    }
    return []
  } catch {
    return devicesStr.split(',').map(d => d.trim()).filter(d => d)
  }
}

function getDeviceIcon(device) {
  const iconMap = {
    '投影仪': Monitor,
    '电视': VideoCamera,
    '白板': DataBoard,
    '会议电话': Phone
  }
  return iconMap[device] || Monitor
}

async function fetchRooms() {
  loading.value = true
  try {
    const res = await api.get('/rooms')
    rooms.value = res.data
  } catch (error) {
    ElMessage.error('获取会议室列表失败')
  } finally {
    loading.value = false
  }
}

function openBookingDrawer(room) {
  selectedRoom.value = room
  form.date = ''
  form.startTime = ''
  form.endTime = ''
  form.purpose = ''
  drawerVisible.value = true
}

async function submitBooking() {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    await api.post('/bookings', {
      room_id: selectedRoom.value.id,
      date: form.date,
      start_time: form.startTime,
      end_time: form.endTime,
      purpose: form.purpose
    })

    ElMessage.success('预订成功！')
    drawerVisible.value = false

    ElMessageBox.confirm(
      '预订成功！是否查看我的预订？',
      '提示',
      {
        confirmButtonText: '去查看',
        cancelButtonText: '留在首页',
        type: 'success'
      }
    ).then(() => {
      router.push('/my-bookings')
    }).catch(() => {})
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '预订失败')
  } finally {
    submitting.value = false
  }
}

function handleLogout() {
  ElMessageBox.confirm(
    '确定要退出登录吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setAuthToken(null)
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  padding-bottom: 80px;
}

.banner {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
  padding: 20px;
  color: white;
}

.banner-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.banner h1 {
  font-size: 24px;
  margin-bottom: 8px;
}

.banner p {
  font-size: 14px;
  opacity: 0.9;
}

.user-info {
  margin-top: 8px;
  font-weight: 500;
}

.logout-btn {
  color: white;
  padding: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.content {
  padding: 20px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 16px;
}

.room-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.room-card {
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 16px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.room-header {
  margin-bottom: 12px;
}

.room-header h3 {
  font-size: 18px;
  color: var(--text-color);
  margin-bottom: 8px;
}

.capacity {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--text-secondary);
  font-size: 14px;
}

.devices {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}

.device-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--text-secondary);
  font-size: 13px;
}

.book-btn {
  width: 100%;
  border-radius: 6px;
  background-color: var(--primary-color);
}

.book-btn:hover {
  background-color: var(--primary-light) !important;
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

.my-bookings-btn {
  width: 100%;
  border-radius: 6px;
  background-color: var(--primary-color);
}

.my-bookings-btn:hover {
  background-color: var(--primary-light) !important;
}

.booking-form {
  padding: 10px 0;
}

.form-buttons {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.form-buttons .el-button {
  flex: 1;
  border-radius: 6px;
}
</style>
