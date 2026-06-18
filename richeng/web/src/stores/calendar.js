import { defineStore } from 'pinia'
import { eventApi, categoryApi } from '../api'
import {
  startOfMonth, endOfMonth,
  startOfWeek, endOfWeek,
  startOfDay, endOfDay,
  addMonths, addWeeks, addDays,
  subMonths, subWeeks, subDays,
  format
} from 'date-fns'
import { zhCN } from 'date-fns/locale'

export const useCalendarStore = defineStore('calendar', {
  state: () => ({
    view: 'month',
    currentDate: new Date(),
    events: [],
    categories: [],
    selectedCategoryId: null,
    keyword: '',
    hoveredEvent: null,
    editingEvent: null,
    showEventModal: false
  }),

  getters: {
    viewTitle: (state) => {
      const d = state.currentDate
      if (state.view === 'month') {
        return format(d, 'yyyy年MM月', { locale: zhCN })
      } else if (state.view === 'week') {
        const start = startOfWeek(d, { weekStartsOn: 1 })
        const end = endOfWeek(d, { weekStartsOn: 1 })
        return `${format(start, 'yyyy年MM月dd日', { locale: zhCN })} - ${format(end, 'MM月dd日', { locale: zhCN })}`
      } else {
        return format(d, 'yyyy年MM月dd日 EEEE', { locale: zhCN })
      }
    },

    dateRange: (state) => {
      const d = state.currentDate
      if (state.view === 'month') {
        return {
          start: startOfWeek(startOfMonth(d), { weekStartsOn: 1 }).toISOString(),
          end: endOfWeek(endOfMonth(d), { weekStartsOn: 1 }).toISOString()
        }
      } else if (state.view === 'week') {
        return {
          start: startOfWeek(d, { weekStartsOn: 1 }).toISOString(),
          end: endOfWeek(d, { weekStartsOn: 1 }).toISOString()
        }
      } else {
        return {
          start: startOfDay(d).toISOString(),
          end: endOfDay(d).toISOString()
        }
      }
    },

    filteredEvents: (state) => {
      return state.events
    },

    categoryMap: (state) => {
      const map = {}
      state.categories.forEach(c => {
        map[c.id] = c
      })
      return map
    }
  },

  actions: {
    setView(view) {
      this.view = view
      this.fetchEvents()
    },

    navigatePrev() {
      if (this.view === 'month') {
        this.currentDate = subMonths(this.currentDate, 1)
      } else if (this.view === 'week') {
        this.currentDate = subWeeks(this.currentDate, 1)
      } else {
        this.currentDate = subDays(this.currentDate, 1)
      }
      this.fetchEvents()
    },

    navigateNext() {
      if (this.view === 'month') {
        this.currentDate = addMonths(this.currentDate, 1)
      } else if (this.view === 'week') {
        this.currentDate = addWeeks(this.currentDate, 1)
      } else {
        this.currentDate = addDays(this.currentDate, 1)
      }
      this.fetchEvents()
    },

    goToday() {
      this.currentDate = new Date()
      this.fetchEvents()
    },

    selectDate(date) {
      this.currentDate = date
      if (this.view !== 'day') {
        this.view = 'day'
      }
      this.fetchEvents()
    },

    async fetchCategories() {
      try {
        const res = await categoryApi.getCategories()
        this.categories = res.data.data || []
      } catch (err) {
        console.error('Failed to fetch categories:', err)
      }
    },

    async fetchEvents() {
      try {
        const params = {
          start_date: this.dateRange.start,
          end_date: this.dateRange.end
        }
        if (this.selectedCategoryId) {
          params.category_id = this.selectedCategoryId
        }
        if (this.keyword) {
          params.keyword = this.keyword
        }
        const res = await eventApi.getEvents(params)
        this.events = res.data.data || []
      } catch (err) {
        console.error('Failed to fetch events:', err)
      }
    },

    async createEvent(data) {
      try {
        const res = await eventApi.createEvent(data)
        await this.fetchEvents()
        return res.data.data
      } catch (err) {
        console.error('Failed to create event:', err)
        throw err
      }
    },

    async updateEvent(id, data) {
      try {
        const res = await eventApi.updateEvent(id, data)
        await this.fetchEvents()
        return res.data.data
      } catch (err) {
        console.error('Failed to update event:', err)
        throw err
      }
    },

    async deleteEvent(id) {
      try {
        await eventApi.deleteEvent(id)
        await this.fetchEvents()
      } catch (err) {
        console.error('Failed to delete event:', err)
        throw err
      }
    },

    setFilter(categoryId, keyword) {
      this.selectedCategoryId = categoryId || null
      this.keyword = keyword || ''
      this.fetchEvents()
    },

    openNewEventModal(defaultData = {}) {
      this.editingEvent = { ...defaultData }
      this.showEventModal = true
    },

    openEditEventModal(event) {
      this.editingEvent = { ...event }
      this.showEventModal = true
    },

    closeEventModal() {
      this.editingEvent = null
      this.showEventModal = false
    },

    setHoveredEvent(event) {
      this.hoveredEvent = event
    }
  }
})
