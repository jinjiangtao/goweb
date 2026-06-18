<template>
  <div class="day-view">
    <div class="day-header">
      <div class="time-header"></div>
      <div :class="['day-info', { today: isToday(calendarStore.currentDate) }]">
        <div class="day-name">{{ format(calendarStore.currentDate, 'EEEE', { locale: zhCN }) }}</div>
        <div class="day-num">{{ format(calendarStore.currentDate, 'MM月dd日', { locale: zhCN }) }}</div>
      </div>
    </div>
    <div class="day-body" ref="bodyRef">
      <div class="time-column">
        <div v-for="h in hours" :key="h" class="time-slot">
          {{ formatHour(h) }}
        </div>
      </div>
      <div
        class="day-column"
        :class="{ 'drag-over': isDragOver }"
        @dragover.prevent="handleDragOver"
        @dragleave="handleDragLeave"
        @drop="handleDrop"
        @dblclick="handleDblClick"
      >
        <div v-for="h in hours" :key="h" class="hour-cell" :data-hour="h" @dblclick.stop="handleDblClickHour(h)"></div>
        <div class="events-layer">
          <EventCard
            v-for="event in dayEvents"
            :key="event.id + '-' + event.start_time"
            :event="event"
            :view="'day'"
            :day-date="calendarStore.currentDate"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useCalendarStore } from '../stores/calendar'
import {
  isToday, format, startOfDay, addHours, addMinutes, setHours, setMinutes,
  differenceInMinutes, parseISO
} from 'date-fns'
import { zhCN } from 'date-fns/locale'
import EventCard from './EventCard.vue'

const calendarStore = useCalendarStore()
const bodyRef = ref(null)

const hours = Array.from({ length: 24 }, (_, i) => i)
const isDragOver = ref(false)
const dragData = ref(null)

const dayEvents = computed(() => {
  const day = calendarStore.currentDate
  const dayStart = startOfDay(day)
  const dayEnd = addHours(dayStart, 24)
  return calendarStore.events.filter(e => {
    const eStart = new Date(e.start_time)
    const eEnd = new Date(e.end_time)
    return eStart < dayEnd && eEnd > dayStart
  })
})

const formatHour = (h) => {
  return `${h.toString().padStart(2, '0')}:00`
}

const handleDragOver = (e) => {
  isDragOver.value = true
  try {
    const raw = e.dataTransfer.getData('application/json')
    if (raw) {
      dragData.value = JSON.parse(raw)
    }
  } catch (err) {
    dragData.value = null
  }
}

const handleDragLeave = () => {
  isDragOver.value = false
}

const calculateNewTime = (e) => {
  if (!bodyRef.value || !dragData.value) return null
  const rect = e.currentTarget.getBoundingClientRect()
  const y = e.clientY - rect.top
  const scrollTop = bodyRef.value.scrollTop
  const totalMinutes = ((y + scrollTop) / 48) * 60
  const snappedMinutes = Math.round(totalMinutes / 15) * 15
  const clampedMinutes = Math.max(0, Math.min(24 * 60 - 15, snappedMinutes))

  const day = calendarStore.currentDate
  const newStart = addMinutes(startOfDay(day), clampedMinutes)
  const origStart = parseISO(dragData.value.originalStart)
  const origEnd = parseISO(dragData.value.originalEnd)
  const duration = differenceInMinutes(origEnd, origStart)

  if (dragData.value.type === 'move') {
    return {
      start_time: newStart.toISOString(),
      end_time: addMinutes(newStart, duration).toISOString()
    }
  } else if (dragData.value.type === 'resize') {
    const newEnd = newStart
    const minEnd = addMinutes(origStart, 15)
    return {
      start_time: dragData.value.originalStart,
      end_time: newEnd < minEnd ? minEnd.toISOString() : newEnd.toISOString()
    }
  }
  return null
}

const handleDrop = async (e) => {
  isDragOver.value = false
  const timeData = calculateNewTime(e)
  if (!dragData.value || !timeData) return

  try {
    await calendarStore.updateEvent(dragData.value.eventId, timeData)
  } catch (err) {
    console.error('Drag update failed:', err)
  }
  dragData.value = null
}

const handleDblClick = () => {
}

const handleDblClickHour = (hour) => {
  const start = setMinutes(setHours(startOfDay(calendarStore.currentDate), hour), 0)
  const end = addHours(start, 1)
  calendarStore.openNewEventModal({
    title: '',
    description: '',
    start_time: start.toISOString(),
    end_time: end.toISOString(),
    category_id: null,
    is_recurring: false,
    recurring_type: '',
    recurring_end: null,
    has_reminder: false,
    reminder_time: null,
    location: ''
  })
}
</script>

<style scoped>
.day-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.day-header {
  display: grid;
  grid-template-columns: 60px 1fr;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
  position: sticky;
  top: 0;
  z-index: 10;
}

.time-header {
  border-right: 1px solid #e5e7eb;
}

.day-info {
  padding: 12px 24px;
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.day-info.today {
  background: #eff6ff;
}

.day-name {
  font-size: 14px;
  color: #6b7280;
}

.day-num {
  font-size: 22px;
  font-weight: 600;
  color: #1f2937;
}

.day-body {
  display: flex;
  flex: 1;
  overflow: auto;
}

.time-column {
  width: 60px;
  flex-shrink: 0;
  border-right: 1px solid #e5e7eb;
  background: white;
}

.time-slot {
  height: 48px;
  font-size: 11px;
  color: #6b7280;
  text-align: right;
  padding: 2px 8px 0 0;
  border-bottom: 1px solid #f3f4f6;
}

.day-column {
  flex: 1;
  position: relative;
  background: white;
  min-width: 0;
}

.day-column.drag-over {
  background: #eff6ff;
}

.hour-cell {
  height: 48px;
  border-bottom: 1px solid #f3f4f6;
  cursor: pointer;
  transition: background 0.15s;
}

.hour-cell:hover {
  background: #f0f9ff;
}

.events-layer {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}
</style>
