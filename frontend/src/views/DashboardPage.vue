<script setup>
import router from '@/router'
import { computed, onMounted, ref, watch } from 'vue'
import dayjs from 'dayjs'
import SwitcherFeatures from '@/components/features/SwitcherFeatures.vue'
/* ====== STATE ====== */
const today = dayjs()
const currentMonth = ref(today.month())
const currentYear = ref(today.year())
const token = localStorage.getItem('access_token')
const URL = import.meta.env.VITE_BACKEND_SERVICE
// dummy data tampilan
const budgetBulanan = ref(0)
const currentBudgetId = ref(null)
const pengeluaranHarian = ref([]) // [{date: "2025-01-02", spent: 50000}, ...]

const editingBudget = ref(false)
const budgetInput = ref(budgetBulanan.value)

// Modal state for pengeluaran
const showPengeluaranModal = ref(false)
const editingPengeluaranId = ref(null)
const editingPengeluaranDate = ref('')
const editingPengeluaranAmount = ref(0)

/* ====== COMPUTED ====== */
watch([currentMonth, currentYear], () => {
  fetchBudgetBulanan()
})

watch(currentBudgetId, (newId) => {
  if (newId) {
    fetchPengeluaranHarian()
  } else {
    pengeluaranHarian.value = []
  }
})

onMounted(async () => {
  await fetchBudgetBulanan()
  if (currentBudgetId.value) {
    await fetchPengeluaranHarian()
  }
})

const daysInThisMonth = computed(() =>
  dayjs().year(currentYear.value).month(currentMonth.value).daysInMonth()
)

const budgetHarian = computed(() =>
  Math.floor(budgetBulanan.value / daysInThisMonth.value)
)

const daysInMonth = computed(() => {
  const days = []
  for (let d = 1; d <= daysInThisMonth.value; d++) {
    const date = dayjs().year(currentYear.value).month(currentMonth.value).date(d)
    const dateKey = date.format('YYYY-MM-DD')

    // Find matching expenditure for this date
    const pengeluaranEntry = pengeluaranHarian.value.find(p => {
      return dayjs(p.date).isSame(dayjs(dateKey), 'day')
    })
    const spent = pengeluaranEntry ? pengeluaranEntry.spent : 0
    const id = pengeluaranEntry ? pengeluaranEntry.id : null

    days.push({
      date: date.toDate(),
      day: d,
      weekday: date.format('ddd'),
      spent,
      id
    })
  }
  return days
})

const totalPengeluaran = computed(() =>
  daysInMonth.value.reduce((sum, d) => sum + d.spent, 0)
)

const monthLabel = computed(() =>
  dayjs().year(currentYear.value).month(currentMonth.value).format('MMMM YYYY')
)

/* ====== METHODS ====== */
function prevMonth() {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

function nextMonth() {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

function isToday(date) {
  return dayjs(date).isSame(dayjs(), 'day')
}

function format(value) {
  return value.toLocaleString('id-ID')
}

function startEditBudget() {
  budgetInput.value = budgetBulanan.value
  editingBudget.value = true
}

async function saveBudget() {
  try {
    if (!token) {
      router.push('/login')
      return
    }

    const payload = {
      budget: `${budgetInput.value}`,
      date: dayjs().year(currentYear.value).month(currentMonth.value).date(1).format('YYYY-MM-DD')
    }

    let res
    if (currentBudgetId.value) {
      // Update existing budget
      res = await fetch(`${URL}/budgets/${currentBudgetId.value}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      })
    } else {
      // Create new budget
      res = await fetch(`${URL}/budgets`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      })
    }

    if (!res.ok) {
      throw new Error('Gagal menyimpan budget')
    }

    // optional: refresh dari server agar state selalu valid
    await fetchBudgetBulanan()

    // refresh pengeluaran harian setelah budget diperbarui
    if (currentBudgetId.value) {
      await fetchPengeluaranHarian()
    }

    budgetBulanan.value = budgetInput.value
    editingBudget.value = false
  } catch (err) {
    console.error(err.message)
  }
}

function cancelEdit() {
  editingBudget.value = false
}

// Open modal to add/edit pengeluaran
function openPengeluaranModal(date, existingEntry = null) {
  const dateStr = dayjs(date).format('YYYY-MM-DD')
  editingPengeluaranDate.value = dateStr

  if (existingEntry) {
    editingPengeluaranId.value = existingEntry.id
    editingPengeluaranAmount.value = existingEntry.spent
  } else {
    editingPengeluaranId.value = null
    editingPengeluaranAmount.value = 0
  }

  showPengeluaranModal.value = true
}

// Close pengeluaran modal
function closePengeluaranModal() {
  showPengeluaranModal.value = false
  editingPengeluaranId.value = null
  editingPengeluaranDate.value = ''
  editingPengeluaranAmount.value = 0
}

// Save pengeluaran (create or update)
async function savePengeluaran() {
  try {
    if (!token || !currentBudgetId.value) {
      throw new Error('Budget ID tidak ditemukan')
    }

    const payload = {
      budget_id: currentBudgetId.value,
      amount: `${editingPengeluaranAmount.value}`,
      date: editingPengeluaranDate.value
    }

    let res
    if (editingPengeluaranId.value) {
      // Update existing pengeluaran
      res = await fetch(`${URL}/histories/${editingPengeluaranId.value}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      })
    } else {
      // Create new pengeluaran
      res = await fetch(`${URL}/histories`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      })
    }

    if (!res.ok) {
      throw new Error('Gagal menyimpan pengeluaran')
    }

    // Refresh pengeluaran list
    await fetchPengeluaranHarian()
    closePengeluaranModal()
  } catch (err) {
    console.error(err.message)
    alert(err.message)
  }
}

// Delete pengeluaran
async function deletePengeluaran(id) {
  try {
    if (!confirm('Yakin ingin menghapus pengeluaran ini?')) return

    if (!token) {
      throw new Error('Token tidak ditemukan')
    }

    const res = await fetch(`${URL}/histories/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      }
    })

    if (!res.ok) {
      throw new Error('Gagal menghapus pengeluaran')
    }

    // Refresh pengeluaran list
    await fetchPengeluaranHarian()
    closePengeluaranModal()
  } catch (err) {
    console.error(err.message)
    alert(err.message)
  }
}
const dateFrom = computed(() =>
  dayjs().year(currentYear.value).month(currentMonth.value).date(1).format('YYYY-MM-DD')
)

const dateTo = computed(() =>
  dayjs().year(currentYear.value).month(currentMonth.value).endOf('month').format('YYYY-MM-DD')
)

async function fetchBudgetBulanan() {
  try {
    if (!token) {
      console.warn('Token tidak ditemukan')
      budgetBulanan.value = 0
      throw new Error('Token tidak ditemukan')
    }

    const res = await fetch(`${URL}/budgets?start_date=${dateFrom.value}&end_date=${dateTo.value}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
    })

    if (!res.ok) throw new Error('Gagal mengambil budget')

    const { data } = await res.json()
    if (data && data.length > 0) {
      budgetBulanan.value = data[0].budget ? parseFloat(data[0].budget) : 0
      currentBudgetId.value = data[0].id
    } else {
      budgetBulanan.value = 0
      currentBudgetId.value = null
    }
  } catch (err) {
    console.error(err)
    router.push('/login')
  }
}

async function fetchPengeluaranHarian() {
  try {
    if (!token || !currentBudgetId.value) {
      pengeluaranHarian.value = []
      return
    }

    const res = await fetch(`${URL}/histories?budget_id=${currentBudgetId.value}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
    })

    if (!res.ok) throw new Error('Gagal mengambil pengeluaran')

    const { data } = await res.json()

    // Convert array to [{id, date, spent}] format
    const pengeluaran = []
    if (data && Array.isArray(data)) {
      data.forEach(item => {
        const dateKey = dayjs(item.date).format('YYYY-MM-DD')
        pengeluaran.push({
          id: item.id,
          date: dateKey,
          spent: parseFloat(item.amount) || 0
        })
      })
    }
    pengeluaranHarian.value = pengeluaran
  } catch (err) {
    console.error(err)
    pengeluaranHarian.value = []
  }
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100">
    <!-- Top Bar -->
    <header class="sticky top-0 z-40 bg-white/80 backdrop-blur-md border-b">
      <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-2xl">💰</span>
          <span class="text-xl font-bold text-slate-800">FinanceTrack</span>
        </div>

        <!-- Month Switcher -->
        <SwitcherFeatures :label="monthLabel" @prev="prevMonth" @next="nextMonth">Periode</SwitcherFeatures>
      </div>
    </header>

    <!-- Summary -->
    <section class="max-w-7xl mx-auto px-4 py-8">
      <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="bg-white rounded-2xl p-5 shadow-sm">
          <p class="text-sm text-slate-500 mb-1">Budget Bulanan</p>

          <div v-if="!editingBudget" class="flex items-center justify-between">
            <p class="text-2xl font-bold text-slate-800">
              Rp {{ format(budgetBulanan) }}
            </p>
            <button @click="startEditBudget" class="text-sm text-blue-600 hover:underline">
              Edit
            </button>
          </div>

          <form v-else class="flex items-center gap-2" @submit.prevent="saveBudget">
            <input v-model.number="budgetInput" type="number"
              class="w-full border rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500" />

            <button type="submit" @click="saveBudget" class="px-3 py-2 bg-blue-600 text-white rounded-lg">
              ✔
            </button>

            <button type="button" @click="cancelEdit" class="px-3 py-2 border rounded-lg">
              ✕
            </button>
          </form>

        </div>


        <div class="bg-white rounded-2xl p-5 shadow-sm">
          <p class="text-sm text-slate-500 mb-1">Budget Harian</p>
          <p class="text-2xl font-bold text-blue-600">
            Rp {{ format(budgetHarian) }}
          </p>
        </div>

        <div class="bg-white rounded-2xl p-5 shadow-sm">
          <p class="text-sm text-slate-500 mb-1">Total Pengeluaran</p>
          <p class="text-2xl font-bold text-red-500">
            Rp {{ format(totalPengeluaran) }}
          </p>
        </div>

        <div class="bg-white rounded-2xl p-5 shadow-sm">
          <p class="text-sm text-slate-500 mb-1">Sisa Budget</p>
          <p class="text-2xl font-bold"
            :class="(budgetBulanan - totalPengeluaran) >= 0 ? 'text-emerald-600' : 'text-red-600'">
            Rp {{ format(budgetBulanan - totalPengeluaran) }}
          </p>
        </div>
      </div>
    </section>

    <!-- Daily Budget Grid -->
    <section class="max-w-7xl mx-auto px-4 pb-16">
      <h2 class="text-xl font-bold text-slate-800 mb-4">
        Pengeluaran Harian
      </h2>

      <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div v-for="day in daysInMonth" :key="day.date"
          class="bg-white rounded-2xl p-4 shadow-sm hover:shadow-md transition cursor-pointer"
          :class="isToday(day.date) ? 'ring-2 ring-blue-500' : ''"
          @click="openPengeluaranModal(day.date, day.id ? { id: day.id, spent: day.spent } : null)">
          <div class="flex items-center justify-between mb-3">
            <div>
              <p class="text-sm text-slate-500">
                {{ day.weekday }}
              </p>
              <p class="text-lg font-semibold text-slate-800">
                {{ day.day }}
              </p>
            </div>

            <span class="text-xs px-2 py-1 rounded-full" :class="day.spent > budgetHarian
              ? 'bg-red-100 text-red-600'
              : 'bg-emerald-100 text-emerald-600'">
              {{ day.spent > budgetHarian ? 'Over' : 'On Track' }}
            </span>
          </div>

          <div class="space-y-1 text-sm">
            <div class="flex justify-between">
              <span class="text-slate-500">Budget</span>
              <span class="font-medium">
                Rp {{ format(budgetHarian) }}
              </span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-500">Pengeluaran</span>
              <span :class="day.spent > budgetHarian ? 'text-red-500' : 'text-slate-700'">
                Rp {{ format(day.spent) }}
              </span>
            </div>
            <div class="flex justify-between border-t pt-1">
              <span class="text-slate-500">Sisa</span>
              <span class="font-semibold" :class="(budgetHarian - day.spent) >= 0
                ? 'text-emerald-600'
                : 'text-red-600'">
                Rp {{ format(budgetHarian - day.spent) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>

  <!-- Pengeluaran Modal -->
  <div v-if="showPengeluaranModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
    <div class="bg-white rounded-2xl p-6 w-full max-w-md mx-4">
      <h3 class="text-xl font-bold text-slate-800 mb-4">
        {{ editingPengeluaranId ? 'Edit Pengeluaran' : 'Tambah Pengeluaran' }}
      </h3>

      <div class="mb-4">
        <label class="block text-sm text-slate-600 mb-1">Tanggal</label>
        <input type="date" v-model="editingPengeluaranDate" readonly
          class="w-full border rounded-lg px-3 py-2 bg-slate-100 cursor-not-allowed" />
      </div>

      <div class="mb-6">
        <label class="block text-sm text-slate-600 mb-1">Jumlah Pengeluaran (Rp)</label>
        <input v-model.number="editingPengeluaranAmount" type="number"
          class="w-full border rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500" placeholder="0" />
      </div>

      <div class="flex gap-2">
        <button @click="savePengeluaran" class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          Simpan
        </button>
        <button v-if="editingPengeluaranId" @click="deletePengeluaran(editingPengeluaranId)"
          class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600">
          Hapus
        </button>
        <button @click="closePengeluaranModal" class="px-4 py-2 border rounded-lg hover:bg-slate-100">
          Batal
        </button>
      </div>
    </div>
  </div>
</template>
