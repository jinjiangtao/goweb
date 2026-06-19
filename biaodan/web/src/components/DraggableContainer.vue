<template>
  <div class="draggable-container" ref="containerRef">
    <slot />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    required: true
  },
  itemKey: {
    type: String,
    default: 'id'
  },
  animation: {
    type: Number,
    default: 200
  },
  ghostClass: {
    type: String,
    default: 'ghost'
  }
})

const emit = defineEmits(['update:modelValue', 'change'])
const containerRef = ref(null)
let dragSrcIndex = -1
let draggingEl = null
let items = []

function getChildren() {
  if (!containerRef.value) return []
  return Array.from(containerRef.value.querySelectorAll(':scope > .field-wrapper'))
}

function getItemIndex(el) {
  const children = getChildren()
  return children.indexOf(el)
}

function onDragStart(e) {
  const target = e.target.closest('.field-wrapper')
  if (!target) return
  dragSrcIndex = getItemIndex(target)
  draggingEl = target
  target.style.opacity = '0.4'
  e.dataTransfer.effectAllowed = 'move'
  try {
    e.dataTransfer.setData('text/plain', String(dragSrcIndex))
  } catch (err) {}
}

function onDragEnd(e) {
  if (draggingEl) {
    draggingEl.style.opacity = ''
  }
  draggingEl = null
  dragSrcIndex = -1
  const children = getChildren()
  children.forEach(el => {
    el.classList.remove(props.ghostClass)
    el.style.transform = ''
    el.style.transition = ''
  })
}

function onDragOver(e) {
  e.preventDefault()
  if (dragSrcIndex < 0 || !draggingEl) return
  e.dataTransfer.dropEffect = 'move'

  const target = e.target.closest('.field-wrapper')
  if (!target || target === draggingEl) return

  const targetIndex = getItemIndex(target)
  if (targetIndex < 0) return

  const children = getChildren()
  const srcEl = children[dragSrcIndex]
  const rect = target.getBoundingClientRect()
  const insertBefore = e.clientY < rect.top + rect.height / 2

  if (insertBefore) {
    target.parentNode.insertBefore(srcEl, target)
  } else {
    target.parentNode.insertBefore(srcEl, target.nextSibling)
  }

  const newIndex = getItemIndex(srcEl)
  if (newIndex !== dragSrcIndex) {
    const newArr = [...props.modelValue]
    const [item] = newArr.splice(dragSrcIndex, 1)
    newArr.splice(newIndex, 0, item)
    emit('update:modelValue', newArr)
    emit('change', newArr)
    dragSrcIndex = newIndex
  }
}

function onDrop(e) {
  e.preventDefault()
  onDragEnd()
}

function bindEvents() {
  if (!containerRef.value) return
  const container = containerRef.value
  container.addEventListener('dragstart', onDragStart)
  container.addEventListener('dragend', onDragEnd)
  container.addEventListener('dragover', onDragOver)
  container.addEventListener('drop', onDrop)
}

function unbindEvents() {
  if (!containerRef.value) return
  const container = containerRef.value
  container.removeEventListener('dragstart', onDragStart)
  container.removeEventListener('dragend', onDragEnd)
  container.removeEventListener('dragover', onDragOver)
  container.removeEventListener('drop', onDrop)
}

watch(() => props.modelValue, () => {
  nextTick(() => {
    const children = getChildren()
    children.forEach(el => {
      el.setAttribute('draggable', 'true')
    })
  })
}, { deep: true, immediate: true })

onMounted(() => {
  nextTick(bindEvents)
})

onBeforeUnmount(() => {
  unbindEvents()
})
</script>

<style scoped>
.draggable-container {
  width: 100%;
}
</style>
