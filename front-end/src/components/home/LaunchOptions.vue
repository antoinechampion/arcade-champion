<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import ArcadeButton from '@/components/design-system/ArcadeButton.vue'
import { lockNavigation, unlockNavigation } from '@/composables/navigation'

export type LaunchMode = 'online' | 'training' | 'arcade'

const emit = defineEmits<{
  select: [mode: LaunchMode]
  cancel: []
}>()

const options: { mode: LaunchMode; label: string }[] = [
  { mode: 'online', label: 'Online' },
  { mode: 'training', label: 'Training' },
  { mode: 'arcade', label: 'Arcade' },
]

const focusedIndex = ref(0)

const NAV_KEYS = new Set(['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', ' ', 'Escape'])

function onKeydown(e: KeyboardEvent) {
  if (!NAV_KEYS.has(e.key)) return
  e.preventDefault()
  e.stopPropagation()
  e.stopImmediatePropagation()

  switch (e.key) {
    case 'ArrowLeft':
      if (focusedIndex.value > 0) focusedIndex.value--
      break
    case 'ArrowRight':
      if (focusedIndex.value < options.length - 1) focusedIndex.value++
      break
    case ' ':
      emit('select', options[focusedIndex.value].mode)
      break
    case 'Escape':
      emit('cancel')
      break
    case 'ArrowUp':
    case 'ArrowDown':
      break
  }
}

onMounted(() => {
  lockNavigation()
  window.addEventListener('keydown', onKeydown, { capture: true })
})

onUnmounted(() => {
  unlockNavigation()
  window.removeEventListener('keydown', onKeydown, { capture: true })
})
</script>

<template>
  <div class="launch-overlay">
    <div class="launch-panel">
      <h2 class="title">Select Mode</h2>
      <div class="options">
        <ArcadeButton
          v-for="(opt, i) in options"
          :key="opt.mode"
          :label="opt.label"
          :focused="i === focusedIndex"
          @click="emit('select', opt.mode)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.launch-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(10, 10, 15, 0.85);
  backdrop-filter: blur(12px);
}

.launch-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2.5rem;
}

.title {
  font-size: 1.5rem;
  font-weight: 600;
  opacity: 0.9;
}

.options {
  display: flex;
  gap: 1.5rem;
}
</style>
