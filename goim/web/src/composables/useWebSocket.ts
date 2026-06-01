import { ref, watch, onUnmounted } from 'vue'
import type { WSMessage } from '@/types'

export function useWebSocket(userID: ref<string>) {
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const messages = ref<WSMessage[]>([])
  const messageQueue = ref<WSMessage[]>([])
  const isConnecting = ref(false)

  const connect = () => {
    if (!userID.value) {
      console.log('WebSocket: userID is empty, skipping connect')
      return
    }

    if (ws.value?.readyState === WebSocket.OPEN) {
      console.log('WebSocket: already connected')
      return
    }

    if (isConnecting.value) {
      console.log('WebSocket: already connecting')
      return
    }

    isConnecting.value = true

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/${userID.value}`
    console.log('WebSocket: connecting to', wsUrl)

    ws.value = new WebSocket(wsUrl)

    ws.value.onopen = () => {
      isConnected.value = true
      isConnecting.value = false
      console.log('WebSocket: connected')
      flushMessageQueue()
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
      isConnecting.value = false
    }

    ws.value.onclose = () => {
      isConnected.value = false
      isConnecting.value = false
      console.log('WebSocket: disconnected, reconnecting in 3s...')
      setTimeout(connect, 3000)
    }
  }

  const flushMessageQueue = () => {
    while (messageQueue.value.length > 0 && ws.value?.readyState === WebSocket.OPEN) {
      const msg = messageQueue.value.shift()
      if (msg) {
        console.log('WebSocket: sending queued message', msg)
        ws.value.send(JSON.stringify(msg))
      }
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
      isConnected.value = false
      isConnecting.value = false
    }
  }

  const sendMessage = async (message: WSMessage) => {
    console.log('WebSocket: sending message', message)
    console.log('WebSocket: readyState', ws.value?.readyState, 'OPEN:', WebSocket.OPEN)

    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(message))
      return true
    }

    if (ws.value?.readyState === WebSocket.CONNECTING || isConnecting.value) {
      console.log('WebSocket: connecting, queueing message')
      messageQueue.value.push(message)
      await waitForConnection(10000)
      return true
    }

    console.log('WebSocket: not connected, queueing message and reconnecting')
    messageQueue.value.push(message)
    connect()
    await waitForConnection(10000)
    return true
  }

  const waitForConnection = (timeout: number): Promise<void> => {
    return new Promise((resolve) => {
      const startTime = Date.now()
      const checkConnection = () => {
        if (isConnected.value && ws.value?.readyState === WebSocket.OPEN) {
          resolve()
          return
        }
        if (Date.now() - startTime > timeout) {
          console.error('WebSocket: connection timeout')
          resolve()
          return
        }
        setTimeout(checkConnection, 100)
      }
      checkConnection()
    })
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

  watch(userID, (newUserID, oldUserID) => {
    console.log('WebSocket: userID changed from', oldUserID, 'to', newUserID)
    
    if (oldUserID && oldUserID !== newUserID) {
      console.log('WebSocket: disconnecting old connection')
      disconnect()
    }
    
    if (newUserID && newUserID !== oldUserID) {
      console.log('WebSocket: initiating connection for new userID')
      setTimeout(() => {
        connect()
      }, 100)
    }
  }, { immediate: false })

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