import { getToken } from '../api/auth'

const whiteList = ['/login', '/register']

export const setupAuthGuard = (router) => {
  router.beforeEach((to, from, next) => {
    const token = getToken()

    if (whiteList.includes(to.path)) {
      next()
    } else if (!token) {
      next('/login')
    } else {
      next()
    }
  })
}

export const isAuthenticated = () => {
  return !!getToken()
}