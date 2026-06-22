<template>
  <div class="dashboard" v-loading="loading">
    <el-row :gutter="20">
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon" style="background: #409EFF;">
              <el-icon size="28"><Tickets /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.overview?.total || 0 }}</div>
              <div class="stat-label">总工单数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon" style="background: #E6A23C;">
              <el-icon size="28"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.overview?.pending || 0 }}</div>
              <div class="stat-label">待处理</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon" style="background: #409EFF;">
              <el-icon size="28"><Loading /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.overview?.processing || 0 }}</div>
              <div class="stat-label">处理中</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon" style="background: #67C23A;">
              <el-icon size="28"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.overview?.completed || 0 }}</div>
              <div class="stat-label">已完成</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon" style="background: #F56C6C;">
              <el-icon size="28"><Warning /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.overview?.overdue || 0 }}</div>
              <div class="stat-label">已逾期</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-item">
            <div class="stat-icon" style="background: #909399;">
              <el-icon size="28"><Star /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ (stats.overview?.avg_rating || 0).toFixed(1) }}</div>
              <div class="stat-label">平均评分</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>工单状态分布</span></template>
          <div ref="statusChart" style="height: 320px;"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header><span>近7日工单趋势</span></template>
          <div ref="trendChart" style="height: 320px;"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="8">
        <el-card shadow="never">
          <template #header><span>故障类型统计</span></template>
          <div ref="typeChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <template #header><span>区域报修统计</span></template>
          <div ref="areaChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <template #header><span>运维效率统计</span></template>
          <el-table :data="stats.technician_stats || []" size="small" stripe>
            <el-table-column prop="real_name" label="运维人员" />
            <el-table-column prop="total_orders" label="总工单" width="70" align="center" />
            <el-table-column prop="completed" label="已完成" width="70" align="center" />
            <el-table-column label="平均耗时" width="90" align="center">
              <template #default="{ row }">{{ row.avg_duration?.toFixed(1) }}h</template>
            </el-table-column>
            <el-table-column label="评分" width="80" align="center">
              <template #default="{ row }">
                <el-rate :model-value="row.avg_rating" disabled size="small" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { getStatistics } from '@/api'

const loading = ref(false)
const stats = ref({})
const statusChart = ref(null)
const trendChart = ref(null)
const typeChart = ref(null)
const areaChart = ref(null)

let statusChartInstance = null
let trendChartInstance = null
let typeChartInstance = null
let areaChartInstance = null

const statusTextMap = {
  pending: '待处理', assigned: '已派单', processing: '维修中',
  completed: '已完成', rejected: '已驳回', cancelled: '已取消'
}

const deviceTypeTextMap = {
  computer: '电脑', printer: '打印机', network: '网络设备',
  projector: '投影仪', aircon: '空调', other: '其他'
}

const loadStats = async () => {
  loading.value = true
  try {
    stats.value = await getStatistics()
    await nextTick()
    initCharts()
  } finally {
    loading.value = false
  }
}

const initCharts = () => {
  if (statusChart.value) {
    statusChartInstance = echarts.init(statusChart.value)
    const data = (stats.value.status_stats || []).map(s => ({
      name: statusTextMap[s.status] || s.status,
      value: s.count
    }))
    statusChartInstance.setOption({
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', right: 10, top: 'center' },
      series: [{
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['40%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
        label: { show: false },
        data,
        color: ['#909399', '#E6A23C', '#409EFF', '#67C23A', '#F56C6C', '#C0C4CC']
      }]
    })
  }

  if (trendChart.value) {
    trendChartInstance = echarts.init(trendChart.value)
    const dailyStats = stats.value.daily_stats || []
    trendChartInstance.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: 40, right: 20, top: 30, bottom: 40 },
      xAxis: {
        type: 'category',
        data: dailyStats.map(d => d.date.slice(5)),
        axisLine: { lineStyle: { color: '#ccc' } }
      },
      yAxis: { type: 'value', minInterval: 1 },
      series: [{
        type: 'line',
        data: dailyStats.map(d => d.count),
        smooth: true,
        areaStyle: { color: 'rgba(64,158,255,0.2)' },
        lineStyle: { color: '#409EFF', width: 3 },
        itemStyle: { color: '#409EFF' }
      }]
    })
  }

  if (typeChart.value) {
    typeChartInstance = echarts.init(typeChart.value)
    const data = (stats.value.type_stats || []).map(t => ({
      name: deviceTypeTextMap[t.device_type] || t.device_type,
      value: t.count
    }))
    typeChartInstance.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: 80, right: 20, top: 20, bottom: 30 },
      xAxis: { type: 'value' },
      yAxis: {
        type: 'category',
        data: data.map(d => d.name),
        axisLine: { lineStyle: { color: '#ccc' } }
      },
      series: [{
        type: 'bar',
        data: data.map(d => d.value),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: '#667eea' },
            { offset: 1, color: '#764ba2' }
          ]),
          borderRadius: [0, 4, 4, 0]
        }
      }]
    })
  }

  if (areaChart.value) {
    areaChartInstance = echarts.init(areaChart.value)
    const data = stats.value.area_stats || []
    areaChartInstance.setOption({
      tooltip: { trigger: 'item' },
      series: [{
        type: 'pie',
        radius: '65%',
        data: data.map(d => ({ name: d.area || '其他', value: d.count })),
        emphasis: {
          itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.5)' }
        },
        label: { formatter: '{b}\n{d}%' }
      }],
      color: ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#9b59b6']
    })
  }
}

const handleResize = () => {
  statusChartInstance?.resize()
  trendChartInstance?.resize()
  typeChartInstance?.resize()
  areaChartInstance?.resize()
}

onMounted(() => {
  loadStats()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  statusChartInstance?.dispose()
  trendChartInstance?.dispose()
  typeChartInstance?.dispose()
  areaChartInstance?.dispose()
})
</script>

<style scoped>
.stat-card {
  text-align: center;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-info {
  flex: 1;
  text-align: left;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}
</style>
