
&lt;template&gt;
  &lt;div class="stats"&gt;
    &lt;div class="cards"&gt;
      &lt;div class="card"&gt;
        &lt;div class="label"&gt;今日预订数&lt;/div&gt;
        &lt;div class="value"&gt;{{ stats.today_bookings || 0 }}&lt;/div&gt;
      &lt;/div&gt;
      &lt;div class="card"&gt;
        &lt;div class="label"&gt;启用会议室&lt;/div&gt;
        &lt;div class="value"&gt;{{ stats.active_rooms || 0 }}&lt;/div&gt;
      &lt;/div&gt;
    &lt;/div&gt;

    &lt;div class="charts"&gt;
      &lt;div class="chart-box"&gt;
        &lt;h3&gt;本周每日预订量&lt;/h3&gt;
        &lt;div ref="dailyChart" style="width: 100%; height: 350px"&gt;&lt;/div&gt;
      &lt;/div&gt;
      &lt;div class="chart-box"&gt;
        &lt;h3&gt;会议室本周预订排行&lt;/h3&gt;
        &lt;div ref="roomChart" style="width: 100%; height: 350px"&gt;&lt;/div&gt;
      &lt;/div&gt;
    &lt;/div&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { ref, onMounted, nextTick } from 'vue'
import { useAuthStore } from '../stores/auth'
import * as echarts from 'echarts'

const authStore = useAuthStore()
const stats = ref({})
const dailyChart = ref(null)
const roomChart = ref(null)
let dailyChartInstance = null
let roomChartInstance = null

const fetchStats = async () =&gt; {
  const res = await authStore.api.get('/stats')
  stats.value = res.data
  nextTick(() =&gt; {
    renderCharts()
  })
}

const renderCharts = () =&gt; {
  if (dailyChart.value) {
    dailyChartInstance = echarts.init(dailyChart.value)
    const days = (stats.value.daily_stats || []).map(d =&gt; d.date.slice(5))
    dailyChartInstance.setOption({
      backgroundColor: 'transparent',
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: days, axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' } },
      yAxis: { type: 'value', axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' }, splitLine: { lineStyle: { color: '#334155' } } },
      series: [{ type: 'bar', data: (stats.value.daily_stats || []).map(d =&gt; d.count), itemStyle: { color: '#3b82f6', borderRadius: [4, 4, 0, 0] } }]
    })
  }

  if (roomChart.value) {
    roomChartInstance = echarts.init(roomChart.value)
    const sortedRooms = (stats.value.room_stats || []).sort((a, b) =&gt; b.count - a.count)
    roomChartInstance.setOption({
      backgroundColor: 'transparent',
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'value', axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' }, splitLine: { lineStyle: { color: '#334155' } } },
      yAxis: { type: 'category', data: sortedRooms.map(r =&gt; r.room_name), axisLine: { lineStyle: { color: '#475569' } }, axisLabel: { color: '#cbd5e1' } },
      series: [{ type: 'bar', data: sortedRooms.map(r =&gt; r.count), itemStyle: { color: '#8b5cf6', borderRadius: [0, 4, 4, 0] } }]
    })
  }
}

onMounted(() =&gt; {
  fetchStats()
  window.addEventListener('resize', () =&gt; {
    dailyChartInstance?.resize()
    roomChartInstance?.resize()
  })
})
&lt;/script&gt;

&lt;style scoped&gt;
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
&lt;/style&gt;

