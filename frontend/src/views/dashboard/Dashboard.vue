<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">Total Users</div>
          <div class="stat-value">{{ stats.users }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">Total Products</div>
          <div class="stat-value">{{ stats.products }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">Total Orders</div>
          <div class="stat-value">{{ stats.orders }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">Revenue</div>
          <div class="stat-value">${{ stats.revenue }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>Order Stats</template>
          <div class="order-stats">
            <p>Pending: {{ orderStats.pending }}</p>
            <p>Paid: {{ orderStats.paid }}</p>
            <p>Shipped: {{ orderStats.shipped }}</p>
            <p>Completed: {{ orderStats.completed }}</p>
            <p>Cancelled: {{ orderStats.cancelled }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>Revenue Stats</template>
          <div class="order-stats">
            <p>Today: ${{ revenueStats.today }}</p>
            <p>This Week: ${{ revenueStats.week }}</p>
            <p>This Month: ${{ revenueStats.month }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOrderStats, getRevenueStats } from '../../api/stats'
import { getUsers } from '../../api/user'
import { getProducts } from '../../api/product'
import { getOrders } from '../../api/order'

const stats = ref({
  users: 0,
  products: 0,
  orders: 0,
  revenue: 0,
})

const orderStats = ref({
  pending: 0,
  paid: 0,
  shipped: 0,
  completed: 0,
  cancelled: 0,
})

const revenueStats = ref({
  today: 0,
  week: 0,
  month: 0,
})

const loadStats = async () => {
  try {
    const [usersRes, productsRes, ordersRes, orderStatsRes, revenueRes] = await Promise.all([
      getUsers(),
      getProducts(),
      getOrders(),
      getOrderStats(),
      getRevenueStats(),
    ])
    stats.value.users = usersRes.data?.length || 0
    stats.value.products = productsRes.data?.length || 0
    stats.value.orders = ordersRes.data?.length || 0
    stats.value.revenue = revenueRes.data?.total || 0
    if (orderStatsRes.data) {
      orderStats.value = orderStatsRes.data
    }
    if (revenueRes.data) {
      revenueStats.value = revenueRes.data
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.stat-card {
  text-align: center;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.order-stats p {
  margin: 8px 0;
  font-size: 14px;
}
</style>