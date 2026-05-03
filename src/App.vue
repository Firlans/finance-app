<script setup>
import { RouterView, useRoute } from 'vue-router'
import { computed } from 'vue'
import MenuSection from '@/views/sections/MenuSection.vue'

const route = useRoute()
const showMenu = computed(() => !['landing', 'login', 'register'].includes(route.name))
</script>

<template>
  <div class="min-h-screen bg-slate-50 flex">
    <MenuSection v-if="showMenu" />
    <main class="flex-1 p-4 md:p-6">
      <RouterView />
    </main>
  </div>
</template>
