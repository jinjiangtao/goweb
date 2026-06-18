<template>
  <div class="snippet-card card" @click="$router.push(`/snippet/${snippet.id}`)">
    <div class="snippet-header mb-16">
      <div class="flex-between" style="width:100%">
        <h3 class="snippet-title line-clamp-2">{{ snippet.title }}</h3>
        <span :class="['language-tag', `lang-${snippet.language}`]">{{ snippet.language }}</span>
      </div>
    </div>
    <p v-if="snippet.description" class="snippet-desc line-clamp-2 mb-16 text-muted text-sm">
      {{ snippet.description }}
    </p>
    <div class="code-preview mb-16">
      <pre><code>{{ previewContent }}</code></pre>
    </div>
    <div v-if="snippet.tags" class="tags mb-16 flex flex-wrap gap-8">
      <el-tag v-for="tag in tagList" :key="tag" size="small" effect="plain" class="tag-btn">
        {{ tag }}
      </el-tag>
    </div>
    <div class="snippet-footer flex-between">
      <div class="snippet-user flex gap-8">
        <el-avatar :size="24" :src="snippet.user?.avatar">
          {{ snippet.user?.username?.charAt(0) }}
        </el-avatar>
        <span class="text-sm">{{ snippet.user?.username }}</span>
      </div>
      <div class="stats flex gap-16">
        <span class="stat-item"><el-icon><View /></el-icon>{{ snippet.views_count }}</span>
        <span class="stat-item"><el-icon><Star /></el-icon>{{ snippet.favs_count }}</span>
        <span class="stat-item"><el-icon><GoodFilled /></el-icon>{{ snippet.likes_count }}</span>
        <span class="stat-item"><el-icon><ChatDotRound /></el-icon>{{ snippet.comments_count }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  snippet: { type: Object, required: true }
})

const previewContent = computed(() => {
  const lines = (props.snippet.content || '').split('\n').slice(0, 6)
  return lines.join('\n') + (lines.length >= 6 ? '\n...' : '')
})

const tagList = computed(() => {
  return (props.snippet.tags || '').split(',').map(t => t.trim()).filter(Boolean)
})
</script>

<style scoped>
.snippet-card {
  cursor: pointer;
  height: 100%;
  display: flex;
  flex-direction: column;
}
.snippet-title {
  font-size: 16px;
  color: #1f2937;
  font-weight: 600;
  flex: 1;
  margin-right: 12px;
}
.code-preview {
  background: #1e1e1e;
  border-radius: 6px;
  padding: 12px 16px;
  flex: 1;
  overflow: hidden;
  max-height: 160px;
}
.code-preview pre {
  margin: 0;
  padding: 0;
}
.code-preview code {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #d4d4d4;
  white-space: pre-wrap;
  line-height: 1.6;
}
.snippet-footer {
  border-top: 1px solid #f0f2f5;
  padding-top: 12px;
  margin-top: auto;
}
.lang-javascript { background: #f7df1e20; color: #b59000; }
.lang-typescript { background: #3178c620; color: #3178c6; }
.lang-python { background: #3776ab20; color: #3776ab; }
.lang-go { background: #00add820; color: #00add8; }
.lang-java { background: #ed8b0020; color: #ed8b00; }
.lang-c { background: #a8b9cc20; color: #00599c; }
.lang-cpp { background: #00599c20; color: #00599c; }
.lang-html { background: #e44d2620; color: #e44d26; }
.lang-css { background: #1572b620; color: #1572b6; }
.lang-json { background: #88888820; color: #888; }
.lang-rust { background: #dea58420; color: #e77219; }
.lang-ruby { background: #cc342d20; color: #cc342d; }
.lang-php { background: #8892bf20; color: #777bb4; }
.lang-sql { background: #f2911120; color: #f29111; }
.lang-shell, .lang-bash { background: #4eaa2520; color: #4eaa25; }
.lang-yaml { background: #cb171e20; color: #cb171e; }
.lang-markdown { background: #083fa120; color: #083fa1; }
:deep(.language-tag) { text-transform: capitalize; }
</style>
