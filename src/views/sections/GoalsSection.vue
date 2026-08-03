<script setup>
import { reactive, ref, computed, onMounted, nextTick } from 'vue'
import { BaseButton, FormFeatures, BaseInput } from '@packages/components'
import { Loading } from '@packages/utils/Loading.js'
import { Notification } from '@packages/utils/Notification.js'
import { getGoals, createGoal, updateGoal, deleteGoal } from '@/DataService.js'

const loading = new Loading()
const notification = new Notification()
const token = localStorage.getItem('access_token')

const goals = ref([])
const isFormOpen = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({
  name: '',
  target_amount: '',
  current_amount: ''
})

const formTitle = computed(() => (editingId.value ? 'Edit Goal Tabungan' : 'Tambah Goal Tabungan'))
const submitLabel = computed(() => (editingId.value ? 'Simpan Perubahan' : 'Tambah Goal'))

const formatMoney = (value) => {
  if (value == null) return '0'
  return value.toLocaleString('id-ID')
}

const getPercentage = (goal) => {
  if (!goal.target_amount || goal.target_amount === 0) return 0
  return Math.min(100, Math.round((goal.current_amount / goal.target_amount) * 100))
}

const loadData = async () => {
  try {
    const res = await getGoals(token)
    goals.value = res || []
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat data goals')
  }
}

const resetForm = () => {
  form.name = ''
  form.target_amount = ''
  form.current_amount = ''
  editingId.value = null
}

const focusFormField = async () => {
  await nextTick()
  formRef.value?.querySelector('input, textarea, select')?.focus()
}

const openNewForm = async () => { resetForm(); isFormOpen.value = true; await focusFormField() }

const openEditForm = async (goal) => {
  form.name = goal.name || ''
  form.target_amount = goal.target_amount || ''
  form.current_amount = goal.current_amount || ''
  editingId.value = goal.id
  isFormOpen.value = true
  await focusFormField()
}

const closeForm = () => { isFormOpen.value = false; resetForm() }

const handleSubmit = async (event) => {
  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data goal')
    return
  }

  event.loading.start()

  const payload = {
    name: form.name.trim(),
    target_amount: parseFloat(form.target_amount),
    current_amount: parseFloat(form.current_amount) || 0
  }

  try {
    if (editingId.value) {
      await updateGoal(token, editingId.value, payload)
    } else {
      await createGoal(token, payload)
    }
    notification.showSuccess('Goal berhasil disimpan')
    await loadData()
    closeForm()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan goal')
  } finally {
    event.loading.stop()
  }
}

const handleDelete = async (id) => {
  if (!confirm('Apakah Anda yakin ingin menghapus goal ini?')) return

  try {
    await deleteGoal(token, id)
    notification.showSuccess('Goal berhasil dihapus')
    if (editingId.value === id) {
      closeForm()
    }
    await loadData()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menghapus goal')
  }
}

onMounted(async () => {
  if (!token) return
  loading.start({ label: 'Memuat data goals...' })
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
        <h2 class="text-lg font-semibold text-slate-900">Goals Tabungan</h2>
        <p class="text-slate-500 text-sm">Kelola target tabungan Anda.</p>
      </div>
      <button @click="openNewForm"
        class="inline-flex items-center justify-center rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
        Tambah Goal
      </button>
    </div>

    <!-- Formulir -->
    <div v-if="isFormOpen" class="bg-white rounded-3xl p-6 shadow-lg">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-slate-200">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">{{ formTitle }}</h3>
          <p class="text-slate-500 text-sm">Isi data goal lalu simpan.</p>
        </div>
        <button @click="closeForm" class="text-sm font-medium text-slate-600 transition hover:text-slate-900">Batal</button>
      </div>
      <FormFeatures ref="formRef" @submit="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <BaseInput v-model="form.name" label="Nama Tabungan" placeholder="Contoh: Liburan Bali" required />
        <BaseInput v-model="form.target_amount" type="money" label="Target (Rp)" placeholder="Contoh: 5000000" required min="1" />
        <BaseInput v-model="form.current_amount" type="money" label="Jumlah Saat Ini (Rp)" placeholder="Contoh: 2000000" min="0" />

        <div class="md:col-span-2 flex flex-col gap-3 sm:flex-row sm:justify-end mt-2 pt-2 border-t border-slate-100">
          <button v-if="editingId" type="button" @click="handleDelete(editingId)"
            class="w-full sm:mr-auto rounded-xl border border-red-200 text-red-600 px-4 py-3 text-sm font-semibold transition hover:bg-red-50 sm:w-auto">
            Hapus
          </button>
          <button type="button" @click="closeForm"
            class="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 sm:w-auto">
            Batal
          </button>
          <BaseButton type="submit"
            buttonClass="w-full rounded-xl bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 sm:w-auto">
            {{ submitLabel }}
          </BaseButton>
        </div>
      </FormFeatures>
    </div>

    <!-- Daftar Goals -->
    <div class="bg-white rounded-3xl p-6 shadow-lg">
      <div v-if="goals.length === 0" class="space-y-3 text-center text-slate-600">
        <p class="text-lg font-medium">Belum ada goal tabungan</p>
        <p class="text-sm">Klik tombol Tambah Goal untuk membuat target tabungan baru.</p>
      </div>
      <div v-else>
        <!-- Mobile -->
        <div class="md:hidden space-y-3">
          <div v-for="goal in goals" :key="goal.id"
            class="rounded-2xl bg-slate-50 p-4 shadow-sm space-y-3 cursor-pointer transition hover:bg-slate-100"
            role="button" @click="openEditForm(goal)">
            <div class="flex justify-between items-start">
              <div>
                <p class="font-semibold text-slate-900 text-sm">
                  {{ goal.name }}
                  <span v-if="goal.is_completed" class="ml-1 text-[10px] bg-green-100 text-green-700 px-2 py-0.5 rounded-full">✓ Tercapai</span>
                </p>
              </div>
              <div class="text-right">
                <p class="text-xs font-semibold text-slate-900">Rp {{ formatMoney(goal.current_amount) }}</p>
                <p class="text-[11px] text-slate-500">/ Rp {{ formatMoney(goal.target_amount) }}</p>
              </div>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
              <div
                class="h-2 rounded-full transition-all duration-500"
                :class="goal.is_completed ? 'bg-green-500' : 'bg-gradient-to-r from-emerald-400 to-blue-500'"
                :style="{ width: `${getPercentage(goal)}%` }"
              ></div>
            </div>
            <p class="text-xs text-slate-400 text-right">{{ getPercentage(goal) }}% Terkumpul</p>
          </div>
        </div>
        <!-- Desktop -->
        <div class="hidden md:block overflow-x-auto">
          <table class="min-w-full border-separate border-spacing-y-3 text-left">
            <thead>
              <tr class="text-sm text-slate-500">
                <th class="px-4 py-3">Nama Goal</th>
                <th class="px-4 py-3">Target</th>
                <th class="px-4 py-3">Saat Ini</th>
                <th class="px-4 py-3">Progress</th>
                <th class="px-4 py-3">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="goal in goals" :key="goal.id"
                class="rounded-3xl bg-slate-50 align-top text-sm shadow-sm transition hover:bg-slate-100">
                <td class="px-4 py-4 text-slate-900 font-medium">
                  {{ goal.name }}
                  <span v-if="goal.is_completed" class="block mt-1 w-max text-[10px] bg-green-100 text-green-700 px-2 py-0.5 rounded-full">✓ Tercapai</span>
                </td>
                <td class="px-4 py-4 font-semibold text-slate-900">Rp {{ formatMoney(goal.target_amount) }}</td>
                <td class="px-4 py-4 text-slate-600">Rp {{ formatMoney(goal.current_amount) }}</td>
                <td class="px-4 py-4 min-w-[160px]">
                  <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden mb-1">
                    <div
                      class="h-2 rounded-full transition-all duration-500"
                      :class="goal.is_completed ? 'bg-green-500' : 'bg-gradient-to-r from-emerald-400 to-blue-500'"
                      :style="{ width: `${getPercentage(goal)}%` }"
                    ></div>
                  </div>
                  <span class="text-xs text-slate-400">{{ getPercentage(goal) }}% Terkumpul</span>
                </td>
                <td class="px-4 py-4 space-x-2">
                  <button @click.stop="openEditForm(goal)"
                    class="rounded-lg bg-slate-200 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-300">Edit</button>
                  <button @click.stop="handleDelete(goal.id)"
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
