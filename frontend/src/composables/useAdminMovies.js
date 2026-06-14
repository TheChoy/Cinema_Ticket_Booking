import { ref, computed } from 'vue'
import api from '../services/api'

export function useAdminMovies() {
  const movies = ref([])
  const loading = ref(false)
  const saving = ref(false)

  const searchQuery = ref('')
  const filterGenre = ref('')
  const filterStatus = ref('')

  const genres = ['Action', 'Drama', 'Comedy', 'Horror', 'Sci-Fi', 'Romance', 'Animation']

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

  return {
    movies, filteredMovies, loading, saving,
    searchQuery, filterGenre, filterStatus, genres,
    fetchMovies, createMovie, updateMovie, deleteMovie, formatDuration
  }
}