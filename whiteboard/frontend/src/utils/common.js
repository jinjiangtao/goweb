export function generateId() {
  return Date.now().toString(36) + Math.random().toString(36).substr(2, 9)
}

export function generateColor() {
  const colors = [
    '#f56c6c', '#e6a23c', '#67c23a', '#409eff', 
    '#909399', '#8e44ad', '#16a085', '#d35400',
    '#c0392b', '#2980b9', '#27ae60', '#f39c12'
  ]
  return colors[Math.floor(Math.random() * colors.length)]
}

export function formatTime(timestamp) {
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

export function debounce(fn, delay) {
  let timer = null
  return function () {
    const args = arguments
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => fn.apply(this, args), delay)
  }
}

export function throttle(fn, delay) {
  let last = 0
  return function () {
    const now = Date.now()
    if (now - last >= delay) {
      last = now
      fn.apply(this, arguments)
    }
  }
}
