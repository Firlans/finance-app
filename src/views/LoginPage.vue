<script setup>
import BaseInput from '@/components/base/BaseInput.vue'
import NotificationFeature from '@/components/features/NotificationFeature.vue'
import { parseApiError } from '@/utils/Error.js'
import { Loading } from '@/utils/Loading.js'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  email: '',
  password: ''
})

const showPassword = ref(false)
const errors = reactive({})
const loading = new Loading()

const notif = reactive({ message: '', type: 'success' })

const handleLogin = async () => {
  // reset error
  errors.email = ''
  errors.password = ''

  // validasi minimal manual (opsional tapi waras)
  if (!form.email) {
    errors.email = 'Email is required'
    return
  }
  if (!form.password) {
    errors.password = 'Password is required'
    return
  }

  loading.start({ label: 'Logging in...' })

  try {
    const API_BASE = import.meta.env.VITE_BACKEND_SERVICE || 'http://localhost:8080/api'
    const res = await fetch(
      `${API_BASE}/users/login`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: form.email,
          password: form.password
        })
      }
    )

    if (!res.ok) {
      const errorMessage = await parseApiError(res)
      throw new Error(errorMessage)
    }

    let data = null
    const contentType = res.headers.get('content-type')
    if (contentType?.includes('application/json')) {
      data = await res.json()
    }

    const accessToken = data?.data?.access_token || data?.access_token
    if (!accessToken) {
      throw new Error('Login berhasil, tetapi token tidak diterima dari server.')
    }

    localStorage.setItem('access_token', accessToken)

    notif.message = data?.message || 'Login successful'
    notif.type = 'success'
    setTimeout(() => router.push('/dashboard'), 800)
  } catch (error) {
    notif.message = error.message || 'Login error'
    notif.type = 'error'
  } finally {
    loading.stop()
  }
}

const handleNotifClose = () => { notif.message = '' }
</script>


<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-100">
    <NotificationFeature :message="notif.message" :type="notif.type" @close="handleNotifClose" />
    <div class="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full">
      <h1 class="text-2xl font-bold text-center mb-6">Login</h1>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <BaseInput v-model="form.email" label="Email" type="email" placeholder="email@example.com" :error="errors.email"
          inputClass="bg-slate-50" />
        <BaseInput v-model="form.password" label="Password" :type="showPassword ? 'text' : 'password'"
          :error="errors.password">
          <template #right>
            <button type="button" class="text-sm text-slate-500" @click="showPassword = !showPassword">
              {{ showPassword ? 'Hide' : 'Show' }}
            </button>
          </template>
        </BaseInput>
        <button type="submit" class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
          Login
        </button>
      </form>
      <div class="flex justify-between items-center mt-4 text-sm text-slate-600">
        <router-link to="/forget-password" class="text-blue-600 hover:underline">Forgot password?</router-link>
        <span>
          Don't have an account?
          <router-link to="/register" class="text-blue-600 hover:underline">Register</router-link>
        </span>
      </div>
    </div>
  </div>
</template>
