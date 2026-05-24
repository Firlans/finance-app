<script setup>
import { RouterView, useRoute } from 'vue-router'
import { computed } from 'vue'
import MenuSection from '@/views/sections/MenuSection.vue'

const route = useRoute()
const showMenu = computed(() => !['landing', 'login', 'register'].includes(route.name))
</script>

<template>
  <div class="min-h-screen bg-slate-50 flex w-full overflow-x-hidden">
    <MenuSection v-if="showMenu" />
    <main class="flex-1 min-w-0 p-4 md:p-6 pb-20 md:pb-6 overflow-x-hidden">
      <RouterView />
    </main>
  </div>
</template>
