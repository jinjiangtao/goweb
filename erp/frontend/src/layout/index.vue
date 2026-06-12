<template>
  <div class="layout">
    <div class="header">
      <h2>ERP管理系统</h2>
      <button @click="handleLogout">退出登录</button>
    </div>
    <div class="main-content">
      <div class="sidebar">
        <ul>
          <template v-for="menu in menuTree" :key="menu.id || menu.ID">
            <li v-if="!menu.hidden">
              <router-link :to="menu.path">{{ menu.name }}</router-link>
              <ul v-if="menu.children && menu.children.length > 0" class="submenu">
                <li v-for="child in menu.children" :key="child.id || child.ID" v-if="!child.hidden">
                  <router-link :to="child.path">{{ child.name }}</router-link>
                </li>
              </ul>
            </li>
          </template>
        </ul>
      </div>
      <div class="content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user.js'

const router = useRouter()
const userStore = useUserStore()

// 构建菜单树
const buildMenuTree = (menus) => {
  const menuMap = {}
  const result = []
  
  menus.forEach(menu => {
    const menuId = menu.id || menu.ID
    menuMap[menuId] = { ...menu, children: [] }
  })
  
  menus.forEach(menu => {
    const menuId = menu.id || menu.ID
    const parentId = menu.parentId || menu.ParentId || 0
    if (!parentId || parentId === 0) {
      result.push(menuMap[menuId])
    } else if (menuMap[parentId]) {
      menuMap[parentId].children.push(menuMap[menuId])
    }
  })
  
  return result
}

const menuTree = computed(() => buildMenuTree(userStore.menus || []))

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
  margin: 0;
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

.sidebar .submenu {
  margin-left: 10px;
  font-size: 14px;
}

.content {
  flex: 1;
  padding: 20px;
  background: #f0f2f5;
}
</style>
