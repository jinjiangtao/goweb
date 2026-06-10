<template>
  <el-container style="height: 100vh">
    <el-aside width="200px" style="background: #304156">
      <div style="color: white; padding: 20px; font-size: 18px; font-weight: bold">抽奖管理</div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/stats">
          <span>统计看板</span>
        </el-menu-item>
        <el-menu-item index="/prizes">
          <span>奖品管理</span>
        </el-menu-item>
        <el-menu-item index="/records">
          <span>抽奖记录</span>
        </el-menu-item>
        <el-menu-item index="/address">
          <span>地址库管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="background: white; border-bottom: 1px solid #eee; display: flex; justify-content: space-between; align-items: center">
        <span style="font-size: 20px; font-weight: bold">抽奖转盘后台管理系统</span>
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

const router = useRouter()
const route = useRoute()

const activeMenu = computed(() => route.path)

const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>
