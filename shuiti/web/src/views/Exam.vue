<template>
  <div class="exam-page">
    <div class="page-container" v-if="questions.length">
      <div class="exam-header-card">
        <div class="exam-info">
          <div class="mode-tag">
            <el-tag :type="modeTagColor" size="large" effect="dark">{{ modeName }}</el-tag>
            <span style="margin-left: 12px; color: #606266; font-size: 14px;">
              共 {{ questions.length }} 题 · 已答 {{ answeredCount }} 题
            </span>
          </div>
        </div>
        <div class="timer">
          <el-icon :size="18" color="#E6A23C"><Timer /></el-icon>
          <span>用时：{{ formatTime(elapsedTime) }}</span>
        </div>
        <div class="header-actions">
          <el-popconfirm title="确定要放弃本次答题吗？未提交的答案会被保留" @confirm="quitExam">
            <template #reference>
              <el-button>
                <el-icon><SwitchButton /></el-icon>
                退出
              </el-button>
            </template>
          </el-popconfirm>
          <el-button type="primary" :disabled="answeredCount === 0" @click="confirmSubmit">
            <el-icon><Check /></el-icon>
            提交试卷
          </el-button>
        </div>
      </div>

      <el-row :gutter="20">
        <el-col :xs="24" :md="17">
          <QuestionCard
            v-if="currentQuestion"
            :question="currentQuestion"
            :index="currentIndex + 1"
            :total="questions.length"
            v-model="currentAnswer"
          />

          <div class="nav-card">
            <div class="progress-bar-wrap">
              <el-progress
                :percentage="Math.round((currentIndex + 1) / questions.length * 100)"
                :stroke-width="8"
                :show-text="false"
                color="#409EFF"
              />
            </div>
            <div class="nav-buttons">
              <el-button size="large" :disabled="currentIndex === 0" @click="prevQuestion">
                <el-icon><ArrowLeft /></el-icon>
                上一题
              </el-button>
              <div class="progress-text">
                {{ currentIndex + 1 }} / {{ questions.length }}
              </div>
              <el-button size="large" type="primary" :disabled="currentIndex === questions.length - 1" @click="nextQuestion">
                下一题
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
          </div>
        </el-col>

        <el-col :xs="24" :md="7">
          <el-card class="answer-sheet" shadow="never">
            <template #header>
              <div class="sheet-header">
                <el-icon color="#409EFF"><Grid /></el-icon>
                <span>答题卡</span>
                <span style="margin-left: auto; font-size: 13px; color: #909399; font-weight: normal;">
                  {{ answeredCount }}/{{ questions.length }}
                </span>
              </div>
            </template>
            <div class="answer-grid">
              <div
                v-for="(q, idx) in questions"
                :key="q.id"
                class="grid-item"
                :class="getGridClass(idx)"
                @click="goToQuestion(idx)"
              >
                {{ idx + 1 }}
              </div>
            </div>
            <div class="legend">
              <div class="legend-item">
                <div class="legend-box current"></div>
                <span>当前</span>
              </div>
              <div class="legend-item">
                <div class="legend-box answered"></div>
                <span>已答</span>
              </div>
              <div class="legend-item">
                <div class="legend-box"></div>
                <span>未答</span>
              </div>
            </div>
          </el-card>

          <el-card class="tip-card" shadow="never">
            <div class="tip-title">
              <el-icon color="#E6A23C"><LightBulb /></el-icon>
              <span>温馨提示</span>
            </div>
            <ul class="tip-list">
              <li>答题进度已实时缓存，刷新也不会丢失</li>
              <li>多选题需要选择所有正确选项才算对</li>
              <li>提交后系统将自动判分并归集错题</li>
            </ul>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <el-empty v-else description="暂无题目，请返回选择组卷参数" style="padding-top: 80px;">
      <el-button type="primary" @click="router.push('/practice')">返回选择</el-button>
    </el-empty>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import QuestionCard from '@/components/QuestionCard.vue'
import { submitExam, getWrongPractice } from '@/api'
import { useUserStore } from '@/stores/user'
import { useExamStore } from '@/stores/exam'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const examStore = useExamStore()

const questions = ref([])
const answers = ref({})
const currentIndex = ref(0)
const elapsedTime = ref(0)
let timer = null

const mode = computed(() => route.params.mode || 'random')
const modeName = computed(() => ({ order: '顺序刷题', random: '随机组卷', smart: '智能组卷', wrong: '错题练习' }[mode.value] || '练习'))
const modeTagColor = computed(() => ({ order: 'primary', random: 'warning', smart: 'success', wrong: 'danger' }[mode.value] || 'info'))

const currentQuestion = computed(() => questions.value[currentIndex.value])
const currentId = computed(() => currentQuestion.value?.id)

const currentAnswer = computed({
  get: () => currentId.value ? (answers.value[currentId.value] || '') : '',
  set: (val) => {
    if (currentId.value) {
      answers.value[currentId.value] = val
      saveProgress()
    }
  }
})

const answeredCount = computed(() => {
  return Object.values(answers.value).filter(v => v && v.toString().trim()).length
})

const loadExam = async () => {
  const saved = examStore.progress
  if (saved && saved.questions && saved.questions.length > 0 && saved.mode === mode.value) {
    questions.value = saved.questions
    answers.value = saved.answers || {}
    currentIndex.value = saved.currentIndex || 0
    elapsedTime.value = saved.elapsedTime || 0
    return
  }

  if (mode.value === 'wrong') {
    const res = await getWrongPractice({ user_id: userStore.userId })
    questions.value = res.questions || []
  } else {
    const cached = sessionStorage.getItem('exam_data')
    if (cached) {
      const data = JSON.parse(cached)
      questions.value = data.questions || []
    }
  }
  answers.value = {}
  currentIndex.value = 0
  elapsedTime.value = 0
}

const saveProgress = () => {
  examStore.saveProgress({
    mode: mode.value,
    questions: questions.value,
    answers: answers.value,
    currentIndex: currentIndex.value,
    elapsedTime: elapsedTime.value
  })
}

const goToQuestion = (idx) => {
  currentIndex.value = idx
  saveProgress()
}

const prevQuestion = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--
    saveProgress()
  }
}

const nextQuestion = () => {
  if (currentIndex.value < questions.value.length - 1) {
    currentIndex.value++
    saveProgress()
  }
}

const getGridClass = (idx) => {
  const classes = []
  if (idx === currentIndex.value) classes.push('current')
  const q = questions.value[idx]
  if (q && answers.value[q.id] && answers.value[q.id].toString().trim()) {
    classes.push('answered')
  }
  return classes
}

const confirmSubmit = () => {
  const unanswered = questions.value.length - answeredCount.value
  const msg = unanswered > 0
    ? `还有 ${unanswered} 道题未作答，确定要提交吗？`
    : `已完成所有 ${questions.value.length} 道题，确定提交吗？`

  ElMessageBox.confirm(msg, '提交确认', {
    type: unanswered > 0 ? 'warning' : 'success',
    confirmButtonText: '确定提交',
    cancelButtonText: '继续答题'
  }).then(() => handleSubmit()).catch(() => {})
}

const handleSubmit = async () => {
  const answersArr = questions.value.map(q => ({
    question_id: q.id,
    answer: answers.value[q.id] || ''
  }))

  const payload = {
    user_id: userStore.userId,
    exam_name: modeName.value,
    mode: mode.value,
    answers: answersArr
  }

  try {
    const res = await submitExam(payload)
    examStore.clearProgress()
    sessionStorage.removeItem('exam_data')
    ElMessage.success('提交成功，正在跳转...')
    router.push(`/result/${res.record_id}`)
  } catch {}
}

const quitExam = () => {
  saveProgress()
  router.push('/practice')
}

const formatTime = (seconds) => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  const pad = (n) => n.toString().padStart(2, '0')
  if (h > 0) return `${pad(h)}:${pad(m)}:${pad(s)}`
  return `${pad(m)}:${pad(s)}`
}

onMounted(async () => {
  await loadExam()
  timer = setInterval(() => {
    elapsedTime.value++
    saveProgress()
  }, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.exam-page {
  min-height: 100%;
}

.page-container {
  max-width: 1400px;
  margin: 0 auto;
}

.exam-header-card {
  background: #fff;
  border-radius: 12px;
  padding: 18px 24px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.exam-info {
  flex: 1;
}

.timer {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-right: 24px;
  padding: 6px 14px;
  background: #fffbeb;
  border-radius: 20px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.nav-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px 24px;
  margin-top: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.progress-bar-wrap {
  margin-bottom: 16px;
}

.nav-buttons {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.progress-text {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.answer-sheet {
  border-radius: 12px;
  margin-bottom: 20px;
}

.sheet-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #303133;
}

.answer-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 10px;
  margin-bottom: 20px;
}

.grid-item {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: #f5f7fa;
  color: #606266;
  font-weight: 500;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.15s ease;
  font-size: 14px;
}

.grid-item:hover {
  border-color: #c6e2ff;
}

.grid-item.current {
  border-color: #409EFF;
  background: #ecf5ff;
  color: #409EFF;
  font-weight: 600;
}

.grid-item.answered {
  background: linear-gradient(135deg, #67C23A 0%, #95d475 100%);
  color: #fff;
}

.grid-item.current.answered {
  background: linear-gradient(135deg, #409EFF 0%, #79bbff 100%);
  color: #fff;
  border-color: #409EFF;
}

.legend {
  display: flex;
  justify-content: space-around;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #606266;
}

.legend-box {
  width: 18px;
  height: 18px;
  border-radius: 4px;
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
}

.legend-box.current {
  background: #ecf5ff;
  border-color: #409EFF;
}

.legend-box.answered {
  background: #67C23A;
  border-color: #67C23A;
}

.tip-card {
  border-radius: 12px;
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
  border: none;
}

.tip-card :deep(.el-card__body) {
  padding: 20px;
}

.tip-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: #92400e;
  margin-bottom: 12px;
}

.tip-list {
  list-style: none;
  padding: 0;
}

.tip-list li {
  padding: 6px 0;
  color: #78350f;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.tip-list li::before {
  content: '•';
  color: #e6a23c;
  font-weight: bold;
}
</style>
