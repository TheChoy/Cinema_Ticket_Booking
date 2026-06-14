import { ref, reactive } from 'vue'
import { signInWithEmailAndPassword } from 'firebase/auth'
import { auth } from '../firebase'
import { useRouter } from 'vue-router'
import api from '../services/api'

export function Login() {
  const router = useRouter()
  const form = reactive({ email: '', password: '' })
  const showPassword = ref(false)
  const loading = ref(false)
  const errorMsg = ref('')

  async function handleLogin() {
    errorMsg.value = ''
    loading.value = true
    try {
      const userCredential = await signInWithEmailAndPassword(auth, form.email, form.password)
      const token = await userCredential.user.getIdToken()

      // ส่ง token ไป backend เพื่อ verify + ดึง user_id/role จาก MongoDB
      const res = await api.post('/auth/login', {}, {
        headers: { Authorization: `Bearer ${token}` }
      })

      const user = res.data // { _id, uid, email, role }

      // เก็บ token + user info
      localStorage.setItem('token', token)
      localStorage.setItem('user_id', user.id)
      localStorage.setItem('uid', user.uid)
      localStorage.setItem('role', user.role)
      localStorage.setItem('email', user.email)

      await router.push('/home')
    } catch (err) {
      console.error(err)
      errorMsg.value =
        err.code === 'auth/invalid-credential'
          ? 'Email หรือ Password ไม่ถูกต้อง'
          : 'เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง'
    } finally {
      loading.value = false
    }
  }

  return { form, showPassword, loading, errorMsg, handleLogin }
}