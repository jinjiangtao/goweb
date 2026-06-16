<template>
  <div class="home-page">
    <header class="page-header">
      <div class="container">
        <h1 class="page-title">找好房，就来租房平台</h1>
        <p class="page-subtitle">海量真实房源，为您找到温馨的家</p>
      </div>
    </header>

    <main class="page-main container">
      <div class="search-section">
        <div class="search-box">
          <input
            type="text"
            v-model="searchKeyword"
            placeholder="搜索小区名称、房源标题..."
            class="search-input"
            @keyup.enter="handleSearch"
          />
          <button class="search-btn" @click="handleSearch">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8" />
              <path d="m21 21-4.35-4.35" />
            </svg>
            搜索
          </button>
        </div>

        <div class="filter-section">
          <div class="filter-row">
            <div class="filter-item">
              <label class="filter-label">区域</label>
              <select v-model="filters.district" class="filter-select">
                <option value="">全部</option>
                <option v-for="item in districtOptions" :key="item" :value="item">
                  {{ item }}
                </option>
              </select>
            </div>

            <div class="filter-item">
              <label class="filter-label">户型</label>
              <select v-model="filters.roomType" class="filter-select">
                <option value="">全部</option>
                <option v-for="item in roomTypeOptions" :key="item" :value="item">
                  {{ item }}
                </option>
              </select>
            </div>

            <div class="filter-item">
              <label class="filter-label">租金</label>
              <select v-model="filters.priceRange" class="filter-select">
                <option value="">全部</option>
                <option v-for="item in priceRangeOptions" :key="item.value" :value="item.value">
                  {{ item.label }}
                </option>
              </select>
            </div>

            <div class="filter-item">
              <label class="filter-label">排序</label>
              <select v-model="sortBy" class="filter-select">
                <option value="latest">最新发布</option>
                <option value="price_asc">价格低到高</option>
                <option value="price_desc">价格高到低</option>
              </select>
            </div>
          </div>

          <div class="filter-actions">
            <button class="btn btn-primary" @click="handleSearch">搜索</button>
            <button class="btn btn-default" @click="handleReset">重置</button>
          </div>
        </div>
      </div>

      <div class="result-section">
        <div class="result-header">
          <span class="result-count">
            共找到 <em>{{ total }}</em> 套房源
          </span>
        </div>

        <div v-if="loading" class="house-grid">
          <HouseCardSkeleton v-for="i in pageSize" :key="i" />
        </div>

        <div v-else-if="houseList.length > 0" class="house-grid">
          <HouseCard v-for="house in houseList" :key="house.id" :house="house" />
        </div>

        <EmptyState v-else text="暂无房源，请稍后再来" />

        <Pagination
          v-if="!loading && total > 0"
          v-model:currentPage="currentPage"
          v-model:pageSize="pageSize"
          :total="total"
          @change="handlePageChange"
        />
      </div>
    </main>

    <footer class="page-footer">
      <div class="container">
        <p>© 2024 租房信息平台 版权所有</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { getHouseList } from '../api/house'
import HouseCard from '../components/HouseCard.vue'
import HouseCardSkeleton from '../components/HouseCardSkeleton.vue'
import Pagination from '../components/Pagination.vue'
import EmptyState from '../components/EmptyState.vue'

const searchKeyword = ref('')
const sortBy = ref('latest')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const loading = ref(false)
const houseList = ref([])

const filters = reactive({
  district: '',
  roomType: '',
  priceRange: '',
})

const districtOptions = [
  '东城区',
  '西城区',
  '朝阳区',
  '海淀区',
  '丰台区',
  '石景山区',
  '通州区',
  '顺义区',
  '大兴区',
  '昌平区',
  '房山区',
  '门头沟区',
]

const roomTypeOptions = ['一室一厅', '两室一厅', '两室两厅', '三室一厅', '三室两厅', '四室一厅', '四室两厅', '一室', '两室', '三室']

const priceRangeOptions = [
  { label: '1000-2000元', value: '1000-2000' },
  { label: '2000-4000元', value: '2000-4000' },
  { label: '4000-6000元', value: '4000-6000' },
  { label: '6000元以上', value: '6000+' },
]

async function fetchHouses() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value,
      district: filters.district,
      room_type: filters.roomType,
      sort_by: sortBy.value,
    }

    if (filters.priceRange) {
      if (filters.priceRange.endsWith('+')) {
        params.min_price = filters.priceRange.replace('+', '')
      } else {
        const parts = filters.priceRange.split('-')
        if (parts.length === 2) {
          params.min_price = parts[0]
          params.max_price = parts[1]
        }
      }
    }

    const res = await getHouseList(params)
    houseList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取房源列表失败:', error)
    houseList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  fetchHouses()
}

function handleReset() {
  searchKeyword.value = ''
  filters.district = ''
  filters.roomType = ''
  filters.priceRange = ''
  sortBy.value = 'latest'
  currentPage.value = 1
  fetchHouses()
}

function handlePageChange({ page, pageSize: size }) {
  currentPage.value = page
  pageSize.value = size
  fetchHouses()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  fetchHouses()
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
  padding: 40px 0 60px;
  color: #fff;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
}

.page-subtitle {
  font-size: 15px;
  opacity: 0.9;
}

.page-main {
  flex: 1;
  padding-top: 24px;
  padding-bottom: 40px;
  margin-top: -40px;
}

.search-section {
  background: var(--bg-white);
  border-radius: var(--radius-xl);
  padding: 20px;
  box-shadow: var(--shadow-md);
  margin-bottom: 24px;
}

.search-box {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.search-input {
  flex: 1;
  height: 44px;
  padding: 0 16px;
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  font-size: 15px;
  transition: border-color var(--transition-fast);
}

.search-input:focus {
  border-color: var(--primary-color);
}

.search-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 44px;
  padding: 0 24px;
  background: var(--primary-color);
  color: #fff;
  border-radius: var(--radius-md);
  font-size: 15px;
  font-weight: 500;
  transition: background var(--transition-fast);
}

.search-btn:hover {
  background: var(--primary-dark);
}

.filter-section {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: flex-end;
  justify-content: space-between;
}

.filter-row {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  flex: 1;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 140px;
}

.filter-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.filter-select {
  height: 38px;
  padding: 0 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  font-size: 14px;
  background: var(--bg-white);
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.filter-select:hover,
.filter-select:focus {
  border-color: var(--primary-color);
}

.filter-actions {
  display: flex;
  gap: 10px;
}

.btn {
  height: 38px;
  padding: 0 20px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  transition: all var(--transition-fast);
}

.btn-primary {
  background: var(--primary-color);
  color: #fff;
}

.btn-primary:hover {
  background: var(--primary-dark);
}

.btn-default {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  color: var(--text-color);
}

.btn-default:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.result-section {
  background: var(--bg-white);
  border-radius: var(--radius-xl);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.result-header {
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}

.result-count {
  font-size: 14px;
  color: var(--text-secondary);
}

.result-count em {
  font-style: normal;
  color: var(--primary-color);
  font-weight: 600;
  font-size: 16px;
}

.house-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.page-footer {
  background: var(--secondary-color);
  color: rgba(255, 255, 255, 0.7);
  padding: 24px 0;
  text-align: center;
  font-size: 13px;
}

@media (max-width: 768px) {
  .page-header {
    padding: 24px 0 50px;
  }

  .page-title {
    font-size: 22px;
  }

  .search-section {
    padding: 16px;
  }

  .search-box {
    flex-direction: column;
  }

  .search-btn {
    width: 100%;
    justify-content: center;
  }

  .filter-section {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-row {
    gap: 12px;
  }

  .filter-item {
    min-width: calc(50% - 6px);
    flex: 1;
  }

  .filter-actions {
    width: 100%;
  }

  .filter-actions .btn {
    flex: 1;
  }

  .house-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .result-section {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .house-grid {
    grid-template-columns: 1fr;
  }

  .filter-item {
    min-width: 100%;
  }
}
</style>
