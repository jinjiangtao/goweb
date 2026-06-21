import { ref, onUnmounted } from 'vue'
import { useBoardStore } from '@/store/board'
import { ElMessage } from 'element-plus'

export function useWebSocket() {
  const store = useBoardStore()
  const ws = ref(null)
  const connected = ref(false)
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5

  function connect(boardId, userName, userColor) {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      disconnect()
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/${boardId}?name=${encodeURIComponent(userName)}&color=${encodeURIComponent(userColor)}`

    try {
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        connected.value = true
        reconnectAttempts.value = 0
        console.log('WebSocket connected')
      }

      ws.value.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          handleMessage(message)
        } catch (e) {
          console.error('Message parse error:', e)
        }
      }

      ws.value.onclose = (event) => {
        connected.value = false
        console.log('WebSocket closed:', event.code, event.reason)
        
        if (reconnectAttempts.value < maxReconnectAttempts) {
          reconnectAttempts.value++
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.value), 10000)
          console.log(`Reconnecting in ${delay}ms... (${reconnectAttempts.value}/${maxReconnectAttempts})`)
          setTimeout(() => {
            connect(boardId, userName, userColor)
          }, delay)
        } else {
          ElMessage.error('连接已断开，无法自动重连')
        }
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (e) {
      console.error('WebSocket connect error:', e)
      ElMessage.error('连接服务器失败')
    }
  }

  function disconnect() {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    connected.value = false
  }

  function handleMessage(message) {
    switch (message.type) {
      case 'draw':
        handleDrawMessage(message.data)
        break
      case 'cursor':
        handleCursorMessage(message.userId, message.data)
        break
      case 'user_join':
        handleUserJoin(message.data)
        break
      case 'user_leave':
        handleUserLeave(message.data)
        break
      case 'users':
        handleUsersList(message.data)
        break
      case 'clear':
        handleClear(message.userId)
        break
      case 'sync':
        handleSync(message.data)
        break
      case 'undo':
        handleUndo(message.userId)
        break
      case 'redo':
        handleRedo(message.userId)
        break
      case 'ping':
        send({ type: 'pong', data: {} })
        break
    }
  }

  function handleDrawMessage(element) {
    if (element.userId !== store.currentUser.id) {
      const existing = store.elements.find(e => e.id === element.id)
      if (existing) {
        store.updateElement(element.id, element)
      } else {
        store.elements.push(element)
      }
    }
  }

  function handleCursorMessage(userId, point) {
    if (userId !== store.currentUser.id) {
      store.updateUserCursor(userId, point)
    }
  }

  function handleUserJoin(user) {
    if (user.id !== store.currentUser.id) {
      user.connected = true
      store.addUser(user)
      ElMessage.info(`${user.name} 加入了白板`)
    }
  }

  function handleUserLeave(user) {
    store.removeUser(user.id)
    ElMessage.info(`${user.name} 离开了白板`)
  }

  function handleUsersList(usersList) {
    const usersWithStatus = usersList.map(u => ({ ...u, connected: true }))
    store.setUsers(usersWithStatus)
  }

  function handleClear(userId) {
    if (userId !== store.currentUser.id) {
      store.elements = []
      ElMessage.info('白板已被清空')
    }
  }

  function handleSync(data) {
    if (data && data.elements) {
      store.setElements(data.elements)
    }
  }

  function handleUndo(userId) {
    if (userId !== store.currentUser.id) {
      store.undo()
    }
  }

  function handleRedo(userId) {
    if (userId !== store.currentUser.id) {
      store.redo()
    }
  }

  function send(message) {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      message.userId = store.currentUser.id
      message.timestamp = Date.now()
      ws.value.send(JSON.stringify(message))
    }
  }

  function sendDraw(element) {
    send({
      type: 'draw',
      data: element
    })
  }

  function sendCursor(point) {
    send({
      type: 'cursor',
      data: point
    })
  }

  function sendClear() {
    send({
      type: 'clear',
      data: {}
    })
  }

  function sendUndo() {
    send({
      type: 'undo',
      data: {}
    })
  }

  function sendRedo() {
    send({
      type: 'redo',
      data: {}
    })
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    connect,
    disconnect,
    sendDraw,
    sendCursor,
    sendClear,
    sendUndo,
    sendRedo
  }
}
