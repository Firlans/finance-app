<script setup>
const props = defineProps({
  show: { type: Boolean, default: false },
  label: { type: String, default: 'Loading...' }
})
</script>

<template>
  <div v-if="props.show" class="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm px-4">
    <div class="bg-white/95 border border-slate-200 rounded-3xl p-6 flex flex-col items-center gap-4 shadow-xl max-w-sm w-full text-center">
      <div class="w-16 h-16 rounded-full border-4 border-blue-600 border-t-transparent animate-spin"></div>
      <div class="text-sm text-slate-700 font-medium">{{ props.label }}</div>
    </div>
  </div>
</template>
