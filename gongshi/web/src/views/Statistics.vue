<template>
  <div class="page-container">
    <div class="card">
      <div class="card-header">
        <div class="card-title">
          <el-icon style="vertical-align: middle; margin-right: 6px; color: #409eff"><DataAnalysis /></el-icon>
          数据统计
        </div>
        <div class="flex-row" style="gap: 12px">
          <el-radio-group v-model="timeRange" size="default" @change="loadAll">
            <el-radio-button value="week">本周</el-radio-button>
            <el-radio-button value="month">本月</el-radio-button>
            <el-radio-button value="quarter">本季度</el-radio-button>
            <el-radio-button value="year">本年</el-radio-button>
          </el-radio-group>
          <el-select v-if="userStore.isAdmin" v-model="filterUserId" placeholder="全部员工" clearable style="width: 140px" @change="loadAll">
            <el-option v-for="u in users" :key="u.id" :label="u.name" :value="u.id" />
          </el-select>
        </div>
      </div>

      <el-row :gutter="20" style="margin-bottom: 20px">
        <el-col :span="6">
          <div class="stat-card primary">
            <el-icon :size="32" style="color: #409eff"><Timer /></el-icon>
            <div class="stat-value">{{ overall.total_hours }}</div>
            <div class="stat-label">总工时 (小时)</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card success">
            <el-icon :size="32" style="color: #67c23a"><CircleCheck /></el-icon>
            <div class="stat-value">{{ getStatusHours('approved') }}</div>
            <div class="stat-label">已审核 (小时)</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card warning">
            <el-icon :size="32" style="color: #e6a23c"><Clock /></el-icon>
            <div class="stat-value">{{ getStatusHours('pending') }}</div>
            <div class="stat-label">待审核 (小时)</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card danger">
            <el-icon :size="32" style="color: #f56c6c"><DocumentDelete /></el-icon>
            <div class="stat-value">{{ overall.total_records }}</div>
            <div class="stat-label">填报记录 (条)</div>
          </div>
        </el-col>
      </el-row>
    </div>

    <el-row :gutter="20">
      <el-col :span="12">
        <div class="card">
          <div class="card-header">
            <div class="card-title">每日工时趋势</div>
          </div>
          <div ref="dailyChartRef" class="chart-container"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card">
          <div class="card-header">
            <div class="card-title">项目工时分布</div>
          </div>
          <div ref="projectChartRef" class="chart-container"></div>
        </div>
      </el-col>
    </el-row>

    <div class="card">
      <div class="card-header">
        <div class="card-title">项目工时明细</div>
      </div>
      <el-table :data="projectStats" border>
        <el-table-column type="index" label="序号" width="70" align="center" />
        <el-table-column prop="project_code" label="项目编号" width="120" />
        <el-table-column prop="project_name" label="项目名称" min-width="160" />
        <el-table-column prop="record_count" label="记录数" width="100" align="center" />
        <el-table-column label="工时分布" min-width="300">
          <template #default="{ row }">
            <el-progress
              :percentage="row.total > 0 ? Math.round(row.total / maxProjectHours * 100) : 0"
              :stroke-width="18"
              :text-inside="true"
              :status="row.total > 0 ? '' : ''"
            />
          </template>
        </el-table-column>
        <el-table-column prop="total" label="合计工时" width="110" align="center">
          <template #default="{ row }">
            <span style="color: #409eff; font-weight: 600">{{ row.total }} h</span>
          </template>
        </el-table-column>
        <el-table-column prop="approved" label="已通过" width="100" align="center">
          <template #default="{ row }">
            <span style="color: #67c23a">{{ row.approved }}h</span>
          </template>
        </el-table-column>
        <el-table-column prop="pending" label="待审核" width="100" align="center">
          <template #default="{ row }">
            <span style="color: #e6a23c">{{ row.pending }}h</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="90" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewProjectDetail(row)">明细</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="detailVisible" :title="dialogTitle" width="720px">
      <el-descriptions :column="2" border style="margin-bottom: 16px" v-if="currentProject">
        <el-descriptions-item label="项目编号">{{ currentProject.project_code }}</el-descriptions-item>
        <el-descriptions-item label="项目名称">{{ currentProject.project_name }}</el-descriptions-item>
        <el-descriptions-item label="总工时">
          <span style="color: #409eff; font-weight: 600">{{ currentProject.total }} 小时</span>
        </el-descriptions-item>
        <el-descriptions-item label="记录数">{{ currentProject.record_count }} 条</el-descriptions-item>
      </el-descriptions>
      <el-table :data="projectRecords" border max-height="360">
        <el-table-column prop="work_date" label="日期" width="110" sortable />
        <el-table-column prop="user.name" label="填报人" width="100" />
        <el-table-column prop="hours" label="工时" width="90" align="center">
          <template #default="{ row }">{{ row.hours }}h</template>
        </el-table-column>
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'pending'" type="warning" size="small">待审</el-tag>
            <el-tag v-else-if="row.status === 'approved'" type="success" size="small">通过</el-tag>
            <el-tag v-else type="danger" size="small">驳回</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import { useUserStore } from '../stores/user'
import { getUsers } from '../api/auth'
import { getWorkRecords, getDailyStats, getProjectStats, getOverallStats } from '../api/workrecord'

const userStore = useUserStore()
const timeRange = ref('month')
const filterUserId = ref(null)
const users = ref([])
const dailyStats = ref([])
const projectStats = ref([])
const overall = reactive({ total_hours: 0, total_records: 0, by_status: [] })
const dailyChartRef = ref()
const projectChartRef = ref()
const detailVisible = ref(false)
const currentProject = ref(null)
const projectRecords = ref([])
const dialogTitle = ref('')

let dailyChart = null
let projectChart = null

const maxProjectHours = computed(() => Math.max(...(projectStats.value || []).map(p => p.total), 1))

const getStatusHours = (status) => {
  const list = overall.by_status || []
  const item = list.find(s => s.status === status)
  return item?.hours || 0
}

const rangeParams = computed(() => {
  let start, end
  const now = dayjs()
  switch (timeRange.value) {
    case 'week':
      start = now.startOf('week').format('YYYY-MM-DD')
      end = now.endOf('week').format('YYYY-MM-DD')
      break
    case 'quarter':
      start = now.startOf('quarter').format('YYYY-MM-DD')
      end = now.endOf('quarter').format('YYYY-MM-DD')
      break
    case 'year':
      start = now.startOf('year').format('YYYY-MM-DD')
      end = now.endOf('year').format('YYYY-MM-DD')
      break
    default:
      start = now.startOf('month').format('YYYY-MM-DD')
      end = now.endOf('month').format('YYYY-MM-DD')
  }
  return { start_date: start, end_date: end }
})

const initCharts = () => {
  if (dailyChartRef.value) {
    dailyChart = echarts.init(dailyChartRef.value)
  }
  if (projectChartRef.value) {
    projectChart = echarts.init(projectChartRef.value)
  }
  window.addEventListener('resize', () => {
    dailyChart?.resize()
    projectChart?.resize()
  })
}

const renderDailyChart = () => {
  if (!dailyChart) return
  const stats = dailyStats.value || []
  const dates = []
  const approved = []
  const pending = []
  const rejected = []
  const { start_date, end_date } = rangeParams.value
  let start = dayjs(start_date)
  const end = dayjs(end_date)

  while (start.isBefore(end) || start.isSame(end, 'day')) {
    const ds = start.format('YYYY-MM-DD')
    dates.push(start.format('MM-DD'))
    const item = stats.find(d => d.work_date === ds)
    approved.push(item?.approved || 0)
    pending.push(item?.pending || 0)
    rejected.push(item?.rejected || 0)
    start = start.add(1, 'day')
  }

  dailyChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['已通过', '待审核', '已驳回'], top: 0 },
    grid: { left: 50, right: 20, top: 40, bottom: 70 },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        fontSize: 10,
        rotate: dates.length > 15 ? 60 : 0,
        interval: dates.length > 30 ? Math.floor(dates.length / 12) : 0,
        margin: 12
      }
    },
    yAxis: { type: 'value', name: '小时', nameLocation: 'middle', nameGap: 30 },
    series: [
      { name: '已通过', type: 'bar', stack: 'total', data: approved, itemStyle: { color: '#67c23a' } },
      { name: '待审核', type: 'bar', stack: 'total', data: pending, itemStyle: { color: '#e6a23c' } },
      { name: '已驳回', type: 'bar', stack: 'total', data: rejected, itemStyle: { color: '#f56c6c' } }
    ]
  })
}

const renderProjectChart = () => {
  if (!projectChart) return
  const data = (projectStats.value || []).map(p => ({
    name: p.project_name,
    value: p.total
  }))
  projectChart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} 小时 ({d}%)'
    },
    legend: { orient: 'vertical', right: 10, top: 'center' },
    series: [{
      type: 'pie',
      radius: ['40%', '65%'],
      center: ['35%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 6,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: true,
        formatter: '{b}: {d}%'
      },
      data: data,
      color: ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#9254de', '#13c2c2']
    }]
  })
}

const loadDaily = async () => {
  const params = { ...rangeParams.value }
  if (filterUserId.value) params.user_id = filterUserId.value
  const res = await getDailyStats(params)
  dailyStats.value = Array.isArray(res) ? res : []
}

const loadProject = async () => {
  const params = { ...rangeParams.value }
  if (filterUserId.value) params.user_id = filterUserId.value
  const res = await getProjectStats(params)
  projectStats.value = Array.isArray(res) ? res : []
}

const loadOverall = async () => {
  const params = { ...rangeParams.value }
  if (filterUserId.value) params.user_id = filterUserId.value
  const res = await getOverallStats(params) || {}
  overall.total_hours = res.total_hours || 0
  overall.total_records = res.total_records || 0
  overall.by_status = Array.isArray(res.by_status) ? res.by_status : []
}

const loadAll = async () => {
  await Promise.all([loadDaily(), loadProject(), loadOverall()])
  await nextTick()
  renderDailyChart()
  renderProjectChart()
}

const viewProjectDetail = async (row) => {
  currentProject.value = row
  dialogTitle.value = `${row.project_name} - 工时明细`
  const params = { ...rangeParams.value, project_id: row.project_id }
  if (filterUserId.value) params.user_id = filterUserId.value
  projectRecords.value = await getWorkRecords(params)
  detailVisible.value = true
}

onMounted(async () => {
  if (userStore.isAdmin) {
    users.value = await getUsers()
  }
  await nextTick()
  initCharts()
  loadAll()
})
</script>
