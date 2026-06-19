<template>
  <div class="page-container">
    <el-card class="main-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-icon :size="22" color="#409EFF"><Document /></el-icon>
            <span>答题记录</span>
            <el-tag type="info" effect="plain" style="margin-left: 12px;">共 {{ total }} 条记录</el-tag>
          </div>
        </div>
      </template>

      <el-table
        :data="list"
        v-loading="loading"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="exam_name" label="练习名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="mode" label="模式" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="modeColor(row.mode)" size="small" effect="plain">{{ modeName(row.mode) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_count" label="总题数" width="90" align="center" />
        <el-table-column label="答对" width="90" align="center">
          <template #default="{ row }">
            <span style="color: #67C23A; font-weight: 500;">{{ row.correct_count }}</span>
          </template>
        </el-table-column>
        <el-table-column label="正确率" width="130" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="Math.round(row.correct_count / row.total_count * 100)"
              :stroke-width="8"
              :color="progressColor(row.score)"
            />
          </template>
        </el-table-column>
        <el-table-column label="得分" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="scoreColor(row.score)" size="small" effect="dark">{{ row.score?.toFixed(1) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="190" align="center" />
        <el-table-column label="操作" width="120" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="router.push(`/record/${row.id}`)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap" v-if="total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>

      <el-empty v-if="!loading && total === 0" description="暂无答题记录" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getExamRecords } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const list = ref([])
const total = ref(0)

const pagination = reactive({
  page: 1,
  page_size: 20
})

const modeName = (m) => ({ order: '顺序刷题', random: '随机组卷', smart: '智能组卷', wrong: '错题练习' }[m] || '练习')
const modeColor = (m) => ({ order: 'primary', random: 'warning', smart: 'success', wrong: 'danger' }[m] || 'info')
const scoreColor = (s) => s >= 80 ? 'success' : s >= 60 ? 'warning' : 'danger'
const progressColor = (s) => s >= 80 ? '#67C23A' : s >= 60 ? '#E6A23C' : '#F56C6C'

const loadData = async () => {
  loading.value = true
  try {
    const res = await getExamRecords({
      user_id: userStore.userId,
      page: pagination.page,
      page_size: pagination.page_size
    })
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.page-container {
  max-width: 1400px;
  margin: 0 auto;
}

.main-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 17px;
  font-weight: 600;
  color: #303133;
}

.pagination-wrap {
  display: flex;
  justify-content: center;
  padding-top: 24px;
}
</style>
