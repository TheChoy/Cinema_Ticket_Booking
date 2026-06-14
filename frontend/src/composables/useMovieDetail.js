import { ref, computed } from 'vue'
import api from '../services/api'

export function useMovieDetail(movieId) {
  const movie = ref(null)
  const showtimes = ref([])
  const selectedShowtime = ref(null)
  const loading = ref(false)
  const error = ref('')

  async function fetchMovie() {
    loading.value = true
    try {
      const res = await api.get(`/movies/${movieId}`)
      movie.value = res.data
    } catch {
      error.value = 'โหลดหนังไม่สำเร็จ'
    } finally {
      loading.value = false
    }
  }

  async function fetchShowtimes() {
    try {
      const res = await api.get('/showtimes', { params: { movie_id: movieId } })
      showtimes.value = res.data || []
    } catch {
      showtimes.value = []
    }
  }

  // จัดกลุ่มตามวันที่
const groupedShowtimes = computed(() => {
  const groups = {}
  
  const sorted = [...showtimes.value].sort((a, b) => 
    new Date(a.start_time) - new Date(b.start_time)
  )
  
  sorted.forEach(s => {
    const date = new Date(s.start_time).toLocaleDateString('th-TH', {
      weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
      timeZone: 'Asia/Bangkok'
    })
    if (!groups[date]) groups[date] = []
    groups[date].push(s)
  })
  return groups
})

  function formatTime(iso) {
    return new Date(iso).toLocaleTimeString('th-TH', {
      hour: '2-digit', minute: '2-digit'
    })
  }

  return {
    movie, showtimes, groupedShowtimes,
    selectedShowtime, loading, error,
    fetchMovie, fetchShowtimes, formatTime
  }
}