import request from '../utils/request'

export function getWorkRecords(params) {
  return request.get('/work-records', { params })
}

export function createWorkRecord(data) {
  return request.post('/work-records', data)
}

export function batchCreateWorkRecords(data) {
  return request.post('/work-records/batch', data)
}

export function updateWorkRecord(id, data) {
  return request.put(`/work-records/${id}`, data)
}

export function deleteWorkRecord(id) {
  return request.delete(`/work-records/${id}`)
}

export function approveWorkRecord(id) {
  return request.post(`/work-records/${id}/approve`)
}

export function rejectWorkRecord(id, data) {
  return request.post(`/work-records/${id}/reject`, data)
}

export function batchApprove(data) {
  return request.post('/work-records/batch-approve', data)
}

export function getDailyStats(params) {
  return request.get('/stats/daily', { params })
}

export function getProjectStats(params) {
  return request.get('/stats/project', { params })
}

export function getOverallStats(params) {
  return request.get('/stats/overall', { params })
}
