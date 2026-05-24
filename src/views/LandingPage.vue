<script setup>
import { useRouter } from 'vue-router'
import SliderFeature from '@packages/components/base/SliderFeatures.vue'
import { onMounted, ref } from 'vue'
import { Loading } from '@packages/utils/Loading.js'

const router = useRouter()
const isMobileMenuOpen = ref(false)
const developers = ref([
  { name: 'firlans', asof: 'Frontend Developer' },
  { name: 'TubagusAldiMY', asof: 'backend Developer' }
])
const loading = new Loading()
onMounted(async () => {
  loading.start({ label: 'Loading developer profiles...' })

  try {
    const githubAPI = import.meta.env.VITE_GITHUB_API || 'https://api.github.com/users/'
    const res = await Promise.all(
      developers.value.map(async developer => {
        const response = await fetch(`${githubAPI}${developer.name}`)
        const data = await response.json()

        return {
          name: data.name || data.login,
          role: developer.asof,
          bio: data.bio,
          github: data.html_url,
          avatar: data.avatar_url,
        }
      })
    )

    contributors.value = res
  } catch (err) {
    console.warn(err)
  } finally {
    loading.stop()
  }
})
// Feature data
const features = [
  {
    icon: '📊',
    title: 'Pencatatan Pengeluaran',
    description: 'Catat setiap pengeluaranmu dengan mudah dan cepat. Kategorisasi otomatis membantu kamu memahami pola pengeluaran.'
  },
  {
    icon: '💰',
    title: 'Perencanaan Budget',
    description: 'Buat budget bulanan untuk berbagai kategori. Dapatkan notifikasi ketika mendekati batas budget.'
  },
  {
    icon: '📈',
    title: 'Laporan & Analisis',
    description: 'Lihat laporan keuangan lengkap dengan grafik dan analisis. Ketahui tren keuanganmu dari waktu ke waktu.'
  },
  {
    icon: '🎯',
    title: 'Target Tabungan',
    description: 'Tetapkan target tabungan dan lacak progresmu. Rayakan setiap pencapaian yang kamu buat.'
  }
]

// Contributors data (based on GitHub contributors)
const contributors = ref([])

const screenshots = [
  {
    id: 1,
    title: 'Dashboard Utama',
    description: 'Lihat ringkasan keuanganmu dalam satu tampilan',
    image: `${import.meta.env.BASE_URL}dashboard.png`
  },
  {
    id: 2,
    title: 'Pencatatan Transaksi',
    description: 'Catat pemasukan dan pengeluaran dengan mudah',
    image: `${import.meta.env.BASE_URL}transaksi.png`
  },
  {
    id: 3,
    title: 'Laporan Keuangan',
    description: 'Analisis keuangan dengan grafik yang intuitif',
    image: `${import.meta.env.BASE_URL}images/laporan.png`
  },
  {
    id: 4,
    title: 'Target Tabungan',
    description: 'Pantau progres pencapaian target tabunganmu',
    image: `${import.meta.env.BASE_URL}images/target.png`
  }
]

const navigateTo = (path) => {
  isMobileMenuOpen.value = false
  router.push(path)
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100 overflow-x-hidden">
    <!-- Navigation -->
    <nav class="fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-md shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center gap-2">
            <span class="text-2xl">💰</span>
            <span class="text-xl font-bold text-slate-800">FinanceTrack</span>
          </div>
          <!-- Desktop nav buttons -->
          <div class="hidden md:flex items-center gap-4">
            <button @click="navigateTo('/login')"
              class="px-4 py-2 text-slate-600 hover:text-slate-800 font-medium transition">
              Login
            </button>
            <button @click="navigateTo('/register')"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition">
              Register
            </button>
          </div>
          <!-- Mobile hamburger -->
          <button @click="isMobileMenuOpen = !isMobileMenuOpen"
            class="md:hidden inline-flex items-center justify-center p-2 rounded-lg text-slate-600 hover:text-slate-800 hover:bg-slate-100 transition"
            :aria-label="isMobileMenuOpen ? 'Tutup menu' : 'Buka menu'">
            <span class="text-2xl leading-none">{{ isMobileMenuOpen ? '×' : '☰' }}</span>
          </button>
        </div>
        <!-- Mobile dropdown menu -->
        <div v-if="isMobileMenuOpen" class="md:hidden border-t border-slate-200 py-3 space-y-2">
          <button @click="navigateTo('/login')"
            class="w-full text-left px-4 py-2 text-slate-700 hover:bg-slate-100 rounded-lg font-medium transition">
            Login
          </button>
          <button @click="navigateTo('/register')"
            class="w-full text-left px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition">
            Register
          </button>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="pt-24 pb-12 sm:pt-32 sm:pb-20 px-4">
      <div class="max-w-7xl mx-auto text-center">
        <h1 class="text-3xl sm:text-4xl md:text-6xl font-bold text-slate-800 mb-4 sm:mb-6 break-words">
          Kelola Keuanganmu dengan
          <span class="text-blue-600">Mudah</span> dan
          <span class="text-blue-600">Cerdas</span>
        </h1>
        <p class="text-base sm:text-xl text-slate-600 mb-6 sm:mb-8 max-w-2xl mx-auto">
          FinanceTrack membantu kamu mencatat, menganalisa, dan merencanakan keuangan
          dengan fitur-fitur pintar yang mudah digunakan.
        </p>
        <div class="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center">
          <button @click="navigateTo('/register')"
            class="px-6 sm:px-8 py-3 bg-blue-600 text-white text-base sm:text-lg rounded-xl hover:bg-blue-700 font-semibold transition shadow-lg hover:shadow-xl">
            Mulai Gratis
          </button>
          <button @click="navigateTo('/login')"
            class="px-6 sm:px-8 py-3 bg-white text-blue-600 text-base sm:text-lg rounded-xl hover:bg-slate-50 font-semibold transition border-2 border-blue-600">
            Login
          </button>
        </div>
      </div>
    </section>

    <!-- Screenshot Slideshow Section (Reusable Component) -->
    <SliderFeature title="Fitur Aplikasi" subtitle="Semua yang kamu butuhkan ada di sini" :subjects="screenshots"
      :autoPlayInterval="4000" />

    <!-- Features Section -->
    <section class="py-12 sm:py-20 px-4">
      <div class="max-w-7xl mx-auto">
        <h2 class="text-2xl sm:text-3xl font-bold text-center text-slate-800 mb-3 sm:mb-4">
          Fitur Utama
        </h2>
        <p class="text-sm sm:text-base text-center text-slate-600 mb-8 sm:mb-12">
          Berbagai fitur yang membantu kamu mengelola keuangan dengan lebih baik
        </p>

        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-5 sm:gap-8">
          <div v-for="feature in features" :key="feature.title"
            class="bg-white rounded-2xl p-5 sm:p-6 shadow-lg hover:shadow-xl transition-shadow">
            <div class="text-3xl sm:text-4xl mb-3 sm:mb-4">{{ feature.icon }}</div>
            <h3 class="text-base sm:text-xl font-semibold text-slate-800 mb-2">
              {{ feature.title }}
            </h3>
            <p class="text-sm text-slate-600">
              {{ feature.description }}
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Contributors Section -->
    <section class="py-12 sm:py-20 px-4 bg-white">
      <div class="max-w-7xl mx-auto">
        <h2 class="text-2xl sm:text-3xl font-bold text-center text-slate-800 mb-3 sm:mb-4">
          Tim Pengembang
        </h2>
        <p class="text-sm sm:text-base text-center text-slate-600 mb-8 sm:mb-12">
          Perkenalkan tim di balik FinanceTrack
        </p>

        <div class="flex flex-wrap justify-center gap-6 sm:gap-8">
          <div v-for="contributor in contributors" :key="contributor.name"
            class="bg-slate-50 rounded-2xl p-5 sm:p-6 text-center w-full max-w-xs hover:shadow-lg transition-shadow">
            <img :src="contributor.avatar" :alt="contributor.name"
              class="w-20 h-20 sm:w-24 sm:h-24 rounded-full mx-auto mb-4 border-4 border-white shadow-lg">
            <h3 class="text-base sm:text-xl font-semibold text-slate-800 mb-1">
              {{ contributor.name }}
            </h3>
            <p class="text-blue-600 text-sm font-medium mb-2">{{ contributor.role }}</p>
            <p class="text-slate-600 text-sm mb-4">
              {{ contributor.bio }}
            </p>
            <a :href="contributor.github" target="_blank" rel="noopener noreferrer"
              class="inline-flex items-center gap-2 text-slate-500 hover:text-blue-600 transition">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
              </svg>
              GitHub Profile
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-12 sm:py-20 px-4 bg-blue-600">
      <div class="max-w-4xl mx-auto text-center">
        <h2 class="text-2xl sm:text-3xl font-bold text-white mb-3 sm:mb-4">
          Siap Mengelola Keuanganmu?
        </h2>
        <p class="text-blue-100 text-sm sm:text-lg mb-6 sm:mb-8">
          Bergabung sekarang dan mulai perjalanan finansialmu yang lebih baik
        </p>
        <button @click="navigateTo('/register')"
          class="px-6 sm:px-8 py-3 bg-white text-blue-600 text-base sm:text-lg rounded-xl hover:bg-blue-50 font-semibold transition shadow-lg">
          Daftar Gratis Sekarang
        </button>
      </div>
    </section>

    <!-- Footer -->
    <footer class="py-8 px-4 bg-slate-800">
      <div class="max-w-7xl mx-auto text-center">
        <div class="flex items-center justify-center gap-2 mb-4">
          <span class="text-xl">💰</span>
          <span class="text-lg font-bold text-white">FinanceTrack</span>
        </div>
        <p class="text-slate-400 text-sm">
          © 2024 FinanceTrack. Built with ❤️ for better financial management.
        </p>
        <div class="mt-4 flex justify-center gap-4">
          <a href="#" class="text-slate-400 hover:text-white transition">
            Privacy Policy
          </a>
          <a href="#" class="text-slate-400 hover:text-white transition">
            Terms of Service
          </a>
          <a href="#" class="text-slate-400 hover:text-white transition">
            Contact
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped></style>
