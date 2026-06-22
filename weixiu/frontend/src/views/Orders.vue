<template>
  <div class="orders">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单列表</span>
          <div class="header-actions">
            <el-button v-if="!userStore.isTechnician" type="primary" @click="$router.push('/orders/new')">
              <el-icon><Plus /></el-icon>提交工单
            </el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true" :model="query" class="filter-form">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width: 140px;">
            <el-option label="待处理" value="pending" />
            <el-option label="已派单" value="assigned" />
            <el-option label="维修中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="已驳回" value="rejected" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item label="设备类型">
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
        <el-form-item label="优先级">
          <el-select v-model="query.priority" placeholder="全部" clearable style="width: 120px;">
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="标题/编号/描述" clearable style="width: 200px;" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadOrders">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="orders" v-loading="loading" stripe>
        <el-table-column prop="order_no" label="工单编号" width="180">
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/orders/${row.id}`)">{{ row.order_no }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="故障标题" min-width="200" show-overflow-tooltip />
        <el-table-column label="设备类型" width="100">
          <template #default="{ row }">{{ deviceTypeText(row.device_type) }}</template>
        </el-table-column>
        <el-table-column prop="area" label="区域" width="120" />
        <el-table-column label="优先级" width="80">
          <template #default="{ row }">
            <el-tag :type="priorityType(row.priority)" size="small">{{ priorityText(row.priority) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <div class="status-cell">
              <el-tag :type="statusType(row.status)" effect="dark" size="small">
                {{ statusText(row.status) }}
              </el-tag>
              <el-icon v-if="row.is_overdue" class="overdue-icon" color="#f56c6c"><Warning /></el-icon>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="报修人" width="100">
          <template #default="{ row }">{{ row.user?.real_name || row.user?.username }}</template>
        </el-table-column>
        <el-table-column label="运维人员" width="100">
          <template #default="{ row }">{{ row.technician?.real_name || '-' }}</template>
        </el-table-column>
        <el-table-column label="提交时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/orders/${row.id}`)">查看</el-button>
            <el-button
              v-if="userStore.isAdmin && (row.status === 'pending' || row.status === 'rejected')"
              type="warning"
              link
              @click="showAssignDialog(row)"
            >派单</el-button>
            <el-button
              v-if="userStore.user?.role !== 'technician' && (row.status === 'pending' || row.status === 'assigned') && row.user_id === userStore.user?.id"
              type="danger"
              link
              @click="handleCancel(row)"
            >取消</el-button>
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
          @size-change="loadOrders"
          @current-change="loadOrders"
        />
      </div>
    </el-card>

    <el-dialog v-model="assignVisible" title="派单" width="500px">
      <el-form label-width="100px">
        <el-form-item label="选择运维">
          <el-select v-model="assignTech" placeholder="请选择运维人员" style="width: 100%;">
            <el-option v-for="tech in technicians" :key="tech.id" :label="`${tech.real_name} (${tech.area})`" :value="tech.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAutoAssign">智能派单</el-button>
        <el-button type="success" @click="handleAssign">确认派单</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Warning } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { getOrders, getTechnicians, assignOrder, autoAssignOrder, cancelOrder } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)
const orders = ref([])
const total = ref(0)
const technicians = ref([])
const assignVisible = ref(false)
const assignTech = ref(null)
const currentOrder = ref(null)

const query = reactive({
  page: 1,
  page_size: 20,
  status: '',
  device_type: '',
  area: '',
  priority: '',
  keyword: ''
})

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getOrders(query)
    orders.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const loadTechnicians = async () => {
  try {
    technicians.value = await getTechnicians()
  } catch (e) {}
}

const showAssignDialog = (row) => {
  currentOrder.value = row
  assignTech.value = null
  loadTechnicians()
  assignVisible.value = true
}

const handleAssign = async () => {
  if (!assignTech.value) {
    ElMessage.warning('请选择运维人员')
    return
  }
  await assignOrder(currentOrder.value.id, { technician_id: assignTech.value })
  ElMessage.success('派单成功')
  assignVisible.value = false
  loadOrders()
}

const handleAutoAssign = async () => {
  await autoAssignOrder(currentOrder.value.id)
  ElMessage.success('智能派单成功')
  assignVisible.value = false
  loadOrders()
}

const handleCancel = async (row) => {
  try {
    await ElMessageBox.confirm('确定取消此工单吗？', '提示', { type: 'warning' })
    await cancelOrder(row.id)
    ElMessage.success('工单已取消')
    loadOrders()
  } catch (e) {}
}

const handleReset = () => {
  query.page = 1
  query.status = ''
  query.device_type = ''
  query.area = ''
  query.priority = ''
  query.keyword = ''
  loadOrders()
}

const statusText = (s) => ({
  pending: '待处理', assigned: '已派单', processing: '维修中',
  completed: '已完成', rejected: '已驳回', cancelled: '已取消'
}[s] || s)

const statusType = (s) => ({
  pending: 'info', assigned: 'warning', processing: 'primary',
  completed: 'success', rejected: 'danger', cancelled: 'info'
}[s] || '')

const priorityText = (p) => ({ low: '低', medium: '中', high: '高', urgent: '紧急' }[p] || p)
const priorityType = (p) => ({ low: 'info', medium: '', high: 'warning', urgent: 'danger' }[p] || '')

const deviceTypeText = (t) => ({
  computer: '电脑', printer: '打印机', network: '网络设备',
  projector: '投影仪', aircon: '空调', other: '其他'
}[t] || t)

const formatTime = (t) => t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'

onMounted(() => {
  loadOrders()
  if (userStore.isAdmin) loadTechnicians()
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

.status-cell {
  display: flex;
  align-items: center;
  gap: 5px;
}

.overdue-icon {
  font-size: 18px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
