<script setup>
import { reactive, ref, computed, onMounted, nextTick, h, onUnmounted, watch } from 'vue'
import { BaseInput, BaseLookup, BaseSelect, BaseRoll, ToggleFeature } from '@packages/components'
import dayjs from "dayjs"
import 'dayjs/locale/id'
dayjs.locale('id')
import { Loading } from '@packages/utils/Loading.js'
import { Notification } from '@packages/utils/Notification.js'
import { Dialog } from '@packages/utils/Dialog.js'
import { Config } from '@/Config.js'
import {
  createTransaction,
  deleteTransaction,
  getAccounts,
  getTransactions,
  updateTransaction
} from '@/DataService.js'

const notification = new Notification()
const loading = new Loading()
const transactions = ref([])
const accounts = ref([])
const searchQuery = ref('')
const isFormOpen = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const selectedMobileTransaction = ref(null)
const token = localStorage.getItem('access_token')
const currentPage = ref(1)
const hasMore = ref(true)
const loadMoreObserverRef = ref(null)
let observer = null
const isLoadingMore = ref(false)

const form = reactive({
  transaction_date: dayjs().format('YYYY-MM-DD'),
  description: '',
  amount: '',
  type: 'debit',
  account_id: '',
  category_id: ''
})

const categoriesLookupRoute = `${Config.url}/categories`

const typeOptions = [
  { value: 'credit', label: 'Expense' },
  { value: 'debit', label: 'Income' }
]

const transactionActionsDialogContent = {
  name: 'TransactionActionsDialogContent',
  emits: ['dialog-ok', 'dialog-close'],
  setup(_, { emit }) {
    return () => {
      const transaction = selectedMobileTransaction.value
      const summary = transaction
        ? `${transaction.description || '-'} • ${formatCurrency(transaction.amount)} • ${formatDate(transaction.transaction_date)}`
        : '-'

      return h('div', { class: 'space-y-4' }, [
        h('div', { class: 'space-y-1' }, [
          h('h2', { class: 'text-lg font-semibold text-slate-900' }, 'Aksi Transaksi'),
          h('p', { class: 'text-sm leading-6 text-slate-600' }, summary),
          h('p', { class: 'text-xs text-slate-500' }, 'Pilih tindakan untuk transaksi ini.')
        ]),
        h('div', { class: 'grid gap-3' }, [
          h(
            'button',
            {
              type: 'button',
              class: 'w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50',
              onClick: () => emit('dialog-ok', 'edit')
            },
            'Edit'
          ),
          h(
            'button',
            {
              type: 'button',
              class: 'w-full rounded-xl bg-red-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-red-700',
              onClick: () => emit('dialog-ok', 'delete')
            },
            'Hapus'
          ),
          h(
            'button',
            {
              type: 'button',
              class: 'w-full rounded-xl bg-slate-100 px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-200',
              onClick: () => emit('dialog-close')
            },
            'Batal'
          )
        ])
      ])
    }
  }
}

const mobileActionsDialog = new Dialog(transactionActionsDialogContent, 'bottom')

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

const getDateKey = (value) => {
  if (!value) return 'unknown-date'
  const date = dayjs(value)
  if (!date.isValid()) return 'unknown-date'

  return date.format('YYYY-MM-DD')
}

const groupedTransactions = computed(() => {
  const groups = new Map()

  for (const transaction of filteredTransactions.value) {
    const dateKey = getDateKey(transaction.transaction_date)
    const existingGroup = groups.get(dateKey)
    const amount = Number(transaction.amount) || 0
    const isIncome = transaction.type === 'debit'

    if (existingGroup) {
      existingGroup.items.push(transaction)
      if (isIncome) {
        existingGroup.totalIncome += amount
      } else {
        existingGroup.totalExpense += amount
      }
      continue
    }

    groups.set(dateKey, {
      dateKey,
      dateLabel: formatDate(transaction.transaction_date),
      items: [transaction],
      totalIncome: isIncome ? amount : 0,
      totalExpense: isIncome ? 0 : amount
    })
  }

  return Array.from(groups.values()).map((group) => ({
    ...group,
    totalAmount: group.totalIncome - group.totalExpense
  })).sort((left, right) => {
    if (left.dateKey === 'unknown-date') return 1
    if (right.dateKey === 'unknown-date') return -1
    return right.dateKey.localeCompare(left.dateKey)
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

const loadTransactions = async (isLoadMore = false) => {
  if (isLoadMore) {
    isLoadingMore.value = true
  }
  try {
    if (!isLoadMore) {
      currentPage.value = 1
      hasMore.value = true
    }
    const transactionsData = await getTransactions(token, null, null, currentPage.value)
    
    if (transactionsData.length < 100) {
      hasMore.value = false
    }

    const mapped = transactionsData.map((transaction) => ({
      ...transaction,
      type: transaction.transaction_type
    }))

    if (isLoadMore) {
      transactions.value = [...transactions.value, ...mapped]
    } else {
      transactions.value = mapped
    }
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat transaksi')
  } finally {
    if (isLoadMore) {
      isLoadingMore.value = false
    }
  }
}

const handleLoadMore = async () => {
  if (!hasMore.value || isLoadingMore.value) return
  currentPage.value++
  await loadTransactions(true)
}

const resetForm = () => {
  form.description = ''
  form.amount = ''
  form.type = 'debit'
  form.account_id = accounts.value.length ? String(accounts.value[0].id) : ''
  form.category_id = ''
  editingId.value = null
}

const focusFormField = async () => {
  await nextTick()
  formRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  formRef.value?.querySelector('input, select, textarea, button')?.focus()
}

const openNewForm = async () => {
  if (!accounts.value.length) {
    notification.showError('Silakan buat akun terlebih dahulu sebelum menambahkan transaksi')
    return
  }
  resetForm()
  isFormOpen.value = true
  await focusFormField()
}

const openEditForm = async (transaction) => {
  form.description = transaction.description || ''
  form.amount = transaction.amount != null ? String(transaction.amount) : ''
  form.type = transaction.type || 'debit'
  form.account_id = transaction.account_id ? String(transaction.account_id) : accounts.value.length ? String(accounts.value[0].id) : ''
  form.category_id = transaction.category_id ? String(transaction.category_id) : transaction.category?.id ? String(transaction.category.id) : ''
  form.transaction_date = transaction.transaction_date
    ? dayjs(transaction.transaction_date).format('YYYY-MM-DD')
    : dayjs().format('YYYY-MM-DD')
  editingId.value = transaction.id
  isFormOpen.value = true
  await focusFormField()
}

const closeForm = () => {
  isFormOpen.value = false
  editingId.value = null
  form.description = ''
  form.amount = ''
  form.type = 'debit'
  form.account_id = ''
  form.category_id = ''
  form.transaction_date = dayjs().format('YYYY-MM-DD')
  isFormOpen.value = false
  resetForm()
}

const openMobileActions = async (transaction) => {
  selectedMobileTransaction.value = transaction
  const result = await mobileActionsDialog.open()
  selectedMobileTransaction.value = null

  if (result.action !== 'ok') return

  if (result.data === 'edit') {
    await openEditForm(transaction)
    return
  }

  if (result.data === 'delete' && transaction?.id) {
    await handleDelete(transaction.id)
  }
}

const validDescription = (value) => String(value).trim().length > 0
const validLookupValue = (value) => String(value || '').trim().length > 0
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
    account_id: Number(form.account_id),
    category_id: form.category_id ? Number(form.category_id) : null,
    transaction_date: form.transaction_date ? dayjs(form.transaction_date).toISOString() : dayjs().toISOString()
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
  const confirmDialog = new Dialog(
    {
      name: 'DeleteTransactionDialogContent',
      emits: ['dialog-ok', 'dialog-close'],
      setup() {
        return () =>
          h('div', { class: 'space-y-3' }, [
            h('h2', { class: 'text-lg font-semibold text-slate-900' }, 'Hapus transaksi?'),
            h('p', { class: 'text-sm leading-6 text-slate-600' }, 'Transaksi yang dihapus tidak dapat dikembalikan.')
          ])
      }
    },
    'top',
    { label: 'Hapus', value: true },
    { label: 'Batal', value: false }
  )

  const confirmed = await confirmDialog.open()
  if (confirmed.action !== 'ok') return

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
  return dayjs(value).format('DD MMM YYYY')
}

onMounted(async () => {
  if (!token) {
    notification.showError('Token pengguna tidak ditemukan. Silakan login kembali.')
    return
  }

  loading.start({ label: 'Memuat data transaksi...' })

  try {
    await loadAccounts()
    await loadTransactions()
  } finally {
    loading.stop()
  }

  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && hasMore.value) {
      handleLoadMore()
    }
  }, {
    rootMargin: '100px'
  })

  watch(loadMoreObserverRef, (el) => {
    if (el) {
      observer.observe(el)
    }
  }, { immediate: true })
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
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

      <form ref="formRef" @submit.prevent="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <div class="space-y-1.5 md:col-span-2">
          <label class="block text-sm font-semibold text-slate-700">Tipe Transaksi</label>
          <ToggleFeature v-model="form.type" :options="typeOptions" />
        </div>

        <BaseInput v-model="form.description" label="Deskripsi" placeholder="Contoh: Beli makan siang" required
        :validate="['Deskripsi wajib diisi', validDescription]" />

        <BaseLookup v-model="form.account_id" label="Akun" placeholder="Pilih akun" required
        :validate="['Akun wajib dipilih', validLookupValue]" :route="`${Config.url}/accounts`" item-key="id"
        display="account_name" :auth="true" />

        <BaseLookup v-model="form.category_id" label="Kategori" placeholder="Pilih kategori" required
        :validate="['Kategori wajib dipilih', validLookupValue]" :route="categoriesLookupRoute" item-key="id"
        display="name" :auth="true" />

        <BaseInput v-model="form.amount" label="Jumlah" type="money" placeholder="0" required
        :validate="['Jumlah harus lebih besar dari 0', numeric]" />

        <BaseInput v-model="form.transaction_date" type="date" label="Tanggal Transaksi" required />
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

      <div v-else class="max-h-[70vh] space-y-4 overflow-y-auto pr-1">
        <BaseRoll v-for="group in groupedTransactions" :key="group.dateKey" type="vertical"
          :initially-open="false"
          :label="`${group.dateLabel}`" gap-class="gap-3 px-4 pb-4">
          <template #description>
            <div class="flex flex-wrap items-center gap-2">
              <span class="rounded-full bg-slate-100 px-2 py-1 text-[11px] font-medium text-slate-600">
                Total: {{ group.items.length }} transaksi
              </span>
              <span class="rounded-full bg-emerald-100 px-2 py-1 text-[11px] font-medium text-emerald-700">
                +{{ formatCurrency(group.totalIncome) }}
              </span>
              <span class="rounded-full bg-rose-100 px-2 py-1 text-[11px] font-medium text-rose-700">
                -{{ formatCurrency(group.totalExpense) }}
              </span>
              <span class="rounded-full px-2 py-1 text-[11px] font-semibold"
                :class="group.totalAmount >= 0 ? 'bg-emerald-50 text-emerald-700' : 'bg-rose-50 text-rose-700'">
                Netto: {{ group.totalAmount >= 0 ? '+' : '' }}{{ formatCurrency(group.totalAmount) }}
              </span>
            </div>
          </template>

          <!-- Mobile card layout grouped by date -->
          <div class="md:hidden space-y-3">
            <div v-for="transaction in group.items" :key="transaction.id"
              class="rounded-2xl bg-slate-50 p-4 shadow-sm cursor-pointer transition hover:bg-slate-100 active:scale-[0.99]"
              role="button" tabindex="0" @click="openMobileActions(transaction)"
              @keydown.enter.prevent="openMobileActions(transaction)"
              @keydown.space.prevent="openMobileActions(transaction)">
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0">
                  <p class="font-semibold text-slate-900 text-sm truncate">{{ transaction.description }}</p>
                  <p class="text-xs text-slate-500 mt-0.5">{{ getAccountName(transaction.account_id) }}</p>
                </div>
                <div class="text-right shrink-0">
                  <p class="font-semibold text-sm text-slate-900">{{ formatCurrency(transaction.amount) }}</p>
                  <span :class="transaction.type === 'debit'
                    ? 'rounded-full bg-emerald-100 px-2 py-0.5 text-xs text-emerald-700'
                    : 'rounded-full bg-rose-100 px-2 py-0.5 text-xs text-rose-700'">
                    {{ transaction.type === 'debit' ? 'Income' : 'Expense' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Desktop table layout grouped by date -->
          <div class="hidden md:block overflow-x-auto">
            <table class="min-w-full border-separate border-spacing-y-3 text-left">
              <thead>
                <tr class="text-sm text-slate-500">
                  <th class="px-4 py-3">Deskripsi</th>
                  <th class="px-4 py-3">Akun</th>
                  <th class="px-4 py-3">Tipe</th>
                  <th class="px-4 py-3 text-right">Jumlah</th>
                  <th class="px-4 py-3">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="transaction in group.items" :key="transaction.id"
                  class="rounded-3xl bg-slate-50 align-top text-sm shadow-sm transition hover:bg-slate-100">
                  <td class="px-4 py-4 text-slate-900">{{ transaction.description }}</td>
                  <td class="px-4 py-4 text-slate-600">{{ getAccountName(transaction.account_id) }}</td>
                  <td class="px-4 py-4">
                    <span :class="transaction.type === 'debit'
                      ? 'rounded-full bg-emerald-100 px-3 py-1 text-xs text-emerald-700'
                      : 'rounded-full bg-rose-100 px-3 py-1 text-xs text-rose-700'">
                      {{ transaction.type === 'debit' ? 'Income' : 'Expense' }}
                    </span>
                  </td>
                  <td class="px-4 py-4 text-right font-semibold text-slate-900">{{ formatCurrency(transaction.amount)
                  }}</td>
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
        </BaseRoll>

        <div v-if="hasMore && filteredTransactions.length > 0" ref="loadMoreObserverRef" class="flex justify-center pt-2 pb-6">
          <p class="text-sm font-medium text-slate-500">Memuat lebih banyak...</p>
        </div>
      </div>
    </div>
  </section>
</template>
