<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Notification } from '@packages/utils/Notification.js'
import { createAccount, deleteAccount, getAccounts, updateAccount } from '@/DataService.js'

const notification = new Notification()
const token = localStorage.getItem('access_token')

const accounts = ref([])
const searchQuery = ref('')
const isFormOpen = ref(false)
const editingId = ref(null)

const form = reactive({ account_name: '', description: '', balance: '' })

const filteredAccounts = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return accounts.value.filter((a) =>
    !q || [a.account_name, a.description].join(' ').toLowerCase().includes(q)
  )
})

const formTitle = computed(() => (editingId.value ? 'Edit Akun' : 'Tambah Akun'))
const submitLabel = computed(() => (editingId.value ? 'Simpan Perubahan' : 'Tambah Akun'))

const formatCurrency = (value) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(value) || 0)

const formatDate = (value) => {
  if (!value) return '-'
  return new Date(value).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
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
  form.balance = 0
  editingId.value = null
}

const openNewForm = () => { resetForm(); isFormOpen.value = true }
const openEditForm = (account) => {
  form.account_name = account.account_name || ''
  form.description = account.description || ''
  form.balance = account.balance != null ? account.balance : 0
  editingId.value = account.id
  isFormOpen.value = true
}
const closeForm = () => { isFormOpen.value = false; resetForm() }

const numeric = (value) => {
  const number = Number(value)
  return Number.isFinite(number) && number >= 0
}

const handleSubmit = async (event) => {
  event.preventDefault()
  if (!event.target.reportValidity()) {
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

const handleDelete = async (id) => {
  if (!window.confirm('Hapus akun ini?')) return
  try {
    await deleteAccount(token, id)
    accounts.value = accounts.value.filter((a) => a.id !== id)
    notification.showSuccess('Akun berhasil dihapus')
  } catch (error) {
    notification.showError(error?.message || 'Gagal menghapus akun')
  }
}

onMounted(() => {
  if (token) loadAccounts()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-lg font-semibold text-slate-900">Akun</h2>
        <p class="text-slate-500 text-sm">Kelola akun tabungan dan saldo Anda.</p>
      </div>
      <button @click="openNewForm"
        class="inline-flex items-center justify-center rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
        Tambah Akun
      </button>
    </div>

    <div class="grid gap-4 md:grid-cols-[1fr_auto]">
      <input v-model="searchQuery" type="search" placeholder="Cari akun..."
        class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
      <div class="text-sm text-slate-500 self-end">Total: {{ filteredAccounts.length }} akun</div>
    </div>

    <div v-if="isFormOpen" class="bg-white rounded-3xl p-6 shadow-lg">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-slate-200">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">{{ formTitle }}</h3>
          <p class="text-slate-500 text-sm">Isi data akun lalu simpan.</p>
        </div>
        <button @click="closeForm" class="text-sm font-medium text-slate-600 transition hover:text-slate-900">Batal</button>
      </div>
      <form @submit.prevent="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <BaseInput v-model="form.account_name" label="Nama Akun" placeholder="Contoh: Tabungan Pribadi"
          required
          :validate="['Nama akun wajib diisi', value => String(value).trim().length > 0]" />
        <BaseInput v-model="form.balance" label="Saldo" type="number" placeholder="0"
          required
          :validate="['Saldo harus berupa angka positif', numeric]" />
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
      <div v-else>
        <!-- Mobile -->
        <div class="md:hidden space-y-3">
          <div v-for="account in filteredAccounts" :key="account.id"
            class="rounded-2xl bg-slate-50 p-4 shadow-sm space-y-2">
            <div class="flex items-start justify-between gap-2">
              <div>
                <p class="font-semibold text-slate-900 text-sm">{{ account.account_name }}</p>
                <p class="text-xs text-slate-500 mt-0.5">{{ account.description || '-' }}</p>
              </div>
              <span class="text-sm font-semibold text-slate-900 shrink-0">{{ formatCurrency(account.balance) }}</span>
            </div>
            <div class="text-xs text-slate-400">Dibuat: {{ formatDate(account.created_at) }}</div>
            <div class="flex gap-2 pt-1">
              <button @click="openEditForm(account)"
                class="flex-1 rounded-lg bg-slate-100 px-3 py-2 text-sm text-slate-700 transition hover:bg-slate-200">Edit</button>
              <button @click="handleDelete(account.id)"
                class="flex-1 rounded-lg bg-red-600 px-3 py-2 text-sm font-semibold text-white transition hover:bg-red-700">Hapus</button>
            </div>
          </div>
        </div>
        <!-- Desktop -->
        <div class="hidden md:block overflow-x-auto">
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
                    class="rounded-lg bg-slate-100 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-200">Edit</button>
                  <button @click="handleDelete(account.id)"
                    class="rounded-lg bg-red-600 px-3 py-1 text-sm font-semibold text-white transition hover:bg-red-700">Hapus</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
