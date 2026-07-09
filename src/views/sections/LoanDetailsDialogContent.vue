<script setup>
import { computed, reactive, ref } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Notification } from '@packages/utils/Notification.js'
import { createPayment, updatePayment } from '@/DataService.js'
import { getAccounts, getCategories } from '@/DataService.js'


const props = defineProps({
  loan: { type: Object, default: null },
  payments: { type: Array, default: () => [] }
})

const notification = new Notification()

const loading = ref(false)
const showAddPaymentForm = ref(false)
const accounts = ref([])
const categories = ref([])

const form = reactive({
  amount: '',
  transaction_type: 'debit',
  description: '',
  account_id: '',
  category_id: ''
})

const editingPaymentId = ref(null)

const token = localStorage.getItem('access_token')

const loadLookups = async () => {
  if (!token) return
  try {
    // swagger pembayaran mengizinkan account_id (wajib untuk request ini di UI kita)
    // maka dropdown akun & kategori harus di-load di dialog.
    accounts.value = await getAccounts(token)
    categories.value = await getCategories(token)
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat lookup akun/kategori')
    accounts.value = []
    categories.value = []
  }
}

const openAddPayment = async () => {
  if (!accounts.value.length && !categories.value.length) {
    await loadLookups()
  }

  showAddPaymentForm.value = true
  editingPaymentId.value = null

  // defaults
  form.amount = ''
  form.transaction_type = 'debit'
  form.description = ''
  form.category_id = ''

  if (accounts.value.length && !form.account_id) {
    form.account_id = String(accounts.value[0].id)
  }
}

const openEditPayment = async (payment) => {
  if (!accounts.value.length && !categories.value.length) {
    await loadLookups()
  }

  showAddPaymentForm.value = true
  editingPaymentId.value = payment.id

  form.amount = payment.transaction?.amount ? String(payment.transaction.amount) : ''
  form.transaction_type = payment.transaction?.transaction_type || 'debit'
  form.description = payment.transaction?.description || ''
  form.account_id = payment.transaction?.account_id ? String(payment.transaction.account_id) : ''
  form.category_id = payment.transaction?.category_id ? String(payment.transaction.category_id) : ''
}

const closeAddPayment = () => {
  showAddPaymentForm.value = false
  editingPaymentId.value = null
}

const numericAmount = (value) => {
  const n = Number(value)
  return Number.isFinite(n) && n > 0
}

const submitAddPayment = async () => {
  if (!props.loan?.id) return
  if (!form.account_id) {
    notification.showError('Akun wajib dipilih')
    return
  }
  if (!form.amount || Number(form.amount) <= 0) {
    notification.showError('Jumlah harus lebih besar dari 0')
    return
  }

  loading.value = true
  try {
    const payload = {
      loan_id: props.loan.id,
      transaction: {
        amount: Number(form.amount),
        transaction_type: form.transaction_type,
        ...(form.description.trim() ? { description: form.description.trim() } : {}),
        account_id: Number(form.account_id),
        ...(form.category_id ? { category_id: Number(form.category_id) } : {})
      }
    }

    if (editingPaymentId.value) {
      await updatePayment(token, editingPaymentId.value, payload)
      notification.showSuccess('Payment berhasil diperbarui')
    } else {
      await createPayment(token, payload)
      notification.showSuccess('Payment berhasil ditambahkan')
    }

    emit('payment-added')
    closeAddPayment()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan payment')
  } finally {
    loading.value = false
  }
}

const emit = defineEmits(['payment-added'])


const formatCurrency = (value) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 2 }).format(Number(value) || 0)

const formatDate = (value) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

const outstanding = computed(() => props.loan?.outstanding_amount ?? props.loan?.balance ?? 0)
const loanTitle = computed(() => props.loan?.name || '-')
const paymentsTotal = computed(() => props.payments.length)
</script>

<template>
  <div class="space-y-4">
    <div class="space-y-1">
      <div class="flex items-center justify-between gap-3 flex-wrap">
        <div class="min-w-0">
          <h2 class="text-lg font-semibold text-slate-900">{{ loanTitle }}</h2>
          <div class="flex flex-wrap gap-2 items-center mt-2">
            <span class="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">
              Terhutang/Piutang: {{ formatCurrency(outstanding) }}
            </span>
            <span class="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">
              Payments: {{ paymentsTotal }}
            </span>
          </div>
          <p class="text-sm text-slate-500 mt-2">Dibuat: {{ formatDate(loan?.created_at) }}</p>
        </div>

        <button
          v-if="loan?.id"
          @click="openAddPayment"
          class="rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700"
        >
          Tambah Payment
        </button>
      </div>

      <div v-if="showAddPaymentForm" class="mt-4 rounded-2xl border border-slate-200 bg-slate-50 p-4">
        <p class="text-sm font-semibold text-slate-900 mb-3">{{ editingPaymentId ? 'Edit Payment' : 'Form Payment' }}</p>
        <form @submit.prevent="submitAddPayment" class="grid gap-4 md:grid-cols-2">
          <BaseInput
            v-model="form.amount"
            label="Jumlah"
            type="money"
            placeholder="0"
            required
            :validate="['Jumlah harus > 0', numericAmount]"
          />

          <div class="md:col-span-1">
            <label class="block text-sm font-medium text-slate-700 mb-2">Tipe</label>
            <div class="grid grid-cols-2 gap-3">
              <label class="flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-3 cursor-pointer hover:bg-slate-50">
                <input type="radio" value="debit" v-model="form.transaction_type" class="accent-blue-600" />
                <span class="text-sm font-semibold text-slate-800">Debit</span>
              </label>
              <label class="flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-3 cursor-pointer hover:bg-slate-50">
                <input type="radio" value="credit" v-model="form.transaction_type" class="accent-blue-600" />
                <span class="text-sm font-semibold text-slate-800">Credit</span>
              </label>
            </div>
          </div>

          <div class="md:col-span-2">
            <BaseInput
              v-model="form.description"
              label="Deskripsi"
              placeholder="Opsional: pembayaran hutang"
            />
          </div>

          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-slate-700 mb-2">Akun</label>
            <select
              v-model="form.account_id"
              class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              required
            >
              <option value="">Pilih akun</option>
              <option v-for="a in accounts" :key="a.id" :value="String(a.id)">{{ a.account_name }}</option>
            </select>
          </div>

          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-slate-700 mb-2">Kategori (opsional)</label>
            <select
              v-model="form.category_id"
              class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            >
              <option value="">Tidak ada</option>
              <option v-for="c in categories" :key="c.id" :value="String(c.id)">{{ c.name }}</option>
            </select>
          </div>

          <div class="md:col-span-2 flex flex-col gap-3 sm:flex-row sm:justify-end">
            <button
              type="button"
              @click="closeAddPayment"
              class="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 sm:w-auto"
            >
              Batal
            </button>
            <button
              type="submit"
              class="w-full rounded-xl bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 sm:w-auto"
              :disabled="loading"
            >
              {{ loading ? 'Menyimpan...' : 'Simpan Payment' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <div class="space-y-2" >
      <p class="text-sm font-semibold text-slate-900">Daftar Payments</p>
      <div v-if="payments.length === 0" class="text-sm text-slate-500">Belum ada payment untuk hutang ini.</div>
      <div v-else class="overflow-x-auto">
        <table class="min-w-full border-separate border-spacing-y-2 text-left">
          <thead>
            <tr class="text-xs text-slate-500">
              <th class="px-3 py-2">ID</th>
              <th class="px-3 py-2">Balance</th>
              <th class="px-3 py-2">Created</th>
              <th class="px-3 py-2">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="p in payments"
              :key="p.id"
              class="rounded-xl bg-slate-50 text-sm shadow-sm"
            >
              <td class="px-3 py-2 text-slate-900">{{ p.id }}</td>
              <td class="px-3 py-2 text-slate-700">{{ formatCurrency(p.transaction?.amount) }}</td>
              <td class="px-3 py-2 text-slate-500">{{ formatDate(p.created_at) }}</td>
              <td class="px-3 py-2">
                <button
                  @click="openEditPayment(p)"
                  class="rounded-lg bg-blue-100 px-3 py-1 text-xs text-blue-700 transition hover:bg-blue-200"
                >
                  Edit
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- (bagian bawah sudah digantikan) -->
  </div>
</template>


