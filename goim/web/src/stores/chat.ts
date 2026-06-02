import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, Group, Message, WSMessage, Owner } from '@/types'
import * as api from '@/api'

export const useChatStore = defineStore('chat', () => {
  const currentUser = ref<User | null>(null)
  const currentFriend = ref<User | null>(null)
  const currentGroup = ref<Group | null>(null)
  const currentOwner = ref<Owner | null>(null)
  const messages = ref<Message[]>([])
  const friends = ref<User[]>([])
  const groups = ref<Group[]>([])
  const owners = ref<Owner[]>([])
  const joinedOwners = ref<Owner[]>([])
  const onlineUsers = ref<(User & { is_friend: boolean })[]>([])
  const unreadCounts = ref<Record<string, number>>({})
  const userMap = ref<Record<string, User>>({})

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
    
    // 加载用户信息
    await loadUserProfile()
    
    await loadFriends()
    await loadGroups()
    await loadOnlineUsers()
    await loadOwners()
    await loadJoinedOwners()
  }

  async function register(username: string, password: string, nickname: string, avatar: string = '') {
    const user = await api.register({ username, password, nickname, avatar })
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

  async function loadUserProfile() {
    if (!userId.value) return
    try {
      const user = await api.getUserProfile(userId.value)
      currentUser.value = user
    } catch (error) {
      console.error('Failed to load user profile:', error)
    }
  }

  async function updateProfile(nickname: string, avatar: string) {
    if (!userId.value) return
    try {
      const user = await api.updateProfile(userId.value, nickname, avatar)
      currentUser.value = user
    } catch (error) {
      console.error('Failed to update profile:', error)
    }
  }

  function addUserToMap(user: User) {
    userMap.value[user.id] = user
  }

  async function loadFriends() {
    if (!userId.value) return
    try {
      const result = await api.getFriends(userId.value)
      friends.value = result || []
      result.forEach((friend: User) => addUserToMap(friend))
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

  async function loadOwners() {
    try {
      const result = await api.getAllOwners()
      owners.value = result || []
    } catch {
      owners.value = []
    }
  }

  async function loadJoinedOwners() {
    if (!userId.value) return
    try {
      const result = await api.getOwnersByUserID(userId.value)
      joinedOwners.value = result || []
    } catch {
      joinedOwners.value = []
    }
  }

  async function joinOwner(ownerID: string) {
    if (!userId.value) return
    await api.joinOwner(ownerID, userId.value)
    await loadJoinedOwners()
  }

  async function loadOnlineUsers() {
    if (!userId.value) return
    try {
      const result = await api.getOnlineUsers(userId.value)
      onlineUsers.value = result || []
      result.forEach((user: User) => addUserToMap(user))
      console.log('Loaded online users:', onlineUsers.value.length)
    } catch (error) {
      console.error('Failed to load online users:', error)
      onlineUsers.value = []
    }
  }

  function getNickname(userId: string): string {
    if (currentUser.value && currentUser.value.id === userId) {
      return currentUser.value.nickname
    }
    if (userMap.value[userId]) {
      return userMap.value[userId].nickname
    }
    return '未知用户'
  }

  function getAvatar(userId: string): string {
    if (currentUser.value && currentUser.value.id === userId) {
      return currentUser.value.avatar
    }
    if (userMap.value[userId]) {
      return userMap.value[userId].avatar
    }
    return ''
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
    currentOwner.value = null
    loadMessages(group.id, 1)
  }

  function selectOwner(owner: Owner) {
    currentOwner.value = owner
    currentFriend.value = null
    currentGroup.value = null
    loadMessages(owner.id, 1)
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
                           (msg.to === currentFriend.value.id && msg.to_type === 0))) ||
                        (currentGroup.value && msg.to === currentGroup.value.id && msg.to_type === 1) ||
                        (currentOwner.value && msg.to === currentOwner.value.id && msg.to_type === 1)
    
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
    currentOwner,
    messages,
    friends,
    groups,
    owners,
    joinedOwners,
    onlineUsers,
    unreadCounts,
    userMap,
    token,
    userId,
    isLoggedIn,
    login,
    register,
    logout,
    loadUserProfile,
    updateProfile,
    loadFriends,
    loadGroups,
    loadOwners,
    loadJoinedOwners,
    joinOwner,
    loadOnlineUsers,
    refreshFriendsOnlineStatus,
    loadMessages,
    selectFriend,
    selectGroup,
    selectOwner,
    addMessage,
    sendMessage,
    addFriend,
    createGroup,
    getNickname,
    getAvatar,
    addUserToMap
  }
})
