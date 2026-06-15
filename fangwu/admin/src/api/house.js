import request from './request'

export function getHouses(params) {
  return request.get('/api/houses', { params })
}

export function getHouse(id) {
  return request.get(`/api/houses/${id}`)
}

export function createHouse(data) {
  return request.post('/api/houses', data)
}

export function updateHouse(id, data) {
  return request.put(`/api/houses/${id}`, data)
}

export function deleteHouse(id) {
  return request.delete(`/api/houses/${id}`)
}

export function toggleHouseStatus(id) {
  return request.put(`/api/houses/${id}/status`)
}

export function toggleHouseRecommend(id) {
  return request.put(`/api/houses/${id}/recommend`)
}

export function uploadImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/api/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
