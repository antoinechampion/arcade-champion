<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { setInputLocked } from '@/gamepad'

function swallow(e: KeyboardEvent) {
  e.preventDefault()
  e.stopPropagation()
}

onMounted(() => {
  setInputLocked(true)
  window.addEventListener('keydown', swallow, { capture: true })
})
onUnmounted(() => {
  setInputLocked(false)
  window.removeEventListener('keydown', swallow, { capture: true })
})
</script>

<template>
  <div class="playing" />
</template>

<style scoped>
.playing {
  position: fixed;
  inset: 0;
  background: #000;
}
</style>
