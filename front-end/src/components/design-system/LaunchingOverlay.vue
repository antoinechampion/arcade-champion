<script setup lang="ts">
</script>

<template>
  <div class="launching-overlay">
    <div class="spinner" />
    <p class="text">Launching…</p>
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

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
