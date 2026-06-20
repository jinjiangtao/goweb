import request from '../utils/request'

export function login(data) {
  return request.post('/login', data)
}

export function getCurrentUser() {
  return request.get('/user/current')
}

export function getUsers() {
  return request.get('/users')
}
