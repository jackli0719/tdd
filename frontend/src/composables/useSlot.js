import { ref } from 'vue'
import { getSlots } from '../api/slot'

export const useSlot = () => {
  const loading = ref(false)
  const slots = ref([])

  const fetchSlots = async (date) => {
    if (!date) {
      slots.value = []
      return
    }
    loading.value = true
    try {
      const res = await getSlots(date)
      slots.value = res.data.slots || []
    } catch {
      slots.value = []
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    slots,
    fetchSlots,
  }
}