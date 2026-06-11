<template>
  <div class="menu-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>菜单管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增菜单
          </el-button>
        </div>
      </template>
      
      <el-table :data="tableData" :row-key="(row) => row.id || row.ID" default-expand-all style="width: 100%">
        <el-table-column prop="name" label="菜单名称" width="200" />
        <el-table-column prop="path" label="路由路径" width="200" />
        <el-table-column prop="icon" label="图标" width="100">
          <template #default="{ row }">
            {{ row.icon || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.hidden ? 'info' : 'success'">
              {{ row.hidden ? '隐藏' : '显示' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt || row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="250">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="success" size="small" @click="handleAddChild(row)">
              子菜单
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
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="路由路径" prop="path">
          <el-input v-model="form.path" placeholder="请输入路由路径" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="form.icon" placeholder="请输入图标名称" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="父菜单" prop="parentId">
          <el-select v-model="form.parentId" placeholder="请选择父菜单" clearable style="width: 100%">
            <el-option label="顶级菜单" :value="0" />
            <el-option
              v-for="menu in menuOptions"
              :key="menu.id"
              :label="menu.name"
              :value="menu.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="是否隐藏" prop="hidden">
          <el-switch v-model="form.hidden" :active-value="true" :inactive-value="false" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getMenus, getMenu, createMenu, updateMenu, deleteMenu } from '../api/menu'

const tableData = ref([])
const menuOptions = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const loading = ref(false)
const formRef = ref(null)

const form = reactive({
  id: null,
  name: '',
  path: '',
  icon: '',
  sort: 0,
  parentId: 0,
  hidden: false
})

const rules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  path: [{ required: true, message: '请输入路由路径', trigger: 'blur' }]
}

const formatDate = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleString('zh-CN')
}

const buildMenuTree = (menus) => {
  const menuMap = {}
  const result = []
  
  menus.forEach(menu => {
    const menuId = menu.id || menu.ID
    menuMap[menuId] = { ...menu, children: [] }
  })
  
  menus.forEach(menu => {
    const menuId = menu.id || menu.ID
    const parentId = menu.parentId || menu.ParentId || 0
    if (!parentId || parentId === 0) {
      result.push(menuMap[menuId])
    } else if (menuMap[parentId]) {
      menuMap[parentId].children.push(menuMap[menuId])
    }
  })
  
  return result
}

const buildMenuOptions = (menus) => {
  const result = []
  
  const traverse = (list, level = 0) => {
    list.forEach(menu => {
      const menuId = menu.id || menu.ID
      const prefix = level > 0 ? '　'.repeat(level) + '└ ' : ''
      result.push({
        id: menuId,
        name: prefix + menu.name
      })
      if (menu.children && menu.children.length > 0) {
        traverse(menu.children, level + 1)
      }
    })
  }
  
  traverse(buildMenuTree(menus))
  return result
}

const fetchData = async () => {
  try {
    const res = await getMenus()
    const menus = Array.isArray(res.data) ? res.data : []
    tableData.value = buildMenuTree(menus)
    menuOptions.value = buildMenuOptions(menus)
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增菜单'
  isEdit.value = false
  form.parentId = 0
  dialogVisible.value = true
}

const handleAddChild = (row) => {
  const rowId = row.id || row.ID
  dialogTitle.value = '新增子菜单'
  isEdit.value = false
  form.parentId = rowId
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  try {
    const rowId = row.id || row.ID
    const res = await getMenu(rowId)
    const menu = res.data
    Object.assign(form, {
      id: menu.id || menu.ID,
      name: menu.name,
      path: menu.path,
      icon: menu.icon,
      sort: menu.sort,
      parentId: menu.parentId || menu.ParentId || 0,
      hidden: menu.hidden || false
    })
    dialogTitle.value = '编辑菜单'
    isEdit.value = true
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取菜单信息失败')
  }
}

const handleDelete = async (row) => {
  try {
    const rowId = row.id || row.ID
    await ElMessageBox.confirm('确定要删除该菜单吗？', '提示', { type: 'warning' })
    await deleteMenu(rowId)
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
          await updateMenu(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createMenu(form)
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

const resetForm = () => {
  form.id = null
  form.name = ''
  form.path = ''
  form.icon = ''
  form.sort = 0
  form.parentId = 0
  form.hidden = false
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
