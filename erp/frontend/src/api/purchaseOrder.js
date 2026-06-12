import request from '../utils/request.js'

export const getPurchaseOrders = (params) => {
  return request({
    url: '/purchase-orders',
    method: 'get',
    params
  })
}

export const getPurchaseOrder = (id) => {
  return request({
    url: `/purchase-orders/${id}`,
    method: 'get'
  })
}

export const createPurchaseOrder = (data) => {
  return request({
    url: '/purchase-orders',
    method: 'post',
    data
  })
}

export const updatePurchaseOrder = (id, data) => {
  return request({
    url: `/purchase-orders/${id}`,
    method: 'put',
    data
  })
}

export const updatePurchaseOrderStatus = (id, data) => {
  return request({
    url: `/purchase-orders/${id}/status`,
    method: 'put',
    data
  })
}

export const deletePurchaseOrder = (id) => {
  return request({
    url: `/purchase-orders/${id}`,
    method: 'delete'
  })
}
