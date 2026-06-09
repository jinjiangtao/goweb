<template>
  <el-container>
    <el-aside width="200px">
      <div class="logo">内推系统</div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <template v-if="authStore.isEmployee()">
          <el-menu-item index="/jobs">
            <el-icon><Document /></el-icon>
            <span>职位管理</span>
          </el-menu-item>
          <el-menu-item index="/referrals">
            <el-icon><User /></el-icon>
            <span>我的推荐</span>
          </el-menu-item>
        </template>
        <template v-if="authStore.isHR()">
          <el-menu-item index="/admin/referrals">
            <el-icon><Document /></el-icon>
            <span>内推列表</span>
          </el-menu-item>
          <el-menu-item index="/admin/jobs">
            <el-icon><Briefcase /></el-icon>
            <span>职位列表</span>
          </el-menu-item>
          <el-menu-item index="/admin/stats">
            <el-icon><DataAnalysis /></el-icon>
            <span>统计看板</span>
          </el-menu-item>
        </template>
        <template v-if="authStore.isAdmin()">
          <el-menu-item index="/admin/users">
            <el-icon><Setting /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-content">
          <span>欢迎，{{ authStore.user?.real_name }}</span>
          <el-button type="text" @click="handleLogout">退出登录</el-button>
        </div>
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
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, User, Briefcase, DataAnalysis, Setting } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeMenu = computed(() => route.path)

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗?', '提示')
    authStore.logout()
    router.push('/login')
    ElMessage.success('已退出')
  } catch {}
}
</script>

<style scoped>
.el-container {
  height: 100vh;
}
.el-aside {
  background-color: #304156;
  color: #fff;
}
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
}
.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
}
.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.el-main {
  background-color: #f0f2f5;
}
</style>
