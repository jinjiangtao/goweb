<template>
  <div class="home-container">
    <aside class="sidebar">
      <div class="logo">
        <h2>管理后台</h2>
      </div>
      <el-menu :default-active="activeMenu" mode="vertical">
        <template v-for="menu in authStore.menus" :key="menu.id">
          <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.id.toString()">
            <template #title>
              <el-icon><component :is="getIcon(menu.icon)" /></el-icon>
              <span>{{ menu.name }}</span>
            </template>
            <el-menu-item
              v-for="child in menu.children"
              :key="child.id"
              :index="child.id.toString()"
              @click="navigate(child.path)"
            >
              <el-icon><component :is="getIcon(child.icon)" /></el-icon>
              <span>{{ child.name }}</span>
            </el-menu-item>
          </el-sub-menu>
          <el-menu-item v-else :index="menu.id.toString()" @click="navigate(menu.path)">
            <el-icon><component :is="getIcon(menu.icon)" /></el-icon>
            <span>{{ menu.name }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </aside>
    <main class="main-content">
      <header class="header">
        <span class="title">{{ pageTitle }}</span>
        <div class="header-right">
          <span class="user-info">{{ authStore.user?.nickname }} ({{ authStore.user?.role === 'super_admin' ? '超级管理员' : '普通管理员' }})</span>
          <el-button type="text" @click="handleLogout">退出登录</el-button>
        </div>
      </header>
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, markRaw } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { List, DataAnalysis, School, User, Menu } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const iconMap = {
  List: markRaw(List),
  DataAnalysis: markRaw(DataAnalysis),
  School: markRaw(School),
  User: markRaw(User),
  Menu: markRaw(Menu)
}

const getIcon = (iconName) => {
  return iconMap[iconName] || List
}

const activeMenu = computed(() => {
  const currentPath = route.path
  for (const menu of authStore.menus) {
    if (menu.path === currentPath) {
      return menu.id.toString()
    }
    if (menu.children) {
      for (const child of menu.children) {
        if (child.path === currentPath) {
          return child.id.toString()
        }
      }
    }
  }
  return ''
})

const pageTitle = computed(() => {
  const currentPath = route.path
  for (const menu of authStore.menus) {
    if (menu.path === currentPath) {
      return menu.name
    }
    if (menu.children) {
      for (const child of menu.children) {
        if (child.path === currentPath) {
          return child.name
        }
      }
    }
  }
  return '管理后台'
})

const navigate = (path) => {
  router.push(path)
}

const handleLogout = () => {
  authStore.logout()
  ElMessage.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.home-container {
  display: flex;
  min-height: 100vh;
  background: #f5f5f5;
}
.sidebar {
  width: 200px;
  background: #34495e;
  color: white;
  padding: 20px 0;
}
.logo {
  text-align: center;
  padding: 0 20px 20px;
  border-bottom: 1px solid #2c3e50;
}
.logo h2 {
  margin: 0;
  font-size: 18px;
}
.el-menu {
  border-right: none;
  background: transparent;
}
.el-menu-item,
.el-sub-menu__title {
  color: #bdc3c7;
  background: transparent;
}
.el-menu-item:hover,
.el-sub-menu__title:hover {
  color: #ecf0f1;
  background: rgba(255, 255, 255, 0.1);
}
.el-menu-item.is-active {
  color: #ffffff;
  background: #2c3e50;
  border-left: 3px solid #3498db;
}
.el-sub-menu__title.is-active {
  color: #ffffff;
  background: #2c3e50;
}
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}
.user-info {
  font-size: 14px;
  color: #666;
}
.title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}
</style>