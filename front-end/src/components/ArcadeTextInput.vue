<script setup lang="ts">
defineProps<{
  focused?: boolean
  placeholder?: string
}>()

const model = defineModel<string>()
</script>

<template>
  <div class="arcade-input-wrapper plasma-border" :class="{ focused, 'plasma-border-active': focused }">
    <input
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
