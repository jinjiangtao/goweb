<template>
  <div class="menu-list">
    <div class="search-bar">
      <el-button type="success" @click="openAddModal">添加菜单</el-button>
    </div>

    <el-table :data="tableData" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="菜单名称" />
      <el-table-column prop="path" label="路由路径" />
      <el-table-column prop="icon" label="图标" width="100">
        <template #default="scope">
          <el-icon><component :is="getIcon(scope.row.icon)" /></el-icon>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序号" width="100" />
      <el-table-column prop="parent_name" label="父级菜单" width="120" />
      <el-table-column prop="visible" label="是否显示" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.visible === 1 ? 'success' : 'warning'">
            {{ scope.row.visible === 1 ? '显示' : '隐藏' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="scope">
          <el-button type="text" @click="openEditModal(scope.row)">编辑</el-button>
          <el-button type="text" danger @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :title="addForm.id ? '编辑菜单' : '添加菜单'" :visible.sync="addModalVisible" width="450px">
      <el-form :model="addForm" ref="addFormRef" label-width="100px">
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="addForm.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="路由路径" prop="path">
          <el-input v-model="addForm.path" placeholder="请输入路由路径" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-select v-model="addForm.icon">
            <el-option label="列表" value="List" />
            <el-option label="数据分析" value="DataAnalysis" />
            <el-option label="学校" value="School" />
            <el-option label="用户" value="User" />
            <el-option label="菜单" value="Menu" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序号" prop="sort">
          <el-input v-model.number="addForm.sort" placeholder="数字越小越靠前" />
        </el-form-item>
        <el-form-item label="父级菜单" prop="parent_id">
          <el-select v-model="addForm.parent_id">
            <el-option label="顶级（一级菜单）" :value="0" />
            <el-option v-for="menu in parentMenus" :key="menu.id" :label="menu.name" :value="menu.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否显示" prop="visible">
          <el-select v-model="addForm.visible">
            <el-option label="显示" :value="1" />
            <el-option label="隐藏" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addModalVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, markRaw } from 'vue'
import axios from '../utils/axios'
import { ElMessage } from 'element-plus'
import { List, DataAnalysis, School, User, Menu } from '@element-plus/icons-vue'

const iconMap = {
  List: markRaw(List),
  DataAnalysis: markRaw(DataAnalysis),
  School: markRaw(School),
  User: markRaw(User),
  Menu: markRaw(Menu)
}

const getIcon = (iconName) => {
  return iconMap[iconName] || List
}

const tableData = ref([])
const parentMenus = ref([])
const addModalVisible = ref(false)
const addFormRef = ref()

const addForm = reactive({
  id: 0,
  name: '',
  path: '',
  icon: 'List',
  sort: 0,
  parent_id: 0,
  visible: 1
})

const loadData = async () => {
  try {
    const [menusRes, parentRes] = await Promise.all([
      axios.get('/api/admin/menus'),
      axios.get('/api/admin/menus/parent')
    ])
    
    const menus = menusRes.data.data
    const parents = parentRes.data.data
    
    parentMenus.value = parents
    
    const parentMap = {}
    parents.forEach(p => {
      parentMap[p.id] = p.name
    })
    
    tableData.value = menus.map(menu => ({
      ...menu,
      parent_name: menu.parent_id === 0 ? '顶级' : parentMap[menu.parent_id] || '未知'
    }))
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
  }
}

const openAddModal = () => {
  addForm.id = 0
  addForm.name = ''
  addForm.path = ''
  addForm.icon = 'List'
  addForm.sort = 0
  addForm.parent_id = 0
  addForm.visible = 1
  addModalVisible.value = true
}

const openEditModal = (row) => {
  addForm.id = row.id
  addForm.name = row.name
  addForm.path = row.path
  addForm.icon = row.icon || 'List'
  addForm.sort = row.sort
  addForm.parent_id = row.parent_id
  addForm.visible = row.visible
  addModalVisible.value = true
}

const handleSave = async () => {
  if (!addForm.name) {
    ElMessage.warning('请填写菜单名称')
    return
  }

  try {
    if (addForm.id) {
      await axios.put(`/api/admin/menus/${addForm.id}`, {
        name: addForm.name,
        path: addForm.path,
        icon: addForm.icon,
        sort: addForm.sort,
        parent_id: addForm.parent_id,
        visible: addForm.visible
      })
      ElMessage.success('更新成功')
    } else {
      await axios.post('/api/admin/menus', {
        name: addForm.name,
        path: addForm.path,
        icon: addForm.icon,
        sort: addForm.sort,
        parent_id: addForm.parent_id,
        visible: addForm.visible
      })
      ElMessage.success('添加成功')
    }
    addModalVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

const handleDelete = async (row) => {
  try {
    await axios.delete(`/api/admin/menus/${row.id}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.menu-list {
  padding: 20px;
}
.search-bar {
  margin-bottom: 20px;
}
</style>