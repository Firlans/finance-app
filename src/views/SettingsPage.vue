<script setup>
import { ref } from 'vue'
import { ToggleFeature } from '@packages/components'
import AccountsSection from './sections/AccountsSection.vue'
import CategoriesSection from './sections/CategoriesSection.vue'
import LoansSection from './sections/LoansSection.vue'

const activeTab = ref('accounts')
const tabs = [
  { value: 'accounts', label: 'Akun' },
  { value: 'categories', label: 'Kategori' },
  { value: 'loans', label: 'Hutang' }
]
</script>

<template>
  <section class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-slate-900">Pengaturan</h1>
      <p class="text-slate-600 mt-1">Kelola akun dan kategori transaksi Anda.</p>
    </div>

    <!-- Tabs -->
    <div class="w-full md:max-w-xl">
      <ToggleFeature 
        v-model="activeTab" 
        :options="tabs" 
      />
    </div>

    <AccountsSection v-if="activeTab === 'accounts'" />
    <CategoriesSection v-if="activeTab === 'categories'" />
    <LoansSection v-if="activeTab === 'loans'" />
  </section>
</template>
