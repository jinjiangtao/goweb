<template>
  <div class="editor-wrapper" style="height: 100%; min-height: 400px;">
    <div ref="editorRef" style="width: 100%; height: 100%; min-height: 400px;"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, shallowRef } from 'vue'
import * as monaco from 'monaco-editor'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'

self.MonacoEnvironment = {
  getWorker(_, label) {
    if (label === 'json') return new jsonWorker()
    if (label === 'css' || label === 'scss' || label === 'less') return new cssWorker()
    if (label === 'html' || label === 'handlebars' || label === 'razor') return new htmlWorker()
    if (label === 'typescript' || label === 'javascript') return new tsWorker()
    return new editorWorker()
  }
}

const props = defineProps({
  modelValue: { type: String, default: '' },
  language: { type: String, default: 'javascript' },
  readOnly: { type: Boolean, default: false },
  diffMode: { type: Boolean, default: false },
  originalValue: { type: String, default: '' },
  theme: { type: String, default: 'vs-dark' }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = ref(null)
const editor = shallowRef(null)
const diffEditor = shallowRef(null)

const langMap = {
  cpp: 'cpp', c: 'cpp', csharp: 'csharp', dart: 'dart', go: 'go', java: 'java',
  javascript: 'javascript', js: 'javascript', typescript: 'typescript', ts: 'typescript',
  markdown: 'markdown', md: 'markdown', php: 'php', python: 'python', py: 'python',
  ruby: 'ruby', rust: 'rust', shell: 'shell', bash: 'shell', sql: 'sql',
  yaml: 'yaml', css: 'css', scss: 'scss', html: 'html', json: 'json'
}

const getLang = (lang) => langMap[lang?.toLowerCase()] || 'plaintext'

onMounted(() => {
  monaco.editor.defineTheme('custom-dark', {
    base: 'vs-dark',
    inherit: true,
    rules: [],
    colors: {
      'editor.background': '#1e1e1e'
    }
  })

  if (props.diffMode) {
    diffEditor.value = monaco.editor.createDiffEditor(editorRef.value, {
      readOnly: props.readOnly,
      theme: props.theme,
      automaticLayout: true,
      minimap: { enabled: false },
      fontSize: 14,
      fontFamily: 'Consolas, Monaco, Courier New, monospace',
      renderLineHighlight: 'all',
      scrollBeyondLastLine: false,
      wordWrap: 'on'
    })
    const originalModel = monaco.editor.createModel(props.originalValue, getLang(props.language))
    const modifiedModel = monaco.editor.createModel(props.modelValue, getLang(props.language))
    diffEditor.value.setModel({ original: originalModel, modified: modifiedModel })
  } else {
    editor.value = monaco.editor.create(editorRef.value, {
      value: props.modelValue,
      language: getLang(props.language),
      readOnly: props.readOnly,
      theme: props.theme,
      automaticLayout: true,
      minimap: { enabled: false },
      fontSize: 14,
      fontFamily: 'Consolas, Monaco, Courier New, monospace',
      renderLineHighlight: 'all',
      scrollBeyondLastLine: false,
      wordWrap: 'on',
      tabSize: 2,
      insertSpaces: true,
      formatOnPaste: true,
      formatOnType: true
    })
    editor.value.onDidChangeModelContent(() => {
      emit('update:modelValue', editor.value.getValue())
    })
  }
})

watch(() => props.language, (newLang) => {
  if (!props.diffMode && editor.value) {
    monaco.editor.setModelLanguage(editor.value.getModel(), getLang(newLang))
  }
})

watch(() => props.modelValue, (newVal) => {
  if (!props.diffMode && editor.value) {
    if (editor.value.getValue() !== newVal) {
      editor.value.setValue(newVal)
    }
  }
})

watch(() => props.originalValue, (newVal) => {
  if (props.diffMode && diffEditor.value) {
    const model = diffEditor.value.getModel()
    model.original.setValue(newVal)
  }
})

const formatCode = () => {
  if (!props.diffMode && editor.value) {
    editor.value.getAction('editor.action.formatDocument').run()
  }
}

const getValue = () => editor.value?.getValue() || ''

defineExpose({ formatCode, getValue, editor })

onBeforeUnmount(() => {
  editor.value?.dispose()
  diffEditor.value?.dispose()
})
</script>
