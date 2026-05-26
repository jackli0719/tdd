<template>
  <el-dialog
    v-model="visible"
    title="创建订单"
    width="600px"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="用户">
        <el-select v-model="form.user_id" placeholder="请选择用户">
          <el-option
            v-for="user in users"
            :key="user.id"
            :label="user.username + ' (' + user.email + ')'"
            :value="user.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="预约时间">
        <div class="appointment-row">
          <DatePicker v-model="appointmentDate" @change="onDateChange" />
          <TimeSlotPicker
            v-model="appointmentTime"
            :date="appointmentDate"
            :slots="slots"
            :loading="slotsLoading"
            :disabled="!appointmentDate"
            @slots-updated="onSlotsUpdated"
          />
        </div>
      </el-form-item>
      <el-form-item label="产品">
        <div class="product-items">
          <div v-for="(item, index) in form.items" :key="index" class="product-item-row">
            <el-select v-model="item.product_id" placeholder="选择产品" @change="onProductChange(index)">
              <el-option
                v-for="product in products"
                :key="product.id"
                :label="product.name + ' (¥' + product.price + ', 库存:' + product.stock + ')'"
                :value="product.id"
              />
            </el-select>
            <el-input-number v-model="item.quantity" :min="1" :max="getMaxStock(index)" />
            <el-button type="danger" @click="removeItem(index)">删除</el-button>
          </div>
          <el-button type="primary" @click="addItem">添加产品</el-button>
        </div>
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
import { createOrder } from '../../api/order'
import { getUsers } from '../../api/user'
import { getProducts } from '../../api/product'
import { getSlots } from '../../api/slot'
import DatePicker from '../../components/DatePicker.vue'
import TimeSlotPicker from '../../components/TimeSlotPicker.vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(props.modelValue)
const loading = ref(false)
const users = ref([])
const products = ref([])
const appointmentDate = ref('')
const appointmentTime = ref('')
const slots = ref([])
const slotsLoading = ref(false)
const form = ref({
  user_id: null,
  items: [{ product_id: null, quantity: 1 }],
})

watch(() => props.modelValue, async (val) => {
  visible.value = val
  if (val) {
    await loadUsers()
    await loadProducts()
    appointmentDate.value = ''
    appointmentTime.value = ''
    slots.value = []
    form.value = { user_id: null, items: [{ product_id: null, quantity: 1 }] }
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const loadUsers = async () => {
  try {
    const res = await getUsers()
    users.value = res.data?.users || res?.users || []
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

const loadProducts = async () => {
  try {
    const res = await getProducts()
    products.value = res.data?.products || res?.products || []
  } catch (error) {
    console.error('Failed to load products:', error)
  }
}

const onDateChange = async (date) => {
  appointmentTime.value = ''
  if (!date) {
    slots.value = []
    return
  }
  slotsLoading.value = true
  try {
    const res = await getSlots(date)
    slots.value = res.data?.slots || []
  } catch (error) {
    slots.value = []
  } finally {
    slotsLoading.value = false
  }
}

const onSlotsUpdated = (updatedSlots) => {
  slots.value = updatedSlots
}

const onProductChange = (index) => {
  form.value.items[index].quantity = 1
}

const getMaxStock = (index) => {
  const productId = form.value.items[index].product_id
  const product = products.value.find(p => p.id === productId)
  return product ? product.stock : 999
}

const addItem = () => {
  form.value.items.push({ product_id: null, quantity: 1 })
}

const removeItem = (index) => {
  if (form.value.items.length > 1) {
    form.value.items.splice(index, 1)
  }
}

const buildAppointmentTime = () => {
  if (!appointmentDate.value || !appointmentTime.value) return null
  return appointmentDate.value + ' ' + appointmentTime.value + ':00'
}

const handleSubmit = async () => {
  if (!form.value.user_id) {
    ElMessage.warning('请选择用户')
    return
  }

  const validItems = form.value.items.filter(item => item.product_id && item.quantity > 0)
  if (validItems.length === 0) {
    ElMessage.warning('请添加至少一个产品')
    return
  }

  const appointmentTime = buildAppointmentTime()
  const orderData = {
    user_id: form.value.user_id,
    items: validItems,
  }
  if (appointmentTime) {
    orderData.appointment_time = appointmentTime
  }

  loading.value = true
  try {
    await createOrder(orderData)
    ElMessage.success('创建成功')
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error(error.message || '创建失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.product-items {
  width: 100%;
}

.product-item-row {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  align-items: center;
}

.product-item-row .el-select {
  flex: 1;
}

.appointment-row {
  display: flex;
  gap: 10px;
  align-items: center;
}
</style>