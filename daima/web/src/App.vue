<template>
  <div id="app">
    <el-container class="layout-container">
      <el-header class="app-header">
        <div class="header-content">
          <router-link to="/" class="logo">
            <el-icon><Document /></el-icon>
            <span>代码片段分享平台</span>
          </router-link>
            <div class="header-search">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索代码片段..."
              size="large"
              :prefix-icon="Search"
              class="search-input"
              @keyup.enter="goSearch"
            >
              <template #append>
                <el-button :icon="Search" type="primary" @click="goSearch">搜索</el-button>
              </template>
            </el-input>
          </div>
          <div class="header-nav">
            <template v-if="userStore.isLoggedIn">
              <el-button type="primary" :icon="Plus" size="large" @click="$router.push('/create')">
                创建片段
              </el-button>
              <el-dropdown trigger="click">
                <div class="user-info">
                  <el-avatar :size="36" :src="userStore.user?.avatar">
                    {{ userStore.user?.username?.charAt(0) }}
                  </el-avatar>
                  <span class="username">{{ userStore.user?.username }}</span>
                  <el-icon><ArrowDown /></el-icon>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="$router.push('/my')">
                      <el-icon><Folder /></el-icon>我的片段
                    </el-dropdown-item>
                    <el-dropdown-item @click="$router.push('/favorites')">
                      <el-icon><Star /></el-icon>我的收藏
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="logout">
                      <el-icon><SwitchButton /></el-icon>退出登录
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
            <template v-else>
              <el-button size="large" @click="showLogin = true">登录</el-button>
              <el-button type="primary" size="large" @click="showRegister = true">注册</el-button>
            </template>
          </div>
        </div>
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
    <LoginDialog v-model:visible="showLogin" @success="handleLoginSuccess" />
    <RegisterDialog v-model:visible="showRegister" @success="handleRegisterSuccess" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from './store/user'
import LoginDialog from './components/LoginDialog.vue'
import RegisterDialog from './components/RegisterDialog.vue'

const router = useRouter()
const userStore = useUserStore()
const searchKeyword = ref('')
const showLogin = ref(false)
const showRegister = ref(false)

onMounted(() => {
  userStore.checkAuth()
})

const goSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({ name: 'search', query: { q: searchKeyword.value } })
  }
}

const handleLoginSuccess = () => {
  showLogin.value = false
}

const handleRegisterSuccess = () => {
  showRegister.value = false
}

const logout = () => {
  userStore.logout()
  router.push('/')
  ElMessage.success('已退出登录')
}
</script>

<style scoped>
.layout-container {
  min-height: 100%;
}
.app-header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
  height: 64px;
  display: flex;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}
.header-content {
  width: 100%;
  max-width: 1600px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  gap: 32px;
}
.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  text-decoration: none;
  white-space: nowrap;
}
.logo .el-icon {
  color: #409eff;
  font-size: 24px;
}
.header-search {
  flex: 1;
  max-width: 500px;
}
.search-input :deep(.el-input__wrapper) {
  border-radius: 8px;
}
.header-nav {
  display: flex;
  align-items: center;
  gap: 16px;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.2s;
}
.user-info:hover {
  background: #f5f7fa;
}
.username {
  font-weight: 500;
  color: #1f2937;
}
.app-main {
  background: #f5f7fa;
  padding: 0;
}
</style>
