<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  src: string
  alt?: string
}>()

const loaded = ref(false)
const errored = ref(false)

watch(() => props.src, () => {
  loaded.value = false
  errored.value = false
})
</script>

<template>
  <div class="relative overflow-hidden rounded-lg bg-checker bg-checker">
    <img
      v-if="src && !errored"
      :src="src"
      :alt="alt || 'QR Code'"
      class="block w-full transition-opacity duration-300"
      :class="loaded ? 'opacity-100' : 'opacity-0'"
      @load="loaded = true"
      @error="errored = true"
    />
    <div v-if="!loaded && !errored" class="flex items-center justify-center py-12 text-ink-faint">
      <div class="h-5 w-5 animate-spin rounded-full border-2 border-lime/30 border-t-lime" />
    </div>
    <div v-if="errored" class="flex items-center justify-center py-12 text-ink-faint">
      <span class="text-xs">加载失败</span>
    </div>
    <div
      v-if="loaded"
      class="pointer-events-none absolute inset-x-0 top-0 h-px bg-lime/40"
      style="animation: scan 1.4s ease-in-out"
    />
  </div>
</template>
