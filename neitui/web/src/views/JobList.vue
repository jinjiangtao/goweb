<template>
  <div class="job-list">
    <van-nav-bar title="内推招聘" />
    
    <div class="search-filter">
      <van-search v-model="search" placeholder="搜索职位" @search="onSearch" />
      <van-dropdown-menu>
        <van-dropdown-item v-model="salaryRange" :options="salaryOptions" @change="onFilterChange" />
        <van-dropdown-item v-model="location" :options="locationOptions" @change="onFilterChange" />
      </van-dropdown-menu>
    </div>

    <div class="jobs-container">
      <van-loading v-if="loading" style="text-align: center; padding: 20px;" />
      <div v-else-if="jobs.length === 0" class="empty">
        <van-empty description="暂无职位" />
      </div>
      <div v-else>
        <div v-for="job in jobs" :key="job.id" class="job-card" @click="goToDetail(job.id)">
          <div class="job-title">{{ job.title }}</div>
          <div class="job-info">
            <span class="salary">{{ job.salaryRange || '面议' }}</span>
            <span class="location">{{ job.location || '不限' }}</span>
            <span class="favorite-count">
              <van-icon name="star-o" /> {{ job.favoriteCount || 0 }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <van-tabbar v-model="activeTab" @change="onTabChange">
      <van-tabbar-item icon="apps-o" :to="{ name: 'JobList' }">职位</van-tabbar-item>
      <van-tabbar-item icon="star-o" :to="{ name: 'Favorites' }">收藏</van-tabbar-item>
      <van-tabbar-item icon="notes-o" :to="{ name: 'MySubmissions' }">投递</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import api from '@/api'

const router = useRouter()
const search = ref('')
const salaryRange = ref('')
const location = ref('')
const loading = ref(false)
const jobs = ref([])
const activeTab = ref(0)

const salaryOptions = [
  { text: '全部薪资', value: '' },
  { text: '5k-10k', value: '5k-10k' },
  { text: '10k-20k', value: '10k-20k' },
  { text: '20k以上', value: '20k以上' }
]

const locationOptions = [
  { text: '全部地点', value: '' },
  { text: '北京', value: '北京' },
  { text: '上海', value: '上海' },
  { text: '深圳', value: '深圳' },
  { text: '广州', value: '广州' },
  { text: '杭州', value: '杭州' }
]

const fetchJobs = async () => {
  loading.value = true
  try {
    const params = {}
    if (search.value) params.search = search.value
    if (salaryRange.value) params.salary_range = salaryRange.value
    if (location.value) params.location = location.value
    
    const res = await api.get('/public/jobs', { params })
    jobs.value = res.data || []
  } catch (error) {
    showToast('获取职位失败')
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  fetchJobs()
}

const onFilterChange = () => {
  fetchJobs()
}

const goToDetail = (id) => {
  router.push(`/job/${id}`)
}

const onTabChange = () => {
  // 当点击标签栏会触发路由跳转
}

onMounted(() => {
  fetchJobs()
})
</script>

<style scoped>
.job-list {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.search-filter {
  background-color: white;
  padding-bottom: 10px;
}

.jobs-container {
  padding: 10px;
}

.job-card {
  background-color: white;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.job-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 10px;
}

.job-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.salary {
  color: #ff6034;
  font-size: 14px;
}

.location {
  color: #666;
  font-size: 14px;
}

.favorite-count {
  color: #999;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 2px;
}

.empty {
  padding: 50px 0;
}
</style>
