<template>
  <el-dialog
    :model-value="modelValue"
    title="订单评价"
    width="520px"
    @update:model-value="emit('update:modelValue', $event)"
    @close="resetForm"
  >
    <el-form label-width="90px">
      <el-form-item label="订单ID">
        <span>{{ order?.id || '-' }}</span>
      </el-form-item>
      <el-form-item label="评分" required>
        <el-rate v-model="form.rating" :max="5" />
      </el-form-item>
      <el-form-item label="评价内容">
        <el-input
          v-model="form.comment"
          type="textarea"
          :rows="4"
          maxlength="1000"
          show-word-limit
          placeholder="请输入评价内容"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="submit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createReview } from '../../api/review'

const props = defineProps({
  modelValue: Boolean,
  order: Object,
})

const emit = defineEmits(['update:modelValue', 'success'])

const submitting = ref(false)
const form = reactive({
  rating: 5,
  comment: '',
})

const resetForm = () => {
  form.rating = 5
  form.comment = ''
}

watch(() => props.order?.id, () => {
  resetForm()
})

const submit = async () => {
  if (!props.order?.id) {
    ElMessage.warning('请选择订单')
    return
  }
  if (form.rating < 1 || form.rating > 5) {
    ElMessage.warning('请选择评分')
    return
  }

  submitting.value = true
  try {
    await createReview({
      order_id: props.order.id,
      rating: form.rating,
      comment: form.comment,
    })
    ElMessage.success('评价成功')
    emit('update:modelValue', false)
    emit('success')
  } catch (error) {
    ElMessage.error(error.message || '评价失败')
  } finally {
    submitting.value = false
  }
}
</script>
