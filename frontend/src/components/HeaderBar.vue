<template>
  <header class="header">
    <router-link to="/home" class="header-logo">🎬 CinemaTicket</router-link>

    <div class="header-actions" ref="dropdownRef">
      <button class="btn-profile" @click="toggleDropdown" aria-label="เมนูผู้ใช้">
        👤
      </button>

      <div v-if="open" class="dropdown">
        <router-link to="/profile" class="dropdown-item" @click="open = false">
          👤 โปรไฟล์
        </router-link>
        <router-link to="/history" class="dropdown-item" @click="open = false">
          🎟️ ประวัติการจอง
        </router-link>
        <div class="dropdown-divider" />
        <button class="dropdown-item danger" @click="handleLogout">
          🚪 ออกจากระบบ
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

function toggleDropdown() {
  open.value = !open.value
}

function handleLogout() {
  localStorage.clear()
  router.push('/login')
}

// ปิด dropdown เมื่อคลิกนอก
function handleClickOutside(e) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>