<template>
  <div
    ref="cardRef"
    :class="['event-card', view, { expired: isExpired, 'is-dragging': isDragging }]"
    :style="cardStyle"
    draggable="true"
    @click.stop="handleClick"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    @dragstart="handleDragStart"
    @dragend="handleDragEnd"
  >
    <div class="event-color-bar" :style="{ background: categoryColor }"></div>
    <div class="event-content">
      <div v-if="view !== 'month'" class="event-time">{{ timeLabel }}</div>
      <div class="event-title">{{ event.title || '(无标题)' }}</div>
      <div v-if="view === 'month' && event.location" class="event-location">{{ event.location }}</div>
    </div>
    <div
      v-if="view !== 'month'"
      class="resize-handle"
      draggable="true"
      @dragstart.stop="handleResizeStart"
      @dragend.stop="handleResizeEnd"
      @click.stop
    ></div>
  </div>
</template>

<script setup>import { ref, computed } from 'vue';
import { useCalendarStore } from '../stores/calendar';
import { format, differenceInMinutes, isBefore, startOfDay, differenceInCalendarDays } from 'date-fns';
const props = defineProps({
 event: {
 type: Object,
 required: true
 },
 view: {
 type: String,
 default: 'month'
 },
 dayDate: {
 type: Date,
 default: null
 }
});
const calendarStore = useCalendarStore();
const cardRef = ref(null);
const isDragging = ref(false);
const categoryColor = computed(() => {
 if (props.event.category) {
 return props.event.category.color;
 }
 if (props.event.category_id && calendarStore.categoryMap[props.event.category_id]) {
 return calendarStore.categoryMap[props.event.category_id].color;
 }
 return '#6b7280';
});
const isExpired = computed(() => {
 return isBefore(new Date(props.event.end_time), new Date());
});
const timeLabel = computed(() => {
 const start = new Date(props.event.start_time);
 const end = new Date(props.event.end_time);
 return `${format(start, 'HH:mm')} - ${format(end, 'HH:mm')}`;
});
const cardStyle = computed(() => {
 if (props.view === 'month') {
 return {
 background: categoryColor.value + '20',
 borderLeft: `3px solid ${categoryColor.value}`
 };
 }
 const dayStart = props.dayDate ? startOfDay(props.dayDate) : startOfDay(new Date(props.event.start_time));
 const startMinutes = differenceInMinutes(new Date(props.event.start_time), dayStart);
 const duration = differenceInMinutes(new Date(props.event.end_time), new Date(props.event.start_time));
 const top = (startMinutes / 60) * 48;
 const height = Math.max((duration / 60) * 48, 24);
 return {
 position: 'absolute',
 top: `${top}px`,
 left: '4px',
 right: '4px',
 height: `${height}px`,
 background: categoryColor.value,
 pointerEvents: 'auto',
 zIndex: isDragging.value ? 100 : 10
 };
});
const handleClick = () => {
 calendarStore.openEditEventModal(props.event);
};
const handleMouseEnter = (e) => {
 calendarStore.setHoveredEvent({ event: props.event, x: e.clientX, y: e.clientY });
};
const handleMouseLeave = () => {
 calendarStore.setHoveredEvent(null);
};
const handleDragStart = (e) => {
 if (props.view === 'month')
 return;
 isDragging.value = true;
 e.dataTransfer.effectAllowed = 'move';
 e.dataTransfer.setData('application/json', JSON.stringify({
 type: 'move',
 eventId: props.event.id,
 originalStart: props.event.start_time,
 originalEnd: props.event.end_time
 }));
};
const handleDragEnd = () => {
 isDragging.value = false;
};
const handleResizeStart = (e) => {
 isDragging.value = true;
 e.dataTransfer.effectAllowed = 'link';
 e.dataTransfer.setData('application/json', JSON.stringify({
 type: 'resize',
 eventId: props.event.id,
 originalStart: props.event.start_time,
 originalEnd: props.event.end_time
 }));
};
const handleResizeEnd = () => {
 isDragging.value = false;
};
</script>

<style scoped>
.event-card {
  display: flex;
  align-items: stretch;
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.15s;
  user-select: none;
  min-height: 22px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.event-card:hover {
  filter: brightness(0.95);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  z-index: 20;
}

.event-card.month {
  padding: 2px 6px;
  font-size: 12px;
  align-items: center;
}

.event-card.expired {
  opacity: 0.55;
}

.event-card.is-dragging {
  opacity: 0.7;
}

.event-color-bar {
  width: 4px;
  flex-shrink: 0;
}

.event-card.month .event-color-bar {
  width: 0;
}

.event-content {
  flex: 1;
  min-width: 0;
  padding: 2px 4px;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.event-card.month .event-content {
  padding: 0;
  flex-direction: row;
  align-items: center;
  gap: 6px;
  color: #374151;
}

.event-card:not(.month) .event-content {
  color: white;
}

.event-time {
  font-size: 11px;
  font-weight: 500;
  opacity: 0.9;
  white-space: nowrap;
}

.event-title {
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.event-card.month .event-title {
  font-weight: 500;
}

.event-location {
  font-size: 11px;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.resize-handle {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 6px;
  cursor: ns-resize;
  background: transparent;
}

.resize-handle:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
