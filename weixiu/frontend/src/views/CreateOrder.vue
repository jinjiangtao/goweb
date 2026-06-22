<template>
  <div class="create-order">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>提交报修工单</span>
        </div>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" class="order-form">
        <el-form-item label="故障标题" prop="title">
          <el-input v-model="form.title" placeholder="请简要描述故障问题" maxlength="200" show-word-limit />
        </el-form-item>

        <el-form-item label="设备类型" prop="device_type">
          <el-select v-model="form.device_type" placeholder="请选择设备类型" style="width: 100%;">
            <el-option label="电脑" value="computer" />
            <el-option label="打印机" value="printer" />
            <el-option label="网络设备" value="network" />
            <el-option label="投影仪" value="projector" />
            <el-option label="空调" value="aircon" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>

        <el-form-item label="设备名称">
          <el-input v-model="form.device_name" placeholder="请输入设备名称或型号" />
        </el-form-item>

        <el-form-item label="报修区域" prop="area">
          <el-select v-model="form.area" placeholder="请选择报修区域" style="width: 100%;">
            <el-option label="教学楼A区" value="教学楼A区" />
            <el-option label="教学楼B区" value="教学楼B区" />
            <el-option label="实验楼" value="实验楼" />
            <el-option label="办公楼" value="办公楼" />
            <el-option label="图书馆" value="图书馆" />
            <el-option label="食堂" value="食堂" />
            <el-option label="宿舍区" value="宿舍区" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>

        <el-form-item label="具体位置">
          <el-input v-model="form.location" placeholder="请输入具体位置，如：教学楼A区302室" />
        </el-form-item>

        <el-form-item label="紧急程度">
          <el-radio-group v-model="form.priority">
            <el-radio value="low">低</el-radio>
            <el-radio value="medium">中</el-radio>
            <el-radio value="high">高</el-radio>
            <el-radio value="urgent">紧急</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="故障描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="5"
            placeholder="请详细描述故障现象"
            maxlength="1000"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="故障图片">
          <el-upload
            :action="uploadUrl"
            :headers="uploadHeaders"
            list-type="picture-card"
            :on-success="handleUploadSuccess"
            :on-remove="handleRemove"
            :file-list="fileList"
            multiple
            accept="image/*"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <div class="upload-tip">支持多张图片上传，单张图片不超过5MB</div>
        </el-form-item>

        <el-form-item label="联系人">
          <el-input v-model="form.contact_name" placeholder="请输入联系人姓名" />
        </el-form-item>

        <el-form-item label="联系电话">
          <el-input v-model="form.contact_phone" placeholder="请输入联系电话" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">提交工单</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { createOrder, uploadFile } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)
const fileList = ref([])

const form = reactive({
  title: '',
  description: '',
  device_type: '',
  device_name: '',
  area: '',
  location: '',
  priority: 'medium',
  contact_name: userStore.user?.real_name || '',
  contact_phone: userStore.user?.phone || '',
  fault_images: ''
})

const uploadUrl = computed(() => '/api/upload')
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
}))

const rules = {
  title: [{ required: true, message: '请输入故障标题', trigger: 'blur' }],
  device_type: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
  area: [{ required: true, message: '请选择报修区域', trigger: 'change' }],
  description: [{ required: true, message: '请输入故障描述', trigger: 'blur' }]
}

const handleUploadSuccess = (res, file) => {
  fileList.value.push({ url: res.url, name: file.name, uid: file.uid })
}

const handleRemove = (file) => {
  const index = fileList.value.findIndex(f => f.uid === file.uid)
  if (index > -1) fileList.value.splice(index, 1)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        form.fault_images = fileList.value.map(f => f.url).join(',')
        await createOrder(form)
        ElMessage.success('工单提交成功')
        router.push('/orders')
      } finally {
        loading.value = false
      }
    }
  })
}

const handleReset = () => {
  formRef.value?.resetFields()
  fileList.value = []
}
</script>

<style scoped>
.create-order {
  max-width: 800px;
}

.card-header {
  font-weight: 600;
  font-size: 16px;
}

.order-form {
  padding: 20px 0;
}

.upload-tip {
  color: #999;
  font-size: 12px;
  margin-top: 5px;
}
</style>
