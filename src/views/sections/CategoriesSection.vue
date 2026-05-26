<script setup>
import { reactive, ref, computed, onMounted, nextTick } from 'vue'
import BaseInput from '@packages/components/base/BaseInput.vue'
import { Notification } from '@packages/utils/Notification.js'
import { getCategories, createCategory, updateCategory, deleteCategory } from '@/DataService.js'

const notification = new Notification()
const token = localStorage.getItem('access_token')

const categories = ref([])
const searchQuery = ref('')
const isFormOpen = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const form = reactive({ name: '', description: '' })

const filteredCategories = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return categories.value.filter((c) =>
    !q || [c.name, c.description].join(' ').toLowerCase().includes(q)
  )
})

const formTitle = computed(() => (editingId.value ? 'Edit Kategori' : 'Tambah Kategori'))
const submitLabel = computed(() => (editingId.value ? 'Simpan Perubahan' : 'Tambah Kategori'))

const formatDate = (value) => {
  if (!value) return '-'
  return new Date(value).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

const loadCategories = async () => {
  try {
    categories.value = await getCategories(token)
  } catch (error) {
    notification.showError(error?.message || 'Gagal memuat kategori')
  }
}

const resetForm = () => {
  form.name = ''
  form.description = ''
  editingId.value = null
}

const focusFormField = async () => {
  await nextTick()
  formRef.value?.querySelector('input, textarea, select')?.focus()
}

const openNewForm = async () => { resetForm(); isFormOpen.value = true; await focusFormField() }
const openEditForm = async (category) => {
  form.name = category.name || ''
  form.description = category.description || ''
  editingId.value = category.id
  isFormOpen.value = true
  await focusFormField()
}
const closeForm = () => { isFormOpen.value = false; resetForm() }

const validCategoryName = (value) => String(value).trim().length > 0

const handleSubmit = async (event) => {
  event.preventDefault()
  if (!event.target.reportValidity()) {
    notification.showError('Periksa kembali data kategori')
    return
  }
  const payload = {
    name: form.name.trim(),
    description: form.description.trim()
  }
  try {
    if (editingId.value) {
      await updateCategory(token, editingId.value, payload)
      notification.showSuccess('Kategori berhasil diperbarui')
    } else {
      await createCategory(token, payload)
      notification.showSuccess('Kategori berhasil dibuat')
    }
    await loadCategories()
    closeForm()
  } catch (error) {
    notification.showError(error?.message || 'Gagal menyimpan kategori')
  }
}

const handleDelete = async (id) => {
  if (!window.confirm('Hapus kategori ini?')) return
  try {
    await deleteCategory(token, id)
    categories.value = categories.value.filter((c) => c.id !== id)
    notification.showSuccess('Kategori berhasil dihapus')
  } catch (error) {
    notification.showError(error?.message || 'Gagal menghapus kategori')
  }
}

onMounted(() => {
  if (token) loadCategories()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-lg font-semibold text-slate-900">Kategori</h2>
        <p class="text-slate-500 text-sm">Kelola kategori untuk transaksi Anda.</p>
      </div>
      <button @click="openNewForm"
        class="inline-flex items-center justify-center rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700">
        Tambah Kategori
      </button>
    </div>

    <div class="grid gap-4 md:grid-cols-[1fr_auto]">
      <input v-model="searchQuery" type="search" placeholder="Cari kategori..."
        class="w-full rounded-2xl border border-slate-300 bg-white py-3 px-4 text-sm text-slate-900 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
      <div class="text-sm text-slate-500 self-end">Total: {{ filteredCategories.length }} kategori</div>
    </div>

    <div v-if="isFormOpen" class="bg-white rounded-3xl p-6 shadow-lg">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-slate-200">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">{{ formTitle }}</h3>
          <p class="text-slate-500 text-sm">Isi data kategori lalu simpan.</p>
        </div>
        <button @click="closeForm" class="text-sm font-medium text-slate-600 transition hover:text-slate-900">Batal</button>
      </div>
      <form ref="formRef" @submit.prevent="handleSubmit" class="grid gap-5 pt-6 md:grid-cols-2">
        <BaseInput v-model="form.name" label="Nama Kategori" placeholder="Contoh: Makanan & Minuman"
          required
          :validate="['Nama kategori wajib diisi', validCategoryName]" />
        <div class="md:col-span-2">
          <BaseInput v-model="form.description" label="Deskripsi" placeholder="Contoh: Pengeluaran harian untuk makan" />
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
      <div v-if="filteredCategories.length === 0" class="space-y-3 text-center text-slate-600">
        <p class="text-lg font-medium">Belum ada kategori</p>
        <p class="text-sm">Klik tombol Tambah Kategori untuk menambahkan kategori baru.</p>
      </div>
      <div v-else>
        <!-- Mobile -->
        <div class="md:hidden space-y-3">
          <div v-for="category in filteredCategories" :key="category.id"
            class="rounded-2xl bg-slate-50 p-4 shadow-sm space-y-2">
            <div>
              <p class="font-semibold text-slate-900 text-sm">{{ category.name }}</p>
              <p class="text-xs text-slate-500 mt-0.5">{{ category.description || '-' }}</p>
            </div>
            <div class="text-xs text-slate-400">Dibuat: {{ formatDate(category.created_at) }}</div>
            <div class="flex gap-2 pt-1">
              <button @click="openEditForm(category)"
                class="flex-1 rounded-lg bg-slate-100 px-3 py-2 text-sm text-slate-700 transition hover:bg-slate-200">Edit</button>
              <button @click="handleDelete(category.id)"
                class="flex-1 rounded-lg bg-red-600 px-3 py-2 text-sm font-semibold text-white transition hover:bg-red-700">Hapus</button>
            </div>
          </div>
        </div>
        <!-- Desktop -->
        <div class="hidden md:block overflow-x-auto">
          <table class="min-w-full border-separate border-spacing-y-3 text-left">
            <thead>
              <tr class="text-sm text-slate-500">
                <th class="px-4 py-3">Nama Kategori</th>
                <th class="px-4 py-3">Deskripsi</th>
                <th class="px-4 py-3">Dibuat</th>
                <th class="px-4 py-3">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="category in filteredCategories" :key="category.id"
                class="rounded-3xl bg-slate-50 align-top text-sm shadow-sm transition hover:bg-slate-100">
                <td class="px-4 py-4 text-slate-900">{{ category.name }}</td>
                <td class="px-4 py-4 text-slate-600">{{ category.description || '-' }}</td>
                <td class="px-4 py-4 text-slate-600">{{ formatDate(category.created_at) }}</td>
                <td class="px-4 py-4 space-x-2">
                  <button @click="openEditForm(category)"
                    class="rounded-lg bg-slate-100 px-3 py-1 text-sm text-slate-700 transition hover:bg-slate-200">Edit</button>
                  <button @click="handleDelete(category.id)"
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
