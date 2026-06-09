<template>
  <div class="my-submissions">
    <van-nav-bar title="我的投递" />
    
    <div class="phone-input">
      <van-field
        v-model="phone"
        type="tel"
        placeholder="请输入手机号查询"
        clearable
      >
        <template #button>
          <van-button size="small" type="primary" @click="fetchSubmissions">查询</van-button>
        </template>
      </van-field>
    </div>

    <div class="submissions-container">
      <van-loading v-if="loading" style="text-align: center; padding: 20px;" />
      <div v-else-if="!phone" class="empty">
        <van-empty description="请输入手机号查询投递记录" />
      </div>
      <div v-else-if="submissions.length === 0" class="empty">
        <van-empty description="暂无投递记录" />
      </div>
      <div v-else>
        <div v-for="sub in submissions" :key="sub.id" class="sub-card">
          <div class="sub-header">
            <span class="job-title">{{ sub.jobTitle }}</span>
            <span :class="['status', getStatusClass(sub.status)]">{{ getStatusText(sub.status) }}</span>
          </div>
          <div class="sub-info">
            <div class="info-item">
              <span class="label">投递人：</span>
              <span>{{ sub.name }}</span>
            </div>
            <div class="info-item">
              <span class="label">投递时间：</span>
              <span>{{ formatDate(sub.createdAt) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { showToast } from 'vant'
import api from '@/api'

const phone = ref('')
const loading = ref(false)
const submissions = ref([])

const getStatusText = (status) => {
  const map = {
    'screening': '初筛中',
    'interviewing': '已面试',
    'offer': '已发offer',
    'hired': '已入职',
    'rejected': '已淘汰'
  }
  return map[status] || status
}

const getStatusClass = (status) => {
  const map = {
    'screening': 'pending',
    'interviewing': 'processing',
    'offer': 'success',
    'hired': 'success',
    'rejected': 'danger'
  }
  return map[status] || ''
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

const fetchSubmissions = async () => {
  if (!phone.value) {
    showToast('请输入手机号')
    return
  }
  
  if (!/^1[3-9]\d{9}$/.test(phone.value)) {
    showToast('请输入正确的手机号')
    return
  }
  
  loading.value = true
  try {
    const res = await api.get('/public/submissions', { params: { phone: phone.value } })
    submissions.value = res
  } catch (error) {
    showToast('查询失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.my-submissions {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.phone-input {
  background-color: white;
  padding: 10px 15px;
  margin-bottom: 10px;
}

.submissions-container {
  padding: 0 15px 15px;
}

.sub-card {
  background-color: white;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
}

.sub-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.job-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.status {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
}

.status.pending {
  background-color: #fff7e6;
  color: #ff976a;
}

.status.processing {
  background-color: #e6f7ff;
  color: #1989fa;
}

.status.success {
  background-color: #f0f9eb;
  color: #07c160;
}

.status.danger {
  background-color: #fff0f0;
  color: #ee0a24;
}

.sub-info {
  color: #666;
  font-size: 14px;
}

.info-item {
  margin-bottom: 5px;
}

.label {
  color: #999;
}

.empty {
  padding: 50px 0;
}
</style>
