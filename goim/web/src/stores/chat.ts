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
  const unreadCounts = ref<Record<string, number>>({})

  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(localStorage.getItem('userId') || '')

  const isLoggedIn = computed(() => !!token.value && !!userId.value)

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

  async function loadMessages(targetID: string, targetType: number) {
    if (!userId.value) return
    messages.value = await api.getMessageHistory(userId.value, targetID, targetType)
    messages.value.reverse()
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
    
    const isTargetChat = (currentFriend.value && msg.from === currentFriend.value.id && msg.to_type === 0) ||
                        (currentGroup.value && msg.to === currentGroup.value.id && msg.to_type === 1)
    
    if (isTargetChat) {
      messages.value.push(message)
    } else {
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
    unreadCounts,
    token,
    userId,
    isLoggedIn,
    login,
    register,
    logout,
    loadFriends,
    loadGroups,
    loadMessages,
    selectFriend,
    selectGroup,
    addMessage,
    sendMessage,
    addFriend,
    createGroup
  }
})