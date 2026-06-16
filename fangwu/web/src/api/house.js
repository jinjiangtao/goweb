import request from './request'

export function getHouseList(params) {
  return request.get('/api/public/houses', { params })
}

export function getHouseDetail(id) {
  return request.get(`/api/public/houses/${id}`)
}

export function recordConsult(id) {
  return request.post(`/api/public/houses/${id}/consult`)
}

export function getRecommendHouses(id) {
  return request.get(`/api/public/houses/${id}/recommend`)
}

export function incrementViewCount(id) {
  return request.post(`/api/public/houses/${id}/view`)
}
