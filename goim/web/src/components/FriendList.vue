<script setup lang="ts">
import { ref, computed } from 'vue'
import { useChatStore } from '@/stores/chat'
import type { User, Group } from '@/types'

const chatStore = useChatStore()
const searchQuery = ref('')

const filteredFriends = computed(() => {
  if (!searchQuery.value) return chatStore.friends
  const query = searchQuery.value.toLowerCase()
  return chatStore.friends.filter(f => 
    f.nickname.toLowerCase().includes(query) || 
    f.username.toLowerCase().includes(query)
  )
})

const filteredGroups = computed(() => {
  if (!searchQuery.value) return chatStore.groups
  const query = searchQuery.value.toLowerCase()
  return chatStore.groups.filter(g => g.name.toLowerCase().includes(query))
})

const filteredOnlineUsers = computed(() => {
  const nonFriends = chatStore.onlineUsers.filter(u => !u.is_friend)
  if (!searchQuery.value) return nonFriends
  const query = searchQuery.value.toLowerCase()
  return nonFriends.filter(u => 
    u.nickname.toLowerCase().includes(query) || 
    u.username.toLowerCase().includes(query)
  )
})

function getInitials(name: string): string {
  return name.charAt(0).toUpperCase()
}

function selectFriend(friend: User) {
  chatStore.selectFriend(friend)
}

function selectGroup(group: Group) {
  chatStore.selectGroup(group)
}

async function handleAddFriend(user: User & { is_friend: boolean }) {
  await chatStore.addFriend(user.id)
}
</script>

<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <h3>GoIM</h3>
      <button class="logout-btn" @click="chatStore.logout">退出</button>
    </div>
    
    <div class="search-bar">
      <input v-model="searchQuery" type="text" placeholder="搜索好友或用户..." />
    </div>
    
    <div class="friend-list">
      <div v-if="filteredOnlineUsers.length > 0">
        <div class="group-list">
          <h5>在线用户</h5>
        </div>
        <div 
          v-for="user in filteredOnlineUsers" 
          :key="user.id"
          class="friend-item"
        >
          <div class="avatar">{{ getInitials(user.nickname) }}</div>
          <div class="friend-info">
            <h4>{{ user.nickname }}</h4>
            <p>{{ user.username }}</p>
          </div>
          <div class="friend-status">
            <div class="online-dot"></div>
            <button 
              class="add-friend-btn"
              @click="handleAddFriend(user)"
            >
              +加好友
            </button>
          </div>
        </div>
      </div>

      <div v-if="filteredFriends.length > 0">
        <div class="group-list">
          <h5>好友</h5>
        </div>
        <div 
          v-for="friend in filteredFriends" 
          :key="friend.id"
          class="friend-item"
          :class="{ active: chatStore.currentFriend?.id === friend.id }"
          @click="selectFriend(friend)"
        >
          <div class="avatar">{{ getInitials(friend.nickname) }}</div>
          <div class="friend-info">
            <h4>{{ friend.nickname }}</h4>
            <p>{{ friend.username }}</p>
          </div>
          <div class="friend-status">
            <div :class="['online-dot', { 'offline-dot': !friend.online }]"></div>
            <span v-if="chatStore.unreadCounts[friend.id]" class="unread-badge">
              {{ chatStore.unreadCounts[friend.id] }}
            </span>
          </div>
        </div>
      </div>
      
      <div v-if="filteredGroups.length > 0">
        <div class="group-list">
          <h5>群组</h5>
        </div>
        <div 
          v-for="group in filteredGroups" 
          :key="group.id"
          class="friend-item"
          :class="{ active: chatStore.currentGroup?.id === group.id }"
          @click="selectGroup(group)"
        >
          <div class="avatar">{{ getInitials(group.name) }}</div>
          <div class="friend-info">
            <h4>{{ group.name }}</h4>
            <p>群组</p>
          </div>
          <div class="friend-status">
            <span v-if="chatStore.unreadCounts[group.id]" class="unread-badge">
              {{ chatStore.unreadCounts[group.id] }}
            </span>
          </div>
        </div>
      </div>
      
      <div v-if="filteredOnlineUsers.length === 0 && filteredFriends.length === 0 && filteredGroups.length === 0" class="empty-chat">
        {{ searchQuery ? '未找到匹配的联系人' : '暂无联系人' }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.add-friend-btn {
  background: #667eea;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 11px;
  cursor: pointer;
  margin-top: 4px;
}

.add-friend-btn:hover {
  background: #764ba2;
}
</style>