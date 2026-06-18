<template>
  <div class="app-container">
    <CalendarHeader />
    <div class="calendar-body">
      <Sidebar />
      <div class="calendar-main">
        <MonthView v-if="calendarStore.view === 'month'" />
        <WeekView v-else-if="calendarStore.view === 'week'" />
        <DayView v-else />
      </div>
    </div>
    <EventModal />
    <EventTooltip />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useCalendarStore } from './stores/calendar'
import CalendarHeader from './components/CalendarHeader.vue'
import Sidebar from './components/Sidebar.vue'
import MonthView from './components/MonthView.vue'
import WeekView from './components/WeekView.vue'
import DayView from './components/DayView.vue'
import EventModal from './components/EventModal.vue'
import EventTooltip from './components/EventTooltip.vue'

const calendarStore = useCalendarStore()

onMounted(() => {
  calendarStore.fetchCategories()
  calendarStore.fetchEvents()
})
</script>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.calendar-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.calendar-main {
  flex: 1;
  overflow: auto;
  background: white;
}
</style>
