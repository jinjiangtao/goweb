<template>
  <Layout>
    <div class="albums-page">
      <div class="page-header">
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          创建相册
        </el-button>
        <div class="header-actions">
          <el-button @click="toggleSortMode">
            <el-icon><Sort /></el-icon>
            {{ sortMode ? '完成排序' : '排序' }}
          </el-button>
        </div>
      </div>

      <div v-if="albums.length === 0" class="empty-state">
        <el-empty description="暂无相册，创建您的第一个相册吧" />
        <el-button type="primary" size="large" @click="showCreateDialog = true">
          创建相册
        </el-button>
      </div>

      <draggable
        v-else
        v-model="albums"
        item-key="id"
        handle=".album-card"
        :disabled="!sortMode"
        class="album-grid"
        :animation="200"
        @end="onSortEnd"
      >
        <template #item="{ element }">
          <div class="album-item">
            <div class="album-card" @click="goToAlbum(element.id)">
              <div class="album-cover">
                <img v-if="element.cover" :src="getCoverUrl(element.cover)" />
                <div v-else class="cover-placeholder">
                  <el-icon :size="64" color="#c0c4cc"><Picture /></el-icon>
                </div>
                <div class="album-badge" v-if="element.is_private">
                  <el-icon :size="14"><Lock /></el-icon>
                  私密
                </div>
                <div v-if="sortMode" class="sort-handle">
                  <el-icon :size="24"><Rank /></el-icon>
                </div>
              </div>
              <div class="album-info">
                <div class="album-name text-ellipsis">{{ element.name }}</div>
                <div class="album-desc text-ellipsis">{{ element.description || '暂无描述' }}</div>
                <div class="album-footer">
                  <span class="photo-count">
                    <el-icon><Picture /></el-icon>
                    {{ element.photo_count }} 张
                  </span>
                  <span class="album-date">
                    {{ formatDate(element.created_at) }}
                  </span>
                </div>
              </div>
            </div>
            <div class="album-actions" v-if="!sortMode">
              <el-button size="small" @click.stop="editAlbum(element)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button size="small" type="danger" @click.stop="deleteAlbum(element)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>
          </div>
        </template>
      </draggable>
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
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEditDialog" title="编辑相册" width="500px">
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
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate">保存</el-button>
      </template>
    </el-dialog>
  </Layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import draggable from 'vuedraggable'
import Layout from '@/components/Layout.vue'
import { getAlbums, createAlbum, updateAlbum, deleteAlbum as deleteAlbumApi, updateAlbumSort } from '@/utils/api'

const router = useRouter()

const albums = ref([])
const sortMode = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const editingAlbum = ref(null)

const albumForm = reactive({
  name: '',
  description: '',
  is_private: false,
  password: ''
})

const getCoverUrl = (path) => {
  if (!path) return ''
  const filename = path.split('/').pop()
  return `/uploads/thumbnails/${filename}`
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const goToAlbum = (id) => {
  if (!sortMode.value) {
    router.push(`/album/${id}`)
  }
}

const loadAlbums = async () => {
  try {
    const res = await getAlbums()
    albums.value = res.albums
  } catch {}
}

const handleCreate = async () => {
  if (!albumForm.name.trim()) {
    ElMessage.warning('请输入相册名称')
    return
  }
  
  try {
    await createAlbum({
      name: albumForm.name,
      description: albumForm.description,
      is_private: albumForm.is_private,
      password: albumForm.password
    })
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    resetForm()
    loadAlbums()
  } catch {}
}

const editAlbum = (album) => {
  editingAlbum.value = album
  albumForm.name = album.name
  albumForm.description = album.description || ''
  albumForm.is_private = album.is_private
  albumForm.password = album.password || ''
  showEditDialog.value = true
}

const handleUpdate = async () => {
  if (!albumForm.name.trim()) {
    ElMessage.warning('请输入相册名称')
    return
  }
  
  try {
    await updateAlbum(editingAlbum.value.id, {
      name: albumForm.name,
      description: albumForm.description,
      is_private: albumForm.is_private,
      password: albumForm.password
    })
    ElMessage.success('保存成功')
    showEditDialog.value = false
    resetForm()
    loadAlbums()
  } catch {}
}

const deleteAlbum = async (album) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除相册"${album.name}"吗？相册内的所有照片也会被删除。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteAlbumApi(album.id)
    ElMessage.success('删除成功')
    loadAlbums()
  } catch {
    if (action !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const toggleSortMode = () => {
  sortMode.value = !sortMode.value
}

const onSortEnd = async () => {
  try {
    const albumIds = albums.value.map(a => a.id)
    await updateAlbumSort(albumIds)
    ElMessage.success('排序已保存')
  } catch {
    loadAlbums()
  }
}

const resetForm = () => {
  albumForm.name = ''
  albumForm.description = ''
  albumForm.is_private = false
  albumForm.password = ''
  editingAlbum.value = null
}

onMounted(() => {
  loadAlbums()
})
</script>

<style scoped>
.albums-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.empty-state {
  padding: 80px 0;
  text-align: center;
}

.album-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
}

.album-item {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: box-shadow 0.3s, transform 0.2s;
}

.album-item:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
}

.sort-mode .album-item {
  cursor: move;
}

.album-card {
  cursor: pointer;
}

.album-cover {
  position: relative;
  width: 100%;
  padding-top: 75%;
  background-color: #f0f2f5;
  overflow: hidden;
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
  top: 10px;
  right: 10px;
  padding: 4px 10px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  font-size: 12px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.sort-handle {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60px;
  height: 60px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.album-info {
  padding: 16px;
}

.album-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 6px;
}

.album-desc {
  font-size: 13px;
  color: #909399;
  margin-bottom: 12px;
  height: 18px;
}

.album-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #c0c4cc;
}

.photo-count {
  display: flex;
  align-items: center;
  gap: 4px;
}

.album-actions {
  display: flex;
  gap: 0;
  border-top: 1px solid #f0f2f5;
}

.album-actions .el-button {
  flex: 1;
  border: none;
  border-radius: 0;
}

.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
