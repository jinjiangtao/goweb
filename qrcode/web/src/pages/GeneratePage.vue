<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import {
  QrCode, Upload, RotateCcw, Save, ListPlus, FileText,
  Palette, Settings2, Eye,
} from 'lucide-vue-next'
import { previewQR, generateQR, batchGenerate, listTemplates } from '@/api'
import { defaultStyle } from '@/types'
import type { StyleConfig, Template } from '@/types'

const mode = ref<'single' | 'batch'>('single')
const content = ref('')
const batchText = ref('')
const style = ref<StyleConfig>(defaultStyle())
const previewUrl = ref('')
const previewLoading = ref(false)
const generating = ref(false)
const batchGenerating = ref(false)
const saveTemplate = ref(false)
const templateName = ref('')
const templates = ref<Template[]>([])
const showTemplatePicker = ref(false)
const lastResult = ref<{ id: number; previewUrl: string } | null>(null)

let previewTimer: ReturnType<typeof setTimeout> | null = null

const cornerOptions: { value: StyleConfig['cornerStyle']; label: string }[] = [
  { value: 'square', label: '直角' },
  { value: 'rounded', label: '圆角' },
  { value: 'dot', label: '圆点' },
]

const eccOptions: { value: StyleConfig['eccLevel']; label: string }[] = [
  { value: 'L', label: 'L (7%)' },
  { value: 'M', label: 'M (15%)' },
  { value: 'Q', label: 'Q (25%)' },
  { value: 'H', label: 'H (30%)' },
]

function schedulePreview() {
  if (previewTimer) clearTimeout(previewTimer)
  previewTimer = setTimeout(doPreview, 400)
}

watch([content, style], schedulePreview, { deep: true })

async function doPreview() {
  if (!content.value.trim()) {
    previewUrl.value = ''
    return
  }
  previewLoading.value = true
  try {
    const blob = await previewQR(content.value, style.value)
    previewUrl.value = URL.createObjectURL(blob)
  } catch {
    previewUrl.value = ''
  } finally {
    previewLoading.value = false
  }
}

async function handleGenerate() {
  if (!content.value.trim()) return
  generating.value = true
  try {
    const res = await generateQR(content.value, style.value, {
      saveTemplate: saveTemplate.value || undefined,
      templateName: templateName.value || undefined,
    })
    lastResult.value = res
    previewUrl.value = res.previewUrl
  } finally {
    generating.value = false
  }
}

function getBatchItems() {
  return batchText.value
    .split('\n')
    .map(l => l.trim())
    .filter(Boolean)
    .map(l => ({ content: l }))
}

async function handleBatchGenerate() {
  const items = getBatchItems()
  if (!items.length) return
  batchGenerating.value = true
  try {
    await batchGenerate(items, style.value)
  } finally {
    batchGenerating.value = false
  }
}

async function loadTemplates() {
  templates.value = await listTemplates()
}

function applyTemplate(t: Template) {
  style.value = { ...t.styleConfig }
  showTemplatePicker.value = false
}

function handleLogoUpload(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    style.value.logo = reader.result as string
  }
  reader.readAsDataURL(file)
}

function removeLogo() {
  style.value.logo = ''
}

function resetStyle() {
  style.value = defaultStyle()
}

const previewObjectUrl = computed(() => {
  if (previewUrl.value.startsWith('blob:') || previewUrl.value.startsWith('/')) return previewUrl.value
  return ''
})
</script>

<template>
  <div class="stagger space-y-6">
    <div class="flex gap-2">
      <button class="btn" :class="mode === 'single' ? 'btn-primary' : 'btn-ghost'" @click="mode = 'single'">
        <QrCode class="h-4 w-4" /> 单条生成
      </button>
      <button class="btn" :class="mode === 'batch' ? 'btn-primary' : 'btn-ghost'" @click="mode = 'batch'">
        <ListPlus class="h-4 w-4" /> 批量生成
      </button>
    </div>

    <div class="grid grid-cols-12 gap-6">
      <div class="col-span-4 space-y-5">
        <div class="panel p-5 space-y-4">
          <h3 class="field-label mb-3">
            <FileText class="inline h-3.5 w-3.5" /> 内容输入
          </h3>
          <div v-if="mode === 'single'">
            <textarea
              v-model="content"
              class="field-input h-32 font-mono text-sm"
              placeholder="输入文本或链接..."
            />
          </div>
          <div v-else>
            <textarea
              v-model="batchText"
              class="field-input h-48 font-mono text-sm"
              placeholder="每行输入一条内容&#10;https://example.com&#10;Hello World"
            />
            <p class="mt-1.5 text-xs text-ink-faint">共 {{ getBatchItems().length }} 条</p>
          </div>
        </div>

        <div class="flex gap-2">
          <button
            v-if="mode === 'single'"
            class="btn btn-primary flex-1"
            :disabled="generating || !content.trim()"
            @click="handleGenerate"
          >
            <QrCode class="h-4 w-4" />
            {{ generating ? '生成中...' : '生成并归档' }}
          </button>
          <button
            v-else
            class="btn btn-primary flex-1"
            :disabled="batchGenerating || !getBatchItems().length"
            @click="handleBatchGenerate"
          >
            <ListPlus class="h-4 w-4" />
            {{ batchGenerating ? '批量生成中...' : `批量生成 (${getBatchItems().length})` }}
          </button>
        </div>

        <div v-if="mode === 'single'" class="panel p-4 space-y-3">
          <label class="flex items-center gap-2 text-sm text-ink-secondary">
            <input v-model="saveTemplate" type="checkbox" class="accent-lime" />
            保存为模板
          </label>
          <input
            v-if="saveTemplate"
            v-model="templateName"
            class="field-input text-sm"
            placeholder="模板名称"
          />
        </div>
      </div>

      <div class="col-span-4 space-y-5">
        <div class="panel p-5 space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="field-label mb-0">
              <Settings2 class="inline h-3.5 w-3.5" /> 样式配置
            </h3>
            <div class="flex gap-1.5">
              <button class="btn btn-ghost px-2 py-1" @click="showTemplatePicker = true; loadTemplates()">
                <Palette class="h-3.5 w-3.5" /> 模板
              </button>
              <button class="btn btn-ghost px-2 py-1" @click="resetStyle">
                <RotateCcw class="h-3.5 w-3.5" />
              </button>
            </div>
          </div>

          <div>
            <label class="field-label">尺寸 (px)</label>
            <input v-model.number="style.size" type="number" min="128" max="2048" step="64" class="field-input text-sm" />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="field-label">前景色</label>
              <div class="flex items-center gap-2">
                <input v-model="style.foreground" type="color" class="h-8 w-8 cursor-pointer rounded border border-base-border bg-transparent" />
                <input v-model="style.foreground" class="field-input text-xs font-mono" />
              </div>
            </div>
            <div>
              <label class="field-label">背景色</label>
              <div class="flex items-center gap-2">
                <input v-model="style.background" type="color" class="h-8 w-8 cursor-pointer rounded border border-base-border bg-transparent" />
                <input v-model="style.background" class="field-input text-xs font-mono" />
              </div>
            </div>
          </div>

          <div>
            <label class="field-label">边角样式</label>
            <div class="flex gap-2">
              <button
                v-for="opt in cornerOptions"
                :key="opt.value"
                class="btn px-3 py-1.5 text-xs"
                :class="style.cornerStyle === opt.value ? 'btn-primary' : 'btn-ghost'"
                @click="style.cornerStyle = opt.value"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>

          <div>
            <label class="field-label">容错级别</label>
            <div class="flex gap-2">
              <button
                v-for="opt in eccOptions"
                :key="opt.value"
                class="btn px-3 py-1.5 text-xs"
                :class="style.eccLevel === opt.value ? 'btn-primary' : 'btn-ghost'"
                @click="style.eccLevel = opt.value"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>

          <div>
            <label class="field-label">外边距</label>
            <input v-model.number="style.margin" type="number" min="0" max="20" class="field-input text-sm" />
          </div>

          <div>
            <label class="field-label">内嵌 Logo</label>
            <div v-if="style.logo" class="flex items-center gap-3">
              <img :src="style.logo" alt="Logo" class="h-10 w-10 rounded border border-base-border object-contain" />
              <button class="btn btn-ghost px-2 py-1 text-xs" @click="removeLogo">移除</button>
            </div>
            <label v-else class="btn btn-ghost cursor-pointer">
              <Upload class="h-3.5 w-3.5" /> 上传 Logo
              <input type="file" accept="image/*" class="hidden" @change="handleLogoUpload" />
            </label>
          </div>
        </div>
      </div>

      <div class="col-span-4">
        <div class="panel p-5">
          <h3 class="field-label mb-4">
            <Eye class="inline h-3.5 w-3.5" /> 实时预览
          </h3>
          <div class="mx-auto max-w-[320px]">
            <div v-if="previewObjectUrl" class="relative overflow-hidden rounded-lg bg-checker bg-checker">
              <img :src="previewObjectUrl" alt="Preview" class="block w-full" />
              <div v-if="previewLoading" class="absolute inset-0 flex items-center justify-center bg-base/60">
                <div class="h-6 w-6 animate-spin rounded-full border-2 border-lime/30 border-t-lime" />
              </div>
            </div>
            <div v-else class="flex aspect-square items-center justify-center rounded-lg border border-dashed border-base-border text-sm text-ink-faint">
              输入内容后预览
            </div>
          </div>

          <div v-if="lastResult" class="mt-4 rounded-lg bg-lime/5 px-4 py-3">
            <p class="text-xs text-ink-secondary">
              已归档至记录 <span class="font-mono text-lime">#{{ lastResult.id }}</span>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="showTemplatePicker"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      @click.self="showTemplatePicker = false"
    >
      <div class="panel w-full max-w-lg p-6">
        <h3 class="mb-4 font-display text-base font-semibold text-ink">选择模板</h3>
        <div v-if="templates.length" class="grid max-h-80 grid-cols-3 gap-3 overflow-y-auto">
          <button
            v-for="t in templates"
            :key="t.id"
            class="panel-elevated overflow-hidden text-left transition hover:shadow-glow"
            @click="applyTemplate(t)"
          >
            <div class="aspect-square bg-checker bg-checker p-2">
              <img :src="`/api/templates/${t.id}/preview`" :alt="t.name" class="h-full w-full object-contain" />
            </div>
            <p class="truncate px-2 py-1.5 text-xs font-medium text-ink">{{ t.name }}</p>
          </button>
        </div>
        <p v-else class="py-8 text-center text-sm text-ink-faint">暂无模板，请先在生成时保存</p>
        <button class="btn btn-ghost mt-4 w-full" @click="showTemplatePicker = false">关闭</button>
      </div>
    </div>
  </div>
</template>
