<template>
  <div class="product-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Products</span>
          <el-button type="primary" @click="showForm(null)">Add Product</el-button>
        </div>
      </template>
      <el-table :data="products" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="Name" />
        <el-table-column prop="price" label="Price">
          <template #default="{ row }">
            ${{ row.price }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="Stock" />
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
    <ProductForm
      v-model="formVisible"
      :product="currentProduct"
      @success="loadProducts"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProducts, deleteProduct } from '../../api/product'
import ProductForm from './ProductForm.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const products = ref([])
const loading = ref(false)
const formVisible = ref(false)
const currentProduct = ref(null)

const loadProducts = async () => {
  loading.value = true
  try {
    const res = await getProducts()
    products.value = res.data || []
  } catch (error) {
    ElMessage.error('Failed to load products')
  } finally {
    loading.value = false
  }
}

const showForm = (product) => {
  currentProduct.value = product
  formVisible.value = true
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('Are you sure to delete this product?', 'Warning', {
      type: 'warning',
    })
    await deleteProduct(id)
    ElMessage.success('Product deleted')
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to delete product')
    }
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>