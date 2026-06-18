<template>
  <header class="calendar-header">
    <div class="header-left">
      <h1 class="logo">📅 日程日历</h1>
    </div>
    <div class="header-center">
      <div class="nav-buttons">
        <button class="btn-nav" @click="calendarStore.navigatePrev">‹</button>
        <button class="btn-today" @click="calendarStore.goToday">今天</button>
        <button class="btn-nav" @click="calendarStore.navigateNext">›</button>
      </div>
      <h2 class="current-title">{{ calendarStore.viewTitle }}</h2>
    </div>
    <div class="header-right">
      <div class="view-switcher">
        <button
          v-for="v in views"
          :key="v.value"
          :class="['btn-view', { active: calendarStore.view === v.value }]"
          @click="calendarStore.setView(v.value)"
        >
          {{ v.label }}
        </button>
      </div>
      <button class="btn-new" @click="handleNewEvent">
        + 新建日程
      </button>
    </div>
  </header>
</template>

<script setup>
import { useCalendarStore } from '../stores/calendar'
import { startOfDay, addHours } from 'date-fns'

const calendarStore = useCalendarStore()

const views = [
  { value: 'month', label: '月' },
  { value: 'week', label: '周' },
  { value: 'day', label: '日' }
]

const handleNewEvent = () => {
  const now = new Date()
  const start = addHours(startOfDay(now), now.getHours() + 1)
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
.calendar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.header-left .logo {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.header-center {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-buttons {
  display: flex;
  align-items: center;
  gap: 4px;
}

.btn-nav {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  font-size: 18px;
  color: #4b5563;
  transition: all 0.15s;
}

.btn-nav:hover {
  background: #f3f4f6;
}

.btn-today {
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 14px;
  color: #4b5563;
  border: 1px solid #d1d5db;
  transition: all 0.15s;
}

.btn-today:hover {
  background: #f9fafb;
}

.current-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  min-width: 280px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.view-switcher {
  display: flex;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  overflow: hidden;
}

.btn-view {
  padding: 6px 16px;
  font-size: 14px;
  color: #4b5563;
  transition: all 0.15s;
}

.btn-view:hover {
  background: #f9fafb;
}

.btn-view.active {
  background: #3b82f6;
  color: white;
}

.btn-new {
  padding: 8px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  background: #3b82f6;
  color: white;
  transition: all 0.15s;
}

.btn-new:hover {
  background: #2563eb;
}
</style>
