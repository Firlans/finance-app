<script setup>
import { reactive, ref } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Loading } from '@packages/utils/Loading.js'
import { parseApiError } from '@packages/utils/Error.js'
import { Notification } from '@packages/utils/Notification.js'
import router from '@/router'

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})
const showPassword = ref(false)
const loading = new Loading()
const notification = new Notification()

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const validName = (value) => String(value).trim().length > 0
const validEmail = (value) => emailPattern.test(String(value).trim())
const validPassword = (value) => String(value).length >= 8
const validPasswordConfirmation = (value) => String(value).length > 0 && value === form.password

const handleRegister = async (event) => {
  event.preventDefault()
  const API_BASE = import.meta.env.VITE_BACKEND_SERVICE || 'http://localhost:8080/api'
  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data pendaftaran')
    return
  }

  loading.start({ label: 'Registering account...' })

  try {

    const res = await fetch(`${API_BASE}/users/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: form.name,
        email: form.email,
        password: form.password
      })
    })


    if (!res.ok) {
      const errorMessage = await parseApiError(res)
      throw new Error(errorMessage)
    }

    const data = await res.json()
    notification.showSuccess(data?.message || 'Registrasi berhasil')
    router.push('/login')
  } catch (error) {
    notification.showError(error?.message || 'Gagal mendaftar')
  } finally {
    loading.stop()
  }
}

</script>

<template>
  <div class="min-h-screen bg-slate-100 flex flex-col items-center px-4 py-10">
    <div class="my-auto w-full max-w-md bg-white p-6 sm:p-8 rounded-2xl shadow-lg">
      <h1 class="text-xl sm:text-2xl font-bold text-center mb-5 sm:mb-6">Register</h1>
      <form @submit.prevent="handleRegister" class="space-y-4">
        <BaseInput v-model="form.name" label="Name" placeholder="Your name" required
          :validate="['Nama wajib diisi', validName]" />

        <BaseInput v-model="form.email" label="Email" type="email" placeholder="email@example.com" inputClass="bg-slate-50"
          required
          :validate="['Masukkan email yang valid', validEmail]" />

        <BaseInput v-model="form.password" label="Password" :type="showPassword ? 'text' : 'password'"
          required
          :validate="['Password minimal 8 karakter', validPassword]">
          <template #right>
            <button type="button" class="text-sm text-slate-500" @click="showPassword = !showPassword">
              {{ showPassword ? 'Hide' : 'Show' }}
            </button>
          </template>
        </BaseInput>


        <BaseInput v-model="form.confirmPassword" label="Confirm Password" type="password"
          required
          :validate="['Konfirmasi password harus sama', validPasswordConfirmation]" />

        <button type="submit" class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
          Register
        </button>
      </form>
      <p class="text-center mt-4 text-slate-600">
        Already have an account?
        <router-link to="/login" class="text-blue-600 hover:underline">Login</router-link>
      </p>
    </div>
  </div>
</template>
