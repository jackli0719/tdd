<template>
  <el-dialog
    v-model="visible"
    :title="user ? '编辑用户' : '添加用户'"
    width="500px"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="用户名">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="form.email" type="email" />
      </el-form-item>
      <el-form-item label="电话">
        <el-input v-model="form.phone" />
      </el-form-item>
      <el-form-item label="密码" v-if="!user">
        <el-input v-model="form.password" type="password" />
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
import { createUser, updateUser } from '../../api/user'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
  user: Object,
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(props.modelValue)
const loading = ref(false)
const form = ref({
  username: '',
  email: '',
  phone: '',
  password: '',
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.user) {
    form.value = {
      username: props.user.username || '',
      email: props.user.email || '',
      phone: props.user.phone || '',
      password: '',
    }
  } else if (val) {
    form.value = { username: '', email: '', phone: '', password: '' }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  if (!form.value.username || !form.value.email) {
    ElMessage.warning('请填写必填字段')
    return
  }
  loading.value = true
  try {
    if (props.user) {
      await updateUser(props.user.id, form.value)
    } else {
      await createUser(form.value)
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