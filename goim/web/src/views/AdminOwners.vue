<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { createOwner, getAllOwners, deleteOwner, getOwnerMembers, removeOwnerMember } from '@/api'
import type { Owner, OwnerMember } from '@/types'
import { ElMessage, ElTable, ElTableColumn, ElButton, ElDialog, ElForm, ElFormItem, ElInput, ElTag, ElPopconfirm } from 'element-plus'

const router = useRouter()
const adminStore = useAdminStore()

const owners = ref<Owner[]>([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showMembersDialog = ref(false)
const currentOwner = ref<Owner | null>(null)
const ownerMembers = ref<OwnerMember[]>([])

const form = ref({
  name: '',
  description: '',
  avatar: ''
})

async function loadOwners() {
  loading.value = true
  try {
    owners.value = await getAllOwners()
  } catch (error) {
    ElMessage.error('加载群主列表失败')
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!form.value.name.trim()) {
    ElMessage.warning('请输入群主名称')
    return
  }
  
  loading.value = true
  try {
    await createOwner(form.value.name, form.value.description, form.value.avatar, adminStore.token)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    form.value = { name: '', description: '', avatar: '' }
    await loadOwners()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    loading.value = false
  }
}

async function handleDelete(ownerID: string) {
  loading.value = true
  try {
    await deleteOwner(ownerID, adminStore.token)
    ElMessage.success('删除成功')
    await loadOwners()
  } catch (error) {
    ElMessage.error('删除失败')
  } finally {
    loading.value = false
  }
}

async function openMembersDialog(owner: Owner) {
  currentOwner.value = owner
  loading.value = true
  try {
    ownerMembers.value = await getOwnerMembers(owner.id)
    showMembersDialog.value = true
  } catch (error) {
    ElMessage.error('加载成员列表失败')
  } finally {
    loading.value = false
  }
}

async function handleRemoveMember(member: OwnerMember) {
  loading.value = true
  try {
    await removeOwnerMember(member.owner_id, member.user_id, adminStore.token)
    ownerMembers.value = ownerMembers.value.filter(m => m.id !== member.id)
    ElMessage.success('移除成功')
  } catch (error) {
    ElMessage.error('移除失败')
  } finally {
    loading.value = false
  }
}

function goToUsers() {
  router.push('/admin/users')
}

function goToMessages() {
  router.push('/admin/messages')
}

function handleLogout() {
  adminStore.logout()
  router.push('/admin')
}

watch(() => adminStore.isLoggedIn, (loggedIn) => {
  if (!loggedIn) {
    router.push('/admin')
  }
})

onMounted(() => {
  if (!adminStore.isLoggedIn) {
    router.push('/admin')
    return
  }
  loadOwners()
})
</script>

<template>
  <div class="admin-owners-container">
    <div class="admin-header">
      <h1>GoIM 管理后台 - 群主管理</h1>
      <div class="header-actions">
        <el-button @click="goToUsers">用户管理</el-button>
        <el-button @click="goToMessages">聊天记录</el-button>
        <el-button type="primary" @click="showCreateDialog = true" :loading="loading">
          创建群主
        </el-button>
        <el-button @click="handleLogout">退出</el-button>
      </div>
    </div>
    
    <div class="admin-content">
      <el-table :data="owners" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="群主ID" width="200" />
        <el-table-column prop="name" label="群主名称" width="150" />
        <el-table-column prop="description" label="简介" width="250" />
        <el-table-column prop="avatar" label="头像" width="100">
          <template #default="{ row }">
            <div v-if="row.avatar" class="avatar-preview">
              <img :src="row.avatar" :alt="row.name" />
            </div>
            <div v-else class="avatar-placeholder">
              {{ row.name.charAt(0).toUpperCase() }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="200" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="openMembersDialog(row)">查看成员</el-button>
            <el-popconfirm
              title="确定要删除这个群主吗？"
              confirm-button-text="确定"
              cancel-button-text="取消"
              @confirm="handleDelete(row.id)"
            >
              <el-button size="small" type="danger">删除</el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <ElDialog title="创建群主" v-model="showCreateDialog" @close="form = { name: '', description: '', avatar: '' }">
      <ElForm :model="form" label-width="80px">
        <ElFormItem label="群主名称" required>
          <ElInput v-model="form.name" placeholder="请输入群主名称" />
        </ElFormItem>
        <ElFormItem label="群主简介">
          <ElInput v-model="form.description" type="textarea" placeholder="请输入群主简介" />
        </ElFormItem>
        <ElFormItem label="头像URL">
          <ElInput v-model="form.avatar" placeholder="请输入头像URL（可选）" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="showCreateDialog = false">取消</ElButton>
        <ElButton type="primary" @click="handleCreate" :loading="loading">创建</ElButton>
      </template>
    </ElDialog>

    <ElDialog :title="currentOwner?.name + ' - 成员列表'" v-model="showMembersDialog">
      <el-table :data="ownerMembers" v-loading="loading" style="width: 100%">
        <el-table-column prop="user_id" label="用户ID" width="200" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="nickname" label="昵称" width="150" />
        <el-table-column prop="joined_at" label="加入时间" width="200" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-popconfirm
              title="确定要移除这个成员吗？"
              confirm-button-text="确定"
              cancel-button-text="取消"
              @confirm="handleRemoveMember(row)"
            >
              <el-button size="small" type="danger">移除</el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <ElButton @click="showMembersDialog = false">关闭</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<style scoped>
.admin-owners-container {
  min-height: 100vh;
  background: #f0f2f5;
}

.admin-header {
  background: white;
  padding: 20px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.admin-header h1 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.admin-content {
  padding: 20px 40px;
}

.avatar-preview {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
}
</style>