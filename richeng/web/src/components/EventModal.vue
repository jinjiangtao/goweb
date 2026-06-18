<template>
  <Teleport to="body">
    <div v-if="calendarStore.showEventModal" class="modal-overlay" @click.self="calendarStore.closeEventModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ isEditing ? '编辑日程' : '新建日程' }}</h2>
          <button class="btn-close" @click="calendarStore.closeEventModal">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标题</label>
            <input
              type="text"
              v-model="form.title"
              class="form-input"
              placeholder="请输入日程标题"
            />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>开始时间</label>
              <input
                type="datetime-local"
                v-model="form.start_time_local"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>结束时间</label>
              <input
                type="datetime-local"
                v-model="form.end_time_local"
                class="form-input"
              />
            </div>
          </div>

          <div class="form-group">
            <label>分类</label>
            <select v-model="form.category_id" class="form-input">
              <option :value="null">无分类</option>
              <option v-for="cat in calendarStore.categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>地点</label>
            <input
              type="text"
              v-model="form.location"
              class="form-input"
              placeholder="请输入地点（可选）"
            />
          </div>

          <div class="form-group">
            <label>描述</label>
            <textarea
              v-model="form.description"
              class="form-input form-textarea"
              placeholder="请输入日程描述（可选）"
              rows="3"
            ></textarea>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="form.is_recurring" />
              <span>重复日程</span>
            </label>
          </div>

          <div v-if="form.is_recurring" class="form-row">
            <div class="form-group">
              <label>重复频率</label>
              <select v-model="form.recurring_type" class="form-input">
                <option value="daily">每天</option>
                <option value="weekly">每周</option>
                <option value="monthly">每月</option>
              </select>
            </div>
            <div class="form-group">
              <label>重复结束日期</label>
              <input
                type="date"
                v-model="form.recurring_end_local"
                class="form-input"
              />
            </div>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="form.has_reminder" />
              <span>设置提醒</span>
            </label>
          </div>

          <div v-if="form.has_reminder" class="form-group">
            <label>提醒时间</label>
            <input
              type="datetime-local"
              v-model="form.reminder_time_local"
              class="form-input"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button v-if="isEditing" class="btn btn-danger" @click="handleDelete">删除</button>
          <div class="footer-right">
            <button class="btn btn-secondary" @click="calendarStore.closeEventModal">取消</button>
            <button class="btn btn-primary" @click="handleSave">保存</button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useCalendarStore } from '../stores/calendar'
import { format, parseISO } from 'date-fns'

const calendarStore = useCalendarStore()

const form = ref({
  id: null,
  title: '',
  description: '',
  start_time: '',
  end_time: '',
  start_time_local: '',
  end_time_local: '',
  category_id: null,
  is_recurring: false,
  recurring_type: 'daily',
  recurring_end: null,
  recurring_end_local: '',
  has_reminder: false,
  reminder_time: null,
  reminder_time_local: '',
  location: ''
})

const isEditing = computed(() => !!form.value.id)

const toLocal = (iso) => {
  if (!iso) return ''
  const d = parseISO(iso)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

const toLocalDate = (iso) => {
  if (!iso) return ''
  const d = parseISO(iso)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const fromLocal = (local) => {
  if (!local) return null
  return new Date(local).toISOString()
}

const fromLocalDate = (local) => {
  if (!local) return null
  return new Date(local + 'T23:59:59').toISOString()
}

watch(
  () => calendarStore.editingEvent,
  (evt) => {
    if (evt) {
      form.value = {
        id: evt.id || null,
        title: evt.title || '',
        description: evt.description || '',
        start_time: evt.start_time || '',
        end_time: evt.end_time || '',
        start_time_local: toLocal(evt.start_time),
        end_time_local: toLocal(evt.end_time),
        category_id: evt.category_id || null,
        is_recurring: evt.is_recurring || false,
        recurring_type: evt.recurring_type || 'daily',
        recurring_end: evt.recurring_end || null,
        recurring_end_local: toLocalDate(evt.recurring_end),
        has_reminder: evt.has_reminder || false,
        reminder_time: evt.reminder_time || null,
        reminder_time_local: toLocal(evt.reminder_time),
        location: evt.location || ''
      }
    }
  },
  { immediate: true }
)

const handleSave = async () => {
  const data = {
    title: form.value.title,
    description: form.value.description,
    start_time: fromLocal(form.value.start_time_local),
    end_time: fromLocal(form.value.end_time_local),
    category_id: form.value.category_id || null,
    is_recurring: form.value.is_recurring,
    recurring_type: form.value.is_recurring ? form.value.recurring_type : '',
    recurring_end: form.value.is_recurring ? fromLocalDate(form.value.recurring_end_local) : null,
    has_reminder: form.value.has_reminder,
    reminder_time: form.value.has_reminder ? fromLocal(form.value.reminder_time_local) : null,
    location: form.value.location
  }

  try {
    if (isEditing.value) {
      await calendarStore.updateEvent(form.value.id, data)
    } else {
      await calendarStore.createEvent(data)
    }
    calendarStore.closeEventModal()
  } catch (err) {
    console.error('Save failed:', err)
    alert('保存失败，请重试')
  }
}

const handleDelete = async () => {
  if (!confirm('确定要删除这个日程吗？')) return
  try {
    await calendarStore.deleteEvent(form.value.id)
    calendarStore.closeEventModal()
  } catch (err) {
    console.error('Delete failed:', err)
    alert('删除失败，请重试')
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 12px;
  width: 520px;
  max-width: 90vw;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.btn-close {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  font-size: 24px;
  color: #6b7280;
  line-height: 1;
  transition: all 0.15s;
}

.btn-close:hover {
  background: #f3f4f6;
}

.modal-body {
  padding: 20px 24px;
  overflow-y: auto;
  flex: 1;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.form-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  color: #1f2937;
  outline: none;
  transition: all 0.15s;
}

.form-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.form-textarea {
  resize: vertical;
  min-height: 72px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #374151;
}

.checkbox-label input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  border-top: 1px solid #e5e7eb;
}

.footer-right {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.15s;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: #f3f4f6;
  color: #4b5563;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.btn-danger {
  background: #fee2e2;
  color: #dc2626;
}

.btn-danger:hover {
  background: #fecaca;
}
</style>
