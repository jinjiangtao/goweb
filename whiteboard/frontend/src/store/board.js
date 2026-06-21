import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { generateId, generateColor } from '@/utils/common'

export const useBoardStore = defineStore('board', () => {
  const elements = ref([])
  const users = ref([])
  const currentTool = ref('pen')
  const currentStyle = ref({
    color: '#409eff',
    width: 3,
    opacity: 1,
    fill: 'transparent',
    fillOpacity: 0
  })
  const currentUser = ref({
    id: generateId(),
    name: '用户' + Math.floor(Math.random() * 1000),
    avatar: '',
    color: generateColor(),
    connected: true
  })
  const boardId = ref('')
  const boardName = ref('')
  const background = ref('#ffffff')
  const zoom = ref(1)
  const pan = ref({ x: 0, y: 0 })
  const isDrawing = ref(false)
  const currentElement = ref(null)
  const historyStack = ref([])
  const redoStack = ref([])
  const selectedLayerId = ref(null)

  const onlineUsers = computed(() => users.value.filter(u => u.connected))

  function setBoardId(id) {
    boardId.value = id
  }

  function setBoardName(name) {
    boardName.value = name
  }

  function setBackground(color) {
    background.value = color
  }

  function setTool(tool) {
    currentTool.value = tool
  }

  function setStyle(style) {
    currentStyle.value = { ...currentStyle.value, ...style }
  }

  function setZoom(value) {
    zoom.value = Math.max(0.1, Math.min(5, value))
  }

  function setPan(x, y) {
    pan.value = { x, y }
  }

  function addElement(element) {
    saveHistory()
    elements.value.push(element)
  }

  function updateElement(id, updates) {
    const index = elements.value.findIndex(e => e.id === id)
    if (index !== -1) {
      elements.value[index] = { ...elements.value[index], ...updates }
    }
  }

  function removeElement(id) {
    saveHistory()
    elements.value = elements.value.filter(e => e.id !== id)
  }

  function clearElements() {
    saveHistory()
    elements.value = []
  }

  function setElements(newElements) {
    elements.value = newElements
  }

  function saveHistory() {
    historyStack.value.push(JSON.parse(JSON.stringify(elements.value)))
    redoStack.value = []
    if (historyStack.value.length > 50) {
      historyStack.value.shift()
    }
  }

  function undo() {
    if (historyStack.value.length > 0) {
      redoStack.value.push(JSON.parse(JSON.stringify(elements.value)))
      elements.value = historyStack.value.pop() || []
      return true
    }
    return false
  }

  function redo() {
    if (redoStack.value.length > 0) {
      historyStack.value.push(JSON.parse(JSON.stringify(elements.value)))
      elements.value = redoStack.value.pop() || []
      return true
    }
    return false
  }

  function addUser(user) {
    const existing = users.value.find(u => u.id === user.id)
    if (!existing) {
      users.value.push(user)
    } else {
      Object.assign(existing, user)
    }
  }

  function removeUser(userId) {
    const user = users.value.find(u => u.id === userId)
    if (user) {
      user.connected = false
    }
  }

  function updateUserCursor(userId, cursor) {
    const user = users.value.find(u => u.id === userId)
    if (user) {
      user.cursor = cursor
    }
  }

  function setUsers(newUsers) {
    users.value = newUsers
  }

  function selectLayer(id) {
    selectedLayerId.value = id === selectedLayerId.value ? null : id
  }

  function moveLayerUp(id) {
    const index = elements.value.findIndex(e => e.id === id)
    if (index !== -1 && index < elements.value.length - 1) {
      saveHistory()
      const temp = elements.value[index]
      elements.value[index] = elements.value[index + 1]
      elements.value[index + 1] = temp
    }
  }

  function moveLayerDown(id) {
    const index = elements.value.findIndex(e => e.id === id)
    if (index > 0) {
      saveHistory()
      const temp = elements.value[index]
      elements.value[index] = elements.value[index - 1]
      elements.value[index - 1] = temp
    }
  }

  return {
    elements,
    users,
    currentTool,
    currentStyle,
    currentUser,
    boardId,
    boardName,
    background,
    zoom,
    pan,
    isDrawing,
    currentElement,
    selectedLayerId,
    onlineUsers,
    setBoardId,
    setBoardName,
    setBackground,
    setTool,
    setStyle,
    setZoom,
    setPan,
    addElement,
    updateElement,
    removeElement,
    clearElements,
    setElements,
    undo,
    redo,
    addUser,
    removeUser,
    updateUserCursor,
    setUsers,
    selectLayer,
    moveLayerUp,
    moveLayerDown
  }
})
