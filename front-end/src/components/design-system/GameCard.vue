<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  title: string
  platform: string
  releaseYear: number
  developer: string
  imageUrl: string
  selected?: boolean
}>()

const failed = ref(false)
</script>

<template>
  <div class="game-card plasma-border" :class="{ selected, 'plasma-border-active': selected }">
    <div v-if="failed" class="card-fallback">
      <span class="fallback-title">{{ title }}</span>
    </div>
    <img v-else :src="imageUrl" :alt="title" class="card-img w-full h-full object-cover" @error="failed = true" />
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
  background: rgba(255, 255, 255, 0.03);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.06);
  transition: box-shadow, transform 0.3s ease;
}

.card-fallback {
  display: flex;
  align-items: flex-end;
  width: 100%;
  height: 100%;
  padding: 0.75rem;
  border-radius: inherit;
  background:
    radial-gradient(ellipse 80% 60% at 30% 20%, var(--color-primary-dark) 0%, transparent 70%),
    radial-gradient(ellipse 70% 55% at 75% 90%, var(--color-accent) 0%, transparent 70%),
    var(--color-bg);
}

.fallback-title {
  font-size: 0.875rem;
  font-weight: 700;
  line-height: 1.2;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.6);
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
  transition: opacity 0.3s ease;
}

.game-card.selected {
  box-shadow:
    0 0 20px var(--color-glow),
    0 0 40px rgba(124, 92, 224, 0.15);
  transform: scale(1.05);
}

.game-card.selected::after {
  padding: 3px;
}

.game-card.selected .info {
  opacity: 1;
}
</style>
