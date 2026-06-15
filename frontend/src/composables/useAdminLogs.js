import { ref, computed, onMounted, watch } from 'vue'
import api from '../services/api'
import { auth } from '../firebase'
import { onAuthStateChanged } from 'firebase/auth'

export function useAdminLogs() {
  const logs = ref([])
  const loading = ref(false)

  const filterEvent = ref('')
  const filterUserID = ref('')
  const filterBookingID = ref('')
  const filterDateFrom = ref('')
  const filterDateTo = ref('')

  const eventTypes = ['booking_success', 'booking_timeout', 'seat_release', 'lock_failed']

  async function fetchLogs() {
    loading.value = true
    try {
      const params = {}
      if (filterEvent.value) params.event = filterEvent.value
      if (filterUserID.value) params.user_id = filterUserID.value
      if (filterBookingID.value) params.booking_id = filterBookingID.value

      const res = await api.get('/admin/event-logs', { params })
      logs.value = res.data || []
    } catch (err) {
      console.error('โหลด logs ไม่สำเร็จ', err)
    } finally {
      loading.value = false
    }
  }

  function toBangkokDateString(iso) {
    if (!iso) return ''
    return new Date(iso).toLocaleDateString('en-CA', { timeZone: 'Asia/Bangkok' })
  }

  const filteredLogs = computed(() => {
    return logs.value
      .filter(l => !filterDateFrom.value || toBangkokDateString(l.created_at) >= filterDateFrom.value)
      .filter(l => !filterDateTo.value || toBangkokDateString(l.created_at) <= filterDateTo.value)
  })

  watch([filterEvent, filterUserID, filterBookingID], fetchLogs)

  onMounted(() => {
    if (auth.currentUser) {
      fetchLogs()
    } else {
      const unsubscribe = onAuthStateChanged(auth, (user) => {
        unsubscribe()
        if (user) fetchLogs()
      })
    }
  })

  return {
    logs, filteredLogs, loading, eventTypes,
    filterEvent, filterUserID, filterBookingID, filterDateFrom, filterDateTo,
    fetchLogs
  }
}