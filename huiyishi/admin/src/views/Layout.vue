<template>
  <el-container class="layout">
    <el-aside width="220px">
      <div class="logo">会议室预订</div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="transparent"
        text-color="#cbd5e1"
        active-text-color="#3b82f6"
      >
        <el-menu-item index="/stats">
          <el-icon><DataLine /></el-icon>
          <span>统计看板</span>
        </el-menu-item>
        <el-menu-item index="/rooms">
          <el-icon><OfficeBuilding /></el-icon>
          <span>会议室管理</span>
        </el-menu-item>
        <el-menu-item index="/bookings">
          <el-icon><Calendar /></el-icon>
          <span>预订管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <span class="welcome">欢迎, {{ authStore.admin?.nickname }}</span>
        <el-button type="danger" @click="handleLogout">退出登录</el-button>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { DataLine, OfficeBuilding, Calendar } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeMenu = computed(() => route.path)

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  height: 100vh;
}

.el-aside {
  background: rgba(30, 41, 59, 0.95);
  backdrop-filter: blur(10px);
  border-right: 1px solid #475569;
}

.logo {
  padding: 20px;
  font-size: 18px;
  font-weight: bold;
  color: #e2e8f0;
  border-bottom: 1px solid #475569;
  text-align: center;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.el-menu {
  border-right: none;
}

.el-header {
  background: #1e293b;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #475569;
}

.welcome {
  color: #cbd5e1;
}

.el-main {
  background: #0f172a;
  padding: 24px;
}
</style>
