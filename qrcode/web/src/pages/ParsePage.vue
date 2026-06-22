<script setup lang="ts">
import { ref } from 'vue'
import { ScanLine, FileUp, Copy, QrCode, Check } from 'lucide-vue-next'
import { parseText, parseImage, generateQR } from '@/api'
import { defaultStyle } from '@/types'
import type { ParseResult } from '@/types'

const tab = ref<'text' | 'image'>('text')
const textInput = ref('')
const textResults = ref<ParseResult[]>([])
const textLoading = ref(false)
const imageResults = ref<ParseResult[]>([])
const imageLoading = ref(false)
const copiedIdx = ref<number | null>(null)

async function handleParseText() {
  const lines = textInput.value.split('\n').map(l => l.trim()).filter(Boolean)
  if (!lines.length) return
  textLoading.value = true
  try {
    textResults.value = await parseText(lines)
  } finally {
    textLoading.value = false
  }
}

async function handleParseImage(files: FileList | null) {
  if (!files?.length) return
  imageLoading.value = true
  try {
    imageResults.value = await parseImage(Array.from(files))
  } finally {
    imageLoading.value = false
  }
}

async function handleDrop(e: DragEvent) {
  e.preventDefault()
  const files = e.dataTransfer?.files
  if (files?.length) await handleParseImage(files)
}

function handleFileInput(e: Event) {
  const input = e.target as HTMLInputElement
  handleParseImage(input.files)
}

async function copyContent(text: string, idx: number) {
  await navigator.clipboard.writeText(text)
  copiedIdx.value = idx
  setTimeout(() => { copiedIdx.value = null }, 1500)
}

async function generateFromResult(result: ParseResult) {
  await generateQR(result.content, defaultStyle())
}
</script>

<template>
  <div class="stagger space-y-6">
    <div class="flex gap-2">
      <button
        class="btn"
        :class="tab === 'text' ? 'btn-primary' : 'btn-ghost'"
        @click="tab = 'text'"
      >
        <ScanLine class="h-4 w-4" /> 文本解析
      </button>
      <button
        class="btn"
        :class="tab === 'image' ? 'btn-primary' : 'btn-ghost'"
        @click="tab = 'image'"
      >
        <FileUp class="h-4 w-4" /> 图片识别
      </button>
    </div>

    <div v-if="tab === 'text'" class="panel p-5 space-y-4">
      <div>
        <label class="field-label">输入内容（每行一条）</label>
        <textarea
          v-model="textInput"
          class="field-input h-40 font-mono text-sm"
          placeholder="https://example.com&#10;Hello World&#10;weixin://dl/business/..."
        />
      </div>
      <button class="btn btn-primary" :disabled="textLoading" @click="handleParseText">
        <ScanLine class="h-4 w-4" />
        {{ textLoading ? '解析中...' : '批量解析' }}
      </button>

      <div v-if="textResults.length" class="space-y-2">
        <h3 class="field-label">解析结果</h3>
        <div
          v-for="(r, i) in textResults"
          :key="i"
          class="panel-elevated flex items-center gap-3 px-4 py-3"
        >
          <span class="chip shrink-0">{{ r.type === 'url' ? '链接' : '文本' }}</span>
          <span class="flex-1 truncate text-sm text-ink">{{ r.content }}</span>
          <button class="btn btn-ghost px-2 py-1" @click="copyContent(r.content, i)">
            <component :is="copiedIdx === i ? Check : Copy" class="h-3.5 w-3.5" />
          </button>
          <button class="btn btn-ghost-lime px-2 py-1" @click="generateFromResult(r)">
            <QrCode class="h-3.5 w-3.5" /> 生成
          </button>
        </div>
      </div>
    </div>

    <div v-if="tab === 'image'" class="panel p-5 space-y-4">
      <div
        class="flex min-h-[200px] cursor-pointer flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed border-base-border transition-colors hover:border-lime/40"
        @drop="handleDrop"
        @dragover.prevent
        @click="($refs.fileInput as HTMLInputElement)?.click()"
      >
        <FileUp class="h-8 w-8 text-ink-faint" />
        <p class="text-sm text-ink-muted">拖拽图片到此处或点击上传</p>
        <p class="text-xs text-ink-faint">支持 PNG / JPG / GIF / ZIP 压缩包</p>
        <input
          ref="fileInput"
          type="file"
          accept="image/*,.zip"
          multiple
          class="hidden"
          @change="handleFileInput"
        />
      </div>
      <button
        v-if="imageLoading"
        class="btn btn-ghost"
        disabled
      >
        <div class="h-4 w-4 animate-spin rounded-full border-2 border-lime/30 border-t-lime" />
        识别中...
      </button>

      <div v-if="imageResults.length" class="space-y-2">
        <h3 class="field-label">识别结果</h3>
        <div
          v-for="(r, i) in imageResults"
          :key="i"
          class="panel-elevated flex items-center gap-3 px-4 py-3"
        >
          <span
            class="chip shrink-0"
            :class="r.type === 'error' ? 'border-amber/30 text-amber' : ''"
          >
            {{ r.type === 'url' ? '链接' : r.type === 'error' ? '错误' : '文本' }}
          </span>
          <span class="flex-1 truncate text-sm" :class="r.type === 'error' ? 'text-ink-muted' : 'text-ink'">
            {{ r.content }}
          </span>
          <span class="shrink-0 text-xs text-ink-faint">{{ r.raw }}</span>
          <button
            v-if="r.type !== 'error'"
            class="btn btn-ghost-lime px-2 py-1"
            @click="generateFromResult(r)"
          >
            <QrCode class="h-3.5 w-3.5" /> 生成
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
