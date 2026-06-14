import { ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

export function useBookingConfirm(bookingId) {
  const booking = ref(null)
  const loading = ref(false)
  const paying = ref(false)
  const error = ref('')
  const countdown = ref(300)
  const router = useRouter()
  let ws = null
  let timer = null

  function startTimer() {
    timer = setInterval(() => {
      if (countdown.value > 0) {
        countdown.value--
      } else {
        clearInterval(timer)
        router.push('/home')
      }
    }, 1000)
  }

  async function fetchBooking() {
    loading.value = true
    try {
      const res = await api.get(`/bookings/${bookingId}`)
      booking.value = res.data
      connectWS(res.data.showtime_id)
      startTimer()
    } catch (err) {
      error.value = 'โหลดข้อมูลการจองไม่สำเร็จ'
    } finally {
      loading.value = false
    }
  }

  function connectWS(showtimeId) {
    ws = new WebSocket(`ws://localhost:8081/ws/showtimes/${showtimeId}`)
    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data)
      if (msg.type === 'countdown' && msg.booking_id === bookingId) {
        countdown.value = msg.remaining_seconds
        if (countdown.value <= 0) {
          clearInterval(timer)
          router.push('/home')
        }
      }
    }
  }

  async function payBooking() {
    paying.value = true
    try {
      await api.put(`/bookings/${bookingId}/pay`)
      router.push(`/bookings/${bookingId}/success`)
    } catch {
      error.value = 'ชำระเงินไม่สำเร็จ กรุณาลองใหม่'
    } finally {
      paying.value = false
    }
  }

  function formatDateTime(iso) {
    return new Date(iso).toLocaleString('th-TH', {
      weekday: 'short', year: 'numeric', month: 'short',
      day: 'numeric', hour: '2-digit', minute: '2-digit',
      timeZone: 'Asia/Bangkok'
    })
  }

  function formatCountdown(seconds) {
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    return `${m}:${String(s).padStart(2, '0')}`
  }

  onUnmounted(() => {
    if (ws) ws.close()
    if (timer) clearInterval(timer)
  })

  return {
    booking, loading, paying, error, countdown,
    fetchBooking, payBooking, formatDateTime, formatCountdown
  }
}