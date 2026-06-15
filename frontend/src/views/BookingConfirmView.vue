<template>
  <div class="confirm-page">
    <HeaderBar />

    <div class="confirm-content">
      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="error" class="error-box">⚠️ {{ error }}</div>

      <template v-else-if="booking">
        <div class="confirm-card">
          <!-- Header -->
          <div class="confirm-header">
            <h2>สรุปการจอง</h2>
            <span class="booking-number">{{ booking.booking_number }}</span>
          </div>

          <!-- Countdown -->
          <div class="countdown-bar">
            <div class="countdown-text">
              ⏱ กรุณาชำระเงินภายใน
              <span :class="{ 'countdown-urgent': countdown <= 60 }">
                {{ formatCountdown(countdown) }}
              </span>
            </div>
            <div class="countdown-track">
              <div
                class="countdown-fill"
                :style="{ width: (countdown / 300 * 100) + '%' }"
                :class="{ urgent: countdown <= 60 }"
              ></div>
            </div>
          </div>

          <!-- Movie -->
          <div class="confirm-movie">
            <img
              v-if="booking.poster_url"
              :src="booking.poster_url"
              class="confirm-poster"
              :alt="booking.movie_title"
            />
            <div v-else class="confirm-poster" style="display:flex;align-items:center;justify-content:center">🎬</div>
            <div class="confirm-movie-info">
              <div class="confirm-movie-title">{{ booking.movie_title }}</div>
              <div class="confirm-movie-meta">{{ formatDateTime(booking.start_time) }}</div>
              <div class="confirm-movie-meta">{{ booking.room }}</div>
            </div>
          </div>

          <!-- Details -->
          <div class="confirm-details">
            <div class="confirm-row">
              <span class="confirm-label">ที่นั่ง</span>
              <div class="seats-list">
                <span v-for="s in booking.seats" :key="s" class="seat-tag">{{ s }}</span>
              </div>
            </div>
            <div class="confirm-row">
              <span class="confirm-label">จำนวน</span>
              <span class="confirm-value">{{ booking.seats?.length }} ที่นั่ง</span>
            </div>
            <div class="confirm-row">
              <span class="confirm-label">สถานะ</span>
              <span class="status-badge status-pending">รอชำระเงิน</span>
            </div>
          </div>

          <!-- Total -->
          <div class="confirm-total">
            <span class="total-label">ราคารวม</span>
            <span class="total-price">฿{{ booking.total_price }}</span>
          </div>
        </div>

        <!-- Payment Options -->
        <div class="section-title">เลือกวิธีชำระเงิน</div>
        <div class="payment-options">
          <button
            class="payment-option"
            :class="{ selected: selectedPayment === 'qr' }"
            @click="selectedPayment = 'qr'"
          >
            <span class="payment-icon">📱</span>
            <div class="payment-info">
              <div class="payment-name">QR Code พร้อมเพย์</div>
              <div class="payment-desc">สแกนจ่ายผ่านแอปธนาคาร</div>
            </div>
            <span class="payment-check" v-if="selectedPayment === 'qr'">✓</span>
          </button>
          <button class="payment-option disabled" disabled>
            <span class="payment-icon">💳</span>
            <div class="payment-info">
              <div class="payment-name">บัตรเครดิต</div>
              <div class="payment-desc">ไม่พร้อมใช้งาน</div>
            </div>
            <span class="coming-soon-badge">เร็วๆ นี้</span>
          </button>
          <button class="payment-option disabled" disabled>
            <span class="payment-icon">🏦</span>
            <div class="payment-info">
              <div class="payment-name">โอนเงินผ่านธนาคาร</div>
              <div class="payment-desc">ไม่พร้อมใช้งาน</div>
            </div>
            <span class="coming-soon-badge">เร็วๆ นี้</span>
          </button>
        </div>

        <!-- Actions -->
        <div class="confirm-actions">
          <button
            class="btn-pay"
            :disabled="!selectedPayment"
            @click="showQRModal = true"
          >
            ชำระเงิน
          </button>
          <button class="btn-back" @click="router.back()">ย้อนกลับ</button>
        </div>
      </template>
    </div>

    <!-- QR Modal -->
    <div v-if="showQRModal" class="modal-overlay" @click.self="showQRModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>สแกน QR Code</h3>
          <button class="modal-close" @click="showQRModal = false">✕</button>
        </div>
        <div class="qr-section">
          <div class="qr-amount">฿{{ booking?.total_price }}</div>
          <div class="qr-box">
            <svg viewBox="0 0 100 100" width="180" height="180" xmlns="http://www.w3.org/2000/svg">
              <rect width="100" height="100" fill="white"/>
              <rect x="5" y="5" width="25" height="25" fill="none" stroke="#111" stroke-width="3"/>
              <rect x="10" y="10" width="15" height="15" fill="#111"/>
              <rect x="70" y="5" width="25" height="25" fill="none" stroke="#111" stroke-width="3"/>
              <rect x="75" y="10" width="15" height="15" fill="#111"/>
              <rect x="5" y="70" width="25" height="25" fill="none" stroke="#111" stroke-width="3"/>
              <rect x="10" y="75" width="15" height="15" fill="#111"/>
              <rect x="35" y="5" width="5" height="5" fill="#111"/>
              <rect x="45" y="5" width="5" height="5" fill="#111"/>
              <rect x="55" y="10" width="5" height="5" fill="#111"/>
              <rect x="35" y="15" width="10" height="5" fill="#111"/>
              <rect x="50" y="20" width="5" height="5" fill="#111"/>
              <rect x="40" y="25" width="5" height="10" fill="#111"/>
              <rect x="35" y="35" width="5" height="5" fill="#111"/>
              <rect x="45" y="35" width="20" height="5" fill="#111"/>
              <rect x="70" y="35" width="5" height="5" fill="#111"/>
              <rect x="80" y="35" width="15" height="5" fill="#111"/>
              <rect x="35" y="45" width="10" height="5" fill="#111"/>
              <rect x="50" y="45" width="5" height="10" fill="#111"/>
              <rect x="60" y="50" width="10" height="5" fill="#111"/>
              <rect x="75" y="45" width="5" height="10" fill="#111"/>
              <rect x="85" y="45" width="10" height="5" fill="#111"/>
              <rect x="35" y="55" width="5" height="10" fill="#111"/>
              <rect x="45" y="60" width="15" height="5" fill="#111"/>
              <rect x="65" y="55" width="5" height="5" fill="#111"/>
              <rect x="80" y="55" width="15" height="5" fill="#111"/>
              <rect x="35" y="70" width="5" height="5" fill="#111"/>
              <rect x="45" y="70" width="10" height="5" fill="#111"/>
              <rect x="60" y="70" width="5" height="10" fill="#111"/>
              <rect x="70" y="70" width="5" height="5" fill="#111"/>
              <rect x="80" y="70" width="5" height="5" fill="#111"/>
              <rect x="90" y="70" width="5" height="5" fill="#111"/>
              <rect x="35" y="80" width="5" height="5" fill="#111"/>
              <rect x="45" y="80" width="5" height="5" fill="#111"/>
              <rect x="55" y="80" width="10" height="5" fill="#111"/>
              <rect x="70" y="80" width="10" height="5" fill="#111"/>
              <rect x="85" y="80" width="10" height="5" fill="#111"/>
              <rect x="35" y="90" width="10" height="5" fill="#111"/>
              <rect x="50" y="90" width="5" height="5" fill="#111"/>
              <rect x="60" y="90" width="10" height="5" fill="#111"/>
              <rect x="75" y="90" width="5" height="5" fill="#111"/>
              <rect x="85" y="90" width="10" height="5" fill="#111"/>
            </svg>
          </div>
          <div class="qr-hint">สแกนด้วยแอปธนาคารหรือแอปพร้อมเพย์</div>
          <button class="btn-confirm-pay" :disabled="paying" @click="payBooking">
            {{ paying ? 'กำลังดำเนินการ...' : '✅ ยืนยันการชำระเงิน' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HeaderBar from '../components/HeaderBar.vue'
import { useBookingConfirm } from '../composables/useBookingConfirm.js'
import '../assets/styles/booking-confirm.css'

const route = useRoute()
const router = useRouter()
const bookingId = route.params.id

const {
  booking, loading, paying, error, countdown,
  fetchBooking, payBooking, formatDateTime, formatCountdown
} = useBookingConfirm(bookingId)

const showQRModal = ref(false)
const selectedPayment = ref('qr')

onMounted(() => fetchBooking())
</script>