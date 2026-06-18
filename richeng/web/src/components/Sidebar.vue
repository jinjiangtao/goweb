<template>
  <aside class="sidebar">
    <div class="mini-calendar">
      <div class="mini-header">
        <button @click="miniPrev">‹</button>
        <span>{{ format(miniMonth, 'yyyy年MM月', { locale: zhCN }) }}</span>
        <button @click="miniNext">›</button>
      </div>
      <div class="mini-weekdays">
        <span v-for="d in weekdayLabels" :key="d">{{ d }}</span>
      </div>
      <div class="mini-days">
        <div
          v-for="day in miniDays"
          :key="day.key"
          :class="['mini-day', {
            'other-month': !day.isCurrentMonth,
            'today': day.isToday,
            'selected': day.isSelected
          }]"
          @click="selectDay(day.date)"
        >
          {{ day.day }}
        </div>
      </div>
    </div>

    <div class="filter-section">
      <h3 class="section-title">筛选</h3>
      <div class="filter-item">
        <input
          type="text"
          class="filter-input"
          placeholder="搜索日程..."
          v-model="keyword"
          @input="applyFilter"
        />
      </div>
      <div class="category-list">
        <label
          v-for="cat in [{ id: null, name: '全部', color: '#6b7280' }, ...calendarStore.categories]"
          :key="cat.id ?? 'all'"
          class="category-item"
        >
          <input
            type="radio"
            :name="'category-filter'"
            :checked="calendarStore.selectedCategoryId === cat.id"
            @change="selectCategory(cat.id)"
          />
          <span class="color-dot" :style="{ background: cat.color }"></span>
          <span class="category-name">{{ cat.name }}</span>
        </label>
      </div>
    </div>

    <div class="category-section">
      <h3 class="section-title">分类管理</h3>
      <div class="category-manage-list">
        <div
          v-for="cat in calendarStore.categories"
          :key="cat.id"
          class="category-manage-item"
        >
          <span class="color-dot" :style="{ background: cat.color }"></span>
          <span class="category-name">{{ cat.name }}</span>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useCalendarStore } from '../stores/calendar'
import {
  format, startOfMonth, endOfMonth, startOfWeek, endOfWeek,
  eachDayOfInterval, isSameMonth, isSameDay, isToday,
  addMonths, subMonths
} from 'date-fns'
import { zhCN } from 'date-fns/locale'

const calendarStore = useCalendarStore()

const miniMonth = ref(startOfMonth(new Date()))
const keyword = ref('')

const weekdayLabels = computed(() => {
  return ['一', '二', '三', '四', '五', '六', '日']
})

const miniDays = computed(() => {
  const start = startOfWeek(startOfMonth(miniMonth.value), { weekStartsOn: 1 })
  const end = endOfWeek(endOfMonth(miniMonth.value), { weekStartsOn: 1 })
  const days = eachDayOfInterval({ start, end })
  return days.map(d => ({
    key: d.toISOString(),
    date: d,
    day: d.getDate(),
    isCurrentMonth: isSameMonth(d, miniMonth.value),
    isToday: isToday(d),
    isSelected: isSameDay(d, calendarStore.currentDate)
  }))
})

const miniPrev = () => {
  miniMonth.value = subMonths(miniMonth.value, 1)
}

const miniNext = () => {
  miniMonth.value = addMonths(miniMonth.value, 1)
}

const selectDay = (date) => {
  calendarStore.selectDate(date)
  miniMonth.value = startOfMonth(date)
}

const selectCategory = (id) => {
  calendarStore.setFilter(id, keyword.value)
}

const applyFilter = () => {
  calendarStore.setFilter(calendarStore.selectedCategoryId, keyword.value)
}
</script>

<style scoped>
.sidebar {
  width: 280px;
  background: #f9fafb;
  border-right: 1px solid #e5e7eb;
  padding: 20px 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.mini-calendar {
  background: white;
  border-radius: 8px;
  padding: 12px;
  border: 1px solid #e5e7eb;
}

.mini-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
}

.mini-header button {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  font-size: 16px;
  color: #6b7280;
}

.mini-header button:hover {
  background: #f3f4f6;
}

.mini-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  text-align: center;
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
}

.mini-weekdays span {
  padding: 4px 0;
}

.mini-days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.mini-day {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  border-radius: 4px;
  cursor: pointer;
  color: #1f2937;
  transition: all 0.15s;
}

.mini-day:hover {
  background: #eff6ff;
}

.mini-day.other-month {
  color: #d1d5db;
}

.mini-day.today {
  background: #dbeafe;
  font-weight: 600;
  color: #1d4ed8;
}

.mini-day.selected {
  background: #3b82f6;
  color: white;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
}

.filter-section, .category-section {
  background: white;
  border-radius: 8px;
  padding: 12px;
  border: 1px solid #e5e7eb;
}

.filter-item {
  margin-bottom: 12px;
}

.filter-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 13px;
  outline: none;
  transition: all 0.15s;
}

.filter-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.15s;
}

.category-item:hover {
  background: #f3f4f6;
}

.color-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.category-name {
  color: #374151;
}

.category-manage-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.category-manage-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  font-size: 13px;
}
</style>
