<template>
  <div class="page-container">
    <el-card class="main-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-icon :size="22" color="#409EFF"><EditPen /></el-icon>
            <span>选择刷题模式</span>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8">
          <div class="mode-card" :class="{ active: selectedMode === 'order' }" @click="selectedMode = 'order'">
            <div class="mode-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
              <el-icon :size="36"><Sort /></el-icon>
            </div>
            <h3>顺序刷题</h3>
            <p>按照题目录入顺序逐一练习，全面覆盖题库内容</p>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8">
          <div class="mode-card" :class="{ active: selectedMode === 'random' }" @click="selectedMode = 'random'">
            <div class="mode-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
              <el-icon :size="36"><Shuffle /></el-icon>
            </div>
            <h3>随机组卷</h3>
            <p>从题库中随机抽取题目进行练习，模拟真实考试</p>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8">
          <div class="mode-card" :class="{ active: selectedMode === 'smart' }" @click="selectedMode = 'smart'">
            <div class="mode-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
              <el-icon :size="36"><MagicStick /></el-icon>
            </div>
            <h3>智能组卷</h3>
            <p>均衡搭配各题型比例，全面考察知识掌握情况</p>
          </div>
        </el-col>
      </el-row>

      <el-divider>设置组卷参数</el-divider>

      <el-form :model="form" label-width="100px" class="config-form">
        <el-row :gutter="24">
          <el-col :xs="24" :md="12">
            <el-form-item label="题目数量">
              <el-slider
                v-model="form.count"
                :min="5"
                :max="sliderMax"
                :step="5"
                show-input
                :marks="marks"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="题目类型">
              <el-checkbox-group v-model="form.types">
                <el-checkbox value="single">单选题</el-checkbox>
                <el-checkbox value="multiple">多选题</el-checkbox>
                <el-checkbox value="truefalse">判断题</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="题目分类">
              <el-select v-model="form.categories" multiple placeholder="全部分类" clearable style="width: 100%">
                <el-option v-for="c in categories" :key="c" :label="c" :value="c" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="难度等级">
              <el-rate v-model="form.difficulty" :max="3" show-text :texts="['全部','简单','中等','困难']" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <div class="config-tip">
        <el-alert
          :title="tips"
          type="info"
          show-icon
          :closable="false"
        />
      </div>

      <div class="action-bar">
        <el-button size="large" @click="router.back()">
          <el-icon><Back /></el-icon>
          返回
        </el-button>
        <el-button type="primary" size="large" :loading="loading" @click="startExam">
          <el-icon><Right /></el-icon>
          开始答题
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { generateExam, getQuestionCategories } from '@/api'

const router = useRouter()
const selectedMode = ref('order')
const loading = ref(false)
const categories = ref([])
const totalQuestions = ref(50)

const form = reactive({
  count: 10,
  types: [],
  categories: [],
  difficulty: 0
})

const sliderMax = computed(() => Math.max(5, Math.min(totalQuestions.value, 100)))

const marks = computed(() => {
  const allMarks = {
    5: '5题',
    10: '10题',
    20: '20题',
    50: '50题',
    100: '100题'
  }
  const valid = {}
  for (const key in allMarks) {
    const numKey = Number(key)
    if (!isNaN(numKey) && numKey >= 5 && numKey <= sliderMax.value) {
      valid[numKey] = allMarks[key]
    }
  }
  return valid
})

const tips = computed(() => {
  const modeName = { order: '顺序刷题', random: '随机组卷', smart: '智能组卷' }[selectedMode.value]
  return `当前模式：${modeName} | 共 ${form.count} 题 | ${form.types.length ? form.types.map(t => ({ single: '单选', multiple: '多选', truefalse: '判断' }[t])).join('、') : '全部题型'}`
})

const loadCategories = async () => {
  const res = await getQuestionCategories()
  categories.value = res.categories || []
  totalQuestions.value = res.total || 50
}

const startExam = async () => {
  loading.value = true
  try {
    const payload = {
      mode: selectedMode.value,
      count: form.count,
      types: form.types,
      categories: form.categories,
      difficulty: form.difficulty
    }
    const res = await generateExam(payload)
    if (!res.questions || res.questions.length === 0) {
      ElMessage.warning('没有符合条件的题目，请调整筛选条件')
      return
    }
    const examData = {
      mode: selectedMode.value,
      questions: res.questions,
      answers: {},
      startTime: Date.now()
    }
    sessionStorage.setItem('exam_data', JSON.stringify(examData))
    router.push(`/exam/${selectedMode.value}`)
  } finally {
    loading.value = false
  }
}

onMounted(loadCategories)
</script>

<style scoped>
.page-container {
  max-width: 1100px;
  margin: 0 auto;
}

.main-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 17px;
  color: #303133;
}

.mode-card {
  padding: 28px 24px;
  border-radius: 12px;
  border: 2px solid #e4e7ed;
  cursor: pointer;
  transition: all 0.25s ease;
  text-align: center;
  background: #fff;
}

.mode-card:hover {
  border-color: #c6e2ff;
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.12);
}

.mode-card.active {
  border-color: #409EFF;
  background: #ecf5ff;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.2);
}

.mode-icon {
  width: 72px;
  height: 72px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin: 0 auto 16px;
}

.mode-card h3 {
  font-size: 18px;
  margin-bottom: 8px;
  color: #303133;
}

.mode-card p {
  font-size: 13px;
  color: #909399;
  line-height: 1.6;
}

.config-form {
  padding: 8px 0;
}

.config-tip {
  margin: 20px 0 28px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  padding-top: 24px;
  border-top: 1px solid #ebeef5;
}
</style>
