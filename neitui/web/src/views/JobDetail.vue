<template>
  <div class="job-detail">
    <van-nav-bar title="职位详情" left-arrow @click-left="router.back()" />
    
    <div v-if="loading" class="loading">
      <van-loading />
    </div>
    
    <div v-else-if="job" class="content">
      <div class="job-header">
        <h1>{{ job.title }}</h1>
        <div class="job-meta">
          <span class="salary">{{ job.salaryRange || '面议' }}</span>
          <span class="location">{{ job.location || '不限' }}</span>
        </div>
      </div>
      
      <div class="job-section">
        <h3>职位要求</h3>
        <p>{{ job.requirement }}</p>
      </div>
    </div>

    <div class="footer">
      <van-button type="primary" block @click="showApplyDialog = true">立即投递</van-button>
    </div>

    <van-dialog
      v-model:show="showApplyDialog"
      title="投递简历"
      show-cancel-button
      @confirm="submitApplication"
    >
      <div class="apply-form">
        <van-field
          v-model="form.name"
          label="姓名"
          placeholder="请输入姓名"
          required
        />
        <van-field
          v-model="form.phone"
          label="手机号"
          type="tel"
          placeholder="请输入手机号"
          required
        />
        <van-uploader
          v-model="form.resumeFile"
          :max-count="1"
          accept=".pdf,.doc,.docx"
          :before-read="beforeRead"
        />
      </div>
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import api from '@/api'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const job = ref(null)
const showApplyDialog = ref(false)
const form = ref({
  name: '',
  phone: '',
  resumeFile: null
})

const fetchJobDetail = async () => {
  loading.value = true
  try {
    const res = await api.get(`/public/jobs/${route.params.id}`)
    job.value = res
  } catch (error) {
    showToast('获取职位详情失败')
  } finally {
    loading.value = false
  }
}

const beforeRead = (file) => {
  const isValid = file.type === 'application/pdf' || 
                  file.type === 'application/msword' || 
                  file.type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
  if (!isValid) {
    showToast('请上传PDF或Word文件')
    return false
  }
  if (file.size > 5 * 1024 * 1024) {
    showToast('文件大小不能超过5MB')
    return false
  }
  return true
}

const submitApplication = async () => {
  if (!form.value.name || !form.value.phone || !form.value.resumeFile) {
    showToast('请填写完整信息')
    return false
  }
  
  if (!/^1[3-9]\d{9}$/.test(form.value.phone)) {
    showToast('请输入正确的手机号')
    return false
  }
  
  try {
    const formData = new FormData()
    formData.append('job_id', route.params.id)
    formData.append('name', form.value.name)
    formData.append('phone', form.value.phone)
    formData.append('resume_file', form.value.resumeFile.file)
    
    await api.post('/public/submit', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    showToast('投递成功，HR会尽快联系您')
    showApplyDialog.value = false
    form.value = { name: '', phone: '', resumeFile: null }
  } catch (error) {
    if (error.response?.data?.error) {
      showToast(error.response.data.error)
    } else {
      showToast('投递失败，请稍后重试')
    }
  }
}

onMounted(() => {
  fetchJobDetail()
})
</script>

<style scoped>
.job-detail {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 70px;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50vh;
}

.content {
  padding: 15px;
}

.job-header {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 15px;
}

.job-header h1 {
  font-size: 20px;
  color: #333;
  margin-bottom: 15px;
}

.job-meta {
  display: flex;
  gap: 15px;
}

.salary {
  color: #ff6034;
  font-size: 16px;
  font-weight: 600;
}

.location {
  color: #666;
  font-size: 14px;
}

.job-section {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
}

.job-section h3 {
  font-size: 16px;
  color: #333;
  margin-bottom: 15px;
}

.job-section p {
  color: #666;
  line-height: 1.6;
}

.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 10px 15px;
  background-color: white;
  box-shadow: 0 -1px 3px rgba(0, 0, 0, 0.1);
}

.apply-form {
  padding-top: 10px;
}
</style>
