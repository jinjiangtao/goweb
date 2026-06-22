<template>
  <div class="profile">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="never" class="user-card">
          <div class="avatar-section">
            <el-avatar :size="100" style="background: #409EFF; font-size: 40px;">
              {{ user?.real_name?.charAt(0) || user?.username?.charAt(0) }}
            </el-avatar>
            <h2>{{ user?.real_name || user?.username }}</h2>
            <el-tag :type="roleType" size="large">{{ roleText }}</el-tag>
          </div>
          <el-descriptions :column="1" class="desc" size="small">
            <el-descriptions-item label="用户名">{{ user?.username }}</el-descriptions-item>
            <el-descriptions-item label="电话">{{ user?.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ user?.email || '-' }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ user?.area || '-' }}</el-descriptions-item>
            <el-descriptions-item label="技能">{{ user?.skills || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="never">
          <template #header><span>我的工单统计</span></template>
          <el-row :gutter="20">
            <el-col :span="6">
              <el-statistic title="总工单" :value="myStats.total" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="待处理" :value="myStats.pending">
                <template #suffix><el-icon style="color: #E6A23C;"><Clock /></el-icon></template>
              </el-statistic>
            </el-col>
            <el-col :span="6">
              <el-statistic title="处理中" :value="myStats.processing">
                <template #suffix><el-icon style="color: #409EFF;"><Loading /></el-icon></template>
              </el-statistic>
            </el-col>
            <el-col :span="6">
              <el-statistic title="已完成" :value="myStats.completed">
                <template #suffix><el-icon style="color: #67C23A;"><CircleCheck /></el-icon></template>
              </el-statistic>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getOrders } from '@/api'

const userStore = useUserStore()
const user = ref(userStore.user)

const myStats = reactive({ total: 0, pending: 0, processing: 0, completed: 0 })

const roleText = computed(() => {
  const role = user.value?.role
  if (role === 'admin') return '管理员'
  if (role === 'technician') return '运维人员'
  return '普通用户'
})

const roleType = computed(() => {
  const role = user.value?.role
  if (role === 'admin') return 'danger'
  if (role === 'technician') return 'warning'
  return 'success'
})

const loadStats = async () => {
  try {
    const uid = userStore.user?.id
    if (!uid) return
    const params = userStore.isTechnician ? { technician_id: uid, page_size: 1000 } : { user_id: uid, page_size: 1000 }
    const res = await getOrders(params)
    myStats.total = res.total
    myStats.pending = res.list.filter(o => o.status === 'pending' || o.status === 'assigned').length
    myStats.processing = res.list.filter(o => o.status === 'processing').length
    myStats.completed = res.list.filter(o => o.status === 'completed').length
  } catch (e) {}
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.user-card {
  text-align: center;
}

.avatar-section {
  padding: 20px 0;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 20px;
}

.avatar-section h2 {
  margin: 15px 0 10px;
  font-size: 20px;
  color: #333;
}

.desc {
  text-align: left;
}
</style>
