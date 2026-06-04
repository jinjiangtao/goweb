<template>
  <div class="page-container">
    <div class="page-header">
      <el-input v-model="searchForm.username" placeholder="请输入用户名搜索" class="search-input" @keyup.enter="loadUsers" />
      <el-button type="primary" @click="handleAdd">新增</el-button>
    </div>
    <el-table :data="userList" border style="width: 100%">
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="角色">
        <template #default="scope">
          <el-tag :type="getRoleTagType(scope.row.role)">
            {{ getRoleName(scope.row.role) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-switch :value="scope.row.status === 1" @change="handleStatusChange(scope.row)" />
        </template>
      </el-table-column>
      <el-table-column prop="last_login_at" label="最后登录时间" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="text" @click="handleEdit(scope.row)" :disabled="scope.row.role === 'super'">编辑</el-button>
          <el-button type="text" @click="handleDelete(scope.row)" :disabled="scope.row.role === 'super'">删除</el-button>
          <el-button type="text" @click="handleResetPassword(scope.row)">重置密码</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      :current-page="pagination.page"
      :page-size="pagination.size"
      :total="pagination.total"
      @current-change="handlePageChange"
    />
    <el-dialog :visible="dialogVisible" :title="dialogTitle" width="400px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="isEdit && form.role === 'super'" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" :disabled="form.role === 'super'">
            <el-option label="超级管理员" value="super" />
            <el-option label="管理员" value="admin" />
            <el-option label="操作员" value="operator" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>import { ref, reactive, computed, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getUsers, createUser, updateUser, deleteUser, updateUserStatus, resetPassword } from '../api';
const userList = ref([]);
const pagination = reactive({
 page: 1,
 size: 10,
 total: 0
});
const searchForm = reactive({
 username: ''
});
const dialogVisible = ref(false);
const dialogTitle = ref('');
const isEdit = ref(false);
const formRef = ref(null);
const form = reactive({
 id: '',
 username: '',
 nickname: '',
 password: '',
 role: 'admin'
});
const rules = computed(() => ({
 username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
 nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
 password: isEdit.value ? [] : [{ required: true, message: '请输入密码', trigger: 'blur' }],
 role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}));
const roleMap = {
 super: '超级管理员',
 admin: '管理员',
 operator: '操作员'
};
const getRoleName = (role) => roleMap[role] || role;
const getRoleTagType = (role) => {
 switch (role) {
 case 'super': return 'danger';
 case 'admin': return 'warning';
 default: return 'info';
 }
};
const loadUsers = async () => {
 const res = await getUsers({
 page: pagination.page,
 size: pagination.size,
 username: searchForm.username
 });
 if (res.code === 200) {
 userList.value = res.data.list;
 pagination.total = res.data.total;
 }
};
const handlePageChange = (page) => {
 pagination.page = page;
 loadUsers();
};
const handleAdd = () => {
 dialogTitle.value = '新增管理员';
 isEdit.value = false;
 form.id = '';
 form.username = '';
 form.nickname = '';
 form.password = '';
 form.role = 'admin';
 dialogVisible.value = true;
};
const handleEdit = (row) => {
 dialogTitle.value = '编辑管理员';
 isEdit.value = true;
 form.id = row.id;
 form.username = row.username;
 form.nickname = row.nickname;
 form.role = row.role;
 dialogVisible.value = true;
};
const handleDelete = (row) => {
 if (row.role === 'super') {
 ElMessage.warning('超级管理员不可删除');
 return;
 }
 ElMessage.confirm('确定删除该管理员?', '提示', {
 confirmButtonText: '确定',
 cancelButtonText: '取消'
 }).then(async () => {
 const res = await deleteUser(row.id);
 if (res.code === 200) {
 ElMessage.success('删除成功');
 loadUsers();
 }
 });
};
const handleStatusChange = async (row) => {
 if (row.role === 'super') {
 row.status = 1;
 ElMessage.warning('超级管理员状态不可修改');
 return;
 }
 const res = await updateUserStatus(row.id);
 if (res.code === 200) {
 row.status = res.data.status;
 ElMessage.success('状态更新成功');
 }
};
const handleResetPassword = (row) => {
 ElMessage.confirm('确定重置密码为123456?', '提示').then(async () => {
 const res = await resetPassword(row.id);
 if (res.code === 200) {
 ElMessage.success('密码重置成功');
 }
 });
};
const handleSubmit = async () => {
 const valid = await formRef.value.validate();
 if (!valid)
 return;
 try {
 let res;
 if (isEdit.value) {
 res = await updateUser(form.id, {
 nickname: form.nickname,
 role: form.role,
 status: 1
 });
 }
 else {
 res = await createUser(form);
 }
 if (res.code === 200) {
 ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
 dialogVisible.value = false;
 loadUsers();
 }
 }
 catch (error) {
 ElMessage.error('操作失败');
 }
};
onMounted(() => {
 loadUsers();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.page-header {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}
.search-input {
  width: 300px;
}
</style>