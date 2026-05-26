import api from './index'

export const getOrders = (params) => api.get('/orders', { params })
export const getOrder = (id) => api.get(`/orders/${id}`)
export const createOrder = (data) => api.post('/orders', data)
export const deleteOrder = (id) => api.delete(`/orders/${id}`)
export const updateOrderStatus = (id, action) => api.post(`/orders/${id}/${action}`)
export const assignStaff = (id, staffId) => api.put(`/orders/${id}/staff`, { staff_id: staffId })