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
              <div v-if="record.province" class="address-info">
                <span>📦 {{ record.province }}{{ record.city }}{{ record.district }}{{ record.detailAddress }}</span>
              </div>
              <div v-if="record.shippingStatus" class="shipping-status">
                <el-tag size="small" :type="record.shippingStatus === '已发货' ? 'success' : 'info'">
                  {{ record.shippingStatus }}
                </el-tag>
                <span v-if="record.trackingNumber" class="tracking-number">物流：{{ record.trackingNumber }}</span>
              </div>
            </div>
            <el-button 
              v-if="record.isWin && record.status !== '已领取'" 
              type="warning" 
              size="small"
              @click="openAddressForm(record)"
              :loading="submittingId === record.id"
            >
              填写地址
            </el-button>
            <el-button 
              v-if="record.isWin && record.status === '已领取' && !record.province" 
              type="warning" 
              size="small"
              @click="openAddressForm(record)"
              :loading="submittingId === record.id"
            >
              填写地址
            </el-button>
            <el-button 
              v-if="record.isWin && record.status === '已领取' && record.province" 
              type="primary" 
              size="small"
              @click="openAddressForm(record)"
              :loading="submittingId === record.id"
            >
              修改地址
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="addressFormVisible" title="填写收货地址" width="90%">
      <el-form :model="addressForm" label-width="80px">
        <el-form-item label="收货人">
          <el-input v-model="addressForm.receiverName" placeholder="请输入收货人姓名"></el-input>
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="addressForm.receiverPhone" placeholder="请输入手机号" maxlength="11"></el-input>
        </el-form-item>
        <el-form-item label="所在地区">
          <AddressPicker v-model="addressRegion" />
        </el-form-item>
        <el-form-item label="详细地址">
          <el-input v-model="addressForm.detailAddress" type="textarea" :rows="2" placeholder="请输入详细地址" maxlength="100"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addressFormVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitAddress" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyRecords, submitAddress as submitAddressApi } from '../utils/api'
import AddressPicker from '../components/AddressPicker.vue'

const phone = ref('')
const inputPhone = ref('')
const records = ref([])
const submittingId = ref(null)
const addressFormVisible = ref(false)
const submitting = ref(false)
const currentRecord = ref(null)
const addressRegion = ref({})
const addressForm = ref({
  receiverName: '',
  receiverPhone: '',
  detailAddress: ''
})

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

const openAddressForm = (record) => {
  currentRecord.value = record
  addressForm.value = {
    receiverName: record.receiverName || record.name,
    receiverPhone: record.receiverPhone || record.phone,
    detailAddress: record.detailAddress || ''
  }
  addressRegion.value = {
    province: record.province || '',
    city: record.city || '',
    district: record.district || ''
  }
  addressFormVisible.value = true
}

const handleSubmitAddress = async () => {
  if (!addressForm.value.receiverName) {
    ElMessage.warning('请输入收货人姓名')
    return
  }
  if (!addressForm.value.receiverPhone) {
    ElMessage.warning('请输入手机号')
    return
  }
  if (!/^1\d{10}$/.test(addressForm.value.receiverPhone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }
  if (!addressRegion.value.province || !addressRegion.value.city || !addressRegion.value.district) {
    ElMessage.warning('请选择完整的省市区')
    return
  }
  if (!addressForm.value.detailAddress) {
    ElMessage.warning('请输入详细地址')
    return
  }

  submitting.value = true
  submittingId.value = currentRecord.value.id
  try {
    await submitAddressApi(currentRecord.value.id, {
      receiverName: addressForm.value.receiverName,
      receiverPhone: addressForm.value.receiverPhone,
      province: addressRegion.value.province,
      city: addressRegion.value.city,
      district: addressRegion.value.district,
      detailAddress: addressForm.value.detailAddress
    })
    ElMessage.success('地址提交成功，奖品将尽快寄出！')
    addressFormVisible.value = false
    loadRecords()
  } catch (err) {
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
    submittingId.value = null
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
  margin-bottom: 5px;
}

.prize-status.claimed {
  color: #67c23a;
}

.address-info {
  font-size: 12px;
  color: #666;
  margin-top: 5px;
}

.shipping-status {
  margin-top: 5px;
  font-size: 12px;
}

.tracking-number {
  margin-left: 8px;
  color: #666;
}
</style>
