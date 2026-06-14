<template>
  <div class="home-page">
    <HeaderBar />

    <div class="home-content">
      <!-- Tabs -->
      <div class="tabs">
        <button
          class="tab-btn"
          :class="{ active: status === 'now_showing' }"
          @click="status = 'now_showing'"
        >Now Showing</button>
        <button
          class="tab-btn"
          :class="{ active: status === 'coming_soon' }"
          @click="status = 'coming_soon'"
        >Coming Soon</button>
      </div>

      <!-- Filter bar -->
      <div class="filter-bar">
        <button
          class="btn-filter-toggle"
          :class="{ active: showFilter }"
          @click="showFilter = !showFilter"
        >
          ☰ Genre
        </button>
        <input
          class="filter-search"
          v-model="search"
          placeholder="ค้นหาชื่อหนัง..."
        />

      </div>

      <!-- Genre drawer -->
      <div v-if="showFilter" class="filter-drawer">
        <button
          v-for="g in genres"
          :key="g"
          class="genre-chip"
          :class="{ active: genre === g }"
          @click="selectGenre(g)"
        >{{ g }}</button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="state-box">
        <span>⏳</span> กำลังโหลด...
      </div>

      <!-- Error -->
      <div v-else-if="error" class="state-box">
        <span>⚠️</span> {{ error }}
      </div>

      <!-- Empty -->
      <div v-else-if="!movies.length" class="state-box">
        <span>🎬</span> ไม่พบหนัง
      </div>

      <!-- Grid -->
      <div v-else class="movie-grid">
        <div
          v-for="movie in movies"
          :key="movie.id"
          class="movie-card"
          @click="router.push(`/movies/${movie.id}`)"
        >
          <img
            v-if="movie.poster_url"
            :src="movie.poster_url"
            :alt="movie.title"
            class="movie-poster"
          />
          <div v-else class="movie-poster-placeholder">🎬</div>
            <div class="movie-info">
              <div class="movie-title" :title="movie.title">
                {{ movie.title }} · {{ formatDuration(movie.duration) }}
              </div>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import HeaderBar from '../components/HeaderBar.vue'
import { useMovies } from '../composables/useMovies.js'
import '../assets/styles/home.css'
import { useRouter } from 'vue-router'

const router = useRouter()

const { movies, loading, error, status, search, genre, fetchMovies, formatDuration } = useMovies()

const showFilter = ref(false)
const genres = ['Action', 'Drama', 'Comedy', 'Horror', 'Sci-Fi', 'Romance', 'Animation']

function selectGenre(g) {
  genre.value = genre.value === g ? '' : g  // toggle
}

onMounted(() => fetchMovies())
</script>