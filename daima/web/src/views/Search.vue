<template>
  <div class="page-container">
    <div class="flex-between mb-20 flex-wrap gap-12">
      <h2 style="font-size:20px;">搜索: {{ keyword }}</h2>
      <el-select v-model="language" placeholder="语言过滤" size="large" clearable style="width:160px;">
        <el-option label="全部语言" value="all" />
        <el-option v-for="lang in languages" :key="lang" :label="lang" :value="lang" />
      </el-select>
    </div>
    <el-row :gutter="20">
      <el-col v-for="item in snippets" :key="item.id" :xs="24" :sm="12" :md="8" :lg="8" :xl="6" class="mb-20">
        <SnippetCard :snippet="item" />
      </el-col>
    </el-row>
    <div v-if="snippets.length === 0 && loading === false" class="flex-center" style="padding:40px 0;">
      <el-empty description="没有找到相关的代码片段" />
    </div>
    <div v-if="total > pageSize" class="pagination flex-center mt-24">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        background
        @current-change="fetchData"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { snippetApi } from '../api'
import SnippetCard from '../components/SnippetCard.vue'

const route = useRoute()
const snippets = ref([])
const languages = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)
const keyword = ref(route.query.q || '')
const language = ref('all')

const fetchLanguages = async () => {
  try {
    const res = await snippetApi.getLanguages()
    languages.value = res.data
  } catch (e) {}
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await snippetApi.search({
      q: keyword.value,
      page: page.value,
      page_size: pageSize.value,
      language: language.value
    })
    snippets.value = res.data
    total.value = res.total
  } finally {
    loading.value = false
  }
}

watch(() => route.query.q, (q) => {
  keyword.value = q || ''
  page.value = 1
  fetchData()
})

watch(language, () => {
  page.value = 1
  fetchData()
})

onMounted(() => {
  fetchLanguages()
  fetchData()
})
</script>
