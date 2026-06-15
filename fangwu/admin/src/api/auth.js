import request from './request'

export function login(data) {
  return request.post('/api/login', data)
}

export function getProfile() {
  return request.get('/api/profile')
}

export function changePassword(data) {
  return request.put('/api/change-password', data)
}
