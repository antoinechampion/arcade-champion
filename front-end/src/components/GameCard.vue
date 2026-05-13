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
  <div class="game-card plasma-border" :class="{ selected, 'plasma-border-active': selected }">
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
  border-radius: inherit;
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
  border-radius: 0 0 0.75rem 0.75rem;
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

.game-card.selected {
  z-index: 1;
  overflow: visible;
  transform: scale(1.20);
  box-shadow:
    0 0 20px var(--color-glow),
    0 0 40px rgba(124, 92, 224, 0.15);
}

.game-card.selected .info {
  opacity: 1;
  transform: translateY(0);
}
</style>
