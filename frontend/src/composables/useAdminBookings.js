import { ref, computed, onMounted } from 'vue'
import api from '../services/api'
import { auth } from '../firebase'
import { onAuthStateChanged } from 'firebase/auth'

export function useAdminBookings() {
  const bookings = ref([])
  const loading = ref(false)
  const filterStatus = ref('')
  const searchQuery = ref('')
  const filterDate = ref('')

  async function fetchBookings() {
    loading.value = true
    try {
      const res = await api.get('/admin/bookings')
      bookings.value = res.data || []
    } catch (err) {
      console.error('โหลด bookings ไม่สำเร็จ', err)
    } finally {
      loading.value = false
    }
  }

  function toBangkokDateString(iso) {
    if (!iso) return ''
    return new Date(iso).toLocaleDateString('en-CA', { timeZone: 'Asia/Bangkok' })
  }

  const filteredBookings = computed(() => {
    const q = searchQuery.value.toLowerCase()
    return bookings.value
      .filter(b => !filterStatus.value || b.status === filterStatus.value)
      .filter(b => !q ||
        b.user_email?.toLowerCase().includes(q) ||
        b.movie_title?.toLowerCase().includes(q)
      )
      .filter(b => !filterDate.value || toBangkokDateString(b.created_at) === filterDate.value)
  })

  async function updateStatus(id, status) {
    try {
      await api.put(`/admin/bookings/${id}`, { status })
      const b = bookings.value.find(x => x._id === id)
      if (b) b.status = status
      return true
    } catch (err) {
      console.error('อัพเดทสถานะไม่สำเร็จ', err)
      return false
    }
  }

  async function cancelBooking(id) {
    try {
      await api.delete(`/admin/bookings/${id}`)
      const b = bookings.value.find(x => x._id === id)
      if (b) b.status = 'cancelled'
      return true
    } catch (err) {
      console.error('ยกเลิก booking ไม่สำเร็จ', err)
      return false
    }
  }

  onMounted(() => {
    if (auth.currentUser) {
      fetchBookings()
    } else {
      const unsubscribe = onAuthStateChanged(auth, (user) => {
        unsubscribe()
        if (user) fetchBookings()
      })
    }
  })

  return {
    bookings, filteredBookings, loading,
    filterStatus, searchQuery, filterDate,
    fetchBookings, updateStatus, cancelBooking
  }
}