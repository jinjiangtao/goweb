import type { WSMessage } from '@/types'

class WebSocketManager {
  private ws: WebSocket | null = null
  private messageQueue: WSMessage[] = []
  private isConnecting = false
  private isConnected = false
  private messageHandlers: ((msg: WSMessage) => void)[] = []
  private userId: string = ''
  private reconnectTimeout: ReturnType<typeof setTimeout> | null = null

  connect(userID: string) {
    if (!userID) {
      console.log('WebSocketManager: userID is empty, skipping connect')
      return
    }

    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('WebSocketManager: already connected')
      return
    }

    if (this.isConnecting) {
      console.log('WebSocketManager: already connecting')
      return
    }

    this.userId = userID
    this.isConnecting = true

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/${userID}`
    console.log('WebSocketManager: connecting to', wsUrl)

    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      this.isConnected = true
      this.isConnecting = false
      console.log('WebSocketManager: connected')
      this.flushMessageQueue()
    }

    this.ws.onmessage = (event) => {
      try {
        console.log('WebSocketManager: received message', event.data)
        const message: WSMessage = JSON.parse(event.data)
        if (message.type !== 'pong') {
          this.messageHandlers.forEach(handler => handler(message))
        }
      } catch (error) {
        console.error('WebSocketManager: failed to parse message', error)
      }
    }

    this.ws.onerror = (error) => {
      console.error('WebSocketManager: error', error)
      this.isConnecting = false
    }

    this.ws.onclose = () => {
      this.isConnected = false
      this.isConnecting = false
      console.log('WebSocketManager: disconnected, reconnecting in 3s...')
      this.scheduleReconnect()
    }
  }

  private scheduleReconnect() {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
    }
    this.reconnectTimeout = setTimeout(() => {
      if (this.userId) {
        this.connect(this.userId)
      }
    }, 3000)
  }

  disconnect() {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
    }
    if (this.ws) {
      this.ws.close()
      this.ws = null
      this.isConnected = false
      this.isConnecting = false
    }
  }

  flushMessageQueue() {
    while (this.messageQueue.length > 0 && this.ws?.readyState === WebSocket.OPEN) {
      const msg = this.messageQueue.shift()
      if (msg) {
        console.log('WebSocketManager: sending queued message', msg)
        this.ws.send(JSON.stringify(msg))
      }
    }
  }

  async sendMessage(message: WSMessage): Promise<boolean> {
    console.log('WebSocketManager: sending message', message)
    console.log('WebSocketManager: readyState', this.ws?.readyState, 'OPEN:', WebSocket.OPEN)

    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
      return true
    }

    if (this.ws?.readyState === WebSocket.CONNECTING || this.isConnecting) {
      console.log('WebSocketManager: connecting, queueing message')
      this.messageQueue.push(message)
      await this.waitForConnection(10000)
      return true
    }

    console.log('WebSocketManager: not connected, queueing message and reconnecting')
    this.messageQueue.push(message)
    if (this.userId) {
      this.connect(this.userId)
    }
    await this.waitForConnection(10000)
    return true
  }

  private waitForConnection(timeout: number): Promise<void> {
    return new Promise((resolve) => {
      const startTime = Date.now()
      const checkConnection = () => {
        if (this.isConnected && this.ws?.readyState === WebSocket.OPEN) {
          resolve()
          return
        }
        if (Date.now() - startTime > timeout) {
          console.error('WebSocketManager: connection timeout')
          resolve()
          return
        }
        setTimeout(checkConnection, 100)
      }
      checkConnection()
    })
  }

  onMessage(handler: (msg: WSMessage) => void) {
    this.messageHandlers.push(handler)
  }

  removeMessageHandler(handler: (msg: WSMessage) => void) {
    const index = this.messageHandlers.indexOf(handler)
    if (index !== -1) {
      this.messageHandlers.splice(index, 1)
    }
  }

  get connected(): boolean {
    return this.isConnected
  }
}

const wsManager = new WebSocketManager()

export function useWebSocket() {
  return {
    sendMessage: (msg: WSMessage) => wsManager.sendMessage(msg),
    connect: (userID: string) => wsManager.connect(userID),
    disconnect: () => wsManager.disconnect(),
    onMessage: (handler: (msg: WSMessage) => void) => wsManager.onMessage(handler),
    removeMessageHandler: (handler: (msg: WSMessage) => void) => wsManager.removeMessageHandler(handler),
    connected: () => wsManager.connected
  }
}

export { wsManager }
