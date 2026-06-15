import { ref, watch } from 'vue'
import api from '../services/api'

export function useMovies() {
  const movies = ref([])
  const loading = ref(false)
  const error = ref('')

  const status = ref('now_showing')
  const search = ref('')
  const genre = ref('')

  async function fetchMovies() {
    loading.value = true
    error.value = ''
    try {
      const params = { status: status.value }
      if (search.value) params.search = search.value
      if (genre.value) params.genre = genre.value

      const res = await api.get('/movies', { params })
      movies.value = res.data || []
    } catch (err) {
      console.error(err)
      error.value = 'โหลดหนังไม่สำเร็จ'
    } finally {
      loading.value = false
    }
  }

  watch([status, genre], () => fetchMovies())

  let searchTimer = null
  watch(search, () => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => fetchMovies(), 400)
  })

  function formatDuration(minutes) {
  if (!minutes) return '-'
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return `${h}h ${String(m).padStart(2, '0')}m`
}

  return { movies, loading, error, status, search, genre, fetchMovies, formatDuration }
}