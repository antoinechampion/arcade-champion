<script setup lang="ts">
import { ref, watch } from 'vue'
import { openKeyboard } from '@/api/client'

const props = defineProps<{
  focused?: boolean
  placeholder?: string
}>()

const model = defineModel<string>()
const inputRef = ref<HTMLInputElement | null>(null)

watch(() => props.focused, (focused) => {
  if (focused) {
    inputRef.value?.focus()
    openKeyboard()
  } else {
    inputRef.value?.blur()
  }
})
</script>

<template>
  <div class="arcade-input-wrapper plasma-border" :class="{ focused, 'plasma-border-active': focused }">
    <input
      ref="inputRef"
      v-model="model"
      type="text"
      :placeholder="placeholder"
      class="arcade-input"
    />
  </div>
</template>

<style scoped>
.arcade-input-wrapper {
  display: inline-flex;
  border-radius: 8px;
  overflow: hidden;
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12);
}

.arcade-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: inherit;
  background: rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(24px);
  color: var(--color-text);
  font-size: 1rem;
  font-family: inherit;
  outline: none;
}

.arcade-input::placeholder {
  color: var(--color-text);
  opacity: 0.4;
}

.arcade-input-wrapper:focus-within,
.arcade-input-wrapper.focused {
  box-shadow:
    0 0 30px var(--color-glow),
    0 0 60px rgba(124, 92, 224, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}
</style>
