<template>
  <div class="login-page">
    <div class="login-card">
      <div class="brand">
        <h1>Welcome to Cinema Ticket</h1>
        <p>Sign in to your account</p>
      </div>

      <div class="form">
        <div v-if="errorMsg" class="error-box">{{ errorMsg }}</div>

        <div class="field">
          <label for="email">Email</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            placeholder="you@example.com"
            autocomplete="email"
          />
        </div>

        <div class="field">
          <label for="password">
            Password
            <a href="#" class="forgot">Forgot password?</a>
          </label>
          <div class="input-wrap">
            <input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="••••••••"
              autocomplete="current-password"
            />
            <button type="button" class="toggle-pw" @click="showPassword = !showPassword">
              {{ showPassword ? 'Hide' : 'Show' }}
            </button>
          </div>
        </div>

        <button class="btn-login" :disabled="loading" @click="handleLogin">
          {{ loading ? 'Signing in...' : 'Sign in' }}
        </button>
      </div>

      <p class="register-hint">
        Don't have an account? <a href="#">Create one</a>
      </p>
    </div>
  </div>
</template>

<script setup>
import { Login } from '../composables/Login.js'
import '../assets/styles/login.css'

const { form, showPassword, loading, errorMsg, handleLogin } = Login()
</script>