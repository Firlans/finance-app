<script setup>
import { reactive, ref, computed, watch, nextTick } from 'vue'
import dayjs from 'dayjs'
import { BaseButton, FormFeatures, BaseInput, BaseSelect } from '@packages/components'
import { Loading } from '@packages/utils/Loading.js'
import { Notification } from '@packages/utils/Notification.js'
import { createTransfer, getTransferByID, updateTransfer } from '@/DataService.js'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  accounts: {
    type: Array,
    default: () => []
  },
  editingId: {
    type: [Number, String],
    default: null
  }
})

const emit = defineEmits(['close', 'success'])

const loading = new Loading()
const notification = new Notification()
const token = localStorage.getItem('access_token')

const formRef = ref(null)

const form = reactive({
  from_account_id: '',
  to_account_id: '',
  amount: '',
  admin_fee: 0,
  description: '',
  transaction_date: dayjs().format('YYYY-MM-DD')
})

const modalTitle = computed(() => (props.editingId ? 'Edit Transfer / Tarik Tunai' : 'Transfer / Tarik Tunai'))
const submitLabel = computed(() => (props.editingId ? 'Simpan Perubahan' : 'Kirim Transfer'))

const fromAccountOptions = computed(() => {
  return props.accounts.map(acc => ({
    value: acc.id,
    label: `${acc.account_name} (${formatCurrency(acc.balance)})`
  }))
})

const toAccountOptions = computed(() => {
  return props.accounts
    .filter(acc => String(acc.id) !== String(form.from_account_id))
    .map(acc => ({
      value: acc.id,
      label: `${acc.account_name} (${formatCurrency(acc.balance)})`
    }))
})

const selectedFromAccount = computed(() => {
  return props.accounts.find(a => String(a.id) === String(form.from_account_id))
})

const totalDeducted = computed(() => {
  const amountNum = Number(form.amount) || 0
  const feeNum = Number(form.admin_fee) || 0
  return amountNum + feeNum
})

const isBalanceSufficient = computed(() => {
  if (!selectedFromAccount.value) return true
  return Number(selectedFromAccount.value.balance) >= totalDeducted.value
})

const formatCurrency = (val) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(val) || 0)

watch([() => props.isOpen, () => props.editingId], async ([open, editId]) => {
  if (open) {
    if (editId) {
      loading.start({ label: 'Memuat data transfer...' })
      try {
        const transferData = await getTransferByID(token, editId)
        if (transferData) {
          form.from_account_id = transferData.from_account_id
          form.to_account_id = transferData.to_account_id
          form.amount = transferData.amount
          form.admin_fee = transferData.admin_fee || 0
          form.description = transferData.description || ''
          if (transferData.transaction_date) {
            form.transaction_date = dayjs(transferData.transaction_date).format('YYYY-MM-DD')
          }
        }
      } catch (err) {
        notification.showError(err?.message || 'Gagal memuat detail transfer')
      } finally {
        loading.stop()
      }
    } else {
      resetForm()
      if (props.accounts.length >= 2 && !form.from_account_id) {
        form.from_account_id = props.accounts[0].id
        form.to_account_id = props.accounts[1].id
      }
    }
    await nextTick()
    formRef.value?.querySelector('input, select')?.focus()
  } else {
    resetForm()
  }
})

const resetForm = () => {
  form.from_account_id = ''
  form.to_account_id = ''
  form.amount = ''
  form.admin_fee = 0
  form.description = ''
  form.transaction_date = dayjs().format('YYYY-MM-DD')
}

const handleSubmit = async (event) => {
  if (!form.from_account_id) {
    notification.showError('Pilih Akun Asal')
    return
  }
  if (!form.to_account_id) {
    notification.showError('Pilih Akun Tujuan')
    return
  }
  if (String(form.from_account_id) === String(form.to_account_id)) {
    notification.showError('Akun Asal dan Akun Tujuan tidak boleh sama')
    return
  }
  if (!form.amount || Number(form.amount) <= 0) {
    notification.showError('Masukkan nominal transfer yang valid')
    return
  }
  if (!isBalanceSufficient.value) {
    notification.showError(`Saldo ${selectedFromAccount.value?.account_name || 'Akun Asal'} tidak mencukupi (Tersedia: ${formatCurrency(selectedFromAccount.value?.balance)})`)
    return
  }

  event.loading.start()
  try {
    const payload = {
      from_account_id: Number(form.from_account_id),
      to_account_id: Number(form.to_account_id),
      amount: Number(form.amount),
      admin_fee: Number(form.admin_fee) || 0,
      description: form.description || '',
      transaction_date: new Date(form.transaction_date).toISOString()
    }
    if (props.editingId) {
      await updateTransfer(token, props.editingId, payload)
      notification.showSuccess('Transfer berhasil diperbarui!')
    } else {
      await createTransfer(token, payload)
      notification.showSuccess('Transfer berhasil dilakukan!')
    }
    emit('success')
    emit('close')
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan transfer')
  } finally {
    event.loading.stop()
  }
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm animate-fade-in">
    <div ref="formRef" class="w-full max-w-lg overflow-hidden rounded-2xl bg-white shadow-2xl transition-all border border-slate-100">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-slate-100 bg-gradient-to-r from-blue-600 to-indigo-600 px-6 py-4 text-white">
        <div class="flex items-center space-x-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-white/20 backdrop-blur-md">
            <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-bold">{{ modalTitle }}</h3>
            <p class="text-xs text-blue-100">Pindahkan saldo antar akun milik Anda</p>
          </div>
        </div>
        <button type="button" @click="$emit('close')" class="rounded-lg p-1 text-white/80 transition hover:bg-white/20 hover:text-white">
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Form Body -->
      <FormFeatures @submit="handleSubmit" class="space-y-4 p-6">
        <!-- Akun Asal & Akun Tujuan -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-500">Akun Asal (Keluar)</label>
            <BaseSelect v-model="form.from_account_id" :options="fromAccountOptions" placeholder="Pilih Akun Asal" />
          </div>
          <div>
            <label class="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-500">Akun Tujuan (Masuk)</label>
            <BaseSelect v-model="form.to_account_id" :options="toAccountOptions" placeholder="Pilih Akun Tujuan" />
          </div>
        </div>

        <!-- Nominal & Admin Fee -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-500">Nominal Transfer (Rp)</label>
            <BaseInput v-model="form.amount" type="number" placeholder="Contoh: 500000" min="1" step="any" required />
          </div>
          <div>
            <label class="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-500">Biaya Admin (Opsional)</label>
            <BaseInput v-model="form.admin_fee" type="number" placeholder="Contoh: 2500" min="0" step="any" />
          </div>
        </div>

        <!-- Tanggal & Catatan -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-500">Tanggal</label>
            <BaseInput v-model="form.transaction_date" type="date" required />
          </div>
          <div>
            <label class="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-500">Catatan / Keterangan</label>
            <BaseInput v-model="form.description" type="text" placeholder="Misal: Tarik tunai di ATM BCA" />
          </div>
        </div>

        <!-- Info Summary Box -->
        <div v-if="form.amount" class="rounded-xl border p-3 text-sm transition-all" :class="isBalanceSufficient ? 'border-blue-100 bg-blue-50/50 text-blue-900' : 'border-red-200 bg-red-50 text-red-800'">
          <div class="flex items-center justify-between font-medium">
            <span>Total Potongan Saldo Asal:</span>
            <span class="font-bold">{{ formatCurrency(totalDeducted) }}</span>
          </div>
          <div v-if="form.admin_fee > 0" class="mt-1 flex items-center justify-between text-xs opacity-80">
            <span>(Nominal {{ formatCurrency(form.amount) }} + Fee {{ formatCurrency(form.admin_fee) }})</span>
          </div>
          <div v-if="!isBalanceSufficient" class="mt-2 text-xs font-semibold text-red-600">
            ⚠️ Saldo akun asal tidak mencukupi (Tersedia: {{ formatCurrency(selectedFromAccount?.balance) }})
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-end space-x-3 pt-2">
          <button type="button" @click="$emit('close')" class="rounded-xl border border-slate-200 bg-white px-5 py-2.5 text-sm font-semibold text-slate-700 transition hover:bg-slate-50">
            Batal
          </button>
          <BaseButton type="submit" :disabled="!isBalanceSufficient" buttonClass="inline-flex items-center space-x-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:opacity-50">
            <span>{{ submitLabel }}</span>
          </BaseButton>
        </div>
      </FormFeatures>
    </div>
  </div>
</template>
