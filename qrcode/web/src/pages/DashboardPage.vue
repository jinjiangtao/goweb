<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { QrCode, ScanLine, Palette, Archive, ArrowRight, AlertTriangle } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { recordImageUrl } from '@/api'

const store = useAppStore()

onMounted(() => {
  store.refreshStats()
})

const statCards = [
  { key: 'total' as const, label: '累计生成', icon: QrCode, color: 'text-lime' },
  { key: 'today' as const, label: '今日生成', icon: QrCode, color: 'text-lime' },
  { key: 'templates' as const, label: '模板数', icon: Palette, color: 'text-ink-secondary' },
  { key: 'invalid' as const, label: '失效标记', icon: AlertTriangle, color: 'text-amber' },
]

const shortcuts = [
  { to: '/generate', icon: QrCode, label: '生成二维码', desc: '单条/批量生成' },
  { to: '/parse', icon: ScanLine, label: '解析识别', desc: '文本/图片解析' },
  { to: '/templates', icon: Palette, label: '模板管理', desc: '保存复用样式' },
  { to: '/archive', icon: Archive, label: '归档导出', desc: '记录/打包下载' },
]
</script>

<template>
  <div class="stagger space-y-6">
    <section class="grid grid-cols-4 gap-4">
      <div
        v-for="card in statCards"
        :key="card.key"
        class="panel flex items-center gap-4 px-5 py-4"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-lg"
          :class="card.key === 'invalid' ? 'bg-amber/10' : card.key === 'templates' ? 'bg-base-hover' : 'bg-lime/10'"
        >
          <component :is="card.icon" class="h-5 w-5" :class="card.color" />
        </div>
        <div>
          <p class="font-mono text-2xl font-semibold tracking-tight" :class="card.color">
            {{ store.stats?.[card.key] ?? '—' }}
          </p>
          <p class="text-xs text-ink-muted">{{ card.label }}</p>
        </div>
      </div>
    </section>

    <section class="grid grid-cols-4 gap-4">
      <RouterLink
        v-for="s in shortcuts"
        :key="s.to"
        :to="s.to"
        class="panel group flex items-center gap-4 px-5 py-4 transition hover:shadow-glow"
      >
        <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-lime/10">
          <component :is="s.icon" class="h-4.5 w-4.5 text-lime" />
        </div>
        <div class="flex-1">
          <p class="text-sm font-medium text-ink">{{ s.label }}</p>
          <p class="text-xs text-ink-muted">{{ s.desc }}</p>
        </div>
        <ArrowRight class="h-4 w-4 text-ink-faint transition group-hover:text-lime" />
      </RouterLink>
    </section>

    <section>
      <h2 class="mb-3 font-display text-sm font-semibold text-ink-secondary">近期生成</h2>
      <div v-if="store.stats?.recent?.length" class="grid grid-cols-4 gap-4">
        <div
          v-for="r in store.stats.recent"
          :key="r.id"
          class="panel-elevated overflow-hidden"
        >
          <div class="aspect-square bg-checker bg-checker p-2">
            <img
              :src="recordImageUrl(r.id)"
              :alt="r.content"
              class="h-full w-full object-contain"
            />
          </div>
          <div class="px-3 py-2.5">
            <p class="truncate text-xs font-medium text-ink">{{ r.content }}</p>
            <div class="mt-1 flex items-center gap-2">
              <span class="chip">
                {{ r.contentType === 'url' ? '链接' : '文本' }}
              </span>
              <span
                class="chip"
                :class="r.status === 'invalid' ? 'border-amber/30 text-amber' : ''"
              >
                {{ r.status === 'active' ? '有效' : '失效' }}
              </span>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="panel flex items-center justify-center py-16 text-sm text-ink-faint">
        暂无生成记录
      </div>
    </section>
  </div>
</template>
