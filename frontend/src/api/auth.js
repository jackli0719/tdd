import api from './index'
import { getToken, setToken, removeToken } from './index'

export const login = (data) => {
  return api.post('/auth/login', data)
}

export const register = (data) => {
  return api.post('/auth/register', data)
}

export const getCurrentUser = () => {
  return api.get('/auth/me')
}

export { getToken, setToken, removeToken }