<template>
  <div class="image-carousel" ref="carouselRef">
    <div
      class="carousel-container"
      :style="containerStyle"
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
    >
      <div
        v-for="(image, index) in images"
        :key="index"
        class="carousel-slide"
      >
        <img :src="image" :alt="`图片${index + 1}`" />
      </div>
    </div>

    <button
      v-if="images.length > 1"
      class="carousel-btn prev-btn"
      @click="prevSlide"
    >
      <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="15 18 9 12 15 6"></polyline>
      </svg>
    </button>

    <button
      v-if="images.length > 1"
      class="carousel-btn next-btn"
      @click="nextSlide"
    >
      <svg viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="9 18 15 12 9 6"></polyline>
      </svg>
    </button>

    <div v-if="images.length > 1" class="carousel-dots">
      <span
        v-for="(_, index) in images"
        :key="index"
        class="dot"
        :class="{ active: index === currentIndex }"
        @click="goToSlide(index)"
      ></span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  images: {
    type: Array,
    default: () => [],
  },
  autoPlay: {
    type: Boolean,
    default: false,
  },
  interval: {
    type: Number,
    default: 4000,
  },
})

const currentIndex = ref(0)
let touchStartX = 0
let touchEndX = 0
let autoPlayTimer = null

const containerStyle = computed(() => ({
  transform: `translateX(-${currentIndex.value * 100}%)`,
}))

function nextSlide() {
  if (currentIndex.value < props.images.length - 1) {
    currentIndex.value++
  } else {
    currentIndex.value = 0
  }
}

function prevSlide() {
  if (currentIndex.value > 0) {
    currentIndex.value--
  } else {
    currentIndex.value = props.images.length - 1
  }
}

function goToSlide(index) {
  currentIndex.value = index
}

function handleTouchStart(e) {
  touchStartX = e.touches[0].clientX
}

function handleTouchMove(e) {
  touchEndX = e.touches[0].clientX
}

function handleTouchEnd() {
  const diff = touchStartX - touchEndX
  if (Math.abs(diff) > 50) {
    if (diff > 0) {
      nextSlide()
    } else {
      prevSlide()
    }
  }
}

function startAutoPlay() {
  if (props.autoPlay && props.images.length > 1) {
    autoPlayTimer = setInterval(nextSlide, props.interval)
  }
}

function stopAutoPlay() {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer)
    autoPlayTimer = null
  }
}

onMounted(() => {
  startAutoPlay()
})

onUnmounted(() => {
  stopAutoPlay()
})
</script>

<style scoped>
.image-carousel {
  position: relative;
  width: 100%;
  overflow: hidden;
  background: #f0f0f0;
}

.carousel-container {
  display: flex;
  transition: transform 0.4s ease;
  touch-action: pan-y;
}

.carousel-slide {
  flex-shrink: 0;
  width: 100%;
}

.carousel-slide img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.carousel-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.4);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background var(--transition-fast);
  z-index: 2;
}

.carousel-btn:hover {
  background: rgba(0, 0, 0, 0.6);
}

.prev-btn {
  left: 16px;
}

.next-btn {
  right: 16px;
}

.carousel-dots {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  z-index: 2;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.dot.active {
  background: #fff;
  width: 24px;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .carousel-btn {
    width: 32px;
    height: 32px;
  }

  .prev-btn {
    left: 8px;
  }

  .next-btn {
    right: 8px;
  }
}
</style>
