<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409EFF">
              <el-icon :size="30"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67C23A">
              <el-icon :size="30"><Box /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.products }}</div>
              <div class="stat-label">产品总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #E6A23C">
              <el-icon :size="30"><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.roles }}</div>
              <div class="stat-label">角色总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #F56C6C">
              <el-icon :size="30"><Menu /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.menus }}</div>
              <div class="stat-label">菜单总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>欢迎使用ERP管理系统</span>
            </div>
          </template>
          <div class="welcome-content">
            <h2>您好，{{ userStore.userInfo.username }}！</h2>
            <p>欢迎来到ERP管理系统，这是一个基于Vue3 + Element Plus开发的企业资源管理平台。</p>
            <el-divider />
            <h3>主要功能</h3>
            <ul>
              <li>用户管理：管理系统用户，支持新增、编辑、删除和禁用操作</li>
              <li>角色管理：管理用户角色，为角色分配菜单权限</li>
              <li>菜单管理：管理系统菜单，支持多级菜单</li>
              <li>产品管理：管理产品信息，支持搜索和分页</li>
            </ul>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { getUsers } from '@/api/user'
import { getRoles } from '@/api/role'
import { getMenus } from '@/api/menu'
import { getProducts } from '@/api/product'

const userStore = useUserStore()

const stats = ref({
  users: 0,
  products: 0,
  roles: 0,
  menus: 0
})

const fetchStats = async () => {
  try {
    const [usersRes, rolesRes, menusRes, productsRes] = await Promise.all([
      getUsers(),
      getRoles(),
      getMenus(),
      getProducts()
    ])
    stats.value.users = Array.isArray(usersRes) ? usersRes.length : 0
    stats.value.roles = Array.isArray(rolesRes) ? rolesRes.length : 0
    stats.value.menus = Array.isArray(menusRes) ? menusRes.length : 0
    stats.value.products = Array.isArray(productsRes) ? productsRes.length : 0
  } catch (error) {
    console.error('获取统计数据失败', error)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.stat-card {
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-5px);
  transition: all 0.3s;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 5px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-content h2 {
  margin-bottom: 10px;
  color: #333;
}

.welcome-content p {
  color: #666;
  line-height: 1.6;
}

.welcome-content h3 {
  margin: 20px 0 10px;
  color: #333;
}

.welcome-content ul {
  padding-left: 20px;
  color: #666;
  line-height: 2;
}
</style>
