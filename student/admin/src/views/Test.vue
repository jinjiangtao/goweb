<template>
  <div class="test">
    <el-button @click="handleClick">Test Button</el-button>
  </div>
</template>

<script setup>
import { ElMessage } from "element-plus"

const handleClick = () => {
  console.log("Button clicked")
  ElMessage.success("Button clicked")
}
</script>

<style scoped>
.test {
  padding: 20px;
}
</style>
