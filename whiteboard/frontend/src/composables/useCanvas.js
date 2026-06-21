import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useBoardStore } from '@/store/board'
import { generateId, throttle } from '@/utils/common'
import { useWebSocket } from './useWebSocket'

export function useCanvas(canvasRef) {
  const store = useBoardStore()
  const ws = useWebSocket()
  
  const ctx = ref(null)
  const tempCanvas = ref(null)
  const tempCtx = ref(null)
  const isDrawing = ref(false)
  const startPoint = ref(null)
  const currentPoints = ref([])
  const textInputPosition = ref(null)
  const textInputVisible = ref(false)
  const textInputValue = ref('')
  
  let canvasWidth = 0
  let canvasHeight = 0

  function initCanvas() {
    if (!canvasRef.value) return
    
    const canvas = canvasRef.value
    const container = canvas.parentElement
    if (!container) return
    
    canvasWidth = container.clientWidth
    canvasHeight = container.clientHeight
    
    canvas.width = canvasWidth * window.devicePixelRatio
    canvas.height = canvasHeight * window.devicePixelRatio
    canvas.style.width = canvasWidth + 'px'
    canvas.style.height = canvasHeight + 'px'
    
    ctx.value = canvas.getContext('2d')
    if (ctx.value) {
      ctx.value.scale(window.devicePixelRatio, window.devicePixelRatio)
    }
    
    tempCanvas.value = document.createElement('canvas')
    tempCanvas.value.width = canvas.width
    tempCanvas.value.height = canvas.height
    tempCtx.value = tempCanvas.value.getContext('2d')
    if (tempCtx.value) {
      tempCtx.value.scale(window.devicePixelRatio, window.devicePixelRatio)
    }
    
    render()
  }

  function getMousePos(e) {
    if (!canvasRef.value) return { x: 0, y: 0 }
    
    const rect = canvasRef.value.getBoundingClientRect()
    let clientX, clientY
    
    if ('touches' in e) {
      clientX = e.touches[0].clientX
      clientY = e.touches[0].clientY
    } else {
      clientX = e.clientX
      clientY = e.clientY
    }
    
    return {
      x: (clientX - rect.left - store.pan.x) / store.zoom,
      y: (clientY - rect.top - store.pan.y) / store.zoom
    }
  }

  function startDrawing(e) {
    e.preventDefault()
    
    const pos = getMousePos(e)
    
    if (store.currentTool === 'text') {
      textInputPosition.value = pos
      textInputValue.value = ''
      textInputVisible.value = true
      return
    }
    
    if (store.currentTool === 'pan') {
      startPoint.value = { x: 'touches' in e ? e.touches[0].clientX : e.clientX, y: 'touches' in e ? e.touches[0].clientY : e.clientY }
      return
    }
    
    isDrawing.value = true
    startPoint.value = pos
    currentPoints.value = [pos]
    
    if (store.currentTool === 'pen' || store.currentTool === 'eraser') {
      store.currentElement = createElement(pos)
    }
    
    ws.sendCursor(pos)
  }

  function draw(e) {
    e.preventDefault()
    
    const pos = getMousePos(e)
    
    throttleCursor(pos)
    
    if (store.currentTool === 'pan' && startPoint.value) {
      const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX
      const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY
      const dx = clientX - startPoint.value.x
      const dy = clientY - startPoint.value.y
      store.setPan(store.pan.x + dx, store.pan.y + dy)
      startPoint.value = { x: clientX, y: clientY }
      render()
      return
    }
    
    if (!isDrawing.value) return
    
    currentPoints.value.push(pos)
    
    if (store.currentTool === 'pen' || store.currentTool === 'eraser') {
      if (store.currentElement) {
        store.currentElement.points = [...currentPoints.value]
      }
      render()
      drawTempElement()
    } else {
      render()
      drawTempShape(pos)
    }
  }

  function stopDrawing(e) {
    if (store.currentTool === 'pan') {
      startPoint.value = null
      return
    }
    
    if (store.currentTool === 'text') {
      return
    }
    
    if (!isDrawing.value) return
    isDrawing.value = false
    
    if (currentPoints.value.length > 0 || startPoint.value) {
      let element
      
      if (store.currentElement) {
        element = store.currentElement
      } else {
        element = createElement(startPoint.value)
        element.points = currentPoints.value.length > 0 ? currentPoints.value : [startPoint.value, startPoint.value]
      }
      
      if (store.currentTool === 'eraser') {
        eraseAtPoints(currentPoints.value)
      } else {
        store.addElement(element)
        ws.sendDraw(element)
      }
    }
    
    store.currentElement = null
    currentPoints.value = []
    startPoint.value = null
    clearTempCanvas()
    render()
  }

  function createElement(pos) {
    return {
      id: generateId(),
      type: store.currentTool,
      points: [pos],
      style: { ...store.currentStyle },
      userId: store.currentUser.id,
      timestamp: Date.now()
    }
  }

  function drawTempElement() {
    if (!tempCtx.value || !store.currentElement) return
    clearTempCanvas()
    drawElement(store.currentElement, tempCtx.value)
  }

  function drawTempShape(endPos) {
    if (!tempCtx.value || !startPoint.value) return
    clearTempCanvas()
    
    const element = {
      id: 'temp',
      type: store.currentTool,
      points: [startPoint.value, endPos],
      style: { ...store.currentStyle },
      userId: store.currentUser.id,
      timestamp: Date.now()
    }
    
    drawElement(element, tempCtx.value)
  }

  function clearTempCanvas() {
    if (tempCtx.value) {
      tempCtx.value.clearRect(0, 0, canvasWidth, canvasHeight)
    }
  }

  function render() {
    if (!ctx.value || !canvasRef.value) return
    
    ctx.value.save()
    ctx.value.fillStyle = store.background
    ctx.value.fillRect(0, 0, canvasWidth, canvasHeight)
    
    ctx.value.translate(store.pan.x, store.pan.y)
    ctx.value.scale(store.zoom, store.zoom)
    
    for (const element of store.elements) {
      drawElement(element, ctx.value)
    }
    
    if (store.currentElement) {
      drawElement(store.currentElement, ctx.value)
    }
    
    ctx.value.restore()
    
    if (tempCanvas.value) {
      ctx.value.drawImage(tempCanvas.value, 0, 0)
    }
  }

  function drawElement(element, context) {
    if (!element.points || element.points.length === 0) return
    
    context.save()
    context.globalAlpha = element.style.opacity
    
    switch (element.type) {
      case 'pen':
        drawPen(element, context)
        break
      case 'line':
        drawLine(element, context)
        break
      case 'rect':
        drawRect(element, context)
        break
      case 'circle':
        drawCircle(element, context)
        break
      case 'ellipse':
        drawEllipse(element, context)
        break
      case 'triangle':
        drawTriangle(element, context)
        break
      case 'arrow':
        drawArrow(element, context)
        break
      case 'text':
        drawText(element, context)
        break
      case 'image':
        drawImageElement(element, context)
        break
    }
    
    context.restore()
  }

  function drawPen(element, context) {
    if (element.points.length < 2) return
    
    context.strokeStyle = element.style.color
    context.lineWidth = element.style.width
    context.lineCap = 'round'
    context.lineJoin = 'round'
    
    context.beginPath()
    context.moveTo(element.points[0].x, element.points[0].y)
    
    for (let i = 1; i < element.points.length; i++) {
      const p1 = element.points[i - 1]
      const p2 = element.points[i]
      const midX = (p1.x + p2.x) / 2
      const midY = (p1.y + p2.y) / 2
      context.quadraticCurveTo(p1.x, p1.y, midX, midY)
    }
    
    context.stroke()
  }

  function drawLine(element, context) {
    if (element.points.length < 2) return
    
    const start = element.points[0]
    const end = element.points[element.points.length - 1]
    
    context.strokeStyle = element.style.color
    context.lineWidth = element.style.width
    context.lineCap = 'round'
    
    context.beginPath()
    context.moveTo(start.x, start.y)
    context.lineTo(end.x, end.y)
    context.stroke()
  }

  function drawRect(element, context) {
    if (element.points.length < 2) return
    
    const start = element.points[0]
    const end = element.points[element.points.length - 1]
    const x = Math.min(start.x, end.x)
    const y = Math.min(start.y, end.y)
    const w = Math.abs(end.x - start.x)
    const h = Math.abs(end.y - start.y)
    
    if (element.style.fill && element.style.fill !== 'transparent') {
      context.globalAlpha = (element.style.fillOpacity ?? 1) * element.style.opacity
      context.fillStyle = element.style.fill
      context.fillRect(x, y, w, h)
      context.globalAlpha = element.style.opacity
    }
    
    context.strokeStyle = element.style.color
    context.lineWidth = element.style.width
    context.strokeRect(x, y, w, h)
  }

  function drawCircle(element, context) {
    if (element.points.length < 2) return
    
    const start = element.points[0]
    const end = element.points[element.points.length - 1]
    const radius = Math.sqrt(Math.pow(end.x - start.x, 2) + Math.pow(end.y - start.y, 2))
    
    context.beginPath()
    context.arc(start.x, start.y, radius, 0, Math.PI * 2)
    
    if (element.style.fill && element.style.fill !== 'transparent') {
      context.globalAlpha = (element.style.fillOpacity ?? 1) * element.style.opacity
      context.fillStyle = element.style.fill
      context.fill()
      context.globalAlpha = element.style.opacity
    }
    
    context.strokeStyle = element.style.color
    context.lineWidth = element.style.width
    context.stroke()
  }

  function drawEllipse(element, context) {
    if (element.points.length < 2) return
    
    const start = element.points[0]
    const end = element.points[element.points.length - 1]
    const rx = Math.abs(end.x - start.x)
    const ry = Math.abs(end.y - start.y)
    const cx = (start.x + end.x) / 2
    const cy = (start.y + end.y) / 2
    
    context.beginPath()
    context.ellipse(cx, cy, rx, ry, 0, 0, Math.PI * 2)
    
    if (element.style.fill && element.style.fill !== 'transparent') {
      context.globalAlpha = (element.style.fillOpacity ?? 1) * element.style.opacity
      context.fillStyle = element.style.fill
      context.fill()
      context.globalAlpha = element.style.opacity
    }
    
    context.strokeStyle = element.style.color
    context.lineWidth = element.style.width
    context.stroke()
  }

  function drawTriangle(element, context) {
    if (element.points.length < 2) return
    
    const start = element.points[0]
    const end = element.points[element.points.length - 1]
    const x = Math.min(start.x, end.x)
    const y = Math.min(start.y, end.y)
    const w = Math.abs(end.x - start.x)
    const h = Math.abs(end.y - start.y)
    
    context.beginPath()
    context.moveTo(x + w / 2, y)
    context.lineTo(x, y + h)
    context.lineTo(x + w, y + h)
    context.closePath()
    
    if (element.style.fill && element.style.fill !== 'transparent') {
      context.globalAlpha = (element.style.fillOpacity ?? 1) * element.style.opacity
      context.fillStyle = element.style.fill
      context.fill()
      context.globalAlpha = element.style.opacity
    }
    
    context.strokeStyle = element.style.color
    context.lineWidth = element.style.width
    context.stroke()
  }

  function drawArrow(element, context) {
    if (element.points.length < 2) return
    
    const start = element.points[0]
    const end = element.points[element.points.length - 1]
    const headLength = element.style.width * 4
    const angle = Math.atan2(end.y - start.y, end.x - start.x)
    
    context.strokeStyle = element.style.color
    context.lineWidth = element.style.width
    context.lineCap = 'round'
    context.lineJoin = 'round'
    
    context.beginPath()
    context.moveTo(start.x, start.y)
    context.lineTo(end.x, end.y)
    context.stroke()
    
    context.beginPath()
    context.moveTo(end.x, end.y)
    context.lineTo(
      end.x - headLength * Math.cos(angle - Math.PI / 6),
      end.y - headLength * Math.sin(angle - Math.PI / 6)
    )
    context.moveTo(end.x, end.y)
    context.lineTo(
      end.x - headLength * Math.cos(angle + Math.PI / 6),
      end.y - headLength * Math.sin(angle + Math.PI / 6)
    )
    context.stroke()
  }

  function drawText(element, context) {
    if (!element.text || element.points.length === 0) return
    
    const pos = element.points[0]
    const fontSize = element.fontSize || 14
    const fontFamily = element.fontFamily || 'sans-serif'
    
    context.font = `${fontSize}px ${fontFamily}`
    context.fillStyle = element.style.color
    context.globalAlpha = element.style.opacity
    context.fillText(element.text, pos.x, pos.y)
  }

  function drawImageElement(element, context) {
    if (!element.imageData) return
    
    const img = new window.Image()
    img.src = element.imageData
    
    const x = element.x ?? element.points[0]?.x ?? 0
    const y = element.y ?? element.points[0]?.y ?? 0
    const w = element.width ?? 200
    const h = element.height ?? 200
    
    img.onload = () => {
      context.drawImage(img, x, y, w, h)
    }
  }

  function eraseAtPoints(points) {
    const eraseRadius = store.currentStyle.width * 2
    
    store.elements = store.elements.filter(elem => {
      if (elem.type === 'eraser') return false
      return !isElementIntersecting(elem, points, eraseRadius)
    })
  }

  function isElementIntersecting(element, points, radius) {
    if (!element.points) return false
    
    for (const p1 of element.points) {
      for (const p2 of points) {
        const dist = Math.sqrt(Math.pow(p1.x - p2.x, 2) + Math.pow(p1.y - p2.y, 2))
        if (dist < radius) return true
      }
    }
    
    if (element.type === 'rect' || element.type === 'triangle' || element.type === 'ellipse' || element.type === 'circle') {
      const start = element.points[0]
      const end = element.points[element.points.length - 1]
      const minX = Math.min(start.x, end.x) - radius
      const maxX = Math.max(start.x, end.x) + radius
      const minY = Math.min(start.y, end.y) - radius
      const maxY = Math.max(start.y, end.y) + radius
      
      for (const p of points) {
        if (p.x >= minX && p.x <= maxX && p.y >= minY && p.y <= maxY) {
          return true
        }
      }
    }
    
    return false
  }

  function addText(text, pos) {
    const element = {
      id: generateId(),
      type: 'text',
      points: [pos],
      style: { ...store.currentStyle },
      userId: store.currentUser.id,
      timestamp: Date.now(),
      text,
      fontSize: 14,
      fontFamily: 'sans-serif'
    }
    
    store.addElement(element)
    ws.sendDraw(element)
    textInputVisible.value = false
    render()
  }

  function addImage(imageData, pos, width, height) {
    const element = {
      id: generateId(),
      type: 'image',
      points: [pos],
      style: { ...store.currentStyle },
      userId: store.currentUser.id,
      timestamp: Date.now(),
      imageData,
      x: pos.x,
      y: pos.y,
      width,
      height
    }
    
    store.addElement(element)
    ws.sendDraw(element)
    render()
  }

  const throttleCursor = throttle((pos) => {
    ws.sendCursor(pos)
  }, 50)

  function handleResize() {
    initCanvas()
  }

  function handleWheel(e) {
    if (e.ctrlKey || e.metaKey) {
      e.preventDefault()
      const delta = e.deltaY > 0 ? 0.9 : 1.1
      store.setZoom(store.zoom * delta)
      render()
    }
  }

  watch(() => store.elements, () => {
    nextTick(() => render())
  }, { deep: true })

  watch(() => store.background, () => {
    render()
  })

  watch(() => [store.zoom, store.pan], () => {
    render()
  }, { deep: true })

  onMounted(() => {
    initCanvas()
    window.addEventListener('resize', handleResize)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
  })

  return {
    initCanvas,
    render,
    startDrawing,
    draw,
    stopDrawing,
    addText,
    addImage,
    isDrawing,
    textInputPosition,
    textInputVisible,
    textInputValue,
    handleWheel
  }
}
