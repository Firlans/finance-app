<script setup>
import { reactive, ref, computed, onMounted, nextTick, h } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Loading } from '@packages/utils/Loading.js'
import { Notification } from '@packages/utils/Notification.js'
import { Dialog } from '@packages/utils/Dialog.js'
import { getAccounts } from '@/DataService.js'
import { getPaymentsByLoan } from '@/DataService.js'
import {
  createLoan,
  deleteLoan,
  getLoans,
  updateLoan
} from '@/DataService.js'
import LoanDetailsDialogContent from './LoanDetailsDialogContent.vue'


const loading = new Loading()
const notification = new Notification()

const dialog = new Dialog(
  {
    name: 'DeleteLoanDialogContent',
    setup() {
      return () =>
        h('div', { class: 'space-y-3' }, [
          h('h2', { class: 'text-lg font-semibold text-slate-900' }, 'Hapus hutang?'),
          h('p', { class: 'text-sm leading-6 text-slate-600' }, 'Hutang yang dihapus tidak dapat dikembalikan.')
        ])
    }
  },
  'top',
  { label: 'Hapus', value: true },
  { label: 'Batal', value: false }
)

const token = localStorage.getItem('access_token')

const loans = ref([])
const accounts = ref([])
const searchQuery = ref('')

const isFormOpen = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const selectedLoan = ref(null)
const payments = ref([])

const isMobileActionDialogOpen = ref(false)
const isDetailsDialogOpen = ref(false)

const openDetailsForLoan = async (loan) => {
  selectedLoan.value = loan
  isDetailsDialogOpen.value = true
  try {
    payments.value = await getPaymentsByLoan(token, loan.id)
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat payments')
    payments.value = []
  }
}

const closeDetailsDialog = () => {
  isDetailsDialogOpen.value = false
  selectedLoan.value = null
  payments.value = []
}

const form = reactive({
  name: '',
  balance: '',
  loan_type: 'debt',
  account_id: ''
})

const formatCurrency = (value) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 2 }).format(Number(value) || 0)

const filteredLoans = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return loans.value
  return loans.value.filter((l) =>
    [l.name, l.loan_type].filter(Boolean).join(' ').toLowerCase().includes(q)
  )
})

const formTitle = computed(() => (editingId.value ? 'Edit Hutang' : 'Tambah Hutang'))
const submitLabel = computed(() => (editingId.value ? 'Simpan Perubahan' : 'Tambah Hutang'))

const resetForm = () => {
  form.name = ''
  form.balance = ''
  form.loan_type = 'debt'
  form.account_id = ''
  editingId.value = null
}

const focusFormField = async () => {
  await nextTick()
  formRef.value?.querySelector('input, textarea, select')?.focus()
}

const openNewForm = async () => {
  resetForm()
  isFormOpen.value = true
  await focusFormField()
}

const openEditForm = async (loan) => {
  form.name = loan?.name || ''
  form.balance = loan?.balance != null ? String(loan.balance) : ''
  form.loan_type = loan?.loan_type || 'debt'
  form.account_id = loan?.account_id != null ? String(loan.account_id) : ''
  editingId.value = loan?.id
  isFormOpen.value = true
  await focusFormField()
}

const closeForm = () => {
  isFormOpen.value = false
  resetForm()
}

const loadLoans = async () => {
  if (!token) return
  loans.value = await getLoans(token)
}

const loadAccountsForLookup = async () => {
  if (!token) return
  // optional: hutang bisa dikaitkan dengan account
  accounts.value = await getAccounts(token)
}

const numericBalance = (value) => {
  const number = Number(value)
  return Number.isFinite(number) && number >= 0
}

const handleSubmit = async (event) => {
  event.preventDefault()

  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data hutang')
    return
  }

  const payload = {
    name: form.name.trim(),
    balance: Number(form.balance),
    loan_type: form.loan_type,
    ...(form.account_id ? { account_id: Number(form.account_id) } : {})
  }

  try {
    if (editingId.value) {
      await updateLoan(token, editingId.value, payload)
      notification.showSuccess('Hutang berhasil diperbarui')
    } else {
      await createLoan(token, payload)
      notification.showSuccess('Hutang berhasil dibuat')
    }

    await loadLoans()
    closeForm()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan hutang')
  }
}

const handleDelete = async (id) => {
  const result = await dialog.open()
  if (result.action !== 'ok') return

  try {
    await deleteLoan(token, id)
    loans.value = loans.value.filter((l) => l.id !== id)
    notification.showSuccess('Hutang berhasil dihapus')
  } catch (error) {
    notification.showError(error?.message || 'Gagal menghapus hutang')
  }
}

const openMobileActions = (loan) => {
  selectedLoan.value = loan
  isMobileActionDialogOpen.value = true
}

const viewFromMobileActions = async () => {
  if (!selectedLoan.value) return
  const loan = selectedLoan.value
  closeMobileActions()
  await openDetailsForLoan(loan)
}

const closeMobileActions = () => {
  isMobileActionDialogOpen.value = false
  selectedLoan.value = null
}

const deleteFromMobileActions = async () => {
  if (!selectedLoan.value) return
  const id = selectedLoan.value.id
  closeMobileActions()
  await handleDelete(id)
}

const formatDate = (value) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

onMounted(async () => {
  if (!token) return
  loading.start({ label: 'Memuat data hutang...' })
  try {
    await Promise.all([loadLoans(), loadAccountsForLookup()])
  } finally {
    loading.stop()
  }
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-lg font-semibold text-slate-900">Hutang</h2>
        <p class="text-slate-500 text-sm">Kelola daftar hutang dan saldo terhutang Anda.</p>
      </div>

      <button
        @click="openNewForm"
        class="inline-flex items-center justify-center rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700"
      >
        Tambah Hutang
      </button>
    </div>

    <div class="grid gap-4 md:grid-cols-[1fr_auto]">
      <input
        v-model="searchQuery"
        type="search"
        placeholder="Cari hutang..."
        class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
      />
      <div class="text-sm text-slate-500 self-end">Total: {{ filteredLoans.length }} hutang</div>
    </div>

    <div v-if="isFormOpen" class="bg-white rounded-3xl p-6 shadow-lg">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-slate-200">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">{{ formTitle }}</h3>
          <p class="text-slate-500 text-sm">Isi data hutang lalu simpan.</p>
        </div>
        <button @click="closeForm" class="text-sm font-medium text-slate-600 transition hover:text-slate-900">Batal</button>
      </div>

      <form ref="formRef" @submit.prevent="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <BaseInput
          v-model="form.name"
          label="Nama Hutang"
          placeholder="Contoh: KPR / Cicilan Motor"
          required
          :validate="['Nama hutang wajib diisi', (value) => String(value).trim().length > 0]"
        />

        <BaseInput
          v-model="form.balance"
          label="Balance"
          type="number"
          placeholder="0"
          required
          :validate="['Balance harus angka positif', numericBalance]"
        />

        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-slate-700 mb-2">Tipe</label>
          <div class="grid grid-cols-2 gap-3">
            <label class="flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-3 cursor-pointer hover:bg-slate-50">
              <input type="radio" value="debt" v-model="form.loan_type" class="accent-blue-600" />
              <span class="text-sm font-semibold text-slate-800">Debt</span>
            </label>
            <label class="flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-3 cursor-pointer hover:bg-slate-50">
              <input type="radio" value="receivable" v-model="form.loan_type" class="accent-blue-600" />
              <span class="text-sm font-semibold text-slate-800">Receivable</span>
            </label>
          </div>
        </div>

        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-slate-700 mb-2">Akun (opsional)</label>
          <select
            v-model="form.account_id"
            class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
          >
            <option value="">Pilih akun (tidak ada)</option>
            <option v-for="a in accounts" :key="a.id" :value="String(a.id)">
              {{ a.account_name }}
            </option>
          </select>
        </div>

        <div class="md:col-span-2 flex flex-col gap-3 sm:flex-row sm:justify-end">
          <button
            type="button"
            @click="closeForm"
            class="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 sm:w-auto"
          >
            Batal
          </button>
          <button
            type="submit"
            class="w-full rounded-xl bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 sm:w-auto"
          >
            {{ submitLabel }}
          </button>
        </div>
      </form>
    </div>

    <div class="bg-white rounded-3xl p-6 shadow-lg">
      <div v-if="filteredLoans.length === 0" class="space-y-3 text-center text-slate-600">
        <p class="text-lg font-medium">Belum ada hutang</p>
        <p class="text-sm">Klik tombol Tambah Hutang untuk menambahkan hutang baru.</p>
      </div>

      <div v-else>
        <!-- Mobile -->
        <div class="md:hidden space-y-3">
          <div
            v-for="loan in filteredLoans"
            :key="loan.id"
            class="rounded-2xl bg-slate-50 p-4 shadow-sm space-y-2 cursor-pointer transition hover:bg-slate-100 active:scale-[0.99]"
            role="button"
            tabindex="0"
            @click="openMobileActions(loan)"
            @keydown.enter.prevent="openMobileActions(loan)"
            @keydown.space.prevent="openMobileActions(loan)"
          >
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <p class="font-semibold text-slate-900 text-sm truncate">{{ loan.name }}</p>
                <p class="text-xs text-slate-500 mt-0.5">Tipe: {{ loan.loan_type === 'debt' ? 'Debt' : 'Receivable' }}</p>
              </div>
              <div class="text-right shrink-0">
                <p class="font-semibold text-sm text-slate-900">{{ formatCurrency(loan.outstanding_amount ?? loan.balance) }}</p>
                <span
                  class="rounded-full px-2 py-0.5 text-[11px] font-semibold"
                  :class="loan.loan_type === 'debt' ? 'bg-rose-100 text-rose-700' : 'bg-emerald-100 text-emerald-700'"
                >
                  {{ loan.loan_type === 'debt' ? 'Terhutang' : 'Piutang' }}
                </span>
              </div>
            </div>
            <div class="text-xs text-slate-400">Dibuat: {{ formatDate(loan.created_at) }}</div>
          </div>
        </div>

        <!-- Desktop -->
        <div class="hidden md:block overflow-x-auto">
          <table class="min-w-full border-separate border-spacing-y-3 text-left">
            <thead>
              <tr class="text-sm text-slate-500">
                <th class="px-4 py-3">Nama</th>
                <th class="px-4 py-3">Balance</th>
                <th class="px-4 py-3">Terhutang / Piutang</th>
                <th class="px-4 py-3">Tipe</th>
                <th class="px-4 py-3">Dibuat</th>
                <th class="px-4 py-3">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="loan in filteredLoans"
                :key="loan.id"
                class="rounded-3xl bg-slate-50 align-top text-sm shadow-sm transition hover:bg-slate-100"
              >
<td class="px-4 py-4 text-slate-900">{{ loan.name }}</td>
                <td class="px-4 py-4 text-slate-900">{{ formatCurrency(loan.balance) }}</td>
                <td class="px-4 py-4 text-slate-900">{{ formatCurrency(loan.outstanding_amount ?? 0) }}</td>
                <td class="px-4 py-4">
                  <span
                    class="rounded-full px-3 py-1 text-xs font-semibold"
                    :class="loan.loan_type === 'debt' ? 'bg-rose-100 text-rose-700' : 'bg-emerald-100 text-emerald-700'"
                  >
                    {{ loan.loan_type === 'debt' ? 'Debt' : 'Receivable' }}
                  </span>
                </td>
                <td class="px-4 py-4 text-slate-600">{{ formatDate(loan.created_at) }}</td>
                <td class="px-4 py-4 space-x-2">
                  <button
                    @click="openDetailsForLoan(loan)"
                    class="rounded-lg bg-slate-100 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-200"
                  >
                    View
                  </button>
                  <button
                    @click="openEditForm(loan)"
                    class="rounded-lg bg-slate-100 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-200"
                  >
                    Edit
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Details Dialog -->
    <div
      v-if="isDetailsDialogOpen"
      class="fixed inset-0 z-50 flex items-end justify-center bg-slate-900/40 p-4 md:flex md:items-center"
      @click.self="closeDetailsDialog"
    >
      <div class="w-full max-w-2xl rounded-2xl bg-white p-5 shadow-xl md:p-7 overflow-hidden">
        <div class="flex items-center justify-between gap-3 pb-4 border-b border-slate-200">
          <p class="text-sm font-semibold text-slate-900">Preview Hutang</p>
          <button
            @click="closeDetailsDialog"
            class="rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
          >
            Tutup
          </button>
        </div>

        <LoanDetailsDialogContent :loan="selectedLoan" :payments="payments" />
      </div>
    </div>

    <div
      v-if="isMobileActionDialogOpen"
      class="fixed inset-0 z-50 flex items-end justify-center bg-slate-900/40 p-4 md:hidden"
      @click.self="closeMobileActions"
    >
      <div class="w-full max-w-sm rounded-2xl bg-white p-5 shadow-xl">
        <p class="text-xs font-medium uppercase tracking-wide text-slate-500">Aksi Hutang</p>
        <p class="mt-2 truncate text-base font-semibold text-slate-900">{{ selectedLoan?.name || '-' }}</p>
        <p class="mt-1 text-sm text-slate-500">{{ formatCurrency(selectedLoan?.outstanding_amount ?? selectedLoan?.balance) }}</p>

        <div class="mt-4 grid grid-cols-2 gap-3">
          <button
            @click="viewFromMobileActions"
            class="rounded-xl bg-slate-100 px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-200"
          >
            View
          </button>
          <button
            @click="deleteFromMobileActions"
            class="rounded-xl bg-red-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-red-700"
          >
            Hapus
          </button>
        </div>

        <button
          @click="closeMobileActions"
          class="mt-3 w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
        >
          Batal
        </button>
      </div>
    </div>
  </div>
</template>

