<template>
  <div class="stats">
    <div class="cards">
      <div class="card">
        <div class="label">今日预订数</div>
        <div class="value">{{ stats.today_bookings || 0 }}</div>
      </div>
      <div class="card">
        <div class="label">启用会议室</div>
        <div class="value">{{ stats.active_rooms || 0 }}</div>
      </div>
    </div>

    <div class="charts">
      <div class="chart-box">
        <h3>本周每日预订量</h3>
        <div ref="dailyChart" style="width: 100%; height: 350px"></div>
      </div>
      <div class="chart-box">
        <h3>会议室本周预订排行</h3>
        <div ref="roomChart" style="width: 100%; height: 350px"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import api from '../api'

const stats = ref({})
const dailyChart = ref(null)
const roomChart = ref(null)
let dailyChartInstance = null
let roomChartInstance = null

const fetchStats = async () => {
  try {
    const res = await api.get('/stats')
    stats.value = res.data
    nextTick(() => {
      renderCharts()
    })
  } catch (e) {
    console.error(e)
  }
}

const renderCharts = () => {
  if (dailyChart.value) {
    dailyChartInstance = echarts.init(dailyChart.value)
    const days = (stats.value.daily_stats || []).map(d => d.date.slice(5))
    dailyChartInstance.setOption({
      backgroundColor: 'transparent',
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: days, axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' } },
      yAxis: { type: 'value', axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' }, splitLine: { lineStyle: { color: '#334155' } } },
      series: [{ type: 'bar', data: (stats.value.daily_stats || []).map(d => d.count), itemStyle: { color: '#3b82f6', borderRadius: [4, 4, 0, 0] } }]
    })
  }

  if (roomChart.value) {
    roomChartInstance = echarts.init(roomChart.value)
    const sortedRooms = (stats.value.room_stats || []).sort((a, b) => b.count - a.count)
    roomChartInstance.setOption({
      backgroundColor: 'transparent',
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'value', axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' }, splitLine: { lineStyle: { color: '#334155' } } },
      yAxis: { type: 'category', data: sortedRooms.map(r => r.room_name), axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' } },
      series: [{ type: 'bar', data: sortedRooms.map(r => r.count), itemStyle: { color: '#8b5cf6', borderRadius: [0, 4, 4, 0] } }]
    })
  }
}

onMounted(() => {
  fetchStats()
  window.addEventListener('resize', () => {
    dailyChartInstance?.resize()
    roomChartInstance?.resize()
  })
})
</script>

<style scoped>
.stats {
  padding: 0;
}

.cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 24px;
}

.card {
  background: linear-gradient(135deg, #1e293b, #334155);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid #475569;
}

.card .label {
  color: #94a3b8;
  font-size: 14px;
  margin-bottom: 8px;
}

.card .value {
  color: #e2e8f0;
  font-size: 36px;
  font-weight: bold;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.charts {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.chart-box {
  background: #1e293b;
  border-radius: 12px;
  padding: 24px;
  border: 1px solid #475569;
}

.chart-box h3 {
  color: #e2e8f0;
  font-size: 16px;
  margin-bottom: 16px;
}
</style>
