import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../components/Layout.vue'

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
        path: 'stats',
        name: 'Stats',
        component: () => import('../views/stats/Stats.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router