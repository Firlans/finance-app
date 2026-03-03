<template>
  <div class="flex items-center gap-3">
    <button @click="$emit('prev')"
      class="w-9 h-9 flex items-center justify-center rounded-lg border bg-white hover:bg-slate-50"
      aria-label="Previous month">
      ‹
    </button>

    <div class="text-center min-w-[160px]">
      <p class="text-sm text-slate-500">
        <slot></slot>
      </p>
      <p class="font-semibold text-slate-800">
        {{ label }}
      </p>
    </div>

    <button @click="$emit('next')"
      class="w-9 h-9 flex items-center justify-center rounded-lg border bg-white hover:bg-slate-50"
      aria-label="Next month">
      ›
    </button>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true
  }
})

defineEmits(['prev', 'next'])
</script>
