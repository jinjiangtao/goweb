import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export const getOrders = (params) => api.get('/orders', { params })
export const getOrder = (id) => api.get(`/orders/${id}`)
export const createOrder = (data) => api.post('/orders', data)
export const updateOrder = (id, data) => api.put(`/orders/${id}`, data)
export const deleteOrder = (id) => api.delete(`/orders/${id}`)
export const searchOrder = (orderNo) => api.get(`/orders/search?order_no=${orderNo}`)

export const getNodes = (orderId) => api.get(`/orders/${orderId}/nodes`)
export const addNode = (orderId, data) => api.post(`/orders/${orderId}/nodes`, data)
export const updateNode = (nodeId, data) => api.put(`/nodes/${nodeId}`, data)
export const deleteNode = (nodeId) => api.delete(`/nodes/${nodeId}`)
export const getProgress = (orderId) => api.get(`/orders/${orderId}/progress`)

export default api
