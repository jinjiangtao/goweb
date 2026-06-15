<template>
  <el-dialog
    :model-value="visible"
    :title="editData ? '编辑房源' : '发布房源'"
    width="800px"
    :close-on-click-modal="false"
    @close="emit('close')"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
      <el-divider content-position="left">基本信息</el-divider>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="标题" prop="title">
            <el-input v-model="form.title" placeholder="请输入房源标题" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="小区名称" prop="community">
            <el-input v-model="form.community" placeholder="请输入小区名称" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="区域" prop="district">
            <el-select v-model="form.district" placeholder="请选择区域" style="width: 100%">
              <el-option v-for="d in districtOptions" :key="d" :label="d" :value="d" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="详细地址" prop="address">
            <el-input v-model="form.address" placeholder="请输入详细地址" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="户型" prop="room_type">
            <el-select v-model="form.room_type" placeholder="请选择" style="width: 100%">
              <el-option v-for="r in roomTypeOptions" :key="r" :label="r" :value="r" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="面积(㎡)" prop="area">
            <el-input-number v-model="form.area" :min="1" :precision="1" style="width: 100%" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="租金(元/月)" prop="price">
            <el-input-number v-model="form.price" :min="0" :precision="0" style="width: 100%" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="楼层">
            <el-input v-model="form.floor" placeholder="如：15/30" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="朝向">
            <el-select v-model="form.orientation" placeholder="请选择" clearable style="width: 100%">
              <el-option v-for="o in orientationOptions" :key="o" :label="o" :value="o" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="装修">
            <el-select v-model="form.decoration" placeholder="请选择" clearable style="width: 100%">
              <el-option v-for="d in decorationOptions" :key="d" :label="d" :value="d" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-divider content-position="left">详情信息</el-divider>

      <el-form-item label="配套设施">
        <el-checkbox-group v-model="facilityList">
          <el-checkbox v-for="f in facilityOptions" :key="f" :label="f" :value="f" />
        </el-checkbox-group>
      </el-form-item>

      <el-form-item label="房源描述" prop="description">
        <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入房源描述" />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="联系人" prop="contact_name">
            <el-input v-model="form.contact_name" placeholder="请输入联系人" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="联系电话" prop="contact_phone">
            <el-input v-model="form.contact_phone" placeholder="请输入联系电话" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-divider content-position="left">房源图片</el-divider>

      <el-form-item label="房源图片" prop="images" :error="imagesError">
        <div class="image-upload-area">
          <draggable
            v-model="imageList"
            item-key="url"
            class="image-list"
            :animation="200"
            @end="onDragEnd"
          >
            <template #item="{ element, index }">
              <div class="image-item">
                <el-image :src="element.url" fit="cover" class="image-thumb" />
                <div v-if="index === 0" class="cover-badge">封面</div>
                <div class="image-actions">
                  <el-icon class="action-icon" @click="removeImage(index)"><Delete /></el-icon>
                </div>
              </div>
            </template>
          </draggable>
          <el-upload
            v-if="imageList.length < 9"
            :show-file-list="false"
            :before-upload="beforeUpload"
            :http-request="handleUpload"
            accept="image/*"
            class="image-uploader"
          >
            <el-icon class="uploader-icon"><Plus /></el-icon>
          </el-upload>
        </div>
        <div class="image-tip">最多上传9张图片，第一张为封面图，可拖拽排序</div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="emit('close')">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, nextTick } from 'vue'
import { Delete, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import draggable from 'vuedraggable'
import { createHouse, updateHouse, uploadImage } from '../api/house'

const props = defineProps({
  visible: Boolean,
  editData: Object,
})

const emit = defineEmits(['close', 'saved'])

const formRef = ref(null)
const submitting = ref(false)
const imagesError = ref('')

const districtOptions = ['东城区', '西城区', '朝阳区', '海淀区', '丰台区', '石景山区', '通州区', '顺义区', '大兴区', '昌平区', '房山区', '门头沟区', '其他']
const roomTypeOptions = ['一室一厅', '两室一厅', '两室两厅', '三室一厅', '三室两厅', '四室一厅', '四室两厅', '一室', '两室', '三室', '其他']
const orientationOptions = ['东', '南', '西', '北', '东南', '西南', '东北', '西北', '南北']
const decorationOptions = ['毛坯', '简装', '精装', '豪装']
const facilityOptions = ['WiFi', '空调', '冰箱', '洗衣机', '暖气', '热水器', '电视', '微波炉', '油烟机', '床', '衣柜', '沙发', '书桌']

const form = reactive({
  title: '',
  community: '',
  district: '',
  address: '',
  room_type: '',
  area: 0,
  price: 0,
  floor: '',
  orientation: '',
  decoration: '',
  facilities: '',
  description: '',
  contact_name: '',
  contact_phone: '',
  images: '',
})

const facilityList = ref([])
const imageList = ref([])

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  community: [{ required: true, message: '请输入小区名称', trigger: 'blur' }],
  district: [{ required: true, message: '请选择区域', trigger: 'change' }],
  address: [{ required: true, message: '请输入详细地址', trigger: 'blur' }],
  room_type: [{ required: true, message: '请选择户型', trigger: 'change' }],
  area: [{ required: true, message: '请输入面积', trigger: 'blur' }],
  price: [{ required: true, message: '请输入租金', trigger: 'blur' }],
  description: [{ required: true, message: '请输入房源描述', trigger: 'blur' }],
  contact_name: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  contact_phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
}

watch(() => props.visible, (val) => {
  if (val) {
    imagesError.value = ''
    if (props.editData) {
      const d = props.editData
      form.title = d.title || ''
      form.community = d.community || ''
      form.district = d.district || ''
      form.address = d.address || ''
      form.room_type = d.room_type || ''
      form.area = d.area || 0
      form.price = d.price || 0
      form.floor = d.floor || ''
      form.orientation = d.orientation || ''
      form.decoration = d.decoration || ''
      form.description = d.description || ''
      form.contact_name = d.contact_name || ''
      form.contact_phone = d.contact_phone || ''

      facilityList.value = d.facilities ? d.facilities.split(',').filter(Boolean) : []

      imageList.value = (d.images || []).map((url) => ({ url }))
    } else {
      Object.assign(form, {
        title: '', community: '', district: '', address: '',
        room_type: '', area: 0, price: 0, floor: '', orientation: '',
        decoration: '', facilities: '', description: '', contact_name: '',
        contact_phone: '', images: '',
      })
      facilityList.value = []
      imageList.value = []
    }
    nextTick(() => {
      formRef.value?.clearValidate()
    })
  }
}, { immediate: true })

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB')
    return false
  }
  return true
}

const handleUpload = async (options) => {
  try {
    const res = await uploadImage(options.file)
    imageList.value.push({ url: res.data.url })
    imagesError.value = ''
  } catch {
    ElMessage.error('图片上传失败')
  }
}

const removeImage = (index) => {
  imageList.value.splice(index, 1)
}

const onDragEnd = () => {}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  if (imageList.value.length === 0) {
    imagesError.value = '请至少上传1张房源图片'
    return
  }
  imagesError.value = ''

  form.facilities = facilityList.value.join(',')
  form.images = JSON.stringify(imageList.value.map((item) => item.url))

  submitting.value = true
  try {
    if (props.editData) {
      await updateHouse(props.editData.id, form)
      ElMessage.success('编辑成功')
    } else {
      await createHouse(form)
      ElMessage.success('发布成功')
    }
    emit('saved')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.image-upload-area {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: flex-start;
}

.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.image-item {
  position: relative;
  width: 104px;
  height: 104px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  overflow: hidden;
  cursor: move;
}

.image-thumb {
  width: 100%;
  height: 100%;
}

.cover-badge {
  position: absolute;
  top: 0;
  left: 0;
  background: #409eff;
  color: #fff;
  font-size: 12px;
  padding: 0 6px;
  line-height: 20px;
}

.image-actions {
  position: absolute;
  top: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.5);
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s;
}

.image-item:hover .image-actions {
  opacity: 1;
}

.action-icon {
  color: #fff;
  font-size: 14px;
}

.image-uploader {
  width: 104px;
  height: 104px;
  border: 1px dashed #dcdfe6;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s;
}

.image-uploader:hover {
  border-color: #409eff;
}

.uploader-icon {
  font-size: 28px;
  color: #8c939d;
}

.image-tip {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
</style>
