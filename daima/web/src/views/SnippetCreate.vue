<template>
  <div class="page-container">
    <div class="card">
      <h2 class="mb-20" style="font-size:20px;">{{ isEdit ? '编辑代码片段' : '创建代码片段' }}</h2>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入代码片段标题" size="large" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="2"
            placeholder="简要描述这个代码片段"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="语言" prop="language">
              <el-select v-model="form.language" placeholder="选择编程语言" size="large" filterable style="width:100%">
                <el-option v-for="lang in languages" :key="lang" :label="lang" :value="lang" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="可见性" prop="visibility">
              <el-radio-group v-model="form.visibility" size="large">
                <el-radio-button value="public">公开</el-radio-button>
                <el-radio-button value="private">私密</el-radio-button>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="标签">
              <el-input
                v-model="form.tags"
                placeholder="用逗号分隔，如: go,gin,api"
                size="large"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="代码" prop="content">
          <div style="width:100%">
            <div class="flex-between mb-12">
              <div class="flex gap-8">
                <el-button size="small" :icon="MagicStick" type="primary" plain @click="formatCode">
                  格式化代码
                </el-button>
              </div>
            </div>
            <div style="height: 500px;">
              <CodeEditor ref="editorRef" v-model="form.content" :language="form.language" />
            </div>
          </div>
        </el-form-item>
        <el-form-item>
          <div class="flex gap-12 justify-end">
            <el-button size="large" @click="$router.back()">取消</el-button>
            <el-button size="large" type="primary" :loading="saving" @click="submit">
              {{ isEdit ? '保存修改' : '创建片段' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { snippetApi } from '../api'
import CodeEditor from '../components/CodeEditor.vue'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const formRef = ref(null)
const editorRef = ref(null)
const saving = ref(false)
const languages = ref([])

const form = reactive({
  title: '',
  description: '',
  language: 'javascript',
  visibility: 'public',
  content: '',
  tags: ''
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  language: [{ required: true, message: '请选择语言', trigger: 'change' }],
  content: [{ required: true, message: '请输入代码内容', trigger: 'change' }]
}

const formatCode = () => editorRef.value?.formatCode()

const fetchLanguages = async () => {
  try {
    const res = await snippetApi.getLanguages()
    languages.value = res.data
  } catch (e) {}
}

const loadSnippet = async () => {
  try {
    const res = await snippetApi.get(route.params.id)
    const s = res.data.snippet
    form.title = s.title
    form.description = s.description
    form.language = s.language
    form.visibility = s.visibility
    form.content = s.content
    form.tags = s.tags
  } catch (e) {
    ElMessage.error('加载片段失败')
    router.push('/')
  }
}

const submit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    let res
    if (isEdit.value) {
      res = await snippetApi.update(route.params.id, form)
      ElMessage.success('保存成功')
    } else {
      res = await snippetApi.create(form)
      ElMessage.success('创建成功')
    }
    router.push(`/snippet/${res.data.id}`)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchLanguages()
  if (isEdit.value) {
    loadSnippet()
  }
})
</script>
