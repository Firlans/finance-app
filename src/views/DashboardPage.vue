<script setup>
import { ref, onMounted, shallowRef } from 'vue'
import Chart from 'chart.js/auto'
import dayjs from 'dayjs'
// Import DataService yang sudah ada di repo (apps/finance-app/src/DataService.js)
import { getTransactions } from '@/DataService.js'
import { ToggleFeature } from '@packages/components'

const transactionType = ref('credit') // Default: credit (Pengeluaran)
const selectedRange = ref(30) // Default awal: 30 hari

const timeRanges = [
  { label: '1 Hari', value: 1 },
  { label: '7 Hari', value: 7 },
  { label: '30 Hari', value: 30 },
  { label: '90 Hari', value: 90 },
]

const transactionTypeOptions = [
  { label: 'Pengeluaran', value: 'credit' },
  { label: 'Pemasukan', value: 'debit' }
]

const chartCanvas = ref(null)
const chartInstance = shallowRef(null)

// 1. Load dari LocalStorage saat pertama kali dirender
onMounted(() => {
  const savedRange = localStorage.getItem('dashboard_time_range')
  if (savedRange && !isNaN(parseInt(savedRange))) {
    selectedRange.value = parseInt(savedRange)
  }

  fetchDataAndRender()
})

// 2. Fungsi merubah tipe transaksi (Credit/Debit)
const setTransactionType = (type) => {
  if (transactionType.value === type) return
  transactionType.value = type
  fetchDataAndRender()
}

// 3. Fungsi merubah dan menyimpan range hari (LocalStorage)
const setTimeRange = (range) => {
  if (selectedRange.value === range) return
  selectedRange.value = range
  localStorage.setItem('dashboard_time_range', range.toString())
  fetchDataAndRender()
}

// 4. Fetch Data ke API sesuai Swagger
const fetchDataAndRender = async () => {
  const token = localStorage.getItem('access_token') || ''

  const toDate = dayjs().format('YYYY-MM-DD')
  const fromDate = dayjs().subtract(selectedRange.value, 'day').format('YYYY-MM-DD')

  try {
    const rawTransactions = await getTransactions(token, fromDate, toDate)

    // 1. Filter data berdasarkan transaction_type ('credit' atau 'debit')
    const filteredData = rawTransactions.filter(
      (item) => item.transaction_type === transactionType.value
    )

    // 2. Kelompokkan data berdasarkan Tanggal (transaction_date)
    const dateMap = {}

    filteredData.forEach((tx) => {
      // Ubah format tanggal API menjadi format yang enak dibaca (Contoh: "25 Jun 2026")
      const dateStr = dayjs(tx.transaction_date).format('DD MMM YYYY')

      // Jumlahkan nominal transaksi pada tanggal yang sama
      dateMap[dateStr] = (dateMap[dateStr] || 0) + tx.amount
    })

    // 3. Urutkan tanggal dari yang terlama ke terbaru (kronologis) agar grafik tidak mundur
    const sortedLabels = Object.keys(dateMap).sort((a, b) => {
      return dayjs(a, 'DD MMM YYYY').valueOf() - dayjs(b, 'DD MMM YYYY').valueOf()
    })

    // 4. Petakan total amount (dataValues) sesuai urutan tanggal yang sudah disortir
    const dataValues = sortedLabels.map(dateLabel => dateMap[dateLabel])

    // 5. Render ke Chart
    renderChart(sortedLabels, dataValues)

  } catch (error) {
    console.error('Gagal memuat data chart dashboard:', error)
  }
}

// 5. Render Grafik dengan Chart.js
const renderChart = (labels, dataValues) => {
  // 1. Validasi jika element canvas belum siap di DOM, batalkan proses
  if (!chartCanvas.value) return

  // 2. Konfigurasi warna dinamis berdasarkan tipe transaksi yang aktif
  const isCredit = transactionType.value === 'credit'
  const borderColor = isCredit ? 'rgb(239, 68, 68)' : 'rgb(34, 197, 94)'       // Merah (Credit) atau Hijau (Debit)
  const backgroundColor = isCredit ? 'rgba(239, 68, 68, 0.1)' : 'rgba(34, 197, 94, 0.1)'
  const labelText = isCredit ? 'Total Pengeluaran (Rp)' : 'Total Pemasukan (Rp)'

  // 3. Hancurkan (destroy) instance grafik lama jika sudah ada
  // Ini SANGAT PENTING dalam Vue agar grafik tidak tumpang tindih atau freeze saat ganti rentang hari
  if (chartInstance.value) {
    chartInstance.value.destroy()
  }

  // 4. Buat grafik baru dan simpan di chartInstance
  chartInstance.value = new Chart(chartCanvas.value, {
    type: 'bar', // Anda bisa mengubahnya menjadi 'line' jika ingin melihat tren berupa garis
    data: {
      labels: labels, // Array tanggal hasil format dayjs (misal: ['25 Jun 2026', '26 Jun 2026'])
      datasets: [
        {
          label: labelText,
          data: dataValues, // Array nominal (misal: [2601, 15000])
          backgroundColor: backgroundColor,
          borderColor: borderColor,
          borderWidth: 2,
          borderRadius: 4, // Membuat ujung bar sedikit tumpul agar modern
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: true,
          position: 'top',
        },
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            // Mengubah format angka di sumbu Y menjadi format mata uang Rupiah
            callback: function (value) {
              return 'Rp ' + value.toLocaleString('id-ID')
            },
          },
        },
      },
    },
  })
}
</script>

<template>
  <div class="p-6 bg-white rounded-lg shadow-sm">
    <div class="flex flex-col mb-6 space-y-4">
      <h2 class="text-xl font-semibold text-gray-800">Ringkasan Transaksi</h2>

      <div class="w-full">
        <ToggleFeature 
          :model-value="transactionType" 
          :options="transactionTypeOptions" 
          @change="setTransactionType" 
        />
      </div>
    </div>

    <div class="flex space-x-2 mb-6">
      <button v-for="range in timeRanges" :key="range.value" @click="setTimeRange(range.value)" :class="[
        'px-3 py-1.5 rounded-full text-xs font-medium transition-colors',
        selectedRange === range.value
          ? 'bg-blue-100 text-blue-700 border border-blue-200'
          : 'bg-gray-50 text-gray-600 border border-gray-200 hover:bg-gray-100'
      ]">
        {{ range.label }}
      </button>
    </div>

    <div class="relative w-full h-[300px]">
      <canvas ref="chartCanvas"></canvas>
    </div>
  </div>
</template>
