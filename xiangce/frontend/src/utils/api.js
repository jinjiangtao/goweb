import request from '@/utils/request'

export const login = (data) => {
  return request.post('/auth/login', data)
}

export const register = (data) => {
  return request.post('/auth/register', data)
}

export const getUserInfo = () => {
  return request.get('/user/info')
}

export const getUserStats = () => {
  return request.get('/user/stats')
}

export const getAlbums = () => {
  return request.get('/albums')
}

export const getAlbum = (id) => {
  return request.get(`/albums/${id}`)
}

export const createAlbum = (data) => {
  return request.post('/albums', data)
}

export const updateAlbum = (id, data) => {
  return request.put(`/albums/${id}`, data)
}

export const deleteAlbum = (id) => {
  return request.delete(`/albums/${id}`)
}

export const updateAlbumSort = (albumIds) => {
  return request.post('/albums/sort', { album_ids: albumIds })
}

export const getMediaList = (albumId, params = {}) => {
  return request.get(`/media/album/${albumId}`, { params })
}

export const uploadMedia = (albumId, files, onProgress) => {
  const formData = new FormData()
  formData.append('album_id', albumId)
  files.forEach(file => {
    formData.append('files', file)
  })

  return request.post('/media/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        onProgress(percent)
      }
    }
  })
}

export const deleteMedia = (id) => {
  return request.delete(`/media/${id}`)
}

export const batchDeleteMedia = (mediaIds) => {
  return request.post('/media/batch/delete', { media_ids: mediaIds })
}

export const downloadMedia = (id) => {
  window.open(`/api/media/${id}/download`, '_blank')
}

export const batchDownloadMedia = (mediaIds) => {
  return request.post('/media/batch/download', { media_ids: mediaIds }, {
    responseType: 'blob'
  })
}

export const moveMedia = (mediaIds, targetAlbumId) => {
  return request.post('/media/move', {
    media_ids: mediaIds,
    target_album_id: targetAlbumId
  })
}

export const updateMediaSort = (albumId, mediaIds) => {
  return request.post(`/media/sort/${albumId}`, { media_ids: mediaIds })
}

export const getShareLinks = (albumId = null) => {
  if (albumId) {
    return request.get(`/shares/album/${albumId}`)
  }
  return request.get('/shares')
}

export const createShareLink = (data) => {
  return request.post('/shares', data)
}

export const deleteShareLink = (id) => {
  return request.delete(`/shares/${id}`)
}

export const getVisitRecords = (shareId) => {
  return request.get(`/shares/${shareId}/visits`)
}

export const getShareByToken = (token) => {
  return request.get(`/share/${token}`)
}

export const verifySharePassword = (token, password) => {
  return request.post(`/share/${token}/verify`, { password })
}

export const recordShareVisit = (token) => {
  return request.post(`/share/${token}/visit`)
}
