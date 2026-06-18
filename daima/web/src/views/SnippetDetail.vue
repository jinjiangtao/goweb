<template>
  <div class="page-container">
    <div v-if="loading" class="flex-center" style="padding: 60px 0;">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>
    <div v-else-if="!snippet" class="flex-center" style="padding: 60px 0;">
      <el-empty description="片段不存在" />
    </div>
    <template v-else>
      <div class="detail-header mb-20 card">
        <div class="flex-between flex-wrap gap-12 mb-16">
          <h1 class="title">{{ snippet.title }}</h1>
          <div class="actions flex gap-12">
            <el-button size="large" :icon="CopyDocument" @click="copyCode">复制代码</el-button>
            <el-button
              v-if="userStore.isLoggedIn"
              size="large"
              :type="liked ? 'danger' : 'default'"
              :icon="liked ? GoodFilled : Good"
              @click="toggleLike"
            >
              点赞 {{ snippet.likes_count }}
            </el-button>
            <el-button
              v-if="userStore.isLoggedIn"
              size="large"
              :type="favorited ? 'warning' : 'default'"
              :icon="favorited ? StarFilled : Star"
              @click="toggleFavorite"
            >
              收藏 {{ snippet.favs_count }}
            </el-button>
            <el-button
              v-if="userStore.isLoggedIn"
              size="large"
              type="success"
              :icon="CopyDocument"
              @click="handleFork"
            >
              Fork
            </el-button>
            <el-button
              v-if="isOwner && userStore.isLoggedIn"
              size="large"
              type="primary"
              :icon="Edit"
              @click="$router.push(`/edit/${snippet.id}`)"
            >
              编辑
            </el-button>
            <el-popconfirm title="确定删除此片段？" @confirm="handleDelete">
              <template #reference>
                <el-button
                  v-if="isOwner && userStore.isLoggedIn"
                  size="large"
                  type="danger"
                  :icon="Delete"
                >
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>
        <div class="flex flex-wrap gap-16 mb-16 items-center">
          <span :class="['language-tag', `lang-${snippet.language}`]" style="padding:4px 12px;font-size:13px;">
            {{ snippet.language }}
          </span>
          <div class="flex gap-8 items-center">
            <el-avatar :size="28" :src="snippet.user?.avatar">
              {{ snippet.user?.username?.charAt(0) }}
            </el-avatar>
            <span class="text-sm">{{ snippet.user?.username }}</span>
          </div>
          <span class="stat-item"><el-icon><View /></el-icon>{{ snippet.views_count }} 浏览</span>
          <span class="text-sm text-muted">创建于 {{ formatDate(snippet.created_at) }}</span>
        </div>
        <p v-if="snippet.description" class="desc mb-12">{{ snippet.description }}</p>
        <div v-if="tags.length" class="tags flex flex-wrap gap-8">
          <el-tag v-for="tag in tags" :key="tag" effect="plain" class="tag-btn">{{ tag }}</el-tag>
        </div>
      </div>

      <div class="card mb-20">
        <div class="flex-between mb-12">
          <div class="flex gap-12 items-center">
            <h3 style="font-size:16px;">代码</h3>
            <el-button size="small" :icon="MagicStick" @click="formatCode">格式化</el-button>
            <span class="text-sm text-muted">{{ snippet.language }}</span>
          </div>
        </div>
        <div style="height: 500px;">
          <CodeEditor ref="editorRef" v-model="code" :language="snippet.language" :read-only="true" />
        </div>
      </div>

      <div class="card">
        <h3 class="mb-16" style="font-size:16px;">评论 ({{ comments.length }})</h3>
        <div v-if="userStore.isLoggedIn" class="mb-20">
          <el-input
            v-model="commentContent"
            type="textarea"
            :rows="3"
            placeholder="发表你的评论..."
            maxlength="1000"
            show-word-limit
          />
          <div class="flex justify-end mt-12">
            <el-button type="primary" :loading="commentLoading" @click="submitComment">
              发布评论
            </el-button>
          </div>
        </div>
        <div v-else class="mb-20 text-muted text-center" style="padding:20px;">
          请先登录后发表评论
        </div>
        <div class="comments">
          <div v-if="comments.length === 0" class="text-muted text-center" style="padding:20px;">
            暂无评论
          </div>
          <div v-for="c in comments" :key="c.id" class="comment-item mb-16" style="padding:16px;border-bottom:1px solid #f0f2f5;">
            <div class="flex-between mb-8">
              <div class="flex gap-8 items-center">
                <el-avatar :size="28" :src="c.user?.avatar">
                  {{ c.user?.username?.charAt(0) }}
                </el-avatar>
                <span style="font-weight:500;">{{ c.user?.username }}</span>
                <span class="text-xs text-muted">{{ formatDate(c.created_at) }}</span>
              </div>
              <el-button
                v-if="c.user?.id === userStore.user?.id"
                link
                type="danger"
                size="small"
                @click="deleteComment(c.id)"
              >
                删除
              </el-button>
            </div>
            <div class="comment-content" style="padding-left:36px;white-space:pre-wrap;">{{ c.content }}</div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { snippetApi, commentApi } from '../api'
import { useUserStore } from '../store/user'
import CodeEditor from '../components/CodeEditor.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const snippet = ref(null)
const comments = ref([])
const loading = ref(false)
const liked = ref(false)
const favorited = ref(false)
const code = ref('')
const commentContent = ref('')
const commentLoading = ref(false)
const editorRef = ref(null)

const tags = computed(() => (snippet.value?.tags || '').split(',').map(t => t.trim()).filter(Boolean))
const isOwner = computed(() => snippet.value?.user_id === userStore.user?.id)

const formatDate = (d) => dayjs(d).format('YYYY-MM-DD HH:mm')

const fetchData = async () => {
  loading.value = true
  try {
    const res = await snippetApi.get(route.params.id)
    snippet.value = res.data.snippet
    code.value = snippet.value.content
    liked.value = !!res.data.liked
    favorited.value = !!res.data.favorited
    const cRes = await commentApi.list(route.params.id)
    comments.value = cRes.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(code.value)
    ElMessage.success('代码已复制到剪贴板')
  } catch (e) {
    ElMessage.error('复制失败')
  }
}

const formatCode = () => editorRef.value?.formatCode()

const toggleLike = async () => {
  try {
    const res = await snippetApi.like(route.params.id)
    liked.value = res.data.liked
    snippet.value.likes_count = res.data.likes_count
  } catch (e) {}
}

const toggleFavorite = async () => {
  try {
    const res = await snippetApi.favorite(route.params.id)
    favorited.value = res.data.favorited
    snippet.value.favs_count = res.data.favs_count
  } catch (e) {}
}

const handleFork = async () => {
  try {
    const res = await snippetApi.fork(route.params.id)
    ElMessage.success('Fork成功')
    router.push(`/edit/${res.data.id}`)
  } catch (e) {}
}

const handleDelete = async () => {
  try {
    await snippetApi.delete(route.params.id)
    ElMessage.success('删除成功')
    router.push('/')
  } catch (e) {}
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }
  commentLoading.value = true
  try {
    const res = await commentApi.create(route.params.id, { content: commentContent.value })
    comments.value.unshift(res.data)
    snippet.value.comments_count++
    commentContent.value = ''
    ElMessage.success('评论成功')
  } finally {
    commentLoading.value = false
  }
}

const deleteComment = async (id) => {
  try {
    await commentApi.delete(id)
    comments.value = comments.value.filter(c => c.id !== id)
    snippet.value.comments_count--
    ElMessage.success('删除成功')
  } catch (e) {}
}

onMounted(fetchData)
</script>

<style scoped>
.title {
  font-size: 22px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}
.desc {
  color: #6b7280;
  font-size: 14px;
  line-height: 1.6;
}
.lang-javascript { background: #f7df1e20; color: #b59000; }
.lang-typescript { background: #3178c620; color: #3178c6; }
.lang-python { background: #3776ab20; color: #3776ab; }
.lang-go { background: #00add820; color: #00add8; }
.lang-java { background: #ed8b0020; color: #ed8b00; }
.lang-html { background: #e44d2620; color: #e44d26; }
.lang-css { background: #1572b620; color: #1572b6; }
:deep(.language-tag) { text-transform: capitalize; border-radius: 6px; font-weight: 500; }
</style>
