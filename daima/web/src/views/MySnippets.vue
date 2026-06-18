<template>
  <div class="page-container">
    <div class="flex-between mb-20 flex-wrap gap-12">
      <h2 style="font-size:20px;">我的片段</h2>
      <div class="flex gap-12">
        <el-select v-model="visibility" size="large" style="width:140px;">
          <el-option label="全部" value="all" />
          <el-option label="公开" value="public" />
          <el-option label="私密" value="private" />
        </el-select>
        <el-button type="primary" size="large" :icon="Plus" @click="$router.push('/create')">新建片段</el-button>
      </div>
    </div>
    <el-row :gutter="20">
      <el-col v-for="item in snippets" :key="item.id" :xs="24" :sm="12" :md="8" :lg="8" :xl="6" class="mb-20">
        <SnippetCard :snippet="item" />
      </el-col>
    </el-row>
    <div v-if="snippets.length === 0 && loading === false" class="flex-center" style="padding:40px 0;">
      <el-empty description="暂无代码片段，点击右上角创建你的第一个片段吧！" />
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
import { snippetApi } from '../api'
import SnippetCard from '../components/SnippetCard.vue'

const snippets = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)
const visibility = ref('all')

const fetchData = async () => {
  loading.value = true
  try {
    const res = await snippetApi.my({
      page: page.value,
      page_size: pageSize.value,
      visibility: visibility.value
    })
    snippets.value = res.data
    total.value = res.total
  } finally {
    loading.value = false
  }
}

watch(visibility, () => {
  page.value = 1
  fetchData()
})

onMounted(fetchData)
</script>
