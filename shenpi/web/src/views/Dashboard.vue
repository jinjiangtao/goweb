<script setup>
import { ref, onMounted } from 'vue'
import { getApprovalStats } from '@/api/approval'
import { useRouter } from 'vue-router'

const router = useRouter()

const stats = ref([
  { title: '待我审批', value: 0, icon: 'Clock', color: '#409eff', path: '/approval/todo' },
  { title: '我发起的', value: 0, icon: 'Edit', color: '#67c23a', path: '/approval' },
  { title: '已通过', value: 0, icon: 'Check', color: '#e6a23c', path: '/approval' },
  { title: '流程模板', value: 0, icon: 'Files', color: '#f56c6c', path: '/template' }
])

const recentApprovals = ref([])
const myRequests = ref([])

const statusMap = {
  pending: '审批中',
  approved: '已通过',
  rejected: '已驳回',
  revoked: '已撤回'
}

const typeMap = {
  leave: '请假申请',
  expense: '报销申请',
  general: '通用申请'
}

function handleStatClick(item) {
  if (item.path) {
    router.push(item.path)
  }
}

async function loadStats() {
  try {
    const res = await getApprovalStats()
    if (res.code === 200) {
      const data = res.data
      stats.value[0].value = data.todo_count || 0
      stats.value[1].value = data.my_count || 0
      stats.value[2].value = data.approved_count || 0
      stats.value[3].value = data.template_count || 0
      
      recentApprovals.value = (data.recent_approvals || []).map(item => ({
        title: item.title,
        type: typeMap[item.type] || item.type,
        status: statusMap[item.status] || item.status,
        time: item.created_at
      }))
      
      myRequests.value = (data.my_requests || []).map(item => ({
        title: item.title,
        status: statusMap[item.status] || item.status,
        createTime: item.created_at
      }))
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(() => {
  loadStats()
})
</script>

<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6" v-for="(item, index) in stats" :key="index">
        <el-card class="stat-card" @click="handleStatClick(item)">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-value" :style="{ color: item.color }">{{ item.value }}</div>
              <div class="stat-title">{{ item.title }}</div>
            </div>
            <div class="stat-icon" :style="{ backgroundColor: item.color + '20' }">
              <el-icon :size="32" :color="item.color">
                <component :is="item.icon" />
              </el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>待我审批</span>
          </template>
          <el-table :data="recentApprovals" empty-text="暂无数据" style="cursor: pointer">
            <el-table-column prop="title" label="标题" />
            <el-table-column prop="type" label="类型" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag size="small" :type="row.status === '审批中' ? 'warning' : 'success'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="时间" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>我的申请</span>
          </template>
          <el-table :data="myRequests" empty-text="暂无数据">
            <el-table-column prop="title" label="标题" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag size="small" :type="{
                  '审批中': 'warning',
                  '已通过': 'success',
                  '已驳回': 'danger',
                  '已撤回': 'info'
                }[row.status]">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
.stat-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .stat-content {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .stat-info {
      .stat-value {
        font-size: 32px;
        font-weight: bold;
        margin-bottom: 8px;
      }

      .stat-title {
        color: #909399;
        font-size: 14px;
      }
    }

    .stat-icon {
      width: 64px;
      height: 64px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>
