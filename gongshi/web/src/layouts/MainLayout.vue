<template>
  <el-container class="main-layout">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <el-icon :size="24"><Timer /></el-icon>
        <span>工时管理</span>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        background-color="transparent"
        text-color="#c0c4cc"
        active-text-color="#ffffff"
        class="menu"
      >
        <el-menu-item index="/timesheet">
          <el-icon><Calendar /></el-icon>
          <span>工时填报</span>
        </el-menu-item>
        <el-menu-item index="/history">
          <el-icon><Document /></el-icon>
          <span>历史记录</span>
        </el-menu-item>
        <el-menu-item index="/statistics">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据统计</span>
        </el-menu-item>
        <el-menu-item v-if="userStore.isAdmin" index="/review">
          <el-icon><Check /></el-icon>
          <span>工时审核</span>
        </el-menu-item>
        <el-menu-item v-if="userStore.isAdmin" index="/projects">
          <el-icon><Folder /></el-icon>
          <span>项目管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <h2>{{ $route.meta.title || '工时管理系统' }}</h2>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" class="avatar">
                {{ userStore.user?.name?.charAt(0) }}
              </el-avatar>
              <span class="username">{{ userStore.user?.name }}</span>
              <el-tag v-if="userStore.isAdmin" type="danger" size="small" style="margin-left: 8px">管理员</el-tag>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout" icon="SwitchButton">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const handleCommand = (cmd) => {
  if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
}
.sidebar {
  background: linear-gradient(180deg, #304156 0%, #1f2d3d 100%);
  overflow: hidden;
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  gap: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.menu {
  border: none;
  padding-top: 10px;
}
.menu :deep(.el-menu-item) {
  height: 48px;
  line-height: 48px;
}
.menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(90deg, #409eff 0%, #66b1ff 100%);
  border-radius: 4px;
  margin: 4px 10px;
}
.header {
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
}
.header-left h2 {
  font-size: 18px;
  color: #303133;
  font-weight: 600;
}
.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 10px;
  color: #606266;
}
.avatar {
  background: #409eff;
  color: #fff;
}
.username {
  font-weight: 500;
}
.main-content {
  padding: 0;
  background: #f5f7fa;
}
</style>
