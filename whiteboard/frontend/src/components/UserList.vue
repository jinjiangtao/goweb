<template>
  <div class="user-list">
    <div class="user-list-header">
      <h3 class="title">
        <el-icon color="#409eff"><User /></el-icon>
        在线成员 ({{ onlineUsers.length }})
      </h3>
    </div>
    
    <div class="user-list-content">
      <div 
        v-for="user in users" 
        :key="user.id" 
        class="user-item"
        :class="{ offline: !user.connected }"
      >
        <div class="user-avatar" :style="{ background: user.color }">
          {{ user.name.charAt(0) }}
        </div>
        <div class="user-info">
          <span class="user-name">{{ user.name }}</span>
          <span class="user-status">
            <span class="status-dot" :class="{ connected: user.connected }"></span>
            {{ user.connected ? '在线' : '离线' }}
          </span>
        </div>
        <div 
          v-if="user.id === currentUser.id" 
          class="user-tag"
        >
          我
        </div>
      </div>
      
      <el-empty 
        v-if="users.length === 0" 
        description="暂无成员" 
        :image-size="60"
      >
        <template #image>
          <el-icon :size="48" color="#dcdfe6"><UserFilled /></el-icon>
        </template>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useBoardStore } from '@/store/board'

const store = useBoardStore()

const users = computed(() => store.users)
const onlineUsers = computed(() => store.onlineUsers)
const currentUser = computed(() => store.currentUser)
</script>

<style lang="scss" scoped>
.user-list {
  width: 200px;
  background: #fff;
  border-left: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;
  z-index: 50;
  
  .user-list-header {
    padding: 12px 16px;
    border-bottom: 1px solid #e4e7ed;
    
    .title {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 14px;
      margin: 0;
      color: #303133;
      font-weight: 500;
    }
  }
  
  .user-list-content {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }
  
  .user-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px;
    border-radius: 6px;
    transition: background 0.2s;
    
    &:hover {
      background: #f5f7fa;
    }
    
    &.offline {
      opacity: 0.5;
    }
    
    .user-avatar {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      font-size: 14px;
      font-weight: 500;
      flex-shrink: 0;
    }
    
    .user-info {
      flex: 1;
      min-width: 0;
      
      .user-name {
        display: block;
        font-size: 13px;
        color: #303133;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
      
      .user-status {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 11px;
        color: #909399;
        
        .status-dot {
          width: 6px;
          height: 6px;
          border-radius: 50%;
          background: #dcdfe6;
          
          &.connected {
            background: #67c23a;
          }
        }
      }
    }
    
    .user-tag {
      background: #409eff;
      color: #fff;
      font-size: 10px;
      padding: 2px 6px;
      border-radius: 4px;
    }
  }
}

.user-list-content::-webkit-scrollbar {
  width: 4px;
}

.user-list-content::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 2px;
}
</style>
