import api from './index'

export const getStaffList = (params) => api.get('/staff', { params })
export const getStaff = (id) => api.get(`/staff/${id}`)
export const createStaff = (data) => api.post('/staff', data)
export const updateStaff = (id, data) => api.put(`/staff/${id}`, data)
export const deleteStaff = (id) => api.delete(`/staff/${id}`)
export const updateStaffStatus = (id, status) => api.put(`/staff/${id}/status`, { status })
export const getAvailableStaff = () => api.get('/staff', { params: { status: 'available' } })