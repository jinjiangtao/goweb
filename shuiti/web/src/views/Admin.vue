<template>
  <div class="page-container">
    <el-card class="main-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-icon :size="22" color="#409EFF"><Management /></el-icon>
            <span>题库管理</span>
          </div>
          <div class="header-actions">
            <el-button type="success" @click="showBatchAdd = true">
              <el-icon><DocumentAdd /></el-icon>
              批量新增
            </el-button>
            <el-button type="primary" @click="openAddDialog">
              <el-icon><Plus /></el-icon>
              新增题目
            </el-button>
            <el-button type="danger" :disabled="selectedIds.length === 0" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
      </template>

      <div class="filter-bar">
        <el-form :inline="true" :model="filters">
          <el-form-item label="题目类型">
            <el-select v-model="filters.type" placeholder="全部" clearable style="width: 140px">
              <el-option label="单选题" value="single" />
              <el-option label="多选题" value="multiple" />
              <el-option label="判断题" value="truefalse" />
            </el-select>
          </el-form-item>
          <el-form-item label="分类">
            <el-select v-model="filters.category" placeholder="全部" clearable filterable style="width: 160px">
              <el-option v-for="c in categories" :key="c" :label="c" :value="c" />
            </el-select>
          </el-form-item>
          <el-form-item label="关键词">
            <el-input v-model="filters.keyword" placeholder="搜索题目内容..." clearable style="width: 240px" @keyup.enter="loadData" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="resetFilters">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
        :data="tableData"
        v-loading="loading"
        border
        stripe
        @selection-change="handleSelectionChange"
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="typeTagColor(row.type)" size="small">{{ typeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="题目内容" min-width="280" show-overflow-tooltip />
        <el-table-column label="答案" width="100" align="center">
          <template #default="{ row }">
            <el-tag type="success" effect="plain">{{ row.answer }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120" align="center">
          <template #default="{ row }">
            <span v-if="row.category">{{ row.category }}</span>
            <span v-else style="color:#c0c4cc">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="difficulty" label="难度" width="80" align="center">
          <template #default="{ row }">
            <el-rate v-model="row.difficulty" disabled :max="3" size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" align="center" />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="680px" destroy-on-close>
      <el-form :model="form" label-width="90px" ref="formRef" :rules="formRules">
        <el-form-item label="题目类型" prop="type">
          <el-radio-group v-model="form.type" @change="onTypeChange">
            <el-radio value="single">单选题</el-radio>
            <el-radio value="multiple">多选题</el-radio>
            <el-radio value="truefalse">判断题</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="题目内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="3" placeholder="请输入题目内容" />
        </el-form-item>
        <template v-if="form.type !== 'truefalse'">
          <el-form-item v-for="i in optionCount" :key="i" :label="`选项${labelOf(i)}`">
            <el-input v-model="form[`option_${labelOf(i)}`]" :placeholder="`请输入选项 ${labelOf(i)} 内容`" />
          </el-form-item>
        </template>
        <template v-else>
          <el-form-item label="选项A">
            <el-input v-model="form.option_A" disabled value="正确" />
          </el-form-item>
          <el-form-item label="选项B">
            <el-input v-model="form.option_B" disabled value="错误" />
          </el-form-item>
        </template>
        <el-form-item label="正确答案" prop="answer">
          <template v-if="form.type === 'multiple'">
            <el-checkbox-group v-model="answerArr">
              <el-checkbox v-for="i in optionCount" :key="i" :label="labelOf(i)">{{ labelOf(i) }}</el-checkbox>
            </el-checkbox-group>
          </template>
          <template v-else-if="form.type === 'truefalse'">
            <el-radio-group v-model="form.answer">
              <el-radio value="A">正确 (A)</el-radio>
              <el-radio value="B">错误 (B)</el-radio>
            </el-radio-group>
          </template>
          <template v-else>
            <el-radio-group v-model="form.answer">
              <el-radio v-for="i in optionCount" :key="i" :label="labelOf(i)">{{ labelOf(i) }}</el-radio>
            </el-radio-group>
          </template>
        </el-form-item>
        <el-form-item label="题目解析">
          <el-input v-model="form.analysis" type="textarea" :rows="3" placeholder="请输入题目解析（可选）" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分类">
              <el-input v-model="form.category" placeholder="如：Go基础" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="难度">
              <el-rate v-model="form.difficulty" :max="3" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showBatchAdd" title="批量新增题目" width="760px" destroy-on-close>
      <el-alert
        title="每行一道题，格式：类型|内容|选项|答案|解析|分类|难度"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 16px;"
      />
      <div style="margin-bottom:12px; color:#606266; font-size:13px;">
        示例：single|Go的声明变量方式？|var x int=10#variable x=10#int x=10#x :=: 10|A|Go使用var或:=|Go基础|1
      </div>
      <el-input
        v-model="batchText"
        type="textarea"
        :rows="14"
        placeholder="请按格式输入..."
      />
      <template #footer>
        <el-button @click="showBatchAdd = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleBatchSubmit">提交导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getQuestions, getQuestionCategories, createQuestion, updateQuestion,
  deleteQuestion, batchDeleteQuestions, batchCreateQuestions
} from '@/api'

const loading = ref(false)
const saving = ref(false)
const tableData = ref([])
const categories = ref([])
const selectedIds = ref([])
const dialogVisible = ref(false)
const showBatchAdd = ref(false)
const batchText = ref('')
const formRef = ref()
const editingId = ref(null)

const filters = reactive({
  type: '',
  category: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const form = reactive({
  type: 'single',
  content: '',
  option_A: '',
  option_B: '',
  option_C: '',
  option_D: '',
  option_E: '',
  option_F: '',
  answer: '',
  analysis: '',
  category: '',
  difficulty: 1
})

const answerArr = ref([])

const optionCount = computed(() => {
  if (form.type === 'truefalse') return 2
  return 4
})

const dialogTitle = computed(() => editingId.value ? '编辑题目' : '新增题目')

const typeText = (t) => ({ single: '单选题', multiple: '多选题', truefalse: '判断题' }[t] || t)
const typeTagColor = (t) => ({ single: 'primary', multiple: 'warning', truefalse: 'success' }[t] || 'info')
const labelOf = (i) => String.fromCharCode(64 + i)

watch(answerArr, (val) => {
  if (form.type === 'multiple') {
    form.answer = [...val].sort().join('')
  }
})

const formRules = {
  type: [{ required: true, message: '请选择题型', trigger: 'change' }],
  content: [{ required: true, message: '请输入题目内容', trigger: 'blur' }],
  answer: [{ required: true, message: '请选择正确答案', trigger: 'change' }]
}

const onTypeChange = () => {
  form.answer = ''
  answerArr.value = []
  if (form.type === 'truefalse') {
    form.option_A = '正确'
    form.option_B = '错误'
  }
}

const loadCategories = async () => {
  const res = await getQuestionCategories()
  categories.value = res.categories || []
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getQuestions({
      page: pagination.page,
      page_size: pagination.page_size,
      type: filters.type,
      category: filters.category,
      keyword: filters.keyword
    })
    tableData.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.type = ''
  filters.category = ''
  filters.keyword = ''
  pagination.page = 1
  loadData()
}

const handleSelectionChange = (val) => {
  selectedIds.value = val.map(x => x.id)
}

const resetForm = () => {
  Object.assign(form, {
    type: 'single',
    content: '',
    option_A: '',
    option_B: '',
    option_C: '',
    option_D: '',
    option_E: '',
    option_F: '',
    answer: '',
    analysis: '',
    category: '',
    difficulty: 1
  })
  answerArr.value = []
  editingId.value = null
}

const openAddDialog = () => {
  resetForm()
  dialogVisible.value = true
}

const parseOptions = (optionsStr) => {
  try {
    const arr = JSON.parse(optionsStr)
    const labels = ['A', 'B', 'C', 'D', 'E', 'F']
    arr.forEach((opt, i) => {
      const parts = opt.split('.')
      if (parts.length > 1) {
        form[`option_${labels[i]}`] = parts.slice(1).join('.').trim()
      }
    })
  } catch {}
}

const buildOptions = () => {
  const count = form.type === 'truefalse' ? 2 : optionCount.value
  const opts = []
  for (let i = 1; i <= count; i++) {
    const label = labelOf(i)
    const text = form.type === 'truefalse'
      ? (label === 'A' ? '正确' : '错误')
      : form[`option_${label}`] || ''
    opts.push(`${label}. ${text}`)
  }
  return JSON.stringify(opts)
}

const openEditDialog = (row) => {
  resetForm()
  editingId.value = row.id
  Object.assign(form, {
    type: row.type,
    content: row.content,
    answer: row.answer,
    analysis: row.analysis || '',
    category: row.category || '',
    difficulty: row.difficulty || 1
  })
  if (row.type === 'multiple') {
    answerArr.value = row.answer.split('')
  }
  parseOptions(row.options)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  saving.value = true
  try {
    const payload = {
      type: form.type,
      content: form.content,
      options: buildOptions(),
      answer: form.answer,
      analysis: form.analysis,
      category: form.category,
      difficulty: form.difficulty
    }
    if (editingId.value) {
      await updateQuestion(editingId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await createQuestion(payload)
      ElMessage.success('新增成功')
    }
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定删除 ID:${row.id} 的题目？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteQuestion(row.id)
      ElMessage.success('删除成功')
      loadData()
    }).catch(() => {})
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确定删除选中的 ${selectedIds.value.length} 道题目？`, '确认', { type: 'warning' })
    .then(async () => {
      await batchDeleteQuestions(selectedIds.value)
      ElMessage.success('删除成功')
      loadData()
    }).catch(() => {})
}

const handleBatchSubmit = async () => {
  if (!batchText.value.trim()) {
    ElMessage.warning('请输入内容')
    return
  }
  saving.value = true
  try {
    const lines = batchText.value.trim().split('\n').filter(l => l.trim())
    const questions = lines.map(line => {
      const parts = line.split('|')
      const [type, content, optsRaw, answer, analysis = '', category = '', difficulty = '1'] = parts
      const opts = optsRaw.split('#').map((o, i) => `${labelOf(i + 1)}. ${o.trim()}`)
      return {
        type,
        content: content.trim(),
        options: JSON.stringify(opts),
        answer: answer.trim(),
        analysis: analysis.trim(),
        category: category.trim(),
        difficulty: parseInt(difficulty) || 1
      }
    })
    const res = await batchCreateQuestions(questions)
    ElMessage.success(`成功导入 ${res.count} 道题`)
    showBatchAdd.value = false
    batchText.value = ''
    loadData()
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadData()
  loadCategories()
})
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

.filter-bar {
  padding: 16px 0;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 20px;
}

.pagination-wrap {
  display: flex;
  justify-content: center;
  padding-top: 24px;
}
</style>
