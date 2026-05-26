<template>
  <el-date-picker
    :model-value="modelValue"
    type="date"
    placeholder="选择日期"
    :disabled-date="disabledDate"
    format="YYYY-MM-DD"
    value-format="YYYY-MM-DD"
    @update:model-value="handleUpdate"
    @change="handleChange"
  />
</template>

<script setup>
const props = defineProps({
  modelValue: String,
})
const emit = defineEmits(['update:modelValue', 'change'])

const handleUpdate = (val) => {
  emit('update:modelValue', val)
  emit('change', val)
}

const handleChange = (val) => {
  emit('change', val)
}

// Disable past dates
const disabledDate = (date) => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return date.getTime() < today.getTime()
}
</script>