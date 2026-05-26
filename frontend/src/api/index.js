/* global localStorage, window */
import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// Token management
const TOKEN_KEY = 'auth_token'

export const getToken = () => localStorage.getItem(TOKEN_KEY)

export const setToken = (token) => localStorage.setItem(TOKEN_KEY, token)

export const removeToken = () => localStorage.removeItem(TOKEN_KEY)

// Request interceptor - add auth token
api.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      removeToken()
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    const message = error.response?.data?.message || error.message || 'Request failed'
    return Promise.reject(new Error(message))
  }
)

export default api
