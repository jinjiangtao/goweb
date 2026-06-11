<template>
  <div class="product-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>产品管理</span>
          <div class="header-right">
            <el-input
              v-model="searchKeyword"
              placeholder="请输入产品名称"
              clearable
              style="width: 200px; margin-right: 10px"
              @clear="handleSearch"
            >
              <template #append>
                <el-button :icon="Search" @click="handleSearch" />
              </template>
            </el-input>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              新增产品
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="tableData" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="产品名称" width="200" />
        <el-table-column prop="code" label="产品编码" width="150" />
        <el-table-column prop="price" label="价格" width="120">
          <template #default="{ row }">
            ¥{{ row.price?.toFixed(2) || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="spec" label="规格" />
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
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

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="resetForm"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="产品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入产品名称" />
        </el-form-item>
        <el-form-item label="产品编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入产品编码" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="规格" prop="spec">
          <el-input v-model="form.spec" type="textarea" placeholder="请输入规格" />
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
import { Plus, Search } from '@element-plus/icons-vue'
import { getProducts, getProduct, createProduct, updateProduct, deleteProduct } from '../api/product'

const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const loading = ref(false)
const formRef = ref(null)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')

const form = reactive({
  id: null,
  name: '',
  code: '',
  price: 0,
  spec: ''
})

const rules = {
  name: [{ required: true, message: '请输入产品名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入产品编码', trigger: 'blur' }]
}

const formatDate = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleString('zh-CN')
}

const fetchData = async () => {
  try {
    const params = { page: currentPage.value, pageSize: pageSize.value }
    if (searchKeyword.value) {
      params.name = searchKeyword.value
    }
    const res = await getProducts(params)
    tableData.value = res.data?.list || []
    total.value = res.data?.total || 0
    console.log('产品列表数据:', tableData.value)
  } catch (error) {
    ElMessage.error('获取产品列表失败')
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

const handleAdd = () => {
  dialogTitle.value = '新增产品'
  isEdit.value = false
  dialogVisible.value = true
}

const handleEdit = async (row) => {
  try {
    const res = await getProduct(row.id)
    const product = res.data
    Object.assign(form, {
      id: product.id,
      name: product.name,
      code: product.code,
      price: product.price,
      spec: product.spec
    })
    dialogTitle.value = '编辑产品'
    isEdit.value = true
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取产品信息失败')
  }
}

const handleDelete = async (row) => {
  try {
    if (!row || !row.id) {
      ElMessage.error('产品ID不存在，无法删除')
      return
    }
    await ElMessageBox.confirm('确定要删除该产品吗？', '提示', { type: 'warning' })
    await deleteProduct(row.id)
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
          await updateProduct(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createProduct(form)
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
  form.name = ''
  form.code = ''
  form.price = 0
  form.spec = ''
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

.header-right {
  display: flex;
  align-items: center;
}
</style>
