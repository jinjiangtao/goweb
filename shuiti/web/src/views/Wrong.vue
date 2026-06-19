<template>
  <div class="page-container">
    <el-card class="main-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-icon :size="22" color="#F56C6C"><Warning /></el-icon>
            <span>错题本</span>
            <el-tag type="danger" effect="plain" style="margin-left: 12px;">共 {{ total }} 道错题</el-tag>
          </div>
          <div class="header-actions">
            <el-button type="warning" :disabled="list.length === 0" @click="startPractice">
              <el-icon><EditPen /></el-icon>
              重做错题
            </el-button>
            <el-popconfirm title="确定清空所有错题吗？此操作不可恢复" @confirm="handleClearAll">
              <template #reference>
                <el-button type="danger" :disabled="list.length === 0">
                  <el-icon><Delete /></el-icon>
                  清空错题
                </el-button>
              </template>
            </el-popconfirm>
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
        <el-table-column label="序号" width="70" align="center" type="index">
          <template #default="{ $index }">
            {{ (pagination.page - 1) * pagination.page_size + $index + 1 }}
          </template>
        </el-table-column>
        <el-table-column label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="typeColor(row.question?.type)" size="small">{{ typeName(row.question?.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="题目内容" min-width="300" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.question?.content || '题目已删除' }}
          </template>
        </el-table-column>
        <el-table-column label="分类" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.question?.category">{{ row.question.category }}</span>
            <span v-else style="color:#c0c4cc">-</span>
          </template>
        </el-table-column>
        <el-table-column label="错误次数" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.wrong_count >= 3 ? 'danger' : (row.wrong_count >= 2 ? 'warning' : 'info')" effect="plain">
              {{ row.wrong_count }} 次
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="你的答案" width="100" align="center">
          <template #default="{ row }">
            <span style="color: #F56C6C; font-weight: 500;">{{ row.user_answer || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="正确答案" width="100" align="center">
          <template #default="{ row }">
            <span style="color: #67C23A; font-weight: 500;">{{ row.question?.answer || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_wrong_at" label="最近错误" width="180" align="center" />
        <el-table-column label="操作" width="160" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="showDetail(row)">查看</el-button>
            <el-popconfirm title="确定从错题本中移除此题？" @confirm="handleRemove(row)">
              <template #reference>
                <el-button link type="danger" size="small">移除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap" v-if="total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>

      <el-empty v-if="!loading && total === 0" description="暂无错题，继续加油！" />
    </el-card>

    <el-dialog v-model="detailVisible" title="错题详情" width="680px" destroy-on-close>
      <QuestionCard
        v-if="currentQuestion"
        :question="currentQuestion"
        :index="1"
        :total="1"
        :model-value="currentUserAnswer"
        :show-result="true"
        :result-detail="currentDetail"
      />
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import QuestionCard from '@/components/QuestionCard.vue'
import { getWrongQuestions, deleteWrongQuestion, clearWrongQuestions, getWrongPractice } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const detailVisible = ref(false)
const currentQuestion = ref(null)
const currentUserAnswer = ref('')
const currentDetail = ref(null)

const pagination = reactive({
  page: 1,
  page_size: 20
})

const typeName = (t) => ({ single: '单选题', multiple: '多选题', truefalse: '判断题' }[t] || t)
const typeColor = (t) => ({ single: 'primary', multiple: 'warning', truefalse: 'success' }[t] || 'info')

const loadData = async () => {
  loading.value = true
  try {
    const res = await getWrongQuestions({
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

const showDetail = (row) => {
  currentQuestion.value = row.question
  currentUserAnswer.value = row.user_answer
  currentDetail.value = {
    user_answer: row.user_answer,
    is_correct: false
  }
  detailVisible.value = true
}

const handleRemove = async (row) => {
  await deleteWrongQuestion(row.id, userStore.userId)
  ElMessage.success('已移除')
  loadData()
}

const handleClearAll = async () => {
  await clearWrongQuestions(userStore.userId)
  ElMessage.success('已清空所有错题')
  loadData()
}

const startPractice = async () => {
  const res = await getWrongPractice({ user_id: userStore.userId })
  if (!res.questions || res.questions.length === 0) {
    ElMessage.warning('暂无可练习的错题')
    return
  }
  const examData = {
    mode: 'wrong',
    questions: res.questions,
    answers: {},
    startTime: Date.now()
  }
  sessionStorage.setItem('exam_data', JSON.stringify(examData))
  router.push('/exam/wrong')
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
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 17px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.pagination-wrap {
  display: flex;
  justify-content: center;
  padding-top: 24px;
}
</style>
