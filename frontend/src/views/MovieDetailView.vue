<template>
  <div>
    <div v-if="loading">loading...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else>
      <h1>{{ movie.title }}</h1>
      <p>ID หนัง: {{ movie.id }}</p>
      <p>Description: {{ movie.description }}</p>
      <p>Genre: {{ movie.genre }}</p>
      <p>Duration: {{ movie.duration }} นาที</p>
      <img :src="movie.poster_url" :alt="movie.title" width="200" />

      <hr />
      <h2>รอบฉาย</h2>
      <div v-if="loadingShowtimes">loading showtimes...</div>
      <div v-else-if="!showtimes.length">ไม่มีรอบฉาย</div>
      <div v-for="s in showtimes" :key="s.id">
        <p>ID รอบ: {{ s.id }}</p>
        <p>ห้อง: {{ s.room }}</p>
        <p>เริ่ม: {{ s.start_time }}</p>
        <p>จบ: {{ s.end_time }}</p>
        <p>จำนวนที่นั่ง: {{ s.seat_count }}</p>
        <p>ราคา: {{ s.price }} บาท</p>
        <hr />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../services/api'

const route = useRoute()
const id = route.params.id

const movie = ref({})
const showtimes = ref([])
const loading = ref(false)
const loadingShowtimes = ref(false)
const error = ref('')

async function fetchMovie() {
  loading.value = true
  try {
    const res = await api.get(`/movies/${id}`)
    movie.value = res.data
  } catch (err) {
    error.value = 'โหลดหนังไม่สำเร็จ'
  } finally {
    loading.value = false
  }
}

async function fetchShowtimes() {
  loadingShowtimes.value = true
  try {
    const res = await api.get(`/showtimes`, { params: { movie_id: id } })
    showtimes.value = res.data || []
  } finally {
    loadingShowtimes.value = false
  }
}

onMounted(() => {
  fetchMovie()
  fetchShowtimes()
})
</script>