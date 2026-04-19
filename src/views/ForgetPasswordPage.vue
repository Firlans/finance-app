<script setup>
import BaseInput from '@/components/base/BaseInput.vue'
import NotificationFeature from '@/components/features/NotificationFeature.vue'
import { parseApiError } from '@/utils/Error.js'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const form = reactive({
  email: ''
})
const errors = reactive({})
const notif = reactive({ message: '', type: 'success' })
const loading = ref(false)
const loadingLabel = ref('')

const handleNotifClose = () => {
  notif.message = ''
}

const handleSubmit = async () => {
  errors.email = ''

  if (!form.email) {
    errors.email = 'Email is required'
    return
  }

  loading.value = true
  loadingLabel.value = 'Sending reset link...'

  try {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_SERVICE}/users/forget-password`, {
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
    loading.value = false
    loadingLabel.value = ''
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-100">
    <NotificationFeature :message="notif.message" :type="notif.type" @close="handleNotifClose" />
    <LoadingFeature :show="loading" :label="loadingLabel" />
    <div class="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full">
      <h1 class="text-2xl font-bold text-center mb-6">Forget Password</h1>

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
