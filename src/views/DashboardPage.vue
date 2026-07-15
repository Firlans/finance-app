<script setup>
import { ref, onMounted, shallowRef } from 'vue'
import Chart from 'chart.js/auto'
import dayjs from 'dayjs'
// Import DataService yang sudah ada di repo (apps/finance-app/src/DataService.js)
import { getTransactions } from '@/DataService.js'
import { ToggleFeature } from '@packages/components'

const transactionType = ref('credit') // Default: credit (Pengeluaran)
const selectedRange = ref('daily') // Default awal: daily

const timeRanges = [
  { label: 'Harian', value: 'daily' },
  { label: 'Mingguan', value: 'weekly' },
  { label: 'Bulanan', value: 'monthly' },
  { label: 'Tahunan', value: 'yearly' },
  { label: 'Semua', value: 'all' },
]

const transactionTypeOptions = [
  { label: 'Pengeluaran', value: 'credit' },
  { label: 'Pemasukan', value: 'debit' }
]

const chartCanvas = ref(null)
const chartInstance = shallowRef(null)

const chartTransactions = shallowRef([])
const chartLabels = ref([])
const selectedTransactions = ref(null)
const selectedLabel = ref('')

// 1. Load dari LocalStorage saat pertama kali dirender
onMounted(() => {
  const savedRange = localStorage.getItem('dashboard_time_range')
  if (savedRange && ['daily', 'weekly', 'monthly', 'yearly', 'all'].includes(savedRange)) {
    selectedRange.value = savedRange
  } else {
    selectedRange.value = 'daily'
  }

  fetchDataAndRender()
})

// 2. Fungsi merubah tipe transaksi (Credit/Debit)
const setTransactionType = (type) => {
  if (transactionType.value === type) return
  transactionType.value = type
  fetchDataAndRender()
}

// 3. Fungsi merubah dan menyimpan range (LocalStorage)
const setTimeRange = (range) => {
  if (selectedRange.value === range) return
  selectedRange.value = range
  localStorage.setItem('dashboard_time_range', range)
  fetchDataAndRender()
}

// 4. Fetch Data ke API sesuai Swagger
const fetchDataAndRender = async () => {
  const token = localStorage.getItem('access_token') || ''

  let toDate = ''
  let fromDate = ''

  if (selectedRange.value === 'daily') {
    toDate = dayjs().format('YYYY-MM-DD')
    fromDate = dayjs().subtract(6, 'day').format('YYYY-MM-DD')
  } else if (selectedRange.value === 'weekly') {
    toDate = dayjs().endOf('month').format('YYYY-MM-DD')
    fromDate = dayjs().startOf('month').format('YYYY-MM-DD')
  } else if (selectedRange.value === 'monthly') {
    toDate = dayjs().endOf('year').format('YYYY-MM-DD')
    fromDate = dayjs().startOf('year').format('YYYY-MM-DD')
  } else {
    // yearly atau all -> fetch seluruh data
    toDate = ''
    fromDate = ''
  }

  try {
    let rawTransactions = []
    if (fromDate && toDate) {
      rawTransactions = await getTransactions(token, fromDate, toDate)
    } else {
      rawTransactions = await getTransactions(token)
    }

    // 1. Filter data berdasarkan transaction_type ('credit' atau 'debit')
    const filteredData = rawTransactions.filter(
      (item) => item.transaction_type === transactionType.value
    )

    // 2. Tentukan jenis grouping
    let groupType = selectedRange.value
    if (groupType === 'all') {
      if (filteredData.length > 0) {
        const dates = filteredData.map(tx => dayjs(tx.transaction_date).valueOf())
        const minDate = dayjs(Math.min(...dates))
        const maxDate = dayjs(Math.max(...dates))
        const diffMonths = maxDate.diff(minDate, 'month', true)
        if (diffMonths > 12) groupType = 'yearly'
        else if (diffMonths > 1) groupType = 'monthly_all'
        else groupType = 'daily_all'
      } else {
        groupType = 'daily_all'
      }
    }

    // 3. Kelompokkan data
    const dateMap = {}
    const sortMap = {}
    const formatLabelMap = {}
    const transactionsGroupMap = {}

    const initKey = (k, sortVal, label) => {
      dateMap[k] = 0
      sortMap[k] = sortVal
      formatLabelMap[k] = label
      transactionsGroupMap[k] = []
    }

    // Inisialisasi default keys agar grafik konsisten (untuk range eksplisit)
    if (groupType === 'daily') {
      for (let i = 6; i >= 0; i--) {
        const d = dayjs().subtract(i, 'day')
        initKey(d.format('YYYY-MM-DD'), d.startOf('day').valueOf(), d.format('DD MMM'))
      }
    } else if (groupType === 'weekly') {
      ['Minggu 1', 'Minggu 2', 'Minggu 3', 'Minggu 4'].forEach((m, idx) => {
        initKey(m, idx, m)
      })
    } else if (groupType === 'monthly') {
      const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des']
      monthNames.forEach((m, idx) => {
        initKey(m, idx, m)
      })
    }

    filteredData.forEach((tx) => {
      const txDate = dayjs(tx.transaction_date)
      let key = ''
      let sortVal = 0
      let label = ''

      if (groupType === 'daily') {
        key = txDate.format('YYYY-MM-DD')
        sortVal = txDate.startOf('day').valueOf()
        label = txDate.format('DD MMM')
      } else if (groupType === 'weekly') {
        const date = txDate.date()
        if (date <= 7) { key = 'Minggu 1'; sortVal = 1; label = key }
        else if (date <= 14) { key = 'Minggu 2'; sortVal = 2; label = key }
        else if (date <= 21) { key = 'Minggu 3'; sortVal = 3; label = key }
        else { key = 'Minggu 4'; sortVal = 4; label = key }
      } else if (groupType === 'monthly') {
        const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des']
        key = monthNames[txDate.month()]
        sortVal = txDate.month()
        label = key
      } else if (groupType === 'yearly') {
        key = txDate.format('YYYY')
        sortVal = txDate.year()
        label = key
      } else if (groupType === 'monthly_all') {
        key = txDate.format('YYYY-MM')
        sortVal = txDate.startOf('month').valueOf()
        label = txDate.format('MMM YYYY')
      } else if (groupType === 'daily_all') {
        key = txDate.format('YYYY-MM-DD')
        sortVal = txDate.startOf('day').valueOf()
        label = txDate.format('DD MMM YYYY')
      }

      if (dateMap[key] === undefined) {
        initKey(key, sortVal, label)
      }
      dateMap[key] += tx.amount
      transactionsGroupMap[key].push(tx)
    })

    const sortedKeys = Object.keys(dateMap).sort((a, b) => sortMap[a] - sortMap[b])
    const finalLabels = sortedKeys.map(k => formatLabelMap[k])
    const dataValues = sortedKeys.map(k => dateMap[k])
    
    chartTransactions.value = sortedKeys.map(k => transactionsGroupMap[k])
    chartLabels.value = finalLabels
    selectedTransactions.value = null
    selectedLabel.value = ''

    // 5. Render ke Chart
    renderChart(finalLabels, dataValues)

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
  const labelText = isCredit ? 'Total Pengeluaran' : 'Total Pemasukan'

  // 3. Hancurkan (destroy) instance grafik lama jika sudah ada
  if (chartInstance.value) {
    chartInstance.value.destroy()
  }

  // 4. Custom plugin untuk garis silang (crosshair) dengan gaya putus-putus
  const crosshairPlugin = {
    id: 'crosshair',
    afterDraw: (chart) => {
      if (chart.tooltip?._active && chart.tooltip._active.length) {
        const activePoint = chart.tooltip._active[0]
        const ctx = chart.ctx
        const x = activePoint.element.x
        const y = activePoint.element.y
        const topY = chart.scales.y.top
        const bottomY = chart.scales.y.bottom
        const leftX = chart.scales.x.left
        const rightX = chart.scales.x.right

        ctx.save()
        ctx.beginPath()
        ctx.setLineDash([5, 5])
        ctx.lineWidth = 1
        ctx.strokeStyle = 'rgba(156, 163, 175, 0.8)' // abu-abu modern

        // Garis vertikal
        ctx.moveTo(x, topY)
        ctx.lineTo(x, bottomY)
        
        // Garis horizontal
        ctx.moveTo(leftX, y)
        ctx.lineTo(rightX, y)

        ctx.stroke()
        ctx.restore()
      }
    }
  }

  // 5. Buat grafik baru dan simpan di chartInstance
  chartInstance.value = new Chart(chartCanvas.value, {
    type: 'bar', // Anda bisa mengubahnya menjadi 'line' jika ingin melihat tren berupa garis
    data: {
      labels: labels, 
      datasets: [
        {
          label: labelText,
          data: dataValues, 
          backgroundColor: backgroundColor,
          borderColor: borderColor,
          borderWidth: 2,
          borderRadius: 4, 
        },
      ],
    },
    plugins: [crosshairPlugin],
    options: {
      onClick: (event, elements, chart) => {
        // Gunakan mode 'index' dengan intersect: false agar bisa diklik di manapun pada kolom (tidak harus pas di blok grafik)
        const activeElements = chart.getElementsAtEventForMode(event, 'index', { intersect: false }, true)
        if (activeElements && activeElements.length > 0) {
          const index = activeElements[0].index
          if (chartTransactions.value[index]) {
            selectedTransactions.value = chartTransactions.value[index]
            selectedLabel.value = chartLabels.value[index]
          }
        }
      },
      interaction: {
        mode: 'index',
        intersect: false,
      },
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: true,
          position: 'top',
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || ''
              if (label) {
                label += ': '
              }
              if (context.parsed.y !== null) {
                label += 'Rp ' + context.parsed.y.toLocaleString('id-ID')
              }
              return label
            }
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
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

    <!-- Transaction List Below Chart -->
    <div v-if="selectedTransactions !== null" class="mt-8 border-t border-gray-100 pt-6">
      <h3 class="text-lg font-semibold text-gray-800 mb-4">
        Transaksi: {{ selectedLabel }}
      </h3>
      <div v-if="selectedTransactions.length === 0" class="text-sm text-gray-500 bg-gray-50 p-4 rounded-lg text-center">
        Tidak ada transaksi pada periode ini.
      </div>
      <div v-else class="space-y-3">
        <div v-for="tx in selectedTransactions" :key="tx.id" class="flex justify-between items-center p-4 bg-gray-50 hover:bg-gray-100 transition-colors rounded-lg border border-gray-100">
          <div>
            <div class="font-medium text-gray-800">{{ tx.description || 'Tanpa Keterangan' }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ dayjs(tx.transaction_date).format('DD MMM YYYY') }}</div>
          </div>
          <div :class="tx.transaction_type === 'credit' ? 'text-red-600' : 'text-green-600'" class="font-semibold text-right">
            {{ tx.transaction_type === 'credit' ? '-' : '+' }}Rp {{ tx.amount.toLocaleString('id-ID') }}
            <div class="text-xs text-gray-500 mt-1 font-normal">{{ tx.account_name || 'Akun' }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
