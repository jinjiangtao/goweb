<template>
  <div class="home-container">
    <aside class="sidebar">
      <div class="logo">
        <h2>管理后台</h2>
      </div>
      <el-menu :default-active="activeMenu" mode="vertical">
        <el-menu-item index="1" @click="navigate('/')">
          <el-icon><component :is="List" /></el-icon>
          <span>报名列表</span>
        </el-menu-item>
        <el-menu-item index="2" @click="navigate('/stats')">
          <el-icon><component :is="DataAnalysis" /></el-icon>
          <span>统计看板</span>
        </el-menu-item>
      </el-menu>
    </aside>
    <main class="main-content">
      <header class="header">
        <span class="title">{{ pageTitle }}</span>
        <el-button type="text" @click="handleLogout">退出登录</el-button>
      </header>
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { List, DataAnalysis } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeMenu = computed(() => {
  return route.path === '/' ? '1' : '2'
})

const pageTitle = computed(() => {
  return route.path === '/' ? '报名列表' : '统计看板'
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
  background: #2c3e50;
  color: white;
  padding: 20px 0;
}
.logo {
  text-align: center;
  padding: 0 20px 20px;
  border-bottom: 1px solid #34495e;
}
.logo h2 {
  margin: 0;
  font-size: 18px;
}
.el-menu {
  border-right: none;
}
.el-menu-item {
  color: #ecf0f1;
}
.el-menu-item:hover, .el-menu-item.is-active {
  background: #34495e;
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
.title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}
</style>
