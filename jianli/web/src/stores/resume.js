import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getResumeDetail,
  updateResume,
  createResume,
  getVersions,
  restoreVersion as apiRestoreVersion,
  getVersion
} from '@/api/resume'

const MODULE_PRESETS = [
  { type: 'basic', name: '基本信息', removable: false, visible: true },
  { type: 'education', name: '教育经历', removable: true, visible: true },
  { type: 'experience', name: '工作经历', removable: true, visible: true },
  { type: 'projects', name: '项目经历', removable: true, visible: true },
  { type: 'skills', name: '技能特长', removable: true, visible: true },
  { type: 'summary', name: '个人总结', removable: true, visible: true },
  { type: 'evaluation', name: '自我评价', removable: true, visible: true }
]

const DEFAULT_CONTENT = {
  basic: {
    name: '',
    phone: '',
    email: '',
    city: '',
    job_intention: '',
    github: '',
    website: ''
  },
  education: [],
  experience: [],
  projects: [],
  skills: [],
  summary: '',
  evaluation: ''
}

const DEFAULT_STYLE = {
  primaryColor: '#409eff',
  accentColor: '#0f3460',
  fontSize: 14,
  lineSpacing: 1.6,
  marginTop: 20,
  marginBottom: 20,
  marginLeft: 25,
  marginRight: 25
}

function generateId() {
  return Date.now().toString(36) + Math.random().toString(36).substr(2)
}

function cloneDeep(obj) {
  return JSON.parse(JSON.stringify(obj))
}

export const useResumeStore = defineStore('resume', () => {
  const resumeList = ref([])
  const templates = ref([])
  const userIdentity = ref('default_user')
  const currentResume = ref(null)
  const resumeId = ref(null)
  const isSaving = ref(false)
  const versions = ref([])
  const selectedModuleId = ref(null)

  const modules = ref([
    { id: generateId(), ...MODULE_PRESETS[0] },
    { id: generateId(), ...MODULE_PRESETS[1] },
    { id: generateId(), ...MODULE_PRESETS[2] }
  ])

  const content = ref(cloneDeep(DEFAULT_CONTENT))
  const styleConfig = ref(cloneDeep(DEFAULT_STYLE))

  const visibleModules = computed(() => modules.value.filter(m => m.visible))

  function setResumeList(list) {
    resumeList.value = list
  }

  function setTemplates(list) {
    templates.value = list
  }

  function setUserIdentity(identity) {
    userIdentity.value = identity
  }

  function setCurrentResume(resume) {
    currentResume.value = resume
  }

  function clearCurrentResume() {
    currentResume.value = null
    resumeId.value = null
    modules.value = [
      { id: generateId(), ...MODULE_PRESETS[0] },
      { id: generateId(), ...MODULE_PRESETS[1] },
      { id: generateId(), ...MODULE_PRESETS[2] }
    ]
    content.value = cloneDeep(DEFAULT_CONTENT)
    styleConfig.value = cloneDeep(DEFAULT_STYLE)
    versions.value = []
    selectedModuleId.value = null
  }

  function getAvailableModules() {
    const existingTypes = modules.value.map(m => m.type)
    return MODULE_PRESETS.filter(m => !existingTypes.includes(m.type) && m.removable)
  }

  function selectModule(id) {
    selectedModuleId.value = id
  }

  function toggleModuleVisibility(id) {
    const module = modules.value.find(m => m.id === id)
    if (module) {
      module.visible = !module.visible
    }
  }

  function removeModule(id) {
    const index = modules.value.findIndex(m => m.id === id)
    if (index > -1 && modules.value[index].removable) {
      modules.value.splice(index, 1)
      if (selectedModuleId.value === id) {
        selectedModuleId.value = modules.value[0]?.id || null
      }
    }
  }

  function addModule(type) {
    const preset = MODULE_PRESETS.find(m => m.type === type)
    if (!preset) return
    const existingTypes = modules.value.map(m => m.type)
    if (existingTypes.includes(type)) return

    const newModule = {
      id: generateId(),
      ...preset
    }
    modules.value.push(newModule)
    selectedModuleId.value = newModule.id

    if (!content.value[type]) {
      if (type === 'summary' || type === 'evaluation') {
        content.value[type] = ''
      } else {
        content.value[type] = []
      }
    }
  }

  function reorderModules(newModules) {
    modules.value = newModules
  }

  function updateContent(type, data) {
    content.value[type] = data
  }

  function updateStyle(key, value) {
    styleConfig.value[key] = value
  }

  function updateModule(type, data) {
    content.value[type] = data
  }

  function addArrayItem(moduleType) {
    if (!content.value[moduleType]) {
      content.value[moduleType] = []
    }
    const newItem = {}
    if (moduleType === 'education') {
      Object.assign(newItem, { school: '', degree: '', major: '', start_date: '', end_date: '', description: '' })
    } else if (moduleType === 'experience') {
      Object.assign(newItem, { company: '', position: '', start_date: '', end_date: '', description: '', highlights: [] })
    } else if (moduleType === 'skills') {
      Object.assign(newItem, { category: '', items: '' })
    } else if (moduleType === 'projects') {
      Object.assign(newItem, { name: '', role: '', start_date: '', end_date: '', description: '', highlights: [] })
    }
    content.value[moduleType].push(newItem)
  }

  function removeArrayItem(moduleType, index) {
    if (content.value[moduleType] && content.value[moduleType].length > index) {
      content.value[moduleType].splice(index, 1)
    }
  }

  function addHighlight(moduleType, itemIndex) {
    if (content.value[moduleType] && content.value[moduleType][itemIndex]) {
      if (!content.value[moduleType][itemIndex].highlights) {
        content.value[moduleType][itemIndex].highlights = []
      }
      content.value[moduleType][itemIndex].highlights.push('')
    }
  }

  function removeHighlight(moduleType, itemIndex, highlightIndex) {
    if (content.value[moduleType] && content.value[moduleType][itemIndex]?.highlights) {
      content.value[moduleType][itemIndex].highlights.splice(highlightIndex, 1)
    }
  }

  async function loadResume(id) {
    if (id === 'new') {
      resumeId.value = 'new'
      return
    }
    try {
      const data = await getResumeDetail(id)
      resumeId.value = id
      currentResume.value = data

      if (data.title) {
        currentResume.value.title = data.title
      }

      if (data.content && data.content.modules && Array.isArray(data.content.modules)) {
        modules.value = data.content.modules.map(m => ({
          id: m.id || generateId(),
          type: m.type,
          name: m.name || MODULE_PRESETS.find(p => p.type === m.type)?.name || m.type,
          removable: m.removable !== false,
          visible: m.visible !== false,
          order: m.order || 0
        }))
        modules.value.sort((a, b) => (a.order || 0) - (b.order || 0))
        delete data.content.modules
      }

      if (data.content) {
        content.value = { ...cloneDeep(DEFAULT_CONTENT), ...data.content }
      }

      if (data.style_config) {
        styleConfig.value = { ...cloneDeep(DEFAULT_STYLE), ...data.style_config }
      }

      if (modules.value.length > 0) {
        selectedModuleId.value = modules.value[0].id
      }
    } catch (e) {
      console.error('加载简历失败:', e)
      throw e
    }
  }

  async function saveResume() {
    if (!resumeId.value || resumeId.value === 'new') {
      return
    }
    isSaving.value = true
    try {
      const data = {
        title: currentResume.value?.title || '未命名简历',
        modules: modules.value,
        content: content.value,
        style: styleConfig.value
      }
      const result = await updateResume(resumeId.value, data)
      currentResume.value = { ...currentResume.value, ...data }
      return result
    } catch (e) {
      console.error('保存简历失败:', e)
      throw e
    } finally {
      isSaving.value = false
    }
  }

  async function createNewResume(title) {
    isSaving.value = true
    try {
      const data = {
        title: title || '未命名简历',
        template_id: 1,
        user_identity: userIdentity.value
      }
      const result = await createResume(data)
      resumeId.value = result?.id
      currentResume.value = result
      return result?.id
    } catch (e) {
      console.error('创建简历失败:', e)
      throw e
    } finally {
      isSaving.value = false
    }
  }

  async function loadVersions(id) {
    try {
      const data = await getVersions(id)
      versions.value = data || []
      return versions.value
    } catch (e) {
      console.error('加载版本列表失败:', e)
      throw e
    }
  }

  async function restoreVersion(versionId) {
    try {
      const restoredData = await apiRestoreVersion(versionId)

      if (restoredData && restoredData.resume) {
        resumeId.value = restoredData.resume.id
        await loadResume(restoredData.resume.id)
      }

      return restoredData
    } catch (e) {
      console.error('恢复版本失败:', e)
      throw e
    }
  }

  const selectedModule = computed(() => modules.value.find(m => m.id === selectedModuleId.value))

  return {
    resumeList,
    templates,
    userIdentity,
    currentResume,
    resumeId,
    isSaving,
    versions,
    modules,
    selectedModuleId,
    selectedModule,
    content,
    styleConfig,
    visibleModules,
    setResumeList,
    setTemplates,
    setUserIdentity,
    setCurrentResume,
    clearCurrentResume,
    getAvailableModules,
    selectModule,
    toggleModuleVisibility,
    removeModule,
    addModule,
    reorderModules,
    updateContent,
    updateStyle,
    updateModule,
    addArrayItem,
    removeArrayItem,
    addHighlight,
    removeHighlight,
    loadResume,
    saveResume,
    createNewResume,
    loadVersions,
    restoreVersion
  }
})
