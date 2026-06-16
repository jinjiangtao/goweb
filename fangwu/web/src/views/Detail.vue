<template>
  <div class="detail-page" v-if="house">
    <div class="detail-header">
      <button class="back-btn" @click="goBack">
        <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
        返回
      </button>
    </div>

    <div class="carousel-wrapper">
      <ImageCarousel :images="houseImages" :auto-play="true" />
    </div>

    <main class="detail-main">
      <section class="info-section">
        <h1 class="house-title">{{ house.title }}</h1>
        <div class="house-price">
          <span class="price">¥{{ formatPrice }}</span>
          <span class="unit">元/月</span>
        </div>

        <div class="house-tags">
          <span class="tag-item" v-if="house.room_type">{{ house.room_type }}</span>
          <span class="tag-divider">|</span>
          <span class="tag-item" v-if="house.area">{{ house.area }}㎡</span>
          <span class="tag-divider" v-if="house.floor">|</span>
          <span class="tag-item" v-if="house.floor">{{ house.floor }}</span>
          <span class="tag-divider" v-if="house.orientation">|</span>
          <span class="tag-item" v-if="house.orientation">{{ house.orientation }}</span>
          <span class="tag-divider" v-if="house.decoration">|</span>
          <span class="tag-item" v-if="house.decoration">{{ house.decoration }}</span>
        </div>

        <div class="house-address">
          <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
            <circle cx="12" cy="10" r="3"></circle>
          </svg>
          <span>{{ house.district }} · {{ house.community }} · {{ house.address }}</span>
        </div>
      </section>

      <section class="facilities-section">
        <h2 class="section-title">配套设施</h2>
        <div class="facilities-grid">
          <div
            v-for="facility in allFacilities"
            :key="facility.name"
            class="facility-item"
            :class="{ active: hasFacility(facility.name) }"
          >
            <div class="facility-icon" v-html="facility.icon"></div>
            <span class="facility-name">{{ facility.label }}</span>
          </div>
        </div>
      </section>

      <section class="description-section">
        <h2 class="section-title">房源描述</h2>
        <div class="description-content">
          <p>{{ house.description }}</p>
        </div>
      </section>

      <section class="contact-section">
        <h2 class="section-title">联系房东</h2>
        <div class="contact-card">
          <div class="contact-avatar">
            <svg viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
          </div>
          <div class="contact-info">
            <div class="contact-name">{{ house.contact_name }}</div>
            <div class="contact-phone">
              <span>{{ displayPhone }}</span>
              <button v-if="!showPhone" class="show-phone-btn" @click="handleShowPhone">
                显示号码
              </button>
            </div>
          </div>
          <a :href="`tel:${house.contact_phone}`" v-if="showPhone" class="call-btn">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
            </svg>
            拨打电话
          </a>
        </div>
      </section>

      <section class="recommend-section" v-if="recommendHouses.length > 0">
        <h2 class="section-title">推荐房源</h2>
        <div class="recommend-scroll">
          <div class="recommend-list">
            <div
              v-for="item in recommendHouses"
              :key="item.id"
              class="recommend-card"
              @click="goDetail(item.id)"
            >
              <div class="recommend-image">
                <img :src="getCoverImage(item)" :alt="item.title" />
              </div>
              <div class="recommend-info">
                <h3 class="recommend-title text-ellipsis">{{ item.title }}</h3>
                <p class="recommend-meta">{{ item.room_type }} · {{ item.area }}㎡</p>
                <p class="recommend-price">¥{{ Math.round(item.price) }}<span>/月</span></p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <div class="view-count">
        浏览次数：{{ viewCount }} 次
      </div>
    </main>
  </div>

  <div class="detail-page loading" v-else>
    <div class="skeleton-carousel skeleton"></div>
    <div class="detail-main">
      <div class="skeleton-title skeleton"></div>
      <div class="skeleton-price skeleton"></div>
      <div class="skeleton-tags skeleton"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getHouseDetail, recordConsult, getRecommendHouses, incrementViewCount } from '../api/house'
import { useUserStore } from '../stores/user'
import ImageCarousel from '../components/ImageCarousel.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const house = ref(null)
const showPhone = ref(false)
const viewCount = ref(0)
const recommendHouses = ref([])
const facilitiesList = ref([])

const allFacilities = [
  { name: 'wifi', label: 'WiFi', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12.55a11 11 0 0 1 14.08 0"></path><path d="M1.42 9a16 16 0 0 1 21.16 0"></path><path d="M8.53 16.11a6 6 0 0 1 6.95 0"></path><line x1="12" y1="20" x2="12.01" y2="20"></line></svg>' },
  { name: '空调', label: '空调', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="5" width="20" height="10" rx="2"></rect><line x1="6" y1="10" x2="6.01" y2="10"></line><line x1="10" y1="10" x2="10.01" y2="10"></line><line x1="14" y1="10" x2="14.01" y2="10"></line><path d="M8 15v5"></path><path d="M16 15v5"></path></svg>' },
  { name: '冰箱', label: '冰箱', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><rect x="5" y="2" width="14" height="20" rx="2"></rect><path d="M5 10h14"></path><circle cx="9" cy="6" r="1"></circle></svg>' },
  { name: '洗衣机', label: '洗衣机', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="2" width="18" height="20" rx="2"></rect><circle cx="12" cy="13" r="5"></circle><circle cx="12" cy="13" r="1"></circle></svg>' },
  { name: '暖气', label: '暖气', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2v6"></path><path d="M12 22v-2"></path><path d="M5 8a7 7 0 0 1 14 0"></path><path d="M3 12a9 9 0 0 1 18 0"></path></svg>' },
  { name: '热水器', label: '热水器', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="12" rx="2"></rect><path d="M8 20v-4"></path><path d="M16 20v-4"></path><line x1="8" y1="10" x2="16" y2="10"></line></svg>' },
  { name: '床', label: '床', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><path d="M2 4v16"></path><path d="M2 8h18a2 2 0 0 1 2 2v10"></path><path d="M2 17h20"></path><path d="M6 8v9"></path></svg>' },
  { name: '衣柜', label: '衣柜', icon: '<svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="2" width="18" height="20" rx="2"></rect><line x1="12" y1="2" x2="12" y2="22"></line><circle cx="9" cy="12" r="1"></circle><circle cx="15" cy="12" r="1"></circle></svg>' },
]

const houseImages = computed(() => {
  if (house.value && house.value.images && house.value.images.length > 0) {
    return house.value.images
  }
  return ['data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 600"%3E%3Crect fill="%23f0f0f0" width="800" height="600"/%3E%3Ctext fill="%23999" font-family="sans-serif" font-size="24" x="50%25" y="50%25" text-anchor="middle" dy=".3em"%3E暂无图片%3C/text%3E%3C/svg%3E']
})

const formatPrice = computed(() => {
  if (!house.value) return 0
  const price = house.value.price
  if (price >= 10000) {
    return (price / 10000).toFixed(1) + '万'
  }
  return Math.round(price)
})

const displayPhone = computed(() => {
  if (!house.value) return ''
  if (showPhone.value) {
    return house.value.contact_phone
  }
  const phone = house.value.contact_phone
  if (phone.length >= 7) {
    return phone.slice(0, 3) + '****' + phone.slice(-4)
  }
  return '****'
})

function hasFacility(name) {
  if (!facilitiesList.value || facilitiesList.value.length === 0) return false
  return facilitiesList.value.includes(name)
}

function getCoverImage(item) {
  if (item.images && item.images.length > 0) {
    return item.images[0]
  }
  return 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 300"%3E%3Crect fill="%23f0f0f0" width="400" height="300"/%3E%3Ctext fill="%23999" font-family="sans-serif" font-size="14" x="50%25" y="50%25" text-anchor="middle" dy=".3em"%3E暂无图片%3C/text%3E%3C/svg%3E'
}

function goBack() {
  router.back()
}

function goDetail(id) {
  router.push(`/detail/${id}`)
}

async function handleShowPhone() {
  if (!userStore.isLogin) {
    alert('请先登录后查看完整号码')
    return
  }

  if (showPhone.value) return

  try {
    await recordConsult(route.params.id)
    showPhone.value = true
  } catch (error) {
    console.error('记录咨询失败:', error)
    showPhone.value = true
  }
}

async function fetchHouseDetail() {
  try {
    const res = await getHouseDetail(route.params.id)
    house.value = res.data

    if (res.data.facilities) {
      facilitiesList.value = res.data.facilities.split(',').filter(f => f.trim() !== '')
    }

    viewCount.value = res.data.view_count || Math.floor(Math.random() * 500) + 50

    try {
      await incrementViewCount(route.params.id)
    } catch (e) {
      console.log('增加浏览量失败')
    }
  } catch (error) {
    console.error('获取房源详情失败:', error)
  }
}

async function fetchRecommendHouses() {
  try {
    const res = await getRecommendHouses(route.params.id)
    if (Array.isArray(res.data)) {
      recommendHouses.value = res.data
    } else if (res.data && res.data.list) {
      recommendHouses.value = res.data.list
    } else {
      recommendHouses.value = []
    }
  } catch (error) {
    console.error('获取推荐房源失败:', error)
    recommendHouses.value = []
  }
}

onMounted(() => {
  fetchHouseDetail()
  fetchRecommendHouses()
})
</script>

<style scoped>
.detail-page {
  min-height: 100vh;
  background: var(--bg-color);
  padding-bottom: 40px;
}

.detail-page.loading {
  background: var(--bg-white);
}

.detail-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: 12px 16px;
  background: linear-gradient(to bottom, rgba(0, 0, 0, 0.5), transparent);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  color: var(--text-color);
  font-size: 14px;
  transition: background var(--transition-fast);
}

.back-btn:hover {
  background: #fff;
}

.carousel-wrapper {
  width: 100%;
  aspect-ratio: 16 / 10;
  max-height: 450px;
  background: #f0f0f0;
}

.detail-main {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 16px;
}

.info-section {
  background: var(--bg-white);
  margin: -20px 0 12px;
  border-radius: var(--radius-xl);
  padding: 24px 20px;
  position: relative;
  box-shadow: var(--shadow-sm);
}

.house-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-color);
  margin-bottom: 12px;
  line-height: 1.4;
}

.house-price {
  display: flex;
  align-items: baseline;
  margin-bottom: 16px;
}

.house-price .price {
  font-size: 32px;
  font-weight: 700;
  color: var(--price-color);
}

.house-price .unit {
  font-size: 14px;
  color: var(--text-secondary);
  margin-left: 4px;
}

.house-tags {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.tag-item {
  font-size: 14px;
  color: var(--text-secondary);
}

.tag-divider {
  color: var(--border-color);
}

.house-address {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: var(--text-secondary);
}

.house-address svg {
  color: var(--primary-color);
  flex-shrink: 0;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 16px;
  padding-left: 10px;
  border-left: 4px solid var(--primary-color);
}

.facilities-section,
.description-section,
.contact-section,
.recommend-section {
  background: var(--bg-white);
  border-radius: var(--radius-xl);
  padding: 20px;
  margin-bottom: 12px;
  box-shadow: var(--shadow-sm);
}

.facilities-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.facility-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: #ccc;
  transition: color var(--transition-fast);
}

.facility-item.active {
  color: var(--primary-color);
}

.facility-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-color);
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
}

.facility-item.active .facility-icon {
  background: rgba(255, 107, 53, 0.1);
}

.facility-name {
  font-size: 12px;
}

.description-content {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.8;
}

.description-content p {
  white-space: pre-wrap;
  word-break: break-all;
}

.contact-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--bg-color);
  border-radius: var(--radius-lg);
}

.contact-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--primary-color);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.contact-info {
  flex: 1;
  min-width: 0;
}

.contact-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 4px;
}

.contact-phone {
  font-size: 14px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 10px;
}

.show-phone-btn {
  padding: 4px 12px;
  background: var(--primary-color);
  color: #fff;
  border-radius: 12px;
  font-size: 12px;
}

.call-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: var(--primary-color);
  color: #fff;
  border-radius: var(--radius-md);
  font-size: 14px;
  flex-shrink: 0;
  transition: background var(--transition-fast);
}

.call-btn:hover {
  background: var(--primary-dark);
}

.recommend-scroll {
  overflow-x: auto;
  margin: 0 -20px;
  padding: 0 20px;
  -webkit-overflow-scrolling: touch;
}

.recommend-scroll::-webkit-scrollbar {
  display: none;
}

.recommend-list {
  display: flex;
  gap: 12px;
  padding-bottom: 4px;
}

.recommend-card {
  flex-shrink: 0;
  width: 200px;
  background: var(--bg-white);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: transform var(--transition-fast);
}

.recommend-card:hover {
  transform: translateY(-2px);
}

.recommend-image {
  width: 100%;
  height: 130px;
  overflow: hidden;
}

.recommend-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-normal);
}

.recommend-card:hover .recommend-image img {
  transform: scale(1.05);
}

.recommend-info {
  padding: 10px 12px;
}

.recommend-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 4px;
}

.recommend-meta {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.recommend-price {
  font-size: 16px;
  font-weight: 700;
  color: var(--price-color);
}

.recommend-price span {
  font-size: 12px;
  font-weight: normal;
  color: var(--text-secondary);
}

.view-count {
  text-align: center;
  padding: 20px 0;
  font-size: 13px;
  color: var(--text-light);
}

.skeleton-carousel {
  width: 100%;
  aspect-ratio: 16 / 10;
  max-height: 450px;
}

.skeleton-title {
  height: 28px;
  width: 80%;
  border-radius: 6px;
  margin: 24px 0 16px;
}

.skeleton-price {
  height: 36px;
  width: 30%;
  border-radius: 6px;
  margin-bottom: 16px;
}

.skeleton-tags {
  height: 18px;
  width: 60%;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .info-section {
    padding: 20px 16px;
  }

  .house-title {
    font-size: 18px;
  }

  .house-price .price {
    font-size: 26px;
  }

  .facilities-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
  }

  .facility-icon {
    width: 36px;
    height: 36px;
  }

  .recommend-card {
    width: 160px;
  }

  .recommend-image {
    height: 100px;
  }

  .contact-card {
    flex-wrap: wrap;
  }

  .call-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .facilities-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;
  }

  .facility-icon {
    width: 32px;
    height: 32px;
  }

  .facility-name {
    font-size: 11px;
  }
}
</style>
