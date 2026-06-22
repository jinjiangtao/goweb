<template>
  <Layout>
    <div class="home-page">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon album">
                <el-icon :size="32"><FolderOpened /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.album_count || 0 }}</div>
                <div class="stat-label">相册数量</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon photo">
                <el-icon :size="32"><Picture /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stats.photo_count || 0 }}</div>
                <div class="stat-label">照片数量</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon share">
                <el-icon :size="32"><Share /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ shareCount }}</div>
                <div class="stat-label">分享链接</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-card class="section-card">
            <template #header>
              <div class="card-header">
                <span>最近相册</span>
                <el-button type="primary" link @click="$router.push('/albums')">
                  查看全部
                  <el-icon><ArrowRight /></el-icon>
                </el-button>
              </div>
            </template>
            
            <div v-if="recentAlbums.length === 0" class="empty-state">
              <el-empty description="暂无相册，快去创建第一个相册吧" />
              <el-button type="primary" @click="showCreateDialog = true">
                创建相册
              </el-button>
            </div>
            
            <div v-else class="album-grid">
              <div
                v-for="album in recentAlbums"
                :key="album.id"
                class="album-item"
                @click="goToAlbum(album.id)"
              >
                <div class="album-cover">
                  <img v-if="album.cover" :src="getCoverUrl(album.cover)" />
                  <div v-else class="cover-placeholder">
                    <el-icon :size="48" color="#c0c4cc"><Picture /></el-icon>
                  </div>
                  <div class="album-badge" v-if="album.is_private">
                    <el-icon :size="14"><Lock /></el-icon>
                    私密
                  </div>
                </div>
                <div class="album-info">
                  <div class="album-name text-ellipsis">{{ album.name }}</div>
                  <div class="album-meta">
                    {{ album.photo_count }} 张照片
                  </div>
                </div>
              </div>
            </div>
          </el-card>

      <el-card class="section-card">
        <template #header>
          <div class="card-header">
            <span>最近照片</span>
          </div>
        </template>
        
        <div v-if="recentPhotos.length === 0" class="empty-state">
          <el-empty description="暂无照片" />
        </div>
        
        <div v-else class="photo-grid">
          <div
            v-for="photo in recentPhotos"
            :key="photo.id"
            class="photo-item"
          >
            <img :src="getThumbUrl(photo.thumb_path)" class="photo-thumb" />
          </div>
        </div>
      </el-card>
    </div>
    
    <el-dialog v-model="showCreateDialog" title="创建相册" width="500px">
      <el-form :model="albumForm" label-width="80px">
        <el-form-item label="相册名称">
          <el-input v-model="albumForm.name" placeholder="请输入相册名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="albumForm.description" type="textarea" :rows="3" placeholder="相册描述（可选）" />
        </el-form-item>
        <el-form-item label="私密相册">
          <el-switch v-model="albumForm.is_private" />
        </el-form-item>
        <el-form-item label="访问密码" v-if="albumForm.is_private">
          <el-input v-model="albumForm.password" type="password" placeholder="设置访问密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createAlbum">确定</el-button>
      </template>
    </el-dialog>
  </Layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Layout from '@/components/Layout.vue'
import { useUserStore } from '@/store/user'
import { getAlbums, createAlbum as createAlbumApi, getMediaList, getShareLinks } from '@/utils/api'

const router = useRouter()
const userStore = useUserStore()

const stats = reactive({
  album_count: 0,
  photo_count: 0
})

const recentAlbums = ref([])
const recentPhotos = ref([])
const shareCount = ref(0)
const showCreateDialog = ref(false)

const albumForm = reactive({
  name: '',
  description: '',
  is_private: false,
  password: ''
})

const getCoverUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `/uploads/thumbnails/${path.split('/').pop()}`
}

const getThumbUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  const filename = path.split('/').pop()
  return `/uploads/thumbnails/${filename}`
}

const goToAlbum = (id) => {
  router.push(`/album/${id}`)
}

const loadStats = async () => {
  try {
    const res = await userStore.fetchStats()
    Object.assign(stats, res)
  } catch {}
}

const loadAlbums = async () => {
  try {
    const res = await getAlbums()
    recentAlbums.value = res.albums.slice(0, 6)
    
    if (res.albums.length > 0) {
      loadRecentPhotos(res.albums[0].id)
    }
  } catch {}
}

const loadRecentPhotos = async (albumId) => {
  try {
    const res = await getMediaList(albumId)
    recentPhotos.value = res.media.slice(0, 12)
  } catch {}
}

const loadShares = async () => {
  try {
    const res = await getShareLinks()
    shareCount.value = res.shares.length
  } catch {}
}

const createAlbum = async () => {
  if (!albumForm.name) {
    ElMessage.warning('请输入相册名称')
    return
  }
  
  try {
    await createAlbumApi(albumForm)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    albumForm.name = ''
    albumForm.description = ''
    albumForm.is_private = false
    albumForm.password = ''
    loadAlbums()
    loadStats()
  } catch {}
}

onMounted(() => {
  loadStats()
  loadAlbums()
  loadShares()
})
</script>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stat-card {
  border-radius: 8px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-icon.album {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.photo {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.share {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.section-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}

.album-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
}

.album-item {
  cursor: pointer;
  transition: transform 0.2s;
}

.album-item:hover {
  transform: translateY(-2px);
}

.album-cover {
  position: relative;
  width: 100%;
  padding-top: 75%;
  border-radius: 8px;
  overflow: hidden;
  background-color: #f0f2f5;
}

.album-cover img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}

.album-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 2px 8px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  font-size: 12px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.album-info {
  margin-top: 8px;
}

.album-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.album-meta {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.photo-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 8px;
}

.photo-item {
  aspect-ratio: 1;
  border-radius: 4px;
  overflow: hidden;
  background: #f0f2f5;
}

.photo-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
