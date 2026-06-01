import { ref, onMounted, onUnmounted } from 'vue'
import type { WSMessage } from '@/types'

export function useWebSocket(userID: string) {
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const messages = ref<WSMessage[]>([])

  const connect = () => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      return
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/${userID}`
    
    ws.value = new WebSocket(wsUrl)

    ws.value.onopen = () => {
      isConnected.value = true
      console.log('WebSocket connected')
    }

    ws.value.onmessage = (event) => {
      try {
        const message: WSMessage = JSON.parse(event.data)
        if (message.type !== 'pong') {
          messages.value.push(message)
          showNotification(message)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.value.onclose = () => {
      isConnected.value = false
      console.log('WebSocket disconnected, reconnecting...')
      setTimeout(connect, 3000)
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
      isConnected.value = false
    }
  }

  const sendMessage = (message: WSMessage) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(message))
    }
  }

  const sendPing = () => {
    sendMessage({ type: 'ping', from: userID, to: '', to_type: 0, content: '', msg_type: 0, timestamp: Date.now() })
  }

  const showNotification = (message: WSMessage) => {
    if (!('Notification' in window)) {
      return
    }

    if (Notification.permission === 'granted') {
      new Notification('新消息', {
        body: message.content,
        icon: '/vite.svg'
      })
    } else if (Notification.permission !== 'denied') {
      Notification.requestPermission().then(permission => {
        if (permission === 'granted') {
          new Notification('新消息', {
            body: message.content,
            icon: '/vite.svg'
          })
        }
      })
    }
  }

  let pingInterval: ReturnType<typeof setInterval> | null = null

  onMounted(() => {
    connect()
    pingInterval = setInterval(sendPing, 30000)
  })

  onUnmounted(() => {
    disconnect()
    if (pingInterval) {
      clearInterval(pingInterval)
    }
  })

  return {
    ws,
    isConnected,
    messages,
    sendMessage,
    connect,
    disconnect
  }
}