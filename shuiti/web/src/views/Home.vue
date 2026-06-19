<template>
  <div class="home-page">
    <div class="hero-section">
      <div class="hero-content">
        <div class="hero-text">
          <h1>在线刷题学习系统</h1>
          <p class="subtitle">海量题库 · 智能组卷 · 实时判分 · 错题归集</p>
          <p class="desc">助力高效学习，查漏补缺，科学提升学习效果</p>
          <div class="hero-actions">
            <el-button type="primary" size="large" @click="router.push('/practice')">
              <el-icon><EditPen /></el-icon>
              立即开始刷题
            </el-button>
            <el-button size="large" @click="router.push('/stats')">
              <el-icon><TrendCharts /></el-icon>
              查看学习统计
            </el-button>
          </div>
        </div>
        <div class="hero-illustration">
          <div class="illustration-box">
            <el-icon :size="120" color="#409EFF"><Reading /></el-icon>
            <div class="float-tag tag-1">
              <el-icon><Star /></el-icon>
              <span>智能组卷</span>
            </div>
            <div class="float-tag tag-2">
              <el-icon><CircleCheck /></el-icon>
              <span>自动判分</span>
            </div>
            <div class="float-tag tag-3">
              <el-icon><Trophy /></el-icon>
              <span>数据统计</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-row :gutter="20" class="stat-cards">
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon icon-1">
              <el-icon :size="28"><Collection /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-num">{{ dashboard.total_questions || 0 }}</div>
              <div class="stat-label">题目总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon icon-2">
              <el-icon :size="28"><Tickets /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-num">{{ dashboard.total_exams || 0 }}</div>
              <div class="stat-label">答题次数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon icon-3">
              <el-icon :size="28"><Sunny /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-num">{{ dashboard.today_exams || 0 }}</div>
              <div class="stat-label">今日练习</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon icon-4">
              <el-icon :size="28"><Medal /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-num">{{ typeCount }}</div>
              <div class="stat-label">题目类型</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="feature-section">
      <el-col :span="16">
        <el-card class="feature-card" shadow="never">
          <template #header>
            <div class="card-title">
              <el-icon color="#67C23A"><List /></el-icon>
              <span>题目类型分布</span>
            </div>
          </template>
          <div class="type-list">
            <div v-for="t in dashboard.type_count || []" :key="t.type" class="type-item">
              <div class="type-head">
                <el-tag :type="typeColor(t.type)" size="large" effect="dark">{{ typeName(t.type) }}</el-tag>
                <span class="type-count">{{ t.count }} 题</span>
              </div>
              <el-progress
                :percentage="Math.round(t.count / (dashboard.total_questions || 1) * 100)"
                :color="progressColor(t.type)"
                :stroke-width="10"
              />
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="feature-card" shadow="never">
          <template #header>
            <div class="card-title">
              <el-icon color="#409EFF"><Lightning /></el-icon>
              <span>快速开始</span>
            </div>
          </template>
          <div class="quick-actions">
            <div class="quick-item" @click="startExam('order')">
              <div class="q-icon primary">
                <el-icon :size="24"><Sort /></el-icon>
              </div>
              <div class="q-info">
                <div class="q-title">顺序刷题</div>
                <div class="q-desc">按题库顺序练习</div>
              </div>
            </div>
            <div class="quick-item" @click="startExam('random')">
              <div class="q-icon warning">
                <el-icon :size="24"><Shuffle /></el-icon>
              </div>
              <div class="q-info">
                <div class="q-title">随机组卷</div>
                <div class="q-desc">随机抽取题目</div>
              </div>
            </div>
            <div class="quick-item" @click="startExam('smart')">
              <div class="q-icon success">
                <el-icon :size="24"><MagicStick /></el-icon>
              </div>
              <div class="q-info">
                <div class="q-title">智能组卷</div>
                <div class="q-desc">各题型均衡搭配</div>
              </div>
            </div>
            <div class="quick-item" @click="router.push('/wrong')">
              <div class="q-icon danger">
                <el-icon :size="24"><Warning /></el-icon>
              </div>
              <div class="q-info">
                <div class="q-title">错题回顾</div>
                <div class="q-desc">重点攻克薄弱点</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="feature-card recent-card" shadow="never">
      <template #header>
        <div class="card-title">
          <el-icon color="#E6A23C"><Clock /></el-icon>
          <span>最近答题记录</span>
          <el-button link type="primary" @click="router.push('/records')" style="margin-left:auto">查看全部</el-button>
        </div>
      </template>
      <el-table :data="dashboard.recent_exams || []" v-if="(dashboard.recent_exams || []).length" stripe>
        <el-table-column prop="exam_name" label="练习名称" />
        <el-table-column prop="mode" label="模式" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ modeName(row.mode) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="正确率" width="140" align="center">
          <template #default="{ row }">
            {{ row.correct_count }} / {{ row.total_count }}
            <el-tag :type="scoreColor(row.score)" size="small" style="margin-left:8px">{{ row.score?.toFixed(1) }}分</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="200" align="center" />
        <el-table-column label="操作" width="120" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="router.push(`/record/${row.id}`)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="暂无答题记录" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboard } from '@/api'

const router = useRouter()
const dashboard = ref({})

const loadDashboard = async () => {
  const res = await getDashboard()
  dashboard.value = res
}

const typeCount = computed(() => (dashboard.value.type_count || []).length)

const typeName = (t) => ({ single: '单选题', multiple: '多选题', truefalse: '判断题' }[t] || t)
const typeColor = (t) => ({ single: 'primary', multiple: 'warning', truefalse: 'success' }[t] || 'info')
const progressColor = (t) => ({ single: '#409EFF', multiple: '#E6A23C', truefalse: '#67C23A' }[t] || '#909399')
const modeName = (m) => ({ order: '顺序', random: '随机', smart: '智能', wrong: '错题' }[m] || '练习')
const scoreColor = (s) => s >= 80 ? 'success' : s >= 60 ? 'warning' : 'danger'

const startExam = (mode) => {
  router.push(`/exam/${mode}`)
}

onMounted(loadDashboard)
</script>

<style scoped>
.home-page {
  max-width: 1400px;
  margin: 0 auto;
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 48px;
  margin-bottom: 24px;
  color: #fff;
  overflow: hidden;
  position: relative;
}

.hero-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 40px;
}

.hero-text h1 {
  font-size: 38px;
  font-weight: 700;
  margin-bottom: 16px;
  color: #fff;
}

.hero-text .subtitle {
  font-size: 20px;
  opacity: 0.95;
  margin-bottom: 12px;
}

.hero-text .desc {
  font-size: 15px;
  opacity: 0.85;
  margin-bottom: 32px;
}

.hero-actions {
  display: flex;
  gap: 16px;
}

.hero-illustration {
  flex-shrink: 0;
}

.illustration-box {
  position: relative;
  width: 280px;
  height: 240px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
}

.float-tag {
  position: absolute;
  background: #fff;
  color: #303133;
  padding: 8px 14px;
  border-radius: 20px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.tag-1 {
  top: -12px;
  right: -10px;
  color: #E6A23C;
}

.tag-2 {
  bottom: 20px;
  left: -20px;
  color: #67C23A;
}

.tag-3 {
  top: 60px;
  left: -30px;
  color: #F56C6C;
}

.stat-cards {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 12px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.icon-1 { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.icon-2 { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.icon-3 { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.icon-4 { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }

.stat-num {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.feature-section {
  margin-bottom: 24px;
}

.feature-card {
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

.type-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.type-item {
  padding: 0 4px;
}

.type-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.type-count {
  color: #606266;
  font-size: 14px;
  font-weight: 500;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.quick-item {
  display: flex;
  align-items: center;
  padding: 14px 16px;
  border-radius: 10px;
  background: #f8fafc;
  cursor: pointer;
  transition: all 0.2s ease;
  gap: 14px;
}

.quick-item:hover {
  background: #f0f5ff;
  transform: translateX(4px);
}

.q-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.q-icon.primary { background: #409EFF; }
.q-icon.warning { background: #E6A23C; }
.q-icon.success { background: #67C23A; }
.q-icon.danger { background: #F56C6C; }

.q-title {
  font-weight: 600;
  color: #303133;
  font-size: 15px;
}

.q-desc {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.recent-card {
  border-radius: 12px;
}
</style>
