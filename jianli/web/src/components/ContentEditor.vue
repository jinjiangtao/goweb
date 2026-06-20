<template>
  <div class="content-editor">
    <div class="editor-header">
      <span class="title">{{ selectedModule?.name || '请选择模块' }}</span>
    </div>

    <div class="editor-body" v-if="selectedModule">
      <template v-if="selectedModule.type === 'basic'">
        <el-form :model="basicForm" label-width="100px" class="edit-form">
          <el-form-item label="姓名">
            <el-input v-model="basicForm.name" placeholder="请输入姓名" />
          </el-form-item>
          <el-form-item label="电话">
            <el-input v-model="basicForm.phone" placeholder="请输入电话" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="basicForm.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="城市">
            <el-input v-model="basicForm.city" placeholder="请输入城市" />
          </el-form-item>
          <el-form-item label="求职意向">
            <el-input v-model="basicForm.job_intention" placeholder="请输入求职意向" />
          </el-form-item>
          <el-form-item label="GitHub">
            <el-input v-model="basicForm.github" placeholder="请输入GitHub链接" />
          </el-form-item>
          <el-form-item label="个人网站">
            <el-input v-model="basicForm.website" placeholder="请输入个人网站链接" />
          </el-form-item>
        </el-form>
      </template>

      <template v-else-if="selectedModule.type === 'education'">
        <div class="array-list">
          <div v-for="(item, index) in educationList" :key="index" class="array-item">
            <div class="item-header">
              <span class="item-title">教育经历 {{ index + 1 }}</span>
              <el-button type="danger" link @click="removeArrayItem(selectedModule.type, index)">
                删除
              </el-button>
            </div>
            <el-form :model="item" label-width="80px" class="edit-form">
              <el-form-item label="学校">
                <el-input v-model="item.school" placeholder="请输入学校名称" />
              </el-form-item>
              <el-form-item label="学历">
                <el-input v-model="item.degree" placeholder="请输入学历" />
              </el-form-item>
              <el-form-item label="专业">
                <el-input v-model="item.major" placeholder="请输入专业" />
              </el-form-item>
              <el-form-item label="开始时间">
                <el-date-picker
                  v-model="item.start_date"
                  type="month"
                  placeholder="选择开始时间"
                  value-format="YYYY-MM"
                />
              </el-form-item>
              <el-form-item label="结束时间">
                <el-date-picker
                  v-model="item.end_date"
                  type="month"
                  placeholder="选择结束时间"
                  value-format="YYYY-MM"
                />
              </el-form-item>
              <el-form-item label="描述">
                <el-input
                  v-model="item.description"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入描述"
                />
              </el-form-item>
            </el-form>
          </div>
          <el-button type="primary" plain class="add-item-btn" @click="addArrayItem(selectedModule.type)">
            <el-icon><Plus /></el-icon>
            添加教育经历
          </el-button>
        </div>
      </template>

      <template v-else-if="selectedModule.type === 'experience'">
        <div class="array-list">
          <div v-for="(item, index) in experienceList" :key="index" class="array-item">
            <div class="item-header">
              <span class="item-title">工作经历 {{ index + 1 }}</span>
              <el-button type="danger" link @click="removeArrayItem(selectedModule.type, index)">
                删除
              </el-button>
            </div>
            <el-form :model="item" label-width="80px" class="edit-form">
              <el-form-item label="公司">
                <el-input v-model="item.company" placeholder="请输入公司名称" />
              </el-form-item>
              <el-form-item label="职位">
                <el-input v-model="item.position" placeholder="请输入职位" />
              </el-form-item>
              <el-form-item label="开始时间">
                <el-date-picker
                  v-model="item.start_date"
                  type="month"
                  placeholder="选择开始时间"
                  value-format="YYYY-MM"
                />
              </el-form-item>
              <el-form-item label="结束时间">
                <el-date-picker
                  v-model="item.end_date"
                  type="month"
                  placeholder="选择结束时间"
                  value-format="YYYY-MM"
                />
              </el-form-item>
              <el-form-item label="工作描述">
                <el-input
                  v-model="item.description"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入工作描述"
                />
              </el-form-item>
              <el-form-item label="工作亮点">
                <div class="highlights-list">
                  <div v-for="(hl, hlIndex) in item.highlights" :key="hlIndex" class="highlight-item">
                    <el-input v-model="item.highlights[hlIndex]" placeholder="请输入工作亮点" />
                    <el-button type="danger" link @click="removeHighlight(selectedModule.type, index, hlIndex)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                  <el-button type="primary" link @click="addHighlight(selectedModule.type, index)">
                    <el-icon><Plus /></el-icon>
                    添加亮点
                  </el-button>
                </div>
              </el-form-item>
            </el-form>
          </div>
          <el-button type="primary" plain class="add-item-btn" @click="addArrayItem(selectedModule.type)">
            <el-icon><Plus /></el-icon>
            添加工作经历
          </el-button>
        </div>
      </template>

      <template v-else-if="selectedModule.type === 'skills'">
        <div class="array-list">
          <div v-for="(item, index) in skillsList" :key="index" class="array-item">
            <div class="item-header">
              <span class="item-title">技能分类 {{ index + 1 }}</span>
              <el-button type="danger" link @click="removeArrayItem(selectedModule.type, index)">
                删除
              </el-button>
            </div>
            <el-form :model="item" label-width="80px" class="edit-form">
              <el-form-item label="分类名称">
                <el-input v-model="item.category" placeholder="如：前端开发、后端开发" />
              </el-form-item>
              <el-form-item label="技能内容">
                <el-input
                  v-model="item.items"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入技能，用逗号分隔"
                />
              </el-form-item>
            </el-form>
          </div>
          <el-button type="primary" plain class="add-item-btn" @click="addArrayItem(selectedModule.type)">
            <el-icon><Plus /></el-icon>
            添加技能分类
          </el-button>
        </div>
      </template>

      <template v-else-if="selectedModule.type === 'projects'">
        <div class="array-list">
          <div v-for="(item, index) in projectsList" :key="index" class="array-item">
            <div class="item-header">
              <span class="item-title">项目经验 {{ index + 1 }}</span>
              <el-button type="danger" link @click="removeArrayItem(selectedModule.type, index)">
                删除
              </el-button>
            </div>
            <el-form :model="item" label-width="80px" class="edit-form">
              <el-form-item label="项目名称">
                <el-input v-model="item.name" placeholder="请输入项目名称" />
              </el-form-item>
              <el-form-item label="担任角色">
                <el-input v-model="item.role" placeholder="请输入担任角色" />
              </el-form-item>
              <el-form-item label="开始时间">
                <el-date-picker
                  v-model="item.start_date"
                  type="month"
                  placeholder="选择开始时间"
                  value-format="YYYY-MM"
                />
              </el-form-item>
              <el-form-item label="结束时间">
                <el-date-picker
                  v-model="item.end_date"
                  type="month"
                  placeholder="选择结束时间"
                  value-format="YYYY-MM"
                />
              </el-form-item>
              <el-form-item label="项目描述">
                <el-input
                  v-model="item.description"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入项目描述"
                />
              </el-form-item>
              <el-form-item label="项目亮点">
                <div class="highlights-list">
                  <div v-for="(hl, hlIndex) in item.highlights" :key="hlIndex" class="highlight-item">
                    <el-input v-model="item.highlights[hlIndex]" placeholder="请输入项目亮点" />
                    <el-button type="danger" link @click="removeHighlight(selectedModule.type, index, hlIndex)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                  <el-button type="primary" link @click="addHighlight(selectedModule.type, index)">
                    <el-icon><Plus /></el-icon>
                    添加亮点
                  </el-button>
                </div>
              </el-form-item>
            </el-form>
          </div>
          <el-button type="primary" plain class="add-item-btn" @click="addArrayItem(selectedModule.type)">
            <el-icon><Plus /></el-icon>
            添加项目经验
          </el-button>
        </div>
      </template>

      <template v-else-if="selectedModule.type === 'summary' || selectedModule.type === 'evaluation'">
        <el-form label-width="80px" class="edit-form">
          <el-form-item :label="selectedModule.type === 'summary' ? '个人简介' : '自我评价'">
            <el-input
              :model-value="summaryText"
              @update:model-value="updateSummary"
              type="textarea"
              :rows="10"
              placeholder="请输入内容"
            />
          </el-form-item>
        </el-form>
      </template>
    </div>

    <div v-else class="empty-state">
      <el-empty description="请选择左侧模块进行编辑" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useResumeStore } from '@/stores/resume'
import { Plus, Delete } from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const store = useResumeStore()

const selectedModule = computed(() => store.selectedModule)

const basicForm = computed({
  get: () => store.content.basic || {},
  set: (val) => store.updateModule('basic', val)
})

const educationList = computed({
  get: () => store.content.education || [],
  set: (val) => store.updateModule('education', val)
})

const experienceList = computed({
  get: () => store.content.experience || [],
  set: (val) => store.updateModule('experience', val)
})

const skillsList = computed({
  get: () => store.content.skills || [],
  set: (val) => store.updateModule('skills', val)
})

const projectsList = computed({
  get: () => store.content.projects || [],
  set: (val) => store.updateModule('projects', val)
})

const summaryText = computed({
  get: () => {
    const type = selectedModule.value?.type
    return type ? (store.content[type] || '') : ''
  },
  set: (val) => {
    const type = selectedModule.value?.type
    if (type) store.updateModule(type, val)
  }
})

function updateSummary(val) {
  const type = selectedModule.value?.type
  if (type) store.updateModule(type, val)
}

function addArrayItem(moduleType) {
  store.addArrayItem(moduleType)
}

function removeArrayItem(moduleType, index) {
  ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    store.removeArrayItem(moduleType, index)
  }).catch(() => {})
}

function addHighlight(moduleType, itemIndex) {
  store.addHighlight(moduleType, itemIndex)
}

function removeHighlight(moduleType, itemIndex, highlightIndex) {
  store.removeHighlight(moduleType, itemIndex, highlightIndex)
}
</script>

<style lang="scss" scoped>
.content-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-right: 1px solid #e4e7ed;

  .editor-header {
    padding: 16px;
    border-bottom: 1px solid #e4e7ed;
    font-weight: 600;
    font-size: 14px;
  }

  .editor-body {
    flex: 1;
    padding: 20px;
    overflow-y: auto;
  }

  .edit-form {
    .el-form-item {
      margin-bottom: 18px;
    }
  }

  .array-list {
    .array-item {
      border: 1px solid #e4e7ed;
      border-radius: 8px;
      padding: 16px;
      margin-bottom: 16px;
      background: #fafafa;

      .item-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;

        .item-title {
          font-weight: 600;
          font-size: 14px;
          color: #303133;
        }
      }
    }

    .add-item-btn {
      width: 100%;
    }
  }

  .highlights-list {
    .highlight-item {
      display: flex;
      gap: 8px;
      margin-bottom: 8px;
    }
  }

  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
