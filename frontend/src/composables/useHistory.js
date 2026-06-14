import { ref } from 'vue'
import api from '../services/api'

export function useHistory() {
  const bookings = ref([])
  const loading = ref(false)
  const error = ref('')

  async function fetchHistory() {
    loading.value = true
    try {
      const res = await api.get('/bookings/me')
      bookings.value = res.data || []
    } catch {
      error.value = 'โหลดประวัติไม่สำเร็จ'
    } finally {
      loading.value = false
    }
  }

  function formatDateTime(iso) {
    return new Date(iso).toLocaleString('th-TH', {
      year: 'numeric', month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit',
      timeZone: 'Asia/Bangkok'
    })
  }

  return { bookings, loading, error, fetchHistory, formatDateTime }
}