<template>
  <div class="admin-page">
    <AdminHeaderBar />

    <div class="admin-content">
      <div class="admin-top">
        <input v-model="searchQuery" placeholder="ค้นหาชื่อหนัง..." />
        <select v-model="filterGenre">
          <option value="">ทุก Genre</option>
          <option v-for="g in genres" :key="g" :value="g">{{ g }}</option>
        </select>
        <select v-model="filterStatus">
          <option value="">ทุก Status</option>
          <option value="now_showing">Now Showing</option>
          <option value="coming_soon">Coming Soon</option>
        </select>
        <button class="btn-add" @click="openCreate">+ เพิ่มหนัง</button>
      </div>

      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="!filteredMovies.length" class="state-box">ไม่พบหนัง</div>
      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th style="width:50px">โปสเตอร์</th>
              <th>ชื่อหนัง</th>
              <th>Genre</th>
              <th>ความยาว</th>
              <th>Status</th>
              <th style="width:180px">จัดการ</th>
            </tr>
          </thead>
          <tbody v-for="movie in filteredMovies" :key="movie.id">
            <!-- แถวหนัง -->
            <tr>
              <td>
                <div class="admin-movie-poster">
                  <img v-if="movie.poster_url" :src="movie.poster_url" :alt="movie.title" />
                  <span v-else>🎬</span>
                </div>
              </td>
              <td style="font-weight:500">{{ movie.title }}</td>
              <td><span class="badge badge-genre">{{ movie.genre }}</span></td>
              <td>{{ formatDuration(movie.duration) }}</td>
              <td>
                <span class="badge" :class="movie.status === 'now_showing' ? 'badge-showing' : 'badge-soon'">
                  {{ movie.status === 'now_showing' ? 'Now Showing' : 'Coming Soon' }}
                </span>
              </td>
              <td>
                <div class="actions">
                  <button class="btn-showtime" @click="toggleShowtime(movie)">
                    🕐 {{ expandedMovie === movie.id ? 'ปิด' : 'รอบฉาย' }}
                  </button>
                  <button class="btn-edit" @click="openEdit(movie)">✏️</button>
                  <button class="btn-delete" @click="confirmDelete(movie)">🗑</button>
                </div>
              </td>
            </tr>

            <!-- Accordion รอบฉาย -->
            <tr v-if="expandedMovie === movie.id">
              <td colspan="6" class="showtime-panel">
                <div class="showtime-header">
                  <span style="font-weight:600;font-size:0.875rem">รอบฉาย</span>
                  <button class="btn-add-sm" @click="openCreateShowtime(movie)">+ เพิ่มรอบ</button>
                </div>

                <div v-if="loadingShowtimes" style="padding:1rem;color:#9ca3af;font-size:0.85rem">⏳ กำลังโหลด...</div>
                <div v-else-if="!showtimes.length" style="padding:1rem;color:#9ca3af;font-size:0.85rem">ยังไม่มีรอบฉาย</div>
                <table v-else class="showtime-table">
                  <thead>
                    <tr>
                      <th>ห้อง</th>
                      <th>เริ่ม</th>
                      <th>จบ</th>
                      <th>ที่นั่ง</th>
                      <th>ราคา</th>
                      <th style="width:80px"></th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="s in showtimes" :key="s.id">
                      <td>{{ s.room }}</td>
                      <td>{{ formatDateTime(s.start_time) }}</td>
                      <td>{{ formatDateTime(s.end_time) }}</td>
                      <td>{{ s.seat_count }}</td>
                      <td>฿{{ s.price }}</td>
                      <td>
                        <div class="actions">
                          <button class="btn-edit" @click="openEditShowtime(s)">✏️</button>
                          <button class="btn-delete" @click="confirmDeleteShowtime(s)">🗑</button>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal หนัง -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editingMovie ? 'แก้ไขหนัง' : 'เพิ่มหนังใหม่' }}</h3>
          <button class="modal-close" @click="showModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>ชื่อหนัง *</label>
            <input v-model="form.title" placeholder="ชื่อหนัง" />
          </div>
          <div class="form-field">
            <label>Description</label>
            <textarea v-model="form.description" placeholder="รายละเอียดหนัง" />
          </div>
          <div class="form-field">
            <label>Genre</label>
            <select v-model="form.genre">
              <option value="">เลือก Genre</option>
              <option v-for="g in genres" :key="g" :value="g">{{ g }}</option>
            </select>
          </div>
          <div class="form-field">
            <label>ความยาว (นาที)</label>
            <input v-model.number="form.duration" type="number" placeholder="120" />
          </div>
          <div class="form-field">
            <label>Status</label>
            <select v-model="form.status">
              <option value="now_showing">Now Showing</option>
              <option value="coming_soon">Coming Soon</option>
            </select>
          </div>
          <div class="form-field">
            <label>โปสเตอร์</label>
            <input type="file" accept="image/*" @change="onFileChange" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showModal = false">ยกเลิก</button>
          <button class="btn-save" :disabled="saving || !form.title" @click="saveMovie">
            {{ saving ? 'กำลังบันทึก...' : 'บันทึก' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Modal รอบฉาย -->
    <div v-if="showShowtimeModal" class="modal-overlay" @click.self="showShowtimeModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editingShowtime ? 'แก้ไขรอบฉาย' : 'เพิ่มรอบฉาย' }}</h3>
          <button class="modal-close" @click="showShowtimeModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label>ห้อง *</label>
            <input v-model="showtimeForm.room" placeholder="Room 1" />
          </div>
          <div class="form-field">
            <label>เวลาเริ่ม *</label>
            <input v-model="showtimeForm.start_time" type="datetime-local" />
          </div>
          <div class="form-field">
            <label>เวลาจบ *</label>
            <input v-model="showtimeForm.end_time" type="datetime-local" />
          </div>
          <div class="form-field">
            <label>จำนวนที่นั่ง</label>
            <input v-model.number="showtimeForm.seat_count" type="number" placeholder="120" />
          </div>
          <div class="form-field">
            <label>ราคา (บาท)</label>
            <input v-model.number="showtimeForm.price" type="number" placeholder="250" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="showShowtimeModal = false">ยกเลิก</button>
          <button class="btn-save" :disabled="savingShowtime" @click="saveShowtime">
            {{ savingShowtime ? 'กำลังบันทึก...' : 'บันทึก' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import AdminHeaderBar from '../components/AdminHeaderBar.vue'
import { useAdminMovies } from '../composables/useAdminMovies.js'
import '../assets/styles/admin-movies.css'

const {
  filteredMovies, loading, saving,
  searchQuery, filterGenre, filterStatus, genres,
  fetchMovies, createMovie, updateMovie, deleteMovie, formatDuration,
  showtimes, loadingShowtimes, savingShowtime,
  expandedMovie, showShowtimeModal, editingShowtime, showtimeForm,
  toggleShowtime, openCreateShowtime, openEditShowtime,
  saveShowtime, confirmDeleteShowtime, formatDateTime
} = useAdminMovies()

const showModal = ref(false)
const editingMovie = ref(null)
const posterFile = ref(null)
const form = reactive({
  title: '', description: '', genre: '', duration: '', status: 'now_showing'
})

function openCreate() {
  editingMovie.value = null
  Object.assign(form, { title: '', description: '', genre: '', duration: '', status: 'now_showing' })
  posterFile.value = null
  showModal.value = true
}

function openEdit(movie) {
  editingMovie.value = movie
  Object.assign(form, {
    title: movie.title,
    description: movie.description,
    genre: movie.genre,
    duration: movie.duration,
    status: movie.status
  })
  posterFile.value = null
  showModal.value = true
}

function onFileChange(e) {
  posterFile.value = e.target.files[0] || null
}

async function saveMovie() {
  const formData = new FormData()
  formData.append('title', form.title)
  formData.append('description', form.description)
  formData.append('genre', form.genre)
  formData.append('duration', String(form.duration))
  formData.append('status', form.status)
  if (posterFile.value) formData.append('poster', posterFile.value)

  let ok
  if (editingMovie.value) {
    ok = await updateMovie(editingMovie.value.id, formData)
  } else {
    ok = await createMovie(formData)
  }
  if (ok) showModal.value = false
}

async function confirmDelete(movie) {
  if (!confirm(`ลบ "${movie.title}" ใช่ไหม?`)) return
  await deleteMovie(movie.id)
}

onMounted(() => fetchMovies())
</script>