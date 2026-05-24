<script setup>
import { reactive, ref } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Loading } from '@packages/utils/Loading.js'
import { Validator } from '@packages/utils/Validator.js'
import { required, email, minLength, sameAs } from '@packages/utils/Validator.js'
import { parseApiError } from '@/utils/Error.js'
import router from '@/router'

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})
const showPassword = ref(false)
const errors = reactive({})
const loading = new Loading()

const handleRegister = async () => {
  const API_BASE = import.meta.env.VITE_BACKEND_SERVICE || 'http://localhost:8080/api'
  console.log(`${API_BASE}/users/register`)
  const validator = new Validator(form, {
    name: [required()],
    email: [required(), email()],
    password: [minLength(8)],
    confirmPassword: [sameAs('password', 'Password not match')]
  })

  if (!validator.validate()) {
    Object.assign(errors, validator.getErrors())
    return
  }

  Object.keys(errors).forEach(k => delete errors[k])


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
    console.log('REGISTER SUCCESS', data)
    router.push('/login')
  } catch (error) {
    alert('gagal mendaftar: ' + (error?.message || error))
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
        <BaseInput v-model="form.name" label="Name" placeholder="Your name" :error="errors.name" />

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


        <BaseInput v-model="form.confirmPassword" label="Confirm Password" type="password"
          :error="errors.confirmPassword" />

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
