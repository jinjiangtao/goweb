<template>
  <div class="house-card" @click="goDetail">
    <div class="card-image">
      <img :src="coverImage" :alt="house.title" />
      <div class="card-tags" v-if="hasTags">
        <span class="tag tag-recommend" v-if="house.is_recommended">推荐</span>
        <span class="tag tag-new" v-if="isNew">新上</span>
      </div>
    </div>
    <div class="card-content">
      <h3 class="card-title text-ellipsis">{{ house.title }}</h3>
      <div class="card-info">
        <span class="community">{{ house.community }}</span>
        <span class="district">{{ house.district }}</span>
      </div>
      <div class="card-meta">
        <span class="room-type">{{ house.room_type }}</span>
        <span class="area">{{ house.area }}㎡</span>
      </div>
      <div class="card-price">
        <span class="price">¥{{ formatPrice }}</span>
        <span class="unit">/月</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  house: {
    type: Object,
    required: true,
  },
})

const router = useRouter()

const coverImage = computed(() => {
  if (props.house.images && props.house.images.length > 0) {
    return props.house.images[0]
  }
  return 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 300"%3E%3Crect fill="%23f0f0f0" width="400" height="300"/%3E%3Ctext fill="%23999" font-family="sans-serif" font-size="14" x="50%25" y="50%25" text-anchor="middle" dy=".3em"%3E暂无图片%3C/text%3E%3C/svg%3E'
})

const formatPrice = computed(() => {
  const price = props.house.price
  if (price >= 10000) {
    return (price / 10000).toFixed(1) + '万'
  }
  return Math.round(price)
})

const isNew = computed(() => {
  const createdAt = new Date(props.house.created_at)
  const now = new Date()
  const diffDays = Math.floor((now - createdAt) / (1000 * 60 * 60 * 24))
  return diffDays <= 7
})

const hasTags = computed(() => {
  return props.house.is_recommended || isNew.value
})

function goDetail() {
  router.push(`/detail/${props.house.id}`)
}
</script>

<style scoped>
.house-card {
  background: var(--bg-white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: transform var(--transition-normal), box-shadow var(--transition-normal);
}

.house-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.card-image {
  position: relative;
  width: 100%;
  padding-top: 66.67%;
  overflow: hidden;
}

.card-image img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-normal);
}

.house-card:hover .card-image img {
  transform: scale(1.05);
}

.card-tags {
  position: absolute;
  top: 12px;
  left: 12px;
  display: flex;
  gap: 8px;
  z-index: 1;
}

.tag {
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  color: #fff;
  font-weight: 500;
}

.tag-recommend {
  background: var(--primary-color);
}

.tag-new {
  background: var(--success-color);
}

.card-content {
  padding: 12px 14px 16px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 8px;
}

.card-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.card-info .community::after {
  content: '·';
  margin-left: 8px;
  color: var(--text-light);
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 10px;
}

.card-price {
  display: flex;
  align-items: baseline;
}

.card-price .price {
  font-size: 22px;
  font-weight: 700;
  color: var(--price-color);
}

.card-price .unit {
  font-size: 13px;
  color: var(--text-secondary);
  margin-left: 4px;
}

@media (max-width: 480px) {
  .card-content {
    padding: 10px 12px 14px;
  }

  .card-title {
    font-size: 15px;
  }

  .card-price .price {
    font-size: 20px;
  }
}
</style>
