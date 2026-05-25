<template>
  <el-dialog
    v-model="visible"
    :title="user ? 'Edit User' : 'Add User'"
    width="500px"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="Username">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="Email">
        <el-input v-model="form.email" type="email" />
      </el-form-item>
      <el-form-item label="Password" v-if="!user">
        <el-input v-model="form.password" type="password" />
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
  password: '',
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.user) {
    form.value = {
      username: props.user.username || '',
      email: props.user.email || '',
      password: '',
    }
  } else if (val) {
    form.value = { username: '', email: '', password: '' }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  if (!form.value.username || !form.value.email) {
    ElMessage.warning('Please fill in required fields')
    return
  }
  loading.value = true
  try {
    if (props.user) {
      await updateUser(props.user.id, form.value)
    } else {
      await createUser(form.value)
    }
    ElMessage.success('Success')
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error(error.message || 'Failed to save user')
  } finally {
    loading.value = false
  }
}
</script>