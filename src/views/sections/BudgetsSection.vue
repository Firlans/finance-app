<script setup>
import { reactive, ref, computed, onMounted, nextTick } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Loading } from '@packages/utils/Loading.js'
import { Notification } from '@packages/utils/Notification.js'
import { getBudgets, upsertBudget, deleteBudget, getCategories } from '@/DataService.js'
import dayjs from 'dayjs'

const loading = new Loading()
const notification = new Notification()
const token = localStorage.getItem('access_token')

const budgets = ref([])
const categories = ref([])
const isFormOpen = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({ 
  id: null,
  name: '', 
  amount: '', 
  interval_type: 'monthly',
  day: 1,
  date: 1,
  month: 1,
  repeat: true,
  category_ids: []
})

const intervalOptions = [
  { label: 'Mingguan', value: 'weekly' },
  { label: 'Bulanan', value: 'monthly' },
  { label: 'Tahunan', value: 'yearly' },
]

const dayOptions = [
  { label: 'Senin', value: 1 },
  { label: 'Selasa', value: 2 },
  { label: 'Rabu', value: 3 },
  { label: 'Kamis', value: 4 },
  { label: 'Jumat', value: 5 },
  { label: 'Sabtu', value: 6 },
  { label: 'Minggu', value: 7 },
]

const monthOptions = [
  { label: 'Januari', value: 1 },
  { label: 'Februari', value: 2 },
  { label: 'Maret', value: 3 },
  { label: 'April', value: 4 },
  { label: 'Mei', value: 5 },
  { label: 'Juni', value: 6 },
  { label: 'Juli', value: 7 },
  { label: 'Agustus', value: 8 },
  { label: 'September', value: 9 },
  { label: 'Oktober', value: 10 },
  { label: 'November', value: 11 },
  { label: 'Desember', value: 12 },
]

const formTitle = computed(() => (editingId.value ? 'Edit Anggaran' : 'Tambah Anggaran'))
const submitLabel = computed(() => (editingId.value ? 'Simpan Perubahan' : 'Tambah Anggaran'))

const formatMoney = (value) => {
  if (value == null) return '0'
  return value.toLocaleString('id-ID')
}

const formatDate = (value) => {
  if (!value) return '-'
  return dayjs(value).format('DD MMM YYYY')
}

const getCategoryNames = (ids) => {
  if (!ids || ids.length === 0) return 'Semua Kategori'
  return ids.map(id => categories.value.find(c => c.id === id)?.name).filter(Boolean).join(', ')
}

const formatIntervalInfo = (budget) => {
  if (budget.interval_type === 'weekly') {
    const dayName = dayOptions.find(d => d.value === budget.day)?.label || ''
    return `Setiap Hari ${dayName}`
  }
  if (budget.interval_type === 'monthly') {
    return `Setiap Tanggal ${budget.date}`
  }
  if (budget.interval_type === 'yearly') {
    const monthName = monthOptions.find(m => m.value === budget.month)?.label || ''
    return `Setiap ${budget.date} ${monthName}`
  }
  return '-'
}

const loadData = async () => {
  try {
    const [bRes, cRes] = await Promise.all([
      getBudgets(token),
      getCategories(token)
    ])
    budgets.value = bRes || []
    categories.value = cRes || []
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat data')
  }
}

const resetForm = () => {
  form.id = null
  form.name = ''
  form.amount = ''
  form.interval_type = 'monthly'
  form.day = 1
  form.date = 1
  form.month = 1
  form.repeat = true
  form.category_ids = []
  editingId.value = null
}

const focusFormField = async () => {
  await nextTick()
  formRef.value?.querySelector('input, textarea, select')?.focus()
}

const openNewForm = async () => { resetForm(); isFormOpen.value = true; await focusFormField() }

const openEditForm = async (budget) => {
  form.id = budget.id
  form.name = budget.name || ''
  form.amount = budget.amount || ''
  form.interval_type = budget.interval_type || 'monthly'
  form.day = budget.day || 1
  form.date = budget.date || 1
  form.month = budget.month || 1
  form.repeat = budget.repeat !== false
  form.category_ids = budget.category_ids || []
  editingId.value = budget.id
  isFormOpen.value = true
  await focusFormField()
}

const closeForm = () => { isFormOpen.value = false; resetForm() }

const toggleCategory = (catId) => {
  const index = form.category_ids.indexOf(catId)
  if (index > -1) {
    form.category_ids.splice(index, 1)
  } else {
    form.category_ids.push(catId)
  }
}

const handleSubmit = async (event) => {
  event.preventDefault()
  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data anggaran')
    return
  }
  
  const payload = {
    name: form.name.trim(),
    amount: parseFloat(form.amount),
    interval_type: form.interval_type,
    day: parseInt(form.day),
    date: parseInt(form.date),
    month: parseInt(form.month),
    repeat: form.repeat,
    category_ids: form.category_ids
  }
  
  if (editingId.value) {
    payload.id = editingId.value
  }

  try {
    await upsertBudget(token, payload)
    notification.showSuccess('Anggaran berhasil disimpan')
    await loadData()
    closeForm()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan anggaran')
  }
}

const handleDelete = async (id) => {
  if (!confirm('Apakah Anda yakin ingin menghapus anggaran ini?')) return
  
  try {
    await deleteBudget(token, id)
    notification.showSuccess('Anggaran berhasil dihapus')
    if (editingId.value === id) {
      closeForm()
    }
    await loadData()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menghapus anggaran')
  }
}

onMounted(async () => {
  if (!token) return
  loading.start({ label: 'Memuat data anggaran...' })
  try {
    await loadData()
  } finally {
    loading.stop()
  }
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-lg font-semibold text-slate-900">Anggaran (Budget)</h2>
        <p class="text-slate-500 text-sm">Atur batas pengeluaran Anda.</p>
      </div>
      <button @click="openNewForm"
        class="inline-flex items-center justify-center rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
        Tambah Anggaran
      </button>
    </div>

    <!-- Formulir -->
    <div v-if="isFormOpen" class="bg-white rounded-3xl p-6 shadow-lg">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-slate-200">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">{{ formTitle }}</h3>
          <p class="text-slate-500 text-sm">Isi data anggaran lalu simpan.</p>
        </div>
        <button @click="closeForm" class="text-sm font-medium text-slate-600 transition hover:text-slate-900">Batal</button>
      </div>
      <form ref="formRef" @submit.prevent="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <BaseInput v-model="form.name" label="Nama Anggaran" placeholder="Contoh: Makan Bulanan" required />
        <BaseInput v-model="form.amount" type="number" label="Batas Nominal (Rp)" placeholder="Contoh: 1500000" required min="1" />
        
        <div class="flex flex-col space-y-1.5">
          <label class="text-sm font-medium text-slate-700">Tipe Siklus</label>
          <select v-model="form.interval_type" class="rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition-all focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
            <option v-for="opt in intervalOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
        </div>

        <div v-if="form.interval_type === 'weekly'" class="flex flex-col space-y-1.5">
          <label class="text-sm font-medium text-slate-700">Hari (Senin-Minggu)</label>
          <select v-model="form.day" class="rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition-all focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
            <option v-for="opt in dayOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
        </div>

        <div v-if="form.interval_type === 'monthly' || form.interval_type === 'yearly'" class="flex flex-col space-y-1.5">
          <label class="text-sm font-medium text-slate-700">Tanggal (1-31)</label>
          <input type="number" v-model="form.date" min="1" max="31" class="rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition-all focus:border-blue-500 focus:ring-1 focus:ring-blue-500" required />
        </div>

        <div v-if="form.interval_type === 'yearly'" class="flex flex-col space-y-1.5">
          <label class="text-sm font-medium text-slate-700">Bulan</label>
          <select v-model="form.month" class="rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition-all focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
            <option v-for="opt in monthOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
        </div>
        
        <div class="md:col-span-2 flex items-center mt-2">
          <input type="checkbox" id="repeat" v-model="form.repeat" class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 cursor-pointer">
          <label for="repeat" class="ml-2 text-sm font-medium text-slate-700 cursor-pointer">Otomatis Perbarui Tiap Siklus Baru (Repeat)</label>
        </div>

        <div class="md:col-span-2">
          <label class="text-sm font-medium text-slate-700 block mb-2">Kategori (Kosongkan jika untuk semua kategori)</label>
          <div class="flex flex-wrap gap-2">
            <button 
              type="button" 
              v-for="cat in categories" 
              :key="cat.id" 
              @click="toggleCategory(cat.id)"
              :class="[
                'px-3 py-1.5 rounded-full text-xs font-medium border transition-colors',
                form.category_ids.includes(cat.id) ? 'bg-blue-100 text-blue-700 border-blue-200' : 'bg-gray-50 text-gray-600 border-gray-200 hover:bg-gray-100'
              ]"
            >
              {{ cat.name }}
            </button>
          </div>
        </div>

        <div class="md:col-span-2 flex flex-col gap-3 sm:flex-row sm:justify-end mt-2 pt-2 border-t border-slate-100">
          <button v-if="editingId" type="button" @click="handleDelete(editingId)"
            class="w-full sm:mr-auto rounded-xl border border-red-200 text-red-600 px-4 py-3 text-sm font-semibold transition hover:bg-red-50 sm:w-auto">
            Hapus
          </button>
          <button type="button" @click="closeForm"
            class="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 sm:w-auto">
            Batal
          </button>
          <button type="submit"
            class="w-full rounded-xl bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 sm:w-auto">
            {{ submitLabel }}
          </button>
        </div>
      </form>
    </div>

    <!-- Daftar Anggaran -->
    <div class="bg-white rounded-3xl p-6 shadow-lg">
      <div v-if="budgets.length === 0" class="space-y-3 text-center text-slate-600">
        <p class="text-lg font-medium">Belum ada anggaran</p>
        <p class="text-sm">Klik tombol Tambah Anggaran untuk membuat budget.</p>
      </div>
      <div v-else>
        <!-- Mobile -->
        <div class="md:hidden space-y-3">
          <div v-for="budget in budgets" :key="budget.id"
            class="rounded-2xl bg-slate-50 p-4 shadow-sm space-y-2 cursor-pointer transition hover:bg-slate-100"
            role="button" @click="openEditForm(budget)">
            <div>
              <p class="font-semibold text-slate-900 text-sm">{{ budget.name }} <span v-if="budget.repeat" class="ml-1 text-[10px] bg-green-100 text-green-700 px-2 py-0.5 rounded-full">Repeat</span></p>
              <p class="text-xs text-slate-500 mt-0.5">Rp {{ formatMoney(budget.amount) }} • {{ formatIntervalInfo(budget) }}</p>
              <p class="text-[11px] text-blue-600 mt-1 font-medium">Periode Aktif: {{ formatDate(budget.current_period_start) }} - {{ formatDate(budget.current_period_end) }}</p>
            </div>
            <div class="text-xs text-slate-500 truncate">Kategori: {{ getCategoryNames(budget.category_ids) }}</div>
          </div>
        </div>
        <!-- Desktop -->
        <div class="hidden md:block overflow-x-auto">
          <table class="min-w-full border-separate border-spacing-y-3 text-left">
            <thead>
              <tr class="text-sm text-slate-500">
                <th class="px-4 py-3">Nama Anggaran</th>
                <th class="px-4 py-3">Batas Nominal</th>
                <th class="px-4 py-3">Siklus</th>
                <th class="px-4 py-3">Periode Aktif</th>
                <th class="px-4 py-3">Kategori</th>
                <th class="px-4 py-3">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="budget in budgets" :key="budget.id"
                class="rounded-3xl bg-slate-50 align-top text-sm shadow-sm transition hover:bg-slate-100">
                <td class="px-4 py-4 text-slate-900 font-medium">{{ budget.name }} <span v-if="budget.repeat" class="block mt-1 w-max text-[10px] bg-green-100 text-green-700 px-2 py-0.5 rounded-full">Repeat</span></td>
                <td class="px-4 py-4 font-semibold text-slate-900">Rp {{ formatMoney(budget.amount) }}</td>
                <td class="px-4 py-4 text-slate-600">{{ formatIntervalInfo(budget) }}</td>
                <td class="px-4 py-4 text-slate-600 text-[11px]">
                  <span class="block text-slate-800 font-medium mb-1">Terpakai: Rp {{ formatMoney(budget.total_spent) }}</span>
                  {{ formatDate(budget.current_period_start) }} - {{ formatDate(budget.current_period_end) }}
                </td>
                <td class="px-4 py-4 text-slate-600 text-xs max-w-xs truncate" :title="getCategoryNames(budget.category_ids)">
                  {{ getCategoryNames(budget.category_ids) }}
                </td>
                <td class="px-4 py-4 space-x-2">
                  <button @click.stop="openEditForm(budget)"
                    class="rounded-lg bg-slate-200 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-300">Edit</button>
                  <button @click.stop="handleDelete(budget.id)"
                    class="rounded-lg bg-red-100 px-3 py-1 text-sm text-red-600 transition hover:bg-red-200">Hapus</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
