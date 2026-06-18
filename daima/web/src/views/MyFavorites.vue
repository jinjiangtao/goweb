<template>
  <div class="page-container">
    <h2 class="mb-20" style="font-size:20px;">我的收藏</h2>
    <el-row :gutter="20">
      <el-col v-for="item in snippets" :key="item.id" :xs="24" :sm="12" :md="8" :lg="8" :xl="6" class="mb-20">
        <SnippetCard :snippet="item" />
      </el-col>
    </el-row>
    <div v-if="snippets.length === 0 && loading === false" class="flex-center" style="padding:40px 0;">
      <el-empty description="还没有收藏任何片段，去首页发现有趣的代码吧！" />
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
import { ref, onMounted } from 'vue'
import { snippetApi } from '../api'
import SnippetCard from '../components/SnippetCard.vue'

const snippets = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await snippetApi.myFavorites({
      page: page.value,
      page_size: pageSize.value
    })
    snippets.value = res.data
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
