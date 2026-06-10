
&lt;template&gt;
  &lt;el-container style="height: 100vh"&gt;
    &lt;el-aside width="200px" style="background: #304156"&gt;
      &lt;div style="color: white; padding: 20px; font-size: 18px; font-weight: bold"&gt;抽奖管理&lt;/div&gt;
      &lt;el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      &gt;
        &lt;el-menu-item index="/stats"&gt;
          &lt;span&gt;统计看板&lt;/span&gt;
        &lt;/el-menu-item&gt;
        &lt;el-menu-item index="/prizes"&gt;
          &lt;span&gt;奖品管理&lt;/span&gt;
        &lt;/el-menu-item&gt;
        &lt;el-menu-item index="/records"&gt;
          &lt;span&gt;抽奖记录&lt;/span&gt;
        &lt;/el-menu-item&gt;
      &lt;/el-menu&gt;
    &lt;/el-aside&gt;
    &lt;el-container&gt;
      &lt;el-header style="background: white; border-bottom: 1px solid #eee; display: flex; justify-content: space-between; align-items: center"&gt;
        &lt;span style="font-size: 20px; font-weight: bold"&gt;抽奖转盘后台管理系统&lt;/span&gt;
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

const router = useRouter()
const route = useRoute()

const activeMenu = computed(() =&gt; route.path)

const handleLogout = () =&gt; {
  localStorage.removeItem('token')
  router.push('/login')
}
&lt;/script&gt;
