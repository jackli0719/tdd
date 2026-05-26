import api from './index'

export const createReview = (data) => api.post('/reviews', data)
export const getReview = (id) => api.get(`/reviews/${id}`)
export const getReviews = (params) => api.get('/reviews', { params })
export const getReviewByOrder = (orderId) => api.get('/reviews', { params: { order_id: orderId } })
export const getStaffReviewSummary = (staffId) => api.get('/reviews/staff-summary', { params: { staff_id: staffId } })
export const getStaffReviewSummaries = () => api.get('/reviews/staff-summary')
