<template>
  <div class="page-container" v-if="loaded">
    <el-card class="result-header-card" shadow="never">
      <div class="result-header">
        <div class="score-circle" :class="scoreClass">
          <div class="score-value">{{ result.Score?.toFixed(1) }}<span class="unit">分</span></div>
          <div class="score-label">综合得分</div>
        </div>
        <div class="result-stats">
          <h2>答题结果</h2>
          <el-descriptions :column="2" border size="large">
            <el-descriptions-item label="总题数">{{ result.TotalCount }} 题</el-descriptions-item>
            <el-descriptions-item label="答对">{{ result.CorrectCount }} 题</el-descriptions-item>
            <el-descriptions-item label="答错">
              <span style="color: #F56C6C">{{ result.TotalCount - result.CorrectCount }} 题</span>
            </el-descriptions-item>
            <el-descriptions-item label="正确率">
              <el-progress
                :percentage="Math.round(result.CorrectCount / result.TotalCount * 100)"
                :color="progressColor"
                :stroke-width="10"
              />
            </el-descriptions-item>
          </el-descriptions>
          <div class="result-actions">
            <el-button size="large" @click="router.push('/practice')">
              <el-icon><Refresh /></el-icon>
              再来一组
            </el-button>
            <el-button size="large" type="warning" @click="router.push('/wrong')">
              <el-icon><Warning /></el-icon>
              查看错题
            </el-button>
            <el-button size="large" type="primary" @click="router.push('/stats')">
              <el-icon><DataAnalysis /></el-icon>
              学习统计
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <el-radio-group v-model="filterMode" size="large">
          <el-radio-button value="all">全部题目</el-radio-button>
          <el-radio-button value="wrong">仅错题</el-radio-button>
          <el-radio-button value="correct">仅答对</el-radio-button>
        </el-radio-group>
        <div class="filter-stats">
          <el-tag type="success" effect="plain">正确 {{ correctCount }} 题</el-tag>
          <el-tag type="danger" effect="plain">错误 {{ wrongCount }} 题</el-tag>
        </div>
      </div>
    </el-card>

    <div class="question-list">
      <div v-for="(item, idx) in filteredDetails" :key="item.question_id" class="question-wrap">
        <QuestionCard
          :question="item.question"
          :index="getOriginalIndex(item) + 1"
          :total="result.TotalCount"
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
    <p>正在加载答题结果...</p>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import QuestionCard from '@/components/QuestionCard.vue'
import { getExamRecord } from '@/api'

const router = useRouter()
const route = useRoute()
const loaded = ref(false)
const filterMode = ref('all')

const record = ref({})
const result = reactive({ TotalCount: 0, CorrectCount: 0, Score: 0 })
const details = ref([])

const correctCount = computed(() => details.value.filter(d => d.is_correct).length)
const wrongCount = computed(() => details.value.filter(d => !d.is_correct).length)

const filteredDetails = computed(() => {
  if (filterMode.value === 'wrong') return details.value.filter(d => !d.is_correct)
  if (filterMode.value === 'correct') return details.value.filter(d => d.is_correct)
  return details.value
})

const scoreClass = computed(() => {
  const s = result.Score
  if (s >= 80) return 'score-good'
  if (s >= 60) return 'score-mid'
  return 'score-low'
})

const progressColor = computed(() => {
  const s = result.Score
  if (s >= 80) return '#67C23A'
  if (s >= 60) return '#E6A23C'
  return '#F56C6C'
})

const getOriginalIndex = (item) => {
  return details.value.findIndex(d => d.question_id === item.question_id)
}

onMounted(async () => {
  const id = route.params.recordId
  if (!id) {
    router.push('/')
    return
  }
  try {
    const res = await getExamRecord(id)
    record.value = res.record
    result.TotalCount = res.record.total_count
    result.CorrectCount = res.record.correct_count
    result.Score = res.record.score
    details.value = res.details
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

.result-header-card {
  border-radius: 12px;
  margin-bottom: 20px;
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
}

.result-header {
  display: flex;
  align-items: center;
  gap: 48px;
}

.score-circle {
  width: 160px;
  height: 160px;
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
}

.score-good { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }
.score-mid { background: linear-gradient(135deg, #fa709a 0%, #fee140 100%); }
.score-low { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }

.score-value {
  font-size: 44px;
  font-weight: 700;
  line-height: 1;
}

.score-value .unit {
  font-size: 16px;
  font-weight: 400;
  margin-left: 2px;
}

.score-label {
  font-size: 14px;
  opacity: 0.9;
  margin-top: 6px;
}

.result-stats {
  flex: 1;
}

.result-stats h2 {
  margin-bottom: 20px;
  color: #303133;
  font-size: 22px;
}

.result-actions {
  margin-top: 24px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-card {
  border-radius: 12px;
  margin-bottom: 20px;
}

.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-stats {
  display: flex;
  gap: 10px;
}

.question-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-bottom: 24px;
}
</style>
