<script setup>
import { RouterView, useRoute } from 'vue-router'
import { computed } from 'vue'
import MenuSection from '@/views/sections/MenuSection.vue'

const route = useRoute()
const showMenu = computed(() => route.name !== 'landing')
</script>

<template>
  <MenuSection v-if="showMenu" />
  <RouterView />
</template>
