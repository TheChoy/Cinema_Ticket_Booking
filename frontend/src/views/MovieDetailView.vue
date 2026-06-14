<template>
  <div class="detail-page">
    <HeaderBar />

    <div class="detail-content">
      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="error" class="state-box">⚠️ {{ error }}</div>

      <template v-else-if="movie">
        <!-- Movie Info -->
        <div class="movie-section">
          <div class="movie-poster-wrap">
            <img v-if="movie.poster_url" :src="movie.poster_url" :alt="movie.title" />
            <div v-else class="poster-placeholder">🎬</div>
          </div>
          <div class="movie-meta">
            <h1>{{ movie.title }}</h1>
            <div class="meta-tags">
              <span class="tag">{{ movie.genre }}</span>
              <span class="tag">{{ movie.duration }} นาที</span>
            </div>
            <p class="meta-desc">{{ movie.description }}</p>
          </div>
        </div>

        <!-- Showtimes -->
        <div class="showtime-section">
          <h2>รอบฉาย</h2>

          <div v-if="!showtimes.length" class="state-box">ไม่มีรอบฉาย</div>

          <div
            v-for="(times, date) in groupedShowtimes"
            :key="date"
            class="date-group"
          >
            <div class="date-label">{{ date }}</div>
            <div class="showtime-list">
              <label v-for="s in times" :key="s.id">
                <input
                  class="showtime-radio"
                  type="radio"
                  name="showtime"
                  :value="s.id"
                  v-model="selectedShowtime"
                />
                <div class="showtime-card">
                  <div class="showtime-left">
                    <span class="showtime-time">
                      {{ formatTime(s.start_time) }} - {{ formatTime(s.end_time) }}
                    </span>
                    <span class="showtime-room">{{ s.room }}</span>
                  </div>
                  <div class="showtime-right">
                    <span class="showtime-price">฿{{ s.price }}</span>
                    <span class="showtime-seats">{{ s.seat_count }} ที่นั่ง</span>
                  </div>
                </div>
              </label>
            </div>
          </div>

          <button
            class="btn-book"
            :disabled="!selectedShowtime"
            @click="goToSeat"
          >
            {{ selectedShowtime ? 'เลือกที่นั่ง' : 'กรุณาเลือกรอบก่อน' }}
          </button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HeaderBar from '../components/HeaderBar.vue'
import { useMovieDetail } from '../composables/useMovieDetail.js'
import '../assets/styles/movie-detail.css'

const route = useRoute()
const router = useRouter()
const id = route.params.id

const {
  movie, showtimes, groupedShowtimes,
  selectedShowtime, loading, error,
  fetchMovie, fetchShowtimes, formatTime
} = useMovieDetail(id)

function goToSeat() {
  router.push(`/showtimes/${selectedShowtime.value}/seats`)
}

onMounted(() => {
  fetchMovie()
  fetchShowtimes()
})
</script>