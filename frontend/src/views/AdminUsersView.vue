<template>
  <div class="admin-page">
    <AdminHeaderBar />

    <div class="admin-content">
      <div class="admin-top">
        <input v-model="filterEmail" placeholder="ค้นหา..." />
        <select v-model="filterRole">
          <option value="">ทุก Role</option>
          <option value="user">User</option>
          <option value="admin">Admin</option>
        </select>
      </div>

      <div v-if="loading" class="state-box">⏳ กำลังโหลด...</div>
      <div v-else-if="!filteredUsers.length" class="state-box">ไม่พบผู้ใช้</div>
      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>ชื่อ</th>
              <th>Email</th>
              <th>เบอร์โทร</th>
              <th>Role</th>
              <th style="width:140px">จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in filteredUsers" :key="u.id">
              <td>{{ u.name || '-' }}</td>
              <td>{{ u.email }}</td>
              <td>{{ u.phone || '-' }}</td>
              <td>
                <span class="badge" :class="u.role === 'admin' ? 'badge-showing' : 'badge-genre'">
                  {{ u.role }}
                </span>
              </td>
              <td>
                <select :value="u.role" @change="onRoleChange(u, $event.target.value)">
                  <option value="user">User</option>
                  <option value="admin">Admin</option>
                </select>
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
import { useAdminUsers } from '../composables/useAdminUsers.js'
import '../assets/styles/admin-movies.css'

const { filteredUsers, loading, filterEmail, filterRole, updateRole } = useAdminUsers()

async function onRoleChange(u, newRole) {
  if (newRole === u.role) return
  if (!confirm(`เปลี่ยน role ของ ${u.email} เป็น ${newRole}?`)) return
  await updateRole(u.id, newRole)
}
</script>