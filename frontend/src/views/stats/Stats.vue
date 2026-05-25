<template>
  <div class="stats-page">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>订单统计</template>
          <div class="stats-content" v-loading="loading">
            <p>订单总数: {{ stats.total }}</p>
            <p>待支付: {{ stats.pending }}</p>
            <p>已支付: {{ stats.paid }}</p>
            <p>已发货: {{ stats.shipped }}</p>
            <p>已完成: {{ stats.completed }}</p>
            <p>已取消: {{ stats.cancelled }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>营收统计</template>
          <div class="stats-content" v-loading="loading">
            <p>总收入: ¥{{ revenue.total_revenue }}</p>
            <p>待支付营收: ¥{{ revenue.pending_revenue }}</p>
            <p>已支付营收: ¥{{ revenue.paid_revenue }}</p>
            <p>已发货营收: ¥{{ revenue.shipped_revenue }}</p>
            <p>已完成营收: ¥{{ revenue.completed_revenue }}</p>
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
  total_revenue: 0,
  pending_revenue: 0,
  paid_revenue: 0,
  shipped_revenue: 0,
  completed_revenue: 0,
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
    ElMessage.error('加载统计数据失败')
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