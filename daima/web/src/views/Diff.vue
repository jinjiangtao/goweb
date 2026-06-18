<template>
  <div class="page-container">
    <div class="card">
      <h2 class="mb-20" style="font-size:20px;">代码差异对比</h2>
      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="编程语言">
            <el-select v-model="language" size="large" filterable style="width:100%">
              <el-option v-for="lang in languages" :key="lang" :label="lang" :value="lang" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="16">
          <div class="flex gap-8 mb-12 justify-end">
            <el-button size="large" :icon="Refresh" plain @click="clearAll">清空</el-button>
            <el-button size="large" type="primary" :icon="Comparison" @click="compareNow">对比</el-button>
          </div>
        </el-col>
      </el-row>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="分别输入" name="split">
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="mb-12" style="font-weight:500;color:#6b7280;">原始代码</div>
              <div style="height: 450px;">
                <CodeEditor v-model="originalCode" :language="language" />
              </div>
            </el-col>
            <el-col :span="12">
              <div class="mb-12" style="font-weight:500;color:#6b7280;">修改后代码</div>
              <div style="height: 450px;">
                <CodeEditor v-model="modifiedCode" :language="language" />
              </div>
            </el-col>
          </el-row>
        </el-tab-pane>
        <el-tab-pane label="对比视图" name="diff">
          <div style="height: 500px;">
            <CodeEditor
              v-model="modifiedCode"
              :original-value="originalCode"
              :language="language"
              :diff-mode="true"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { snippetApi } from '../api'
import CodeEditor from '../components/CodeEditor.vue'

const activeTab = ref('split')
const languages = ref([])
const language = ref('javascript')
const originalCode = ref(`function hello() {
  console.log("Hello World");
}`)
const modifiedCode = ref(`function hello(name) {
  console.log("Hello " + name);
  return true;
}`)

const fetchLanguages = async () => {
  try {
    const res = await snippetApi.getLanguages()
    languages.value = res.data
  } catch (e) {}
}

const compareNow = () => {
  activeTab.value = 'diff'
}

const clearAll = () => {
  originalCode.value = ''
  modifiedCode.value = ''
  activeTab.value = 'split'
}

watch([originalCode, modifiedCode, language], () => {})

onMounted(fetchLanguages)
</script>
