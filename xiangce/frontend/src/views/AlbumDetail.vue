<template>
  <Layout>
    <div class="album-detail-page">
      <div class="page-header">
        <div class="album-info">
          <el-button link @click="$router.back()">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h2 class="album-title">{{ album?.name || '加载中...' }}</h2>
          <span class="photo-count">{{ mediaList.length }} 张照片</span>
        </div>
        <div class="header-actions">
          <el-upload
            ref="uploadRef"
            :action="uploadUrl"
            :headers="uploadHeaders"
            :data="uploadData"
            :show-file-list="false"
            multiple
            accept="image/*,video/*"
            name="file"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            drag
            class="upload-btn"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              上传照片
            </el-button>
          </el-upload>
          
          <el-dropdown>
            <el-button>
              <el-icon><Sort /></el-icon>
              排序
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="changeSort('sort_order')">默认排序</el-dropdown-item>
                <el-dropdown-item @click="changeSort('created_at')">上传时间</el-dropdown-item>
                <el-dropdown-item @click="changeSort('file_name')">文件名称</el-dropdown-item>
                <el-dropdown-item divided @click="toggleDragSort">
                  {{ dragSortMode ? '完成拖拽' : '拖拽排序' }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <el-button type="success" @click="showShareDialog = true">
            <el-icon><Share /></el-icon>
            分享相册
          </el-button>
          
          <el-button v-if="selectedMedia.length > 0" type="danger" @click="batchDelete">
            <el-icon><Delete /></el-icon>
            删除 ({{ selectedMedia.length }})
          </el-button>
        </div>
      </div>

      <div v-if="loading" class="loading-state">
        <el-icon class="is-loading"><Loading /></el-icon>
        <span>加载中...</span>
      </div>

      <div v-else-if="mediaList.length === 0" class="empty-state">
        <el-empty description="相册里还没有照片，快去上传吧" />
        <el-upload
          :action="uploadUrl"
          :headers="uploadHeaders"
          :data="uploadData"
          :show-file-list="false"
          multiple
          accept="image/*,video/*"
          name="file"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :before-upload="beforeUpload"
          drag
          class="upload-area"
        >
          <div class="upload-content">
            <el-icon class="upload-icon"><UploadFilled /></el-icon>
            <p>点击或拖拽文件到此处上传</p>
            <p class="upload-tip">支持 JPG、PNG、MP4 等格式</p>
          </div>
        </el-upload>
      </div>

      <draggable
        v-else
        v-model="mediaList"
        item-key="id"
        :disabled="!dragSortMode"
        class="media-grid"
        :animation="200"
        @end="onMediaSortEnd"
      >
        <template #item="{ element, index }">
          <div
            class="media-item"
            :class="{ selected: isSelected(element.id), 'drag-mode': dragSortMode }"
            @click="handleMediaClick(element, index)"
          >
            <div class="media-checkbox" v-if="!dragSortMode" @click.stop="toggleSelect(element.id)">
              <el-checkbox :model-value="isSelected(element.id)" />
            </div>
            <div class="media-thumb">
              <img v-if="element.file_type === 'image'" :src="getThumbUrl(element.thumb_path)" />
              <div v-else class="video-thumb">
                <el-icon :size="48" color="#fff"><VideoPlay /></el-icon>
              </div>
            </div>
            <div class="media-name text-ellipsis">{{ element.original_name }}</div>
          </div>
        </template>
      </draggable>
    </div>

    <el-dialog
      v-model="showViewer"
      :title="currentMedia?.original_name"
      width="90%"
      :close-on-click-modal="false"
      class="viewer-dialog"
    >
      <div class="viewer-container">
        <div class="viewer-main">
          <div class="viewer-image" @click="toggleSlideshow ? stopSlideshow : null">
            <img
              v-if="currentMedia?.file_type === 'image'"
              :src="getOriginalUrl(currentMedia?.file_path)"
              :style="viewerStyle"
              @wheel.stop="handleZoom"
              @mousedown.stop="startDrag"
              @mousemove.stop="onDrag"
              @mouseup.stop="endDrag"
              @mouseleave.stop="endDrag"
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
              <el-icon v-if="slideshow"><VideoPause /></el-icon>
              <el-icon v-else><VideoPlay /></el-icon>
            </el-button>
            <el-button circle @click="nextImage">
              <el-icon><DArrowRight /></el-icon>
            </el-button>
          </div>
          <div class="toolbar-right">
            <el-button circle @click="zoomOut">
              <el-icon><ZoomOut /></el-icon>
            </el-button>
            <span class="zoom-level">{{ Math.round(scale * 100) }}%</span>
            <el-button circle @click="zoomIn">
              <el-icon><ZoomIn /></el-icon>
            </el-button>
            <el-button circle @click="resetView">
              <el-icon><RefreshRight /></el-icon>
            </el-button>
            <el-button circle @click="downloadCurrent">
              <el-icon><Download /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>

    <el-dialog v-model="showShareDialog" title="分享相册" width="500px">
      <el-form :model="shareForm" label-width="100px">
        <el-form-item label="设置密码">
          <el-switch v-model="shareForm.hasPassword" />
        </el-form-item>
        <el-form-item label="访问密码" v-if="shareForm.hasPassword">
          <el-input v-model="shareForm.password" type="password" placeholder="设置访问密码" />
        </el-form-item>
        <el-form-item label="有效期">
          <el-radio-group v-model="shareForm.expireType">
            <el-radio :value="0">永久有效</el-radio>
            <el-radio :value="7">7天</el-radio>
            <el-radio :value="30">30天</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="最大访问次数">
          <el-input-number v-model="shareForm.maxViews" :min="0" :max="9999" />
          <span class="form-tip">0 表示不限制</span>
        </el-form-item>
      </el-form>
      
      <div v-if="createdShare" class="share-result">
        <el-divider />
        <el-alert type="success" :closable="false">
          <p>分享链接已生成：</p>
          <div class="share-link">
            <el-input :value="shareLinkUrl" readonly />
            <el-button type="primary" @click="copyShareLink">复制链接</el-button>
          </div>
        </el-alert>
      </div>
      
      <template #footer>
        <el-button @click="showShareDialog = false">关闭</el-button>
        <el-button type="primary" @click="createShare" :loading="creatingShare">
          生成分享链接
        </el-button>
      </template>
    </el-dialog>
  </Layout>
</template>

<script setup>import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import draggable from 'vuedraggable'
import Layout from '@/components/Layout.vue'
import { getAlbum, getMediaList, deleteMedia, batchDeleteMedia, updateMediaSort, downloadMedia, createShareLink } from '@/utils/api'
import { getThumbUrl, getOriginalUrl } from '@/utils/url'

const route = useRoute()
const router = useRouter()

const album = ref(null)
const mediaList = ref([])
const loading = ref(false)
const selectedMedia = ref([])
const sortBy = ref('sort_order')
const sortOrder = ref('asc')
const dragSortMode = ref(false)

const showShareDialog = ref(false)
const creatingShare = ref(false)
const createdShare = ref(null)
const shareForm = reactive({
  hasPassword: false,
  password: '',
  expireType: 0,
  maxViews: 0
})

const shareLinkUrl = computed(() => {
  if (!createdShare.value) return ''
  return `${window.location.origin}/share/${createdShare.value.token}`
})

const showViewer = ref(false)
const currentIndex = ref(0)
const currentMedia = computed(() => mediaList.value[currentIndex.value])

const scale = ref(1)
const position = reactive({ x: 0, y: 0 })
const isDragging = ref(false)
const dragStart = reactive({ x: 0, y: 0 })
const slideshow = ref(false)
let slideshowTimer = null

const uploadRef = ref(null)

const albumId = computed(() => parseInt(route.params.id))

const uploadUrl = '/api/media/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))
const uploadData = computed(() => ({
  album_id: albumId.value
}))

const viewerStyle = computed(() => ({
  transform: `scale(${scale.value}) translate(${position.value.x / scale.value}px, ${position.value.y / scale.value}px)`,
  cursor: isDragging.value ? 'grabbing' : 'grab'
}))

const isSelected = (id) => selectedMedia.value.includes(id)

const toggleSelect = (id) => {
  const idx = selectedMedia.value.indexOf(id)
  if (idx > -1) {
    selectedMedia.value.splice(idx, 1)
  } else {
    selectedMedia.value.push(id)
  }
}

const loadAlbum = async () => {
  try {
    const res = await getAlbum(albumId.value)
    album.value = res.album
  } catch {}
}

const loadMedia = async () => {
  loading.value = true
  try {
    const res = await getMediaList(albumId.value, {
      sort_by: sortBy.value,
      sort_order: sortOrder.value
    })
    mediaList.value = res.media
  } catch {
  } finally {
    loading.value = false
  }
}

const handleMediaClick = (media, index) => {
  if (dragSortMode.value) return
  
  currentIndex.value = index
  showViewer.value = true
  resetView()
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isVideo = file.type.startsWith('video/')
  const isLt100M = file.size / 1024 / 1024 < 100
  
  if (!isImage && !isVideo) {
    ElMessage.error('只能上传图片或视频文件!')
    return false
  }
  if (!isLt100M) {
    ElMessage.error('文件大小不能超过 100MB!')
    return false
  }
  return true
}

const handleUploadSuccess = (response) => {
  ElMessage.success(response.message || '上传成功')
  loadMedia()
  loadAlbum()
}

const handleUploadError = (error) => {
  try {
    const err = JSON.parse(error.message)
    ElMessage.error(err.error || '上传失败')
  } catch {
    ElMessage.error('上传失败，请重试')
  }
}

const changeSort = (field) => {
  sortBy.value = field
  loadMedia()
}

const toggleDragSort = () => {
  dragSortMode.value = !dragSortMode.value
  selectedMedia.value = []
}

const onMediaSortEnd = async () => {
  try {
    const mediaIds = mediaList.value.map(m => m.id)
    await updateMediaSort(albumId.value, mediaIds)
    ElMessage.success('排序已保存')
  } catch {
    loadMedia()
  }
}

const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedMedia.value.length} 个文件吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await batchDeleteMedia(selectedMedia.value)
    ElMessage.success('删除成功')
    selectedMedia.value = []
    loadMedia()
    loadAlbum()
  } catch {}
}

const createShare = async () => {
  if (shareForm.hasPassword && !shareForm.password) {
    ElMessage.warning('请设置访问密码')
    return
  }
  
  creatingShare.value = true
  try {
    const res = await createShareLink({
      album_id: albumId.value,
      password: shareForm.hasPassword ? shareForm.password : '',
      max_views: shareForm.maxViews,
      days: shareForm.expireType
    })
    createdShare.value = res.share
    ElMessage.success('分享链接已生成')
  } catch {
  } finally {
    creatingShare.value = false
  }
}

const copyShareLink = async () => {
  try {
    await navigator.clipboard.writeText(shareLinkUrl.value)
    ElMessage.success('链接已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

const prevImage = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--
    resetView()
  }
}

const nextImage = () => {
  if (currentIndex.value < mediaList.value.length - 1) {
    currentIndex.value++
    resetView()
  }
}

const zoomIn = () => {
  scale.value = Math.min(scale.value * 1.2, 5)
}

const zoomOut = () => {
  scale.value = Math.max(scale.value / 1.2, 0.2)
}

const handleZoom = (e) => {
  if (e.deltaY < 0) {
    zoomIn()
  } else {
    zoomOut()
  }
}

const resetView = () => {
  scale.value = 1
  position.x = 0
  position.y = 0
}

const startDrag = (e) => {
  isDragging.value = true
  dragStart.x = e.clientX - position.value.x
  dragStart.y = e.clientY - position.value.y
}

const onDrag = (e) => {
  if (!isDragging.value) return
  position.value.x = e.clientX - dragStart.x
  position.value.y = e.clientY - dragStart.y
}

const endDrag = () => {
  isDragging.value = false
}

const toggleSlideshow = () => {
  slideshow.value = !slideshow.value
  if (slideshow.value) {
    startSlideshow()
  } else {
    stopSlideshow()
  }
}

const startSlideshow = () => {
  slideshowTimer = setInterval(() => {
    if (currentIndex.value < mediaList.value.length - 1) {
      currentIndex.value++
      resetView()
    } else {
      currentIndex.value = 0
      resetView()
    }
  }, 3000)
}

const stopSlideshow = () => {
  if (slideshowTimer) {
    clearInterval(slideshowTimer)
    slideshowTimer = null
  }
  slideshow.value = false
}

const downloadCurrent = () => {
  if (currentMedia.value) {
    downloadMedia(currentMedia.value.id)
  }
}

watch(showViewer, (val) => {
  if (!val) {
    stopSlideshow()
  }
})

onMounted(() => {
  loadAlbum()
  loadMedia()
})

onUnmounted(() => {
  stopSlideshow()
})
</script>

<style scoped>
.album-detail-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: calc(100vh - 100px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fff;
  border-radius: 8px;
}

.album-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.album-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
}

.photo-count {
  color: #909399;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.loading-state,
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #909399;
}

.upload-area {
  width: 400px;
}

.upload-content {
  padding: 40px;
  text-align: center;
}

.upload-icon {
  font-size: 67px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.upload-tip {
  font-size: 12px;
  color: #c0c4cc;
  margin-top: 8px;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 12px;
  overflow-y: auto;
  padding: 4px;
}

.media-item {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.media-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.media-item.selected {
  outline: 3px solid #409EFF;
}

.media-item.drag-mode {
  cursor: move;
}

.media-checkbox {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 10;
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
  padding: 8px 10px;
  font-size: 12px;
  color: #606266;
  background: #fff;
}

.viewer-dialog :deep(.el-dialog) {
  max-width: 90vw;
  height: 90vh;
  margin: 5vh auto;
}

.viewer-dialog :deep(.el-dialog__body) {
  padding: 0;
  height: calc(90vh - 60px - 20px);
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
  overflow: hidden;
}

.viewer-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  transition: transform 0.1s ease-out;
  user-select: none;
}

.viewer-video {
  max-width: 100%;
  max-height: 100%;
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

.zoom-level {
  font-size: 12px;
  min-width: 50px;
  text-align: center;
}

.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.share-result {
  margin-top: 16px;
}

.share-link {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}
</style>
