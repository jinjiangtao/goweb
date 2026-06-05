<template>
  <div class="school-list">
    <div class="toolbar">
      <div class="filter-area">
        <el-input
          v-model="keyword"
          placeholder="按学校名称搜索"
          class="filter-input"
          @keyup.enter="loadData"
        />
        <el-button type="primary" @click="loadData">搜索</el-button>
        <el-button @click="resetFilter">重置</el-button>
      </div>
      <div class="action-area">
        <el-button type="success" @click="showAddDialog = true">添加学校</el-button>
      </div>
    </div>

    <el-table :data="tableData" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="学校名称" />
      <el-table-column prop="description" label="学校简介" />
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button
            type="text"
            @click="openEditDialog(scope.row)"
          >编辑</el-button>
          <el-button
            type="text"
            @click="handleDelete(scope.row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="pagination.page"
      :page-sizes="[10, 20, 50]"
      :page-size="pagination.pageSize"
      :total="pagination.total"
      layout="total, sizes, prev, pager, next, jumper"
    />

    <el-dialog title="添加学校" v-model="showAddDialog" width="400px">
      <el-form :model="addForm" label-width="80px">
        <el-form-item label="学校名称" prop="name">
          <el-input v-model="addForm.name" />
        </el-form-item>
        <el-form-item label="学校简介" prop="description">
          <el-input v-model="addForm.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddSubmit">确认添加</el-button>
      </template>
    </el-dialog>

    <el-dialog title="编辑学校" v-model="showEditDialog" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="学校名称" prop="name">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="学校简介" prop="description">
          <el-input v-model="editForm.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import axios from '../utils/axios';
import { ElMessage, ElMessageBox } from 'element-plus';

const keyword = ref('');
const tableData = ref([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});
const showAddDialog = ref(false);
const showEditDialog = ref(false);
const addForm = reactive({
  name: '',
  description: ''
});
const editForm = reactive({
  id: '',
  name: '',
  description: ''
});

const loadData = async () => {
  try {
    const response = await axios.get('/api/admin/schools', {
      params: {
        page: pagination.page,
        pageSize: pagination.pageSize,
        keyword: keyword.value
      }
    });
    tableData.value = response.data.list;
    pagination.total = response.data.total;
  } catch (error) {
    ElMessage.error('获取数据失败');
  }
};

const resetFilter = () => {
  keyword.value = '';
  pagination.page = 1;
  loadData();
};

const handleSizeChange = (val) => {
  pagination.pageSize = val;
  pagination.page = 1;
  loadData();
};

const handleCurrentChange = (val) => {
  pagination.page = val;
  loadData();
};

const openEditDialog = (row) => {
  editForm.id = row.id;
  editForm.name = row.name;
  editForm.description = row.description;
  showEditDialog.value = true;
};

const handleAddSubmit = async () => {
  if (!addForm.name || !addForm.description) {
    ElMessage.error('请填写完整信息');
    return;
  }
  try {
    await axios.post('/api/admin/schools', {
      name: addForm.name,
      description: addForm.description
    });
    ElMessage.success('添加成功');
    showAddDialog.value = false;
    addForm.name = '';
    addForm.description = '';
    loadData();
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '添加失败');
  }
};

const handleEditSubmit = async () => {
  if (!editForm.name || !editForm.description) {
    ElMessage.error('请填写完整信息');
    return;
  }
  try {
    await axios.put(`/api/admin/schools/${editForm.id}`, {
      name: editForm.name,
      description: editForm.description
    });
    ElMessage.success('修改成功');
    showEditDialog.value = false;
    loadData();
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '修改失败');
  }
};

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该学校吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );
    const response = await axios.delete(`/api/admin/schools/${row.id}`);
    ElMessage.success(response.data.message);
    loadData();
  } catch (error) {
    if (error.response) {
      ElMessage.error(error.response.data.error);
    } else if (error !== 'cancel') {
      ElMessage.error('删除失败');
    }
  }
};

onMounted(() => {
  loadData();
});
</script>

<style scoped>
.school-list {
  padding: 20px;
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.filter-area {
  display: flex;
  gap: 12px;
}
.filter-input {
  width: 250px;
}
.action-area {
  display: flex;
  gap: 12px;
}
</style>