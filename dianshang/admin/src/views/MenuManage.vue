<template>
  <div class="page-container">
    <div class="page-header">
      <el-button type="primary" @click="handleAdd">新增菜单</el-button>
    </div>
    <el-tree
      :data="menuTree"
      :props="treeProps"
      :expand-on-click-node="false"
      default-expand-all
    >
      <template #default="{ node, data }">
        <span>{{ data.name }}</span>
        <span class="tree-actions">
          <el-button type="text" size="small" @click="handleAddChild(data)">新增子菜单</el-button>
          <el-button type="text" size="small" @click="handleEdit(data)">编辑</el-button>
          <el-button type="text" size="small" @click="handleDelete(data)">删除</el-button>
        </span>
      </template>
    </el-tree>
    <el-dialog :visible="dialogVisible" :title="dialogTitle" width="450px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="父级菜单" prop="parent_id">
          <el-select v-model="form.parent_id">
            <el-option label="无" :value="0" />
            <el-option v-for="menu in parentMenuOptions" :key="menu.id" :label="menu.name" :value="menu.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="路由路径" prop="path">
          <el-input v-model="form.path" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-select v-model="form.icon">
            <el-option label="用户" value="user" />
            <el-option label="菜单" value="menu" />
            <el-option label="锁" value="lock" />
            <el-option label="设置" value="settings" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input v-model.number="form.sort" type="number" />
        </el-form-item>
        <el-form-item label="是否显示" prop="visible">
          <el-switch v-model="form.visible" :active-value="1" :inactive-value="0" />
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
import { getMenus, createMenu, updateMenu, deleteMenu } from '../api';
const menuTree = ref([]);
const dialogVisible = ref(false);
const dialogTitle = ref('');
const isEdit = ref(false);
const formRef = ref(null);
const form = reactive({
 id: '',
 parent_id: 0,
 name: '',
 path: '',
 icon: 'menu',
 sort: 0,
 visible: 1
});
const rules = {
 name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }]
};
const treeProps = {
 children: 'children',
 label: 'name'
};
const parentMenuOptions = computed(() => {
 const options = [];
 const flatten = (menus) => {
 for (const menu of menus) {
 if (!menu.children || menu.children.length === 0) {
 options.push({ id: menu.id, name: menu.name });
 }
 flatten(menu.children || []);
 }
 };
 flatten(menuTree.value);
 return options;
});
const loadMenus = async () => {
 const res = await getMenus();
 if (res.code === 200) {
 menuTree.value = res.data;
 }
};
const handleAdd = () => {
 dialogTitle.value = '新增菜单';
 isEdit.value = false;
 form.id = '';
 form.parent_id = 0;
 form.name = '';
 form.path = '';
 form.icon = 'menu';
 form.sort = 0;
 form.visible = 1;
 dialogVisible.value = true;
};
const handleAddChild = (parent) => {
 dialogTitle.value = '新增子菜单';
 isEdit.value = false;
 form.id = '';
 form.parent_id = parent.id;
 form.name = '';
 form.path = '';
 form.icon = 'menu';
 form.sort = 0;
 form.visible = 1;
 dialogVisible.value = true;
};
const handleEdit = (data) => {
 dialogTitle.value = '编辑菜单';
 isEdit.value = true;
 form.id = data.id;
 form.parent_id = data.parent_id || 0;
 form.name = data.name;
 form.path = data.path || '';
 form.icon = data.icon || 'menu';
 form.sort = data.sort || 0;
 form.visible = data.visible || 1;
 dialogVisible.value = true;
};
const handleDelete = (data) => {
 ElMessage.confirm('确定删除该菜单及所有子菜单?', '提示').then(async () => {
 const res = await deleteMenu(data.id);
 if (res.code === 200) {
 ElMessage.success('删除成功');
 loadMenus();
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
 res = await updateMenu(form.id, {
 name: form.name,
 path: form.path,
 icon: form.icon,
 sort: form.sort,
 visible: form.visible
 });
 }
 else {
 res = await createMenu({
 parent_id: form.parent_id,
 name: form.name,
 path: form.path,
 icon: form.icon,
 sort: form.sort,
 visible: form.visible
 });
 }
 if (res.code === 200) {
 ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
 dialogVisible.value = false;
 loadMenus();
 }
 }
 catch (error) {
 ElMessage.error('操作失败');
 }
};
onMounted(() => {
 loadMenus();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}
.page-header {
  margin-bottom: 20px;
}
.tree-actions {
  float: right;
  margin-left: 20px;
}
</style>