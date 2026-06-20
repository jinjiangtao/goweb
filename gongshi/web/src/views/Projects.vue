<template>
  <div class="page-container">
    <div class="card">
      <div class="card-header">
        <div class="card-title">
          <el-icon style="vertical-align: middle; margin-right: 6px; color: #409eff"><Folder /></el-icon>
          项目管理
        </div>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>新增项目
        </el-button>
      </div>

      <el-table :data="projects" border v-loading="loading">
        <el-table-column type="index" label="序号" width="70" align="center" />
        <el-table-column prop="code" label="项目编号" width="140" />
        <el-table-column prop="name" label="项目名称" min-width="160" />
        <el-table-column prop="description" label="项目描述" min-width="240" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'active'" type="success" size="small">启用</el-tag>
            <el-tag v-else type="info" size="small">停用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑项目' : '新增项目'" width="520px" destroy-on-close>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="项目编号" prop="code">
          <el-input v-model="form.code" placeholder="如: PRJ-001" :disabled="!!editing" />
        </el-form-item>
        <el-form-item label="项目名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入项目名称" maxlength="100" show-word-limit />
        </el-form-item>
        <el-form-item label="项目描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="4" placeholder="请输入项目描述" maxlength="500" show-word-limit />
        </el-form-item>
        <el-form-item label="状态" v-if="editing">
          <el-radio-group v-model="form.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="inactive">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getProjects, createProject, updateProject, deleteProject } from '../api/project'

const projects = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editing = ref(null)
const submitting = ref(false)
const formRef = ref()

const form = reactive({
  code: '',
  name: '',
  description: '',
  status: 'active'
})

const rules = {
  code: [{ required: true, message: '请输入项目编号', trigger: 'blur', min: 2, max: 50 }],
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur', min: 2, max: 100 }]
}

const loadProjects = async () => {
  loading.value = true
  try {
    projects.value = await getProjects()
  } finally {
    loading.value = false
  }
}

const openAddDialog = () => {
  editing.value = null
  Object.assign(form, { code: '', name: '', description: '', status: 'active' })
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  editing.value = row
  Object.assign(form, {
    code: row.code,
    name: row.name,
    description: row.description,
    status: row.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  submitting.value = true
  try {
    if (editing.value) {
      await updateProject(editing.value.id, { ...form })
      ElMessage.success('修改成功')
    } else {
      await createProject({ ...form })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadProjects()
  } finally {
    submitting.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除项目「${row.name}吗？关联的工时记录将无法使用该项目`, '提示', { type: 'warning' })
    .then(async () => {
      await deleteProject(row.id)
      ElMessage.success('删除成功')
      loadProjects()
    }).catch(() => {})
}

onMounted(() => {
  loadProjects()
})
</script>
