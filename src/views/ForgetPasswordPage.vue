<script setup>
import BaseInput from '@packages/components/base/BaseInput.vue'
import NotificationFeature from '@packages/components/features/NotificationFeature.vue'
import { Loading } from '@packages/utils/Loading.js'
import { parseApiError } from '@/utils/Error.js'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const form = reactive({
  email: ''
})
const errors = reactive({})
const notif = reactive({ message: '', type: 'success' })
const loading = new Loading()

const handleNotifClose = () => {
  notif.message = ''
}

const handleSubmit = async () => {
  errors.email = ''

  if (!form.email) {
    errors.email = 'Email is required'
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

    notif.message = 'Password reset link sent to your email'
    notif.type = 'success'

    setTimeout(() => {
      router.push('/login')
    }, 1200)
  } catch (error) {
    notif.message = error.message || 'Failed to send reset link'
    notif.type = 'error'
  } finally {
    loading.stop()
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 flex flex-col items-center px-4 py-10">
    <NotificationFeature :message="notif.message" :type="notif.type" @close="handleNotifClose" />
    <div class="my-auto w-full max-w-md bg-white p-6 sm:p-8 rounded-2xl shadow-lg">
      <h1 class="text-xl sm:text-2xl font-bold text-center mb-5 sm:mb-6">Forget Password</h1>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <BaseInput v-model="form.email" label="Email" type="email" placeholder="email@example.com" :error="errors.email"
          inputClass="bg-slate-50" />

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
