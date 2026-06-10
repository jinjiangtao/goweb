<template>
  <div>
    <div style="margin-bottom: 20px; display: flex; gap: 10px">
      <el-input v-model="query.name" placeholder="姓名" style="width: 200px" clearable @clear="fetchRecords" />
      <el-input v-model="query.phone" placeholder="手机号" style="width: 200px" clearable @clear="fetchRecords" />
      <el-select v-model="query.prizeName" placeholder="奖品" style="width: 200px" clearable @clear="fetchRecords">
        <el-option v-for="prize in prizeNames" :key="prize" :label="prize" :value="prize" />
      </el-select>
      <el-select v-model="query.isWin" placeholder="是否中奖" style="width: 150px" clearable @clear="fetchRecords">
        <el-option label="是" value="true" />
        <el-option label="否" value="false" />
      </el-select>
      <el-button type="primary" @click="fetchRecords">搜索</el-button>
      <el-button type="success" @click="exportRecords">导出Excel</el-button>
    </div>

    <el-table :data="records" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column prop="phone" label="手机号" width="150" />
      <el-table-column prop="prizeName" label="奖品" />
      <el-table-column label="是否中奖" width="100">
        <template #default="{ row }">
          <el-tag :type="row.isWin ? 'success' : 'info'">{{ row.isWin ? '是' : '否' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="row.status === '已领取' ? 'success' : 'warning'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="收货地址" width="200">
        <template #default="{ row }">
          <span v-if="row.province">{{ row.province }}{{ row.city }}{{ row.district }}{{ row.detailAddress }}</span>
          <span v-else class="text-gray">未填写</span>
        </template>
      </el-table-column>
      <el-table-column label="发货状态" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.shippingStatus" :type="row.shippingStatus === '已发货' ? 'success' : 'info'">
            {{ row.shippingStatus }}
          </el-tag>
          <span v-else class="text-gray">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.isWin && row.status === '待领取'" type="primary" size="small" @click="claimRecord(row)">标记领取</el-button>
          <el-button v-if="row.isWin" type="warning" size="small" @click="openAddressDialog(row)">
            {{ row.province ? '编辑地址' : '填写地址' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="addressDialogVisible" title="编辑收货地址" width="600px">
      <el-form :model="addressForm" label-width="100px">
        <el-form-item label="收货人">
          <el-input v-model="addressForm.receiverName" placeholder="请输入收货人姓名"></el-input>
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="addressForm.receiverPhone" placeholder="请输入手机号" maxlength="11"></el-input>
        </el-form-item>
        <el-form-item label="省份">
          <el-select v-model="addressForm.province" placeholder="请选择省份" @change="loadCities" style="width: 100%">
            <el-option v-for="province in provinces" :key="province.name" :label="province.name" :value="province.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="城市">
          <el-select v-model="addressForm.city" placeholder="请选择城市" @change="loadDistricts" style="width: 100%" :disabled="!addressForm.province">
            <el-option v-for="city in cities" :key="city.name" :label="city.name" :value="city.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="区县">
          <el-select v-model="addressForm.district" placeholder="请选择区县" style="width: 100%" :disabled="!addressForm.city">
            <el-option v-for="district in districts" :key="district.name" :label="district.name" :value="district.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="详细地址">
          <el-input v-model="addressForm.detailAddress" type="textarea" :rows="2" placeholder="请输入详细地址"></el-input>
        </el-form-item>
        <el-form-item label="发货状态">
          <el-select v-model="addressForm.shippingStatus" placeholder="请选择发货状态" style="width: 100%">
            <el-option label="未发货" value="未发货" />
            <el-option label="已发货" value="已发货" />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号">
          <el-input v-model="addressForm.trackingNumber" placeholder="请输入物流单号"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addressDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveAddress" :loading="addressSubmitting">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const records = ref([])
const query = ref({
  name: '',
  phone: '',
  prizeName: '',
  isWin: ''
})
const addressDialogVisible = ref(false)
const addressSubmitting = ref(false)
const currentRecord = ref(null)
const provinces = ref([])
const cities = ref([])
const districts = ref([])
const addressForm = ref({
  receiverName: '',
  receiverPhone: '',
  province: '',
  city: '',
  district: '',
  detailAddress: '',
  shippingStatus: '',
  trackingNumber: ''
})

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})

const prizeNames = computed(() => {
  const names = new Set()
  records.value.forEach(r => {
    if (r.prizeName && r.prizeName !== '未中奖') {
      names.add(r.prizeName)
    }
  })
  return Array.from(names)
})

const fetchRecords = async () => {
  const params = {}
  if (query.value.name) params.name = query.value.name
  if (query.value.phone) params.phone = query.value.phone
  if (query.value.prizeName) params.prizeName = query.value.prizeName
  if (query.value.isWin) params.isWin = query.value.isWin
  const res = await api.get('/admin/records', { params })
  records.value = res.data
}

const claimRecord = async (row) => {
  await api.put(`/admin/records/${row.id}/claim`)
  ElMessage.success('已标记领取')
  fetchRecords()
}

const exportRecords = () => {
  window.open(`http://localhost:8080/api/admin/records/export?token=${localStorage.getItem('token')}`)
}

const loadProvinces = async () => {
  try {
    const res = await api.get('/admin/address/provinces')
    provinces.value = res.data
  } catch (err) {
    console.error('加载省份失败', err)
  }
}

const loadCities = async () => {
  addressForm.value.city = ''
  addressForm.value.district = ''
  cities.value = []
  districts.value = []
  const province = provinces.value.find(p => p.name === addressForm.value.province)
  if (province) {
    try {
      const res = await api.get(`/admin/address/cities/${province.id}`)
      cities.value = res.data
    } catch (err) {
      console.error('加载城市失败', err)
    }
  }
}

const loadDistricts = async () => {
  addressForm.value.district = ''
  districts.value = []
  const city = cities.value.find(c => c.name === addressForm.value.city)
  if (city) {
    try {
      const res = await api.get(`/admin/address/districts/${city.id}`)
      districts.value = res.data
    } catch (err) {
      console.error('加载区县失败', err)
    }
  }
}

const openAddressDialog = async (row) => {
  currentRecord.value = row
  addressForm.value = {
    receiverName: row.receiverName || row.name,
    receiverPhone: row.receiverPhone || row.phone,
    province: row.province || '',
    city: row.city || '',
    district: row.district || '',
    detailAddress: row.detailAddress || '',
    shippingStatus: row.shippingStatus || '',
    trackingNumber: row.trackingNumber || ''
  }
  
  // 先加载省份
  await loadProvinces()
  
  // 如果已有地址，加载对应的城市和区县
  if (addressForm.value.province) {
    const province = provinces.value.find(p => p.name === addressForm.value.province)
    if (province) {
      try {
        const res = await api.get(`/admin/address/cities/${province.id}`)
        cities.value = res.data
      } catch (err) {
        console.error('加载城市失败', err)
      }
      
      if (addressForm.value.city) {
        const city = cities.value.find(c => c.name === addressForm.value.city)
        if (city) {
          try {
            const res = await api.get(`/admin/address/districts/${city.id}`)
            districts.value = res.data
          } catch (err) {
            console.error('加载区县失败', err)
          }
        }
      }
    }
  }
  
  addressDialogVisible.value = true
}

const saveAddress = async () => {
  if (!addressForm.value.receiverName) {
    ElMessage.warning('请输入收货人姓名')
    return
  }
  addressSubmitting.value = true
  try {
    await api.put(`/admin/records/${currentRecord.value.id}/address`, addressForm.value)
    ElMessage.success('保存成功')
    addressDialogVisible.value = false
    fetchRecords()
  } catch (err) {
    ElMessage.error('保存失败')
  } finally {
    addressSubmitting.value = false
  }
}

onMounted(() => {
  fetchRecords()
})
</script>

<style scoped>
.text-gray {
  color: #999;
}
</style>
