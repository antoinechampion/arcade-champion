<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { lockNavigation, unlockNavigation } from '@/composables/navigation'

const props = defineProps<{
  visible: boolean
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  close: []
}>()

const ROWS = [
  ['1', '2', '3', '4', '5', '6', '7', '8', '9', '0'],
  ['Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'],
  ['A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L', "'"],
  ['Z', 'X', 'C', 'V', 'B', 'N', 'M', ',', '.', '-'],
  ['SPACE', 'BACKSPACE', 'DONE'],
]

const row = ref(0)
const col = ref(0)

const currentKey = computed(() => ROWS[row.value][col.value])

watch(() => props.visible, (visible) => {
  if (visible) {
    row.value = 1
    col.value = 0
  }
})

function onKeydown(e: KeyboardEvent) {
  if (!props.visible) return

  const key = e.key
  if (key === 'ArrowUp' || key === 'ArrowDown' || key === 'ArrowLeft' || key === 'ArrowRight' || key === ' ') {
    e.preventDefault()
    e.stopPropagation()
  }

  switch (key) {
    case 'ArrowUp':
      if (row.value > 0) row.value--
      col.value = Math.min(col.value, ROWS[row.value].length - 1)
      break
    case 'ArrowDown':
      if (row.value < ROWS.length - 1) row.value++
      col.value = Math.min(col.value, ROWS[row.value].length - 1)
      break
    case 'ArrowLeft':
      if (col.value > 0) col.value--
      break
    case 'ArrowRight':
      if (col.value < ROWS[row.value].length - 1) col.value++
      break
    case ' ':
      pressKey()
      break
  }
}

function pressKey() {
  const key = currentKey.value
  if (key === 'DONE') {
    emit('close')
  } else if (key === 'BACKSPACE') {
    emit('update:modelValue', props.modelValue.slice(0, -1))
  } else if (key === 'SPACE') {
    emit('update:modelValue', props.modelValue + ' ')
  } else {
    emit('update:modelValue', props.modelValue + key.toLowerCase())
  }
}

watch(() => props.visible, (visible) => {
  if (visible) {
    lockNavigation()
    window.addEventListener('keydown', onKeydown, true)
  } else {
    unlockNavigation()
    window.removeEventListener('keydown', onKeydown, true)
  }
}, { immediate: true })
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="visible" class="keyboard-overlay">
        <div class="keyboard-input-preview">{{ modelValue }}<span class="cursor">|</span></div>
        <div class="keyboard-grid">
          <div v-for="(keys, rIdx) in ROWS" :key="rIdx" class="keyboard-row">
            <button
              v-for="(key, cIdx) in keys"
              :key="key"
              class="keyboard-key plasma-border"
              :class="{
                active: row === rIdx && col === cIdx,
                'plasma-border-active': row === rIdx && col === cIdx,
                wide: key === 'SPACE' || key === 'BACKSPACE' || key === 'DONE',
              }"
              @click="row = rIdx; col = cIdx; pressKey()"
            >
              <template v-if="key === 'SPACE'">␣</template>
              <template v-else-if="key === 'BACKSPACE'">⌫</template>
              <template v-else-if="key === 'DONE'">Done</template>
              <template v-else>{{ key }}</template>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.keyboard-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(8px);
}

.keyboard-input-preview {
  font-size: 1.5rem;
  color: var(--color-text);
  margin-bottom: 2rem;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  min-width: 300px;
  text-align: center;
}

.cursor {
  animation: blink 1s step-end infinite;
  opacity: 0.7;
}

@keyframes blink {
  50% { opacity: 0; }
}

.keyboard-grid {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.keyboard-row {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
}

.keyboard-key {
  min-width: 3rem;
  height: 3rem;
  border: none;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.08);
  color: var(--color-text);
  font-size: 1.1rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: background 0.15s ease, box-shadow 0.15s ease;
}

.keyboard-key.wide {
  min-width: 6rem;
  padding: 0 1.5rem;
}

.keyboard-key::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background:
    radial-gradient(ellipse 80% 60% at 30% 80%, var(--color-primary-dark) 0%, transparent 70%),
    radial-gradient(ellipse 60% 50% at 70% 20%, var(--color-primary-light) 0%, transparent 60%),
    radial-gradient(ellipse 70% 55% at 60% 90%, var(--color-accent) 0%, transparent 65%);
  opacity: 0;
  transition: opacity 0.2s ease;
  pointer-events: none;
}

.keyboard-key.active::before {
  opacity: 1;
}

.keyboard-key.active {
  background: rgba(255, 255, 255, 0.06);
  box-shadow:
    0 0 30px var(--color-glow),
    0 0 60px rgba(124, 92, 224, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
