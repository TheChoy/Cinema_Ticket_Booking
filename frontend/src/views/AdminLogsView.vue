<template>
  <div class="admin-page">
    <AdminHeaderBar />

    <div class="admin-content">
      <div class="admin-top">
        <select v-model="filterEvent">
          <option value="">ทุก Event</option>
          <option v-for="e in eventTypes" :key="e" :value="e">{{ e }}</option>
        </select>
        <input v-model="filterUserID" placeholder="User ID (ObjectID)" />
        <input v-model="filterBookingID" placeholder="Booking ID (ObjectID)" />
        <input v-model="filterDateFrom" type="date" />
        <input v-model="filterDateTo" type="date" />
      </div>

      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="!filteredLogs.length" class="state-box">ไม่พบ log</div>
      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>เวลา</th>
              <th>Event</th>
              <th>Message</th>
              <th>User ID</th>
              <th>Booking ID</th>
              <th>Seat IDs</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="l in filteredLogs" :key="l.id">
              <td>{{ formatDate(l.created_at) }}</td>
              <td>
                <span class="badge" :class="eventClass(l.event)">{{ l.event }}</span>
              </td>
              <td>{{ l.message }}</td>
              <td class="mono">{{ l.user_id || '-' }}</td>
              <td class="mono">{{ l.booking_id || '-' }}</td>
              <td>{{ l.seat_ids?.length || 0 }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import AdminHeaderBar from '../components/AdminHeaderBar.vue'
import { useAdminLogs } from '../composables/useAdminLogs.js'
import '../assets/styles/admin-movies.css'

const {
  filteredLogs, loading, eventTypes,
  filterEvent, filterUserID, filterBookingID, filterDateFrom, filterDateTo
} = useAdminLogs()

function eventClass(event) {
  if (event === 'booking_success') return 'badge-showing'
  if (event === 'lock_failed' || event === 'booking_timeout') return 'badge-soon'
  return 'badge-genre'
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('th-TH', {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
    timeZone: 'Asia/Bangkok'
  })
}
</script>