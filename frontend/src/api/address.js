import api from './index'

export const getAddresses = (userId) => api.get('/addresses', { params: { user_id: userId } })
export const getAddress = (id) => api.get(`/addresses/${id}`)
export const createAddress = (data) => api.post('/addresses', data)
export const updateAddress = (id, data) => api.put(`/addresses/${id}`, data)
export const deleteAddress = (id) => api.delete(`/addresses/${id}`)
export const setDefaultAddress = (id, userId) => api.put(`/addresses/${id}/default`, { user_id: userId })