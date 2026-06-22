<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import {
  LayoutDashboard,
  QrCode,
  ScanLine,
  Palette,
  Archive,
  Zap,
} from 'lucide-vue-next'

const route = useRoute()
const pageTitle = computed(() => (route.meta.title as string) || 'QRForge')

const navItems = [
  { to: '/', icon: LayoutDashboard, label: '概览' },
  { to: '/generate', icon: QrCode, label: '生成' },
  { to: '/parse', icon: ScanLine, label: '解析' },
  { to: '/templates', icon: Palette, label: '模板' },
  { to: '/archive', icon: Archive, label: '归档' },
]
</script>

<template>
  <div class="flex min-h-screen">
    <aside class="fixed inset-y-0 left-0 z-30 flex w-[220px] flex-col border-r border-base-border bg-base-panel">
      <div class="flex h-14 items-center gap-2.5 px-5">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-lime/15">
          <Zap class="h-4.5 w-4.5 text-lime" />
        </div>
        <span class="font-display text-lg font-bold tracking-tight text-ink">QRForge</span>
      </div>

      <nav class="mt-2 flex flex-1 flex-col gap-0.5 px-3">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
          :class="[
            route.path === item.to || (item.to !== '/' && route.path.startsWith(item.to))
              ? 'bg-lime/10 text-lime'
              : 'text-ink-secondary hover:bg-base-hover hover:text-ink',
          ]"
        >
          <component :is="item.icon" class="h-[18px] w-[18px]" />
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="border-t border-base-border px-5 py-3">
        <p class="font-mono text-[10px] leading-tight text-ink-faint">
          QRForge v1.0<br />轻量化二维码工具
        </p>
      </div>
    </aside>

    <main class="ml-[220px] flex-1">
      <header class="sticky top-0 z-20 flex h-14 items-center border-b border-base-border bg-base/80 px-6 backdrop-blur-md">
        <h1 class="font-display text-base font-semibold text-ink">{{ pageTitle }}</h1>
      </header>
      <div class="p-6">
        <router-view />
      </div>
    </main>
  </div>
</template>
