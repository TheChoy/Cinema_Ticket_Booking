<template>
  <div class="profile-page">
    <HeaderBar v-if="role !== 'admin'" />
    <AdminHeaderBar v-else />

    <div class="profile-content">
      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>

      <div v-else-if="user" class="profile-card">

        <div class="profile-header">
        <div class="profile-avatar">👤</div>
        <div>{{ user.name || 'ไม่ระบุชื่อ' }}</div>
        <div class="profile-email">{{ user.email }}</div>
        <span v-if="role === 'admin'" class="profile-role">{{ user.role }}</span>
        </div>

        <div class="profile-body">
          <div v-if="success" class="success-box">✅ บันทึกสำเร็จแล้ว</div>
          <div v-if="error" class="error-box">⚠️ {{ error }}</div>

          <div class="form-field">
            <label>Email</label>
            <input :value="user.email" disabled />
          </div>
          <div class="form-field">
            <label>ชื่อ</label>
            <input v-model="form.name" placeholder="ชื่อ-นามสกุล" />
          </div>
          <div class="form-field">
            <label>เบอร์โทร</label>
            <input v-model="form.phone" placeholder="0812345678" />
          </div>
          <div class="form-field">
            <label>วันเกิด</label>
            <input v-model="form.date_of_birth" type="date" />
          </div>
        </div>

        <div class="profile-footer">
          <button class="btn-save" :disabled="saving" @click="saveProfile">
            {{ saving ? 'กำลังบันทึก...' : 'บันทึก' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import HeaderBar from '../components/HeaderBar.vue'
import AdminHeaderBar from '../components/AdminHeaderBar.vue'
import { useProfile } from '../composables/useProfile.js'
import '../assets/styles/profile.css'

const role = ref(localStorage.getItem('role'))
const { user, form, loading, saving, error, success, fetchMe, saveProfile } = useProfile()

onMounted(() => fetchMe())
</script>