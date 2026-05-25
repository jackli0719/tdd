<template>
  <div class="stats-page">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>Order Statistics</template>
          <div class="stats-content" v-loading="loading">
            <p>Total Orders: {{ stats.total }}</p>
            <p>Pending: {{ stats.pending }}</p>
            <p>Paid: {{ stats.paid }}</p>
            <p>Shipped: {{ stats.shipped }}</p>
            <p>Completed: {{ stats.completed }}</p>
            <p>Cancelled: {{ stats.cancelled }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>Revenue Statistics</template>
          <div class="stats-content" v-loading="loading">
            <p>Total Revenue: ${{ revenue.total }}</p>
            <p>Today: ${{ revenue.today }}</p>
            <p>This Week: ${{ revenue.week }}</p>
            <p>This Month: ${{ revenue.month }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOrderStats, getRevenueStats } from '../../api/stats'
import { ElMessage } from 'element-plus'

const stats = ref({
  total: 0,
  pending: 0,
  paid: 0,
  shipped: 0,
  completed: 0,
  cancelled: 0,
})

const revenue = ref({
  total: 0,
  today: 0,
  week: 0,
  month: 0,
})

const loading = ref(false)

const loadStats = async () => {
  loading.value = true
  try {
    const [orderRes, revenueRes] = await Promise.all([
      getOrderStats(),
      getRevenueStats(),
    ])
    if (orderRes.data) {
      stats.value = { ...stats.value, ...orderRes.data }
    }
    if (revenueRes.data) {
      revenue.value = { ...revenue.value, ...revenueRes.data }
    }
  } catch (error) {
    ElMessage.error('Failed to load stats')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.stats-content p {
  margin: 12px 0;
  font-size: 16px;
}
</style>