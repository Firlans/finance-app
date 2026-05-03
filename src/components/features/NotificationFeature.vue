<script setup>
import { computed, ref, watch, onUnmounted } from 'vue'

const props = defineProps({
  message: { type: String, default: '' },
  type: { type: String, default: 'success' },
  duration: { type: Number, default: 2000 },
  notification: { type: Object, default: null }
})
const emit = defineEmits(['close'])

const localVisible = ref(false)
let timer = null

const message = computed(() => {
  if (props.notification && props.notification.message !== undefined) {
    return props.notification.message.value
  }
  return props.message
})

const type = computed(() => {
  if (props.notification && props.notification.type !== undefined) {
    return props.notification.type
  }
  return props.type
})

const duration = computed(() => {
  if (props.notification && props.notification.duration !== undefined) {
    return props.notification.duration
  }
  return props.duration
})

const show = computed(() => {
  if (props.notification && props.notification.show !== undefined) {
    return props.notification.show.value
  }
  return localVisible.value
})

const handleClose = () => {
  if (props.notification && props.notification.show !== undefined) {
    // eslint-disable-next-line vue/no-mutating-props
    props.notification.show.value = false
  } else {
    localVisible.value = false
  }

  if (timer) { clearTimeout(timer); timer = null }
  emit('close')
}

watch([message, show], ([messageValue, showValue]) => {
  if (showValue && messageValue) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      handleClose()
      timer = null
    }, duration.value)
  } else {
    if (timer) { clearTimeout(timer); timer = null }
  }
})

watch(message, (val) => {
  if (val) {
    if (props.notification && props.notification.show !== undefined) {
      // eslint-disable-next-line vue/no-mutating-props
      props.notification.show.value = true
    } else {
      localVisible.value = true
    }
  } else {
    if (props.notification && props.notification.show !== undefined) {
      // eslint-disable-next-line vue/no-mutating-props
      props.notification.show.value = false
    } else {
      localVisible.value = false
    }
  }
})

onUnmounted(() => { if (timer) clearTimeout(timer) })
</script>

<template>
  <div v-if="show" class="fixed right-4 top-6 z-50 max-w-sm w-full">
    <div :class="[
      'px-4 py-3 rounded-lg shadow-lg flex items-center space-x-3 text-white',
      type === 'success' ? 'bg-green-600' : 'bg-red-600'
    ]" class="transition-opacity duration-200">
      <svg v-if="props.type === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 flex-shrink-0"
        viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd"
          d="M16.707 5.293a1 1 0 00-1.414 0L8 12.586 4.707 9.293a1 1 0 00-1.414 1.414l4 4a1 1 0 001.414 0l8-8a1 1 0 000-1.414z"
          clip-rule="evenodd" />
      </svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 flex-shrink-0" viewBox="0 0 20 20"
        fill="currentColor">
        <path fill-rule="evenodd"
          d="M10 18a8 8 0 100-16 8 8 0 000 16zm-1-9V6a1 1 0 112 0v3a1 1 0 11-2 0zm0 4a1 1 0 112 0 1 1 0 01-2 0z"
          clip-rule="evenodd" />
      </svg>

      <div class="flex-1 text-sm">
        {{ message }}
      </div>

      <button @click="handleClose" class="text-white opacity-90 hover:opacity-100">
        <span class="text-lg leading-none">×</span>
      </button>
    </div>
  </div>
</template>
