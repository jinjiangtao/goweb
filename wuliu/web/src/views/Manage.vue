<template>
  <div class="manage-container">
    <div class="sidebar">
      <div class="logo">
        <el-icon><Van /></el-icon>
        <span>物流管理系统</span>
      </div>
      <div class="menu">
        <div
          class="menu-item"
          :class="{ active: activeTab === 'orders' }"
          @click="activeTab = 'orders'"
        >
          <el-icon><List /></el-icon>
          <span>订单管理</span>
        </div>
        <div
          class="menu-item"
          :class="{ active: activeTab === 'nodes' }"
          @click="activeTab = 'nodes'"
        >
          <el-icon><Connection /></el-icon>
          <span>轨迹管理</span>
        </div>
      </div>
    </div>

    <div class="main-content">
      <div class="header">
        <div class="header-left">
          <el-button @click="goHome">
            <el-icon><ArrowLeft /></el-icon>
            返回首页
          </el-button>
          <h2>{{ activeTab === 'orders' ? '订单管理' : '轨迹管理' }}</h2>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="showOrderDialog = true" v-if="activeTab === 'orders'">
            <el-icon><Plus /></el-icon>
            新增订单
          </el-button>
        </div>
      </div>

      <div class="content-body" v-if="activeTab === 'orders'">
        <div class="search-bar">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索订单号、收发件人、地址、货物名称"
            clearable
            style="width: 360px"
            @keyup.enter="handleSearchOrder"
            @clear="handleClearSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="handleSearchOrder">查询</el-button>
        </div>

        <div class="table-card">
          <el-table :data="orders" v-loading="loading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="order_no" label="订单号" min-width="150" />
            <el-table-column prop="sender" label="发货人" width="120" />
            <el-table-column prop="receiver" label="收件人" width="120" />
            <el-table-column prop="origin" label="发货地" min-width="150" />
            <el-table-column prop="destination" label="目的地" min-width="150" />
            <el-table-column prop="goods_name" label="货物" width="120" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button text type="primary" @click="viewOrder(row)">查看</el-button>
                <el-button text type="primary" @click="editOrder(row)">编辑</el-button>
                <el-button text type="danger" @click="handleDeleteOrder(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="page"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="loadOrders"
              @current-change="loadOrders"
            />
          </div>
        </div>
      </div>

      <div class="content-body" v-if="activeTab === 'nodes'">
        <div class="search-bar">
          <el-input
            v-model="nodeSearchOrderNo"
            placeholder="输入订单号查看轨迹"
            style="width: 300px"
            @keyup.enter="searchNodes"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="searchNodes">查询</el-button>
          <el-button type="success" @click="showNodeDialog = true" :disabled="!currentOrder">
            <el-icon><Plus /></el-icon>
            添加节点
          </el-button>
        </div>

        <div v-if="currentOrder" class="order-summary">
          <span class="label">当前订单：</span>
          <span class="order-no">{{ currentOrder.order_no }}</span>
          <el-tag :type="getStatusType(currentOrder.status)" style="margin-left: 15px">
            {{ getStatusText(currentOrder.status) }}
          </el-tag>
        </div>

        <div class="table-card" v-if="currentOrder">
          <el-table :data="nodes" v-loading="nodesLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="node_name" label="节点名称" width="150" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="location" label="位置" min-width="150" />
            <el-table-column prop="description" label="描述" min-width="200" />
            <el-table-column prop="operator" label="操作人" width="120" />
            <el-table-column label="异常" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.is_abnormal" type="danger" size="small">是</el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="occurred_at" label="发生时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.occurred_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button text type="primary" @click="editNode(row)">编辑</el-button>
                <el-button text type="danger" @click="handleDeleteNode(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <el-empty v-else description="请输入订单号查询轨迹节点" />
      </div>
    </div>

    <el-dialog
      v-model="showOrderDialog"
      :title="editingOrder ? '编辑订单' : '新增订单'"
      width="600px"
      @closed="resetOrderForm"
    >
      <el-form :model="orderForm" label-width="100px">
        <el-form-item label="订单号">
          <el-input v-model="orderForm.order_no" :disabled="!!editingOrder" />
        </el-form-item>
        <el-form-item label="发货人">
          <el-input v-model="orderForm.sender" />
        </el-form-item>
        <el-form-item label="收件人">
          <el-input v-model="orderForm.receiver" />
        </el-form-item>
        <el-form-item label="发货地">
          <el-input v-model="orderForm.origin" />
        </el-form-item>
        <el-form-item label="目的地">
          <el-input v-model="orderForm.destination" />
        </el-form-item>
        <el-form-item label="货物名称">
          <el-input v-model="orderForm.goods_name" />
        </el-form-item>
        <el-form-item label="重量(kg)">
          <el-input-number v-model="orderForm.weight" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="状态" v-if="editingOrder">
          <el-select v-model="orderForm.status" style="width: 100%">
            <el-option label="待揽收" value="pending" />
            <el-option label="已揽收" value="picked" />
            <el-option label="运输中" value="in_transit" />
            <el-option label="已到达" value="arrived" />
            <el-option label="已签收" value="delivered" />
            <el-option label="异常" value="exception" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showOrderDialog = false">取消</el-button>
        <el-button type="primary" @click="submitOrder" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="showNodeDialog"
      :title="editingNode ? '编辑节点' : '添加节点'"
      width="600px"
      @closed="resetNodeForm"
    >
      <el-form :model="nodeForm" label-width="120px">
        <el-form-item label="节点名称">
          <el-input v-model="nodeForm.node_name" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="nodeForm.status" style="width: 100%">
            <el-option label="待处理" value="pending" />
            <el-option label="已揽收" value="picked" />
            <el-option label="运输中" value="in_transit" />
            <el-option label="已到达" value="arrived" />
            <el-option label="已签收" value="delivered" />
            <el-option label="异常" value="exception" />
          </el-select>
        </el-form-item>
        <el-form-item label="位置">
          <el-input v-model="nodeForm.location" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="nodeForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="操作人">
          <el-input v-model="nodeForm.operator" />
        </el-form-item>
        <el-form-item label="是否异常">
          <el-switch v-model="nodeForm.is_abnormal" />
        </el-form-item>
        <el-form-item label="异常原因" v-if="nodeForm.is_abnormal">
          <el-input v-model="nodeForm.abnormal_reason" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="发生时间">
          <el-date-picker
            v-model="nodeForm.occurred_at"
            type="datetime"
            placeholder="选择日期时间"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showNodeDialog = false">取消</el-button>
        <el-button type="primary" @click="submitNode" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getOrders,
  createOrder,
  updateOrder,
  deleteOrder,
  searchOrder,
  getNodes,
  addNode,
  updateNode,
  deleteNode
} from '../api'
import {
  Van,
  List,
  Connection,
  Plus,
  Search,
  ArrowLeft
} from '@element-plus/icons-vue'

const router = useRouter()
const activeTab = ref('orders')
const loading = ref(false)
const submitting = ref(false)

const orders = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')

const showOrderDialog = ref(false)
const editingOrder = ref(null)
const orderForm = ref({
  order_no: '',
  sender: '',
  receiver: '',
  origin: '',
  destination: '',
  goods_name: '',
  weight: 0,
  status: 'pending'
})

const nodeSearchOrderNo = ref('')
const currentOrder = ref(null)
const nodes = ref([])
const nodesLoading = ref(false)
const showNodeDialog = ref(false)
const editingNode = ref(null)
const nodeForm = ref({
  node_name: '',
  status: 'in_transit',
  location: '',
  description: '',
  operator: '',
  is_abnormal: false,
  abnormal_reason: '',
  occurred_at: new Date()
})

const getStatusText = (status) => {
  const statusMap = {
    pending: '待揽收',
    picked: '已揽收',
    in_transit: '运输中',
    arrived: '已到达',
    delivered: '已签收',
    exception: '异常'
  }
  return statusMap[status] || status
}

const getStatusType = (status) => {
  const typeMap = {
    pending: 'info',
    picked: 'success',
    in_transit: 'primary',
    arrived: 'warning',
    delivered: 'success',
    exception: 'danger'
  }
  return typeMap[status] || 'info'
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    const res = await getOrders(params)
    orders.value = res.data.data
    total.value = res.data.total
  } catch (err) {
    ElMessage.error('加载订单失败')
  } finally {
    loading.value = false
  }
}

const handleSearchOrder = () => {
  page.value = 1
  loadOrders()
}

const handleClearSearch = () => {
  searchKeyword.value = ''
  page.value = 1
  loadOrders()
}

const viewOrder = (row) => {
  router.push(`/order/${row.id}`)
}

const editOrder = (row) => {
  editingOrder.value = row
  orderForm.value = { ...row }
  showOrderDialog.value = true
}

const resetOrderForm = () => {
  editingOrder.value = null
  orderForm.value = {
    order_no: '',
    sender: '',
    receiver: '',
    origin: '',
    destination: '',
    goods_name: '',
    weight: 0,
    status: 'pending'
  }
}

const submitOrder = async () => {
  if (!orderForm.value.order_no) {
    ElMessage.warning('请输入订单号')
    return
  }
  submitting.value = true
  try {
    if (editingOrder.value) {
      await updateOrder(editingOrder.value.id, orderForm.value)
      ElMessage.success('订单更新成功')
    } else {
      await createOrder(orderForm.value)
      ElMessage.success('订单创建成功')
    }
    showOrderDialog.value = false
    loadOrders()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDeleteOrder = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除订单 ${row.order_no} 吗？`, '提示', {
      type: 'warning'
    })
    await deleteOrder(row.id)
    ElMessage.success('删除成功')
    loadOrders()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const searchNodes = async () => {
  if (!nodeSearchOrderNo.value.trim()) {
    ElMessage.warning('请输入订单号')
    return
  }
  nodesLoading.value = true
  try {
    const res = await searchOrder(nodeSearchOrderNo.value.trim())
    currentOrder.value = res.data.data
    nodes.value = res.data.data.nodes || []
  } catch (err) {
    currentOrder.value = null
    nodes.value = []
    ElMessage.error('未找到该订单')
  } finally {
    nodesLoading.value = false
  }
}

const editNode = (row) => {
  editingNode.value = row
  nodeForm.value = {
    ...row,
    occurred_at: new Date(row.occurred_at)
  }
  showNodeDialog.value = true
}

const resetNodeForm = () => {
  editingNode.value = null
  nodeForm.value = {
    node_name: '',
    status: 'in_transit',
    location: '',
    description: '',
    operator: '',
    is_abnormal: false,
    abnormal_reason: '',
    occurred_at: new Date()
  }
}

const submitNode = async () => {
  if (!nodeForm.value.node_name) {
    ElMessage.warning('请输入节点名称')
    return
  }
  submitting.value = true
  try {
    const data = {
      ...nodeForm.value,
      occurred_at: nodeForm.value.occurred_at?.toISOString()
    }
    if (editingNode.value) {
      await updateNode(editingNode.value.id, data)
      ElMessage.success('节点更新成功')
    } else {
      await addNode(currentOrder.value.id, data)
      ElMessage.success('节点添加成功')
    }
    showNodeDialog.value = false
    searchNodes()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDeleteNode = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除节点 ${row.node_name} 吗？`, '提示', {
      type: 'warning'
    })
    await deleteNode(row.id)
    ElMessage.success('删除成功')
    searchNodes()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const goHome = () => {
  router.push('/')
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.manage-container {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 240px;
  background: #001529;
  color: white;
  display: flex;
  flex-direction: column;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  border-bottom: 1px solid #1f3a56;
}

.logo .el-icon {
  font-size: 24px;
  color: #409eff;
}

.menu {
  flex: 1;
  padding: 20px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 24px;
  cursor: pointer;
  transition: all 0.3s;
  color: rgba(255, 255, 255, 0.65);
}

.menu-item:hover {
  background: #1f3a56;
  color: white;
}

.menu-item.active {
  background: #409eff;
  color: white;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f0f2f5;
}

.header {
  height: 64px;
  background: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header h2 {
  font-size: 18px;
  margin: 0;
  color: #303133;
}

.content-body {
  flex: 1;
  padding: 24px;
}

.search-bar {
  background: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.table-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.order-summary {
  background: #ecf5ff;
  padding: 15px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.order-summary .label {
  color: #606266;
}

.order-summary .order-no {
  font-weight: 600;
  color: #409eff;
}
</style>
