<script setup lang="ts">
import { ref, computed } from 'vue'
import { useChatStore } from '@/stores/chat'
import type { User, Group, Owner } from '@/types'
import ProfileEdit from './ProfileEdit.vue'

const chatStore = useChatStore()
const searchQuery = ref('')
const showProfileEdit = ref(false)

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

function selectOwner(owner: Owner) {
  chatStore.selectOwner(owner)
}

async function handleAddFriend(user: User & { is_friend: boolean }) {
  await chatStore.addFriend(user.id)
}

async function handleJoinOwner(owner: Owner) {
  await chatStore.joinOwner(owner.id)
}

const filteredOwners = computed(() => {
  const unjoinedOwners = chatStore.owners.filter(o => 
    !chatStore.joinedOwners.some(jo => jo.id === o.id)
  )
  if (!searchQuery.value) return unjoinedOwners
  const query = searchQuery.value.toLowerCase()
  return unjoinedOwners.filter(o => 
    o.name.toLowerCase().includes(query) || 
    o.description.toLowerCase().includes(query)
  )
})

const filteredJoinedOwners = computed(() => {
  if (!searchQuery.value) return chatStore.joinedOwners
  const query = searchQuery.value.toLowerCase()
  return chatStore.joinedOwners.filter(o => 
    o.name.toLowerCase().includes(query) || 
    o.description.toLowerCase().includes(query)
  )
})
</script>

<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <div class="user-info" @click="showProfileEdit = true">
        <div class="avatar">
          <img v-if="chatStore.currentUser?.avatar" :src="chatStore.currentUser.avatar" alt="头像" />
          <span v-else>{{ getInitials(chatStore.currentUser?.nickname || 'U') }}</span>
        </div>
        <div>
          <h4>{{ chatStore.currentUser?.nickname }}</h4>
          <p>{{ chatStore.currentUser?.username }}</p>
        </div>
      </div>
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
          <div class="avatar">
            <img v-if="user.avatar" :src="user.avatar" alt="头像" />
            <span v-else>{{ getInitials(user.nickname) }}</span>
          </div>
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
          <div class="avatar">
            <img v-if="friend.avatar" :src="friend.avatar" alt="头像" />
            <span v-else>{{ getInitials(friend.nickname) }}</span>
          </div>
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
          <div class="avatar">
            <img v-if="group.avatar" :src="group.avatar" alt="头像" />
            <span v-else>{{ getInitials(group.name) }}</span>
          </div>
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

      <div v-if="filteredJoinedOwners.length > 0">
        <div class="group-list">
          <h5>群主群聊</h5>
        </div>
        <div 
          v-for="owner in filteredJoinedOwners" 
          :key="owner.id"
          class="friend-item owner-item"
          :class="{ active: chatStore.currentOwner?.id === owner.id }"
          @click="selectOwner(owner)"
        >
          <div class="avatar">
            <img v-if="owner.avatar" :src="owner.avatar" alt="头像" />
            <span v-else>{{ getInitials(owner.name) }}</span>
          </div>
          <div class="friend-info">
            <div class="owner-name">
              <h4>{{ owner.name }}</h4>
              <span class="official-badge">官方</span>
            </div>
            <p>{{ owner.description || '群主' }}</p>
          </div>
          <div class="friend-status">
            <span v-if="chatStore.unreadCounts[owner.id]" class="unread-badge">
              {{ chatStore.unreadCounts[owner.id] }}
            </span>
          </div>
        </div>
      </div>

      <div v-if="filteredOwners.length > 0">
        <div class="group-list">
          <h5>官方群主</h5>
        </div>
        <div 
          v-for="owner in filteredOwners" 
          :key="owner.id"
          class="friend-item owner-item"
        >
          <div class="avatar">
            <img v-if="owner.avatar" :src="owner.avatar" alt="头像" />
            <span v-else>{{ getInitials(owner.name) }}</span>
          </div>
          <div class="friend-info">
            <div class="owner-name">
              <h4>{{ owner.name }}</h4>
              <span class="official-badge">官方</span>
            </div>
            <p>{{ owner.description || '群主' }}</p>
          </div>
          <div class="friend-status">
            <button 
              class="join-owner-btn"
              @click="handleJoinOwner(owner)"
            >
              加入
            </button>
          </div>
        </div>
      </div>
      
      <div v-if="filteredOnlineUsers.length === 0 && filteredFriends.length === 0 && filteredGroups.length === 0 && filteredJoinedOwners.length === 0 && filteredOwners.length === 0" class="empty-chat">
        {{ searchQuery ? '未找到匹配的联系人' : '暂无联系人' }}
      </div>
    </div>
    <ProfileEdit :visible="showProfileEdit" @close="showProfileEdit = false" />
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

.join-owner-btn {
  background: #10b981;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 11px;
  cursor: pointer;
  margin-top: 4px;
}

.join-owner-btn:hover {
  background: #059669;
}

.user-info {
  cursor: pointer;
}

.owner-item {
  background: linear-gradient(90deg, rgba(59, 130, 246, 0.05) 0%, transparent 100%);
}

.owner-name {
  display: flex;
  align-items: center;
  gap: 6px;
}

.official-badge {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  color: white;
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 10px;
  font-weight: 500;
}

.user-info:hover {
  opacity: 0.8;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}
</style>