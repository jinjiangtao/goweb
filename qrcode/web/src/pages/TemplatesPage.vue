<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Trash2, Edit3, QrCode, X, Check } from 'lucide-vue-next'
import { listTemplates, createTemplate, updateTemplate, deleteTemplate } from '@/api'
import { defaultStyle } from '@/types'
import type { Template, StyleConfig } from '@/types'

const templates = ref<Template[]>([])
const searchKey = ref('')
const showCreate = ref(false)
const editId = ref<number | null>(null)
const formName = ref('')
const formStyle = ref<StyleConfig>(defaultStyle())

async function load() {
  templates.value = await listTemplates(searchKey.value || undefined)
}

onMounted(load)

function openCreate() {
  formName.value = ''
  formStyle.value = defaultStyle()
  showCreate.value = true
  editId.value = null
}

function openEdit(t: Template) {
  formName.value = t.name
  formStyle.value = { ...t.styleConfig }
  editId.value = t.id
  showCreate.value = true
}

async function handleSave() {
  if (!formName.value.trim()) return
  if (editId.value) {
    await updateTemplate(editId.value, formName.value, formStyle.value)
  } else {
    await createTemplate(formName.value, formStyle.value)
  }
  showCreate.value = false
  editId.value = null
  await load()
}

async function handleDelete(id: number) {
  await deleteTemplate(id)
  await load()
}
</script>

<template>
  <div class="stagger space-y-6">
    <div class="flex items-center gap-3">
      <input
        v-model="searchKey"
        class="field-input max-w-xs text-sm"
        placeholder="搜索模板..."
        @keyup.enter="load"
      />
      <div class="flex-1" />
      <button class="btn btn-primary" @click="openCreate">
        <Plus class="h-4 w-4" /> 新建模板
      </button>
    </div>

    <div v-if="templates.length" class="grid grid-cols-4 gap-5">
      <div
        v-for="t in templates"
        :key="t.id"
        class="panel overflow-hidden transition hover:shadow-glow"
      >
        <div class="aspect-square bg-checker bg-checker p-3">
          <img
            :src="`/api/templates/${t.id}/preview`"
            :alt="t.name"
            class="h-full w-full object-contain"
          />
        </div>
        <div class="px-4 py-3">
          <p class="truncate text-sm font-medium text-ink">{{ t.name }}</p>
          <div class="mt-1.5 flex gap-1.5">
            <span class="chip">{{ t.styleConfig.cornerStyle }}</span>
            <span class="chip">ECC {{ t.styleConfig.eccLevel }}</span>
            <span class="chip">{{ t.styleConfig.size }}px</span>
          </div>
          <div class="mt-3 flex gap-2">
            <button class="btn btn-ghost-lime flex-1 px-2 py-1.5 text-xs" @click="openEdit(t)">
              <Edit3 class="h-3.5 w-3.5" /> 编辑
            </button>
            <button class="btn btn-danger px-2 py-1.5 text-xs" @click="handleDelete(t.id)">
              <Trash2 class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="panel flex items-center justify-center py-20 text-sm text-ink-faint">
      暂无模板，可在生成中心保存样式为模板
    </div>

    <div
      v-if="showCreate"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      @click.self="showCreate = false"
    >
      <div class="panel w-full max-w-md p-6 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="font-display text-base font-semibold text-ink">
            {{ editId ? '编辑模板' : '新建模板' }}
          </h3>
          <button class="text-ink-faint hover:text-ink" @click="showCreate = false">
            <X class="h-5 w-5" />
          </button>
        </div>
        <div>
          <label class="field-label">模板名称</label>
          <input v-model="formName" class="field-input text-sm" placeholder="输入模板名称" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="field-label">前景色</label>
            <div class="flex items-center gap-2">
              <input v-model="formStyle.foreground" type="color" class="h-8 w-8 cursor-pointer rounded border border-base-border bg-transparent" />
              <input v-model="formStyle.foreground" class="field-input text-xs font-mono" />
            </div>
          </div>
          <div>
            <label class="field-label">背景色</label>
            <div class="flex items-center gap-2">
              <input v-model="formStyle.background" type="color" class="h-8 w-8 cursor-pointer rounded border border-base-border bg-transparent" />
              <input v-model="formStyle.background" class="field-input text-xs font-mono" />
            </div>
          </div>
        </div>
        <div>
          <label class="field-label">边角样式</label>
          <div class="flex gap-2">
            <button
              v-for="opt in [{ value: 'square', label: '直角' }, { value: 'rounded', label: '圆角' }, { value: 'dot', label: '圆点' }]"
              :key="opt.value"
              class="btn px-3 py-1.5 text-xs"
              :class="formStyle.cornerStyle === opt.value ? 'btn-primary' : 'btn-ghost'"
              @click="formStyle.cornerStyle = opt.value as StyleConfig['cornerStyle']"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>
        <div>
          <label class="field-label">容错级别</label>
          <select v-model="formStyle.eccLevel" class="field-input text-sm">
            <option value="L">L (7%)</option>
            <option value="M">M (15%)</option>
            <option value="Q">Q (25%)</option>
            <option value="H">H (30%)</option>
          </select>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="field-label">尺寸</label>
            <input v-model.number="formStyle.size" type="number" class="field-input text-sm" />
          </div>
          <div>
            <label class="field-label">边距</label>
            <input v-model.number="formStyle.margin" type="number" class="field-input text-sm" />
          </div>
        </div>
        <button class="btn btn-primary w-full" @click="handleSave">
          <Check class="h-4 w-4" /> {{ editId ? '更新模板' : '创建模板' }}
        </button>
      </div>
    </div>
  </div>
</template>
