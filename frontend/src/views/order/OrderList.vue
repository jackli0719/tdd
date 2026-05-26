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
        <el-table-column prop="staff_id" label="服务人员">
          <template #default="{ row }">
            {{ row.staff_id || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="appointment_time" label="预约时间">
          <template #default="{ row }">
            {{ formatAppointmentTime(row.appointment_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="address" label="服务地址">
          <template #default="{ row }">
            {{ row.address || (row.address_id ? '地址#' + row.address_id : '-') }}
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
              @click="updateStatus(row.id, 'paid')"
            >支付</el-button>
            <el-button
              v-if="row.status === 'paid'"
              size="small"
              type="warning"
              @click="updateStatus(row.id, 'ship')"
            >发货</el-button>
            <el-button
              v-if="row.status === 'shipped'"
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
            <el-button
              v-if="row.status === 'completed'"
              size="small"
              type="success"
              @click="showReviewForm(row)"
            >评价</el-button>
            <el-button
              size="small"
              @click="showAssignStaff(row)"
            >分配</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <OrderForm
      v-model="formVisible"
      @success="loadOrders"
    />
    <ReviewForm
      v-model="reviewDialogVisible"
      :order="reviewOrder"
      @success="loadOrders"
    />

    <!-- Assign Staff Dialog -->
    <el-dialog v-model="assignDialogVisible" title="分配服务人员" width="500px">
      <el-table :data="availableStaff" v-loading="staffLoading" @row-click="selectStaff">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === 'available' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="assignDialogVisible = false">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOrders, updateOrderStatus, assignStaff } from '../../api/order'
import { getAvailableStaff } from '../../api/staff'
import OrderForm from './OrderForm.vue'
import ReviewForm from '../review/ReviewForm.vue'
import { ElMessage } from 'element-plus'

const orders = ref([])
const loading = ref(false)
const formVisible = ref(false)

// Assign staff dialog
const assignDialogVisible = ref(false)
const availableStaff = ref([])
const staffLoading = ref(false)
const currentOrder = ref(null)
const reviewDialogVisible = ref(false)
const reviewOrder = ref(null)

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
    paid: 'success',
    shipped: 'warning',
    completed: '',
    cancelled: 'danger',
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待支付',
    paid: '已支付',
    shipped: '已发货',
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

const showAssignStaff = async (order) => {
  currentOrder.value = order
  assignDialogVisible.value = true
  staffLoading.value = true
  try {
    const res = await getAvailableStaff()
    availableStaff.value = res.data.staffs || []
  } catch (error) {
    ElMessage.error('加载服务人员列表失败')
  } finally {
    staffLoading.value = false
  }
}

const showReviewForm = (order) => {
  reviewOrder.value = order
  reviewDialogVisible.value = true
}

const selectStaff = async (row) => {
  if (!currentOrder.value) return
  try {
    await assignStaff(currentOrder.value.id, row.id)
    ElMessage.success('分配成功')
    assignDialogVisible.value = false
    loadOrders()
  } catch (error) {
    ElMessage.error(error.message || '分配失败')
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

const formatAppointmentTime = (date) => {
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
