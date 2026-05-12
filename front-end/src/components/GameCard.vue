<script setup lang="ts">
defineProps<{
  title: string
  platform: string
  releaseYear: number
  developer: string
  imageUrl: string
  selected?: boolean
}>()
</script>

<template>
  <div class="game-card" :class="{ selected }">
    <img :src="imageUrl" :alt="title" class="card-img w-full h-full object-cover" />
    <div class="info">
      <span class="text-xs font-medium uppercase tracking-widest opacity-70">{{ platform }}</span>
      <span class="text-sm font-bold leading-tight">{{ title }}</span>
      <span class="text-xs opacity-60">{{ releaseYear }} - {{ developer }}</span>
    </div>
  </div>
</template>

<style scoped>
.game-card {
  position: relative;
  width: 200px;
  min-width: 200px;
  aspect-ratio: 3 / 4;
  border-radius: 0.75rem;
  overflow: hidden;
  cursor: pointer;
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
}

.card-img {
  mask-image: linear-gradient(to top, rgba(0, 0, 0, 0.3) 0%, black 30%);
  -webkit-mask-image: linear-gradient(to top, rgba(0, 0, 0, 0.3) 0%, black 30%);
  transition: mask-image 0.3s ease, -webkit-mask-image 0.3s ease;
}

.game-card.selected .card-img {
  mask-image: none;
  -webkit-mask-image: none;
}

.game-card .info {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 0.75rem;
  padding-top: 3rem;
  gap: 0.25rem;
  background: linear-gradient(to top, rgba(0, 0, 0, 1) 0%, transparent 150%);
  opacity: 0;
  transform: translateY(0.5rem);
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.game-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  padding: 1.5px;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.25),
    var(--color-primary-light),
    var(--color-accent),
    rgba(255, 255, 255, 0.08),
    var(--color-primary-dark),
    rgba(255, 255, 255, 0.2)
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
  opacity: 0;
  transition: opacity 0.3s ease;
}

.game-card.selected {
  z-index: 1;
  transform: scale(1.20);
  box-shadow:
    0 0 20px var(--color-glow),
    0 0 40px rgba(124, 92, 224, 0.15);
}

.game-card.selected::after {
  opacity: 0.8;
}

.game-card.selected .info {
  opacity: 1;
  transform: translateY(0);
}
</style>
