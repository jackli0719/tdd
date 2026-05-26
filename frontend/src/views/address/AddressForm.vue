<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? '编辑地址' : '新增地址'"
    width="500px"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="联系人">
        <el-input v-model="form.name" placeholder="请输入联系人姓名" />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="form.phone" placeholder="请输入手机号" />
      </el-form-item>
      <el-form-item label="省份">
        <el-input v-model="form.province" placeholder="请输入省份" />
      </el-form-item>
      <el-form-item label="城市">
        <el-input v-model="form.city" placeholder="请输入城市" />
      </el-form-item>
      <el-form-item label="区县">
        <el-input v-model="form.district" placeholder="请输入区县" />
      </el-form-item>
      <el-form-item label="详细地址">
        <el-input v-model="form.detail" type="textarea" placeholder="请输入详细地址" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { createAddress, updateAddress } from '../../api/address'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
  address: Object,
  userId: {
    type: [Number, String],
    default: () => 1,
  },
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(props.modelValue)
const loading = ref(false)
const effectiveUserId = ref(Number(props.userId) || 1)

onMounted(() => {
  effectiveUserId.value = Number(props.userId) || 1
})

watch(() => props.userId, (newId) => {
  effectiveUserId.value = Number(newId) || 1
})

const form = ref({
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
})

const isEdit = computed(() => !!props.address?.id)

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    if (props.address) {
      form.value = {
        name: props.address.name || '',
        phone: props.address.phone || '',
        province: props.address.province || '',
        city: props.address.city || '',
        district: props.address.district || '',
        detail: props.address.detail || '',
      }
    } else {
      form.value = {
        name: '',
        phone: '',
        province: '',
        city: '',
        district: '',
        detail: '',
      }
    }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleSubmit = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入联系人姓名')
    return
  }
  if (!form.value.phone) {
    ElMessage.warning('请输入手机号')
    return
  }

  loading.value = true
  try {
    const data = {
      user_id: effectiveUserId.value,
      ...form.value,
    }
    if (isEdit.value) {
      await updateAddress(props.address.id, data)
      ElMessage.success('更新成功')
    } else {
      await createAddress(data)
      ElMessage.success('创建成功')
    }
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}
</script>