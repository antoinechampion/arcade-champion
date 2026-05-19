<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import ArcadeButton from '@/components/design-system/ArcadeButton.vue'
import type { NavCommand } from '@/composables/navigation'

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

const KEY_MAP: Record<string, NavCommand> = {
  ArrowLeft: 'left',
  ArrowRight: 'right',
  ' ': 'confirm',
  Escape: 'up',
}

function onKeydown(e: KeyboardEvent) {
  const command = KEY_MAP[e.key]
  if (!command) return
  e.preventDefault()
  e.stopPropagation()

  switch (command) {
    case 'left':
      if (focusedIndex.value > 0) focusedIndex.value--
      break
    case 'right':
      if (focusedIndex.value < options.length - 1) focusedIndex.value++
      break
    case 'confirm':
      emit('select', options[focusedIndex.value].mode)
      break
    case 'up':
      emit('cancel')
      break
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown, { capture: true }))
onUnmounted(() => window.removeEventListener('keydown', onKeydown, { capture: true }))
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
