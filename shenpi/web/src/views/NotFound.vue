<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/')
}
</script>

<template>
  <div class="not-found">
    <div class="content">
      <h1 class="code">404</h1>
      <p class="message">抱歉，您访问的页面不存在</p>
      <el-button type="primary" @click="goHome">返回首页</el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.not-found {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f2f5;

  .content {
    text-align: center;

    .code {
      font-size: 120px;
      color: #409eff;
      margin: 0;
      line-height: 1;
    }

    .message {
      font-size: 20px;
      color: #909399;
      margin: 20px 0 30px;
    }
  }
}
</style>
