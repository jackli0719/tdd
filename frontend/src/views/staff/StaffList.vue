<template>
  <div class="staff-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>服务人员管理</span>
          <el-button type="primary" @click="showForm(null)">添加人员</el-button>
        </div>
      </template>
      <el-table :data="staffs" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="avatar" label="头像">
          <template #default="{ row }">
            <el-avatar v-if="row.avatar" :src="row.avatar" size="small" />
            <span v-else>-</span>
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
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="showForm(row)">编辑</el-button>
            <el-button
              v-if="row.status !== 'busy'"
              size="small"
              type="warning"
              @click="handleStatusChange(row.id, 'busy')"
            >设为忙碌</el-button>
            <el-button
              v-if="row.status !== 'available'"
              size="small"
              type="success"
              @click="handleStatusChange(row.id, 'available')"
            >设为空闲</el-button>
            <el-button
              v-if="row.status !== 'off'"
              size="small"
              type="info"
              @click="handleStatusChange(row.id, 'off')"
            >休息</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <StaffForm
      v-model="formVisible"
      :staff="currentStaff"
      @success="loadStaffs"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getStaffList, deleteStaff, updateStaffStatus } from '../../api/staff'
import StaffForm from './StaffForm.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const staffs = ref([])
const loading = ref(false)
const formVisible = ref(false)
const currentStaff = ref(null)

const loadStaffs = async () => {
  loading.value = true
  try {
    const res = await getStaffList()
    staffs.value = res.data?.staffs || []
  } catch (error) {
    ElMessage.error('加载人员列表失败')
  } finally {
    loading.value = false
  }
}

const showForm = (staff) => {
  currentStaff.value = staff
  formVisible.value = true
}

const handleStatusChange = async (id, status) => {
  try {
    await updateStaffStatus(id, status)
    ElMessage.success('状态更新成功')
    loadStaffs()
  } catch (error) {
    ElMessage.error(error.message || '状态更新失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除该人员吗？', '提示', {
      type: 'warning',
    })
    await deleteStaff(id)
    ElMessage.success('删除成功')
    loadStaffs()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const getStatusType = (status) => {
  const types = {
    available: 'success',
    busy: 'warning',
    off: 'info',
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    available: '空闲',
    busy: '忙碌',
    off: '休息',
  }
  return texts[status] || status
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadStaffs()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>