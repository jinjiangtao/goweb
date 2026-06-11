import request from '../utils/request.js'

export const getRoles = (params) => {
  return request({
    url: '/roles',
    method: 'get',
    params
  })
}

export const getRole = (id) => {
  return request({
    url: `/roles/${id}`,
    method: 'get'
  })
}

export const createRole = (data) => {
  return request({
    url: '/roles',
    method: 'post',
    data
  })
}

export const updateRole = (id, data) => {
  return request({
    url: `/roles/${id}`,
    method: 'put',
    data
  })
}

export const deleteRole = (id) => {
  return request({
    url: `/roles/${id}`,
    method: 'delete'
  })
}

export const assignMenus = (id, menuIds) => {
  return request({
    url: `/roles/${id}/menus`,
    method: 'put',
    data: { menuIds }
  })
}
