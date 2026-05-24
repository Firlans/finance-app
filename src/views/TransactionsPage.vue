<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Notification } from '@packages/utils/Notification.js'
import {
  createTransaction,
  deleteTransaction,
  getAccounts,
  getTransactions,
  updateTransaction
} from '@/DataService.js'

const notification = new Notification()
const transactions = ref([])
const accounts = ref([])
const searchQuery = ref('')
const isFormOpen = ref(false)
const editingId = ref(null)
const token = localStorage.getItem('access_token')

const form = reactive({
  description: '',
  amount: '',
  type: 'debit',
  account_id: ''
})

const typeOptions = [
  { value: 'debit', label: 'Expense' },
  { value: 'credit', label: 'Income' }
]

const filteredTransactions = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  return transactions.value.filter((transaction) => {
    const values = [
      transaction.description,
      transaction.type,
      getAccountName(transaction.account_id)
    ]
      .join(' ')
      .toLowerCase()
    return !query || values.includes(query)
  })
})

const formTitle = computed(() => (editingId.value ? 'Edit Transaksi' : 'Tambah Transaksi'))
const submitLabel = computed(() => (editingId.value ? 'Simpan Perubahan' : 'Tambah Transaksi'))

const loadAccounts = async () => {
  try {
    accounts.value = await getAccounts(token)
    if (accounts.value.length && !form.account_id) {
      form.account_id = String(accounts.value[0].id)
    }
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat akun')
  }
}

const loadTransactions = async () => {
  try {
    transactions.value = await getTransactions(token)
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat transaksi')
  }
}

const resetForm = () => {
  form.description = ''
  form.amount = ''
  form.type = 'debit'
  form.account_id = accounts.value.length ? String(accounts.value[0].id) : ''
  editingId.value = null
}

const openNewForm = () => {
  if (!accounts.value.length) {
    notification.showError('Silakan buat akun terlebih dahulu sebelum menambahkan transaksi')
    return
  }
  resetForm()
  isFormOpen.value = true
}

const openEditForm = (transaction) => {
  form.description = transaction.description || ''
  form.amount = transaction.amount != null ? String(transaction.amount) : ''
  form.type = transaction.type || 'debit'
  form.account_id = transaction.account_id ? String(transaction.account_id) : accounts.value.length ? String(accounts.value[0].id) : ''
  editingId.value = transaction.id
  isFormOpen.value = true
}

const closeForm = () => {
  isFormOpen.value = false
  resetForm()
}

const validDescription = (value) => String(value).trim().length > 0
const numeric = (value) => {
  const number = Number(value)
  return Number.isFinite(number) && number > 0
}

const handleSubmit = async (event) => {
  event.preventDefault()
  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data transaksi')
    return
  }

  const payload = {
    description: form.description.trim(),
    amount: Number(form.amount),
    transaction_type: form.type,
    account_id: form.account_id
  }

  try {
    if (editingId.value) {
      await updateTransaction(token, editingId.value, payload)
      notification.showSuccess('Transaksi berhasil diperbarui')
    } else {
      await createTransaction(token, payload)
      notification.showSuccess('Transaksi berhasil ditambahkan')
    }
    await loadTransactions()
    closeForm()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan transaksi')
  }
}

const handleDelete = async (transactionId) => {
  const confirmed = window.confirm('Hapus transaksi ini?')
  if (!confirmed) return

  try {
    await deleteTransaction(token, transactionId)
    transactions.value = transactions.value.filter((transaction) => transaction.id !== transactionId)
    notification.showSuccess('Transaksi berhasil dihapus')
  } catch (error) {
    notification.showError(error?.message || 'Gagal menghapus transaksi')
  }
}

const getAccountName = (accountId) => {
  return accounts.value.find((account) => String(account.id) === String(accountId))?.account_name || '-'
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
  }).format(Number(value) || 0)
}

const formatDate = (value) => {
  if (!value) return '-'
  return new Date(value).toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

const accountOptions = computed(() => accounts.value)

onMounted(async () => {
  if (!token) {
    notification.showError('Token pengguna tidak ditemukan. Silakan login kembali.')
    return
  }
  await loadAccounts()
  await loadTransactions()
})
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-slate-900">Transaksi</h1>
        <p class="text-slate-600 mt-1">Kelola transaksi pemasukan dan pengeluaran Anda.</p>
      </div>

      <button @click="openNewForm"
        class="inline-flex items-center justify-center rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
        Tambah Transaksi
      </button>
    </div>

    <div class="grid gap-4 md:grid-cols-[1fr_auto]">
      <div class="relative">
        <input v-model="searchQuery" type="search" placeholder="Cari transaksi..."
          class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
      </div>
      <div class="text-sm text-slate-500 self-end">Total: {{ filteredTransactions.length }} transaksi</div>
    </div>

    <div v-if="isFormOpen" class="bg-white rounded-3xl p-6 shadow-lg">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-slate-200">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">{{ formTitle }}</h2>
          <p class="text-slate-500 text-sm">Isi data transaksi lalu simpan.</p>
        </div>
        <button @click="closeForm"
          class="text-sm font-medium text-slate-600 transition hover:text-slate-900">Batal</button>
      </div>

      <form @submit.prevent="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <BaseInput v-model="form.description" label="Deskripsi" placeholder="Contoh: Beli makan siang"
          required
          :validate="['Deskripsi wajib diisi', validDescription]" />

        <div class="space-y-1">
          <label class="block text-sm font-medium text-slate-700">Akun</label>
          <select v-model="form.account_id" required
            class="w-full rounded-lg border border-slate-300 bg-white px-4 py-2 text-slate-900 focus:border-blue-500 focus:ring-2 focus:ring-blue-100">
            <option v-for="account in accountOptions" :key="account.id" :value="String(account.id)">
              {{ account.account_name }}
            </option>
          </select>
        </div>

        <BaseInput v-model="form.amount" label="Jumlah" type="number" placeholder="0" required
          :validate="['Jumlah harus lebih besar dari 0', numeric]" />

        <div class="space-y-1">
          <label class="block text-sm font-medium text-slate-700">Tipe</label>
          <select v-model="form.type"
            class="w-full rounded-lg border border-slate-300 bg-white px-4 py-2 text-slate-900 focus:border-blue-500 focus:ring-2 focus:ring-blue-100">
            <option v-for="option in typeOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
        </div>

        <div class="md:col-span-2 flex flex-col gap-3 sm:flex-row sm:justify-end">
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

    <div class="bg-white rounded-3xl p-6 shadow-lg">
      <div v-if="filteredTransactions.length === 0" class="space-y-3 text-center text-slate-600">
        <p class="text-lg font-medium">Belum ada transaksi</p>
        <p class="text-sm">Klik tombol Tambah Transaksi untuk membuat daftar transaksi baru.</p>
      </div>

      <div v-else>
        <!-- Mobile card layout -->
        <div class="md:hidden space-y-3">
          <div v-for="transaction in filteredTransactions" :key="transaction.id"
            class="rounded-2xl bg-slate-50 p-4 shadow-sm space-y-2">
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <p class="font-semibold text-slate-900 text-sm truncate">{{ transaction.description }}</p>
                <p class="text-xs text-slate-500 mt-0.5">{{ getAccountName(transaction.account_id) }}</p>
              </div>
              <div class="text-right shrink-0">
                <p class="font-semibold text-sm text-slate-900">{{ formatCurrency(transaction.amount) }}</p>
                <span :class="transaction.type === 'credit'
                  ? 'rounded-full bg-emerald-100 px-2 py-0.5 text-xs text-emerald-700'
                  : 'rounded-full bg-rose-100 px-2 py-0.5 text-xs text-rose-700'">
                  {{ transaction.type === 'credit' ? 'Income' : 'Expense' }}
                </span>
              </div>
            </div>
            <div class="text-xs text-slate-400">{{ formatDate(transaction.created_at) }}</div>
            <div class="flex gap-2 pt-1">
              <button @click="openEditForm(transaction)"
                class="flex-1 rounded-lg bg-slate-100 px-3 py-2 text-sm text-slate-700 transition hover:bg-slate-200">
                Edit
              </button>
              <button @click="handleDelete(transaction.id)"
                class="flex-1 rounded-lg bg-red-600 px-3 py-2 text-sm font-semibold text-white transition hover:bg-red-700">
                Hapus
              </button>
            </div>
          </div>
        </div>

        <!-- Desktop table layout -->
        <div class="hidden md:block overflow-x-auto">
          <table class="min-w-full border-separate border-spacing-y-3 text-left">
            <thead>
              <tr class="text-sm text-slate-500">
                <th class="px-4 py-3">Deskripsi</th>
                <th class="px-4 py-3">Akun</th>
                <th class="px-4 py-3">Tipe</th>
                <th class="px-4 py-3">Tanggal</th>
                <th class="px-4 py-3 text-right">Jumlah</th>
                <th class="px-4 py-3">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="transaction in filteredTransactions" :key="transaction.id"
                class="rounded-3xl bg-slate-50 align-top text-sm shadow-sm transition hover:bg-slate-100">
                <td class="px-4 py-4 text-slate-900">{{ transaction.description }}</td>
                <td class="px-4 py-4 text-slate-600">{{ getAccountName(transaction.account_id) }}</td>
                <td class="px-4 py-4">
                  <span :class="transaction.type === 'credit'
                    ? 'rounded-full bg-emerald-100 px-3 py-1 text-xs text-emerald-700'
                    : 'rounded-full bg-rose-100 px-3 py-1 text-xs text-rose-700'">
                    {{ transaction.type === 'credit' ? 'Income' : 'Expense' }}
                  </span>
                </td>
                <td class="px-4 py-4 text-slate-600">{{ formatDate(transaction.created_at) }}</td>
                <td class="px-4 py-4 text-right font-semibold text-slate-900">{{ formatCurrency(transaction.amount) }}
                </td>
                <td class="px-4 py-4 space-x-2">
                  <button @click="openEditForm(transaction)"
                    class="rounded-lg bg-slate-100 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-200">
                    Edit
                  </button>
                  <button @click="handleDelete(transaction.id)"
                    class="rounded-lg bg-red-600 px-3 py-1 text-sm font-semibold text-white transition hover:bg-red-700">
                    Hapus
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </section>
</template>
