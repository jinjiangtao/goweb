<template>
  <div class="address-picker">
    <el-select v-model="provinceId" placeholder="选择省份" @change="handleProvinceChange" style="width: 100%; margin-bottom: 10px;">
      <el-option v-for="province in provinces" :key="province.id" :label="province.name" :value="province.id" />
    </el-select>
    <el-select v-model="cityId" placeholder="选择城市" @change="handleCityChange" style="width: 100%; margin-bottom: 10px;" :disabled="!provinceId">
      <el-option v-for="city in cities" :key="city.id" :label="city.name" :value="city.id" />
    </el-select>
    <el-select v-model="districtId" placeholder="选择区县" @change="handleDistrictChange" style="width: 100%;" :disabled="!cityId">
      <el-option v-for="district in districts" :key="district.id" :label="district.name" :value="district.id" />
    </el-select>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { getProvinces, getCities, getDistricts } from '../utils/api'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue'])

const provinceId = ref(null)
const cityId = ref(null)
const districtId = ref(null)
const provinces = ref([])
const cities = ref([])
const districts = ref([])
const provinceName = ref('')
const cityName = ref('')
const districtName = ref('')

onMounted(() => {
  loadProvinces()
})

watch(() => props.modelValue, (val) => {
  if (val && Object.keys(val).length > 0) {
    initFromValue(val)
  }
}, { immediate: true, deep: true })

const initFromValue = async (val) => {
  if (val.province && provinces.value.length > 0) {
    const province = provinces.value.find(p => p.name === val.province)
    if (province) {
      provinceId.value = province.id
      await handleProvinceChange(province.id)
      if (val.city) {
        const city = cities.value.find(c => c.name === val.city)
        if (city) {
          cityId.value = city.id
          await handleCityChange(city.id)
          if (val.district) {
            const district = districts.value.find(d => d.name === val.district)
            if (district) {
              districtId.value = district.id
            }
          }
        }
      }
    }
  }
}

const loadProvinces = async () => {
  try {
    const res = await getProvinces()
    provinces.value = res.data
  } catch (err) {
    console.error('Failed to load provinces', err)
  }
}

const handleProvinceChange = async (id) => {
  provinceId.value = id
  cityId.value = null
  districtId.value = null
  cities.value = []
  districts.value = []
  
  const province = provinces.value.find(p => p.id === id)
  if (province) {
    provinceName.value = province.name
  }
  
  try {
    const res = await getCities(id)
    cities.value = res.data
  } catch (err) {
    console.error('Failed to load cities', err)
  }
  
  emitValue()
}

const handleCityChange = async (id) => {
  cityId.value = id
  districtId.value = null
  districts.value = []
  
  const city = cities.value.find(c => c.id === id)
  if (city) {
    cityName.value = city.name
  }
  
  try {
    const res = await getDistricts(id)
    districts.value = res.data
  } catch (err) {
    console.error('Failed to load districts', err)
  }
  
  emitValue()
}

const handleDistrictChange = (id) => {
  districtId.value = id
  
  const district = districts.value.find(d => d.id === id)
  if (district) {
    districtName.value = district.name
  }
  
  emitValue()
}

const emitValue = () => {
  emit('update:modelValue', {
    province: provinceName.value,
    city: cityName.value,
    district: districtName.value
  })
}
</script>

<style scoped>
.address-picker {
  width: 100%;
}
</style>
