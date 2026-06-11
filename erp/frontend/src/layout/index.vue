<template>
  <div class="layout">
    <div class="header">
      <h2>ERP管理系统</h2>
      <button @click="handleLogout">退出登录</button>
    </div>
    <div class="main-content">
      <div class="sidebar">
        <ul>
          <li><router-link to="/dashboard">仪表板</router-link></li>
          <li><router-link to="/system/user">用户管理</router-link></li>
          <li><router-link to="/system/role">角色管理</router-link></li>
          <li><router-link to="/system/menu">菜单管理</router-link></li>
          <li><router-link to="/product">产品管理</router-link></li>
        </ul>
      </div>
      <div class="content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user.js'

const router = useRouter()
const userStore = useUserStore()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.header {
  background: #304156;
  color: white;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.main-content {
  display: flex;
  flex: 1;
}

.sidebar {
  width: 200px;
  background: #304156;
  padding: 20px 0;
}

.sidebar ul {
  list-style: none;
  padding: 0;
}

.sidebar li {
  margin: 10px 0;
}

.sidebar a {
  display: block;
  padding: 10px 20px;
  color: #bfcbd9;
  text-decoration: none;
}

.sidebar a:hover,
.sidebar a.router-link-active {
  background: #263445;
  color: #409EFF;
}

.content {
  flex: 1;
  padding: 20px;
  background: #f0f2f5;
}
</style>
