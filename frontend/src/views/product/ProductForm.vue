<template>
  <el-dialog
    v-model="visible"
    :title="product ? '编辑产品' : '添加产品'"
    width="500px"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="品类" required>
        <el-select v-model="form.category_id" placeholder="请选择品类" style="width: 100%">
          <el-option
            v-for="cat in categories"
            :key="cat.id"
            :label="cat.name"
            :value="cat.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="名称" required>
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="价格" required>
        <el-input v-model="form.price" type="number" />
      </el-form-item>
      <el-form-item label="库存">
        <el-input v-model="form.stock" type="number" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">提交</el-button>
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
  categories: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(props.modelValue)
const loading = ref(false)
const form = ref({
  category_id: null,
  name: '',
  price: '',
  stock: '',
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.product) {
    form.value = {
      category_id: props.product.category_id || null,
      name: props.product.name || '',
      price: props.product.price || '',
      stock: props.product.stock || '',
    }
  } else if (val) {
    form.value = { category_id: null, name: '', price: '', stock: '' }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  if (!form.value.category_id || !form.value.name || !form.value.price) {
    ElMessage.warning('请填写必填字段')
    return
  }
  loading.value = true
  try {
    const data = {
      category_id: Number(form.value.category_id),
      name: form.value.name,
      price: Number(form.value.price),
      stock: Number(form.value.stock) || 0,
    }
    if (props.product) {
      await updateProduct(props.product.id, data)
    } else {
      await createProduct(data)
    }
    ElMessage.success('保存成功')
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    loading.value = false
  }
}
</script>