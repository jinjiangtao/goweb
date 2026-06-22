<template>
  <div class="devices">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>设备管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="showAddDialog = true">
              <el-icon><Plus /></el-icon>新增设备
            </el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true" :model="query" class="filter-form">
        <el-form-item label="类型">
          <el-select v-model="query.device_type" placeholder="全部" clearable style="width: 140px;">
            <el-option label="电脑" value="computer" />
            <el-option label="打印机" value="printer" />
            <el-option label="网络设备" value="network" />
            <el-option label="投影仪" value="projector" />
            <el-option label="空调" value="aircon" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="query.area" placeholder="全部" clearable style="width: 140px;">
            <el-option label="教学楼A区" value="教学楼A区" />
            <el-option label="教学楼B区" value="教学楼B区" />
            <el-option label="实验楼" value="实验楼" />
            <el-option label="办公楼" value="办公楼" />
            <el-option label="图书馆" value="图书馆" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="设备名称/编号" clearable style="width: 200px;" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadDevices">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="devices" v-loading="loading" stripe>
        <el-table-column prop="device_code" label="设备编号" width="140" />
        <el-table-column prop="device_name" label="设备名称" min-width="180" />
        <el-table-column label="设备类型" width="100">
          <template #default="{ row }">{{ deviceTypeText(row.device_type) }}</template>
        </el-table-column>
        <el-table-column prop="area" label="区域" width="120" />
        <el-table-column prop="location" label="位置" min-width="180" show-overflow-tooltip />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'normal' ? 'success' : 'danger'" size="small">
              {{ row.status === 'normal' ? '正常' : '故障' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadDevices"
          @current-change="loadDevices"
        />
      </div>
    </el-card>

    <el-dialog v-model="showAddDialog" :title="isEdit ? '编辑设备' : '新增设备'" width="550px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="设备编号" prop="device_code">
          <el-input v-model="form.device_code" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="设备名称" prop="device_name">
          <el-input v-model="form.device_name" />
        </el-form-item>
        <el-form-item label="设备类型" prop="device_type">
          <el-select v-model="form.device_type" style="width: 100%;">
            <el-option label="电脑" value="computer" />
            <el-option label="打印机" value="printer" />
            <el-option label="网络设备" value="network" />
            <el-option label="投影仪" value="projector" />
            <el-option label="空调" value="aircon" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="form.area" style="width: 100%;" clearable>
            <el-option label="教学楼A区" value="教学楼A区" />
            <el-option label="教学楼B区" value="教学楼B区" />
            <el-option label="实验楼" value="实验楼" />
            <el-option label="办公楼" value="办公楼" />
            <el-option label="图书馆" value="图书馆" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="位置" prop="location">
          <el-input v-model="form.location" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio value="normal">正常</el-radio>
            <el-radio value="fault">故障</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { createDevice, getDevices, updateDevice, deleteDevice } from '@/api'

const loading = ref(false)
const devices = ref([])
const total = ref(0)
const showAddDialog = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const editId = ref(null)

const query = reactive({
  page: 1,
  page_size: 20,
  device_type: '',
  area: '',
  keyword: ''
})

const form = reactive({
  device_code: '',
  device_name: '',
  device_type: '',
  area: '',
  location: '',
  status: 'normal',
  description: ''
})

const rules = {
  device_code: [{ required: true, message: '请输入设备编号', trigger: 'blur' }],
  device_name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  device_type: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
  location: [{ required: true, message: '请输入位置', trigger: 'blur' }]
}

const loadDevices = async () => {
  loading.value = true
  try {
    const res = await getDevices(query)
    devices.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  Object.assign(form, {
    device_code: '',
    device_name: '',
    device_type: '',
    area: '',
    location: '',
    status: 'normal',
    description: ''
  })
  isEdit.value = false
  editId.value = null
  formRef.value?.resetFields()
}

const handleEdit = (row) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(form, row)
  showAddDialog.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          await updateDevice(editId.value, form)
          ElMessage.success('更新成功')
        } else {
          await createDevice(form)
          ElMessage.success('创建成功')
        }
        showAddDialog.value = false
        resetForm()
        loadDevices()
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除此设备吗？', '提示', { type: 'warning' })
    await deleteDevice(row.id)
    ElMessage.success('删除成功')
    loadDevices()
  } catch (e) {}
}

const handleReset = () => {
  query.page = 1
  query.device_type = ''
  query.area = ''
  query.keyword = ''
  loadDevices()
}

const deviceTypeText = (t) => ({
  computer: '电脑', printer: '打印机', network: '网络设备',
  projector: '投影仪', aircon: '空调', other: '其他'
}[t] || t)

const formatTime = (t) => t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
}

.filter-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
