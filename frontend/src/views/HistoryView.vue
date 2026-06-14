<template>
  <div class="history-page">
    <HeaderBar />
    <Toast v-if="showToast" message="จองตั๋วสำเร็จแล้ว 🎬" />

    <div class="history-content">
      <div class="history-title">ประวัติการจอง</div>

      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="error" class="state-box">⚠️ {{ error }}</div>
      <div v-else-if="!bookings.length" class="state-box">
        <span style="font-size:2rem">🎬</span><br>ยังไม่มีประวัติการจอง
      </div>

      <div v-else>
        <div
          v-for="b in bookings"
          :key="b._id"
          class="booking-card"
        >
          <div class="booking-poster">
            <img v-if="b.poster_url" :src="b.poster_url" :alt="b.movie_title" />
            <span v-else>🎬</span>
          </div>
          <div class="booking-info">
            <div class="booking-movie">{{ b.movie_title }}</div>
            <div class="booking-meta">{{ formatDateTime(b.showtime) }}</div>
            <div class="booking-meta">ที่นั่ง: {{ b.seats?.join(', ') }}</div>
            <div class="booking-meta">{{ b.booking_number }}</div>
            <div class="booking-footer">
              <span class="booking-price">฿{{ b.total_price }}</span>
              <span
                class="status-badge"
                :class="{
                  'status-paid': b.status === 'paid',
                  'status-pending': b.status === 'pending',
                  'status-cancelled': b.status === 'cancelled'
                }"
              >
                {{
                  b.status === 'paid' ? 'ชำระแล้ว' :
                  b.status === 'pending' ? 'รอชำระ' : 'ยกเลิก'
                }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import HeaderBar from '../components/HeaderBar.vue'
import Toast from '../components/Toast.vue'
import { useHistory } from '../composables/useHistory.js'
import '../assets/styles/history.css'

const route = useRoute()
const showToast = ref(false)
const { bookings, loading, error, fetchHistory, formatDateTime } = useHistory()

onMounted(() => {
  fetchHistory()
  if (route.query.success === '1') {
    showToast.value = true
  }
})
</script>