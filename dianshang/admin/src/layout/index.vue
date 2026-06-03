<template>
  <div class="layout-container">
    <aside class="sidebar" :class="{ collapsed: isCollapsed }">
      <div class="logo">
        <span v-if="!isCollapsed">后台管理</span>
      </div>
      <el-menu :default-active="activeMenu" mode="vertical" :collapse="isCollapsed">
        <template v-for="menu in menus" :key="menu.id">
          <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="String(menu.id)">
            <template #title>
              <el-icon><component :is="getIcon(menu.icon)" /></el-icon>
              <span>{{ menu.name }}</span>
            </template>
            <template v-for="child in menu.children" :key="child.id">
              <el-sub-menu v-if="child.children && child.children.length > 0" :index="String(child.id)">
                <template #title>
                  <el-icon><component :is="getIcon(child.icon)" /></el-icon>
                  <span>{{ child.name }}</span>
                </template>
                <el-menu-item
                  v-for="grandchild in child.children"
                  :key="grandchild.id"
                  :index="grandchild.path"
                  @click="handleMenuClick(grandchild.path)"
                >
                  <span>{{ grandchild.name }}</span>
                </el-menu-item>
              </el-sub-menu>
              <el-menu-item
                v-else
                :index="child.path"
                @click="handleMenuClick(child.path)"
              >
                <el-icon><component :is="getIcon(child.icon)" /></el-icon>
                <span>{{ child.name }}</span>
              </el-menu-item>
            </template>
          </el-sub-menu>
          <el-menu-item v-else :index="menu.path" @click="handleMenuClick(menu.path)">
            <el-icon><component :is="getIcon(menu.icon)" /></el-icon>
            <span>{{ menu.name }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </aside>
    <main class="main-content">
      <header class="header">
        <button class="collapse-btn" @click="toggleSidebar">
          <el-icon><Menu v-if="isCollapsed" /><Fold v-else /></el-icon>
        </button>
        <el-breadcrumb separator="/">
          <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path" :to="item.path">
            {{ item.name }}
          </el-breadcrumb-item>
        </el-breadcrumb>
        <div class="header-right">
          <span>{{ user?.nickname }}</span>
          <el-button type="text" @click="handleLogout">退出</el-button>
        </div>
      </header>
      <router-view />
    </main>
  </div>
</template>

<script setup>import { ref, computed, onMounted } from 'vue';
import { Menu, Fold, User, Settings, Lock } from '@element-plus/icons-vue';
import { useUserStore } from '../stores/user';
import { logout } from '../api';
import { ElMessage } from 'element-plus';
import router from '../router';
const userStore = useUserStore();
const isCollapsed = ref(false);
const activeMenu = ref('');
const user = computed(() => userStore.user);
const menus = computed(() => userStore.menus);
const breadcrumbs = ref([]);
const iconMap = {
 user: User,
 menu: Menu,
 lock: Lock
};
const getIcon = (iconName) => {
 return iconMap[iconName] || Settings;
};
const toggleSidebar = () => {
 isCollapsed.value = !isCollapsed.value;
};
const handleMenuClick = (path) => {
 activeMenu.value = path;
 router.push(path);
 updateBreadcrumbs(path);
};
const updateBreadcrumbs = (path) => {
 const crumbs = [];
 const findMenu = (items, currentPath) => {
 for (const menu of items) {
 if (menu.path === currentPath) {
 crumbs.unshift({ name: menu.name, path: menu.path });
 return true;
 }
 if (menu.children && menu.children.length > 0) {
 if (findMenu(menu.children, currentPath)) {
 crumbs.unshift({ name: menu.name, path: menu.path });
 return true;
 }
 }
 }
 return false;
 };
 findMenu(menus.value, path);
 breadcrumbs.value = crumbs;
};
const handleLogout = async () => {
 try {
 await logout();
 userStore.logout();
 ElMessage.success('退出成功');
 window.location.href = '/login';
 }
 catch (error) {
 ElMessage.error('退出失败');
 }
};
onMounted(() => {
 const currentPath = router.currentRoute.value.path;
 activeMenu.value = currentPath;
 updateBreadcrumbs(currentPath);
});
</script>

<style scoped>
.layout-container {
  display: flex;
  min-height: 100vh;
  background: #f5f5f5;
}
.sidebar {
  width: 200px;
  background: #2d3748;
  color: white;
  transition: width 0.3s;
}
.sidebar.collapsed {
  width: 60px;
}
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid #4a5568;
}
.el-menu {
  border-right: none;
}
.el-menu-item, .el-sub-menu__title {
  color: #e2e8f0;
}
.el-menu-item:hover, .el-sub-menu__title:hover {
  background: #4a5568;
}
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.header {
  height: 60px;
  background: white;
  display: flex;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.collapse-btn {
  border: none;
  background: none;
  font-size: 20px;
  cursor: pointer;
  margin-right: 20px;
}
.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 20px;
}
</style>