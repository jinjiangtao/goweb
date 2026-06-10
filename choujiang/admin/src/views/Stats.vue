
&lt;template&gt;
  &lt;div&gt;
    &lt;el-row :gutter="20" style="margin-bottom: 20px"&gt;
      &lt;el-col :span="6"&gt;
        &lt;el-card&gt;
          &lt;div style="text-align: center"&gt;
            &lt;div style="font-size: 32px; font-weight: bold; color: #409EFF"&gt;{{ stats.totalDraws || 0 }}&lt;/div&gt;
            &lt;div style="color: #909399; margin-top: 10px"&gt;总抽奖次数&lt;/div&gt;
          &lt;/div&gt;
        &lt;/el-card&gt;
      &lt;/el-col&gt;
      &lt;el-col :span="6"&gt;
        &lt;el-card&gt;
          &lt;div style="text-align: center"&gt;
            &lt;div style="font-size: 32px; font-weight: bold; color: #67C23A"&gt;{{ stats.winCount || 0 }}&lt;/div&gt;
            &lt;div style="color: #909399; margin-top: 10px"&gt;中奖次数&lt;/div&gt;
          &lt;/div&gt;
        &lt;/el-card&gt;
      &lt;/el-col&gt;
      &lt;el-col :span="6"&gt;
        &lt;el-card&gt;
          &lt;div style="text-align: center"&gt;
            &lt;div style="font-size: 32px; font-weight: bold; color: #E6A23C"&gt;{{ (stats.winRate || 0).toFixed(2) }}%&lt;/div&gt;
            &lt;div style="color: #909399; margin-top: 10px"&gt;中奖率&lt;/div&gt;
          &lt;/div&gt;
        &lt;/el-card&gt;
      &lt;/el-col&gt;
      &lt;el-col :span="6"&gt;
        &lt;el-card&gt;
          &lt;div style="text-align: center"&gt;
            &lt;div style="font-size: 32px; font-weight: bold; color: #F56C6C"&gt;{{ stats.pendingCount || 0 }}&lt;/div&gt;
            &lt;div style="color: #909399; margin-top: 10px"&gt;待领取数量&lt;/div&gt;
          &lt;/div&gt;
        &lt;/el-card&gt;
      &lt;/el-col&gt;
    &lt;/el-row&gt;

    &lt;el-row :gutter="20"&gt;
      &lt;el-col :span="12"&gt;
        &lt;el-card&gt;
          &lt;template #header&gt;
            &lt;span&gt;奖品中奖统计&lt;/span&gt;
          &lt;/template&gt;
          &lt;div ref="pieChart" style="height: 400px"&gt;&lt;/div&gt;
        &lt;/el-card&gt;
      &lt;/el-col&gt;
      &lt;el-col :span="12"&gt;
        &lt;el-card&gt;
          &lt;template #header&gt;
            &lt;span&gt;近7天抽奖趋势&lt;/span&gt;
          &lt;/template&gt;
          &lt;div ref="barChart" style="height: 400px"&gt;&lt;/div&gt;
        &lt;/el-card&gt;
      &lt;/el-col&gt;
    &lt;/el-row&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { ref, onMounted, nextTick } from 'vue'
import axios from 'axios'
import * as echarts from 'echarts'

const stats = ref({})
const pieChart = ref(null)
const barChart = ref(null)

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})

const fetchStats = async () =&gt; {
  const res = await api.get('/stats')
  stats.value = res.data
  nextTick(() =&gt; {
    renderPieChart()
    renderBarChart()
  })
}

const renderPieChart = () =&gt; {
  const chart = echarts.init(pieChart.value)
  const data = (stats.value.prizeStats || []).map(item =&gt; ({
    name: item.prizeName,
    value: item.count
  }))
  chart.setOption({
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', left: 'left' },
    series: [
      {
        type: 'pie',
        radius: '50%',
        data: data,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  })
}

const renderBarChart = () =&gt; {
  const chart = echarts.init(barChart.value)
  const data = stats.value.dailyStats || []
  chart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['抽奖次数', '中奖次数'] },
    xAxis: {
      type: 'category',
      data: data.map(item =&gt; item.date)
    },
    yAxis: { type: 'value' },
    series: [
      {
        name: '抽奖次数',
        type: 'bar',
        data: data.map(item =&gt; item.draws)
      },
      {
        name: '中奖次数',
        type: 'bar',
        data: data.map(item =&gt; item.wins)
      }
    ]
  })
}

onMounted(() =&gt; {
  fetchStats()
})
&lt;/script&gt;
