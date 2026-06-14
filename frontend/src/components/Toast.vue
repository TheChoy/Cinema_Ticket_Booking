<template>
  <transition name="toast">
    <div v-if="visible" class="toast">
      ✅ {{ message }}
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  message: { type: String, default: 'สำเร็จแล้ว' },
  duration: { type: Number, default: 3000 }
})

const visible = ref(false)

onMounted(() => {
  visible.value = true
  setTimeout(() => visible.value = false, props.duration)
})
</script>

<style scoped>
.toast {
  position: fixed;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  background: #111;
  color: #fff;
  padding: 12px 24px;
  border-radius: 999px;
  font-size: 0.9rem;
  font-weight: 500;
  z-index: 9999;
  white-space: nowrap;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
}

.toast-enter-active, .toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(20px);
}
</style>