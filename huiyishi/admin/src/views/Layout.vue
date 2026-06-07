
&lt;template&gt;
  &lt;el-container class="layout"&gt;
    &lt;el-aside width="220px"&gt;
      &lt;div class="logo"&gt;会议室预订&lt;/div&gt;
      &lt;el-menu
        :default-active="activeMenu"
        router
        background-color="transparent"
        text-color="#cbd5e1"
        active-text-color="#3b82f6"
      &gt;
        &lt;el-menu-item index="/stats"&gt;
          &lt;el-icon&gt;&lt;DataLine /&gt;&lt;/el-icon&gt;
          &lt;span&gt;统计看板&lt;/span&gt;
        &lt;/el-menu-item&gt;
        &lt;el-menu-item index="/users"&gt;
          &lt;el-icon&gt;&lt;User /&gt;&lt;/el-icon&gt;
          &lt;span&gt;用户管理&lt;/span&gt;
        &lt;/el-menu-item&gt;
        &lt;el-menu-item index="/rooms"&gt;
          &lt;el-icon&gt;&lt;OfficeBuilding /&gt;&lt;/el-icon&gt;
          &lt;span&gt;会议室管理&lt;/span&gt;
        &lt;/el-menu-item&gt;
        &lt;el-menu-item index="/bookings"&gt;
          &lt;el-icon&gt;&lt;Calendar /&gt;&lt;/el-icon&gt;
          &lt;span&gt;预订管理&lt;/span&gt;
        &lt;/el-menu-item&gt;
      &lt;/el-menu&gt;
    &lt;/el-aside&gt;
    &lt;el-container&gt;
      &lt;el-header&gt;
        &lt;span class="welcome"&gt;欢迎, {{ authStore.admin?.nickname }}&lt;/span&gt;
        &lt;el-button type="danger" @click="handleLogout"&gt;退出登录&lt;/el-button&gt;
      &lt;/el-header&gt;
      &lt;el-main&gt;
        &lt;router-view /&gt;
      &lt;/el-main&gt;
    &lt;/el-container&gt;
  &lt;/el-container&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { DataLine, OfficeBuilding, Calendar, User } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeMenu = computed(() =&gt; route.path)

const handleLogout = () =&gt; {
  authStore.logout()
  router.push('/login')
}
&lt;/script&gt;

&lt;style scoped&gt;
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
&lt;/style&gt;
