<template>
  <div class="purchase-order-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>采购订单</span>
          <div class="header-right">
            <el-input
              v-model="searchParams.orderNo"
              placeholder="订单编号"
              clearable
              style="width: 150px; margin-right: 10px"
              @clear="handleSearch"
            />
            <el-select
              v-model="searchParams.supplierId"
              placeholder="供应商"
              clearable
              style="width: 150px; margin-right: 10px"
              @clear="handleSearch"
            >
              <el-option
                v-for="supplier in suppliers"
                :key="supplier.id || supplier.ID"
                :label="supplier.name"
                :value="supplier.id || supplier.ID"
              />
            </el-select>
            <el-select
              v-model="searchParams.status"
              placeholder="状态"
              clearable
              style="width: 120px; margin-right: 10px"
              @clear="handleSearch"
            >
              <el-option label="草稿" value="draft" />
              <el-option label="已审核" value="approved" />
              <el-option label="部分收货" value="partial_receive" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
            <el-button :icon="Search" @click="handleSearch">搜索</el-button>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增订单
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="tableData" stripe style="width: 100%">
        <el-table-column prop="orderNo" label="订单编号" width="180" />
        <el-table-column prop="supplier.name" label="供应商" width="150">
          <template #default="{ row }">
            {{ row.supplier?.name || '' }}
          </template>
        </el-table-column>
        <el-table-column label="订单金额" width="120">
          <template #default="{ row }">
            ¥{{ (row.totalAmount || 0).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="预计到货" width="120">
          <template #default="{ row }">
            {{ row.expectedDate ? formatDate(row.expectedDate) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column label="创建人" width="120">
          <template #default="{ row }">
            {{ row.createdBy?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt || row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="300">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button
              v-if="row.status === 'draft'"
              type="primary"
              size="small"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <template v-if="row.status === 'draft'">
              <el-button type="success" size="small" @click="handleApprove(row)">
                审核
              </el-button>
            </template>
            <template v-else-if="row.status === 'approved' || row.status === 'partial_receive'">
              <el-button type="success" size="small" @click="handleReceive(row)">
                收货
              </el-button>
              <el-button type="success" size="small" @click="handleComplete(row)">
                完成
              </el-button>
            </template>
            <template v-if="row.status !== 'completed' && row.status !== 'cancelled'">
              <el-button type="danger" size="small" @click="handleCancel(row)">
                取消
              </el-button>
            </template>
            <el-button
              v-if="row.status === 'draft'"
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 新增/编辑订单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="900px"
      @close="resetForm"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="供应商" prop="supplierId">
          <el-select v-model="form.supplierId" placeholder="请选择供应商" style="width: 100%">
            <el-option
              v-for="supplier in suppliers"
              :key="supplier.id || supplier.ID"
              :label="supplier.name"
              :value="supplier.id || supplier.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="预计到货" prop="expectedDate">
          <el-date-picker
            v-model="form.expectedDate"
            type="date"
            placeholder="选择日期"
            style="width: 100%"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
        <el-form-item label="采购明细">
          <div style="width: 100%">
            <el-table :data="form.items" border style="width: 100%">
              <el-table-column label="产品" width="200">
                <template #default="{ row, $index }">
                  <el-select
                    v-model="row.productId"
                    placeholder="选择产品"
                    style="width: 100%"
                    @change="handleProductChange($index)"
                  >
                    <el-option
                      v-for="product in products"
                      :key="product.id || product.ID"
                      :label="product.name"
                      :value="product.id || product.ID"
                    />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="数量" width="120">
                <template #default="{ row, $index }">
                  <el-input-number
                    v-model="row.quantity"
                    :min="0.01"
                    :precision="2"
                    :step="1"
                    style="width: 100%"
                    @change="calculateAmount($index)"
                  />
                </template>
              </el-table-column>
              <el-table-column label="单价" width="120">
                <template #default="{ row, $index }">
                  <el-input-number
                    v-model="row.unitPrice"
                    :min="0"
                    :precision="2"
                    :step="0.01"
                    style="width: 100%"
                    @change="calculateAmount($index)"
                  />
                </template>
              </el-table-column>
              <el-table-column label="金额" width="120">
                <template #default="{ row }">
                  ¥{{ (row.amount || 0).toFixed(2) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80">
                <template #default="{ $index }">
                  <el-button type="danger" size="small" @click="removeItem($index)">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-button
              type="primary"
              style="margin-top: 10px"
              @click="addItem"
            >
              <el-icon><Plus /></el-icon>
              添加产品
            </el-button>
            <div style="margin-top: 10px; text-align: right; font-weight: bold">
              总计: ¥{{ calculateTotal().toFixed(2) }}
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 查看详情对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="订单详情"
      width="900px"
    >
      <el-descriptions :column="2" border v-if="viewOrder">
        <el-descriptions-item label="订单编号">{{ viewOrder.orderNo }}</el-descriptions-item>
        <el-descriptions-item label="供应商">{{ viewOrder.supplier?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">¥{{ (viewOrder.totalAmount || 0).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(viewOrder.status)">
            {{ getStatusLabel(viewOrder.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="预计到货">
          {{ viewOrder.expectedDate ? formatDate(viewOrder.expectedDate) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建人">{{ viewOrder.createdBy?.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(viewOrder.createdAt || viewOrder.CreatedAt) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(viewOrder.updatedAt || viewOrder.UpdatedAt) }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ viewOrder.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <el-table :data="viewOrder?.items || []" border style="width: 100%; margin-top: 20px">
        <el-table-column label="产品名称" width="200">
          <template #default="{ row }">
            {{ row.product?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="规格">
          <template #default="{ row }">
            {{ row.product?.spec || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="数量" width="120">{{ row.quantity }}</el-table-column>
        <el-table-column label="单价" width="120">¥{{ (row.unitPrice || 0).toFixed(2) }}</el-table-column>
        <el-table-column label="金额" width="120">¥{{ (row.amount || 0).toFixed(2) }}</el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { getPurchaseOrders, getPurchaseOrder, createPurchaseOrder, updatePurchaseOrder, updatePurchaseOrderStatus, deletePurchaseOrder } from '../api/purchaseOrder'
import { getSuppliers } from '../api/supplier'
import { getProducts } from '../api/product'

const tableData = ref([])
const dialogVisible = ref(false)
const viewDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const loading = ref(false)
const formRef = ref(null)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const suppliers = ref([])
const products = ref([])
const viewOrder = ref(null)

const searchParams = reactive({
  orderNo: '',
  supplierId: '',
  status: ''
})

const form = reactive({
  id: null,
  supplierId: null,
  items: [],
  expectedDate: '',
  remark: ''
})

const rules = {
  supplierId: [{ required: true, message: '请选择供应商', trigger: 'change' }]
}

const formatDate = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleString('zh-CN')
}

const getStatusType = (status) => {
  const typeMap = {
    draft: 'info',
    approved: 'primary',
    partial_receive: 'warning',
    completed: 'success',
    cancelled: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusLabel = (status) => {
  const labelMap = {
    draft: '草稿',
    approved: '已审核',
    partial_receive: '部分收货',
    completed: '已完成',
    cancelled: '已取消'
  }
  return labelMap[status] || status
}

const fetchData = async () => {
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    if (searchParams.orderNo) {
      params.orderNo = searchParams.orderNo
    }
    if (searchParams.supplierId) {
      params.supplierId = searchParams.supplierId
    }
    if (searchParams.status) {
      params.status = searchParams.status
    }
    const res = await getPurchaseOrders(params)
    tableData.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    ElMessage.error('获取订单列表失败')
  }
}

const fetchSuppliers = async () => {
  try {
    const res = await getSuppliers({ page: 1, pageSize: 1000 })
    suppliers.value = res.data?.list || []
  } catch (error) {
    console.error('获取供应商列表失败', error)
  }
}

const fetchProducts = async () => {
  try {
    const res = await getProducts({ page: 1, pageSize: 1000 })
    products.value = res.data?.list || []
  } catch (error) {
    console.error('获取产品列表失败', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

const handleAdd = () => {
  dialogTitle.value = '新增采购订单'
  isEdit.value = false
  form.items = [{ productId: null, quantity: 1, unitPrice: 0, amount: 0 }]
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  const rowId = row.id || row.ID
  try {
    const res = await getPurchaseOrder(rowId)
    const order = res.data
    form.id = order.id || order.ID
    form.supplierId = order.supplierId
    form.expectedDate = order.expectedDate ? new Date(order.expectedDate).toISOString().split('T')[0] : ''
    form.remark = order.remark || ''
    form.items = (order.items || []).map(item => ({
      productId: item.productId,
      quantity: item.quantity,
      unitPrice: item.unitPrice,
      amount: item.amount
    }))
    dialogTitle.value = '编辑采购订单'
    isEdit.value = true
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取订单信息失败')
  }
}

const handleView = async (row) => {
  const rowId = row.id || row.ID
  try {
    const res = await getPurchaseOrder(rowId)
    viewOrder.value = res.data
    viewDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取订单详情失败')
  }
}

const handleDelete = async (row) => {
  const rowId = row.id || row.ID
  try {
    await ElMessageBox.confirm('确定要删除该订单吗？', '提示', { type: 'warning' })
    await deletePurchaseOrder(rowId)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

const handleStatusChange = async (row, status) => {
  const rowId = row.id || row.ID
  try {
    await updatePurchaseOrderStatus(rowId, { status })
    ElMessage.success('状态更新成功')
    fetchData()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '状态更新失败')
  }
}

const handleApprove = (row) => {
  ElMessageBox.confirm('确定要审核该订单吗？', '提示', { type: 'warning' })
    .then(() => handleStatusChange(row, 'approved'))
    .catch(() => {})
}

const handleReceive = (row) => {
  ElMessageBox.confirm('确定要记录收货吗？', '提示', { type: 'warning' })
    .then(() => handleStatusChange(row, 'partial_receive'))
    .catch(() => {})
}

const handleComplete = (row) => {
  ElMessageBox.confirm('确定要完成该订单吗？', '提示', { type: 'warning' })
    .then(() => handleStatusChange(row, 'completed'))
    .catch(() => {})
}

const handleCancel = (row) => {
  ElMessageBox.confirm('确定要取消该订单吗？', '提示', { type: 'warning' })
    .then(() => handleStatusChange(row, 'cancelled'))
    .catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      if (form.items.length === 0) {
        ElMessage.warning('请至少添加一个产品')
        return
      }
      for (const item of form.items) {
        if (!item.productId) {
          ElMessage.warning('请选择所有产品')
          return
        }
        if (!item.quantity || item.quantity <= 0) {
          ElMessage.warning('请输入有效的数量')
          return
        }
      }
      
      loading.value = true
      try {
        const data = {
          supplierId: form.supplierId,
          items: form.items,
          expectedDate: form.expectedDate || null,
          remark: form.remark
        }
        
        if (isEdit.value) {
          await updatePurchaseOrder(form.id, data)
          ElMessage.success('更新成功')
        } else {
          await createPurchaseOrder(data)
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

const handleSizeChange = (size) => {
  pageSize.value = size
  fetchData()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchData()
}

const resetForm = () => {
  form.id = null
  form.supplierId = null
  form.items = []
  form.expectedDate = ''
  form.remark = ''
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const addItem = () => {
  form.items.push({
    productId: null,
    quantity: 1,
    unitPrice: 0,
    amount: 0
  })
}

const removeItem = (index) => {
  form.items.splice(index, 1)
}

const handleProductChange = (index) => {
  const item = form.items[index]
  const product = products.value.find(p => (p.id || p.ID) === item.productId)
  if (product) {
    item.unitPrice = product.price || product.Price || 0
    calculateAmount(index)
  }
}

const calculateAmount = (index) => {
  const item = form.items[index]
  item.amount = (item.quantity || 0) * (item.unitPrice || 0)
}

const calculateTotal = () => {
  return form.items.reduce((sum, item) => sum + (item.amount || 0), 0)
}

onMounted(() => {
  fetchData()
  fetchSuppliers()
  fetchProducts()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
}
</style>
