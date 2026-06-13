import { ref, reactive } from 'vue'
import { signInWithEmailAndPassword } from 'firebase/auth'
import { auth } from '../firebase'
import { useRouter } from 'vue-router'

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
      console.log('Firebase token:', token)
      await router.push('/home')
    } catch (err) {
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