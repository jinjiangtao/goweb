<template>
  <div class="home">
    <header class="header">
      <h1 class="title">
        <el-icon :size="32" color="#409eff"><Edit /></el-icon>
        在线协同白板
      </h1>
      <p class="subtitle">实时协作 · 创意无限</p>
    </header>

    <div class="container">
      <div class="actions">
        <el-button type="primary" size="large" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          新建白板
        </el-button>
      </div>

      <div class="board-list">
        <div v-if="loading" class="loading">
          <el-icon :size="48" class="loading-icon"><Loading /></el-icon>
          <p>加载中...</p>
        </div>

        <el-empty v-else-if="boards.length === 0" description="暂无白板，点击上方按钮创建一个吧">
          <template #image>
            <el-icon :size="80" color="#909399"><Document /></el-icon>
          </template>
        </el-empty>

        <div v-else class="board-grid">
          <div 
            v-for="board in boards" 
            :key="board.id" 
            class="board-card"
            @click="openBoard(board.id)"
          >
            <div class="board-thumbnail" :style="{ background: board.background }">
              <img v-if="board.thumbnail" :src="board.thumbnail" alt="" />
              <el-icon v-else :size="64" color="#dcdfe6"><Picture /></el-icon>
            </div>
            <div class="board-info">
              <h3 class="board-name">{{ board.name }}</h3>
              <p class="board-date">{{ formatDate(board.updatedAt) }}</p>
            </div>
            <div class="board-actions">
              <el-button size="small" @click.stop="renameBoard(board)">
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click.stop="deleteBoard(board)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showCreateDialog" title="新建白板" width="400px">
      <el-form :model="createForm" label-width="80px">
        <el-form-item label="白板名称">
          <el-input v-model="createForm.name" placeholder="请输入白板名称" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="背景颜色">
          <el-color-picker v-model="createForm.background" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createBoard">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showRenameDialog" title="重命名白板" width="400px">
      <el-form :model="renameForm" label-width="80px">
        <el-form-item label="白板名称">
          <el-input v-model="renameForm.name" placeholder="请输入白板名称" maxlength="50" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRenameDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmRename">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { boardApi } from '@/utils/request'
import { formatTime } from '@/utils/common'

const router = useRouter()
const boards = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showRenameDialog = ref(false)
const currentBoard = ref(null)

const createForm = ref({
  name: '',
  background: '#ffffff'
})

const renameForm = ref({
  name: ''
})

async function loadBoards() {
  loading.value = true
  try {
    const res = await boardApi.list()
    boards.value = res.data || []
  } catch (e) {
    ElMessage.error('加载白板列表失败')
  } finally {
    loading.value = false
  }
}

async function createBoard() {
  if (!createForm.value.name.trim()) {
    ElMessage.warning('请输入白板名称')
    return
  }
  
  try {
    const res = await boardApi.create(createForm.value)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', background: '#ffffff' }
    router.push(`/board/${res.data.id}`)
  } catch (e) {
    ElMessage.error('创建失败')
  }
}

function openBoard(id) {
  router.push(`/board/${id}`)
}

function renameBoard(board) {
  currentBoard.value = board
  renameForm.value.name = board.name
  showRenameDialog.value = true
}

async function confirmRename() {
  if (!currentBoard.value || !renameForm.value.name.trim()) {
    ElMessage.warning('请输入白板名称')
    return
  }
  
  try {
    await boardApi.update(currentBoard.value.id, { name: renameForm.value.name })
    ElMessage.success('重命名成功')
    showRenameDialog.value = false
    loadBoards()
  } catch (e) {
    ElMessage.error('重命名失败')
  }
}

async function deleteBoard(board) {
  try {
    await ElMessageBox.confirm(`确定要删除白板"${board.name}"吗？`, '提示', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await boardApi.delete(board.id)
    ElMessage.success('删除成功')
    loadBoards()
  } catch {
    // 用户取消
  }
}

function formatDate(date) {
  return formatTime(date)
}

onMounted(() => {
  loadBoards()
})
</script>

<style lang="scss" scoped>
.home {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 0;
  margin: 0;
}

.header {
  text-align: center;
  padding: 60px 20px 40px;
  color: #fff;
  
  .title {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    font-size: 42px;
    margin: 0 0 10px;
    font-weight: 600;
  }
  
  .subtitle {
    font-size: 18px;
    opacity: 0.9;
    margin: 0;
  }
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px 60px;
  
  .actions {
    text-align: center;
    margin-bottom: 40px;
  }
}

.board-list {
  .loading {
    text-align: center;
    padding: 80px 20px;
    color: #fff;
    
    .loading-icon {
      animation: spin 1s linear infinite;
    }
  }
}

.board-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 20px;
}

.board-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
  }
  
  .board-thumbnail {
    height: 160px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-bottom: 1px solid #eee;
    
    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
  
  .board-info {
    padding: 12px 16px;
    
    .board-name {
      font-size: 16px;
      margin: 0 0 4px;
      color: #303133;
      font-weight: 500;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    
    .board-date {
      font-size: 12px;
      color: #909399;
      margin: 0;
    }
  }
  
  .board-actions {
    display: flex;
    gap: 8px;
    padding: 8px 16px 16px;
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
