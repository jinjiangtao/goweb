<template>
  <div class="template-page">
    <div class="page-header">
      <h2>流程模板管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><DocumentAdd /></el-icon>
          新建模板
        </el-button>
        <el-button @click="goToDesigner">
          <el-icon><Edit /></el-icon>
          流程设计器
        </el-button>
      </div>
    </div>

    <div class="search-bar">
      <el-input
        v-model="searchForm.keyword"
        placeholder="搜索模板名称"
        style="width: 250px; margin-right: 15px;"
        clearable
        @keyup.enter="loadTemplates"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select
        v-model="searchForm.type"
        placeholder="模板类型"
        style="width: 150px; margin-right: 15px;"
        clearable
      >
        <el-option label="请假审批" value="leave" />
        <el-option label="报销审批" value="reimburse" />
        <el-option label="通用申请" value="general" />
      </el-select>
      <el-button type="primary" @click="loadTemplates">搜索</el-button>
      <el-button @click="resetSearch">重置</el-button>
    </div>

    <el-table :data="templateList" border stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="模板名称" min-width="150" />
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)">{{ getTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : 'info'">
            {{ row.status === 'published' ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="creator.name" label="创建人" width="120" />
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="goToEdit(row.id)">编辑</el-button>
          <el-button size="small" type="success" @click="goToSubmit(row.id)">提交申请</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="loadTemplates"
      @current-change="loadTemplates"
      style="margin-top: 20px; justify-content: flex-end;"
    />

    <el-dialog v-model="createDialogVisible" title="创建模板" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="模板名称" required>
          <el-input v-model="createForm.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="模板类型" required>
          <el-select v-model="createForm.type" placeholder="请选择模板类型" style="width: 100%;">
            <el-option label="请假审批" value="leave" />
            <el-option label="报销审批" value="reimburse" />
            <el-option label="通用申请" value="general" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateDefault">创建默认模板</el-button>
        <el-button type="success" @click="goToDesignerWithType">空白设计</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Edit, DocumentAdd } from '@element-plus/icons-vue'
import { getTemplateList, deleteTemplate, createDefaultTemplate } from '@/api/template'

const router = useRouter()

const loading = ref(false)
const templateList = ref([])
const createDialogVisible = ref(false)

const searchForm = reactive({
  keyword: '',
  type: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createForm = reactive({
  name: '',
  type: 'general'
})

onMounted(() => {
  loadTemplates()
})

const loadTemplates = async () => {
  loading.value = true
  try {
    const res = await getTemplateList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword,
      type: searchForm.type
    })
    templateList.value = res.data.list
    pagination.total = res.data.total
  } catch (e) {
    ElMessage.error('加载模板列表失败')
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  pagination.page = 1
  loadTemplates()
}

const getTypeTagType = (type) => {
  switch (type) {
    case 'leave':
      return 'primary'
    case 'reimburse':
      return 'success'
    default:
      return 'info'
  }
}

const getTypeText = (type) => {
  switch (type) {
    case 'leave':
      return '请假审批'
    case 'reimburse':
      return '报销审批'
    default:
      return '通用申请'
  }
}

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const showCreateDialog = () => {
  createForm.name = ''
  createForm.type = 'general'
  createDialogVisible.value = true
}

const handleCreateDefault = async () => {
  if (!createForm.name.trim()) {
    ElMessage.warning('请输入模板名称')
    return
  }
  try {
    await createDefaultTemplate({
      name: createForm.name,
      type: createForm.type
    })
    ElMessage.success('模板创建成功')
    createDialogVisible.value = false
    loadTemplates()
  } catch (e) {
    ElMessage.error(e.message || '创建失败')
  }
}

const goToDesigner = () => {
  createDialogVisible.value = false
  router.push('/template/designer')
}

const goToDesignerWithType = () => {
  if (!createForm.name.trim()) {
    ElMessage.warning('请输入模板名称')
    return
  }
  createDialogVisible.value = false
  router.push({
    path: '/template/designer',
    query: { name: createForm.name, type: createForm.type }
  })
}

const goToEdit = (id) => {
  router.push(`/template/designer/${id}`)
}

const goToSubmit = (templateId) => {
  router.push(`/approval/submit/${templateId}`)
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该模板吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteTemplate(row.id)
    ElMessage.success('删除成功')
    loadTemplates()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}
</script>

<style scoped lang="scss">
.template-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 20px;
      color: #303133;
    }
  }

  .search-bar {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
  }
}
</style>
