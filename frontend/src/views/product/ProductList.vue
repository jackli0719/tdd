<template>
  <div class="product-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>产品管理</span>
          <el-select
            v-model="categoryId"
            placeholder="选择品类"
            clearable
            style="width: 200px; margin-right: 10px"
            @change="loadProducts"
          >
            <el-option label="全部" value="" />
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
          <el-button type="primary" @click="showForm(null)">添加产品</el-button>
        </div>
      </template>
      <el-table :data="products" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="category_id" label="品类ID" width="100" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="price" label="价格">
          <template #default="{ row }">
            ¥{{ row.price }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" />
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
    <ProductForm
      v-model="formVisible"
      :product="currentProduct"
      :categories="categories"
      @success="loadProducts"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProducts, deleteProduct } from '../../api/product'
import { getCategories } from '../../api/category'
import ProductForm from './ProductForm.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const products = ref([])
const categories = ref([])
const loading = ref(false)
const formVisible = ref(false)
const currentProduct = ref(null)
const categoryId = ref('')

const loadCategories = async () => {
  try {
    const res = await getCategories()
    categories.value = res.data.categories || []
  } catch (error) {
    console.error('加载品类失败:', error)
  }
}

const loadProducts = async () => {
  loading.value = true
  try {
    const params = {}
    if (categoryId.value) {
      params.category_id = categoryId.value
    }
    const res = await getProducts(params)
    products.value = res.data.products || []
  } catch (error) {
    ElMessage.error('加载产品列表失败')
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
    await ElMessageBox.confirm('确定要删除该产品吗？', '提示', {
      type: 'warning',
    })
    await deleteProduct(id)
    ElMessage.success('删除成功')
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

onMounted(() => {
  loadCategories()
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
