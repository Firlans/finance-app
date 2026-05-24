<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCurrentUser, logout } from '@/utils/DataService'
import { Notification } from '@packages/utils/Notification.js'

const router = useRouter()
const notification = new Notification()
const currentUser = ref(null)
const loading = ref(true)
const logoutLoading = ref(false)

onMounted(async () => {
    const token = localStorage.getItem('access_token')

    if (!token) {
        router.push('/login')
        return
    }

    try {
        const user = await getCurrentUser(token)
        currentUser.value = user
        if (!user) {
            notification.showError('Data pengguna tidak ditemukan. Silakan login kembali.')
            router.push('/login')
        }
    } catch (error) {
        notification.showError(error.message || 'Gagal memuat profil')
        router.push('/login')
    } finally {
        loading.value = false
    }
})

function formatDate(isoString) {
    if (!isoString) return '-'
    return new Date(isoString).toLocaleString()
}

async function handleLogout() {
    const token = localStorage.getItem('access_token')
    logoutLoading.value = true

    try {
        if (!token) {
            throw new Error('Token tidak ditemukan. Silakan login kembali.')
        }
        await logout(token)
        notification.showSuccess('Logout berhasil')
    } catch (error) {
        notification.showError(error.message || 'Gagal logout')
    } finally {
        localStorage.removeItem('access_token')
        logoutLoading.value = false
        router.push('/login')
    }
}
</script>

<template>
    <section class="profile-page p-6">
        <h1 class="text-2xl font-semibold mb-4">Profil</h1>

        <div v-if="loading" class="text-sm text-slate-600">Memuat profil...</div>

        <div v-else-if="currentUser" class="space-y-5 bg-white rounded-lg shadow-sm p-6">
            <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-2">
                    <div class="text-xs uppercase text-slate-500">ID</div>
                    <div class="text-sm text-slate-900">{{ currentUser.id }}</div>
                </div>
                <div class="space-y-2">
                    <div class="text-xs uppercase text-slate-500">Username</div>
                    <div class="text-sm text-slate-900">{{ currentUser.username }}</div>
                </div>
                <div class="space-y-2">
                    <div class="text-xs uppercase text-slate-500">Email</div>
                    <div class="text-sm text-slate-900">{{ currentUser.email }}</div>
                </div>
                <div class="space-y-2">
                    <div class="text-xs uppercase text-slate-500">Dibuat pada</div>
                    <div class="text-sm text-slate-900">{{ formatDate(currentUser.created_at) }}</div>
                </div>
            </div>

            <div class="border-t border-slate-200 pt-5">
                <p class="text-sm text-slate-500 mb-3">Gunakan tombol logout untuk mengakhiri sesi Anda dengan endpoint
                    backend <code>/api/users/logout</code>.</p>
                <button @click="handleLogout" :disabled="logoutLoading"
                    class="inline-flex items-center justify-center rounded-lg bg-red-600 px-5 py-3 text-white hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition">
                    {{ logoutLoading ? 'Logging out...' : 'Logout' }}
                </button>
            </div>
        </div>

        <div v-else class="text-red-600">Tidak ada data profil. Silakan login kembali.</div>
    </section>
</template>
