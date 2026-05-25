import api from './index'

export const getOrderStats = () => api.get('/stats/orders')
export const getRevenueStats = () => api.get('/stats/revenue')