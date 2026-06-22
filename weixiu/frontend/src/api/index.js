import request from '@/utils/request'

export const login = (data) => request.post('/auth/login', data)
export const register = (data) => request.post('/auth/register', data)
export const getCurrentUser = () => request.get('/auth/me')
export const getUsers = (params) => request.get('/users', { params })
export const getTechnicians = (params) => request.get('/technicians', { params })

export const createOrder = (data) => request.post('/orders', data)
export const getOrders = (params) => request.get('/orders', { params })
export const getOrderDetail = (id) => request.get(`/orders/${id}`)
export const assignOrder = (id, data) => request.post(`/orders/${id}/assign`, data)
export const autoAssignOrder = (id) => request.post(`/orders/${id}/auto-assign`)
export const startProcessOrder = (id) => request.post(`/orders/${id}/start`)
export const processOrder = (id, data) => request.post(`/orders/${id}/process`, data)
export const completeOrder = (id, data) => request.post(`/orders/${id}/complete`, data)
export const rejectOrder = (id, data) => request.post(`/orders/${id}/reject`, data)
export const cancelOrder = (id) => request.post(`/orders/${id}/cancel`)
export const evaluateOrder = (id, data) => request.post(`/orders/${id}/evaluate`, data)

export const createDevice = (data) => request.post('/devices', data)
export const getDevices = (params) => request.get('/devices', { params })
export const getDevice = (id) => request.get(`/devices/${id}`)
export const updateDevice = (id, data) => request.put(`/devices/${id}`, data)
export const deleteDevice = (id) => request.delete(`/devices/${id}`)

export const getStatistics = () => request.get('/stats')

export const uploadFile = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
