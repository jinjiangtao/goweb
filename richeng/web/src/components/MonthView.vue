<template>
  <div class="month-view">
    <div class="weekday-header">
      <div v-for="d in weekdayLabels" :key="d" class="weekday-cell">
        {{ d }}
      </div>
    </div>
    <div class="days-grid">
      <div
        v-for="day in monthDays"
        :key="day.key"
        :class="['day-cell', {
          'other-month': !day.isCurrentMonth,
          'today': day.isToday,
          'weekend': day.isWeekend
        }]"
        @dblclick="handleDblClick(day)"
      >
        <div class="day-number">{{ day.day }}</div>
        <div class="day-events">
          <EventCard
            v-for="event in day.events"
            :key="event.id + '-' + event.start_time"
            :event="event"
            :view="'month'"
          />
          <div v-if="day.moreCount > 0" class="more-events" @click.stop="showDay(day)">
            +{{ day.moreCount }} 更多
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useCalendarStore } from '../stores/calendar'
import {
  startOfMonth, endOfMonth, startOfWeek, endOfWeek,
  eachDayOfInterval, isSameMonth, isSameDay, isToday,
  isSameMonth as checkSameMonth, format, startOfDay, addHours
} from 'date-fns'
import { zhCN } from 'date-fns/locale'
import EventCard from './EventCard.vue'

const calendarStore = useCalendarStore()

const weekdayLabels = computed(() => {
  return ['星期一', '星期二', '星期三', '星期四', '星期五', '星期六', '星期日']
})

const monthDays = computed(() => {
  const d = calendarStore.currentDate
  const start = startOfWeek(startOfMonth(d), { weekStartsOn: 1 })
  const end = endOfWeek(endOfMonth(d), { weekStartsOn: 1 })
  const days = eachDayOfInterval({ start, end })

  return days.map(day => {
    const dayStart = startOfDay(day)
    const dayEnd = addHours(dayStart, 24)
    const dayEvents = calendarStore.events.filter(e => {
      const eStart = new Date(e.start_time)
      const eEnd = new Date(e.end_time)
      return eStart < dayEnd && eEnd > dayStart
    }).slice(0, 3)
    const allEvents = calendarStore.events.filter(e => {
      const eStart = new Date(e.start_time)
      const eEnd = new Date(e.end_time)
      return eStart < dayEnd && eEnd > dayStart
    })

    return {
      key: day.toISOString(),
      date: day,
      day: day.getDate(),
      isCurrentMonth: isSameMonth(day, d),
      isToday: isToday(day),
      isWeekend: day.getDay() === 0 || day.getDay() === 6,
      events: dayEvents,
      moreCount: Math.max(0, allEvents.length - 3)
    }
  })
})

const handleDblClick = (day) => {
  const start = addHours(startOfDay(day.date), 9)
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

const showDay = (day) => {
  calendarStore.selectDate(day.date)
}
</script>

<style scoped>
.month-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.weekday-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
  flex-shrink: 0;
}

.weekday-cell {
  padding: 10px 8px;
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
  text-align: center;
  border-right: 1px solid #e5e7eb;
}

.weekday-cell:last-child {
  border-right: none;
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  grid-auto-rows: 1fr;
  flex: 1;
  min-height: 0;
}

.day-cell {
  border-right: 1px solid #e5e7eb;
  border-bottom: 1px solid #e5e7eb;
  padding: 6px;
  min-height: 0;
  display: flex;
  flex-direction: column;
  cursor: pointer;
  transition: background 0.15s;
  overflow: hidden;
}

.day-cell:hover {
  background: #f9fafb;
}

.day-cell.other-month {
  background: #fafafa;
}

.day-cell.today {
  background: #eff6ff;
}

.day-cell.weekend {
  background: #fafafa;
}

.day-cell.today.weekend {
  background: #eff6ff;
}

.day-number {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 4px;
}

.day-cell.other-month .day-number {
  color: #d1d5db;
}

.day-cell.today .day-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #3b82f6;
  color: white;
  font-weight: 600;
}

.day-events {
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
  flex: 1;
}

.more-events {
  font-size: 12px;
  color: #3b82f6;
  padding: 2px 4px;
  border-radius: 4px;
  cursor: pointer;
}

.more-events:hover {
  background: #dbeafe;
}
</style>
