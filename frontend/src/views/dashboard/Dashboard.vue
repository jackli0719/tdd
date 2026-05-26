<template>
  <div class="dashboard">
    <el-row>
      <el-col :span="24">
        <el-card class="time-card">
          <div class="current-time">{{ currentTime }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">用户总数</div>
          <div class="stat-value">{{ stats.users }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">产品总数</div>
          <div class="stat-value">{{ stats.products }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">订单总数</div>
          <div class="stat-value">{{ stats.orders }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-title">总收入</div>
          <div class="stat-value">¥{{ stats.revenue }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>订单统计</template>
          <div class="order-stats">
            <p>待确认: {{ orderStats.pending }}</p>
            <p>已确认: {{ orderStats.confirmed }}</p>
            <p>服务中: {{ orderStats.in_service }}</p>
            <p>已完成: {{ orderStats.completed }}</p>
            <p>已取消: {{ orderStats.cancelled }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>营收统计</template>
          <div class="order-stats">
            <p>总收入: ¥{{ revenueStats.total_revenue }}</p>
            <p>待确认营收: ¥{{ revenueStats.pending_revenue }}</p>
            <p>已确认营收: ¥{{ revenueStats.confirmed_revenue }}</p>
            <p>服务中营收: ¥{{ revenueStats.in_service_revenue }}</p>
            <p>已完成营收: ¥{{ revenueStats.completed_revenue }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getOrderStats, getRevenueStats } from '../../api/stats'
import { getUsers } from '../../api/user'
import { getProducts } from '../../api/product'
import { getOrders } from '../../api/order'

const currentTime = ref('')
let timer = null

const updateTime = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  const seconds = String(now.getSeconds()).padStart(2, '0')
  currentTime.value = `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`
}

const stats = ref({
  users: 0,
  products: 0,
  orders: 0,
  revenue: 0,
})

const orderStats = ref({
  pending: 0,
  confirmed: 0,
  in_service: 0,
  completed: 0,
  cancelled: 0,
})

const revenueStats = ref({
  total_revenue: 0,
  pending_revenue: 0,
  confirmed_revenue: 0,
  in_service_revenue: 0,
  completed_revenue: 0,
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
    stats.value.users = usersRes.data.users?.length || 0
    stats.value.products = productsRes.data.products?.length || 0
    stats.value.orders = ordersRes.data.orders?.length || 0
    stats.value.revenue = revenueRes.data?.total_revenue || 0
    if (orderStatsRes.data) {
      orderStats.value = orderStatsRes.data
    }
    if (revenueRes.data) {
      revenueStats.value = revenueRes.data
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
  loadStats()
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.time-card {
  text-align: center;
  margin-bottom: 20px;
}

.current-time {
  font-size: 36px;
  font-weight: bold;
  color: #303133;
  padding: 20px 0;
}

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