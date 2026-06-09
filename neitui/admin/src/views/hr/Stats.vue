<template>
  <div>
    <el-row :gutter="20" style="margin-bottom: 20px">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num">{{ stats.total_referrals }}</div>
            <div class="stat-label">总内推数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num" style="color: #67C23A">{{ stats.hired }}</div>
            <div class="stat-label">已入职</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num" style="color: #E6A23C">{{ stats.interview_rate?.toFixed(1) }}%</div>
            <div class="stat-label">面试通过率</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-num" style="color: #F56C6C">{{ stats.average_score?.toFixed(1) || '暂无' }}</div>
            <div class="stat-label">平均评分</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-bottom: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>员工内推成功排行</span>
          </template>
          <div ref="barChartRef" style="height: 400px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>员工平均评分排行</span>
          </template>
          <div ref="scoreBarChartRef" style="height: 400px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>近30天内推趋势</span>
          </template>
          <div ref="lineChartRef" style="height: 400px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import api from '@/api'
import * as echarts from 'echarts'

const stats = ref({})
const barChartRef = ref(null)
const scoreBarChartRef = ref(null)
const lineChartRef = ref(null)
let barChart = null
let scoreBarChart = null
let lineChart = null

const loadStats = async () => {
  const res = await api.get('/admin/stats')
  stats.value = res
  nextTick(() => {
    renderBarChart(res.top_employees || [])
    renderScoreBarChart(res.employee_score_ranking || [])
    renderLineChart(res.thirty_days_trend || [])
  })
}

const renderBarChart = (data) => {
  if (!barChartRef.value) return
  if (barChart) barChart.dispose()
  barChart = echarts.init(barChartRef.value)
  barChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: data.map(d => d.name) },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: data.map(d => d.count), itemStyle: { color: '#409EFF' } }]
  })
}

const renderScoreBarChart = (data) => {
  if (!scoreBarChartRef.value) return
  if (scoreBarChart) scoreBarChart.dispose()
  scoreBarChart = echarts.init(scoreBarChartRef.value)
  scoreBarChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: data.map(d => d.name) },
    yAxis: { type: 'value', max: 5 },
    series: [
      { 
        type: 'bar', 
        data: data.map(d => d.avg_score.toFixed(1)), 
        itemStyle: { color: '#F56C6C' },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}分'
        }
      }
    ]
  })
}

const renderLineChart = (data) => {
  if (!lineChartRef.value) return
  if (lineChart) lineChart.dispose()
  lineChart = echarts.init(lineChartRef.value)
  lineChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: data.map(d => d.date) },
    yAxis: { type: 'value' },
    series: [{ type: 'line', data: data.map(d => d.count), smooth: true, itemStyle: { color: '#67C23A' } }]
  })
}

onMounted(() => {
  loadStats()
  window.addEventListener('resize', () => {
    barChart?.resize()
    scoreBarChart?.resize()
    lineChart?.resize()
  })
})
</script>

<style scoped>
.stat-item {
  text-align: center;
}
.stat-num {
  font-size: 36px;
  font-weight: bold;
  color: #409EFF;
}
.stat-label {
  color: #999;
  margin-top: 10px;
  font-size: 16px;
}
</style>
