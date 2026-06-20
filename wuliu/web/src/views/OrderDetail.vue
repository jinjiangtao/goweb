<template>
  <div class="detail-container">
    <div class="header">
      <div class="header-content">
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <h1>订单详情</h1>
      </div>
    </div>

    <div class="main-content" v-loading="loading">
      <div v-if="order" class="order-card">
        <div class="card-header">
          <div class="order-no">
            <span class="label">物流单号</span>
            <span class="value">{{ order.order_no }}</span>
          </div>
          <el-tag :type="getStatusType(order.status)" effect="dark" size="large">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>

        <div class="order-info">
          <div class="info-section">
            <h3>货物信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">货物名称</span>
                <span class="info-value">{{ order.goods_name || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">重量</span>
                <span class="info-value">{{ order.weight || 0 }} kg</span>
              </div>
            </div>
          </div>

          <div class="info-section">
            <h3>收发货信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">发货人</span>
                <span class="info-value">{{ order.sender || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">收件人</span>
                <span class="info-value">{{ order.receiver || '-' }}</span>
              </div>
              <div class="info-item full-width">
                <span class="info-label">发货地</span>
                <span class="info-value">{{ order.origin || '-' }}</span>
              </div>
              <div class="info-item full-width">
                <span class="info-label">目的地</span>
                <span class="info-value">{{ order.destination || '-' }}</span>
              </div>
            </div>
          </div>

          <div class="info-section">
            <h3>时间信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">创建时间</span>
                <span class="info-value">{{ formatTime(order.created_at) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">更新时间</span>
                <span class="info-value">{{ formatTime(order.updated_at) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="progress-section">
          <div class="progress-header">
            <span class="progress-title">运输进度</span>
            <span class="progress-value">{{ progress.toFixed(0) }}%</span>
          </div>
          <div class="progress-steps">
            <div
              v-for="(step, index) in steps"
              :key="index"
              class="step"
              :class="{ active: index <= currentStep, current: index === currentStep }"
            >
              <div class="step-dot">
                <el-icon v-if="index < currentStep"><CircleCheckFilled /></el-icon>
                <span v-else>{{ index + 1 }}</span>
              </div>
              <div class="step-label">{{ step }}</div>
            </div>
          </div>
        </div>
      </div>

      <div class="timeline-card">
        <div class="card-header">
          <h3>物流轨迹</h3>
          <div class="timeline-actions">
            <el-button size="small" @click="toggleSort">
              <el-icon><Sort /></el-icon>
              {{ isAsc ? '倒序' : '正序' }}
            </el-button>
            <el-button size="small" type="primary" @click="goToManage">
              管理节点
            </el-button>
          </div>
        </div>

        <div class="timeline-wrapper">
          <div
            v-for="(node, index) in displayNodes"
            :key="node.id"
            class="timeline-node"
            :class="{
              abnormal: node.is_abnormal,
              first: index === 0
            }"
          >
            <div class="node-left">
              <div class="node-dot" :class="getNodeClass(node)">
                <el-icon v-if="node.is_abnormal"><WarningFilled /></el-icon>
                <el-icon v-else-if="node.status === 'delivered'"><CircleCheckFilled /></el-icon>
                <el-icon v-else><CircleCheck /></el-icon>
              </div>
              <div v-if="index < displayNodes.length - 1" class="node-line"></div>
            </div>
            <div class="node-content">
              <div class="node-header">
                <span class="node-title">{{ node.node_name }}</span>
                <el-tag
                  v-if="node.is_abnormal"
                  type="danger"
                  size="small"
                  effect="light"
                >
                  异常
                </el-tag>
                <span class="node-time">{{ formatTime(node.occurred_at) }}</span>
              </div>
              <div class="node-desc">{{ node.description }}</div>
              <div class="node-meta">
                <span v-if="node.location" class="meta-item">
                  <el-icon><Location /></el-icon>
                  {{ node.location }}
                </span>
                <span v-if="node.operator" class="meta-item">
                  <el-icon><User /></el-icon>
                  {{ node.operator }}
                </span>
              </div>
              <div v-if="node.is_abnormal && node.abnormal_reason" class="node-abnormal">
                <div class="abnormal-label">异常原因：</div>
                <div class="abnormal-content">{{ node.abnormal_reason }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getOrder, getProgress } from '../api'
import {
  ArrowLeft,
  Location,
  User,
  WarningFilled,
  CircleCheckFilled,
  CircleCheck,
  Sort
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const order = ref(null)
const loading = ref(false)
const progress = ref(0)
const currentStep = ref(0)
const isAsc = ref(false)

const steps = ['创建订单', '已揽收', '运输中', '已到达', '已签收']

const displayNodes = computed(() => {
  if (!order.value?.nodes) return []
  const nodes = [...order.value.nodes]
  if (isAsc.value) {
    return nodes
  }
  return nodes.reverse()
})

const getStatusText = (status) => {
  const statusMap = {
    pending: '待揽收',
    picked: '已揽收',
    in_transit: '运输中',
    arrived: '已到达',
    delivered: '已签收',
    exception: '异常'
  }
  return statusMap[status] || status
}

const getStatusType = (status) => {
  const typeMap = {
    pending: 'info',
    picked: 'success',
    in_transit: 'primary',
    arrived: 'warning',
    delivered: 'success',
    exception: 'danger'
  }
  return typeMap[status] || 'info'
}

const getNodeClass = (node) => {
  if (node.is_abnormal) return 'node-exception'
  return `node-${node.status}`
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const loadOrder = async () => {
  loading.value = true
  try {
    const id = route.params.id
    const res = await getOrder(id)
    order.value = res.data.data
    await loadProgress()
    updateCurrentStep()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const loadProgress = async () => {
  try {
    const id = route.params.id
    const res = await getProgress(id)
    progress.value = res.data.progress
  } catch (err) {
    console.error(err)
  }
}

const updateCurrentStep = () => {
  if (!order.value) return
  const statusMap = {
    pending: 0,
    picked: 1,
    in_transit: 2,
    arrived: 3,
    delivered: 4
  }
  currentStep.value = statusMap[order.value.status] ?? 0
}

const toggleSort = () => {
  isAsc.value = !isAsc.value
}

const goBack = () => {
  router.back()
}

const goToManage = () => {
  router.push('/manage')
}

onMounted(() => {
  loadOrder()
})
</script>

<style scoped>
.detail-container {
  min-height: 100vh;
  background: #f0f2f5;
}

.header {
  background: white;
  border-bottom: 1px solid #e4e7ed;
  padding: 16px 24px;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-content h1 {
  font-size: 20px;
  margin: 0;
  color: #303133;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  display: flex;
  gap: 24px;
}

.order-card {
  flex: 1;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.card-header {
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.timeline-card .card-header {
  background: white;
  color: #303133;
  border-bottom: 1px solid #ebeef5;
  padding: 20px 24px;
}

.order-no .label {
  font-size: 14px;
  opacity: 0.9;
  display: block;
  margin-bottom: 6px;
}

.order-no .value {
  font-size: 24px;
  font-weight: 600;
}

.order-info {
  padding: 24px;
}

.info-section {
  margin-bottom: 24px;
}

.info-section:last-child {
  margin-bottom: 0;
}

.info-section h3 {
  font-size: 16px;
  color: #303133;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-label {
  font-size: 13px;
  color: #909399;
}

.info-value {
  font-size: 14px;
  color: #303133;
}

.progress-section {
  padding: 24px;
  background: #fafafa;
  border-top: 1px solid #ebeef5;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.progress-title {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.progress-value {
  font-size: 18px;
  font-weight: 600;
  color: #409eff;
}

.progress-steps {
  display: flex;
  justify-content: space-between;
  position: relative;
}

.progress-steps::before {
  content: '';
  position: absolute;
  top: 20px;
  left: 20px;
  right: 20px;
  height: 3px;
  background: #e4e7ed;
  z-index: 0;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  position: relative;
  z-index: 1;
}

.step-dot {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #e4e7ed;
  color: #c0c4cc;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s;
}

.step.active .step-dot {
  background: #409eff;
  color: white;
}

.step.current .step-dot {
  background: #667eea;
  color: white;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.2);
}

.step-label {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
}

.step.active .step-label {
  color: #409eff;
  font-weight: 500;
}

.timeline-card {
  flex: 1;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  max-height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}

.timeline-card .card-header h3 {
  font-size: 16px;
  margin: 0;
}

.timeline-actions {
  display: flex;
  gap: 10px;
}

.timeline-wrapper {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.timeline-node {
  display: flex;
  gap: 16px;
}

.node-left {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
  width: 32px;
}

.node-dot {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
  transition: transform 0.3s;
}

.node-dot.node-pending {
  background: #909399;
  color: white;
}

.node-dot.node-picked {
  background: #67c23a;
  color: white;
}

.node-dot.node-in_transit {
  background: #409eff;
  color: white;
}

.node-dot.node-arrived {
  background: #e6a23c;
  color: white;
}

.node-dot.node-delivered {
  background: #67c23a;
  color: white;
}

.node-dot.node-exception {
  background: #f56c6c;
  color: white;
}

.timeline-node.first .node-dot {
  transform: scale(1.1);
}

.node-line {
  width: 2px;
  flex: 1;
  background: #e4e7ed;
  min-height: 40px;
}

.timeline-node.abnormal .node-line {
  background: #fde2e2;
}

.node-content {
  flex: 1;
  padding-bottom: 24px;
}

.node-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.node-title {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.timeline-node.first .node-title {
  color: #409eff;
}

.timeline-node.abnormal .node-title {
  color: #f56c6c;
}

.node-time {
  font-size: 12px;
  color: #909399;
  margin-left: auto;
}

.node-desc {
  color: #606266;
  font-size: 13px;
  margin-bottom: 8px;
  line-height: 1.5;
}

.node-meta {
  display: flex;
  gap: 20px;
  color: #909399;
  font-size: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.node-abnormal {
  margin-top: 10px;
  padding: 12px;
  background: #fef0f0;
  border-radius: 6px;
  border-left: 3px solid #f56c6c;
}

.abnormal-label {
  color: #f56c6c;
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 4px;
}

.abnormal-content {
  color: #f56c6c;
  font-size: 13px;
  line-height: 1.5;
}

@media (max-width: 900px) {
  .main-content {
    flex-direction: column;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
