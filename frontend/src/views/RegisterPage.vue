<script setup>
import { reactive, ref } from 'vue'
import BaseInput from '@/components/base/BaseInput.vue'
import { Validator } from '@/utils/Validator'
import { required, email, minLength, sameAs } from '@/utils/Validator'
import router from '@/router'

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})
const showPassword = ref(false)
const errors = reactive({})

const handleRegister = async () => {
  console.log(`${import.meta.env.VITE_BACKEND_SERVICE}/users/register`)
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


  try {

    const res = await fetch(`${import.meta.env.VITE_BACKEND_SERVICE}/users/register`, {
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
      // error dari server (contoh: email sudah dipakai)
      throw new Error(res.message || 'Register failed')
    }

    const data = await res.json()
    console.log('REGISTER SUCCESS', data)
    router.push('/login')
  } catch (error) {
    alert('gagal mendaftar' + error)
  }
}

</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-100">
    <div class="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full">
      <h1 class="text-2xl font-bold text-center mb-6">Register</h1>
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
