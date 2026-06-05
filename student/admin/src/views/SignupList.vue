<template>
  <div class="signup-list">
    <div class="toolbar">
      <div class="filter-area">
        <el-input
          v-model="filters.name"
          placeholder="按姓名搜索"
          class="filter-input"
          @keyup.enter="loadData"
        />
        <el-input
          v-model="filters.phone"
          placeholder="按手机号搜索"
          class="filter-input"
          @keyup.enter="loadData"
        />
        <el-select
          v-model="filters.status"
          placeholder="选择状态"
          class="filter-select"
        >
          <el-option label="全部" value="all" />
          <el-option label="报名中" value="pending" />
          <el-option label="报名成功" value="approved" />
          <el-option label="报名失败" value="rejected" />
        </el-select>
        <el-button type="primary" @click="loadData">筛选</el-button>
        <el-button @click="resetFilters">重置</el-button>
      </div>
      <div class="action-area">
        <el-button type="success" @click="showAddDialog = true">新增报名</el-button>
        <el-button type="warning" @click="exportExcel">导出Excel</el-button>
      </div>
    </div>

    <el-table :data="tableData" border style="width: 100%">
      <el-table-column prop="name" label="姓名" />
      <el-table-column prop="phone" label="手机号" />
      <el-table-column prop="age" label="年龄" />
      <el-table-column prop="hukou" label="户口地址" />
      <el-table-column prop="school" label="学校" />
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-select
            v-if="editingId === scope.row.id"
            v-model="scope.row.status"
            class="status-select"
            @change="handleStatusChange(scope.row)"
          >
            <el-option label="报名中" value="pending" />
            <el-option label="报名成功" value="approved" />
            <el-option label="报名失败" value="rejected" />
          </el-select>
          <span v-else :class="getStatusClass(scope.row.status)">
            {{ getStatusText(scope.row.status) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="提交时间" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button
            type="text"
            @click="openEditDialog(scope.row)"
          >编辑</el-button>
          <el-button
            type="text"
            @click="startEdit(scope.row.id)"
            v-if="editingId !== scope.row.id"
          >修改状态</el-button>
          <el-button
            type="text"
            @click="cancelEdit"
            v-else
          >取消</el-button>
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

    <el-dialog title="编辑报名信息" v-model="showEditDialog" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="editForm.phone" />
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input v-model.number="editForm.age" type="number" />
        </el-form-item>
        <el-form-item label="户口地址" prop="hukou">
          <el-input v-model="editForm.hukou" />
        </el-form-item>
        <el-form-item label="学校" prop="school">
          <el-input v-model="editForm.school" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit">确认修改</el-button>
      </template>
    </el-dialog>

    <el-dialog title="新增报名" v-model="showAddDialog" width="400px">
      <el-form :model="addForm" label-width="80px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="addForm.name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="addForm.phone" />
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input v-model.number="addForm.age" type="number" />
        </el-form-item>
        <el-form-item label="户口地址" prop="hukou">
          <el-input v-model="addForm.hukou" />
        </el-form-item>
        <el-form-item label="学校" prop="school">
          <el-input v-model="addForm.school" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="addForm.status">
            <el-option label="报名中" value="pending" />
            <el-option label="报名成功" value="approved" />
            <el-option label="报名失败" value="rejected" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddSubmit">确认添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>import { ref, reactive, onMounted } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';
const filters = reactive({
 name: '',
 phone: '',
 status: 'all'
});
const tableData = ref([]);
const editingId = ref(null);
const pagination = reactive({
 page: 1,
 pageSize: 10,
 total: 0
});
const showEditDialog = ref(false);
const showAddDialog = ref(false);
const editForm = reactive({
 id: '',
 name: '',
 phone: '',
 age: '',
 hukou: '',
 school: ''
});
const addForm = reactive({
 name: '',
 phone: '',
 age: '',
 hukou: '',
 school: '',
 status: 'pending'
});
const loadData = async () => {
 try {
 const response = await axios.get('/api/admin/signups', {
 params: {
 page: pagination.page,
 pageSize: pagination.pageSize,
 keyword: filters.name || filters.phone,
 status: filters.status
 }
 });
 tableData.value = response.data.list;
 pagination.total = response.data.total;
 }
 catch (error) {
 ElMessage.error('获取数据失败');
 }
};
const resetFilters = () => {
 filters.name = '';
 filters.phone = '';
 filters.status = 'all';
 pagination.page = 1;
 loadData();
};
const startEdit = (id) => {
 editingId.value = id;
};
const cancelEdit = () => {
 editingId.value = null;
 loadData();
};
const handleStatusChange = async (row) => {
 try {
 await axios.put(`/api/admin/signups/${row.id}/status`, {
 status: row.status
 });
 ElMessage.success('状态更新成功');
 editingId.value = null;
 }
 catch (error) {
 ElMessage.error('更新失败');
 loadData();
 }
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
const getStatusText = (status) => {
 const map = {
 pending: '报名中',
 approved: '报名成功',
 rejected: '报名失败'
 };
 return map[status] || status;
};
const getStatusClass = (status) => {
 const map = {
 pending: 'status-pending',
 approved: 'status-approved',
 rejected: 'status-rejected'
 };
 return map[status] || '';
};
const openEditDialog = (row) => {
 editForm.id = row.id;
 editForm.name = row.name;
 editForm.phone = row.phone;
 editForm.age = row.age;
 editForm.hukou = row.hukou;
 editForm.school = row.school;
 showEditDialog.value = true;
};
const handleEditSubmit = async () => {
 try {
 await axios.put(`/api/admin/signups/${editForm.id}`, {
 name: editForm.name,
 phone: editForm.phone,
 age: editForm.age,
 hukou: editForm.hukou,
 school: editForm.school
 });
 ElMessage.success('修改成功');
 showEditDialog.value = false;
 loadData();
 }
 catch (error) {
 ElMessage.error('修改失败');
 }
};
const handleAddSubmit = async () => {
 try {
 await axios.post('/api/admin/signups', {
 name: addForm.name,
 phone: addForm.phone,
 age: addForm.age,
 hukou: addForm.hukou,
 school: addForm.school,
 status: addForm.status
 });
 ElMessage.success('添加成功');
 showAddDialog.value = false;
 addForm.name = '';
 addForm.phone = '';
 addForm.age = '';
 addForm.hukou = '';
 addForm.school = '';
 addForm.status = 'pending';
 loadData();
 }
 catch (error) {
 ElMessage.error('添加失败');
 }
};
const exportExcel = async () => {
 try {
 const params = new URLSearchParams({
 keyword: filters.name || filters.phone,
 status: filters.status
 });
 const response = await axios.get(`/api/admin/signups/export?${params.toString()}`, {
 responseType: 'blob'
 });
 const url = window.URL.createObjectURL(new Blob([response.data]));
 const a = document.createElement('a');
 a.href = url;
 a.download = 'signups.xlsx';
 document.body.appendChild(a);
 a.click();
 window.URL.revokeObjectURL(url);
 document.body.removeChild(a);
 ElMessage.success('导出成功');
 }
 catch (error) {
 ElMessage.error('导出失败');
 }
};
onMounted(() => {
 loadData();
});
</script>

<style scoped>
.signup-list {
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
  width: 200px;
}
.filter-select {
  width: 150px;
}
.action-area {
  display: flex;
  gap: 12px;
}
.status-select {
  width: 120px;
}
.status-pending {
  color: #e6a23c;
}
.status-approved {
  color: #67c23a;
}
.status-rejected {
  color: #f56c6c;
}
</style>