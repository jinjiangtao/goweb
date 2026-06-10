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
const isInitializing = ref(false) // 新增：用于防止初始化时的循环

// 先定义所有函数
const emitValue = () => {
  if (isInitializing.value) {
    return // 初始化过程中不 emit
  }
  emit('update:modelValue', {
    province: provinceName.value,
    city: cityName.value,
    district: districtName.value
  })
}

const handleDistrictChange = (id) => {
  districtId.value = id
  
  const district = districts.value.find(d => d.id === id)
  if (district) {
    districtName.value = district.name
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

const initFromValue = async (val) => {
  if (!val.province || provinces.value.length === 0) {
    return
  }

  isInitializing.value = true // 开始初始化，设置标志

  try {
    const province = provinces.value.find(p => p.name === val.province)
    if (province) {
      provinceId.value = province.id
      provinceName.value = province.name
      
      try {
        const cityRes = await getCities(province.id)
        cities.value = cityRes.data
        
        if (val.city) {
          const city = cities.value.find(c => c.name === val.city)
          if (city) {
            cityId.value = city.id
            cityName.value = city.name
            
            try {
              const districtRes = await getDistricts(city.id)
              districts.value = districtRes.data
              
              if (val.district) {
                const district = districts.value.find(d => d.name === val.district)
                if (district) {
                  districtId.value = district.id
                  districtName.value = district.name
                }
              }
            } catch (err) {
              console.error('Failed to load districts during init', err)
            }
          }
        }
      } catch (err) {
        console.error('Failed to load cities during init', err)
      }
    }
  } finally {
    isInitializing.value = false // 初始化完成，清除标志
  }
}

const loadProvinces = async () => {
  try {
    const res = await getProvinces()
    provinces.value = res.data
    
    // 省份加载完成后，检查是否需要初始化选中值
    if (props.modelValue && Object.keys(props.modelValue).length > 0) {
      await initFromValue(props.modelValue)
    }
  } catch (err) {
    console.error('Failed to load provinces', err)
  }
}

onMounted(() => {
  loadProvinces()
})

// 监听 modelValue 变化（但不 immediate）
watch(() => props.modelValue, async (newVal, oldVal) => {
  // 只在真正有变化时才处理，且避开初始化阶段
  if (isInitializing.value) {
    return
  }
  
  if (newVal && Object.keys(newVal).length > 0 && provinces.value.length > 0) {
    // 检查是否真的有变化，避免相同值导致的循环
    const hasProvinceChange = newVal.province !== oldVal?.province
    const hasCityChange = newVal.city !== oldVal?.city
    const hasDistrictChange = newVal.district !== oldVal?.district
    
    if (hasProvinceChange || hasCityChange || hasDistrictChange) {
      await initFromValue(newVal)
    }
  }
}, { deep: true })
</script>

<style scoped>
.address-picker {
  width: 100%;
}
</style>
