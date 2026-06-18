<template>
  <div class="approval-progress">
    <div class="progress-header">
      <h3>审批进度</h3>
      <el-tag :type="getStatusType(approvalStatus)">{{ getStatusText(approvalStatus) }}</el-tag>
    </div>

    <div class="flow-chart">
      <svg class="flow-svg" :width="canvasWidth" :height="canvasHeight">
        <defs>
          <marker
            id="progress-arrowhead"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
          >
            <polygon points="0 0, 10 3.5, 0 7" fill="#909399" />
          </marker>
        </defs>

        <g v-for="edge in flowEdges" :key="edge.id">
          <path
            :d="getEdgePath(edge)"
            :stroke="getEdgeColor(edge)"
            stroke-width="2"
            fill="none"
            marker-end="url(#progress-arrowhead)"
          />
          <text
            v-if="edge.label"
            :x="getEdgeLabelPos(edge).x"
            :y="getEdgeLabelPos(edge).y"
            text-anchor="middle"
            fill="#606266"
            font-size="12"
          >{{ edge.label }}</text>
        </g>

        <g v-for="node in flowNodes" :key="node.id">
          <rect
            :x="node.x"
            :y="node.y"
            :width="nodeWidth"
            :height="nodeHeight"
            :rx="6"
            :fill="getNodeColor(node)"
            :stroke="getNodeBorderColor(node)"
            stroke-width="2"
          />
          <text
            :x="node.x + nodeWidth / 2"
            :y="node.y + nodeHeight / 2 - 6"
            text-anchor="middle"
            fill="#303133"
            font-size="13"
            font-weight="500"
          >{{ node.label }}</text>
          <text
            :x="node.x + nodeWidth / 2"
            :y="node.y + nodeHeight / 2 + 12"
            text-anchor="middle"
            :fill="getNodeStatusColor(node)"
            font-size="11"
          >{{ getNodeStatusText(node) }}</text>
          <circle
            v-if="node.status === 'processing'"
            :cx="node.x + nodeWidth - 8"
            :cy="node.y + 8"
            r="6"
            fill="#409eff"
            class="pulse-animation"
          />
        </g>
      </svg>
    </div>

    <div class="approval-records">
      <div class="records-title">审批记录</div>
      <el-timeline>
        <el-timeline-item
          v-for="(record, index) in records"
          :key="record.id"
          :timestamp="formatTime(record.created_at)"
          :type="getRecordType(record.action)"
          :icon="getRecordIcon(record.action)"
          :color="getRecordColor(record.action)"
        >
          <div class="record-content">
            <div class="record-header">
              <span class="record-node">{{ record.node_name }}</span>
              <span class="record-action">{{ getActionText(record.action) }}</span>
            </div>
            <div class="record-body">
              <span class="record-user">{{ record.approver?.name || record.approver_name || '系统' }}</span>
              <span v-if="record.comment" class="record-comment">{{ record.comment }}</span>
            </div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </div>

    <div v-if="showActions && canOperate" class="approval-actions">
      <el-form :model="actionForm" label-width="80px">
        <el-form-item label="审批意见">
          <el-input
            v-model="actionForm.comment"
            type="textarea"
            :rows="3"
            placeholder="请输入审批意见"
          />
        </el-form-item>
        <el-form-item v-if="showAddSign">
          <el-select
            v-model="actionForm.signUser"
            filterable
            placeholder="选择加签人员"
            style="width: 200px; margin-right: 10px;"
          >
            <el-option
              v-for="user in userList"
              :key="user.id"
              :label="user.name"
              :value="user.id"
            />
          </el-select>
          <el-button type="warning" @click="handleAddSign">加签</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="handleApprove">同意</el-button>
          <el-button type="danger" @click="handleReject">驳回</el-button>
          <el-button type="warning" @click="handleRevoke" v-if="canRevoke">撤回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  CircleCheck,
  CircleClose,
  Clock,
  UserFilled,
  EditPen,
  Back
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  approveApproval,
  rejectApproval,
  revokeApproval,
  addSignApproval,
  getApprovalProgress
} from '@/api/approval'
import { getUserList } from '@/api/user'
import { useUserStore } from '@/store/user'

const props = defineProps({
  approvalId: {
    type: [String, Number],
    required: true
  },
  showActions: {
    type: Boolean,
    default: true
  },
  showAddSign: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['approved', 'rejected', 'revoked'])

const userStore = useUserStore()

const instance = ref(null)
const records = ref([])
const approvalStatus = ref('')
const userList = ref([])
const actionForm = ref({
  comment: '',
  signUser: null
})

const nodeWidth = 140
const nodeHeight = 60
const nodeSpacingX = 40
const nodeSpacingY = 80

const flowNodes = computed(() => {
  if (!instance.value?.nodes) return []
  return instance.value.nodes.map((ni, index) => {
    const nodeData = instance.value.nodes_data?.find(n => n.id === ni.node_id) || {}
    return {
      id: ni.node_id,
      label: nodeData.label || ni.node_id,
      type: nodeData.type,
      status: ni.status,
      x: 40 + index * (nodeWidth + nodeSpacingX),
      y: 40
    }
  })
})

const flowEdges = computed(() => {
  if (!instance.value?.edges) return []
  return instance.value.edges
})

const canvasWidth = computed(() => {
  return flowNodes.value.length * (nodeWidth + nodeSpacingX) + 80
})

const canvasHeight = computed(() => {
  return nodeHeight + 120
})

const canOperate = computed(() => {
  if (!instance.value || approvalStatus.value !== 'pending') return false
  const currentNode = instance.value.nodes.find(n => n.status === 'processing')
  if (!currentNode) return false
  return true
})

const canRevoke = computed(() => {
  if (!instance.value || approvalStatus.value !== 'pending') return false
  return true
})

onMounted(() => {
  loadProgress()
  loadUsers()
})

const loadProgress = async () => {
  try {
    const res = await getApprovalProgress(props.approvalId)
    instance.value = res.data.instance
    records.value = res.data.records
    approvalStatus.value = res.data.status

    if (res.data.instance?.nodes && res.data.instance?.edges) {
      const templateNodes = []
      res.data.instance.nodes.forEach(ni => {
        templateNodes.push({
          id: ni.node_id,
          label: ni.node_id,
          type: 'approval'
        })
      })
      instance.value.nodes_data = templateNodes
    }
  } catch (e) {
    console.error('加载审批进度失败', e)
  }
}

const loadUsers = async () => {
  try {
    const res = await getUserList()
    userList.value = res.data
  } catch (e) {
    console.error('加载用户列表失败', e)
  }
}

const getEdgePath = (edge) => {
  const sourceNode = flowNodes.value.find(n => n.id === edge.source)
  const targetNode = flowNodes.value.find(n => n.id === edge.target)
  if (!sourceNode || !targetNode) return ''

  const startX = sourceNode.x + nodeWidth
  const startY = sourceNode.y + nodeHeight / 2
  const endX = targetNode.x
  const endY = targetNode.y + nodeHeight / 2
  const midX = (startX + endX) / 2

  return `M ${startX} ${startY} C ${midX} ${startY}, ${midX} ${endY}, ${endX} ${endY}`
}

const getEdgeLabelPos = (edge) => {
  const sourceNode = flowNodes.value.find(n => n.id === edge.source)
  const targetNode = flowNodes.value.find(n => n.id === edge.target)
  if (!sourceNode || !targetNode) return { x: 0, y: 0 }
  return {
    x: (sourceNode.x + targetNode.x) / 2 + nodeWidth / 2,
    y: (sourceNode.y + targetNode.y) / 2 - 5
  }
}

const getEdgeColor = (edge) => {
  const sourceNode = flowNodes.value.find(n => n.id === edge.source)
  if (!sourceNode) return '#909399'
  if (sourceNode.status === 'approved' || sourceNode.status === 'processing') {
    return '#67c23a'
  }
  if (sourceNode.status === 'rejected') {
    return '#f56c6c'
  }
  return '#909399'
}

const getNodeColor = (node) => {
  switch (node.status) {
    case 'approved':
      return '#f0f9eb'
    case 'rejected':
      return '#fef0f0'
    case 'processing':
      return '#ecf5ff'
    default:
      return '#f5f7fa'
  }
}

const getNodeBorderColor = (node) => {
  switch (node.status) {
    case 'approved':
      return '#67c23a'
    case 'rejected':
      return '#f56c6c'
    case 'processing':
      return '#409eff'
    default:
      return '#dcdfe6'
  }
}

const getNodeStatusColor = (node) => {
  switch (node.status) {
    case 'approved':
      return '#67c23a'
    case 'rejected':
      return '#f56c6c'
    case 'processing':
      return '#409eff'
    default:
      return '#909399'
  }
}

const getNodeStatusText = (node) => {
  switch (node.status) {
    case 'approved':
      return '已通过'
    case 'rejected':
      return '已驳回'
    case 'processing':
      return '审批中'
    case 'skipped':
      return '已跳过'
    default:
      return '待审批'
  }
}

const getStatusType = (status) => {
  switch (status) {
    case 'approved':
      return 'success'
    case 'rejected':
      return 'danger'
    case 'revoked':
      return 'info'
    default:
      return 'warning'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'approved':
      return '已通过'
    case 'rejected':
      return '已驳回'
    case 'revoked':
      return '已撤回'
    case 'pending':
      return '审批中'
    default:
      return status
  }
}

const getRecordType = (action) => {
  switch (action) {
    case 'approve':
    case 'auto_pass':
    case 'end':
      return 'success'
    case 'reject':
      return 'danger'
    case 'start':
      return 'primary'
    case 'revoke':
      return 'info'
    default:
      return 'warning'
  }
}

const getRecordIcon = (action) => {
  switch (action) {
    case 'approve':
    case 'auto_pass':
    case 'end':
      return CircleCheck
    case 'reject':
      return CircleClose
    case 'start':
      return EditPen
    case 'revoke':
      return Back
    default:
      return Clock
  }
}

const getRecordColor = (action) => {
  switch (action) {
    case 'approve':
    case 'auto_pass':
    case 'end':
      return '#67c23a'
    case 'reject':
      return '#f56c6c'
    case 'start':
      return '#409eff'
    case 'revoke':
      return '#909399'
    default:
      return '#e6a23c'
  }
}

const getActionText = (action) => {
  switch (action) {
    case 'start':
      return '发起申请'
    case 'approve':
      return '同意'
    case 'reject':
      return '驳回'
    case 'revoke':
      return '撤回'
    case 'add_sign':
      return '加签'
    case 'transfer':
      return '转交'
    case 'auto_pass':
      return '自动通过'
    case 'end':
      return '流程结束'
    default:
      return action
  }
}

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const handleApprove = async () => {
  try {
    await ElMessageBox.confirm('确定同意该申请吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'success'
    })
    await approveApproval(props.approvalId, { comment: actionForm.value.comment })
    ElMessage.success('审批通过')
    emit('approved')
    loadProgress()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '审批失败')
    }
  }
}

const handleReject = async () => {
  if (!actionForm.value.comment) {
    ElMessage.warning('请输入驳回原因')
    return
  }
  try {
    await ElMessageBox.confirm('确定驳回该申请吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await rejectApproval(props.approvalId, { comment: actionForm.value.comment })
    ElMessage.success('已驳回')
    emit('rejected')
    loadProgress()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '驳回失败')
    }
  }
}

const handleRevoke = async () => {
  try {
    await ElMessageBox.confirm('确定撤回该申请吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await revokeApproval(props.approvalId, { comment: actionForm.value.comment })
    ElMessage.success('已撤回')
    emit('revoked')
    loadProgress()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '撤回失败')
    }
  }
}

const handleAddSign = async () => {
  if (!actionForm.value.signUser) {
    ElMessage.warning('请选择加签人员')
    return
  }
  try {
    await addSignApproval(props.approvalId, {
      sign_user_id: actionForm.value.signUser,
      comment: actionForm.value.comment
    })
    ElMessage.success('加签成功')
    loadProgress()
    actionForm.value.signUser = null
  } catch (e) {
    ElMessage.error(e.message || '加签失败')
  }
}

defineExpose({
  loadProgress
})
</script>

<style scoped lang="scss">
.approval-progress {
  background: #fff;
  border-radius: 8px;
  padding: 20px;

  .progress-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 15px;
    border-bottom: 1px solid #ebeef5;

    h3 {
      margin: 0;
      font-size: 16px;
      color: #303133;
    }
  }

  .flow-chart {
    overflow-x: auto;
    padding: 20px 0;
    margin-bottom: 30px;
    background: #fafafa;
    border-radius: 6px;

    .flow-svg {
      min-width: 600px;
    }

    :deep(.pulse-animation) {
      animation: pulse 1.5s infinite;
    }

    @keyframes pulse {
      0% {
        transform: scale(1);
        opacity: 1;
      }
      50% {
        transform: scale(1.3);
        opacity: 0.7;
      }
      100% {
        transform: scale(1);
        opacity: 1;
      }
    }
  }

  .approval-records {
    .records-title {
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 15px;
    }

    :deep(.el-timeline-item__content) {
      color: #606266;
    }

    .record-content {
      .record-header {
        margin-bottom: 5px;

        .record-node {
          font-weight: 500;
          color: #303133;
          margin-right: 10px;
        }

        .record-action {
          color: #409eff;
          font-size: 12px;
        }
      }

      .record-body {
        font-size: 13px;

        .record-user {
          color: #606266;
          margin-right: 10px;
        }

        .record-comment {
          color: #909399;
        }
      }
    }
  }

  .approval-actions {
    margin-top: 30px;
    padding-top: 20px;
    border-top: 1px solid #ebeef5;

    :deep(.el-form-item) {
      margin-bottom: 15px;
    }
  }
}
</style>
