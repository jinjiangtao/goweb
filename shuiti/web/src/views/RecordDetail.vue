<template>
  <div class="page-container" v-if="loaded">
    <el-card class="info-card" shadow="never">
      <el-descriptions :column="4" border title="答题记录详情">
        <el-descriptions-item label="练习名称">{{ record.exam_name }}</el-descriptions-item>
        <el-descriptions-item label="模式">
          <el-tag :type="modeColor(record.mode)" effect="plain">{{ modeName(record.mode) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ record.created_at }}</el-descriptions-item>
        <el-descriptions-item label="得分">
          <el-tag :type="scoreColor(record.score)" effect="dark" size="large">{{ record.score?.toFixed(1) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="总题数">{{ record.total_count }} 题</el-descriptions-item>
        <el-descriptions-item label="答对">
          <span style="color: #67C23A; font-weight: 600;">{{ record.correct_count }} 题</span>
        </el-descriptions-item>
        <el-descriptions-item label="答错">
          <span style="color: #F56C6C; font-weight: 600;">{{ record.total_count - record.correct_count }} 题</span>
        </el-descriptions-item>
        <el-descriptions-item label="正确率">
          <el-progress
            :percentage="Math.round(record.correct_count / record.total_count * 100)"
            :stroke-width="10"
            :color="progressColor(record.score)"
          />
        </el-descriptions-item>
      </el-descriptions>
      <div style="margin-top: 16px;">
        <el-radio-group v-model="filterMode" size="default">
          <el-radio-button value="all">全部题目</el-radio-button>
          <el-radio-button value="wrong">仅错题</el-radio-button>
          <el-radio-button value="correct">仅答对</el-radio-button>
        </el-radio-group>
      </div>
    </el-card>

    <div class="question-list">
      <div v-for="(item, idx) in filteredDetails" :key="item.question_id" class="question-wrap">
        <QuestionCard
          :question="item.question"
          :index="idx + 1"
          :total="filteredDetails.length"
          :model-value="item.user_answer"
          :show-result="true"
          :result-detail="item"
        />
      </div>
      <el-empty v-if="filteredDetails.length === 0" description="暂无题目" />
    </div>
  </div>

  <div v-else class="loading-wrap">
    <el-icon :size="40" class="is-loading"><Loading /></el-icon>
    <p>正在加载...</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import QuestionCard from '@/components/QuestionCard.vue'
import { getExamRecord } from '@/api'

const route = useRoute()
const router = useRouter()
const loaded = ref(false)
const record = ref({})
const details = ref([])
const filterMode = ref('all')

const filteredDetails = computed(() => {
  if (filterMode.value === 'wrong') return details.value.filter(d => !d.is_correct)
  if (filterMode.value === 'correct') return details.value.filter(d => d.is_correct)
  return details.value
})

const modeName = (m) => ({ order: '顺序刷题', random: '随机组卷', smart: '智能组卷', wrong: '错题练习' }[m] || '练习')
const modeColor = (m) => ({ order: 'primary', random: 'warning', smart: 'success', wrong: 'danger' }[m] || 'info')
const scoreColor = (s) => s >= 80 ? 'success' : s >= 60 ? 'warning' : 'danger'
const progressColor = (s) => s >= 80 ? '#67C23A' : s >= 60 ? '#E6A23C' : '#F56C6C'

onMounted(async () => {
  const id = route.params.id
  if (!id) {
    router.push('/records')
    return
  }
  try {
    const res = await getExamRecord(id)
    record.value = res.record
    details.value = res.details || []
  } finally {
    loaded.value = true
  }
})
</script>

<style scoped>
.page-container {
  max-width: 1100px;
  margin: 0 auto;
}

.loading-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  color: #909399;
  gap: 16px;
}

.info-card {
  border-radius: 12px;
  margin-bottom: 20px;
}

.question-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-bottom: 24px;
}
</style>
