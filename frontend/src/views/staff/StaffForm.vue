<template>
  <el-dialog
    v-model="visible"
    :title="staff ? '编辑人员' : '添加人员'"
    width="500px"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="姓名" prop="name">
        <el-input v-model="form.name" placeholder="请输入姓名" />
      </el-form-item>
      <el-form-item label="手机号" prop="phone">
        <el-input v-model="form.phone" placeholder="请输入手机号" />
      </el-form-item>
      <el-form-item label="头像URL" prop="avatar">
        <el-input v-model="form.avatar" placeholder="请输入头像URL" />
      </el-form-item>
      <el-form-item v-if="staff" label="状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择状态">
          <el-option label="空闲" value="available" />
          <el-option label="忙碌" value="busy" />
          <el-option label="休息" value="off" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { createStaff, updateStaff } from '../../api/staff'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
  staff: Object,
})
const emit = defineEmits(['update:modelValue', 'success'])

const formRef = ref(null)
const loading = ref(false)

const form = ref({
  name: '',
  phone: '',
  avatar: '',
  status: 'available',
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
}

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

watch(() => props.staff, (newStaff) => {
  if (newStaff) {
    form.value = {
      name: newStaff.name || '',
      phone: newStaff.phone || '',
      avatar: newStaff.avatar || '',
      status: newStaff.status || 'available',
    }
  } else {
    form.value = { name: '', phone: '', avatar: '', status: 'available' }
  }
}, { immediate: true })

const handleClose = () => {
  visible.value = false
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    if (props.staff) {
      await updateStaff(props.staff.id, form.value)
      ElMessage.success('更新成功')
    } else {
      await createStaff(form.value)
      ElMessage.success('创建成功')
    }
    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}
</script>