<template>
  <div class="tags-page">
    <div class="page-header flex items-center justify-between mb-md">
      <h2 class="page-title">标签管理</h2>
      <div class="flex gap-sm">
        <el-button @click="showCreateGroupDialog">
          <el-icon><FolderAdd /></el-icon>
          新建分组
        </el-button>
        <el-button type="primary" @click="showCreateTagDialog(null)">
          <el-icon><Plus /></el-icon>
          新建标签
        </el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="tags-tabs">
      <el-tab-pane label="树形结构" name="tree">
        <div class="tree-panel">
          <div class="tree-header text-muted text-small mb-sm">
            提示：拖动标签可调整排序，右键可进行更多操作
          </div>
          <el-tree
            ref="treeRef"
            :data="treeData"
            node-key="id"
            default-expand-all
            :expand-on-click-node="false"
            draggable
            :allow-drop="allowDrop"
            :allow-drag="allowDrag"
            @node-drop="handleNodeDrop"
          >
            <template #default="{ node, data }">
              <div class="tree-node flex items-center justify-between w-full">
                <div class="flex items-center gap-sm">
                  <span
                    class="tag-color-dot"
                    :style="{ backgroundColor: data.color }"
                  />
                  <span>{{ data.name }}</span>
                  <el-tag size="small" type="info" effect="plain">{{ getTagCount(data.id) }}</el-tag>
                </div>
                <div class="node-actions flex gap-sm">
                  <el-button
                    link
                    size="small"
                    @click.stop="showCreateTagDialog(data)"
                  >
                    <el-icon><Plus /></el-icon>
                    子标签
                  </el-button>
                  <el-button
                    link
                    size="small"
                    @click.stop="showEditTagDialog(data)"
                  >
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-button>
                  <el-button
                    link
                    size="small"
                    type="danger"
                    @click.stop="handleDeleteTag(data)"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </div>
              </div>
            </template>
          </el-tree>
        </div>
      </el-tab-pane>

      <el-tab-pane label="分组管理" name="group">
        <div class="groups-panel">
          <el-empty v-if="groups.length === 0" description="暂无分组，点击右上角新建分组" />
          <div
            v-for="group in groups"
            :key="group.id"
            class="group-card"
          >
            <div class="group-header flex items-center justify-between">
              <div class="flex items-center gap-sm">
                <el-icon :size="20"><Folder /></el-icon>
                <span class="group-name">{{ group.name }}</span>
                <el-tag size="small" type="info" effect="plain">{{ group.tags?.length || 0 }}</el-tag>
              </div>
              <div class="flex gap-sm">
                <el-button link size="small" @click="showEditGroupDialog(group)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button
                  link
                  size="small"
                  type="danger"
                  @click="handleDeleteGroup(group)"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </div>
            <div class="group-tags flex gap-sm" style="flex-wrap: wrap; padding: 12px 0">
              <el-tag
                v-for="tag in group.tags"
                :key="tag.id"
                :color="tag.color"
                effect="light"
                closable
                @close="unassignTagFromGroup(tag, group)"
              >
                {{ tag.name }}
              </el-tag>
              <el-button
                link
                type="primary"
                size="small"
                @click="showAssignTagDialog(group)"
              >
                <el-icon><Plus /></el-icon>
                添加标签
              </el-button>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog
      v-model="createTagDialogVisible"
      :title="editingTag ? '编辑标签' : (parentTag ? `创建子标签 - ${parentTag.name}` : '新建标签')"
      width="420px"
    >
      <el-form :model="tagForm" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="tagForm.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-color-picker v-model="tagForm.color" />
        </el-form-item>
        <el-form-item v-if="!editingTag" label="父标签">
          <el-tree-select
            v-model="tagForm.parentId"
            :data="treeData"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            check-strictly
            clearable
            placeholder="不选则为根标签"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createTagDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveTag">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="createGroupDialogVisible"
      :title="editingGroup ? '编辑分组' : '新建分组'"
      width="400px"
    >
      <el-form :model="groupForm" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="groupForm.name" placeholder="请输入分组名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createGroupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveGroup">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="assignTagDialogVisible"
      title="添加标签到分组"
      width="420px"
    >
      <el-select
        v-model="selectedTagToAssign"
        placeholder="选择要添加的标签"
        style="width: 100%"
        filterable
      >
        <el-option
          v-for="tag in unassignedTags"
          :key="tag.id"
          :label="tag.name"
          :value="tag.id"
        >
          <div class="flex items-center gap-sm">
            <span
              class="tag-color-dot"
              :style="{ backgroundColor: tag.color }"
            />
            <span>{{ tag.name }}</span>
          </div>
        </el-option>
      </el-select>
      <template #footer>
        <el-button @click="assignTagDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAssignTag">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { tagApi, tagGroupApi } from '@/api'

const activeTab = ref('tree')
const treeData = ref([])
const flatTagList = ref([])
const groups = ref([])
const tagCounts = ref({})

const treeRef = ref(null)

const createTagDialogVisible = ref(false)
const editingTag = ref(null)
const parentTag = ref(null)
const tagForm = ref({ name: '', color: '#409EFF', parentId: null })

const createGroupDialogVisible = ref(false)
const editingGroup = ref(null)
const groupForm = ref({ name: '' })

const assignTagDialogVisible = ref(false)
const currentGroup = ref(null)
const selectedTagToAssign = ref(null)

const flatAllTags = computed(() => {
  const result = []
  const flatten = (tags) => {
    tags.forEach(tag => {
      result.push(tag)
      if (tag.children && tag.children.length) {
        flatten(tag.children)
      }
    })
  }
  flatten(treeData.value)
  return result
})

const unassignedTags = computed(() => {
  const assignedIds = new Set()
  groups.value.forEach(g => {
    (g.tags || []).forEach(t => assignedIds.add(t.id))
  })
  return flatAllTags.value.filter(t => !assignedIds.has(t.id))
})

const getTagCount = (tagId) => {
  return tagCounts.value[tagId] || 0
}

const loadTags = async () => {
  treeData.value = await tagApi.list()
  flatTagList.value = await tagApi.list({ flat: 'true' })

  for (const tag of flatTagList.value) {
    try {
      const result = await tagApi.getCount(tag.id)
      tagCounts.value[tag.id] = result.count
    } catch (e) {}
  }
}

const loadGroups = async () => {
  groups.value = await tagGroupApi.list()
}

const allowDrop = (draggingNode, dropNode, type) => {
  return type !== 'inner' ? true : true
}

const allowDrag = () => true

const handleNodeDrop = async (draggingNode, dropNode, dropType) => {
  const tagIds = []
  const collectIds = (nodes) => {
    nodes.forEach(node => {
      tagIds.push(node.id)
      if (node.children && node.children.length) {
        collectIds(node.children)
      }
    })
  }

  const reorderedIds = []
  const rootNodes = treeRef.value?.store?.root?.childNodes || []
  rootNodes.forEach(node => {
    collectNodeIds(node, reorderedIds)
  })

  try {
    await tagApi.reorder(reorderedIds)
    ElMessage.success('排序已更新')
  } catch (e) {}

  loadTags()
}

const collectNodeIds = (node, ids) => {
  ids.push(node.data.id)
  if (node.childNodes && node.childNodes.length) {
    node.childNodes.forEach(child => collectNodeIds(child, ids))
  }
}

const showCreateTagDialog = (parent) => {
  editingTag.value = null
  parentTag.value = parent
  tagForm.value = {
    name: '',
    color: '#409EFF',
    parentId: parent ? parent.id : null
  }
  createTagDialogVisible.value = true
}

const showEditTagDialog = (tag) => {
  editingTag.value = tag
  parentTag.value = null
  tagForm.value = {
    name: tag.name,
    color: tag.color,
    parentId: tag.parentId
  }
  createTagDialogVisible.value = true
}

const handleSaveTag = async () => {
  if (!tagForm.value.name.trim()) {
    ElMessage.warning('请输入标签名称')
    return
  }

  if (editingTag.value) {
    await tagApi.update(editingTag.value.id, tagForm.value)
    ElMessage.success('标签已更新')
  } else {
    await tagApi.create(tagForm.value)
    ElMessage.success('标签已创建')
  }

  createTagDialogVisible.value = false
  loadTags()
}

const handleDeleteTag = async (tag) => {
  try {
    await ElMessageBox.confirm(
      `确定删除标签"${tag.name}"？子标签也会被删除，已关联的笔记将取消关联。`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await tagApi.remove(tag.id)
    ElMessage.success('标签已删除')
    loadTags()
  } catch (e) {}
}

const showCreateGroupDialog = () => {
  editingGroup.value = null
  groupForm.value = { name: '' }
  createGroupDialogVisible.value = true
}

const showEditGroupDialog = (group) => {
  editingGroup.value = group
  groupForm.value = { name: group.name }
  createGroupDialogVisible.value = true
}

const handleSaveGroup = async () => {
  if (!groupForm.value.name.trim()) {
    ElMessage.warning('请输入分组名称')
    return
  }

  if (editingGroup.value) {
    await tagGroupApi.update(editingGroup.value.id, groupForm.value)
    ElMessage.success('分组已更新')
  } else {
    await tagGroupApi.create(groupForm.value)
    ElMessage.success('分组已创建')
  }

  createGroupDialogVisible.value = false
  loadGroups()
}

const handleDeleteGroup = async (group) => {
  try {
    await ElMessageBox.confirm(
      `确定删除分组"${group.name}"？分组内的标签不会被删除，只会取消分组关联。`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await tagGroupApi.remove(group.id)
    ElMessage.success('分组已删除')
    loadGroups()
  } catch (e) {}
}

const showAssignTagDialog = (group) => {
  currentGroup.value = group
  selectedTagToAssign.value = null
  assignTagDialogVisible.value = true
}

const handleAssignTag = async () => {
  if (!selectedTagToAssign.value) {
    ElMessage.warning('请选择标签')
    return
  }

  await tagApi.update(selectedTagToAssign.value, { groupId: currentGroup.value.id })
  assignTagDialogVisible.value = false
  ElMessage.success('已添加到分组')
  loadGroups()
}

const unassignTagFromGroup = async (tag, group) => {
  try {
    await tagApi.update(tag.id, { groupId: null })
    ElMessage.success('已从分组移除')
    loadGroups()
  } catch (e) {}
}

onMounted(() => {
  loadTags()
  loadGroups()
})
</script>

<style lang="scss" scoped>
.tags-page {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  min-height: 100%;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.tags-tabs {
  margin-top: 16px;
}

.tree-panel {
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  min-height: 400px;
}

.tag-color-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.tree-node {
  padding: 4px 0;
  width: 100%;
}

.node-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.tree-node:hover .node-actions {
  opacity: 1;
}

.groups-panel {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 16px;
}

.group-card {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 16px;
  transition: all 0.2s;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }
}

.group-header {
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.group-name {
  font-size: 16px;
  font-weight: 500;
}

.group-tags {
  min-height: 40px;
}
</style>
