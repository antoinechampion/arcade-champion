<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const time = ref('')

function updateTime() {
  time.value = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

let timer: ReturnType<typeof setInterval>

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 10_000)
})

onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div class="hud-frame" aria-hidden="true">
    <header class="hud-bar">
      <span class="clock">{{ time }}</span>
    </header>

    <span class="corner corner-tl"></span>
    <span class="corner corner-tr"></span>
    <span class="corner corner-bl"></span>
    <span class="corner corner-br"></span>
  </div>
</template>

<style scoped>
.hud-frame {
  position: fixed;
  inset: 0;
  z-index: 50;
  pointer-events: none;
}

.hud-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 1rem 3.5rem;
}

.clock {
  font-size: 0.875rem;
  font-weight: 600;
  letter-spacing: 0.15em;
  opacity: 0.6;
  font-variant-numeric: tabular-nums;
}

.corner {
  position: absolute;
  width: 1.75rem;
  height: 1.75rem;
  border: 2px solid var(--color-primary-light);
  opacity: 0.3;
}

.corner-tl {
  top: 1rem;
  left: 1rem;
  border-right: none;
  border-bottom: none;
}

.corner-tr {
  top: 1rem;
  right: 1rem;
  border-left: none;
  border-bottom: none;
}

.corner-bl {
  bottom: 1rem;
  left: 1rem;
  border-right: none;
  border-top: none;
}

.corner-br {
  bottom: 1rem;
  right: 1rem;
  border-left: none;
  border-top: none;
}
</style>
