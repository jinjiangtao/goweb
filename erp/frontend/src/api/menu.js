import request from '@/utils/request'

export const getMenus = (params) => {
  return request({
    url: '/menus',
    method: 'get',
    params
  })
}

export const getMenu = (id) => {
  return request({
    url: `/menus/${id}`,
    method: 'get'
  })
}

export const createMenu = (data) => {
  return request({
    url: '/menus',
    method: 'post',
    data
  })
}

export const updateMenu = (id, data) => {
  return request({
    url: `/menus/${id}`,
    method: 'put',
    data
  })
}

export const deleteMenu = (id) => {
  return request({
    url: `/menus/${id}`,
    method: 'delete'
  })
}

export const getMenusByRole = () => {
  return request({
    url: '/menus',
    method: 'get'
  })
}
