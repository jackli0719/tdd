<template>
  <el-dialog
    v-model="visible"
    :title="category ? '编辑品类' : '添加品类'"
    width="500px"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="名称" required>
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="3" />
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
import { createCategory, updateCategory } from '../../api/category'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
  category: Object,
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(props.modelValue)
const loading = ref(false)
const form = ref({
  name: '',
  description: '',
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.category) {
    form.value = {
      name: props.category.name || '',
      description: props.category.description || '',
    }
  } else if (val) {
    form.value = { name: '', description: '' }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  if (!form.value.name) {
    ElMessage.warning('请填写名称')
    return
  }
  loading.value = true
  try {
    if (props.category) {
      await updateCategory(props.category.id, form.value)
    } else {
      await createCategory(form.value)
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
