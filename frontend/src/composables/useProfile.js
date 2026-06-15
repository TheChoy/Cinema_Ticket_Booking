import { ref, reactive } from 'vue'
import api from '../services/api'

export function useProfile() {
  const user = ref(null)
  const loading = ref(false)
  const saving = ref(false)
  const error = ref('')
  const success = ref(false)

  const form = reactive({
    name: '',
    phone: '',
    date_of_birth: ''
  })

  async function fetchMe() {
    loading.value = true
    try {
      const res = await api.get('/users/me')
      user.value = res.data
      form.name = res.data.name || ''
      form.phone = res.data.phone || ''
      form.date_of_birth = res.data.date_of_birth
        ? res.data.date_of_birth.slice(0, 10)
        : ''
    } catch {
      error.value = 'โหลดข้อมูลไม่สำเร็จ'
    } finally {
      loading.value = false
    }
  }

async function saveProfile() {
  saving.value = true
  success.value = false
  error.value = ''
  try {
    const body = {
      name: form.name,
      phone: form.phone,
      date_of_birth: form.date_of_birth ? new Date(form.date_of_birth).toISOString() : null
    }
    await api.put('/users/me', body)
    await fetchMe()
    success.value = true
    setTimeout(() => success.value = false, 3000)
  } catch {
    error.value = 'บันทึกไม่สำเร็จ'
  } finally {
    saving.value = false
  }
}

  return { user, form, loading, saving, error, success, fetchMe, saveProfile }
}