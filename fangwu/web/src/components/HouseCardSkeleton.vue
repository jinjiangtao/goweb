<template>
  <div class="skeleton-card">
    <div class="skeleton-image skeleton"></div>
    <div class="skeleton-content">
      <div class="skeleton-title skeleton"></div>
      <div class="skeleton-info skeleton"></div>
      <div class="skeleton-meta skeleton"></div>
      <div class="skeleton-price skeleton"></div>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.skeleton-card {
  background: var(--bg-white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.skeleton-image {
  width: 100%;
  padding-top: 66.67%;
}

.skeleton-content {
  padding: 12px 14px 16px;
}

.skeleton-title {
  height: 20px;
  width: 80%;
  border-radius: 4px;
  margin-bottom: 10px;
}

.skeleton-info {
  height: 16px;
  width: 60%;
  border-radius: 4px;
  margin-bottom: 10px;
}

.skeleton-meta {
  height: 16px;
  width: 50%;
  border-radius: 4px;
  margin-bottom: 12px;
}

.skeleton-price {
  height: 24px;
  width: 35%;
  border-radius: 4px;
}
</style>
