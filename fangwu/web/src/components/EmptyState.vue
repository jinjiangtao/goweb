<template>
  <div class="empty-state">
    <div class="empty-icon">
      <svg viewBox="0 0 200 200" width="120" height="120">
        <path
          fill="#f0f0f0"
          d="M100 40 L160 80 L160 160 L40 160 L40 80 Z"
        />
        <path
          fill="#e0e0e0"
          d="M100 40 L40 80 L160 80 Z"
        />
        <rect fill="#d0d0d0" x="70" y="100" width="25" height="35" rx="2" />
        <rect fill="#d0d0d0" x="105" y="100" width="25" height="35" rx="2" />
        <rect fill="#d0d0d0" x="85" y="140" width="30" height="20" rx="2" />
      </svg>
    </div>
    <p class="empty-text">{{ text }}</p>
    <slot v-if="$slots.default"></slot>
  </div>
</template>

<script setup>
defineProps({
  text: {
    type: String,
    default: '暂无数据',
  },
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  margin-bottom: 16px;
  opacity: 0.8;
}

.empty-text {
  font-size: 15px;
  color: var(--text-light);
}
</style>
