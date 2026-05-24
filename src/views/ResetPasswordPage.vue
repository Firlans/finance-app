<script setup>
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Loading } from '@packages/utils/Loading.js'
import { parseApiError } from '@packages/utils/Error.js'
import { Notification } from '@packages/utils/Notification.js'
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
const loading = new Loading()
const showPassword = ref(false)
const notification = new Notification()

const validPassword = (value) => String(value).length >= 8
const validPasswordConfirmation = (value) => String(value).length > 0 && value === form.new_password

const handleSubmit = async (event) => {
  event.preventDefault()
  if (!token.value) {
    notification.showError('Link reset tidak valid. Silakan gunakan link dari email.')
    return
  }

  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data password baru')
    return
  }

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

    notification.showSuccess('Password berhasil direset. Silakan login kembali.')

    setTimeout(() => {
      router.push('/login')
    }, 1200)
  } catch (error) {
    notification.showError(error?.message || 'Reset password gagal')
  } finally {
    loading.stop()
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 flex flex-col items-center px-4 py-10">
    <div class="my-auto w-full max-w-md bg-white p-6 sm:p-8 rounded-2xl shadow-lg">
      <h1 class="text-xl sm:text-2xl font-bold text-center mb-5 sm:mb-6">Reset Password</h1>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <BaseInput v-model="form.new_password" label="New Password" :type="showPassword ? 'text' : 'password'"
          inputClass="bg-slate-50" required
          :validate="['Password minimal 8 karakter', validPassword]">
          <template #right>
            <button type="button" class="text-sm text-slate-500" @click="showPassword = !showPassword">
              {{ showPassword ? 'Hide' : 'Show' }}
            </button>
          </template>
        </BaseInput>

        <BaseInput v-model="form.confirm_password" label="Confirm Password" :type="showPassword ? 'text' : 'password'"
          inputClass="bg-slate-50" required
          :validate="['Konfirmasi password harus sama', validPasswordConfirmation]" />

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
