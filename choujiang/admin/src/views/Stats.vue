<template>
  <div>
    <el-row :gutter="20" style="margin-bottom: 20px">
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #409EFF">{{ stats.totalDraws || 0 }}</div>
            <div style="color: #909399; margin-top: 10px">总抽奖次数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #67C23A">{{ stats.winCount || 0 }}</div>
            <div style="color: #909399; margin-top: 10px">中奖次数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #E6A23C">{{ (stats.winRate || 0).toFixed(2) }}%</div>
            <div style="color: #909399; margin-top: 10px">中奖率</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #F56C6C">{{ stats.pendingCount || 0 }}</div>
            <div style="color: #909399; margin-top: 10px">待领取数量</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>奖品中奖统计</span>
          </template>
          <div ref="pieChart" style="height: 400px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>近7天抽奖趋势</span>
          </template>
          <div ref="barChart" style="height: 400px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
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
  const res = await api.get('/admin/stats')
  stats.value = res.data
  nextTick(() =&gt; {
    renderPieChart()
    renderBarChart()
  })
}

const renderPieChart = () => {
  const chart = echarts.init(pieChart.value)
  const data = (stats.value.prizeStats || []).map(item => ({
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

const renderBarChart = () => {
  const chart = echarts.init(barChart.value)
  const data = stats.value.dailyStats || []
  chart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['抽奖次数', '中奖次数'] },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date)
    },
    yAxis: { type: 'value' },
    series: [
      {
        name: '抽奖次数',
        type: 'bar',
        data: data.map(item => item.draws)
      },
      {
        name: '中奖次数',
        type: 'bar',
        data: data.map(item => item.wins)
      }
    ]
  })
}

onMounted(() => {
  fetchStats()
})
</script>