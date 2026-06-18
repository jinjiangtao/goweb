<template>
  <Teleport to="body">
    <div
      v-if="hoverInfo"
      class="tooltip"
      :style="tooltipStyle"
    >
      <div class="tooltip-header">
        <span class="color-dot" :style="{ background: categoryColor }"></span>
        <h3 class="tooltip-title">{{ hoverInfo.event.title || '(无标题)' }}</h3>
      </div>
      <div class="tooltip-body">
        <div class="tooltip-row">
          <span class="row-label">时间</span>
          <span class="row-value">{{ timeRange }}</span>
        </div>
        <div v-if="hoverInfo.event.location" class="tooltip-row">
          <span class="row-label">地点</span>
          <span class="row-value">{{ hoverInfo.event.location }}</span>
        </div>
        <div v-if="hoverInfo.event.category" class="tooltip-row">
          <span class="row-label">分类</span>
          <span class="row-value">{{ hoverInfo.event.category.name }}</span>
        </div>
        <div v-if="hoverInfo.event.is_recurring" class="tooltip-row">
          <span class="row-label">重复</span>
          <span class="row-value">{{ recurringLabel }}</span>
        </div>
        <div v-if="hoverInfo.event.has_reminder" class="tooltip-row">
          <span class="row-label">提醒</span>
          <span class="row-value">{{ reminderLabel }}</span>
        </div>
        <div v-if="hoverInfo.event.description" class="tooltip-row tooltip-desc">
          <span class="row-label">描述</span>
          <span class="row-value">{{ hoverInfo.event.description }}</span>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { useCalendarStore } from '../stores/calendar'
import { format, parseISO } from 'date-fns'
import { zhCN } from 'date-fns/locale'

const calendarStore = useCalendarStore()

const hoverInfo = computed(() => calendarStore.hoveredEvent)

const categoryColor = computed(() => {
  if (!hoverInfo.value) return '#6b7280'
  const e = hoverInfo.value.event
  if (e.category) return e.category.color
  if (e.category_id && calendarStore.categoryMap[e.category_id]) {
    return calendarStore.categoryMap[e.category_id].color
  }
  return '#6b7280'
})

const timeRange = computed(() => {
  if (!hoverInfo.value) return ''
  const e = hoverInfo.value.event
  const start = parseISO(e.start_time)
  const end = parseISO(e.end_time)
  return `${format(start, 'yyyy-MM-dd HH:mm', { locale: zhCN })} ~ ${format(end, 'HH:mm', { locale: zhCN })}`
})

const recurringLabel = computed(() => {
  if (!hoverInfo.value) return ''
  const type = hoverInfo.value.event.recurring_type
  const map = { daily: '每天', weekly: '每周', monthly: '每月' }
  return map[type] || type
})

const reminderLabel = computed(() => {
  if (!hoverInfo.value || !hoverInfo.value.event.reminder_time) return ''
  return format(parseISO(hoverInfo.value.event.reminder_time), 'yyyy-MM-dd HH:mm', { locale: zhCN })
})

const tooltipStyle = computed(() => {
  if (!hoverInfo.value) return {}
  let x = hoverInfo.value.x + 16
  let y = hoverInfo.value.y + 16
  const maxX = window.innerWidth - 300
  const maxY = window.innerHeight - 300
  if (x > maxX) x = hoverInfo.value.x - 300
  if (y > maxY) y = hoverInfo.value.y - 200
  return {
    left: `${x}px`,
    top: `${y}px`
  }
})
</script>

<style scoped>
.tooltip {
  position: fixed;
  z-index: 2000;
  width: 280px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: 1px solid #e5e7eb;
  pointer-events: none;
  overflow: hidden;
}

.tooltip-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.color-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.tooltip-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.tooltip-body {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tooltip-row {
  display: flex;
  gap: 12px;
  font-size: 13px;
}

.tooltip-desc {
  flex-direction: column;
  gap: 4px;
}

.row-label {
  width: 40px;
  flex-shrink: 0;
  color: #6b7280;
  font-weight: 500;
}

.tooltip-desc .row-label {
  width: auto;
}

.row-value {
  color: #374151;
  flex: 1;
  word-break: break-all;
  white-space: pre-wrap;
}
</style>
