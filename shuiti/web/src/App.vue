<template>
  <el-container class="layout-container">
    <el-header class="layout-header">
      <div class="header-content">
        <div class="logo" @click="goHome">
          <el-icon :size="28" color="#409EFF"><Reading /></el-icon>
          <span class="title">刷题学习系统</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          router
          class="nav-menu"
          background-color="transparent"
          text-color="#606266"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/practice">
            <el-icon><EditPen /></el-icon>
            <span>开始刷题</span>
          </el-menu-item>
          <el-menu-item index="/wrong">
            <el-icon><Warning /></el-icon>
            <span>错题本</span>
          </el-menu-item>
          <el-menu-item index="/stats">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据统计</span>
          </el-menu-item>
          <el-menu-item index="/records">
            <el-icon><Document /></el-icon>
            <span>答题记录</span>
          </el-menu-item>
          <el-menu-item index="/admin">
            <el-icon><Setting /></el-icon>
            <span>题库管理</span>
          </el-menu-item>
        </el-menu>
        <div class="user-area">
          <el-input
            v-model="userId"
            placeholder="用户ID"
            size="small"
            style="width: 140px; margin-right: 12px;"
            @change="saveUserId"
          />
          <el-tag type="success" effect="plain">在线</el-tag>
        </div>
      </div>
    </el-header>
    <el-main class="layout-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>
    <el-footer class="layout-footer">
      <span>刷题学习系统 © 2026 | Go Gin + Vue3 + SQLite</span>
    </el-footer>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const userId = ref(userStore.userId || 'user_' + Math.random().toString(36).slice(2, 8))
const activeMenu = computed(() => route.path)

onMounted(() => {
  if (!userStore.userId) {
    userStore.setUserId(userId.value)
  } else {
    userId.value = userStore.userId
  }
})

const saveUserId = () => {
  if (userId.value) {
    userStore.setUserId(userId.value)
  }
}

const goHome = () => {
  router.push('/')
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
}

.layout-container {
  height: 100vh;
}

.layout-header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
  height: 64px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  z-index: 100;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 24px;
}

.logo {
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-right: 48px;
}

.logo .title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin-left: 10px;
}

.nav-menu {
  flex: 1;
  border-bottom: none !important;
}

.nav-menu .el-menu-item {
  height: 64px;
  line-height: 64px;
}

.user-area {
  display: flex;
  align-items: center;
}

.layout-main {
  background: #f5f7fa;
  padding: 24px;
  overflow-y: auto;
}

.layout-footer {
  background: #fff;
  border-top: 1px solid #e4e7ed;
  text-align: center;
  color: #909399;
  font-size: 13px;
  line-height: 50px;
  height: 50px !important;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
