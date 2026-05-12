<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import BaseInput from '@/components/base/BaseInput.vue'
import { Notification } from '@/utils/Notification.js'
import { Validator, required } from '@/utils/Validator.js'
import {
  createAccount,
  deleteAccount,
  getAccounts,
  updateAccount
} from '@/utils/DataService.js'

const notification = new Notification()
const accounts = ref([])
const searchQuery = ref('')
const isFormOpen = ref(false)
const editingId = ref(null)
const token = localStorage.getItem('access_token')

const form = reactive({
  account_name: '',
  description: '',
  balance: ''
})

const errors = reactive({})

const filteredAccounts = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  return accounts.value.filter((account) => {
    const values = [account.account_name, account.description].join(' ').toLowerCase()
    return !query || values.includes(query)
  })
})

const formTitle = computed(() => (editingId.value ? 'Edit Akun' : 'Tambah Akun'))
const submitLabel = computed(() => (editingId.value ? 'Simpan Perubahan' : 'Tambah Akun'))

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

const loadAccounts = async () => {
  try {
    accounts.value = await getAccounts(token)
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat akun')
  }
}

const resetForm = () => {
  form.account_name = ''
  form.description = ''
  form.balance = ''
  editingId.value = null
  Object.keys(errors).forEach((key) => delete errors[key])
}

const openNewForm = () => {
  resetForm()
  isFormOpen.value = true
}

const openEditForm = (account) => {
  form.account_name = account.account_name || ''
  form.description = account.description || ''
  form.balance = account.balance != null ? String(account.balance) : ''
  editingId.value = account.id
  isFormOpen.value = true
  Object.keys(errors).forEach((key) => delete errors[key])
}

const closeForm = () => {
  isFormOpen.value = false
  resetForm()
}

const numeric = (message = 'Saldo harus berupa angka positif') => (value) => {
  const number = Number(value)
  return Number.isFinite(number) && number >= 0 ? '' : message
}

const validateForm = () => {
  const validator = new Validator(form, {
    account_name: [required('Nama akun wajib diisi')],
    balance: [required('Saldo wajib diisi'), numeric('Saldo harus berupa angka positif')]
  })

  if (!validator.validate()) {
    Object.assign(errors, validator.getErrors())
    return false
  }

  Object.keys(errors).forEach((key) => delete errors[key])
  return true
}

const handleSubmit = async () => {
  if (!validateForm()) {
    notification.showError('Periksa kembali data akun')
    return
  }

  const payload = {
    account_name: form.account_name.trim(),
    description: form.description.trim(),
    balance: Number(form.balance)
  }

  try {
    if (editingId.value) {
      await updateAccount(token, editingId.value, payload)
      notification.showSuccess('Akun berhasil diperbarui')
    } else {
      await createAccount(token, payload)
      notification.showSuccess('Akun berhasil dibuat')
    }
    await loadAccounts()
    closeForm()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan akun')
  }
}

const handleDelete = async (accountId) => {
  const confirmed = window.confirm('Hapus akun ini?')
  if (!confirmed) return

  try {
    await deleteAccount(token, accountId)
    accounts.value = accounts.value.filter((account) => account.id !== accountId)
    notification.showSuccess('Akun berhasil dihapus')
  } catch (error) {
    notification.showError(error?.message || 'Gagal menghapus akun')
  }
}

onMounted(() => {
  if (!token) {
    notification.showError('Token pengguna tidak ditemukan. Silakan login kembali.')
    return
  }
  loadAccounts()
})
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-slate-900">Akun</h1>
        <p class="text-slate-600 mt-1">Kelola akun tabungan dan saldo Anda.</p>
      </div>

      <button @click="openNewForm"
        class="inline-flex items-center justify-center rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
        Tambah Akun
      </button>
    </div>

    <div class="grid gap-4 md:grid-cols-[1fr_auto]">
      <div class="relative">
        <input v-model="searchQuery" type="search" placeholder="Cari akun..."
          class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
      </div>
      <div class="text-sm text-slate-500 self-end">Total: {{ filteredAccounts.length }} akun</div>
    </div>

    <div v-if="isFormOpen" class="bg-white rounded-3xl p-6 shadow-lg">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-slate-200">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">{{ formTitle }}</h2>
          <p class="text-slate-500 text-sm">Isi data akun lalu simpan.</p>
        </div>
        <button @click="closeForm"
          class="text-sm font-medium text-slate-600 transition hover:text-slate-900">Batal</button>
      </div>

      <form @submit.prevent="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <BaseInput v-model="form.account_name" label="Nama Akun" placeholder="Contoh: Tabungan Pribadi"
          :error="errors.account_name" />

        <BaseInput v-model="form.balance" label="Saldo" type="number" placeholder="0" :error="errors.balance" />

        <div class="md:col-span-2">
          <BaseInput v-model="form.description" label="Deskripsi" placeholder="Contoh: Tabungan bulanan" />
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
      <div v-if="filteredAccounts.length === 0" class="space-y-3 text-center text-slate-600">
        <p class="text-lg font-medium">Belum ada akun</p>
        <p class="text-sm">Klik tombol Tambah Akun untuk menambahkan akun baru.</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full border-separate border-spacing-y-3 text-left">
          <thead>
            <tr class="text-sm text-slate-500">
              <th class="px-4 py-3">Nama Akun</th>
              <th class="px-4 py-3">Saldo</th>
              <th class="px-4 py-3">Deskripsi</th>
              <th class="px-4 py-3">Dibuat</th>
              <th class="px-4 py-3">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="account in filteredAccounts" :key="account.id"
              class="rounded-3xl bg-slate-50 align-top text-sm shadow-sm transition hover:bg-slate-100">
              <td class="px-4 py-4 text-slate-900">{{ account.account_name }}</td>
              <td class="px-4 py-4 text-slate-900">{{ formatCurrency(account.balance) }}</td>
              <td class="px-4 py-4 text-slate-600">{{ account.description || '-' }}</td>
              <td class="px-4 py-4 text-slate-600">{{ formatDate(account.created_at) }}</td>
              <td class="px-4 py-4 space-x-2">
                <button @click="openEditForm(account)"
                  class="rounded-lg bg-slate-100 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-200">
                  Edit
                </button>
                <button @click="handleDelete(account.id)"
                  class="rounded-lg bg-red-600 px-3 py-1 text-sm font-semibold text-white transition hover:bg-red-700">
                  Hapus
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>
