<template>
  <div>
    <h1>Hello Admin</h1>
    <p>UID: {{ uid }}</p>
    <p>Email: {{ email }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const uid = ref(localStorage.getItem('uid'))
const email = ref(localStorage.getItem('email'))
</script>   