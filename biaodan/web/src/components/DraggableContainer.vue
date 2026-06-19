<template>
  <div class="draggable-container" ref="containerRef">
    <div
      v-for="(item, index) in innerList"
      :key="item[itemKey] || index"
      class="draggable-item"
      :draggable="true"
      @dragstart="onItemDragStart($event, index)"
      @dragover="onItemDragOver($event, index)"
      @drop="onItemDrop($event, index)"
      @dragend="onDragEnd"
      @dragenter.prevent
    >
      <slot name="item" :element="item" :index="index" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'

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

const innerList = ref([...props.modelValue])
let dragSrcIndex = -1
let lastOverIndex = -1

watch(
  () => props.modelValue,
  (val) => {
    innerList.value = Array.isArray(val) ? [...val] : []
  },
  { deep: true }
)

function onItemDragStart(e, index) {
  dragSrcIndex = index
  lastOverIndex = index
  e.dataTransfer.effectAllowed = 'move'
  try {
    e.dataTransfer.setData('text/plain', String(index))
  } catch (err) {}
  requestAnimationFrame(() => {
    const el = e.currentTarget
    if (el) {
      el.classList.add(props.ghostClass)
      el.style.opacity = '0.4'
    }
  })
}

function onItemDragOver(e, index) {
  e.preventDefault()
  if (dragSrcIndex < 0) return
  e.dataTransfer.dropEffect = 'move'
  if (index === dragSrcIndex || index === lastOverIndex) return
  lastOverIndex = index

  const list = innerList.value
  const srcItem = list[dragSrcIndex]
  list.splice(dragSrcIndex, 1)
  list.splice(index, 0, srcItem)
  dragSrcIndex = index
}

function onItemDrop(e, index) {
  e.preventDefault()
  if (dragSrcIndex < 0) return
  e.stopPropagation()
  commitChange()
}

function onDragEnd() {
  if (dragSrcIndex >= 0) {
    commitChange()
  }
}

function commitChange() {
  const newArr = [...innerList.value]
  dragSrcIndex = -1
  lastOverIndex = -1
  emit('update:modelValue', newArr)
  emit('change', newArr)
  setTimeout(() => {
    document.querySelectorAll('.' + props.ghostClass).forEach((el) => {
      el.classList.remove(props.ghostClass)
      el.style.opacity = ''
    })
  }, 0)
}

onMounted(() => {})
onBeforeUnmount(() => {})
</script>

<style scoped>
.draggable-container {
  width: 100%;
  display: flex;
  flex-direction: column;
}

.draggable-item {
  transition: transform 0.2s ease;
}

.ghost {
  opacity: 0.5;
}
</style>
