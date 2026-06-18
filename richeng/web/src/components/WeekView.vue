<template>
  <div class="week-view" ref="viewRef">
    <div class="week-header">
      <div class="time-header"></div>
      <div
        v-for="day in weekDays"
        :key="day.key"
        :class="['day-header', { today: day.isToday, weekend: day.isWeekend }]"
      >
        <div class="day-name">{{ day.weekdayName }}</div>
        <div class="day-num">{{ day.day }}</div>
      </div>
    </div>
    <div class="week-body" ref="bodyRef">
      <div class="time-column">
        <div v-for="h in hours" :key="h" class="time-slot">
          {{ formatHour(h) }}
        </div>
      </div>
      <div class="days-container">
        <div
          v-for="day in weekDays"
          :key="day.key"
          :class="['day-column', { today: day.isToday, weekend: day.isWeekend, 'drag-over': dragOverDay === day.key }]"
          :data-day-key="day.key"
          @dragover.prevent="handleDragOver($event, day)"
          @dragleave="handleDragLeave"
          @drop="handleDrop($event, day)"
          @dblclick="handleDblClickDay($event, day)"
        >
          <div v-for="h in hours" :key="h" class="hour-cell" :data-hour="h" @dblclick.stop="handleDblClickHour($event, day, h)"></div>
          <div class="events-layer">
            <EventCard
              v-for="event in day.events"
              :key="event.id + '-' + event.start_time"
              :event="event"
              :view="'week'"
              :day-date="day.date"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useCalendarStore } from '../stores/calendar'
import {
  startOfWeek, endOfWeek, eachDayOfInterval,
  isToday, format, startOfDay, addHours, addMinutes, setHours, setMinutes,
  differenceInMinutes, parseISO
} from 'date-fns'
import { zhCN } from 'date-fns/locale'
import EventCard from './EventCard.vue'

const calendarStore = useCalendarStore()
const viewRef = ref(null)
const bodyRef = ref(null)

const hours = Array.from({ length: 24 }, (_, i) => i)
const dragOverDay = ref(null)
const dragData = ref(null)

const weekDays = computed(() => {
  const d = calendarStore.currentDate
  const start = startOfWeek(d, { weekStartsOn: 1 })
  const end = endOfWeek(d, { weekStartsOn: 1 })
  const days = eachDayOfInterval({ start, end })

  return days.map(day => {
    const dayStart = startOfDay(day)
    const dayEnd = addHours(dayStart, 24)
    const dayEvents = calendarStore.events.filter(e => {
      const eStart = new Date(e.start_time)
      const eEnd = new Date(e.end_time)
      return eStart < dayEnd && eEnd > dayStart
    })

    return {
      key: day.toISOString(),
      date: day,
      day: day.getDate(),
      weekdayName: format(day, 'EEE', { locale: zhCN }),
      isToday: isToday(day),
      isWeekend: day.getDay() === 0 || day.getDay() === 6,
      events: dayEvents
    }
  })
})

const formatHour = (h) => {
  return `${h.toString().padStart(2, '0')}:00`
}

const handleDragOver = (e, day) => {
  dragOverDay.value = day.key
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
}

const calculateNewTime = (e, day) => {
  if (!bodyRef.value || !dragData.value) return null
  const rect = e.currentTarget.getBoundingClientRect()
  const y = e.clientY - rect.top
  const scrollTop = bodyRef.value.scrollTop
  const totalMinutes = ((y + scrollTop) / 48) * 60
  const snappedMinutes = Math.round(totalMinutes / 15) * 15
  const clampedMinutes = Math.max(0, Math.min(24 * 60 - 15, snappedMinutes))

  const newStart = addMinutes(startOfDay(day.date), clampedMinutes)
  const origStart = parseISO(dragData.value.originalStart)
  const origEnd = parseISO(dragData.value.originalEnd)
  const duration = differenceInMinutes(origEnd, origStart)

  if (dragData.value.type === 'move') {
    return {
      start_time: newStart.toISOString(),
      end_time: addMinutes(newStart, duration).toISOString()
    }
  } else if (dragData.value.type === 'resize') {
    const newEnd = addMinutes(newStart, 0)
    const minEnd = addMinutes(origStart, 15)
    return {
      start_time: dragData.value.originalStart,
      end_time: newEnd < minEnd ? minEnd.toISOString() : newEnd.toISOString()
    }
  }
  return null
}

const handleDrop = async (e, day) => {
  dragOverDay.value = null
  const timeData = calculateNewTime(e, day)
  if (!dragData.value || !timeData) return

  try {
    await calendarStore.updateEvent(dragData.value.eventId, timeData)
  } catch (err) {
    console.error('Drag update failed:', err)
  }
  dragData.value = null
}

const handleDblClickDay = (e, day) => {
}

const handleDblClickHour = (e, day, hour) => {
  const start = setMinutes(setHours(startOfDay(day.date), hour), 0)
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
.week-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 100%;
}

.week-header {
  display: grid;
  grid-template-columns: 60px repeat(7, 1fr);
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
  position: sticky;
  top: 0;
  z-index: 10;
}

.time-header {
  border-right: 1px solid #e5e7eb;
}

.day-header {
  padding: 10px 4px;
  text-align: center;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.day-header:last-child {
  border-right: none;
}

.day-header.today {
  background: #eff6ff;
}

.day-name {
  font-size: 12px;
  color: #6b7280;
}

.day-num {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.day-header.today .day-num {
  background: #3b82f6;
  color: white;
}

.week-body {
  display: flex;
  flex: 1;
  overflow: auto;
  min-height: 0;
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

.days-container {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  flex: 1;
  min-width: 0;
}

.day-column {
  position: relative;
  border-right: 1px solid #e5e7eb;
  background: white;
}

.day-column:last-child {
  border-right: none;
}

.day-column.weekend {
  background: #fafafa;
}

.day-column.today {
  background: #f8fafc;
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
