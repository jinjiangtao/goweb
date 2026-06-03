<template>
  <div class="page-container">
    <div class="role-tabs">
      <el-tabs v-model="activeRole" @change="handleRoleChange">
        <el-tab-pane label="超级管理员" name="super" :disabled="true" />
        <el-tab-pane label="管理员" name="admin" />
        <el-tab-pane label="操作员" name="operator" />
      </el-tabs>
    </div>
    <div class="permission-content">
      <div v-if="activeRole === 'super'" class="super-tip">
        <el-alert title="超级管理员默认拥有所有权限" type="info" :closable="false" />
      </div>
      <div v-else>
        <el-tree
          ref="treeRef"
          :data="menuTree"
          :props="treeProps"
          show-checkbox
          default-expand-all
          :check-strictly="true"
          :checked-keys="checkedKeys"
          @check-change="handleCheckChange"
        />
        <div class="button-group">
          <el-button type="primary" @click="handleSave">保存权限</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>import { ref, watch, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getRoleMenus, setRoleMenus } from '../api';
const activeRole = ref('super');
const menuTree = ref([]);
const checkedKeys = ref([]);
const treeRef = ref(null);
const treeProps = {
 children: 'children',
 label: 'name'
};
const loadRoleMenus = async () => {
 const res = await getRoleMenus(activeRole.value);
 if (res.code === 200) {
 menuTree.value = res.data.tree;
 checkedKeys.value = res.data.ids.map(id => id.toString());
 }
};
const handleRoleChange = () => {
 loadRoleMenus();
};
const handleCheckChange = () => {
};
const handleSave = async () => {
 const checkedNodes = treeRef.value.getCheckedNodes();
 const menuIds = checkedNodes.map(node => node.id);
 const res = await setRoleMenus(activeRole.value, { menu_ids: menuIds });
 if (res.code === 200) {
 ElMessage.success('权限设置成功');
 }
};
const handleReset = () => {
 loadRoleMenus();
};
watch(activeRole, () => {
 loadRoleMenus();
});
onMounted(() => {
 loadRoleMenus();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.role-tabs {
  margin-bottom: 20px;
}
.permission-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
}
.super-tip {
  text-align: center;
}
.button-group {
  margin-top: 20px;
  text-align: right;
}
</style>