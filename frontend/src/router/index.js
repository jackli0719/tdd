import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../components/Layout.vue'
import { setupAuthGuard } from './guards'

const routes = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/Dashboard.vue'),
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/user/UserList.vue'),
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('../views/product/ProductList.vue'),
      },
      {
        path: 'categories',
        name: 'Categories',
        component: () => import('../views/category/CategoryList.vue'),
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('../views/order/OrderList.vue'),
      },
      {
        path: 'staff',
        name: 'Staff',
        component: () => import('../views/staff/StaffList.vue'),
      },
      {
        path: 'stats',
        name: 'Stats',
        component: () => import('../views/stats/Stats.vue'),
      },
    ],
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue'),
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/auth/Register.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

setupAuthGuard(router)

export default router