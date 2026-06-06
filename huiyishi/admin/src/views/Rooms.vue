
&lt;template&gt;
  &lt;div class="rooms"&gt;
    &lt;div class="header"&gt;
      &lt;h2&gt;会议室管理&lt;/h2&gt;
      &lt;div class="actions"&gt;
        &lt;el-input v-model="search" placeholder="搜索会议室" style="width: 240px" clearable /&gt;
        &lt;el-button type="primary" @click="openDialog()"&gt;新增会议室&lt;/el-button&gt;
      &lt;/div&gt;
    &lt;/div&gt;

    &lt;el-table :data="filteredRooms" stripe style="width: 100%"&gt;
      &lt;el-table-column prop="name" label="会议室名称" /&gt;
      &lt;el-table-column prop="capacity" label="容纳人数" width="120" /&gt;
      &lt;el-table-column label="设备" width="280"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-tag v-for="device in (row.devices || '').split(',')" :key="device" size="small" style="margin-right: 4px" v-if="device"&gt;
            {{ device }}
          &lt;/el-tag&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column label="状态" width="100"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-tag :type="row.status === 1 ? 'success' : 'danger'"&gt;
            {{ row.status === 1 ? '启用' : '禁用' }}
          &lt;/el-tag&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
      &lt;el-table-column label="操作" width="180"&gt;
        &lt;template #default="{ row }"&gt;
          &lt;el-button size="small" @click="openDialog(row)"&gt;编辑&lt;/el-button&gt;
          &lt;el-button size="small" type="danger" @click="handleDelete(row)"&gt;删除&lt;/el-button&gt;
        &lt;/template&gt;
      &lt;/el-table-column&gt;
    &lt;/el-table&gt;

    &lt;el-dialog v-model="dialogVisible" :title="editingRoom ? '编辑会议室' : '新增会议室'" width="500px"&gt;
      &lt;el-form :model="form" label-width="80px"&gt;
        &lt;el-form-item label="会议室名"&gt;
          &lt;el-input v-model="form.name" placeholder="请输入会议室名称" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="容纳人数"&gt;
          &lt;el-input-number v-model="form.capacity" :min="1" :max="100" /&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="设备"&gt;
          &lt;el-checkbox-group v-model="selectedDevices"&gt;
            &lt;el-checkbox label="投影" /&gt;
            &lt;el-checkbox label="电视" /&gt;
            &lt;el-checkbox label="白板" /&gt;
            &lt;el-checkbox label="会议电话" /&gt;
          &lt;/el-checkbox-group&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="状态"&gt;
          &lt;el-radio-group v-model="form.status"&gt;
            &lt;el-radio :label="1"&gt;启用&lt;/el-radio&gt;
            &lt;el-radio :label="0"&gt;禁用&lt;/el-radio&gt;
          &lt;/el-radio-group&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
      &lt;template #footer&gt;
        &lt;el-button @click="dialogVisible = false"&gt;取消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="handleSave" :loading="saving"&gt;保存&lt;/el-button&gt;
      &lt;/template&gt;
    &lt;/el-dialog&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
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

const filteredRooms = computed(() =&gt; {
  if (!search.value) return rooms.value
  return rooms.value.filter(r =&gt; r.name.includes(search.value))
})

const fetchRooms = async () =&gt; {
  const res = await authStore.api.get('/rooms')
  rooms.value = res.data
}

const openDialog = (room = null) =&gt; {
  editingRoom.value = room
  if (room) {
    form.value = { ...room }
    selectedDevices.value = (room.devices || '').split(',').filter(d =&gt; d)
  } else {
    form.value = { name: '', capacity: 10, devices: '', status: 1 }
    selectedDevices.value = []
  }
  dialogVisible.value = true
}

const handleSave = async () =&gt; {
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

const handleDelete = async (room) =&gt; {
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

onMounted(() =&gt; {
  fetchRooms()
})
&lt;/script&gt;

&lt;style scoped&gt;
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
&lt;/style&gt;

