<template>
  <el-container class="layout-container">
    <el-aside width="240px" class="sidebar">
      <div class="logo">
        <el-icon :size="24"><Document /></el-icon>
        <span>个人笔记</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="menu"
        @select="handleMenuSelect"
      >
        <el-menu-item index="/notes">
          <el-icon><Document /></el-icon>
          <span>全部笔记</span>
          <el-badge v-if="stats && stats.total > 0" :value="stats.total" class="badge" />
        </el-menu-item>
        <el-menu-item index="/tags">
          <el-icon><PriceTag /></el-icon>
          <span>标签管理</span>
        </el-menu-item>
        <el-menu-item index="/archived">
          <el-icon><FolderOpened /></el-icon>
          <span>归档</span>
          <el-badge v-if="stats && stats.archived > 0" :value="stats.archived" class="badge" />
        </el-menu-item>
        <el-menu-item index="/trash">
          <el-icon><Delete /></el-icon>
          <span>回收站</span>
          <el-badge v-if="stats && stats.deleted > 0" :value="stats.deleted" class="badge" type="danger" />
        </el-menu-item>
      </el-menu>

      <div class="sidebar-section">
        <div class="sidebar-title flex items-center justify-between">
          <span>标签筛选</span>
          <el-button link type="primary" size="small" @click="goToTags">管理</el-button>
        </div>
        <div class="tag-filter-list">
          <el-tag
            v-for="tag in flatTags"
            :key="tag.id"
            :color="tag.color"
            effect="light"
            size="small"
            class="tag-filter-item"
            :class="{ active: selectedTagIds.includes(tag.id) }"
            @click="toggleTagFilter(tag.id)"
            closable
            @close="removeTagFilter(tag.id)"
          >
            {{ tag.name }}
          </el-tag>
          <span v-if="flatTags.length === 0" class="text-muted text-small">暂无标签</span>
        </div>
      </div>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left flex items-center gap-md">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索笔记标题或内容..."
            style="width: 320px"
            clearable
            @input="handleSearch"
            @clear="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            size="default"
            @change="handleDateChange"
          />
          <el-button
            v-if="selectedTagIds.length > 0 || searchKeyword || dateRange"
            @click="clearFilters"
          >
            清除筛选
          </el-button>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="createNote">
            <el-icon><Plus /></el-icon>
            新建笔记
          </el-button>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view
          :search-keyword="searchKeyword"
          :selected-tag-ids="selectedTagIds"
          :date-range="dateRange"
          @filters-changed="handleFiltersChanged"
        />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { noteApi, tagApi } from '@/api'

const router = useRouter()
const route = useRoute()

const stats = ref(null)
const allTags = ref([])
const searchKeyword = ref('')
const selectedTagIds = ref([])
const dateRange = ref(null)

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/notes')) return '/notes'
  if (path.startsWith('/editor')) return '/notes'
  return path.split('/')[1] ? '/' + path.split('/')[1] : '/notes'
})

const flatTags = computed(() => {
  const result = []
  const flatten = (tags) => {
    tags.forEach(tag => {
      result.push(tag)
      if (tag.children && tag.children.length) {
        flatten(tag.children)
      }
    })
  }
  flatten(allTags.value)
  return result
})

const loadStats = async () => {
  stats.value = await noteApi.stats()
}

const loadTags = async () => {
  allTags.value = await tagApi.list()
}

const handleMenuSelect = (index) => {
  if (index !== route.path) {
    router.push(index)
  }
}

const createNote = () => {
  router.push('/editor')
}

const goToTags = () => {
  router.push('/tags')
}

const toggleTagFilter = (tagId) => {
  const idx = selectedTagIds.value.indexOf(tagId)
  if (idx > -1) {
    selectedTagIds.value.splice(idx, 1)
  } else {
    selectedTagIds.value.push(tagId)
  }
}

const removeTagFilter = (tagId) => {
  const idx = selectedTagIds.value.indexOf(tagId)
  if (idx > -1) {
    selectedTagIds.value.splice(idx, 1)
  }
}

const handleSearch = () => {
}

const handleDateChange = () => {
}

const clearFilters = () => {
  searchKeyword.value = ''
  selectedTagIds.value = []
  dateRange.value = null
}

const handleFiltersChanged = (filters) => {
  if (filters.keyword !== undefined) searchKeyword.value = filters.keyword
  if (filters.tagIds !== undefined) selectedTagIds.value = filters.tagIds
}

watch(route, () => {
  if (route.path === '/notes' || route.path === '/archived' || route.path === '/trash') {
    loadStats()
    loadTags()
  }
})

onMounted(() => {
  loadStats()
  loadTags()
})
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background: #fff;
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;

  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    font-size: 18px;
    font-weight: 600;
    border-bottom: 1px solid #e4e7ed;
    color: #409EFF;
  }

  .menu {
    border-right: none;
    padding: 8px;
  }

  .badge {
    margin-left: auto;
  }

  .sidebar-section {
    padding: 16px;
    border-top: 1px solid #f0f0f0;
    flex: 1;
    overflow: auto;
  }

  .sidebar-title {
    margin-bottom: 12px;
    font-weight: 500;
    font-size: 14px;
  }

  .tag-filter-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .tag-filter-item {
    cursor: pointer;
    transition: all 0.2s;

    &.active {
      opacity: 1;
      font-weight: 500;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    &:not(.active) {
      opacity: 0.7;
    }

    &:hover {
      opacity: 1;
    }
  }
}

.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
}

.main-content {
  padding: 24px;
  background: #f5f7fa;
  overflow: auto;
}
</style>
