import axios from 'axios'

const instance = axios.create({
  baseURL: '/api/admin',
  timeout: 5000
})

instance.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(response => {
  return response.data
}, error => {
  if (error.response?.status === 401) {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('menus')
    window.location.href = '/login'
  }
  return Promise.reject(error)
})

export const login = (data) => instance.post('/login', data)
export const getInfo = () => instance.get('/info')
export const logout = () => instance.post('/logout')

export const getUsers = (params) => instance.get('/users', { params })
export const createUser = (data) => instance.post('/users', data)
export const updateUser = (id, data) => instance.put(`/users/${id}`, data)
export const deleteUser = (id) => instance.delete(`/users/${id}`)
export const updateUserStatus = (id) => instance.put(`/users/${id}/status`)
export const resetPassword = (id) => instance.put(`/users/${id}/password`)

export const getMenus = () => instance.get('/menus')
export const createMenu = (data) => instance.post('/menus', data)
export const updateMenu = (id, data) => instance.put(`/menus/${id}`, data)
export const deleteMenu = (id) => instance.delete(`/menus/${id}`)

export const getRoles = () => instance.get('/roles')
export const getRoleMenus = (role) => instance.get(`/roles/${role}/menus`)
export const setRoleMenus = (role, data) => instance.put(`/roles/${role}/menus`, data)