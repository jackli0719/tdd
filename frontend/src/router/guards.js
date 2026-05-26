import { getToken } from '../api/auth'

const whiteList = ['/login', '/register']

export const setupAuthGuard = (router) => {
  router.beforeEach((to, from) => {
    const token = getToken()
    const whiteList = ['/login', '/register']

    if (whiteList.includes(to.path)) {
      return true
    } else if (!token) {
      return '/login'
    }
    return true
  })
}

export const isAuthenticated = () => {
  return !!getToken()
}