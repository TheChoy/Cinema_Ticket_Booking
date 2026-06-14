<template>
  <div class="admin-page">
    <AdminHeaderBar />

    <div class="admin-content">
      <div class="admin-top">
        <input v-model="filterEmail" placeholder="ค้นหาอีเมล..." />
        <input v-model="filterMovie" placeholder="ค้นหาชื่อหนัง..." />
        <input v-model="filterDate" type="date" />
        <select v-model="filterStatus">
          <option value="">ทุกสถานะ</option>
          <option value="pending">Pending</option>
          <option value="paid">Paid</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>

      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="!filteredBookings.length" class="state-box">ไม่พบรายการจอง</div>
      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Booking No.</th>
              <th>Email</th>
              <th>หนัง</th>
              <th>รอบฉาย</th>
              <th>จำนวนที่นั่ง</th>
              <th>ราคารวม</th>
              <th>สถานะ</th>
              <th>วันที่จอง</th>
              <th style="width:160px">จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in filteredBookings" :key="b._id">
              <td>{{ b.booking_number }}</td>
              <td>{{ b.user_email }}</td>
              <td>{{ b.movie_title }}</td>
              <td>{{ formatDate(b.start_time) }}</td>
              <td>{{ b.seat_ids?.length || 0 }}</td>
              <td>฿{{ b.total_price }}</td>
              <td>
                <span class="badge" :class="statusClass(b.status)">{{ b.status }}</span>
              </td>
              <td>{{ formatDate(b.created_at) }}</td>
              <td>
                <div class="actions">
                  <select :value="b.status" @change="onStatusChange(b, $event.target.value)">
                    <option value="pending">Pending</option>
                    <option value="paid">Paid</option>
                    <option value="cancelled">Cancelled</option>
                  </select>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import AdminHeaderBar from '../components/AdminHeaderBar.vue'
import { useAdminBookings } from '../composables/useAdminBookings.js'
import '../assets/styles/admin-movies.css'

const { filteredBookings, loading, filterStatus, filterEmail, filterMovie, filterDate, updateStatus, cancelBooking } = useAdminBookings()

function statusClass(status) {
  if (status === 'paid') return 'badge-showing'
  if (status === 'cancelled') return 'badge-soon'
  return 'badge-genre'
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('th-TH', {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
    timeZone: 'Asia/Bangkok'
  })
}

async function onStatusChange(b, newStatus) {
  if (newStatus === b.status) return
  if (newStatus === 'cancelled') {
    if (!confirm(`ยกเลิก booking ${b.booking_number}? ที่นั่งจะถูกคืน`)) return
    await cancelBooking(b._id)
  } else {
    await updateStatus(b._id, newStatus)
  }
}
</script>