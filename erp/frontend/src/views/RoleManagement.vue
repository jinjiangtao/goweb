<template>
  <div class="role-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增角色
          </el-button>
        </div>
      </template>
      <el-table :data="tableData" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="code" label="角色代码" width="150" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" />
        <el-table-column label="操作" fixed="right" width="300">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="success" size="small" @click="handleAssignMenus(row)">
              分配菜单
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="resetForm"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色代码" prop="code">
          <el-input v-model="form.code" placeholder="请输入角色代码" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          确定
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="menuDialogVisible"
      title="分配菜单"
      width="500px"
    >
      <el-tree
        ref="menuTreeRef"
        :data="menuData"
        :props="menuProps"
        show-checkbox
        node-key="id"
        :default-checked-keys="checkedMenuIds"
      />
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveMenus" :loading="menuLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRoles, createRole, updateRole, deleteRole, assignMenus } from '@/api/role'
import { getMenus } from '@/api/menu'

const tableData = ref([])
const menuData = ref([])
const dialogVisible = ref(false)
const menuDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const loading = ref(false)
const menuLoading = ref(false)
const formRef = ref(null)
const menuTreeRef = ref(null)
const currentRoleId = ref(null)
const checkedMenuIds = ref([])

const menuProps = {
  children: 'children',
  label: 'name'
}

const form = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色代码', trigger: 'blur' }]
}

// 将扁平菜单转换为树形结构
const buildMenuTree = (menus) =&gt; {
  const menuMap = {}
  const result = []
  
  menus.forEach(menu =&gt; {
    menuMap[menu.id] = { ...menu, children: [] }
  })
  
  menus.forEach(menu =&gt; {
    if (menu.parentId === 0 || !menu.parentId) {
      result.push(menuMap[menu.id])
    } else if (menuMap[menu.parentId]) {
      menuMap[menu.parentId].children.push(menuMap[menu.id])
    }
  })
  
  return result
}

const fetchData = async () =&gt; {
  try {
    const res = await getRoles()
    tableData.value = Array.isArray(res.data) ? res.data : []
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  }
}

const fetchMenus = async () =&gt; {
  try {
    const res = await getMenus()
    const menus = Array.isArray(res.data) ? res.data : []
    menuData.value = buildMenuTree(menus)
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增角色'
  isEdit.value = false
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑角色'
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleAssignMenus = async (row) => {
  currentRoleId.value = row.id
  await fetchMenus()
  checkedMenuIds.value = (row.menus || []).map(m => m.id)
  menuDialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该角色吗？', '提示', { type: 'warning' })
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (isEdit.value) {
          await updateRole(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createRole(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '操作失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const handleSaveMenus = async () => {
  if (!menuTreeRef.value) return
  const checkedKeys = menuTreeRef.value.getCheckedKeys()
  const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys()
  const allMenuIds = [...checkedKeys, ...halfCheckedKeys]
  
  menuLoading.value = true
  try {
    await assignMenus(currentRoleId.value, allMenuIds)
    ElMessage.success('分配成功')
    menuDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('分配失败')
  } finally {
    menuLoading.value = false
  }
}

const resetForm = () => {
  form.id = null
  form.name = ''
  form.code = ''
  form.description = ''
  form.status = 1
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
