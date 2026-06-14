import { ref, computed } from 'vue'
import api from '../services/api'

export function useSeat(showtimeId) {
  const seats = ref([])
  const showtime = ref(null)
  const selectedSeats = ref([])
  const loading = ref(false)
  const error = ref('')
  const MAX_SEATS = 8


  async function fetchSeats() {
    loading.value = true
    try {
      const res = await api.get('/seats', { params: { showtime_id: showtimeId } })
      seats.value = res.data || []
    } catch {
      error.value = 'โหลดที่นั่งไม่สำเร็จ'
    } finally {
      loading.value = false
    }
  }

  async function fetchShowtime() {
  try {
    const res = await api.get('/showtimes')
    const all = res.data || []
    showtime.value = all.find(s => s.id === showtimeId) || null
  } catch {}
}

  // จัดกลุ่มตาม row
  const groupedSeats = computed(() => {
    const groups = {}
    seats.value.forEach(s => {
      if (!groups[s.row]) groups[s.row] = []
      groups[s.row].push(s)
    })
    return groups
  })

function toggleSeat(seat) {
  if (seat.status !== 'available') return
  const idx = selectedSeats.value.findIndex(s => s.id === seat.id)
  if (idx === -1) {
    if (selectedSeats.value.length >= MAX_SEATS) return // ไม่ให้เลือกเกิน
    selectedSeats.value.push(seat)
  } else {
    selectedSeats.value.splice(idx, 1)
  }
}
  function isSelected(seat) {
    return selectedSeats.value.some(s => s.id === seat.id)
  }

  const totalPrice = computed(() => {
    if (!showtime.value) return 0
    return selectedSeats.value.length * showtime.value.price
  })

  return {
    seats, showtime, selectedSeats, groupedSeats,
    loading, error,
    fetchSeats, fetchShowtime, toggleSeat, isSelected, totalPrice
  }
}