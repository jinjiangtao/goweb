export interface User {
  id: string
  username: string
  nickname: string
  avatar: string
  online?: boolean
}

export interface Friend {
  id: string
  user_id: string
  friend_id: string
  status: number
  created_at: string
}

export interface Group {
  id: string
  name: string
  avatar: string
  owner_id: string
}

export interface Owner {
  id: string
  name: string
  description: string
  avatar: string
  created_at: string
}

export interface OwnerMember {
  id: string
  owner_id: string
  user_id: string
  username: string
  nickname: string
  joined_at: string
}

export interface Message {
  id: string
  sender_id: string
  receiver_id: string
  receiver_type: number
  content: string
  type: number
  status: number
  created_at: string
}

export interface WSMessage {
  id?: string
  type: string
  from: string
  to: string
  to_type: number
  content: string
  msg_type: number
  timestamp: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  nickname: string
  avatar?: string
}

export interface ChatState {
  currentUser: User | null
  currentFriend: User | null
  currentGroup: Group | null
  messages: Message[]
  friends: User[]
  groups: Group[]
  unreadCounts: Record<string, number>
}