<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Search, Download, Trash2, AlertTriangle, RefreshCw,
  ChevronLeft, ChevronRight, Edit3, Check, X, PackageOpen,
} from 'lucide-vue-next'
import {
  listRecords, updateRecordStatus, updateRecordContent,
  deleteRecord, exportRecords, recordImageUrl,
} from '@/api'
import type { Record } from '@/types'

const records = ref<Record[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const statusFilter = ref('')
const keyword = ref('')
const loading = ref(false)
const selected = ref<Set<number>>(new Set())
const editingId = ref<number | null>(null)
const editContent = ref('')

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size.value)))

async function load() {
  loading.value = true
  try {
    const res = await listRecords({
      status: statusFilter.value || undefined,
      keyword: keyword.value || undefined,
      page: page.value,
      size: size.value,
    })
    records.value = res.list || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(load)

function toggleSelect(id: number) {
  if (selected.value.has(id)) selected.value.delete(id)
  else selected.value.add(id)
}

function toggleAll() {
  if (selected.value.size === records.value.length) {
    selected.value.clear()
  } else {
    records.value.forEach(r => selected.value.add(r.id))
  }
}

async function handleStatusToggle(r: Record) {
  const newStatus = r.status === 'active' ? 'invalid' : 'active'
  await updateRecordStatus(r.id, newStatus)
  await load()
}

async function handleDelete(id: number) {
  await deleteRecord(id)
  selected.value.delete(id)
  await load()
}

function startEdit(r: Record) {
  editingId.value = r.id
  editContent.value = r.content
}

async function saveEdit() {
  if (!editingId.value) return
  await updateRecordContent(editingId.value, editContent.value)
  editingId.value = null
  await load()
}

async function handleExport() {
  if (!selected.value.size) return
  const blob = await exportRecords(Array.from(selected.value))
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'qrcode-export.zip'
  a.click()
  URL.revokeObjectURL(url)
}

async function handleBatchInvalid() {
  for (const id of selected.value) {
    await updateRecordStatus(id, 'invalid')
  }
  selected.value.clear()
  await load()
}

async function handleBatchDelete() {
  for (const id of selected.value) {
    await deleteRecord(id)
  }
  selected.value.clear()
  await load()
}
</script>

<template>
  <div class="stagger space-y-5">
    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2">
        <Search class="h-4 w-4 text-ink-faint" />
        <input
          v-model="keyword"
          class="field-input max-w-xs text-sm"
          placeholder="搜索内容或批次号..."
          @keyup.enter="page = 1; load()"
        />
      </div>
      <select v-model="statusFilter" class="field-input w-28 text-sm" @change="page = 1; load()">
        <option value="">全部状态</option>
        <option value="active">有效</option>
        <option value="invalid">失效</option>
      </select>
      <button class="btn btn-ghost px-3 py-2 text-sm" @click="load">
        <RefreshCw class="h-3.5 w-3.5" />
      </button>
      <div class="flex-1" />
      <span class="text-xs text-ink-faint">共 {{ total }} 条</span>
    </div>

    <div v-if="selected.size" class="flex items-center gap-3 rounded-lg bg-lime/5 px-4 py-2.5">
      <span class="text-sm font-medium text-lime">已选 {{ selected.size }} 条</span>
      <button class="btn btn-primary px-3 py-1.5 text-xs" @click="handleExport">
        <Download class="h-3.5 w-3.5" /> 打包导出
      </button>
      <button class="btn btn-ghost-lime px-3 py-1.5 text-xs" @click="handleBatchInvalid">
        <AlertTriangle class="h-3.5 w-3.5" /> 标记失效
      </button>
      <button class="btn btn-danger px-3 py-1.5 text-xs" @click="handleBatchDelete">
        <Trash2 class="h-3.5 w-3.5" /> 批量删除
      </button>
      <button class="btn btn-ghost px-2 py-1.5" @click="selected.clear()">
        <X class="h-3.5 w-3.5" />
      </button>
    </div>

    <div class="panel overflow-hidden">
      <table class="w-full text-left text-sm">
        <thead>
          <tr class="border-b border-base-border bg-base-elevated/60">
            <th class="w-10 px-4 py-3">
              <input
                type="checkbox"
                class="accent-lime"
                :checked="selected.size === records.length && records.length > 0"
                @change="toggleAll"
              />
            </th>
            <th class="px-4 py-3 font-mono text-[11px] uppercase tracking-wider text-ink-muted">预览</th>
            <th class="px-4 py-3 font-mono text-[11px] uppercase tracking-wider text-ink-muted">内容</th>
            <th class="px-4 py-3 font-mono text-[11px] uppercase tracking-wider text-ink-muted">类型</th>
            <th class="px-4 py-3 font-mono text-[11px] uppercase tracking-wider text-ink-muted">来源</th>
            <th class="px-4 py-3 font-mono text-[11px] uppercase tracking-wider text-ink-muted">状态</th>
            <th class="px-4 py-3 font-mono text-[11px] uppercase tracking-wider text-ink-muted">创建时间</th>
            <th class="px-4 py-3 font-mono text-[11px] uppercase tracking-wider text-ink-muted">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="r in records"
            :key="r.id"
            class="border-b border-base-border/50 transition hover:bg-base-hover/40"
          >
            <td class="px-4 py-3">
              <input type="checkbox" class="accent-lime" :checked="selected.has(r.id)" @change="toggleSelect(r.id)" />
            </td>
            <td class="px-4 py-3">
              <img
                v-if="r.id"
                :src="recordImageUrl(r.id)"
                :alt="r.content"
                class="h-10 w-10 rounded bg-checker bg-checker object-contain p-0.5"
              />
            </td>
            <td class="max-w-[240px] px-4 py-3">
              <div v-if="editingId === r.id" class="flex items-center gap-2">
                <input v-model="editContent" class="field-input text-xs" />
                <button class="text-lime" @click="saveEdit"><Check class="h-4 w-4" /></button>
                <button class="text-ink-muted" @click="editingId = null"><X class="h-4 w-4" /></button>
              </div>
              <span v-else class="block truncate text-xs text-ink">{{ r.content }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="chip">{{ r.contentType === 'url' ? '链接' : '文本' }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="chip">{{ r.source === 'batch' ? '批量' : '单条' }}</span>
            </td>
            <td class="px-4 py-3">
              <span
                class="chip cursor-pointer"
                :class="r.status === 'invalid' ? 'border-amber/30 text-amber' : 'border-lime/20 text-lime'"
                @click="handleStatusToggle(r)"
              >
                {{ r.status === 'active' ? '有效' : '失效' }}
              </span>
            </td>
            <td class="px-4 py-3 font-mono text-xs text-ink-muted">
              {{ r.createdAt?.slice(0, 16) }}
            </td>
            <td class="px-4 py-3">
              <div class="flex items-center gap-1">
                <button class="btn btn-ghost px-1.5 py-1" title="修改内容" @click="startEdit(r)">
                  <Edit3 class="h-3.5 w-3.5" />
                </button>
                <a
                  :href="recordImageUrl(r.id)"
                  download
                  class="btn btn-ghost px-1.5 py-1"
                  title="下载"
                >
                  <Download class="h-3.5 w-3.5" />
                </a>
                <button class="btn btn-ghost px-1.5 py-1 text-amber" title="标记失效" @click="handleStatusToggle(r)">
                  <AlertTriangle class="h-3.5 w-3.5" />
                </button>
                <button class="btn btn-ghost px-1.5 py-1 text-ink-faint hover:text-amber" @click="handleDelete(r.id)">
                  <Trash2 class="h-3.5 w-3.5" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="!records.length && !loading" class="flex items-center justify-center py-16 text-sm text-ink-faint">
        <PackageOpen class="mr-2 h-5 w-5" /> 暂无记录
      </div>
    </div>

    <div v-if="totalPages > 1" class="flex items-center justify-between">
      <p class="text-xs text-ink-faint">第 {{ page }} / {{ totalPages }} 页</p>
      <div class="flex gap-2">
        <button class="btn btn-ghost px-2 py-1" :disabled="page <= 1" @click="page--; load()">
          <ChevronLeft class="h-4 w-4" />
        </button>
        <button class="btn btn-ghost px-2 py-1" :disabled="page >= totalPages" @click="page++; load()">
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
    </div>
  </div>
</template>
