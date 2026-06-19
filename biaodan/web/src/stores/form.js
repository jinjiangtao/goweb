import { defineStore } from 'pinia'

const generateId = () => 'field_' + Math.random().toString(36).substr(2, 9)

export const useFormStore = defineStore('form', {
  state: () => ({
    templateId: null,
    name: '',
    description: '',
    fields: [],
    selectedFieldId: null,
    copiedField: null
  }),

  getters: {
    selectedField: (state) => {
      return state.fields.find(f => f.id === state.selectedFieldId) || null
    }
  },

  actions: {
    reset() {
      this.templateId = null
      this.name = ''
      this.description = ''
      this.fields = []
      this.selectedFieldId = null
      this.copiedField = null
    },

    loadTemplate(template) {
      this.templateId = template.id
      this.name = template.name
      this.description = template.description
      this.fields = template.schema?.fields || []
      this.selectedFieldId = null
    },

    addField(fieldType, index = -1) {
      const field = createField(fieldType)
      if (index < 0 || index >= this.fields.length) {
        this.fields.push(field)
      } else {
        this.fields.splice(index, 0, field)
      }
      this.selectedFieldId = field.id
      return field
    },

    removeField(id) {
      const idx = this.fields.findIndex(f => f.id === id)
      if (idx > -1) {
        this.fields.splice(idx, 1)
        if (this.selectedFieldId === id) {
          this.selectedFieldId = null
        }
      }
    },

    duplicateField(id) {
      const field = this.fields.find(f => f.id === id)
      if (field) {
        const newField = JSON.parse(JSON.stringify(field))
        newField.id = generateId()
        const idx = this.fields.findIndex(f => f.id === id)
        this.fields.splice(idx + 1, 0, newField)
        this.selectedFieldId = newField.id
      }
    },

    moveField(fromIndex, toIndex) {
      if (toIndex < 0 || toIndex >= this.fields.length) return
      const [field] = this.fields.splice(fromIndex, 1)
      this.fields.splice(toIndex, 0, field)
    },

    selectField(id) {
      this.selectedFieldId = id
    },

    updateField(id, updates) {
      const field = this.fields.find(f => f.id === id)
      if (field) {
        Object.assign(field, updates)
      }
    },

    getSchema() {
      return {
        fields: JSON.parse(JSON.stringify(this.fields))
      }
    }
  }
})

function createField(type) {
  const base = {
    id: generateId(),
    type,
    label: getDefaultLabel(type),
    required: false,
    placeholder: '',
    defaultValue: ''
  }

  switch (type) {
    case 'input':
      return { ...base, inputType: 'text', maxLength: 200 }
    case 'textarea':
      return { ...base, rows: 4, maxLength: 1000 }
    case 'number':
      return { ...base, min: null, max: null, step: 1 }
    case 'email':
      return { ...base, inputType: 'email' }
    case 'phone':
      return { ...base, pattern: '' }
    case 'radio':
      return { ...base, options: [{ label: '选项1', value: 'opt1' }, { label: '选项2', value: 'opt2' }] }
    case 'checkbox':
      return { ...base, options: [{ label: '选项1', value: 'opt1' }, { label: '选项2', value: 'opt2' }] }
    case 'select':
      return { ...base, options: [{ label: '选项1', value: 'opt1' }, { label: '选项2', value: 'opt2' }] }
    case 'date':
      return { ...base, dateType: 'date', range: false }
    case 'time':
      return { ...base, format: 'HH:mm' }
    case 'upload':
      return { ...base, multiple: false, accept: '', maxSize: 50 }
    case 'rate':
      return { ...base, max: 5, allowHalf: false }
    case 'switch':
      return { ...base, activeText: '是', inactiveText: '否' }
    case 'slider':
      return { ...base, min: 0, max: 100, step: 1, showInput: true }
    case 'divider':
      return { id: generateId(), type, label: '分割线', content: '' }
    case 'description':
      return { id: generateId(), type, label: '', content: '说明文字内容' }
    default:
      return base
  }
}

function getDefaultLabel(type) {
  const labels = {
    input: '单行文本',
    textarea: '多行文本',
    number: '数字',
    email: '邮箱',
    phone: '手机号',
    radio: '单选',
    checkbox: '多选',
    select: '下拉选择',
    date: '日期',
    time: '时间',
    upload: '文件上传',
    rate: '评分',
    switch: '开关',
    slider: '滑块',
    divider: '',
    description: ''
  }
  return labels[type] || '组件'
}

export const componentTypes = [
  { group: '基础组件', items: [
    { type: 'input', label: '单行文本', icon: 'Edit' },
    { type: 'textarea', label: '多行文本', icon: 'Document' },
    { type: 'number', label: '数字', icon: 'Calculator' },
    { type: 'radio', label: '单选', icon: 'Dot' },
    { type: 'checkbox', label: '多选', icon: 'Finished' },
    { type: 'select', label: '下拉选择', icon: 'ArrowDown' }
  ]},
  { group: '高级组件', items: [
    { type: 'email', label: '邮箱', icon: 'Message' },
    { type: 'phone', label: '手机号', icon: 'Phone' },
    { type: 'date', label: '日期', icon: 'Calendar' },
    { type: 'time', label: '时间', icon: 'Clock' },
    { type: 'upload', label: '文件上传', icon: 'UploadFilled' },
    { type: 'rate', label: '评分', icon: 'Star' },
    { type: 'switch', label: '开关', icon: 'Switch' },
    { type: 'slider', label: '滑块', icon: 'Rank' }
  ]},
  { group: '布局组件', items: [
    { type: 'divider', label: '分割线', icon: 'MiniProgram' },
    { type: 'description', label: '说明文字', icon: 'InfoFilled' }
  ]}
]
