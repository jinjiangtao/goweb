<template>
  <div class="favorites-page">
    <van-nav-bar title="我的收藏" />

    <div class="content">
      <van-loading v-if="loading" style="text-align: center; padding: 20px;" />
      <div v-else-if="favorites.length === 0" class="empty">
        <van-empty description="暂无收藏" />
      </div>
      <div v-else>
        <div v-for="job in favorites" :key="job.id" class="job-card" @click="goToDetail(job.id)">
          <div class="job-title">{{ job.title }}</div>
          <div class="job-info">
            <span class="salary">{{ job.salaryRange || '面议' }}</span>
            <span class="location">{{ job.location || '不限' }}</span>
          </div>
        </div>
      </div>
    </div>

    <van-tabbar v-model="activeTab" @change="onTabChange">
      <van-tabbar-item icon="apps-o" :to="{ name: 'JobList' }">职位</van-tabbar-item>
      <van-tabbar-item icon="star-o" :to="{ name: 'Favorites' }">收藏</van-tabbar-item>
      <van-tabbar-item icon="notes-o" :to="{ name: 'MySubmissions' }">投递</van-tabbar-item>
    </van-tabbar-item>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import api from '@/api'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const favorites = ref([])
const activeTab = ref(1)

const fetchFavorites = async () => {
  const storedFavorites = JSON.parse(localStorage.getItem('favoriteJobs') || '[]')
  if (storedFavorites.length === 0) {
    favorites.value = []
    return
  }

  loading.value = true
  try {
    const res = await api.post('/public/jobs/favorites', { ids: storedFavorites })
    favorites.value = res || []
  } catch (error) {
    showToast('获取收藏失败')
  } finally {
    loading.value = false
  }
}

const goToDetail = (id) => {
  router.push(`/job/${id}`)
}

const onTabChange = () => {
  // 当点击标签栏会触发路由跳转
}

onMounted(() => {
  fetchFavorites()
})
</script>

<style scoped>
.favorites-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.content {
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
}

.salary {
  color: #ff6034;
  font-size: 14px;
}

.location {
  color: #666;
  font-size: 14px;
}

.empty {
  padding: 50px 0;
}
</style>
