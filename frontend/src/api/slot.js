import api from './index'

export const getSlots = (date) => api.get('/slots', { params: { date } })