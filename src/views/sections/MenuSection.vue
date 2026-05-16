<script setup>
import { ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

const sections = [
  { id: 1, name: 'Dashboard', link: '/dashboard', icon: 'D' },
  { id: 2, name: 'Profile', link: '/profile', icon: 'P' },
  { id: 3, name: 'Transactions', link: '/transactions', icon: 'T' },
  { id: 4, name: 'Accounts', link: '/accounts', icon: 'A' },
]

const route = useRoute()
const isOpen = ref(true)
const toggleMenu = () => {
  isOpen.value = !isOpen.value
}
const isActive = (link) => route.path === link
</script>

<template>
  <aside
    :class="['min-h-screen transition-all duration-300 border-r border-slate-800 bg-slate-950 text-slate-100', isOpen ? 'w-64' : 'w-20']">
    <div class="flex items-center justify-between gap-4 px-4 py-4 border-b border-slate-800">
      <div
        :class="['min-w-0 transition-all duration-300', isOpen ? 'flex items-center gap-3' : 'flex flex-col items-center gap-2']">
        <button v-if="!isOpen" @click="toggleMenu" :aria-label="'Buka menu'"
          class="inline-flex flex-none h-11 w-11 items-center justify-center rounded-2xl border border-slate-800 bg-slate-900 text-slate-100 hover:bg-slate-800 transition relative z-10">
          <span class="text-xl">☰</span>
        </button>
        <div class="flex h-11 w-11 items-center justify-center rounded-2xl bg-blue-600 text-base font-bold text-white">
          F
        </div>
        <div v-if="isOpen" class="overflow-hidden transition-all duration-300 min-w-0 max-w-full opacity-100">
          <h1 class="text-lg font-semibold">Finance</h1>
          <p class="text-sm text-slate-400">Control panel</p>
        </div>
      </div>
      <button v-if="isOpen" @click="toggleMenu" :aria-label="'Tutup menu'"
        class="inline-flex flex-none h-11 w-11 items-center justify-center rounded-2xl border border-slate-800 bg-slate-900 text-slate-100 hover:bg-slate-800 transition relative z-10">
        <span class="text-xl">×</span>
      </button>

    </div>

    <nav class="px-2 py-4">
      <ul class="space-y-2">
        <li v-for="section in sections" :key="section.id">
          <RouterLink :to="section.link" class="group flex items-center gap-3 rounded-2xl px-3 py-3 transition-colors"
            :class="isActive(section.link)
              ? 'bg-blue-600 text-white shadow-lg'
              : 'text-slate-200 hover:bg-slate-900 hover:text-white'">
            <span
              class="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-800 text-sm font-semibold text-slate-200 transition group-hover:bg-blue-600">
              {{ section.icon }}
            </span>
            <span v-if="isOpen" class="whitespace-pre transition-all duration-300 opacity-100">
              {{ section.name }}
            </span>
          </RouterLink>
        </li>
      </ul>
    </nav>
  </aside>
</template>
