<template>
  <div class="signup-list">
    <div class="search-bar">
      <el-input
        v-model="keyword"
        placeholder="按姓名或手机号搜索"
        class="search-input"
        @keyup.enter="loadData"
      />
      <el-button type="primary" @click="loadData">搜索</el-button>
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
  </div>
</template>

<script setup>import { ref, reactive, onMounted } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';
const keyword = ref('');
const tableData = ref([]);
const editingId = ref(null);
const pagination = reactive({
 page: 1,
 pageSize: 10,
 total: 0
});
const loadData = async () => {
 try {
 const response = await axios.get('/api/admin/signups', {
 params: {
 page: pagination.page,
 pageSize: pagination.pageSize,
 keyword: keyword.value
 }
 });
 tableData.value = response.data.list;
 pagination.total = response.data.total;
 }
 catch (error) {
 ElMessage.error('获取数据失败');
 }
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
onMounted(() => {
 loadData();
});
</script>

<style scoped>
.signup-list {
  padding: 20px;
}
.search-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}
.search-input {
  width: 300px;
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
