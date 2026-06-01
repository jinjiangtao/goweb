import axios from 'axios'
import type { LoginRequest, RegisterRequest, User, Group, Message } from '@/types'

const API_BASE_URL = '/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
})

export async function login(data: LoginRequest): Promise<{ token: string }> {
  const response = await api.post('/auth/login', data)
  return response.data
}

export async function register(data: RegisterRequest): Promise<User> {
  const response = await api.post('/auth/register', data)
  return response.data
}

export async function getFriends(userID: string): Promise<User[]> {
  const response = await api.get(`/user/friends?user_id=${userID}`)
  return response.data
}

export async function addFriend(userID: string, friendID: string): Promise<void> {
  await api.post('/user/add-friend', { user_id: userID, friend_id: friendID })
}

export async function getOnlineStatus(userID: string): Promise<{ online: boolean }> {
  const response = await api.get(`/user/online-status?user_id=${userID}`)
  return response.data
}

export interface OnlineUser extends User {
  is_friend: boolean
}

export async function getOnlineUsers(userID: string): Promise<OnlineUser[]> {
  const response = await api.get(`/user/online-users?user_id=${userID}`)
  return response.data
}

export async function createGroup(name: string, ownerID: string): Promise<Group> {
  const response = await api.post('/group/create', { name, owner_id: ownerID })
  return response.data
}

export async function getGroups(userID: string): Promise<Group[]> {
  const response = await api.get(`/group/list?user_id=${userID}`)
  return response.data
}

export async function addGroupMember(groupID: string, userID: string): Promise<void> {
  await api.post('/group/add-member', { group_id: groupID, user_id: userID })
}

export async function removeGroupMember(groupID: string, userID: string): Promise<void> {
  await api.post('/group/remove-member', { group_id: groupID, user_id: userID })
}

export async function getGroupMembers(groupID: string): Promise<{ user_id: string; nickname: string }[]> {
  const response = await api.get(`/group/members/${groupID}`)
  return response.data
}

export async function getMessageHistory(
  userID: string,
  targetID: string,
  targetType: number,
  limit = 20,
  offset = 0
): Promise<Message[]> {
  const response = await api.get('/message/history', {
    params: { user_id: userID, target_id: targetID, target_type: targetType, limit, offset }
  })
  return response.data
}

export async function markMessageRead(userID: string, messageID: string): Promise<void> {
  await api.post('/message/read', { user_id: userID, message_id: messageID })
}

export async function getUnreadCounts(userID: string): Promise<Record<string, number>> {
  const response = await api.get(`/message/unread-counts?user_id=${userID}`)
  return response.data
}