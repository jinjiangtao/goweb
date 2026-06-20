<template>
  <div class="page-container">
    <div class="card">
      <div class="card-header">
        <div class="card-title">
          <el-icon style="vertical-align: middle; margin-right: 6px; color: #409eff"><Calendar /></el-icon>
          工时填报
        </div>
        <div class="flex-row" style="gap: 12px">
          <el-radio-group v-model="viewMode" size="default">
            <el-radio-button value="month">月视图</el-radio-button>
            <el-radio-button value="week">周视图</el-radio-button>
          </el-radio-group>
          <div class="flex-row" style="gap: 4px">
            <el-button size="default" @click="prevPeriod">
              <el-icon><ArrowLeft /></el-icon>
            </el-button>
            <el-button size="default" @click="goToday">今天</el-button>
            <el-button size="default" @click="nextPeriod">
              <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
          <span style="font-size: 16px; font-weight: 600; color: #303133; min-width: 140px; text-align: center">{{ periodText }}</span>
        </div>
      </div>

      <div class="flex-row" style="gap: 20px; margin-bottom: 20px; padding: 12px 16px; background: #ecf5ff; border-radius: 6px">
        <div>
          <span style="color: #606266">本月总工时：</span>
          <span style="font-size: 20px; font-weight: 700; color: #409eff">{{ monthStats.total }}</span>
          <span style="color: #606266"> 小时</span>
        </div>
        <div>
          <span style="color: #606266">已审核：</span>
          <span style="font-size: 16px; font-weight: 600; color: #67c23a">{{ monthStats.approved }}</span>
          <span style="color: #606266"> h</span>
        </div>
        <div>
          <span style="color: #606266">待审核：</span>
          <span style="font-size: 16px; font-weight: 600; color: #e6a23c">{{ monthStats.pending }}</span>
          <span style="color: #606266"> h</span>
        </div>
        <div>
          <span style="color: #606266">已驳回：</span>
          <span style="font-size: 16px; font-weight: 600; color: #f56c6c">{{ monthStats.rejected }}</span>
          <span style="color: #606266"> h</span>
        </div>
      </div>

      <div class="calendar-weekdays">
        <div v-for="d in weekdays" :key="d" class="calendar-weekday">{{ d }}</div>
      </div>
      <div class="calendar-grid">
        <div
          v-for="day in calendarDays"
          :key="day.key"
          :class="['calendar-day', {
            'other-month': !day.currentMonth,
            'today': day.isToday,
            'selected': selectedDate === day.dateStr
          }]"
          @click="selectDate(day)"
        >
          <div class="calendar-day-header">
            <span class="calendar-day-date">{{ day.dayNum }}</span>
            <span v-if="day.hours > 0" class="calendar-day-hours">{{ day.hours }}h</span>
          </div>
          <div class="calendar-day-items">
            <div v-for="item in day.items.slice(0, 2)" :key="item.id" class="calendar-day-item" :title="item.project.name + ': ' + item.content">
              <span :class="'status-' + item.status" style="font-weight: 600">●</span>
              {{ item.project.name }} {{ item.hours }}h
            </div>
            <div v-if="day.items.length > 2" style="font-size: 11px; color: #909399; margin-top: 2px">
              +{{ day.items.length - 2 }} 更多...
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="card" v-if="selectedDate">
      <div class="card-header">
        <div class="card-title">
          <el-icon style="vertical-align: middle; margin-right: 6px; color: #409eff"><Edit /></el-icon>
          {{ selectedDate }} 工时详情
          <el-tag v-if="selectedDateData.totalHours > 0" style="margin-left: 10px">{{ selectedDateData.totalHours }} 小时</el-tag>
          <el-tag v-if="selectedDateData.totalHours > 8" type="warning" style="margin-left: 4px">超时</el-tag>
        </div>
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加工时
        </el-button>
      </div>

      <el-table :data="selectedDateData.records" border style="width: 100%">
        <el-table-column prop="project.name" label="项目名称" min-width="140">
          <template #default="{ row }">
            <span><el-tag size="small">{{ row.project?.code }}</el-tag> {{ row.project?.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="hours" label="工时(小时)" width="100" align="center">
          <template #default="{ row }">
            <span style="font-weight: 600; color: #409eff">{{ row.hours }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="工作内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'pending'" type="warning" size="small">待审核</el-tag>
            <el-tag v-else-if="row.status === 'approved'" type="success" size="small">已通过</el-tag>
            <el-tag v-else-if="row.status === 'rejected'" type="danger" size="small">已驳回</el-tag>
          </template>
        </el-table-column>
        <el-table-column v-if="selectedDateData.rejectedCount > 0" prop="rejectReason" label="驳回原因" min-width="140" show-overflow-tooltip />
        <el-table-column label="操作" width="160" align="center" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="openEditDialog(row)" :disabled="row.status === 'approved'">编辑</el-button>
            <el-button size="small" type="danger" link @click="handleDelete(row)" :disabled="row.status === 'approved'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="selectedDateData.records.length === 0" description="暂无工时记录，点击上方按钮添加" />
    </div>

    <el-dialog v-model="dialogVisible" :title="editingRecord ? '编辑工时' : '添加工时'" width="500px" destroy-on-close>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="90px">
        <el-form-item label="工作日期">
          <el-date-picker v-model="form.work_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="所属项目" prop="project_id">
          <el-select v-model="form.project_id" placeholder="请选择项目" style="width: 100%" filterable>
            <el-option v-for="p in projects" :key="p.id" :label="p.code + ' - ' + p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="工作时长" prop="hours">
          <el-input-number v-model="form.hours" :min="0.5" :max="24" :step="0.5" controls-position="right" style="width: 100%" />
          <div style="font-size: 12px; color: #909399; margin-top: 4px">
            当前日期已填: <span :style="{color: currentDayTotal + (form.hours||0) > 8 ? '#f56c6c' : '#67c23a', fontWeight: 600}">{{ currentDayTotal }}</span> 小时，填写后合计: 
            <span :style="{color: currentDayTotal + (form.hours||0) > 24 ? '#f56c6c' : '#409eff', fontWeight: 600}">{{ (currentDayTotal + (form.hours||0)).toFixed(1) }}</span> 小时
          </div>
        </el-form-item>
        <el-form-item label="工作内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="4" placeholder="请详细描述今日工作内容..." maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确认提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import weekday from 'dayjs/plugin/weekday'
import weekOfYear from 'dayjs/plugin/weekOfYear'
import { getProjects } from '../api/project'
import { getWorkRecords, createWorkRecord, updateWorkRecord, deleteWorkRecord, getDailyStats } from '../api/workrecord'

dayjs.extend(weekday)
dayjs.extend(weekOfYear)

const weekdays = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
const viewMode = ref('month')
const currentMonth = ref(dayjs())
const currentWeek = ref(dayjs())
const selectedDate = ref(dayjs().format('YYYY-MM-DD'))
const projects = ref([])
const records = ref([])
const dailyStats = ref([])

const dialogVisible = ref(false)
const editingRecord = ref(null)
const submitting = ref(false)
const formRef = ref()

const form = reactive({
  project_id: null,
  work_date: '',
  hours: 8,
  content: ''
})

const rules = {
  project_id: [{ required: true, message: '请选择项目', trigger: 'change' }],
  hours: [{ required: true, message: '请填写工时', trigger: 'blur' }],
  content: [{ required: true, message: '请填写工作内容', trigger: 'blur', min: 2, max: 500 }]
}

const periodText = computed(() => {
  if (viewMode.value === 'month') {
    return currentMonth.value.format('YYYY年 MM月')
  }
  const start = currentWeek.value.weekday(0)
  const end = currentWeek.value.weekday(6)
  return `${start.format('YYYY/MM/DD')} - ${end.format('YYYY/MM/DD')}`
})

const periodRange = computed(() => {
  let start, end
  if (viewMode.value === 'month') {
    start = currentMonth.value.startOf('month').format('YYYY-MM-DD')
    end = currentMonth.value.endOf('month').format('YYYY-MM-DD')
  } else {
    start = currentWeek.value.weekday(0).format('YYYY-MM-DD')
    end = currentWeek.value.weekday(6).format('YYYY-MM-DD')
  }
  return { start, end }
})

const calendarDays = computed(() => {
  const days = []
  let start, totalDays
  if (viewMode.value === 'month') {
    const firstDay = currentMonth.value.startOf('month')
    start = firstDay.weekday(0)
    totalDays = 42
  } else {
    start = currentWeek.value.weekday(0)
    totalDays = 7
  }

  for (let i = 0; i < totalDays; i++) {
    const d = start.add(i, 'day')
    const dateStr = d.format('YYYY-MM-DD')
    const isToday = dateStr === dayjs().format('YYYY-MM-DD')
    const currentMonthCheck = viewMode.value === 'month' ? d.month() === currentMonth.value.month() : true
    const dayRecords = records.value.filter(r => r.work_date === dateStr)
    const hours = dayRecords.reduce((sum, r) => sum + r.hours, 0)
    days.push({
      key: dateStr,
      dateStr,
      dayNum: d.date(),
      isToday,
      currentMonth: currentMonthCheck,
      items: dayRecords,
      hours: parseFloat(hours.toFixed(1))
    })
  }
  return days
})

const selectedDateData = computed(() => {
  const recs = records.value.filter(r => r.work_date === selectedDate.value)
  const totalHours = recs.reduce((sum, r) => sum + r.hours, 0)
  const rejectedCount = recs.filter(r => r.status === 'rejected').length
  return { records: recs, totalHours: parseFloat(totalHours.toFixed(1)), rejectedCount }
})

const currentDayTotal = computed(() => {
  const wd = form.work_date || selectedDate.value
  let total = 0
  records.value.forEach(r => {
    if (r.work_date === wd) {
      if (!editingRecord.value || r.id !== editingRecord.value.id) {
        total += r.hours
      }
    }
  })
  return parseFloat(total.toFixed(1))
})

const monthStats = computed(() => {
  const monthStart = dayjs().startOf('month').format('YYYY-MM-DD')
  const monthEnd = dayjs().endOf('month').format('YYYY-MM-DD')
  let total = 0, approved = 0, pending = 0, rejected = 0
  records.value.forEach(r => {
    if (r.work_date >= monthStart && r.work_date <= monthEnd) {
      total += r.hours
      if (r.status === 'approved') approved += r.hours
      else if (r.status === 'pending') pending += r.hours
      else if (r.status === 'rejected') rejected += r.hours
    }
  })
  return {
    total: parseFloat(total.toFixed(1)),
    approved: parseFloat(approved.toFixed(1)),
    pending: parseFloat(pending.toFixed(1)),
    rejected: parseFloat(rejected.toFixed(1))
  }
})

const prevPeriod = () => {
  if (viewMode.value === 'month') currentMonth.value = currentMonth.value.subtract(1, 'month')
  else currentWeek.value = currentWeek.value.subtract(1, 'week')
}
const nextPeriod = () => {
  if (viewMode.value === 'month') currentMonth.value = currentMonth.value.add(1, 'month')
  else currentWeek.value = currentWeek.value.add(1, 'week')
}
const goToday = () => {
  currentMonth.value = dayjs()
  currentWeek.value = dayjs()
  selectedDate.value = dayjs().format('YYYY-MM-DD')
}

const selectDate = (day) => {
  selectedDate.value = day.dateStr
}

const openAddDialog = () => {
  editingRecord.value = null
  Object.assign(form, {
    project_id: projects.value[0]?.id || null,
    work_date: selectedDate.value,
    hours: 8,
    content: ''
  })
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  editingRecord.value = row
  Object.assign(form, {
    project_id: row.project_id,
    work_date: row.work_date,
    hours: row.hours,
    content: row.content
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  if (currentDayTotal.value + (form.hours || 0) > 24) {
    ElMessage.error('单日工时合计不能超过24小时')
    return
  }
  submitting.value = true
  try {
    if (editingRecord.value) {
      await updateWorkRecord(editingRecord.value.id, { ...form })
      ElMessage.success('修改成功')
    } else {
      await createWorkRecord({ ...form })
      ElMessage.success('提交成功，等待审核')
    }
    dialogVisible.value = false
    loadRecords()
  } catch (e) {
  } finally {
    submitting.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除 ${row.work_date} ${row.project?.name} 的${row.hours}小时工时记录吗？`, '提示', { type: 'warning' })
    .then(async () => {
      await deleteWorkRecord(row.id)
      ElMessage.success('删除成功')
      loadRecords()
    }).catch(() => {})
}

const loadProjects = async () => {
  projects.value = await getProjects({ status: 'active' })
}

const loadRecords = async () => {
  const { start, end } = periodRange.value
  records.value = await getWorkRecords({ start_date: start, end_date: end })
}

const loadDailyStats = async () => {
  const start = dayjs().startOf('month').format('YYYY-MM-DD')
  const end = dayjs().endOf('month').format('YYYY-MM-DD')
  dailyStats.value = await getDailyStats({ start_date: start, end_date: end })
}

watch(viewMode, loadRecords)
watch([currentMonth, currentWeek], loadRecords)

onMounted(() => {
  loadProjects()
  loadRecords()
  loadDailyStats()
})
</script>
