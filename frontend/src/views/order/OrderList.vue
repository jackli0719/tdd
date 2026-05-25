<template>
  <div class="order-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Orders</span>
        </div>
      </template>
      <el-table :data="orders" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="User ID" />
        <el-table-column prop="total_price" label="Total">
          <template #default="{ row }">
            ${{ row.total_price }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="Status">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created At">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="280">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              size="small"
              type="success"
              @click="updateStatus(row.id, 'paid')"
            >Paid</el-button>
            <el-button
              v-if="row.status === 'paid'"
              size="small"
              type="warning"
              @click="updateStatus(row.id, 'ship')"
            >Ship</el-button>
            <el-button
              v-if="row.status === 'shipped'"
              size="small"
              type="primary"
              @click="updateStatus(row.id, 'complete')"
            >Complete</el-button>
            <el-button
              v-if="row.status !== 'cancelled' && row.status !== 'completed'"
              size="small"
              type="danger"
              @click="updateStatus(row.id, 'cancel')"
            >Cancel</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOrders, updateOrderStatus } from '../../api/order'
import { ElMessage } from 'element-plus'

const orders = ref([])
const loading = ref(false)

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getOrders()
    orders.value = res.data || []
  } catch (error) {
    ElMessage.error('Failed to load orders')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const types = {
    pending: 'info',
    paid: 'success',
    shipped: 'warning',
    completed: '',
    cancelled: 'danger',
  }
  return types[status] || 'info'
}

const updateStatus = async (id, action) => {
  try {
    await updateOrderStatus(id, action)
    ElMessage.success('Status updated')
    loadOrders()
  } catch (error) {
    ElMessage.error(error.message || 'Failed to update status')
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>