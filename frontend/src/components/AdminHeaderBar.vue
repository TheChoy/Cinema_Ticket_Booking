<template>
  <header class="header admin-header">
    <router-link to="/admin/home" class="header-logo">🎬 CinemaTicket Admin</router-link>

    <nav class="admin-nav">
      <router-link to="/admin/movies">Movie</router-link>
      <router-link to="/admin/logs">log</router-link>
      <router-link to="/admin/bookings">การจอง</router-link>
      <router-link to="/admin/users">ผู้ใช้</router-link>
    </nav>

    <div class="header-actions" ref="dropdownRef">
      <button class="btn-profile" @click="toggleDropdown" aria-label="เมนูผู้ใช้">
        👤
      </button>

      <div v-if="open" class="dropdown">
        <div class="dropdown-item">{{ email }}</div>
        <div class="dropdown-divider" />
        <button class="dropdown-item danger" @click="handleLogout">
          ออกจากระบบ
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import '../assets/styles/header.css'

const router = useRouter()
const open = ref(false)
const dropdownRef = ref(null)
const email = ref(localStorage.getItem('email'))

function toggleDropdown() {
  open.value = !open.value
}

function handleLogout() {
  localStorage.clear()
  router.push('/login')
}

function handleClickOutside(e) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>