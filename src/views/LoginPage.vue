<script setup>
import BaseInput from '@/components/base/BaseInput.vue'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  email: '',
  password: ''
})

const showPassword = ref(false)
const errors = reactive({})

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

  try {
    const res = await fetch(
      `${import.meta.env.VITE_BACKEND_SERVICE}/users/login`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: form.email,
          password: form.password
        })
      }
    )

    // backend menolak
    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg || 'Login failed')
    }

    // backend sukses (JSON atau tidak)
    let data = null
    const contentType = res.headers.get('content-type')
    if (contentType?.includes('application/json')) {
      data = await res.json()
    }

    console.log('LOGIN SUCCESS', data)
    localStorage.setItem('access_token', data.data.access_token)
    // redirect
    router.push('/dashboard')
  } catch (error) {
    alert(error.message || 'Login error')
  }
}
</script>


<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-100">
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
      <p class="text-center mt-4 text-slate-600">
        Don't have an account?
        <router-link to="/register" class="text-blue-600 hover:underline">Register</router-link>
      </p>
    </div>
  </div>
</template>
