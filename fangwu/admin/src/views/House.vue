<template>
  <div class="house-page">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="标题">
          <el-input v-model="searchForm.title" placeholder="搜索标题" clearable @clear="fetchList" />
        </el-form-item>
        <el-form-item label="小区">
          <el-input v-model="searchForm.community" placeholder="搜索小区" clearable @clear="fetchList" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 140px" @change="fetchList">
            <el-option label="上架" value="1" />
            <el-option label="下架" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="推荐">
          <el-select v-model="searchForm.is_recommended" placeholder="全部" clearable style="width: 140px" @change="fetchList">
            <el-option label="已推荐" value="1" />
            <el-option label="未推荐" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="action-bar">
      <el-button type="primary" @click="openForm(null)">发布房源</el-button>
    </div>

    <el-table :data="list" border stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip />
      <el-table-column prop="community" label="小区" width="120" show-overflow-tooltip />
      <el-table-column prop="district" label="区域" width="80" />
      <el-table-column prop="room_type" label="户型" width="80" />
      <el-table-column prop="area" label="面积(㎡)" width="90" />
      <el-table-column prop="price" label="租金(元/月)" width="110" />
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="推荐" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_recommended ? 'warning' : 'info'">
            {{ row.is_recommended ? '推荐' : '普通' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="发布时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openForm(row)">编辑</el-button>
          <el-button size="small" :type="row.status === 1 ? 'warning' : 'success'" @click="handleToggleStatus(row)">
            {{ row.status === 1 ? '下架' : '上架' }}
          </el-button>
          <el-button size="small" :type="row.is_recommended ? 'info' : 'warning'" @click="handleToggleRecommend(row)">
            {{ row.is_recommended ? '取消推荐' : '设为推荐' }}
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>

    <HouseForm
      v-if="formVisible"
      :visible="formVisible"
      :edit-data="editData"
      @close="formVisible = false"
      @saved="onSaved"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getHouses, deleteHouse, toggleHouseStatus, toggleHouseRecommend } from '../api/house'
import HouseForm from '../components/HouseForm.vue'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const formVisible = ref(false)
const editData = ref(null)

const searchForm = reactive({
  title: '',
  community: '',
  status: '',
  is_recommended: '',
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (searchForm.title) params.title = searchForm.title
    if (searchForm.community) params.community = searchForm.community
    if (searchForm.status !== '') params.status = searchForm.status
    if (searchForm.is_recommended !== '') params.is_recommended = searchForm.is_recommended

    const res = await getHouses(params)
    list.value = res.data.list || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchForm.title = ''
  searchForm.community = ''
  searchForm.status = ''
  searchForm.is_recommended = ''
  page.value = 1
  fetchList()
}

const openForm = (row) => {
  editData.value = row
  formVisible.value = true
}

const onSaved = () => {
  formVisible.value = false
  fetchList()
}

const handleToggleStatus = async (row) => {
  const action = row.status === 1 ? '下架' : '上架'
  await ElMessageBox.confirm(`确定要${action}该房源吗？`, '提示', { type: 'warning' })
  const res = await toggleHouseStatus(row.id)
  ElMessage.success(`${action}成功`)
  fetchList()
}

const handleToggleRecommend = async (row) => {
  const res = await toggleHouseRecommend(row.id)
  ElMessage.success(row.is_recommended ? '已取消推荐' : '已设为推荐')
  fetchList()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该房源吗？删除后不可恢复。', '警告', { type: 'warning' })
  await deleteHouse(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.house-page {
  padding: 0;
}

.search-bar {
  background: #fff;
  padding: 20px 20px 0;
  border-radius: 4px;
  margin-bottom: 16px;
}

.action-bar {
  margin-bottom: 16px;
}

.pagination-bar {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
