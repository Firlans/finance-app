<script setup>
import { ref, onMounted } from 'vue'
import { getCurrentUser } from '@/utils/DataService'
import { Notification } from '@/utils/Notification'

const notification = new Notification()
const currentUser = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const user = await preprocess()
    currentUser.value = user
  } catch (error) {
    console.error('Error during data preprocessing:', error)
  } finally {
    loading.value = false
  }
})

function preprocess() {
  const token = localStorage.getItem('access_token')
  return getCurrentUser(token).then(user => {
    if (!user) {
      throw new Error('No user data returned')
    }
    return user
  }).catch(error => {
    console.error('Error fetching current user:', error)
    notification.showError('Error fetching current user')
    return null
  })
}

function formatDate(isoString) {
  if (!isoString) return '-'
  return new Date(isoString).toLocaleString()
}
</script>

<template>
  <section class="dashboard-page p-6">
    <h1 class="text-2xl font-semibold mb-4">Dashboard</h1>

    <div v-if="loading" class="text-sm text-gray-600">Loading profile...</div>

    <div v-else-if="currentUser" class="space-y-3 bg-white rounded-lg shadow-sm p-5">
      <div class="text-lg font-medium">Profil Pengguna</div>
      <div class="grid gap-3 sm:grid-cols-2">
        <div>
          <div class="text-xs uppercase text-gray-500">ID</div>
          <div class="text-sm text-gray-900">{{ currentUser.id }}</div>
        </div>
        <div>
          <div class="text-xs uppercase text-gray-500">Username</div>
          <div class="text-sm text-gray-900">{{ currentUser.username }}</div>
        </div>
        <div>
          <div class="text-xs uppercase text-gray-500">Email</div>
          <div class="text-sm text-gray-900">{{ currentUser.email }}</div>
        </div>
        <div>
          <div class="text-xs uppercase text-gray-500">Dibuat pada</div>
          <div class="text-sm text-gray-900">{{ formatDate(currentUser.created_at) }}</div>
        </div>
      </div>
    </div>

    <div v-else class="text-red-600">Data profil tidak tersedia.</div>
  </section>
</template>
