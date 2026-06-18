<template>
  <div class="page-container">
    <div class="home-header mb-24 flex-between flex-wrap gap-16">
      <div class="filter-section flex flex-wrap gap-12 items-center">
        <el-select v-model="language" placeholder="选择语言" size="large" clearable style="width: 160px;">
          <el-option label="全部语言" value="all" />
          <el-option v-for="lang in languages" :key="lang" :label="lang" :value="lang" />
        </el-select>
        <el-radio-group v-model="sort" size="large">
          <el-radio-button value="latest">最新</el-radio-button>
          <el-radio-button value="likes">最热</el-radio-button>
          <el-radio-button value="favs">收藏</el-radio-button>
          <el-radio-button value="views">浏览</el-radio-button>
        </el-radio-group>
      </div>
      <el-button type="primary" size="large" :icon="Comparison" @click="$router.push('/diff')">
        代码对比
      </el-button>
    </div>

    <el-row :gutter="20">
      <el-col v-for="item in snippets" :key="item.id" :xs="24" :sm="12" :md="8" :lg="8" :xl="6" class="mb-20">
        <SnippetCard :snippet="item" />
      </el-col>
    </el-row>

    <div v-if="snippets.length === 0 && loading === false" class="empty-state flex-center">
      <el-empty description="暂无代码片段" />
    </div>

    <div v-if="total > pageSize" class="pagination flex-center mt-24">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        :page-sizes="[10, 20, 50]"
        background
        @current-change="fetchSnippets"
        @size-change="fetchSnippets"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { snippetApi } from '../api'
import SnippetCard from '../components/SnippetCard.vue'

const snippets = ref([])
const languages = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)
const language = ref('all')
const sort = ref('latest')

const fetchLanguages = async () => {
  try {
    const res = await snippetApi.getLanguages()
    languages.value = res.data
  } catch (e) {}
}

const fetchSnippets = async () => {
  loading.value = true
  try {
    const res = await snippetApi.list({
      page: page.value,
      page_size: pageSize.value,
      language: language.value,
      sort: sort.value
    })
    snippets.value = res.data
    total.value = res.total
  } finally {
    loading.value = false
  }
}

watch([language, sort], () => {
  page.value = 1
  fetchSnippets()
})

onMounted(() => {
  fetchLanguages()
  fetchSnippets()
})
</script>
