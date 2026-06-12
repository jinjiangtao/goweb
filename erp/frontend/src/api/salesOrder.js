import request from '../utils/request.js'

export const getSalesOrders = (params) => {
  return request({
    url: '/sales-orders',
    method: 'get',
    params
  })
}

export const getSalesOrder = (id) => {
  return request({
    url: `/sales-orders/${id}`,
    method: 'get'
  })
}

export const createSalesOrder = (data) => {
  return request({
    url: '/sales-orders',
    method: 'post',
    data
  })
}

export const updateSalesOrder = (id, data) => {
  return request({
    url: `/sales-orders/${id}`,
    method: 'put',
    data
  })
}

export const updateSalesOrderStatus = (id, data) => {
  return request({
    url: `/sales-orders/${id}/status`,
    method: 'put',
    data
  })
}

export const deleteSalesOrder = (id) => {
  return request({
    url: `/sales-orders/${id}`,
    method: 'delete'
  })
}
