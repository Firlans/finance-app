<script setup>
import BaseInput from '@packages/components/base/BaseInput.vue'
import { parseApiError } from '@packages/utils/Error.js'
import { Loading } from '@packages/utils/Loading.js'
import { Notification } from '@packages/utils/Notification.js'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  email: '',
  password: ''
})

const showPassword = ref(false)
const loading = new Loading()
const notification = new Notification()

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const validEmail = (value) => emailPattern.test(String(value).trim())

const handleLogin = async (event) => {
  event.preventDefault()
  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data login')
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

    notification.showSuccess(data?.message || 'Login berhasil')
    setTimeout(() => router.push('/dashboard'), 800)
  } catch (error) {
    notification.showError(error?.message || 'Login gagal')
  } finally {
    loading.stop()
  }
}
</script>


<template>
  <div class="min-h-screen bg-slate-100 flex flex-col items-center px-4 py-10">
    <div class="my-auto w-full max-w-md bg-white p-6 sm:p-8 rounded-2xl shadow-lg">
      <h1 class="text-xl sm:text-2xl font-bold text-center mb-5 sm:mb-6">Login</h1>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <BaseInput v-model="form.email" label="Email" type="email" placeholder="email@example.com"
          inputClass="bg-slate-50" required
          :validate="['Masukkan email yang valid', validEmail]" />
        <BaseInput v-model="form.password" label="Password" :type="showPassword ? 'text' : 'password'"
          required
          :validate="['Password wajib diisi', value => String(value).trim().length > 0]">
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
