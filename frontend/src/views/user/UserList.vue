<template>
  <div class="user-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Users</span>
          <el-button type="primary" @click="showForm(null)">Add User</el-button>
        </div>
      </template>
      <el-table :data="users" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="Username" />
        <el-table-column prop="email" label="Email" />
        <el-table-column prop="created_at" label="Created At">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="showForm(row)">Edit</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <UserForm
      v-model="formVisible"
      :user="currentUser"
      @success="loadUsers"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUsers, deleteUser } from '../../api/user'
import UserForm from './UserForm.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const users = ref([])
const loading = ref(false)
const formVisible = ref(false)
const currentUser = ref(null)

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await getUsers()
    users.value = res.data || []
  } catch (error) {
    ElMessage.error('Failed to load users')
  } finally {
    loading.value = false
  }
}

const showForm = (user) => {
  currentUser.value = user
  formVisible.value = true
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('Are you sure to delete this user?', 'Warning', {
      type: 'warning',
    })
    await deleteUser(id)
    ElMessage.success('User deleted')
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to delete user')
    }
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>