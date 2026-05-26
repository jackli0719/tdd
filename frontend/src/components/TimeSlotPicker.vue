<template>
  <div class="time-slot-picker">
    <el-select
      v-model="selectedSlot"
      placeholder="选择时间段"
      :disabled="disabled || slots.length === 0"
      @change="handleSelect"
    >
      <el-option
        v-for="slot in slots"
        :key="slot.start_time"
        :label="`${slot.start_time} - ${slot.end_time}`"
        :value="slot.start_time"
        :disabled="!slot.available"
      />
    </el-select>
    <div v-if="loading" class="loading-tip">加载中...</div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: String,
  date: String,
  disabled: Boolean,
  slots: {
    type: Array,
    default: () => [],
  },
  loading: Boolean,
})

const emit = defineEmits(['update:modelValue', 'slots-updated'])

const selectedSlot = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  selectedSlot.value = val
})

watch(() => props.slots, (val) => {
  emit('slots-updated', val)
}, { immediate: true })

const handleSelect = (val) => {
  emit('update:modelValue', val)
}
</script>

<style scoped>
.time-slot-picker {
  display: inline-block;
}
.loading-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>