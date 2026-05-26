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

export const getUserId = async () => {
  try {
    const res = await getCurrentUser()
    return res?.data?.user_id || res?.user_id || null
  } catch {
    return null
  }
}

export { getToken, setToken, removeToken }