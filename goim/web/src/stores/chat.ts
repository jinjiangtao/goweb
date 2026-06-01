import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, Group, Message, WSMessage } from '@/types'
import * as api from '@/api'

export const useChatStore = defineStore('chat', () => {
  const currentUser = ref<User | null>(null)
  const currentFriend = ref<User | null>(null)
  const currentGroup = ref<Group | null>(null)
  const messages = ref<Message[]>([])
  const friends = ref<User[]>([])
  const groups = ref<Group[]>([])
  const onlineUsers = ref<(User & { is_friend: boolean })[]>([])
  const unreadCounts = ref<Record<string, number>>({})

  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(localStorage.getItem('userId') || '')

  const isLoggedIn = computed(() => !!token.value && !!userId.value)

  const wsSendMessage = ref<((msg: WSMessage) => Promise<boolean>) | null>(null)
  const wsIsConnected = ref(false)

  function setWsSendMessage(sendFunc: (msg: WSMessage) => Promise<boolean>) {
    wsSendMessage.value = sendFunc
  }

  function setWsConnected(connected: boolean) {
    wsIsConnected.value = connected
  }

  async function login(username: string, password: string) {
    const result = await api.login({ username, password })
    token.value = result.token
    localStorage.setItem('token', result.token)
    
    const payload = JSON.parse(atob(result.token.split('.')[1]))
    userId.value = payload.user_id
    localStorage.setItem('userId', payload.user_id)
    
    currentUser.value = {
      id: payload.user_id,
      username: payload.username,
      nickname: payload.username,
      avatar: ''
    }
    
    await loadFriends()
    await loadGroups()
    await loadOnlineUsers()
  }

  async function register(username: string, password: string, nickname: string) {
    const user = await api.register({ username, password, nickname })
    currentUser.value = user
  }

  function logout() {
    token.value = ''
    userId.value = ''
    currentUser.value = null
    currentFriend.value = null
    currentGroup.value = null
    messages.value = []
    friends.value = []
    groups.value = []
    onlineUsers.value = []
    unreadCounts.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
  }

  async function loadFriends() {
    if (!userId.value) return
    try {
      const result = await api.getFriends(userId.value)
      friends.value = result || []
    } catch {
      friends.value = []
    }
  }

  async function loadGroups() {
    if (!userId.value) return
    try {
      const result = await api.getGroups(userId.value)
      groups.value = result || []
    } catch {
      groups.value = []
    }
  }

  async function loadOnlineUsers() {
    if (!userId.value) return
    try {
      const result = await api.getOnlineUsers(userId.value)
      onlineUsers.value = result || []
      console.log('Loaded online users:', onlineUsers.value.length)
    } catch (error) {
      console.error('Failed to load online users:', error)
      onlineUsers.value = []
    }
  }

  async function refreshFriendsOnlineStatus() {
    if (!userId.value) return
    try {
      const result = await api.getFriends(userId.value)
      friends.value = result || []
      console.log('Refreshed friends online status')
    } catch (error) {
      console.error('Failed to refresh friends status:', error)
    }
  }

  async function loadMessages(targetID: string, targetType: number) {
    if (!userId.value) return
    try {
      const result = await api.getMessageHistory(userId.value, targetID, targetType)
      messages.value = result || []
      messages.value.reverse()
    } catch {
      messages.value = []
    }
  }

  function selectFriend(friend: User) {
    currentFriend.value = friend
    currentGroup.value = null
    unreadCounts.value[friend.id] = 0
    loadMessages(friend.id, 0)
  }

  function selectGroup(group: Group) {
    currentGroup.value = group
    currentFriend.value = null
    loadMessages(group.id, 1)
  }

  function addMessage(msg: WSMessage) {
    console.log('addMessage: received', msg)
    console.log('addMessage: currentFriend', currentFriend.value)
    console.log('addMessage: currentUser', currentUser.value)
    
    const message: Message = {
      id: msg.id || '',
      sender_id: msg.from,
      receiver_id: msg.to,
      receiver_type: msg.to_type,
      content: msg.content,
      type: msg.msg_type,
      status: 1,
      created_at: new Date(msg.timestamp * 1000).toISOString()
    }
    
    const isTargetChat = (currentFriend.value && 
                          ((msg.from === currentFriend.value.id && msg.to_type === 0) || 
                           (msg.to === currentFriend.value.id && msg.to === currentUser.value?.id && msg.to_type === 0))) ||
                        (currentGroup.value && msg.to === currentGroup.value.id && msg.to_type === 1)
    
    console.log('addMessage: isTargetChat', isTargetChat)
    
    if (isTargetChat) {
      console.log('addMessage: adding to messages')
      messages.value.push(message)
    } else {
      console.log('addMessage: adding to unreadCounts')
      if (msg.to_type === 0) {
        unreadCounts.value[msg.from] = (unreadCounts.value[msg.from] || 0) + 1
      } else {
        unreadCounts.value[msg.to] = (unreadCounts.value[msg.to] || 0) + 1
      }
    }
  }

  async function sendMessage(content: string, msgType: number = 0) {
    if (!currentUser.value) return

    const target = currentFriend.value || currentGroup.value
    if (!target) return

    const wsMsg: WSMessage = {
      type: 'message',
      from: currentUser.value.id,
      to: target.id,
      to_type: currentFriend.value ? 0 : 1,
      content,
      msg_type: msgType,
      timestamp: Date.now()
    }

    return wsMsg
  }

  async function addFriend(friendID: string) {
    if (!userId.value) return
    await api.addFriend(userId.value, friendID)
    await loadFriends()
    await loadOnlineUsers()
  }

  async function createGroup(name: string) {
    if (!userId.value) return
    await api.createGroup(name, userId.value)
    await loadGroups()
  }

  return {
    currentUser,
    currentFriend,
    currentGroup,
    messages,
    friends,
    groups,
    onlineUsers,
    unreadCounts,
    token,
    userId,
    isLoggedIn,
    wsSendMessage,
    wsIsConnected,
    setWsSendMessage,
    setWsConnected,
    login,
    register,
    logout,
    loadFriends,
    loadGroups,
    loadOnlineUsers,
    refreshFriendsOnlineStatus,
    loadMessages,
    selectFriend,
    selectGroup,
    addMessage,
    sendMessage,
    addFriend,
    createGroup
  }
})
