<template>
  <div class="my-prizes">
    <div class="header">
      <el-button @click="$router.back()" icon="ArrowLeft" circle></el-button>
      <h2>🎁 我的奖品</h2>
      <div style="width: 40px;"></div>
    </div>

    <div class="content">
      <div v-if="!phone" class="phone-form">
        <el-input v-model="inputPhone" placeholder="请输入手机号查询" maxlength="11"></el-input>
        <el-button type="primary" @click="loadRecords" style="margin-top: 15px; width: 100%;">查询</el-button>
      </div>

      <div v-else>
        <div v-if="records.length === 0" class="empty">
          <div class="empty-icon">🎊</div>
          <p>暂无抽奖记录</p>
          <el-button type="primary" @click="$router.push('/')" style="margin-top: 20px;">去抽奖</el-button>
        </div>

        <div v-else class="records-list">
          <div v-for="record in records" :key="record.id" class="record-item">
            <div class="record-info">
              <div class="prize-name">{{ record.prizeName }}</div>
              <div class="prize-time">{{ formatTime(record.createdAt) }}</div>
              <div class="prize-status" :class="{ claimed: record.status === '已领取' }">{{ record.status }}</div>
            </div>
            <el-button 
              v-if="record.isWin && record.status !== '已领取'" 
              type="warning" 
              size="small"
              @click="claimPrize(record)"
              :loading="claimingId === record.id"
            >
              去领取
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyRecords, claimRecord } from '../utils/api'

const phone = ref('')
const inputPhone = ref('')
const records = ref([])
const claimingId = ref(null)

onMounted(() => {
  const saved = localStorage.getItem('lottery_user')
  if (saved) {
    const user = JSON.parse(saved)
    phone.value = user.phone
    inputPhone.value = user.phone
    loadRecords()
  }
})

const loadRecords = async () => {
  if (!inputPhone.value) {
    ElMessage.warning('请输入手机号')
    return
  }
  if (!/^1\d{10}$/.test(inputPhone.value)) {
    ElMessage.warning('请输入正确的11位手机号')
    return
  }
  try {
    const res = await getMyRecords(inputPhone.value)
    records.value = res.data
    phone.value = inputPhone.value
  } catch (err) {
    ElMessage.error('查询失败')
  }
}

const claimPrize = async (record) => {
  claimingId.value = record.id
  try {
    await claimRecord(record.id)
    ElMessage.success('领取成功')
    record.status = '已领取'
  } catch (err) {
    ElMessage.error('领取失败')
  } finally {
    claimingId.value = null
  }
}

const formatTime = (time) => {
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.my-prizes {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 20px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
}

.header h2 {
  color: #fff;
  font-size: 20px;
  margin: 0;
}

.content {
  flex: 1;
  padding: 20px;
}

.phone-form {
  background: #fff;
  padding: 25px;
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.empty {
  text-align: center;
  padding: 60px 20px;
  background: #fff;
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 80px;
  margin-bottom: 20px;
}

.empty p {
  color: #666;
  font-size: 16px;
  margin: 0;
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.record-item {
  background: #fff;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.record-info {
  flex: 1;
}

.prize-name {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
}

.prize-time {
  font-size: 13px;
  color: #999;
  margin-bottom: 5px;
}

.prize-status {
  font-size: 14px;
  color: #ff6b6b;
  font-weight: 500;
}

.prize-status.claimed {
  color: #67c23a;
}
</style>
