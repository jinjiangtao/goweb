
&lt;template&gt;
  &lt;div class="my-prizes"&gt;
    &lt;div class="header"&gt;
      &lt;el-button @click="$router.back()" icon="ArrowLeft" circle /&gt;
      &lt;h2&gt;🎁 我的奖品&lt;/h2&gt;
      &lt;div style="width: 40px;"&gt;&lt;/div&gt;
    &lt;/div&gt;

    &lt;div class="content"&gt;
      &lt;div v-if="!phone" class="phone-form"&gt;
        &lt;el-input v-model="inputPhone" placeholder="请输入手机号查询" maxlength="11" /&gt;
        &lt;el-button type="primary" @click="loadRecords" style="margin-top: 15px; width: 100%;"&gt;查询&lt;/el-button&gt;
      &lt;/div&gt;

      &lt;div v-else&gt;
        &lt;div v-if="records.length === 0" class="empty"&gt;
          &lt;div class="empty-icon"&gt;🎊&lt;/div&gt;
          &lt;p&gt;暂无抽奖记录&lt;/p&gt;
          &lt;el-button type="primary" @click="$router.push('/')" style="margin-top: 20px;"&gt;去抽奖&lt;/el-button&gt;
        &lt;/div&gt;

        &lt;div v-else class="records-list"&gt;
          &lt;div v-for="record in records" :key="record.id" class="record-item"&gt;
            &lt;div class="record-info"&gt;
              &lt;div class="prize-name"&gt;{{ record.prizeName }}&lt;/div&gt;
              &lt;div class="prize-time"&gt;{{ formatTime(record.createdAt) }}&lt;/div&gt;
              &lt;div class="prize-status" :class="{ claimed: record.status === '已领取' }"&gt;{{ record.status }}&lt;/div&gt;
            &lt;/div&gt;
            &lt;el-button 
              v-if="record.isWin &amp;&amp; record.status !== '已领取'" 
              type="warning" 
              size="small"
              @click="claimPrize(record)"
              :loading="claimingId === record.id"
            &gt;
              去领取
            &lt;/el-button&gt;
          &lt;/div&gt;
        &lt;/div&gt;
      &lt;/div&gt;
    &lt;/div&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyRecords, claimRecord } from '../utils/api'

const phone = ref('')
const inputPhone = ref('')
const records = ref([])
const claimingId = ref(null)

onMounted(() =&gt; {
  const saved = localStorage.getItem('lottery_user')
  if (saved) {
    const user = JSON.parse(saved)
    phone.value = user.phone
    inputPhone.value = user.phone
    loadRecords()
  }
})

const loadRecords = async () =&gt; {
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

const claimPrize = async (record) =&gt; {
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

const formatTime = (time) =&gt; {
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}
&lt;/script&gt;

&lt;style scoped&gt;
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
&lt;/style&gt;
