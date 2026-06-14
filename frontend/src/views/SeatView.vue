<template>
  <div class="seat-page">
    <HeaderBar />

    <div class="seat-content">
      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="error" class="state-box">⚠️ {{ error }}</div>

      <template v-else>
        <!-- Showtime info -->
        <div class="seat-info" v-if="showtime">
          <div class="seat-info-item">
            <span class="seat-info-label">ห้อง</span>
            <span class="seat-info-value">{{ showtime.room }}</span>
          </div>
          <div class="seat-info-item">
            <span class="seat-info-label">เริ่ม</span>
            <span class="seat-info-value">{{ formatTime(showtime.start_time) }}</span>
          </div>
          <div class="seat-info-item">
            <span class="seat-info-label">จบ</span>
            <span class="seat-info-value">{{ formatTime(showtime.end_time) }}</span>
          </div>
          <div class="seat-info-item">
            <span class="seat-info-label">ราคา/ที่นั่ง</span>
            <span class="seat-info-value">฿{{ showtime.price }}</span>
          </div>
        </div>

        <!-- Screen -->
        <div class="screen">
          <div class="screen-bar"></div>
          <div class="screen-label">SCREEN</div>
        </div>

        <!-- Legend -->
        <div class="legend">
        <div class="legend-item">
            <div class="legend-box available"></div> ว่าง
        </div>
        <div class="legend-item">
            <div class="legend-box booked"></div> ไม่ว่าง
        </div>
        </div>

        <!-- Seats -->
        <div class="seat-rows">
          <div v-for="(rowSeats, row) in groupedSeats" :key="row" class="seat-row">
            <div class="row-label">{{ row }}</div>
            <div class="seat-grid">
              <button
                v-for="seat in rowSeats"
                :key="seat.id"
                class="seat-btn"
                :class="{
                  selected: isSelected(seat),
                  booked: seat.status === 'booked',
                  locked: seat.status === 'locked'
                }"
                @click="toggleSeat(seat)"
              >
                {{ seat.number }}
              </button>
            </div>
          </div>
        </div>

        <!-- Summary -->
        <div class="seat-summary">
          <div class="summary-text">
            <span class="summary-count">
                เลือก {{ selectedSeats.length }}/8 ที่นั่ง
                {{ selectedSeats.length ? '(' + selectedSeats.map(s => s.label).join(', ') + ')' : '' }}
            </span>
            <span class="summary-price">฿{{ totalPrice }}</span>
          </div>
          <button
            class="btn-confirm"
            :disabled="!selectedSeats.length"
            @click="confirmBooking"
          >
            ยืนยันการจอง
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
import { useSeat } from '../composables/useSeat.js'
import api from '../services/api'  // ← เพิ่มตรงนี้
import '../assets/styles/seat.css'

const route = useRoute()
const router = useRouter()
const showtimeId = route.params.id

console.log('showtimeId:', showtimeId)

const {
  seats, showtime, selectedSeats, groupedSeats,
  loading, error,
  fetchSeats, fetchShowtime,connectWS, toggleSeat, isSelected, totalPrice
} = useSeat(showtimeId)

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('th-TH', {
    hour: '2-digit',
    minute: '2-digit',
    timeZone: 'Asia/Bangkok'
  })
}

async function confirmBooking() {
  try {
    const res = await api.post('/bookings', {
      showtime_id: showtimeId,
      seat_ids: selectedSeats.value.map(s => s.id)
    })
    const booking = res.data
    router.push(`/bookings/${booking.id}`) 
  } catch (err) {
    console.error(err)
    alert('เกิดข้อผิดพลาด กรุณาลองใหม่')
  }
}

onMounted(() => {
  fetchSeats()
  fetchShowtime()
  connectWS()
})
</script>