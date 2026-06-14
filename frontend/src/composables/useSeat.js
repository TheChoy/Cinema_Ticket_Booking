import { ref, computed, onUnmounted } from 'vue'
import api from '../services/api'

export function useSeat(showtimeId) {
  const seats = ref([])
  const showtime = ref(null)
  const selectedSeats = ref([])
  const loading = ref(false)
  const error = ref('')
  let ws = null

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
      showtime.value = res.data?.find(s => s.id === showtimeId) || null
    } catch {}
  }

  function connectWS() {
    ws = new WebSocket(`ws://localhost:8081/ws/showtimes/${showtimeId}`)
    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data)
      if (msg.type === 'seat_update') {
        const seat = seats.value.find(s => s.id === msg.seat_id)
        if (seat) seat.status = msg.status

        // ถ้าที่นั่งที่เลือกไว้โดน lock/booked ให้เอาออก
        if (msg.status !== 'available') {
          selectedSeats.value = selectedSeats.value.filter(s => s.id !== msg.seat_id)
        }
      }
    }
    ws.onclose = () => {
      // reconnect อัตโนมัติถ้าหลุด
      setTimeout(() => connectWS(), 3000)
    }
  }

  const groupedSeats = computed(() => {
    const groups = {}
    seats.value.forEach(s => {
      if (!groups[s.row]) groups[s.row] = []
      groups[s.row].push(s)
    })
    return groups
  })

  const MAX_SEATS = 8

  function toggleSeat(seat) {
    if (seat.status !== 'available') return
    const idx = selectedSeats.value.findIndex(s => s.id === seat.id)
    if (idx === -1) {
      if (selectedSeats.value.length >= MAX_SEATS) return
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

  onUnmounted(() => {
    if (ws) ws.close()
  })

  return {
    seats, showtime, selectedSeats, groupedSeats,
    loading, error,
    fetchSeats, fetchShowtime, connectWS,
    toggleSeat, isSelected, totalPrice
  }
}