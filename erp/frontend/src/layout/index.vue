<template>
  <div class="layout">
    <div class="header">
      <div class="header-left">
        <h2>ERP管理系统</h2>
      </div>
      <div class="header-right">
        <el-button type="danger" @click="handleLogout">退出登录</el-button>
      </div>
    </div>
    <div class="main-content">
      <div class="sidebar">
        <el-menu
          :default-active="activeMenu"
          class="sidebar-menu"
          router
          unique-opened
          :collapse="false"
        >
          <template v-for="menu in menuList" :key="menu.path">
            <el-menu-item v-if="!menu.children || menu.children.length === 0" :index="menu.path">
              <el-icon><component :is="getMenuIcon(menu.icon)" /></el-icon>
              <span>{{ menu.name }}</span>
            </el-menu-item>
            <el-sub-menu v-else :index="menu.path">
              <template #title>
                <el-icon><component :is="getMenuIcon(menu.icon)" /></el-icon>
                <span>{{ menu.name }}</span>
              </template>
              <el-menu-item
                v-for="child in menu.children"
                :key="child.path"
                :index="child.path"
              >
                <span>{{ child.name }}</span>
              </el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
      </div>
      <div class="content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../store/user.js'
import * as ElementPlusIcons from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const iconMap = {
  House: ElementPlusIcons.House,
  Setting: ElementPlusIcons.Setting,
  User: ElementPlusIcons.User,
  UserFilled: ElementPlusIcons.UserFilled,
  Menu: ElementPlusIcons.Menu,
  Goods: ElementPlusIcons.Goods,
  OfficeBuilding: ElementPlusIcons.OfficeBuilding,
  ShoppingCart: ElementPlusIcons.ShoppingCart,
  Sell: ElementPlusIcons.Sell
}

const getMenuIcon = (iconName) => {
  return iconMap[iconName] || ElementPlusIcons.Menu
}

const menuList = computed(() => {
  return [
    {
      name: '首页',
      path: '/dashboard',
      icon: 'House'
    },
    {
      name: '系统设置',
      path: '/system',
      icon: 'Setting',
      children: [
        {
          name: '用户管理',
          path: '/system/user'
        },
        {
          name: '角色管理',
          path: '/system/role'
        },
        {
          name: '菜单管理',
          path: '/system/menu'
        }
      ]
    },
    {
      name: '产品管理',
      path: '/product',
      icon: 'Goods'
    },
    {
      name: '客户管理',
      path: '/customer',
      icon: 'User'
    },
    {
      name: '供应商管理',
      path: '/supplier',
      icon: 'OfficeBuilding'
    },
    {
      name: '采购订单',
      path: '/purchase-order',
      icon: 'ShoppingCart'
    },
    {
      name: '销售订单',
      path: '/sales-order',
      icon: 'Sell'
    }
  ]
})

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
  width: 100%;
  overflow: hidden;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  padding: 0 24px;
  background: #001529;
  color: white;
}

.header-left h2 {
  margin: 0;
  font-size: 20px;
  color: white;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 220px;
  background: #001529;
  flex-shrink: 0;
  overflow-y: auto;
}

.sidebar-menu {
  height: 100%;
  border: none;
  background: #001529;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 220px;
}

.sidebar-menu .el-menu-item,
.sidebar-menu .el-sub-menu__title {
  color: #bfcbd9;
}

.sidebar-menu .el-menu-item:hover,
.sidebar-menu .el-sub-menu__title:hover {
  color: #409EFF !important;
  background: #263445 !important;
}

.sidebar-menu .el-menu-item.is-active {
  color: #409EFF !important;
  background: #263445 !important;
}

.content {
  flex: 1;
  padding: 24px;
  background: #f0f2f5;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
