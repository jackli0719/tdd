<template>
  <el-dialog
    v-model="visible"
    :title="product ? 'Edit Product' : 'Add Product'"
    width="500px"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="Name">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="Price">
        <el-input v-model="form.price" type="number" />
      </el-form-item>
      <el-form-item label="Stock">
        <el-input v-model="form.stock" type="number" />
      </el-form-item>
      <el-form-item label="Description">
        <el-input v-model="form.description" type="textarea" rows="3" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">Cancel</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">Submit</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { createProduct, updateProduct } from '../../api/product'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
  product: Object,
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(props.modelValue)
const loading = ref(false)
const form = ref({
  name: '',
  price: '',
  stock: '',
  description: '',
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.product) {
    form.value = {
      name: props.product.name || '',
      price: props.product.price || '',
      stock: props.product.stock || '',
      description: props.product.description || '',
    }
  } else if (val) {
    form.value = { name: '', price: '', stock: '', description: '' }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  if (!form.value.name || !form.value.price) {
    ElMessage.warning('Please fill in required fields')
    return
  }
  loading.value = true
  try {
    const data = {
      ...form.value,
      price: Number(form.value.price),
      stock: Number(form.value.stock),
    }
    if (props.product) {
      await updateProduct(props.product.id, data)
    } else {
      await createProduct(data)
    }
    ElMessage.success('Success')
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error(error.message || 'Failed to save product')
  } finally {
    loading.value = false
  }
}
</script>