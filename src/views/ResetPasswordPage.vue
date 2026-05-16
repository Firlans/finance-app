<script setup>
import BaseInput from '@/components/base/BaseInput.vue'
import NotificationFeature from '@/components/features/NotificationFeature.vue'
import { Loading } from '@/utils/Loading.js'
import { parseApiError } from '@/utils/Error.js'
import { Validator, required, minLength, sameAs } from '@/utils/Validator'
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const token = ref('')
if (route.query.token) {
  token.value = Array.isArray(route.query.token) ? route.query.token[0] : route.query.token
}

const form = reactive({
  new_password: '',
  confirm_password: ''
})
const errors = reactive({})
const notif = reactive({ message: '', type: 'success' })
const loading = new Loading()
const showPassword = ref(false)

const handleNotifClose = () => {
  notif.message = ''
}

const handleSubmit = async () => {
  errors.new_password = ''
  errors.confirm_password = ''

  if (!token.value) {
    errors.token = 'Invalid reset link. Please use the link from the email.'
    return
  }

  const validator = new Validator(form, {
    new_password: [required('New password is required'), minLength(8, 'Password must be at least 8 characters')],
    confirm_password: [sameAs('new_password', 'Password confirmation does not match')]
  })

  if (!validator.validate()) {
    Object.assign(errors, validator.getErrors())
    return
  }

  Object.keys(errors).forEach((key) => delete errors[key])

  loading.start({ label: 'Resetting password...' })

  try {
    const API_BASE = import.meta.env.VITE_BACKEND_SERVICE || 'http://localhost:8080/api'
    const res = await fetch(`${API_BASE}/users/reset-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        token: token.value,
        new_password: form.new_password,
        confirm_password: form.confirm_password
      })
    })

    if (!res.ok) {
      const errorMessage = await parseApiError(res)
      throw new Error(errorMessage)
    }

    notif.message = 'Password berhasil direset. Silakan login kembali.'
    notif.type = 'success'

    setTimeout(() => {
      router.push('/login')
    }, 1200)
  } catch (error) {
    notif.message = error.message || 'Reset password gagal'
    notif.type = 'error'
  } finally {
    loading.stop()
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-100">
    <NotificationFeature :message="notif.message" :type="notif.type" @close="handleNotifClose" />
    <div class="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full">
      <h1 class="text-2xl font-bold text-center mb-6">Reset Password</h1>

      <p v-if="errors.token" class="text-sm text-red-600 mb-4">{{ errors.token }}</p>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <BaseInput v-model="form.new_password" label="New Password" :type="showPassword ? 'text' : 'password'"
          :error="errors.new_password" inputClass="bg-slate-50">
          <template #right>
            <button type="button" class="text-sm text-slate-500" @click="showPassword = !showPassword">
              {{ showPassword ? 'Hide' : 'Show' }}
            </button>
          </template>
        </BaseInput>

        <BaseInput v-model="form.confirm_password" label="Confirm Password" :type="showPassword ? 'text' : 'password'"
          :error="errors.confirm_password" inputClass="bg-slate-50" />

        <button type="submit"
          class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-60"
          :disabled="!token">
          Reset Password
        </button>
      </form>

      <p class="text-center mt-4 text-slate-600">
        Remembered your password?
        <router-link to="/login" class="text-blue-600 hover:underline">Login</router-link>
      </p>
    </div>
  </div>
</template>
