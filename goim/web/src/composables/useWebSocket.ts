import { ref, watch, onUnmounted } from 'vue'
import type { WSMessage } from '@/types'

export function useWebSocket(userID: ref<string>) {
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const messages = ref<WSMessage[]>([])

  const connect = () => {
    if (!userID.value) {
      console.log('WebSocket: userID is empty, skipping connect')
      return
    }

    if (ws.value?.readyState === WebSocket.OPEN) {
      console.log('WebSocket: already connected')
      return
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/${userID.value}`
    console.log('WebSocket: connecting to', wsUrl)

    ws.value = new WebSocket(wsUrl)

    ws.value.onopen = () => {
      isConnected.value = true
      console.log('WebSocket: connected')
    }

    ws.value.onmessage = (event) => {
      try {
        console.log('WebSocket: received message', event.data)
        const message: WSMessage = JSON.parse(event.data)
        if (message.type !== 'pong') {
          messages.value.push(message)
          showNotification(message)
        }
      } catch (error) {
        console.error('WebSocket: failed to parse message', error)
      }
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket: error', error)
    }

    ws.value.onclose = () => {
      isConnected.value = false
      console.log('WebSocket: disconnected, reconnecting in 3s...')
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
    console.log('WebSocket: sending message', message)
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(message))
    } else {
      console.error('WebSocket: not connected')
    }
  }

  const sendPing = () => {
    sendMessage({ type: 'ping', from: userID.value, to: '', to_type: 0, content: '', msg_type: 0, timestamp: Date.now() })
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

  watch(userID, (newUserID) => {
    console.log('WebSocket: userID changed to', newUserID)
    disconnect()
    if (newUserID) {
      connect()
    }
  }, { immediate: true })

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