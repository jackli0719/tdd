<template>
  <div class="address-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>地址管理</span>
          <el-button type="primary" @click="showForm">新增地址</el-button>
        </div>
      </template>
      <el-table :data="addresses" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="联系人" />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column label="地址">
          <template #default="{ row }">
            {{ formatAddress(row) }}
          </template>
        </el-table-column>
        <el-table-column prop="is_default" label="默认" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="editAddress(row)">编辑</el-button>
            <el-button size="small" type="primary" @click="setDefault(row)" :disabled="row.is_default">设为默认</el-button>
            <el-button size="small" type="danger" @click="deleteAddress(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <AddressForm
      v-model="formVisible"
      :address="currentAddress"
      @success="loadAddresses"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getAddresses, deleteAddress as deleteAddressApi, setDefaultAddress } from '../../api/address'
import AddressForm from './AddressForm.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps({
  userId: {
    type: [Number, String],
    required: true,
  },
})

const route = useRoute()

const addresses = ref([])
const loading = ref(false)
const formVisible = ref(false)
const currentAddress = ref(null)
const effectiveUserId = ref(props.userId)

onMounted(() => {
  effectiveUserId.value = props.userId
  loadAddresses()
})

watch(() => props.userId, (newId) => {
  effectiveUserId.value = newId
  loadAddresses()
})

const loadAddresses = async () => {
  loading.value = true
  try {
    const res = await getAddresses(effectiveUserId.value)
    addresses.value = res.data?.addresses || []
  } catch (error) {
    ElMessage.error('加载地址列表失败')
  } finally {
    loading.value = false
  }
}

const showForm = () => {
  currentAddress.value = null
  formVisible.value = true
}

const editAddress = (address) => {
  currentAddress.value = address
  formVisible.value = true
}

const setDefault = async (address) => {
  try {
    await setDefaultAddress(address.id, effectiveUserId.value)
    ElMessage.success('设置成功')
    loadAddresses()
  } catch (error) {
    ElMessage.error(error.message || '设置失败')
  }
}

const deleteAddress = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除该地址吗？', '提示', {
      type: 'warning',
    })
    await deleteAddressApi(id)
    ElMessage.success('删除成功')
    loadAddresses()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const formatAddress = (row) => {
  const parts = [row.province, row.city, row.district, row.detail].filter(Boolean)
  return parts.join(' ') || '-'
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>