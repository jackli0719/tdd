<template>
  <div class="review-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>评价管理</span>
          <el-input
            v-model="staffIdFilter"
            clearable
            placeholder="服务人员ID"
            class="filter-input"
            @clear="loadReviews"
            @keyup.enter="loadReviews"
          />
        </div>
      </template>
      <el-table :data="reviews" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="order_id" label="订单ID" />
        <el-table-column prop="user_id" label="用户ID" />
        <el-table-column prop="staff_id" label="服务人员ID" />
        <el-table-column prop="rating" label="评分">
          <template #default="{ row }">
            <el-rate :model-value="row.rating" disabled />
          </template>
        </el-table-column>
        <el-table-column prop="comment" label="评价内容" min-width="180" />
        <el-table-column prop="created_at" label="创建时间">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getReviews } from '../../api/review'

const reviews = ref([])
const loading = ref(false)
const staffIdFilter = ref('')

const loadReviews = async () => {
  loading.value = true
  try {
    const params = {}
    if (staffIdFilter.value) {
      params.staff_id = staffIdFilter.value
    }
    const res = await getReviews(params)
    reviews.value = res.data?.reviews || []
  } catch (error) {
    ElMessage.error(error.message || '加载评价列表失败')
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadReviews()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.filter-input {
  width: 180px;
}
</style>
