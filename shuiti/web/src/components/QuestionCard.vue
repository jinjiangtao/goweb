<template>
  <div class="question-card">
    <div class="question-header">
      <el-tag :type="typeTagType" size="large" effect="dark">{{ typeText }}</el-tag>
      <span class="question-index">第 {{ index }} 题 / 共 {{ total }} 题</span>
      <el-tag v-if="question.category" type="info" effect="plain">{{ question.category }}</el-tag>
    </div>

    <div class="question-content">
      <h3>{{ question.content }}</h3>
    </div>

    <div class="question-options" v-if="optionsList.length">
      <div
        v-for="(opt, idx) in optionsList"
        :key="idx"
        class="option-item"
        :class="getOptionClass(opt.label)"
        @click="selectOption(opt.label)"
      >
        <div class="option-label">
          <el-radio
            v-if="question.type === 'single' || question.type === 'truefalse'"
            :model-value="currentAnswer"
            :label="opt.label"
            :disabled="showResult"
          >{{ opt.label }}</el-radio>
          <el-checkbox
            v-else
            :model-value="isOptionSelected(opt.label)"
            :disabled="showResult"
          >{{ opt.label }}</el-checkbox>
        </div>
        <div class="option-text">{{ opt.text }}</div>
        <el-icon v-if="showResult" class="option-icon" :class="getOptionIconClass(opt.label)">
          <CircleCheckFilled v-if="isOptionCorrect(opt.label)" />
          <CircleCloseFilled v-else-if="isOptionWrong(opt.label)" />
        </el-icon>
      </div>
    </div>

    <div v-if="showResult" class="question-result">
      <div class="result-header">
        <el-tag :type="isCorrect ? 'success' : 'danger'" size="large" effect="dark">
          {{ isCorrect ? '✓ 回答正确' : '✗ 回答错误' }}
        </el-tag>
        <div class="answer-info">
          <span>正确答案：<strong class="correct-ans">{{ question.answer }}</strong></span>
          <span v-if="!isCorrect" style="margin-left: 20px;">
            你的答案：<strong class="wrong-ans">{{ currentAnswer || '未作答' }}</strong>
          </span>
        </div>
      </div>
      <div v-if="question.analysis" class="analysis-box">
        <div class="analysis-title">
          <el-icon><LightBulb /></el-icon>
          <span>题目解析</span>
        </div>
        <div class="analysis-content">{{ question.analysis }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  question: {
    type: Object,
    required: true
  },
  index: {
    type: Number,
    default: 1
  },
  total: {
    type: Number,
    default: 1
  },
  modelValue: {
    type: String,
    default: ''
  },
  showResult: {
    type: Boolean,
    default: false
  },
  resultDetail: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const currentAnswer = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  currentAnswer.value = val
})

watch(currentAnswer, (val) => {
  emit('update:modelValue', val)
  emit('change', val)
})

const typeText = computed(() => {
  const map = { single: '单选题', multiple: '多选题', truefalse: '判断题' }
  return map[props.question.type] || '题目'
})

const typeTagType = computed(() => {
  const map = { single: 'primary', multiple: 'warning', truefalse: 'success' }
  return map[props.question.type] || 'info'
})

const optionsList = computed(() => {
  if (!props.question.options) return []
  try {
    const opts = JSON.parse(props.question.options)
    return opts.map(o => {
      const parts = o.split('.')
      const label = parts[0].trim()
      const text = parts.slice(1).join('.').trim()
      return { label, text }
    })
  } catch {
    return []
  }
})

const isOptionSelected = (label) => {
  return currentAnswer.value?.includes(label)
}

const selectOption = (label) => {
  if (props.showResult) return

  if (props.question.type === 'multiple') {
    let arr = currentAnswer.value ? currentAnswer.value.split('') : []
    if (arr.includes(label)) {
      arr = arr.filter(x => x !== label)
    } else {
      arr.push(label)
      arr.sort()
    }
    currentAnswer.value = arr.join('')
  } else {
    currentAnswer.value = label
  }
}

const isCorrect = computed(() => {
  if (props.resultDetail) return props.resultDetail.is_correct
  return props.currentAnswer === props.question.answer
})

const correctAnswerArr = computed(() => {
  return props.question.answer?.split('') || []
})

const userAnswerArr = computed(() => {
  const ans = props.showResult && props.resultDetail
    ? props.resultDetail.user_answer
    : currentAnswer.value
  return ans?.split('') || []
})

const isOptionCorrect = (label) => correctAnswerArr.value.includes(label)
const isOptionWrong = (label) => {
  if (!props.showResult) return false
  return userAnswerArr.value.includes(label) && !correctAnswerArr.value.includes(label)
}

const getOptionClass = (label) => {
  if (!props.showResult) {
    if (isOptionSelected(label)) return 'selected'
    return ''
  }
  if (isOptionCorrect(label)) return 'correct'
  if (isOptionWrong(label)) return 'wrong'
  return ''
}

const getOptionIconClass = (label) => {
  if (isOptionCorrect(label)) return 'icon-correct'
  if (isOptionWrong(label)) return 'icon-wrong'
  return ''
}
</script>

<style scoped>
.question-card {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.question-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  gap: 16px;
}

.question-index {
  color: #909399;
  font-size: 14px;
  flex: 1;
}

.question-content h3 {
  font-size: 18px;
  line-height: 1.7;
  color: #303133;
  font-weight: 500;
  margin-bottom: 28px;
}

.question-options {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.option-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border: 2px solid #e4e7ed;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.option-item:hover {
  border-color: #c6e2ff;
  background: #f0f9ff;
}

.option-item.selected {
  border-color: #409EFF;
  background: #ecf5ff;
}

.option-item.correct {
  border-color: #67C23A;
  background: #f0f9eb;
}

.option-item.wrong {
  border-color: #F56C6C;
  background: #fef0f0;
}

.option-label {
  margin-right: 14px;
  min-width: 24px;
}

.option-text {
  flex: 1;
  font-size: 15px;
  color: #303133;
  line-height: 1.5;
}

.option-icon {
  margin-left: 12px;
  font-size: 22px;
}

.icon-correct {
  color: #67C23A;
}

.icon-wrong {
  color: #F56C6C;
}

.question-result {
  margin-top: 28px;
  padding-top: 24px;
  border-top: 1px dashed #ebeef5;
}

.result-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
}

.answer-info {
  flex: 1;
  display: flex;
  color: #606266;
  font-size: 15px;
}

.correct-ans {
  color: #67C23A;
  font-size: 16px;
  margin-left: 4px;
}

.wrong-ans {
  color: #F56C6C;
  font-size: 16px;
  margin-left: 4px;
}

.analysis-box {
  background: #fffbeb;
  border-radius: 10px;
  padding: 18px 20px;
  border-left: 4px solid #e6a23c;
}

.analysis-title {
  display: flex;
  align-items: center;
  color: #e6a23c;
  font-weight: 600;
  margin-bottom: 10px;
  gap: 6px;
}

.analysis-content {
  color: #606266;
  line-height: 1.8;
  font-size: 14px;
}
</style>
