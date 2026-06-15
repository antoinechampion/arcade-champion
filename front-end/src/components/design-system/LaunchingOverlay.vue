<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { setInputLocked } from '@/gamepad'

const props = withDefaults(defineProps<{ message?: string; cancellable?: boolean }>(), {
  message: 'Launching…',
  cancellable: false,
})

const emit = defineEmits<{ cancel: [] }>()

// Gamepad face buttons all surface as ' ' (confirm); Escape is the keyboard back.
const CANCEL_KEYS = new Set(['Escape', ' '])

function onKeydown(e: KeyboardEvent) {
  e.preventDefault()
  e.stopPropagation()
  if (props.cancellable && CANCEL_KEYS.has(e.key)) {
    emit('cancel')
  }
}

onMounted(() => {
  setInputLocked(true)
  window.addEventListener('keydown', onKeydown, { capture: true })
})
onUnmounted(() => {
  setInputLocked(false)
  window.removeEventListener('keydown', onKeydown, { capture: true })
})
</script>

<template>
  <div class="launching-overlay">
    <div class="spinner" />
    <p class="text">{{ message }}</p>
    <p v-if="cancellable" class="hint">Press a button to cancel</p>
  </div>
</template>

<style scoped>
.launching-overlay {
  position: fixed;
  inset: 0;
  z-index: 200;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1.5rem;
  background: rgba(10, 10, 15, 0.85);
  backdrop-filter: blur(12px);
}

.spinner {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: var(--color-primary-light);
  animation: spin 0.8s linear infinite;
}

.text {
  font-size: 1.125rem;
  font-weight: 500;
  opacity: 0.8;
}

.hint {
  font-size: 0.875rem;
  opacity: 0.5;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
