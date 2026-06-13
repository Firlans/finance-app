<script setup>
import { ref } from 'vue'
import AccountsSection from './sections/AccountsSection.vue'
import CategoriesSection from './sections/CategoriesSection.vue'
import LoansSection from './sections/LoansSection.vue'

const activeTab = ref('accounts')
const tabs = [
  { id: 'accounts', label: 'Akun' },
  { id: 'categories', label: 'Kategori' },
  { id: 'loans', label: 'Hutang' }
]

const requestTabChange = (nextTab) => {
  if (nextTab === activeTab.value) return
  activeTab.value = nextTab
}
</script>

<template>
  <section class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-slate-900">Pengaturan</h1>
      <p class="text-slate-600 mt-1">Kelola akun dan kategori transaksi Anda.</p>
    </div>

    <!-- Tabs -->
    <div class="flex gap-1 rounded-2xl bg-slate-100 p-1 w-fit">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="requestTabChange(tab.id)"
        class="rounded-xl px-5 py-2 text-sm font-semibold transition"
        :class="activeTab === tab.id ? 'bg-white text-slate-900 shadow' : 'text-slate-500 hover:text-slate-700'"
      >
        {{ tab.label }}
      </button>
    </div>

    <AccountsSection v-if="activeTab === 'accounts'" />
    <CategoriesSection v-if="activeTab === 'categories'" />
    <LoansSection v-if="activeTab === 'loans'" />
  </section>
</template>
