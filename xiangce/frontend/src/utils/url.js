export const extractFilename = (path) => {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  
  const normalizedPath = path.replace(/\\/g, '/')
  const parts = normalizedPath.split('/')
  return parts[parts.length - 1]
}

export const getThumbUrl = (path) => {
  const filename = extractFilename(path)
  if (!filename) return ''
  if (filename.startsWith('http')) return filename
  return `/uploads/thumbnails/${filename}`
}

export const getOriginalUrl = (path) => {
  const filename = extractFilename(path)
  if (!filename) return ''
  if (filename.startsWith('http')) return filename
  return `/uploads/${filename}`
}

export const getCoverUrl = (path) => {
  return getThumbUrl(path)
}
