<template>
  <div class="stats-container">
    <div class="stats-cards">
      <div class="stat-card pending">
        <div class="stat-icon">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.pending }}</div>
          <div class="stat-label">报名中</div>
        </div>
      </div>
      <div class="stat-card approved">
        <div class="stat-icon">
          <el-icon><CircleCheck /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.approved }}</div>
          <div class="stat-label">报名成功</div>
        </div>
      </div>
      <div class="stat-card rejected">
        <div class="stat-icon">
          <el-icon><CircleClose /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.rejected }}</div>
          <div class="stat-label">报名失败</div>
        </div>
      </div>
    </div>
    <div class="chart-section">
      <h3>每日报名人数</h3>
      <div class="bar-chart">
        <div
          v-for="(item, index) in dailyData"
          :key="index"
          class="bar-item"
        >
          <div class="bar" :style="{ height: item.count * 10 + 'px' }"></div>
          <div class="bar-label">{{ item.date }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>import { reactive, onMounted } from 'vue';
import axios from 'axios';
import { Clock, CircleCheck, CircleClose } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
const stats = reactive({
 pending: 0,
 approved: 0,
 rejected: 0
});
const dailyData = reactive([
 { date: '06-01', count: 5 },
 { date: '06-02', count: 8 },
 { date: '06-03', count: 12 },
 { date: '06-04', count: 6 },
 { date: '06-05', count: 10 },
 { date: '06-06', count: 7 },
 { date: '06-07', count: 15 }
]);
const loadStats = async () => {
 try {
 const response = await axios.get('/api/admin/stats');
 stats.pending = response.data.pending;
 stats.approved = response.data.approved;
 stats.rejected = response.data.rejected;
 }
 catch (error) {
 ElMessage.error('获取统计数据失败');
 }
};
onMounted(() => {
 loadStats();
});
</script>

<style scoped>
.stats-container {
  padding: 20px;
}
.stats-cards {
  display: flex;
  gap: 20px;
  margin-bottom: 30px;
}
.stat-card {
  flex: 1;
  display: flex;
  align-items: center;
  padding: 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}
.stat-card.pending {
  border-left: 4px solid #e6a23c;
}
.stat-card.approved {
  border-left: 4px solid #67c23a;
}
.stat-card.rejected {
  border-left: 4px solid #f56c6c;
}
.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
}
.pending .stat-icon {
  background: rgba(230, 162, 60, 0.1);
  color: #e6a23c;
}
.approved .stat-icon {
  background: rgba(103, 194, 58, 0.1);
  color: #67c23a;
}
.rejected .stat-icon {
  background: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
}
.stat-value {
  font-size: 32px;
  font-weight: 600;
  color: #333;
}
.stat-label {
  color: #999;
  font-size: 14px;
}
.chart-section {
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}
.chart-section h3 {
  margin: 0 0 20px 0;
  color: #333;
}
.bar-chart {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  height: 200px;
  padding: 20px 0;
  border-bottom: 2px solid #f0f0f0;
}
.bar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}
.bar {
  width: 40px;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.3s ease;
  min-height: 10px;
}
.bar-label {
  font-size: 12px;
  color: #999;
}
</style>
