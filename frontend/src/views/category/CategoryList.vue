<template>
  <div class="category-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>品类管理</span>
          <el-button type="primary" @click="showForm(null)">添加品类</el-button>
        </div>
      </template>
      <el-table :data="categories" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="created_at" label="创建时间">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="showForm(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <CategoryForm
      v-model="formVisible"
      :category="currentCategory"
      @success="loadCategories"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCategories, deleteCategory } from '../../api/category'
import CategoryForm from './CategoryForm.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const categories = ref([])
const loading = ref(false)
const formVisible = ref(false)
const currentCategory = ref(null)

const loadCategories = async () => {
  loading.value = true
  try {
    const res = await getCategories()
    categories.value = res.data.categories || []
  } catch (error) {
    ElMessage.error('加载品类列表失败')
  } finally {
    loading.value = false
  }
}

const showForm = (category) => {
  currentCategory.value = category
  formVisible.value = true
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除该品类吗？', '提示', {
      type: 'warning',
    })
    await deleteCategory(id)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>