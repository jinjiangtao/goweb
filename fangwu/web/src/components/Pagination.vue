<template>
  <div class="pagination" v-if="total > 0">
    <div class="pagination-info">
      共 <span class="total">{{ total }}</span> 条
    </div>
    <div class="pagination-controls">
      <button
        class="page-btn prev"
        :disabled="currentPage <= 1"
        @click="changePage(currentPage - 1)"
      >
        上一页
      </button>
      <div class="page-numbers">
        <template v-for="page in visiblePages" :key="page">
          <span class="page-ellipsis" v-if="page === '...'">...</span>
          <button
            v-else
            class="page-num"
            :class="{ active: page === currentPage }"
            @click="changePage(page)"
          >
            {{ page }}
          </button>
        </template>
      </div>
      <button
        class="page-btn next"
        :disabled="currentPage >= totalPages"
        @click="changePage(currentPage + 1)"
      >
        下一页
      </button>
      <select class="page-size" :value="pageSize" @change="changePageSize">
        <option v-for="size in pageSizeOptions" :key="size" :value="size">
          {{ size }}条/页
        </option>
      </select>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  total: {
    type: Number,
    default: 0,
  },
  currentPage: {
    type: Number,
    default: 1,
  },
  pageSize: {
    type: Number,
    default: 10,
  },
  pageSizeOptions: {
    type: Array,
    default: () => [10, 20, 30, 50],
  },
})

const emit = defineEmits(['update:currentPage', 'update:pageSize', 'change'])

const totalPages = computed(() => {
  return Math.ceil(props.total / props.pageSize) || 1
})

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = props.currentPage

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    pages.push(1)

    if (current > 3) {
      pages.push('...')
    }

    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (current < total - 2) {
      pages.push('...')
    }

    pages.push(total)
  }

  return pages
})

function changePage(page) {
  if (page < 1 || page > totalPages.value || page === props.currentPage) return
  emit('update:currentPage', page)
  emit('change', { page, pageSize: props.pageSize })
}

function changePageSize(e) {
  const size = parseInt(e.target.value)
  emit('update:pageSize', size)
  emit('change', { page: 1, pageSize: size })
}
</script>

<style scoped>
.pagination {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 20px 0;
}

.pagination-info {
  font-size: 14px;
  color: var(--text-secondary);
}

.pagination-info .total {
  color: var(--primary-color);
  font-weight: 600;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.page-btn {
  padding: 8px 16px;
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  font-size: 14px;
  color: var(--text-color);
  transition: all var(--transition-fast);
}

.page-btn:hover:not(:disabled) {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-num {
  min-width: 36px;
  height: 36px;
  padding: 0 12px;
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  font-size: 14px;
  color: var(--text-color);
  transition: all var(--transition-fast);
}

.page-num:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.page-num.active {
  background: var(--primary-color);
  border-color: var(--primary-color);
  color: #fff;
}

.page-ellipsis {
  padding: 0 8px;
  color: var(--text-light);
}

.page-size {
  padding: 8px 12px;
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  font-size: 14px;
  color: var(--text-color);
  cursor: pointer;
}

@media (max-width: 480px) {
  .pagination {
    flex-direction: column;
    align-items: stretch;
  }

  .pagination-info {
    text-align: center;
  }

  .pagination-controls {
    justify-content: center;
  }

  .page-btn {
    padding: 6px 12px;
    font-size: 13px;
  }

  .page-num {
    min-width: 32px;
    height: 32px;
    padding: 0 8px;
    font-size: 13px;
  }
}
</style>
