import request from '@/utils/request'

export function submitApproval(data) {
  return request({
    url: '/approval',
    method: 'post',
    data
  })
}

export function getApprovalDetail(id) {
  return request({
    url: `/approval/${id}`,
    method: 'get'
  })
}

export function getMyRequests(params) {
  return request({
    url: '/approval/my',
    method: 'get',
    params
  })
}

export function getMyApprovals(params) {
  return request({
    url: '/approval/todo',
    method: 'get',
    params
  })
}

export function getApprovalRecords(id) {
  return request({
    url: `/approval/${id}/records`,
    method: 'get'
  })
}

export function getApprovalProgress(id) {
  return request({
    url: `/approval/${id}/progress`,
    method: 'get'
  })
}

export function approveApproval(id, data) {
  return request({
    url: `/approval/${id}/approve`,
    method: 'post',
    data
  })
}

export function rejectApproval(id, data) {
  return request({
    url: `/approval/${id}/reject`,
    method: 'post',
    data
  })
}

export function revokeApproval(id, data) {
  return request({
    url: `/approval/${id}/revoke`,
    method: 'post',
    data
  })
}

export function addSignApproval(id, data) {
  return request({
    url: `/approval/${id}/addsign`,
    method: 'post',
    data
  })
}

export function transferApproval(id, data) {
  return request({
    url: `/approval/${id}/transfer`,
    method: 'post',
    data
  })
}

export function getApprovalStats() {
  return request({
    url: '/approval/stats',
    method: 'get'
  })
}
