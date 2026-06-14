import axios from 'axios'
import { auth } from '../firebase'

const api = axios.create({
  baseURL: 'http://localhost:8081',
})

api.interceptors.request.use(async (config) => {
  await auth.authStateReady() // ← รอ Firebase โหลดก่อน
  const user = auth.currentUser
  if (user) {
    const token = await user.getIdToken()
    config.headers.Authorization = `Bearer ${token}`
    localStorage.setItem('token', token)
  }
  return config
})

export default api