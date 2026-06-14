import { ref, computed, reactive } from 'vue'
import api from '../services/api'

export function useAdminMovies() {
  const movies = ref([])
  const loading = ref(false)
  const saving = ref(false)

  const searchQuery = ref('')
  const filterGenre = ref('')
  const filterStatus = ref('')

  const genres = ['Action', 'Drama', 'Comedy', 'Horror', 'Sci-Fi', 'Romance', 'Animation']

  // Showtime
  const showtimes = ref([])
  const loadingShowtimes = ref(false)
  const savingShowtime = ref(false)
  const expandedMovie = ref(null)
  const currentMovieId = ref(null)
  const showShowtimeModal = ref(false)
  const editingShowtime = ref(null)
  const showtimeForm = reactive({
    room: '', start_time: '', end_time: '', seat_count: 120, price: 250
  })

  async function fetchMovies() {
    loading.value = true
    try {
      const res = await api.get('/movies')
      movies.value = res.data || []
    } finally {
      loading.value = false
    }
  }

  const filteredMovies = computed(() => {
    return movies.value
      .filter(m => !searchQuery.value || m.title.toLowerCase().includes(searchQuery.value.toLowerCase()))
      .filter(m => !filterGenre.value || m.genre.toLowerCase() === filterGenre.value.toLowerCase())
      .filter(m => !filterStatus.value || m.status === filterStatus.value)
  })

  async function createMovie(formData) {
    saving.value = true
    try {
      const res = await api.post('/admin/movies', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      movies.value.unshift(res.data)
      return true
    } catch {
      return false
    } finally {
      saving.value = false
    }
  }

  async function updateMovie(id, formData) {
    saving.value = true
    try {
      await api.put(`/admin/movies/${id}`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      await fetchMovies()
      return true
    } catch {
      return false
    } finally {
      saving.value = false
    }
  }

  async function deleteMovie(id) {
    try {
      await api.delete(`/admin/movies/${id}`)
      movies.value = movies.value.filter(m => m.id !== id)
      return true
    } catch {
      return false
    }
  }

  function formatDuration(min) {
    if (!min) return '-'
    return `${Math.floor(min / 60)}h ${String(min % 60).padStart(2, '0')}m`
  }

  // Showtime functions
  function formatDateTime(iso) {
    return new Date(iso).toLocaleString('th-TH', {
      month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit',
      timeZone: 'Asia/Bangkok'
    })
  }

  function toLocalDatetime(iso) {
    const d = new Date(iso)
    d.setMinutes(d.getMinutes() - d.getTimezoneOffset())
    return d.toISOString().slice(0, 16)
  }

  async function toggleShowtime(movie) {
    if (expandedMovie.value === movie.id) {
      expandedMovie.value = null
      return
    }
    expandedMovie.value = movie.id
    currentMovieId.value = movie.id
    loadingShowtimes.value = true
    try {
      const res = await api.get('/showtimes', { params: { movie_id: movie.id } })
      showtimes.value = res.data || []
    } finally {
      loadingShowtimes.value = false
    }
  }

  function openCreateShowtime(movie) {
    editingShowtime.value = null
    currentMovieId.value = movie.id
    Object.assign(showtimeForm, { room: '', start_time: '', end_time: '', seat_count: 120, price: 250 })
    showShowtimeModal.value = true
  }

  function openEditShowtime(s) {
    editingShowtime.value = s
    Object.assign(showtimeForm, {
      room: s.room,
      start_time: toLocalDatetime(s.start_time),
      end_time: toLocalDatetime(s.end_time),
      seat_count: s.seat_count,
      price: s.price
    })
    showShowtimeModal.value = true
  }

  async function saveShowtime() {
    savingShowtime.value = true
    try {
      const body = {
        movie_id: currentMovieId.value,
        room: showtimeForm.room,
        start_time: new Date(showtimeForm.start_time).toISOString(),
        end_time: new Date(showtimeForm.end_time).toISOString(),
        seat_count: showtimeForm.seat_count,
        price: showtimeForm.price
      }

      if (editingShowtime.value) {
        await api.put(`/admin/showtimes/${editingShowtime.value.id}`, body)
      } else {
        const res = await api.post('/admin/showtimes', body)
        const showtimeId = res.data.id

        const seatsPerRow = 10
        const rowCount = Math.ceil(showtimeForm.seat_count / seatsPerRow)
        const rows = Array.from({ length: rowCount }, (_, i) =>
          String.fromCharCode(65 + i)
        )

        await api.post('/admin/seats/generate', {
          showtime_id: showtimeId,
          rows,
          seats_per_row: seatsPerRow
        })
      }

      const res = await api.get('/showtimes', { params: { movie_id: currentMovieId.value } })
      showtimes.value = res.data || []
      showShowtimeModal.value = false
    } finally {
      savingShowtime.value = false
    }
  }

  async function confirmDeleteShowtime(s) {
    if (!confirm(`ลบรอบ ${formatDateTime(s.start_time)} ใช่ไหม?`)) return
    await api.delete(`/admin/showtimes/${s.id}`)
    showtimes.value = showtimes.value.filter(x => x.id !== s.id)
  }

  return {
    movies, filteredMovies, loading, saving,
    searchQuery, filterGenre, filterStatus, genres,
    fetchMovies, createMovie, updateMovie, deleteMovie, formatDuration,
    showtimes, loadingShowtimes, savingShowtime,
    expandedMovie, showShowtimeModal, editingShowtime, showtimeForm,
    toggleShowtime, openCreateShowtime, openEditShowtime,
    saveShowtime, confirmDeleteShowtime, formatDateTime
  }
}