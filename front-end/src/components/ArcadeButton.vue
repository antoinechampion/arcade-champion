<script setup lang="ts">
defineProps<{
  label?: string
}>()
</script>

<template>
  <button class="arcade-btn" :class="label ? 'has-label' : 'icon-only'">
    <span v-if="$slots.icon" class="icon">
      <slot name="icon" />
    </span>
    <span v-if="label" class="label">{{ label }}</span>
  </button>
</template>

<style scoped>
.arcade-btn {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  border: none;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(24px);
  color: var(--color-text);
  font-size: 1rem;
  font-weight: 500;
  letter-spacing: 0.025em;
  cursor: pointer;
  outline: none;
  overflow: hidden;
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
  box-shadow:
    0 0 20px rgba(255, 255, 255, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.12);
}

.arcade-btn.has-label {
  padding: 0.875rem 1.75rem;
}

.arcade-btn.icon-only {
  padding: 1.125rem;
}

.icon,
.label {
  position: relative;
}

.icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Animated gradient blobs */
.arcade-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  --blob1-x: 30%;
  --blob1-y: 80%;
  --blob2-x: 70%;
  --blob2-y: 20%;
  --blob3-x: 60%;
  --blob3-y: 90%;
  background:
    radial-gradient(ellipse 80% 60% at var(--blob1-x) var(--blob1-y), var(--color-primary-dark) 0%, transparent 70%),
    radial-gradient(ellipse 60% 50% at var(--blob2-x) var(--blob2-y), var(--color-primary-light) 0%, transparent 60%),
    radial-gradient(ellipse 70% 55% at var(--blob3-x) var(--blob3-y), var(--color-accent) 0%, transparent 65%);
  transition:
    --blob1-x 0.6s ease,
    --blob1-y 0.6s ease,
    --blob2-x 0.6s ease,
    --blob2-y 0.6s ease,
    --blob3-x 0.6s ease,
    --blob3-y 0.6s ease;
  pointer-events: none;
}

.arcade-btn:focus-visible::before,
.arcade-btn:focus::before {
  --blob1-x: 55%;
  --blob1-y: 60%;
  --blob2-x: 40%;
  --blob2-y: 40%;
  --blob3-x: 45%;
  --blob3-y: 75%;
}

/* Gradient border */
.arcade-btn::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  padding: 1.5px;
  background: conic-gradient(
    from 135deg,
    transparent 0%,
    var(--color-primary-light) 15%,
    transparent 30%,
    var(--color-accent) 50%,
    transparent 65%,
    var(--color-primary-dark) 80%,
    transparent 100%
  );
  mask:
    linear-gradient(#000 0 0) content-box,
    linear-gradient(#000 0 0);
  mask-composite: exclude;
  -webkit-mask:
    linear-gradient(#000 0 0) content-box,
    linear-gradient(#000 0 0);
  -webkit-mask-composite: xor;
  pointer-events: none;
  opacity: 0.4;
  filter: blur(2px);
  transition:
    opacity 0.3s ease,
    filter 0.3s ease,
    background 0.3s ease;
}

.arcade-btn:focus-visible::after,
.arcade-btn:focus::after {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.25),
    var(--color-primary-light),
    var(--color-accent),
    rgba(255, 255, 255, 0.08),
    var(--color-primary-dark),
    rgba(255, 255, 255, 0.2)
  );
  opacity: 0.8;
  filter: blur(0px);
}

.arcade-btn:focus-visible,
.arcade-btn:focus {
  transform: scale(1.08);
  box-shadow:
    0 0 30px var(--color-glow),
    0 0 60px rgba(124, 92, 224, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.arcade-btn:active {
  transform: scale(0.97);
  transition-duration: 0.1s;
}
</style>
