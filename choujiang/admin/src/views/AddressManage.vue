<template>
  <div>
    <div style="margin-bottom: 20px;">
      <el-button type="primary" @click="openAddDialog('province')">添加省份</el-button>
    </div>

    <el-table :data="addressTree" row-key="id" border default-expand-all>
      <el-table-column prop="name" label="名称"></el-table-column>
      <el-table-column prop="level" label="级别" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.level === 1" type="danger">省</el-tag>
          <el-tag v-else-if="row.level === 2" type="warning">市</el-tag>
          <el-tag v-else type="success">区/县</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="code" label="代码" width="150"></el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button 
            v-if="row.level < 3" 
            type="success" 
            size="small" 
            @click="openAddDialog(row.level === 1 ? 'city' : 'district', row)"
          >
            添加下级
          </el-button>
          <el-popconfirm
            title="确定要删除吗？"
            @confirm="handleDelete(row)"
          >
            <template #reference>
              <el-button type="danger" size="small">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" placeholder="请输入名称"></el-input>
        </el-form-item>
        <el-form-item label="代码">
          <el-input v-model="form.code" placeholder="请输入代码"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})

const addressTree = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const isEdit = ref(false)
const currentRow = ref(null)
const form = ref({
  parentId: 0,
  name: '',
  level: 1,
  code: ''
})

const loadTree = async () => {
  try {
    const res = await api.get('/admin/address/tree')
    addressTree.value = res.data
  } catch (err) {
    ElMessage.error('加载失败')
  }
}

const openAddDialog = (type, parent = null) => {
  isEdit.value = false
  if (type === 'province') {
    dialogTitle.value = '添加省份'
    form.value = { parentId: 0, name: '', level: 1, code: '' }
  } else if (type === 'city') {
    dialogTitle.value = '添加城市'
    form.value = { parentId: parent.id, name: '', level: 2, code: '' }
  } else {
    dialogTitle.value = '添加区县'
    form.value = { parentId: parent.id, name: '', level: 3, code: '' }
  }
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  isEdit.value = true
  currentRow.value = row
  dialogTitle.value = '编辑'
  form.value = {
    parentId: row.parentId,
    name: row.name,
    level: row.level,
    code: row.code
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入名称')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await api.put(`/admin/address/${currentRow.value.id}`, form.value)
      ElMessage.success('编辑成功')
    } else {
      await api.post('/admin/address', form.value)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    loadTree()
  } catch (err) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await api.delete(`/admin/address/${row.id}`)
    ElMessage.success('删除成功')
    loadTree()
  } catch (err) {
    if (err.response && err.response.data && err.response.data.error) {
      ElMessage.error(err.response.data.error)
    } else {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadTree()
})
</script>
