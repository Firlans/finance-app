<script setup>
import { BaseInput } from '@packages/components'
import { Loading } from '@packages/utils/Loading.js'
import { parseApiError } from '@packages/utils/Error.js'
import { Notification } from '@packages/utils/Notification.js'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const form = reactive({
  email: ''
})
const loading = new Loading()
const notification = new Notification()

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const validEmail = (value) => emailPattern.test(String(value).trim())

const handleSubmit = async (event) => {
  event.preventDefault()
  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali email Anda')
    return
  }

  loading.start({ label: 'Sending reset link...' })

  try {
    const API_BASE = import.meta.env.VITE_BACKEND_SERVICE || 'http://localhost:8080/api'
    const res = await fetch(`${API_BASE}/users/forget-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: form.email })
    })

    if (!res.ok) {
      const errorMessage = await parseApiError(res)
      throw new Error(errorMessage)
    }

    notification.showSuccess('Link reset password berhasil dikirim ke email Anda')

    setTimeout(() => {
      router.push('/login')
    }, 1200)
  } catch (error) {
    notification.showError(error?.message || 'Gagal mengirim link reset password')
  } finally {
    loading.stop()
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 flex flex-col items-center px-4 py-10">
    <div class="my-auto w-full max-w-md bg-white p-6 sm:p-8 rounded-2xl shadow-lg">
      <h1 class="text-xl sm:text-2xl font-bold text-center mb-5 sm:mb-6">Forget Password</h1>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <BaseInput v-model="form.email" label="Email" type="email" placeholder="email@example.com" inputClass="bg-slate-50"
          required
          :validate="['Masukkan email yang valid', validEmail]" />

        <button type="submit" class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
          Send Reset Link
        </button>
      </form>

      <p class="text-center mt-4 text-slate-600">
        Remembered your password?
        <router-link to="/login" class="text-blue-600 hover:underline">Login</router-link>
      </p>
    </div>
  </div>
</template>
