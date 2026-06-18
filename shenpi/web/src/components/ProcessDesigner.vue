<template>
  <div class="process-designer">
    <div class="designer-header">
      <div class="header-left">
        <el-input
          v-model="templateName"
          placeholder="请输入流程名称"
          class="name-input"
          clearable
        />
        <el-select v-model="templateType" class="type-select">
          <el-option label="请假审批" value="leave" />
          <el-option label="报销审批" value="reimburse" />
          <el-option label="通用申请" value="general" />
        </el-select>
      </div>
      <div class="header-right">
        <el-button @click="saveTemplate" type="primary">保存模板</el-button>
        <el-button @click="clearCanvas">清空画布</el-button>
      </div>
    </div>

    <div class="designer-body">
      <div class="node-panel">
        <div class="panel-title">节点类型</div>
        <div
          v-for="nodeType in nodeTypes"
          :key="nodeType.type"
          class="node-item"
          draggable="true"
          @dragstart="onDragStart($event, nodeType)"
        >
          <div class="node-icon" :class="nodeType.type">
            <el-icon><component :is="nodeType.icon" /></el-icon>
          </div>
          <span class="node-label">{{ nodeType.label }}</span>
        </div>
      </div>

      <div
        class="canvas-panel"
        ref="canvasRef"
        @drop="onDrop"
        @dragover.prevent
        @click="onCanvasClick"
      >
        <svg class="edge-svg">
          <defs>
            <marker
              id="arrowhead"
              markerWidth="10"
              markerHeight="7"
              refX="9"
              refY="3.5"
              orient="auto"
            >
              <polygon points="0 0, 10 3.5, 0 7" fill="#909399" />
            </marker>
          </defs>
          <g>
            <path
              v-for="edge in edges"
              :key="edge.id"
              :d="getEdgePath(edge)"
              stroke="#909399"
              stroke-width="2"
              fill="none"
              marker-end="url(#arrowhead)"
              :class="{ 'selected': selectedEdge?.id === edge.id }"
              @click.stop="selectEdge(edge)"
            />
            <text
              v-for="edge in edges"
              :key="'label-' + edge.id"
              :x="getEdgeLabelPosition(edge).x"
              :y="getEdgeLabelPosition(edge).y"
              text-anchor="middle"
              fill="#606266"
              font-size="12"
            >{{ edge.label }}</text>
          </g>
        </svg>

        <div
          v-for="node in nodes"
          :key="node.id"
          class="canvas-node"
          :class="[node.type, { 'selected': selectedNode?.id === node.id }]"
          :style="{ left: node.x + 'px', top: node.y + 'px' }"
          @click.stop="selectNode(node)"
          @mousedown="onNodeMouseDown($event, node)"
        >
          <div class="node-content">
            <el-icon class="node-content-icon"><component :is="getNodeIcon(node.type)" /></el-icon>
            <span class="node-content-label">{{ node.label }}</span>
          </div>
          <div
            v-if="node.type !== 'start' && node.type !== 'end'"
            class="connection-point top"
            @mousedown.stop="onConnectionStart($event, node, 'top')"
          ></div>
          <div
            v-if="node.type !== 'start' && node.type !== 'end'"
            class="connection-point bottom"
            @mousedown.stop="onConnectionStart($event, node, 'bottom')"
          ></div>
          <div
            v-if="node.type !== 'start' && node.type !== 'end'"
            class="connection-point left"
            @mousedown.stop="onConnectionStart($event, node, 'left')"
          ></div>
          <div
            v-if="node.type !== 'start' && node.type !== 'end'"
            class="connection-point right"
            @mousedown.stop="onConnectionStart($event, node, 'right')"
          ></div>
        </div>

        <svg v-if="connecting" class="temp-edge-svg">
          <path
            :d="tempEdgePath"
            stroke="#409eff"
            stroke-width="2"
            stroke-dasharray="5,5"
            fill="none"
          />
        </svg>
      </div>

      <div class="property-panel">
        <div class="panel-title">属性配置</div>
        <div v-if="selectedNode" class="property-content">
          <el-form label-width="80px">
            <el-form-item label="节点ID">
              <el-input v-model="selectedNode.id" disabled />
            </el-form-item>
            <el-form-item label="节点名称">
              <el-input v-model="selectedNode.label" @input="updateNodeLabel" />
            </el-form-item>
            <el-form-item label="节点类型">
              <el-input :value="getNodeTypeName(selectedNode.type)" disabled />
            </el-form-item>
            <template v-if="selectedNode.type === 'approval'">
              <el-form-item label="审批人类型">
                <el-select v-model="selectedNode.approverType" @change="onApproverTypeChange">
                  <el-option label="指定人员" value="specified" />
                  <el-option label="申请人" value="applicant" />
                  <el-option label="所有人" value="all" />
                </el-select>
              </el-form-item>
              <el-form-item v-if="selectedNode.approverType === 'specified'" label="审批人员">
                <el-select
                  v-model="selectedApprovers"
                  multiple
                  filterable
                  placeholder="请选择审批人"
                  @change="onApproversChange"
                >
                  <el-option
                    v-for="user in userList"
                    :key="user.id"
                    :label="user.name"
                    :value="user.id"
                  />
                </el-select>
              </el-form-item>
            </template>
            <template v-if="selectedNode.type === 'condition'">
              <el-form-item label="流转条件">
                <el-input
                  v-model="selectedNode.condition"
                  type="textarea"
                  :rows="3"
                  placeholder="例如: amount > 1000"
                />
              </el-form-item>
            </template>
            <el-form-item>
              <el-button type="danger" @click="deleteNode">删除节点</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div v-else-if="selectedEdge" class="property-content">
          <el-form label-width="80px">
            <el-form-item label="连线标签">
              <el-input v-model="selectedEdge.label" placeholder="请输入标签" />
            </el-form-item>
            <el-form-item label="流转条件">
              <el-input
                v-model="selectedEdge.condition"
                type="textarea"
                :rows="3"
                placeholder="例如: days <= 3"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" @click="deleteEdge">删除连线</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div v-else class="empty-tip">
          <el-empty description="请选择节点或连线进行配置" />
        </div>
      </div>
    </div>

    <div class="form-config-panel">
      <div class="panel-title">表单字段配置</div>
      <div class="form-config-content">
        <div class="form-actions">
          <el-button @click="addFormField" size="small" type="primary">添加字段</el-button>
        </div>
        <el-table :data="formFields" border>
          <el-table-column prop="name" label="字段名" width="120">
            <template #default="{ row }">
              <el-input v-model="row.name" size="small" />
            </template>
          </el-table-column>
          <el-table-column prop="label" label="显示名" width="120">
            <template #default="{ row }">
              <el-input v-model="row.label" size="small" />
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              <el-select v-model="row.type" size="small">
                <el-option label="输入框" value="input" />
                <el-option label="文本域" value="textarea" />
                <el-option label="数字" value="number" />
                <el-option label="下拉框" value="select" />
                <el-option label="日期" value="date" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column prop="required" label="必填" width="60">
            <template #default="{ row }">
              <el-switch v-model="row.required" size="small" />
            </template>
          </el-table-column>
          <el-table-column prop="placeholder" label="占位符">
            <template #default="{ row }">
              <el-input v-model="row.placeholder" size="small" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80">
            <template #default="{ $index }">
              <el-button type="danger" size="small" @click="removeFormField($index)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  CirclePlus,
  CircleCheck,
  User,
  Operation,
  RefreshRight,
  DocumentAdd
} from '@element-plus/icons-vue'
import { getTemplateDetail, createTemplate, updateTemplate, createDefaultTemplate } from '@/api/template'
import { getUserList } from '@/api/user'

const route = useRoute()
const router = useRouter()

const templateId = ref(route.params.id)
const templateName = ref('')
const templateType = ref('general')
const nodes = ref([])
const edges = ref([])
const formFields = ref([])
const selectedNode = ref(null)
const selectedEdge = ref(null)
const selectedApprovers = ref([])
const userList = ref([])

const canvasRef = ref(null)
const connecting = ref(false)
const connectionStart = ref(null)
const tempEdgePath = ref('')
const draggingNode = ref(null)
const dragOffset = reactive({ x: 0, y: 0 })

const nodeTypes = [
  { type: 'start', label: '开始节点', icon: CirclePlus },
  { type: 'approval', label: '审批节点', icon: User },
  { type: 'condition', label: '条件节点', icon: Operation },
  { type: 'end', label: '结束节点', icon: CircleCheck }
]

const getNodeIcon = (type) => {
  const nodeType = nodeTypes.find(n => n.type === type)
  return nodeType ? nodeType.icon : Operation
}

const getNodeTypeName = (type) => {
  const nodeType = nodeTypes.find(n => n.type === type)
  return nodeType ? nodeType.label : type
}

onMounted(() => {
  loadUserList()
  if (templateId.value) {
    loadTemplate()
  } else {
    initDefaultNodes()
    if (route.query.name) {
      templateName.value = route.query.name
    }
    if (route.query.type) {
      templateType.value = route.query.type
    }
  }
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
})

onUnmounted(() => {
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)
})

const loadUserList = async () => {
  try {
    const res = await getUserList()
    userList.value = res.data
  } catch (e) {
    console.error('加载用户列表失败', e)
  }
}

const initDefaultNodes = () => {
  nodes.value = [
    { id: 'start', type: 'start', label: '开始', x: 100, y: 200, approverType: 'all', approvers: [] },
    { id: 'end', type: 'end', label: '结束', x: 700, y: 200, approverType: 'all', approvers: [] }
  ]
  formFields.value = [
    { name: 'title', label: '申请标题', type: 'input', required: true, placeholder: '请输入申请标题' },
    { name: 'content', label: '申请内容', type: 'textarea', required: true, placeholder: '请详细说明申请内容' }
  ]
}

const loadTemplate = async () => {
  try {
    const res = await getTemplateDetail(templateId.value)
    const template = res.data
    templateName.value = template.name
    templateType.value = template.type || 'general'
    if (template.nodes) {
      nodes.value = typeof template.nodes === 'string' ? JSON.parse(template.nodes) : template.nodes
    }
    if (template.edges) {
      edges.value = typeof template.edges === 'string' ? JSON.parse(template.edges) : template.edges
    }
    if (template.form_config) {
      formFields.value = typeof template.form_config === 'string' ? JSON.parse(template.form_config) : template.form_config
    } else if (template.FormConfig) {
      formFields.value = typeof template.FormConfig === 'string' ? JSON.parse(template.FormConfig) : template.FormConfig
    }
  } catch (e) {
    console.error('加载模板失败', e)
  }
}

const onDragStart = (event, nodeType) => {
  event.dataTransfer.setData('nodeType', JSON.stringify(nodeType))
}

const onDrop = (event) => {
  event.preventDefault()
  const nodeTypeData = event.dataTransfer.getData('nodeType')
  if (!nodeTypeData) return

  const nodeType = JSON.parse(nodeTypeData)
  const rect = canvasRef.value.getBoundingClientRect()
  const x = event.clientX - rect.left - 60
  const y = event.clientY - rect.top - 25

  const newNode = {
    id: nodeType.type + '_' + Date.now(),
    type: nodeType.type,
    label: nodeType.label,
    x: x,
    y: y,
    approverType: 'specified',
    approvers: [],
    condition: ''
  }

  nodes.value.push(newNode)
  selectNode(newNode)
}

const onCanvasClick = () => {
  selectedNode.value = null
  selectedEdge.value = null
}

const selectNode = (node) => {
  selectedNode.value = node
  selectedEdge.value = null
  if (node.approvers) {
    selectedApprovers.value = [...node.approvers]
  }
}

const selectEdge = (edge) => {
  selectedEdge.value = edge
  selectedNode.value = null
}

const onNodeMouseDown = (event, node) => {
  draggingNode.value = node
  const rect = event.currentTarget.getBoundingClientRect()
  dragOffset.x = event.clientX - rect.left
  dragOffset.y = event.clientY - rect.top
}

const onMouseMove = (event) => {
  if (draggingNode.value) {
    const rect = canvasRef.value.getBoundingClientRect()
    draggingNode.value.x = event.clientX - rect.left - dragOffset.x
    draggingNode.value.y = event.clientY - rect.top - dragOffset.y
  }

  if (connecting.value && connectionStart.value) {
    const rect = canvasRef.value.getBoundingClientRect()
    const endX = event.clientX - rect.left
    const endY = event.clientY - rect.top
    const startPos = getConnectionPointPosition(connectionStart.value.node, connectionStart.value.direction)
    tempEdgePath.value = `M ${startPos.x} ${startPos.y} L ${endX} ${endY}`
  }
}

const onMouseUp = (event) => {
  if (draggingNode.value) {
    draggingNode.value = null
  }

  if (connecting.value) {
    const targetNode = findNodeAtPosition(event.clientX, event.clientY)
    if (targetNode && connectionStart.value && targetNode.id !== connectionStart.value.node.id) {
      const existingEdge = edges.value.find(
        e => e.source === connectionStart.value.node.id && e.target === targetNode.id
      )
      if (!existingEdge) {
        const newEdge = {
          id: 'edge_' + Date.now(),
          source: connectionStart.value.node.id,
          target: targetNode.id,
          label: '',
          condition: ''
        }
        edges.value.push(newEdge)
      }
    }
    connecting.value = false
    connectionStart.value = null
    tempEdgePath.value = ''
  }
}

const findNodeAtPosition = (clientX, clientY) => {
  const rect = canvasRef.value.getBoundingClientRect()
  const x = clientX - rect.left
  const y = clientY - rect.top

  for (const node of nodes.value) {
    if (x >= node.x && x <= node.x + 120 && y >= node.y && y <= node.y + 50) {
      return node
    }
  }
  return null
}

const onConnectionStart = (event, node, direction) => {
  event.stopPropagation()
  connecting.value = true
  connectionStart.value = { node, direction }
  const startPos = getConnectionPointPosition(node, direction)
  tempEdgePath.value = `M ${startPos.x} ${startPos.y} L ${startPos.x} ${startPos.y}`
}

const getConnectionPointPosition = (node, direction) => {
  const nodeWidth = 120
  const nodeHeight = 50
  switch (direction) {
    case 'top':
      return { x: node.x + nodeWidth / 2, y: node.y }
    case 'bottom':
      return { x: node.x + nodeWidth / 2, y: node.y + nodeHeight }
    case 'left':
      return { x: node.x, y: node.y + nodeHeight / 2 }
    case 'right':
      return { x: node.x + nodeWidth, y: node.y + nodeHeight / 2 }
    default:
      return { x: node.x + nodeWidth / 2, y: node.y + nodeHeight / 2 }
  }
}

const getEdgePath = (edge) => {
  const sourceNode = nodes.value.find(n => n.id === edge.source)
  const targetNode = nodes.value.find(n => n.id === edge.target)
  if (!sourceNode || !targetNode) return ''

  const startX = sourceNode.x + 120
  const startY = sourceNode.y + 25
  const endX = targetNode.x
  const endY = targetNode.y + 25

  const midX = (startX + endX) / 2
  return `M ${startX} ${startY} C ${midX} ${startY}, ${midX} ${endY}, ${endX} ${endY}`
}

const getEdgeLabelPosition = (edge) => {
  const sourceNode = nodes.value.find(n => n.id === edge.source)
  const targetNode = nodes.value.find(n => n.id === edge.target)
  if (!sourceNode || !targetNode) return { x: 0, y: 0 }

  return {
    x: (sourceNode.x + targetNode.x) / 2 + 60,
    y: (sourceNode.y + targetNode.y) / 2 - 5
  }
}

const updateNodeLabel = () => {
  if (selectedNode.value) {
    const node = nodes.value.find(n => n.id === selectedNode.value.id)
    if (node) {
      node.label = selectedNode.value.label
    }
  }
}

const onApproverTypeChange = () => {
  if (selectedNode.value && selectedNode.value.approverType !== 'specified') {
    selectedNode.value.approvers = []
    selectedApprovers.value = []
  }
}

const onApproversChange = (val) => {
  if (selectedNode.value) {
    selectedNode.value.approvers = val
  }
}

const deleteNode = () => {
  if (!selectedNode.value) return
  if (selectedNode.value.type === 'start' || selectedNode.value.type === 'end') {
    ElMessage.warning('开始和结束节点不能删除')
    return
  }

  const nodeId = selectedNode.value.id
  nodes.value = nodes.value.filter(n => n.id !== nodeId)
  edges.value = edges.value.filter(e => e.source !== nodeId && e.target !== nodeId)
  selectedNode.value = null
}

const deleteEdge = () => {
  if (!selectedEdge.value) return
  edges.value = edges.value.filter(e => e.id !== selectedEdge.value.id)
  selectedEdge.value = null
}

const clearCanvas = () => {
  ElMessageBox.confirm('确定要清空画布吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    initDefaultNodes()
    edges.value = []
    selectedNode.value = null
    selectedEdge.value = null
    ElMessage.success('画布已清空')
  }).catch(() => {})
}

const addFormField = () => {
  formFields.value.push({
    name: 'field_' + Date.now(),
    label: '新字段',
    type: 'input',
    required: false,
    placeholder: ''
  })
}

const removeFormField = (index) => {
  formFields.value.splice(index, 1)
}

const saveTemplate = async () => {
  if (!templateName.value.trim()) {
    ElMessage.warning('请输入流程名称')
    return
  }

  const hasStart = nodes.value.some(n => n.type === 'start')
  const hasEnd = nodes.value.some(n => n.type === 'end')
  if (!hasStart || !hasEnd) {
    ElMessage.warning('流程必须包含开始和结束节点')
    return
  }

  const data = {
    name: templateName.value,
    type: templateType.value,
    description: '',
    nodes: nodes.value,
    edges: edges.value,
    form_config: formFields.value,
    status: 'published'
  }

  try {
    if (templateId.value) {
      await updateTemplate(templateId.value, data)
      ElMessage.success('模板更新成功')
    } else {
      await createTemplate(data)
      ElMessage.success('模板创建成功')
    }
    router.push('/template')
  } catch (e) {
    console.error('保存模板失败', e)
  }
}
</script>

<style scoped lang="scss">
.process-designer {
  height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.designer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;

  .header-left {
    display: flex;
    gap: 15px;
    align-items: center;

    .name-input {
      width: 300px;
    }

    .type-select {
      width: 150px;
    }
  }
}

.designer-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.node-panel {
  width: 180px;
  background: #fff;
  border-right: 1px solid #e4e7ed;
  padding: 15px;
  overflow-y: auto;

  .panel-title {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 15px;
    color: #303133;
  }

  .node-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    margin-bottom: 10px;
    background: #f5f7fa;
    border-radius: 4px;
    cursor: move;
    transition: all 0.2s;

    &:hover {
      background: #ecf5ff;
      transform: translateX(3px);
    }

    .node-icon {
      width: 32px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 4px;
      color: #fff;
      font-size: 16px;

      &.start {
        background: #67c23a;
      }
      &.approval {
        background: #409eff;
      }
      &.condition {
        background: #e6a23c;
      }
      &.end {
        background: #f56c6c;
      }
    }

    .node-label {
      font-size: 13px;
      color: #303133;
    }
  }
}

.canvas-panel {
  flex: 1;
  position: relative;
  overflow: auto;
  background-image:
    linear-gradient(#e4e7ed 1px, transparent 1px),
    linear-gradient(90deg, #e4e7ed 1px, transparent 1px);
  background-size: 20px 20px;

  .edge-svg,
  .temp-edge-svg {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;

    path {
      pointer-events: stroke;
      cursor: pointer;

      &.selected {
        stroke: #409eff;
        stroke-width: 3;
      }
    }
  }

  .temp-edge-svg {
    pointer-events: none;
  }

  .canvas-node {
    position: absolute;
    width: 120px;
    height: 50px;
    background: #fff;
    border: 2px solid #dcdfe6;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: move;
    transition: all 0.2s;
    user-select: none;
    z-index: 10;

    &.start {
      border-color: #67c23a;
      background: #f0f9eb;

      .node-content-icon {
        color: #67c23a;
      }
    }

    &.end {
      border-color: #f56c6c;
      background: #fef0f0;

      .node-content-icon {
        color: #f56c6c;
      }
    }

    &.approval {
      border-color: #409eff;
      background: #ecf5ff;

      .node-content-icon {
        color: #409eff;
      }
    }

    &.condition {
      border-color: #e6a23c;
      background: #fdf6ec;
      transform: rotate(45deg);

      .node-content {
        transform: rotate(-45deg);
      }

      .node-content-icon {
        color: #e6a23c;
      }
    }

    &.selected {
      box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.3);
    }

    .node-content {
      display: flex;
      align-items: center;
      gap: 5px;
      padding: 0 10px;

      .node-content-icon {
        font-size: 16px;
      }

      .node-content-label {
        font-size: 12px;
        color: #303133;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 80px;
      }
    }

    .connection-point {
      position: absolute;
      width: 10px;
      height: 10px;
      background: #409eff;
      border: 2px solid #fff;
      border-radius: 50%;
      cursor: crosshair;
      z-index: 20;

      &.top {
        top: -6px;
        left: 50%;
        transform: translateX(-50%);
      }

      &.bottom {
        bottom: -6px;
        left: 50%;
        transform: translateX(-50%);
      }

      &.left {
        left: -6px;
        top: 50%;
        transform: translateY(-50%);
      }

      &.right {
        right: -6px;
        top: 50%;
        transform: translateY(-50%);
      }

      &:hover {
        width: 14px;
        height: 14px;
      }
    }
  }
}

.property-panel {
  width: 280px;
  background: #fff;
  border-left: 1px solid #e4e7ed;
  padding: 15px;
  overflow-y: auto;

  .panel-title {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 15px;
    color: #303133;
  }

  .property-content {
    :deep(.el-form-item) {
      margin-bottom: 15px;
    }
  }

  .empty-tip {
    padding: 30px 0;
  }
}

.form-config-panel {
  height: 200px;
  background: #fff;
  border-top: 1px solid #e4e7ed;
  padding: 15px;

  .panel-title {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 10px;
    color: #303133;
  }

  .form-config-content {
    .form-actions {
      margin-bottom: 10px;
    }

    :deep(.el-table) {
      font-size: 12px;
    }
  }
}
</style>
