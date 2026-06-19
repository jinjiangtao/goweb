<template>
  <div class="page-container">
    <el-row :gutter="20" class="overview-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card primary" shadow="hover">
          <div class="icon-wrap">
            <el-icon :size="28"><Tickets /></el-icon>
          </div>
          <div class="num">{{ stats.overview?.exam_count || 0 }}</div>
          <div class="label">答题次数</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card success" shadow="hover">
          <div class="icon-wrap">
            <el-icon :size="28"><Edit /></el-icon>
          </div>
          <div class="num">{{ stats.overview?.total_questions || 0 }}</div>
          <div class="label">累计做题</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card warning" shadow="hover">
          <div class="icon-wrap">
            <el-icon :size="28"><CircleCheck /></el-icon>
          </div>
          <div class="num">{{ accuracy }}%</div>
          <div class="label">整体正确率</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card danger" shadow="hover">
          <div class="icon-wrap">
            <el-icon :size="28"><Warning /></el-icon>
          </div>
          <div class="num">{{ stats.overview?.wrong_count || 0 }}</div>
          <div class="label">待攻克错题</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :xs="24" :md="14">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-title">
              <el-icon color="#409EFF"><TrendCharts /></el-icon>
              <span>每日学习趋势</span>
            </div>
          </template>
          <div ref="dailyChartRef" class="chart-box"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="10">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-title">
              <el-icon color="#67C23A"><PieChart /></el-icon>
              <span>题型正确率分析</span>
            </div>
          </template>
          <div ref="typeChartRef" class="chart-box"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :xs="24" :md="16">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-title">
              <el-icon color="#E6A23C"><DataLine /></el-icon>
              <span>近期答题详情</span>
            </div>
          </template>
          <el-table :data="stats.recent_exams || []" stripe v-loading="loading">
            <el-table-column prop="exam_name" label="练习名称" min-width="160" />
            <el-table-column prop="mode" label="模式" width="80" align="center">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ modeName(row.mode) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="题目数" width="80" align="center" prop="total_count" />
            <el-table-column label="正确数" width="80" align="center" prop="correct_count">
              <template #default="{ row }">
                <span style="color: #67C23A">{{ row.correct_count }}</span>
              </template>
            </el-table-column>
            <el-table-column label="得分" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="scoreColor(row.score)" size="small" effect="dark">
                  {{ row.score?.toFixed(1) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="180" align="center" />
            <el-table-column label="操作" width="100" align="center">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="router.push(`/record/${row.id}`)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!loading && (stats.recent_exams || []).length === 0" description="暂无答题记录" />
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-title">
              <el-icon color="#F56C6C"><DataBoard /></el-icon>
              <span>学习画像</span>
            </div>
          </template>
          <div class="profile-box">
            <div class="profile-item">
              <div class="p-label">平均得分</div>
              <div class="p-value big">{{ avgScore }}</div>
            </div>
            <div class="profile-item">
              <div class="p-label">累计答对</div>
              <div class="p-value">{{ stats.overview?.total_correct || 0 }} 题</div>
            </div>
            <div class="profile-item">
              <div class="p-label">学习天数</div>
              <div class="p-value">{{ (stats.daily_stats || []).length }} 天</div>
            </div>
            <div class="profile-item">
              <div class="p-label">错题待攻克</div>
              <div class="p-value" style="color:#F56C6C">{{ stats.overview?.wrong_count || 0 }} 题</div>
            </div>
          </div>
          <el-divider />
          <div class="encourage-box">
            <el-icon :size="28" color="#E6A23C"><Trophy /></el-icon>
            <div class="encourage-text">
              {{ encourageText }}
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { getStatistics } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const stats = ref({})
const dailyChartRef = ref()
const typeChartRef = ref()
let dailyChart = null
let typeChart = null

const accuracy = computed(() => {
  const ov = stats.value.overview || {}
  return ov.overall_accuracy?.toFixed?.(1) || 0
})

const avgScore = computed(() => {
  return (stats.value.overview?.avg_score || 0).toFixed(1) + ' 分'
})

const encourageText = computed(() => {
  const acc = parseFloat(accuracy.value)
  if (acc >= 90) return '太棒了！你已经掌握得非常好，继续保持！ 🏆'
  if (acc >= 75) return '表现不错！再努力一点就能更上一层楼！ 💪'
  if (acc >= 60) return '已经及格啦！多复习错题会有更大进步！ 📚'
  return '加油！多刷题多练习，一定能取得进步！ 🌟'
})

const modeName = (m) => ({ order: '顺序', random: '随机', smart: '智能', wrong: '错题' }[m] || '练习')
const scoreColor = (s) => s >= 80 ? 'success' : s >= 60 ? 'warning' : 'danger'

const loadStats = async () => {
  loading.value = true
  try {
    const res = await getStatistics({ user_id: userStore.userId })
    stats.value = res
    await nextTick()
    renderCharts()
  } finally {
    loading.value = false
  }
}

const renderCharts = () => {
  renderDailyChart()
  renderTypeChart()
}

const renderDailyChart = () => {
  if (!dailyChartRef.value) return
  if (dailyChart) dailyChart.dispose()
  dailyChart = echarts.init(dailyChartRef.value)

  const daily = stats.value.daily_stats || []
  const dates = daily.map(d => d.date.slice(5))
  const questionData = daily.map(d => d.questions)
  const correctData = daily.map(d => d.correct)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' }
    },
    legend: {
      data: ['做题数', '正确数'],
      bottom: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates.length ? dates : ['暂无数据'],
      axisLabel: { color: '#606266' }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#606266' },
      splitLine: { lineStyle: { color: '#ebeef5' } }
    },
    series: [
      {
        name: '做题数',
        type: 'bar',
        data: questionData.length ? questionData : [0],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#667eea' },
            { offset: 1, color: '#764ba2' }
          ]),
          borderRadius: [4, 4, 0, 0]
        },
        barWidth: 16
      },
      {
        name: '正确数',
        type: 'bar',
        data: correctData.length ? correctData : [0],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#43e97b' },
            { offset: 1, color: '#38f9d7' }
          ]),
          borderRadius: [4, 4, 0, 0]
        },
        barWidth: 16
      }
    ]
  }
  dailyChart.setOption(option)
}

const renderTypeChart = () => {
  if (!typeChartRef.value) return
  if (typeChart) typeChart.dispose()
  typeChart = echarts.init(typeChartRef.value)

  const typeStats = stats.value.type_stats || []

  const data = typeStats.map(t => ({
    name: ({ single: '单选题', multiple: '多选题', truefalse: '判断题' })[t.type] || t.type,
    value: t.total,
    accuracy: t.accuracy
  }))

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: (p) => `${p.name}<br/>做题数: ${p.value}<br/>正确率: ${p.data.accuracy?.toFixed?.(1) || 0}%`
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      top: 'center'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        center: ['65%', '50%'],
        itemStyle: {
          borderRadius: 8,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: '{b}\n{d}%'
        },
        labelLine: { show: true },
        data: data.length ? data : [
          { name: '暂无数据', value: 1, accuracy: 0, itemStyle: { color: '#c0c4cc' } }
        ],
        color: ['#409EFF', '#E6A23C', '#67C23A', '#F56C6C']
      }
    ]
  }
  typeChart.setOption(option)
}

const handleResize = () => {
  dailyChart?.resize()
  typeChart?.resize()
}

onMounted(() => {
  loadStats()
  window.addEventListener('resize', handleResize)
})
</script>

<style scoped>
.page-container {
  max-width: 1400px;
  margin: 0 auto;
}

.overview-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  border: none;
}

.stat-card :deep(.el-card__body) {
  padding: 20px;
}

.stat-card.primary { background: linear-gradient(135deg, #667eea10 0%, #764ba210 100%); }
.stat-card.success { background: linear-gradient(135deg, #43e97b10 0%, #38f9d710 100%); }
.stat-card.warning { background: linear-gradient(135deg, #fa709a10 0%, #fee14010 100%); }
.stat-card.danger { background: linear-gradient(135deg, #f093fb10 0%, #f5576c10 100%); }

.icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 10px;
}

.stat-card.primary .icon-wrap { background: #667eea; color: #fff; }
.stat-card.success .icon-wrap { background: #43e97b; color: #fff; }
.stat-card.warning .icon-wrap { background: #fa709a; color: #fff; }
.stat-card.danger .icon-wrap { background: #f093fb; color: #fff; }

.num {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
}

.label {
  font-size: 13px;
  color: #909399;
  margin-top: 6px;
}

.chart-card {
  border-radius: 12px;
  height: 100%;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

.chart-box {
  width: 100%;
  height: 320px;
}

.profile-box {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.profile-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 8px;
}

.p-label {
  color: #909399;
  font-size: 14px;
}

.p-value {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.p-value.big {
  font-size: 20px;
  color: #409EFF;
}

.encourage-box {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 8px;
}

.encourage-text {
  color: #303133;
  line-height: 1.6;
  font-size: 14px;
}
</style>
