import request from '../utils/request.js'

export const getSuppliers = (params) => {
  return request({
    url: '/suppliers',
    method: 'get',
    params
  })
}

export const getSupplier = (id) => {
  return request({
    url: `/suppliers/${id}`,
    method: 'get'
  })
}

export const createSupplier = (data) => {
  return request({
    url: '/suppliers',
    method: 'post',
    data
  })
}

export const updateSupplier = (id, data) => {
  return request({
    url: `/suppliers/${id}`,
    method: 'put',
    data
  })
}

export const deleteSupplier = (id) => {
  return request({
    url: `/suppliers/${id}`,
    method: 'delete'
  })
}

