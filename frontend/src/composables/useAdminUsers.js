import { ref, computed, onMounted } from 'vue'
import api from '../services/api'
import { auth } from '../firebase'
import { onAuthStateChanged } from 'firebase/auth'

export function useAdminUsers() {
  const users = ref([])
  const loading = ref(false)
  const searchQuery = ref('')
  const filterRole = ref('')

  async function fetchUsers() {
    loading.value = true
    try {
      const res = await api.get('/admin/users')
      users.value = res.data || []
    } catch (err) {
      console.error('โหลด users ไม่สำเร็จ', err)
    } finally {
      loading.value = false
    }
  }

  const filteredUsers = computed(() => {
    const q = searchQuery.value.toLowerCase()
    return users.value
      .filter(u => !q ||
        u.name?.toLowerCase().includes(q) ||
        u.email?.toLowerCase().includes(q) ||
        u.phone?.toLowerCase().includes(q)
      )
      .filter(u => !filterRole.value || u.role === filterRole.value)
  })

  async function updateRole(id, role) {
    try {
      await api.put(`/admin/users/${id}/role`, { role })
      const u = users.value.find(x => x.id === id)
      if (u) u.role = role
      return true
    } catch (err) {
      console.error('อัพเดท role ไม่สำเร็จ', err)
      return false
    }
  }

  onMounted(() => {
    if (auth.currentUser) {
      fetchUsers()
    } else {
      const unsubscribe = onAuthStateChanged(auth, (user) => {
        unsubscribe()
        if (user) fetchUsers()
      })
    }
  })

  return {
    users, filteredUsers, loading,
    searchQuery, filterRole,
    fetchUsers, updateRole
  }
}