<template>
  <div class="home-container">
    <div class="header">
      <div class="header-content">
        <h1 class="title">物流轨迹可视化查询系统</h1>
        <p class="subtitle">实时追踪 · 智能监控 · 全程可视化</p>
      </div>
    </div>

    <div class="main-content">
      <div class="search-section">
        <div class="search-box">
          <el-input
            v-model="searchOrderNo"
            placeholder="请输入物流单号进行查询"
            size="large"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
            <template #append>
              <el-button type="primary" @click="handleSearch" :loading="searching">
                查询
              </el-button>
            </template>
          </el-input>
        </div>
        <div class="quick-links">
          <el-button text @click="goToManage">管理后台</el-button>
          <span class="divider">|</span>
          <el-button text @click="loadSampleOrder">加载示例订单</el-button>
        </div>
      </div>

      <div v-if="orderData" class="result-section">
        <div class="order-info-card">
          <div class="order-header">
            <div class="order-no">
              <span class="label">物流单号：</span>
              <span class="value">{{ orderData.order_no }}</span>
            </div>
            <el-tag :type="getStatusType(orderData.status)" effect="dark" size="large">
              {{ getStatusText(orderData.status) }}
            </el-tag>
          </div>

          <div class="order-info">
            <div class="info-item">
              <span class="info-label">发货人：</span>
              <span class="info-value">{{ orderData.sender || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">收件人：</span>
              <span class="info-value">{{ orderData.receiver || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">货物名称：</span>
              <span class="info-value">{{ orderData.goods_name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">重量：</span>
              <span class="info-value">{{ orderData.weight || 0 }} kg</span>
            </div>
          </div>

          <div class="route-info">
            <div class="route-start">
              <el-icon class="route-icon"><LocationFilled /></el-icon>
              <span>{{ orderData.origin || '发货地' }}</span>
            </div>
            <div class="route-line">
              <div class="route-progress" :style="{ width: progressPercent + '%' }"></div>
              <el-icon class="truck-icon" :style="{ left: progressPercent + '%' }">
                <Van />
              </el-icon>
            </div>
            <div class="route-end">
              <el-icon class="route-icon"><LocationFilled /></el-icon>
              <span>{{ orderData.destination || '目的地' }}</span>
            </div>
          </div>

          <div class="progress-section">
            <div class="progress-header">
              <span class="progress-label">运输进度</span>
              <span class="progress-value">{{ progressPercent.toFixed(0) }}%</span>
            </div>
            <el-progress
              :percentage="progressPercent"
              :status="getProgressStatus()"
              :stroke-width="12"
              :text-inside="false"
            />
          </div>
        </div>

        <div class="timeline-card">
          <div class="timeline-header">
            <h3>物流轨迹</h3>
            <span class="node-count">共 {{ orderData.nodes?.length || 0 }} 个节点</span>
          </div>

          <div class="timeline-container">
            <div
              v-for="(node, index) in reversedNodes"
              :key="node.id"
              class="timeline-item"
              :class="{
                'is-abnormal': node.is_abnormal,
                'is-latest': index === 0
              }"
              @mouseenter="showNodeTooltip(node, $event)"
              @mouseleave="hideNodeTooltip"
            >
              <div class="timeline-dot">
                <div
                  class="dot"
                  :class="getNodeClass(node)"
                >
                  <el-icon v-if="node.is_abnormal"><WarningFilled /></el-icon>
                  <el-icon v-else-if="node.status === 'delivered'"><CircleCheckFilled /></el-icon>
                  <el-icon v-else><CircleCheck /></el-icon>
                </div>
                <div v-if="index < reversedNodes.length - 1" class="timeline-line"></div>
              </div>

              <div class="timeline-content">
                <div class="timeline-time">
                  {{ formatTime(node.occurred_at) }}
                </div>
                <div class="timeline-title">
                  <span class="node-name">{{ node.node_name }}</span>
                  <el-tag
                    v-if="node.is_abnormal"
                    type="danger"
                    size="small"
                    effect="light"
                    style="margin-left: 8px"
                  >
                    异常
                  </el-tag>
                </div>
                <div class="timeline-desc">{{ node.description }}</div>
                <div class="timeline-meta">
                  <span v-if="node.location" class="meta-item">
                    <el-icon><Location /></el-icon>
                    {{ node.location }}
                  </span>
                  <span v-if="node.operator" class="meta-item">
                    <el-icon><User /></el-icon>
                    {{ node.operator }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="!searching && searched && !orderData" class="empty-section">
        <el-empty description="未找到该订单，请检查单号是否正确">
          <el-button type="primary" @click="loadSampleOrder">查看示例订单</el-button>
        </el-empty>
      </div>

      <div v-else class="welcome-section">
        <div class="welcome-card">
          <el-icon class="welcome-icon"><Van /></el-icon>
          <h2>欢迎使用物流轨迹查询系统</h2>
          <p>输入您的物流单号，实时追踪包裹运输状态</p>
          <el-button type="primary" size="large" @click="loadSampleOrder">
            查看示例订单
          </el-button>
        </div>
      </div>
    </div>

    <div
      v-if="tooltipVisible"
      class="node-tooltip"
      :style="{ top: tooltipY + 'px', left: tooltipX + 'px' }"
    >
      <div class="tooltip-header">{{ tooltipNode?.node_name }}</div>
      <div class="tooltip-body">
        <div class="tooltip-row">
          <span class="tooltip-label">状态：</span>
          <el-tag :type="getNodeStatusType(tooltipNode?.status)" size="small">
            {{ getStatusText(tooltipNode?.status) }}
          </el-tag>
        </div>
        <div class="tooltip-row" v-if="tooltipNode?.location">
          <span class="tooltip-label">位置：</span>
          <span>{{ tooltipNode.location }}</span>
        </div>
        <div class="tooltip-row" v-if="tooltipNode?.description">
          <span class="tooltip-label">详情：</span>
          <span>{{ tooltipNode.description }}</span>
        </div>
        <div class="tooltip-row" v-if="tooltipNode?.operator">
          <span class="tooltip-label">操作人：</span>
          <span>{{ tooltipNode.operator }}</span>
        </div>
        <div v-if="tooltipNode?.is_abnormal" class="tooltip-abnormal">
          <div class="abnormal-title">异常原因：</div>
          <div class="abnormal-reason">{{ tooltipNode.abnormal_reason }}</div>
        </div>
        <div class="tooltip-row">
          <span class="tooltip-label">时间：</span>
          <span>{{ formatTime(tooltipNode?.occurred_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { searchOrder, getProgress } from '../api'
import {
  Search,
  LocationFilled,
  Location,
  User,
  Van,
  WarningFilled,
  CircleCheckFilled,
  CircleCheck
} from '@element-plus/icons-vue'

const router = useRouter()
const searchOrderNo = ref('')
const searching = ref(false)
const searched = ref(false)
const orderData = ref(null)
const progressData = ref(null)

const tooltipVisible = ref(false)
const tooltipNode = ref(null)
const tooltipX = ref(0)
const tooltipY = ref(0)

const reversedNodes = computed(() => {
  if (!orderData.value?.nodes) return []
  return [...orderData.value.nodes].reverse()
})

const progressPercent = computed(() => {
  if (!progressData.value) return 0
  return progressData.value.progress || 0
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
    picked: '',
    in_transit: 'primary',
    arrived: 'warning',
    delivered: 'success',
    exception: 'danger'
  }
  return typeMap[status] || 'info'
}

const getNodeStatusType = (status) => {
  return getStatusType(status)
}

const getNodeClass = (node) => {
  if (node.is_abnormal) return 'node-exception'
  return `node-${node.status}`
}

const getProgressStatus = () => {
  if (orderData.value?.status === 'exception') return 'exception'
  if (orderData.value?.status === 'delivered') return 'success'
  return ''
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

const handleSearch = async () => {
  if (!searchOrderNo.value.trim()) {
    return
  }
  searching.value = true
  searched.value = true
  try {
    const res = await searchOrder(searchOrderNo.value.trim())
    orderData.value = res.data.data
    await loadProgress()
  } catch (err) {
    orderData.value = null
  } finally {
    searching.value = false
  }
}

const loadProgress = async () => {
  if (!orderData.value) return
  try {
    const res = await getProgress(orderData.value.id)
    progressData.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const loadSampleOrder = () => {
  searchOrderNo.value = 'WL202401001'
  handleSearch()
}

const goToManage = () => {
  router.push('/manage')
}

const showNodeTooltip = (node, event) => {
  tooltipNode.value = node
  tooltipX.value = event.clientX + 15
  tooltipY.value = event.clientY + 15
  tooltipVisible.value = true
}

const hideNodeTooltip = () => {
  tooltipVisible.value = false
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header {
  padding: 60px 20px 40px;
  text-align: center;
  color: white;
}

.title {
  font-size: 36px;
  font-weight: 600;
  margin-bottom: 10px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.subtitle {
  font-size: 16px;
  opacity: 0.9;
}

.main-content {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px 60px;
}

.search-section {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
}

.search-box {
  max-width: 600px;
  margin: 0 auto 20px;
}

.quick-links {
  text-align: center;
}

.divider {
  color: #dcdfe6;
  margin: 0 10px;
}

.welcome-section {
  display: flex;
  justify-content: center;
}

.welcome-card {
  background: white;
  border-radius: 16px;
  padding: 60px 40px;
  text-align: center;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.welcome-icon {
  font-size: 64px;
  color: #409eff;
  margin-bottom: 20px;
}

.welcome-card h2 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 10px;
}

.welcome-card p {
  color: #909399;
  margin-bottom: 30px;
}

.result-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order-info-card {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.order-no .label {
  color: #909399;
  font-size: 14px;
}

.order-no .value {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin-left: 8px;
}

.order-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 25px;
}

.info-item {
  display: flex;
  align-items: center;
}

.info-label {
  color: #909399;
  font-size: 14px;
  min-width: 80px;
}

.info-value {
  color: #303133;
  font-size: 14px;
}

.route-info {
  display: flex;
  align-items: center;
  margin-bottom: 25px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 12px;
}

.route-start,
.route-end {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
  flex-shrink: 0;
}

.route-icon {
  color: #409eff;
  font-size: 18px;
}

.route-line {
  flex: 1;
  height: 4px;
  background: #e4e7ed;
  margin: 0 15px;
  border-radius: 2px;
  position: relative;
}

.route-progress {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 2px;
  transition: width 0.5s ease;
}

.truck-icon {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  color: #764ba2;
  font-size: 20px;
  transition: left 0.5s ease;
}

.progress-section {
  padding-top: 10px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.progress-label {
  color: #606266;
  font-size: 14px;
}

.progress-value {
  color: #303133;
  font-weight: 600;
  font-size: 14px;
}

.timeline-card {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
}

.timeline-header h3 {
  font-size: 18px;
  color: #303133;
  margin: 0;
}

.node-count {
  color: #909399;
  font-size: 14px;
}

.timeline-container {
  padding-left: 10px;
}

.timeline-item {
  display: flex;
  gap: 20px;
  cursor: pointer;
  transition: all 0.3s;
  padding: 10px 10px 10px 0;
  border-radius: 8px;
}

.timeline-item:hover {
  background: #f5f7fa;
}

.timeline-dot {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
  width: 24px;
}

.dot {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
  transition: transform 0.3s;
}

.timeline-item:hover .dot {
  transform: scale(1.2);
}

.dot.node-pending {
  background: #909399;
  color: white;
}

.dot.node-picked {
  background: #67c23a;
  color: white;
}

.dot.node-in_transit {
  background: #409eff;
  color: white;
}

.dot.node-arrived {
  background: #e6a23c;
  color: white;
}

.dot.node-delivered {
  background: #67c23a;
  color: white;
}

.dot.node-exception {
  background: #f56c6c;
  color: white;
}

.timeline-line {
  width: 2px;
  flex: 1;
  background: #e4e7ed;
  margin-top: 6px;
  min-height: 50px;
}

.timeline-content {
  flex: 1;
  padding-bottom: 20px;
}

.timeline-time {
  color: #909399;
  font-size: 12px;
  margin-bottom: 6px;
}

.timeline-title {
  display: flex;
  align-items: center;
  margin-bottom: 6px;
}

.node-name {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.is-latest .node-name {
  color: #409eff;
}

.is-abnormal .node-name {
  color: #f56c6c;
}

.timeline-desc {
  color: #606266;
  font-size: 13px;
  margin-bottom: 8px;
  line-height: 1.5;
}

.timeline-meta {
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

.empty-section {
  background: white;
  border-radius: 16px;
  padding: 60px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.node-tooltip {
  position: fixed;
  z-index: 9999;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  min-width: 250px;
  pointer-events: none;
}

.tooltip-header {
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.tooltip-body {
  padding: 12px 16px;
}

.tooltip-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 8px;
  font-size: 13px;
  color: #606266;
}

.tooltip-row:last-child {
  margin-bottom: 0;
}

.tooltip-label {
  color: #909399;
  flex-shrink: 0;
  min-width: 60px;
}

.tooltip-abnormal {
  background: #fef0f0;
  border-radius: 6px;
  padding: 10px;
  margin: 10px 0;
}

.abnormal-title {
  color: #f56c6c;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 4px;
}

.abnormal-reason {
  color: #f56c6c;
  font-size: 12px;
  line-height: 1.5;
}
</style>
