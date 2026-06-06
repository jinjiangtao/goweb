<template>
  <div class="rooms">
    <div class="header">
      <h2>会议室管理</h2>
      <div class="actions">
        <el-input v-model="search" placeholder="搜索会议室" style="width: 240px" clearable />
        <el-button type="primary" @click="openDialog()">新增会议室</el-button>
      </div>
    </div>

    <el-table :data="filteredRooms" stripe style="width: 100%">
      <el-table-column prop="name" label="会议室名称" />
      <el-table-column prop="capacity" label="容纳人数" width="120" />
      <el-table-column label="设备" width="280">
        <template #default="{ row }">
          <el-tag v-for="device in (row.devices || '').split(',')" :key="device" size="small" style="margin-right: 4px" v-if="device">
            {{ device }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" @click="openDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingRoom ? '编辑会议室' : '新增会议室'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="会议室名">
          <el-input v-model="form.name" placeholder="请输入会议室名称" />
        </el-form-item>
        <el-form-item label="容纳人数">
          <el-input-number v-model="form.capacity" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="设备">
          <el-checkbox-group v-model="selectedDevices">
            <el-checkbox label="投影" />
            <el-checkbox label="电视" />
            <el-checkbox label="白板" />
            <el-checkbox label="会议电话" />
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const authStore = useAuthStore()
const rooms = ref([])
const search = ref('')
const dialogVisible = ref(false)
const editingRoom = ref(null)
const saving = ref(false)
const form = ref({
  name: '',
  capacity: 10,
  devices: '',
  status: 1
})
const selectedDevices = ref([])

const filteredRooms = computed(() => {
  if (!search.value) return rooms.value
  return rooms.value.filter(r => r.name.includes(search.value))
})

const fetchRooms = async () => {
  try {
    const res = await authStore.api.get('/rooms')
    rooms.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const openDialog = (room = null) => {
  editingRoom.value = room
  if (room) {
    form.value = { ...room }
    selectedDevices.value = (room.devices || '').split(',').filter(d => d)
  } else {
    form.value = { name: '', capacity: 10, devices: '', status: 1 }
    selectedDevices.value = []
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入会议室名称')
    return
  }
  form.value.devices = selectedDevices.value.join(',')
  saving.value = true
  try {
    if (editingRoom.value) {
      await authStore.api.put(`/rooms/${editingRoom.value.id}`, form.value)
    } else {
      await authStore.api.post('/rooms', form.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchRooms()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = async (room) => {
  try {
    await ElMessageBox.confirm('确定要删除该会议室吗?', '提示', { type: 'warning' })
    await authStore.api.delete(`/rooms/${room.id}`)
    ElMessage.success('删除成功')
    fetchRooms()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

onMounted(() => {
  fetchRooms()
})
</script>

<style scoped>
.rooms {
  background: #1e293b;
  border-radius: 12px;
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h2 {
  color: #e2e8f0;
  font-size: 20px;
}

.actions {
  display: flex;
  gap: 12px;
}
</style>
