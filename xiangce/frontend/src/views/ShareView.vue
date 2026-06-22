<template>
  <div class="share-page">
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>

    <div v-else-if="error" class="error-state">
      <el-icon :size="64" color="#f56c6c"><Warning /></el-icon>
      <h2>{{ error }}</h2>
      <el-button type="primary" @click="goHome">返回首页</el-button>
    </div>

    <div v-else-if="needPassword && !verified" class="password-dialog">
      <div class="password-box">
        <el-icon :size="48" color="#409EFF"><Lock /></el-icon>
        <h2>相册已加密</h2>
        <p>请输入访问密码查看相册</p>
        
        <el-form :model="passwordForm" @submit.prevent="verifyPassword">
          <el-form-item>
            <el-input
              v-model="passwordForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              @keyup.enter="verifyPassword"
            />
          </el-form-item>
          <el-button type="primary" size="large" @click="verifyPassword" style="width: 100%">
            解锁查看
          </el-button>
        </el-form>
      </div>
    </div>

    <div v-else class="share-content">
      <div class="share-header">
        <div class="album-info">
          <h1 class="album-name">{{ album?.name }}</h1>
          <p class="album-desc">{{ album?.description || '暂无描述' }}</p>
          <span class="photo-count">{{ mediaList.length }} 张照片</span>
        </div>
        <div class="actions">
          <el-button type="primary" @click="startSlideshow">
            <el-icon><VideoPlay /></el-icon>
            幻灯片播放
          </el-button>
        </div>
      </div>

      <div class="media-grid">
        <div
          v-for="(item, index) in mediaList"
          :key="item.id"
          class="media-item"
          @click="openViewer(index)"
        >
          <div class="media-thumb">
            <img v-if="item.file_type === 'image'" :src="getThumbUrl(item.thumb_path)" />
            <div v-else class="video-thumb">
              <el-icon :size="40" color="#fff"><VideoPlay /></el-icon>
            </div>
          </div>
          <div class="media-name text-ellipsis">{{ item.original_name }}</div>
        </div>
      </div>
    </div>

    <el-dialog
      v-model="showViewer"
      width="90%"
      :close-on-click-modal="false"
      class="viewer-dialog"
      @close="stopSlideshow"
    >
      <div class="viewer-container">
        <div class="viewer-main">
          <div class="viewer-image">
            <img
              v-if="currentMedia?.file_type === 'image'"
              :src="getOriginalUrl(currentMedia?.file_path)"
              class="viewer-img"
            />
            <video
              v-else
              :src="getOriginalUrl(currentMedia?.file_path)"
              controls
              class="viewer-video"
            ></video>
          </div>
          
          <div class="viewer-nav prev" @click="prevImage">
            <el-icon :size="32"><ArrowLeft /></el-icon>
          </div>
          <div class="viewer-nav next" @click="nextImage">
            <el-icon :size="32"><ArrowRight /></el-icon>
          </div>
        </div>
        
        <div class="viewer-toolbar">
          <div class="toolbar-left">
            <span class="page-info">{{ currentIndex + 1 }} / {{ mediaList.length }}</span>
          </div>
          <div class="toolbar-center">
            <el-button circle @click="prevImage">
              <el-icon><DArrowLeft /></el-icon>
            </el-button>
            <el-button circle type="primary" @click="toggleSlideshow">
              <el-icon v-if="isPlaying"><VideoPause /></el-icon>
              <el-icon v-else><VideoPlay /></el-icon>
            </el-button>
            <el-button circle @click="nextImage">
              <el-icon><DArrowRight /></el-icon>
            </el-button>
          </div>
          <div class="toolbar-right">
            <el-button circle @click="downloadCurrent">
              <el-icon><Download /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getShareByToken, verifySharePassword, recordShareVisit } from '@/utils/api'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const error = ref('')
const album = ref(null)
const mediaList = ref([])
const needPassword = ref(false)
const verified = ref(false)
const shareToken = computed(() => route.params.token)

const showViewer = ref(false)
const currentIndex = ref(0)
const currentMedia = computed(() => mediaList.value[currentIndex.value])
const isPlaying = ref(false)
let slideshowTimer = null

const passwordForm = ref({
  password: ''
})

const getThumbUrl = (path) => {
  if (!path) return ''
  const filename = path.split('/').pop()
  return `/uploads/thumbnails/${filename}`
}

const getOriginalUrl = (path) => {
  if (!path) return ''
  const filename = path.split('/').pop()
  return `/uploads/${filename}`
}

const loadShare = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const res = await getShareByToken(shareToken.value)
    album.value = res.album
    mediaList.value = res.media
    needPassword.value = res.share.need_password
    
    if (!needPassword.value) {
      verified.value = true
      recordVisit()
    }
  } catch (err) {
    error.value = err.response?.data?.error || '加载失败'
  } finally {
    loading.value = false
  }
}

const verifyPassword = async () => {
  if (!passwordForm.value.password) {
    ElMessage.warning('请输入密码')
    return
  }
  
  try {
    await verifySharePassword(shareToken.value, passwordForm.value.password)
    verified.value = true
    recordVisit()
  } catch {}
}

const recordVisit = async () => {
  try {
    await recordShareVisit(shareToken.value)
  } catch {}
}

const openViewer = (index) => {
  currentIndex.value = index
  showViewer.value = true
}

const prevImage = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--
  } else {
    currentIndex.value = mediaList.value.length - 1
  }
}

const nextImage = () => {
  if (currentIndex.value < mediaList.value.length - 1) {
    currentIndex.value++
  } else {
    currentIndex.value = 0
  }
}

const startSlideshow = () => {
  if (mediaList.value.length === 0) return
  currentIndex.value = 0
  showViewer.value = true
  isPlaying.value = true
  startTimer()
}

const toggleSlideshow = () => {
  isPlaying.value = !isPlaying.value
  if (isPlaying.value) {
    startTimer()
  } else {
    stopTimer()
  }
}

const startTimer = () => {
  stopTimer()
  slideshowTimer = setInterval(() => {
    nextImage()
  }, 3000)
}

const stopTimer = () => {
  if (slideshowTimer) {
    clearInterval(slideshowTimer)
    slideshowTimer = null
  }
}

const stopSlideshow = () => {
  stopTimer()
  isPlaying.value = false
}

const downloadCurrent = () => {
  if (currentMedia.value) {
    window.open(getOriginalUrl(currentMedia.value.file_path), '_blank')
  }
}

const goHome = () => {
  router.push('/')
}

onMounted(() => {
  loadShare()
})

onUnmounted(() => {
  stopTimer()
})
</script>

<style scoped>
.share-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding: 20px;
}

.loading-state,
.error-state {
  min-height: 60vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #909399;
}

.error-state h2 {
  color: #f56c6c;
  margin: 0;
}

.password-dialog {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.password-box {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.password-box h2 {
  margin: 16px 0 8px;
  color: #303133;
}

.password-box p {
  color: #909399;
  margin-bottom: 24px;
}

.share-content {
  max-width: 1400px;
  margin: 0 auto;
}

.share-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 24px;
  padding: 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.album-name {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px;
  color: #303133;
}

.album-desc {
  color: #606266;
  margin: 0 0 8px;
}

.photo-count {
  color: #909399;
  font-size: 14px;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.media-item {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.media-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.media-thumb {
  position: relative;
  width: 100%;
  padding-top: 100%;
  background: #f0f2f5;
  overflow: hidden;
}

.media-thumb img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-thumb {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.media-name {
  padding: 12px;
  font-size: 13px;
  color: #606266;
}

.viewer-dialog :deep(.el-dialog) {
  max-width: 90vw;
  height: 90vh;
  margin: 5vh auto;
}

.viewer-dialog :deep(.el-dialog__body) {
  padding: 0;
  height: calc(90vh - 40px);
}

.viewer-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1a1a1a;
}

.viewer-main {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.viewer-image {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.viewer-img,
.viewer-video {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.viewer-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  border-radius: 50%;
  cursor: pointer;
  transition: background 0.2s;
}

.viewer-nav:hover {
  background: rgba(0, 0, 0, 0.7);
}

.viewer-nav.prev {
  left: 20px;
}

.viewer-nav.next {
  right: 20px;
}

.viewer-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
}

.toolbar-left,
.toolbar-center,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-info {
  font-size: 14px;
}

.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
