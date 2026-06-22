<template>
  <div class="order-detail" v-loading="loading">
    <template v-if="order">
      <el-page-header @back="$router.back()" class="page-header">
        <template #content>
          <div class="header-content">
            <span>工单详情 - {{ order.order_no }}</span>
            <el-tag :type="statusType(order.status)" effect="dark" style="margin-left: 10px;">
              {{ statusText(order.status) }}
            </el-tag>
            <el-tag v-if="order.is_overdue" type="danger" effect="plain" style="margin-left: 8px;">
              <el-icon><Warning /></el-icon>已逾期
            </el-tag>
          </div>
        </template>
      </el-page-header>

      <el-row :gutter="20" class="content-row">
        <el-col :span="16">
          <el-card shadow="never" class="info-card">
            <template #header><span>基本信息</span></template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="工单标题">{{ order.title }}</el-descriptions-item>
              <el-descriptions-item label="优先级">
                <el-tag :type="priorityType(order.priority)">{{ priorityText(order.priority) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="设备类型">{{ deviceTypeText(order.device_type) }}</el-descriptions-item>
              <el-descriptions-item label="设备名称">{{ order.device_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="报修区域">{{ order.area }}</el-descriptions-item>
              <el-descriptions-item label="具体位置">{{ order.location || '-' }}</el-descriptions-item>
              <el-descriptions-item label="联系人">{{ order.contact_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="联系电话">{{ order.contact_phone || '-' }}</el-descriptions-item>
              <el-descriptions-item label="报修人">{{ order.user?.real_name || order.user?.username }}</el-descriptions-item>
              <el-descriptions-item label="运维人员">{{ order.technician?.real_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="提交时间">{{ formatTime(order.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="期望完成">
                <span :class="{ 'text-danger': order.is_overdue }">{{ formatTime(order.expect_time) }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="派单时间" v-if="order.assign_time">{{ formatTime(order.assign_time) }}</el-descriptions-item>
              <el-descriptions-item label="开始处理" v-if="order.start_process_time">{{ formatTime(order.start_process_time) }}</el-descriptions-item>
              <el-descriptions-item label="完成时间" v-if="order.complete_time">{{ formatTime(order.complete_time) }}</el-descriptions-item>
              <el-descriptions-item label="驳回原因" v-if="order.status === 'rejected'">
                <span class="text-danger">{{ order.reject_reason }}</span>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <el-card shadow="never" class="info-card">
            <template #header><span>故障描述</span></template>
            <p class="description">{{ order.description }}</p>
            <div v-if="order.fault_images" class="images">
              <el-image
                v-for="(img, idx) in parseImages(order.fault_images)"
                :key="idx"
                :src="img"
                :preview-src-list="parseImages(order.fault_images)"
                fit="cover"
                class="fault-image"
              />
            </div>
          </el-card>

          <el-card shadow="never" class="info-card">
            <template #header>
              <div style="display: flex; justify-content: space-between; align-items: center;">
                <span>维修进度</span>
                <div v-if="canProcess">
                  <el-button v-if="order.status === 'assigned'" type="primary" size="small" @click="handleStartProcess">
                    开始处理
                  </el-button>
                  <el-button v-if="order.status === 'processing'" type="primary" size="small" @click="showProcessDialog = true">
                    提交进度
                  </el-button>
                  <el-button v-if="order.status === 'processing'" type="success" size="small" @click="showCompleteDialog = true">
                    完成维修
                  </el-button>
                </div>
                <el-button
                  v-if="userStore.isAdmin && order.status !== 'completed' && order.status !== 'cancelled'"
                  type="danger"
                  size="small"
                  @click="showRejectDialog = true"
                >驳回工单</el-button>
              </div>
            </template>

            <el-steps :active="stepActive" finish-status="success" class="timeline">
              <el-step title="提交工单" :description="formatTime(order.created_at)" />
              <el-step title="派单" :description="formatTime(order.assign_time)" />
              <el-step title="开始维修" :description="formatTime(order.start_process_time)" />
              <el-step title="完成维修" :description="formatTime(order.complete_time)" />
            </el-steps>

            <el-timeline v-if="logs && logs.length" class="logs-timeline">
              <el-timeline-item
                v-for="log in logs"
                :key="log.id"
                :timestamp="formatTime(log.created_at)"
                placement="top"
                :type="log.status === 'completed' ? 'success' : 'primary'"
              >
                <el-card shadow="never" class="log-card">
                  <div class="log-header">
                    <el-tag size="small" type="primary">
                      {{ log.technician?.real_name || '系统' }}
                    </el-tag>
                    <el-tag v-if="log.status" size="small" :type="statusType(log.status)" style="margin-left: 8px;">
                      {{ statusText(log.status) }}
                    </el-tag>
                  </div>
                  <p class="log-content">{{ log.content }}</p>
                  <div v-if="log.images" class="log-images">
                    <el-image
                      v-for="(img, idx) in parseImages(log.images)"
                      :key="idx"
                      :src="img"
                      :preview-src-list="parseImages(log.images)"
                      fit="cover"
                      class="log-image"
                    />
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
            <el-empty v-else description="暂无维修记录" :image-size="80" />
          </el-card>

          <el-card v-if="order.status === 'completed' && userStore.user?.id === order.user_id" shadow="never" class="info-card">
            <template #header><span>工单评价</span></template>
            <div v-if="evaluation">
              <el-rate v-model="evaluation.rating" disabled show-score text-color="#ff9900" />
              <p class="eval-content" v-if="evaluation.content">{{ evaluation.content }}</p>
            </div>
            <el-form v-else label-width="80px">
              <el-form-item label="评分">
                <el-rate v-model="evalForm.rating" />
              </el-form-item>
              <el-form-item label="评价">
                <el-input v-model="evalForm.content" type="textarea" :rows="3" placeholder="请填写评价内容" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleEvaluate" :loading="evalLoading">提交评价</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card shadow="never" class="side-card">
            <template #header><span>快捷操作</span></template>
            <div class="actions">
              <el-button
                v-if="userStore.isAdmin && (order.status === 'pending' || order.status === 'rejected')"
                type="primary"
                block
                @click="showAssignDialog = true"
              >派单给运维</el-button>
              <el-button
                v-if="userStore.isAdmin && (order.status === 'pending' || order.status === 'rejected')"
                type="warning"
                block
                @click="handleAutoAssign"
              >智能派单</el-button>
              <el-button
                v-if="userStore.user?.role !== 'technician' && (order.status === 'pending' || order.status === 'assigned') && order.user_id === userStore.user?.id"
                type="danger"
                block
                @click="handleCancel"
              >取消工单</el-button>
            </div>
          </el-card>

          <el-card shadow="never" class="side-card">
            <template #header><span>统计信息</span></template>
            <el-statistic title="处理时长" :value="calculateDuration()" suffix="小时" />
            <el-divider />
            <el-statistic v-if="evaluation" title="用户评分">
              <template #default>
                <el-rate :model-value="evaluation.rating" disabled show-score text-color="#ff9900" />
              </template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>

      <el-dialog v-model="showAssignDialog" title="派单" width="450px">
        <el-form label-width="100px">
          <el-form-item label="运维人员">
            <el-select v-model="assignTechId" placeholder="请选择运维人员" style="width: 100%;">
              <el-option v-for="tech in technicians" :key="tech.id" :label="`${tech.real_name} (${tech.area})`" :value="tech.id" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showAssignDialog = false">取消</el-button>
          <el-button type="primary" @click="handleAssign">确认派单</el-button>
        </template>
      </el-dialog>

      <el-dialog v-model="showProcessDialog" title="提交维修进度" width="500px">
        <el-form :model="processForm" label-width="100px">
          <el-form-item label="当前状态">
            <el-radio-group v-model="processForm.status">
              <el-radio value="processing">继续维修</el-radio>
              <el-radio value="completed">维修完成</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="进度说明">
            <el-input v-model="processForm.content" type="textarea" :rows="4" placeholder="请描述当前维修进度" />
          </el-form-item>
          <el-form-item label="上传凭证">
            <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              list-type="picture-card"
              :on-success="handleUploadSuccess"
              :on-remove="handleRemove"
              :file-list="processFileList"
              multiple
              accept="image/*"
            >
              <el-icon><Plus /></el-icon>
            </el-upload>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showProcessDialog = false">取消</el-button>
          <el-button type="primary" @click="handleProcess" :loading="processLoading">提交</el-button>
        </template>
      </el-dialog>

      <el-dialog v-model="showCompleteDialog" title="完成维修" width="500px">
        <el-form :model="completeForm" label-width="100px">
          <el-form-item label="维修说明">
            <el-input v-model="completeForm.content" type="textarea" :rows="4" placeholder="请描述维修结果" />
          </el-form-item>
          <el-form-item label="上传凭证">
            <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              list-type="picture-card"
              :on-success="handleCompleteUploadSuccess"
              :on-remove="handleCompleteRemove"
              :file-list="completeFileList"
              multiple
              accept="image/*"
            >
              <el-icon><Plus /></el-icon>
            </el-upload>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showCompleteDialog = false">取消</el-button>
          <el-button type="primary" @click="handleComplete" :loading="completeLoading">确认完成</el-button>
        </template>
      </el-dialog>

      <el-dialog v-model="showRejectDialog" title="驳回工单" width="450px">
        <el-form label-width="100px">
          <el-form-item label="驳回原因">
            <el-input v-model="rejectReason" type="textarea" :rows="3" placeholder="请填写驳回原因" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showRejectDialog = false">取消</el-button>
          <el-button type="danger" @click="handleReject">确认驳回</el-button>
        </template>
      </el-dialog>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Warning } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getOrderDetail, getTechnicians, assignOrder, autoAssignOrder,
  startProcessOrder, processOrder, completeOrder, rejectOrder,
  cancelOrder, evaluateOrder
} from '@/api'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const order = ref(null)
const logs = ref([])
const evaluation = ref(null)
const technicians = ref([])
const showAssignDialog = ref(false)
const showProcessDialog = ref(false)
const showCompleteDialog = ref(false)
const showRejectDialog = ref(false)
const assignTechId = ref(null)
const rejectReason = ref('')
const evalLoading = ref(false)
const processLoading = ref(false)
const completeLoading = ref(false)
const processFileList = ref([])
const completeFileList = ref([])

const evalForm = reactive({ rating: 5, content: '' })
const processForm = reactive({ content: '', status: 'processing' })
const completeForm = reactive({ content: '' })

const uploadUrl = computed(() => '/api/upload')
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))

const canProcess = computed(() => {
  if (!order.value || !userStore.user) return false
  const isTech = userStore.user.id === order.value.technician_id && userStore.userRole === 'technician'
  const isAdmin = userStore.isAdmin
  return (isTech || isAdmin) && ['assigned', 'processing'].includes(order.value.status)
})

const stepActive = computed(() => {
  if (!order.value) return 0
  const s = order.value.status
  if (s === 'completed') return 4
  if (s === 'processing' || order.value.start_process_time) return 3
  if (s === 'assigned' || order.value.assign_time) return 2
  if (s === 'pending' || s === 'rejected' || s === 'cancelled') return 1
  return 0
})

const loadDetail = async () => {
  loading.value = true
  try {
    const res = await getOrderDetail(route.params.id)
    order.value = res.order
    logs.value = res.logs
    evaluation.value = res.evaluation
  } finally {
    loading.value = false
  }
}

const loadTechnicians = async () => {
  try {
    technicians.value = await getTechnicians()
  } catch (e) {}
}

const handleAssign = async () => {
  if (!assignTechId.value) {
    ElMessage.warning('请选择运维人员')
    return
  }
  await assignOrder(order.value.id, { technician_id: assignTechId.value })
  ElMessage.success('派单成功')
  showAssignDialog.value = false
  loadDetail()
}

const handleAutoAssign = async () => {
  await autoAssignOrder(order.value.id)
  ElMessage.success('智能派单成功')
  loadDetail()
}

const handleStartProcess = async () => {
  await startProcessOrder(order.value.id)
  ElMessage.success('已开始处理')
  loadDetail()
}

const handleUploadSuccess = (res, file) => {
  processFileList.value.push({ url: res.url, name: file.name, uid: file.uid })
}

const handleRemove = (file) => {
  const idx = processFileList.value.findIndex(f => f.uid === file.uid)
  if (idx > -1) processFileList.value.splice(idx, 1)
}

const handleCompleteUploadSuccess = (res, file) => {
  completeFileList.value.push({ url: res.url, name: file.name, uid: file.uid })
}

const handleCompleteRemove = (file) => {
  const idx = completeFileList.value.findIndex(f => f.uid === file.uid)
  if (idx > -1) completeFileList.value.splice(idx, 1)
}

const handleProcess = async () => {
  if (!processForm.content) {
    ElMessage.warning('请填写进度说明')
    return
  }
  processLoading.value = true
  try {
    await processOrder(order.value.id, {
      content: processForm.content,
      status: processForm.status,
      images: processFileList.value.map(f => f.url).join(',')
    })
    ElMessage.success('提交成功')
    showProcessDialog.value = false
    processForm.content = ''
    processForm.status = 'processing'
    processFileList.value = []
    loadDetail()
  } finally {
    processLoading.value = false
  }
}

const handleComplete = async () => {
  if (!completeForm.content) {
    ElMessage.warning('请填写维修说明')
    return
  }
  completeLoading.value = true
  try {
    await completeOrder(order.value.id, {
      content: completeForm.content,
      images: completeFileList.value.map(f => f.url).join(',')
    })
    ElMessage.success('工单已完成')
    showCompleteDialog.value = false
    completeForm.content = ''
    completeFileList.value = []
    loadDetail()
  } finally {
    completeLoading.value = false
  }
}

const handleReject = async () => {
  if (!rejectReason.value) {
    ElMessage.warning('请填写驳回原因')
    return
  }
  await rejectOrder(order.value.id, { reason: rejectReason.value })
  ElMessage.success('工单已驳回')
  showRejectDialog.value = false
  rejectReason.value = ''
  loadDetail()
}

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定取消此工单吗？', '提示', { type: 'warning' })
    await cancelOrder(order.value.id)
    ElMessage.success('工单已取消')
    loadDetail()
  } catch (e) {}
}

const handleEvaluate = async () => {
  if (!evalForm.rating) {
    ElMessage.warning('请选择评分')
    return
  }
  evalLoading.value = true
  try {
    await evaluateOrder(order.value.id, evalForm)
    ElMessage.success('评价成功')
    loadDetail()
  } finally {
    evalLoading.value = false
  }
}

const parseImages = (str) => str ? str.split(',').filter(Boolean) : []
const formatTime = (t) => t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'

const calculateDuration = () => {
  if (!order.value) return 0
  const start = order.value.created_at
  const end = order.value.complete_time || new Date().toISOString()
  return dayjs(end).diff(dayjs(start), 'hour', true).toFixed(1)
}

const statusText = (s) => ({
  pending: '待处理', assigned: '已派单', processing: '维修中',
  completed: '已完成', rejected: '已驳回', cancelled: '已取消'
}[s] || s)

const statusType = (s) => ({
  pending: 'info', assigned: 'warning', processing: 'primary',
  completed: 'success', rejected: 'danger', cancelled: 'info'
}[s] || '')

const priorityText = (p) => ({ low: '低', medium: '中', high: '高', urgent: '紧急' }[p] || p)
const priorityType = (p) => ({ low: 'info', medium: '', high: 'warning', urgent: 'danger' }[p] || '')
const deviceTypeText = (t) => ({
  computer: '电脑', printer: '打印机', network: '网络设备',
  projector: '投影仪', aircon: '空调', other: '其他'
}[t] || t)

onMounted(() => {
  loadDetail()
  if (userStore.isAdmin) loadTechnicians()
})
</script>

<style scoped>
.page-header {
  background: #fff;
  border-radius: 4px;
  margin-bottom: 16px;
  padding: 12px 20px;
}

.header-content {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.content-row {
  min-height: calc(100vh - 160px);
}

.info-card {
  margin-bottom: 16px;
}

.side-card {
  margin-bottom: 16px;
}

.description {
  line-height: 1.8;
  color: #333;
  white-space: pre-wrap;
}

.images {
  margin-top: 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.fault-image {
  width: 140px;
  height: 140px;
  border-radius: 4px;
  cursor: pointer;
}

.timeline {
  margin: 20px 0;
}

.logs-timeline {
  margin-top: 30px;
}

.log-card {
  border: 1px solid #ebeef5;
}

.log-header {
  margin-bottom: 8px;
}

.log-content {
  margin: 0;
  color: #333;
  line-height: 1.6;
  white-space: pre-wrap;
}

.log-images {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.log-image {
  width: 100px;
  height: 100px;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.text-danger {
  color: #f56c6c;
  font-weight: 600;
}

.eval-content {
  margin-top: 10px;
  color: #666;
  line-height: 1.6;
}
</style>
