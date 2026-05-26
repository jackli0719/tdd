<template>
  <div class="order-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>订单管理</span>
          <el-button type="primary" @click="showForm">创建订单</el-button>
        </div>
      </template>
      <el-table :data="orders" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" />
        <el-table-column prop="total_amount" label="总价">
          <template #default="{ row }">
            ¥{{ row.total_amount }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              size="small"
              type="success"
              @click="updateStatus(row.id, 'confirm')"
            >确认</el-button>
            <el-button
              v-if="row.status === 'confirmed'"
              size="small"
              type="warning"
              @click="updateStatus(row.id, 'start')"
            >开始服务</el-button>
            <el-button
              v-if="row.status === 'in_service'"
              size="small"
              type="primary"
              @click="updateStatus(row.id, 'complete')"
            >完成</el-button>
            <el-button
              v-if="row.status !== 'cancelled' && row.status !== 'completed'"
              size="small"
              type="danger"
              @click="updateStatus(row.id, 'cancel')"
            >取消</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <OrderForm
      v-model="formVisible"
      @success="loadOrders"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOrders, updateOrderStatus } from '../../api/order'
import OrderForm from './OrderForm.vue'
import { ElMessage } from 'element-plus'

const orders = ref([])
const loading = ref(false)
const formVisible = ref(false)

const showForm = () => {
  formVisible.value = true
}

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getOrders()
    orders.value = res.data.orders || []
  } catch (error) {
    ElMessage.error('加载订单列表失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const types = {
    pending: 'info',
    confirmed: 'success',
    in_service: 'warning',
    completed: '',
    cancelled: 'danger',
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待确认',
    confirmed: '已确认',
    in_service: '服务中',
    completed: '已完成',
    cancelled: '已取消',
  }
  return texts[status] || status
}

const updateStatus = async (id, action) => {
  try {
    await updateOrderStatus(id, action)
    ElMessage.success('状态更新成功')
    loadOrders()
  } catch (error) {
    ElMessage.error(error.message || '状态更新失败')
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