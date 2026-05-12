<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import GameCard from './GameCard.vue'
import type { Game } from '@/api/types'

const props = defineProps<{
  games: Game[]
}>()

const selectedIndex = ref(0)
const containerRef = ref<HTMLElement | null>(null)

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowRight') {
    selectedIndex.value = Math.min(selectedIndex.value + 1, props.games.length - 1)
    scrollToSelected()
  } else if (e.key === 'ArrowLeft') {
    selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
    scrollToSelected()
  }
}

function scrollToSelected() {
  nextTick(() => {
    const container = containerRef.value
    if (!container) return
    const card = container.children[selectedIndex.value] as HTMLElement | undefined
    if (!card) return
    card.scrollIntoView({ behavior: 'smooth', inline: 'center', block: 'nearest' })
  })
}
</script>

<template>
  <section class="py-6 px-12">
    <h2 class="text-lg font-bold mb-4 opacity-80 px-12">Recently Played</h2>
    <div ref="containerRef" class="recently-played-list">
      <GameCard
        v-for="(game, index) in games"
        :key="game.id"
        :title="game.title"
        :platform="game.platform"
        :release-year="game.releaseYear"
        :developer="game.developer"
        :image-url="game.imageUrl"
        :selected="index === selectedIndex"
      />
    </div>
  </section>
</template>

<style scoped>
.recently-played-list {
  display: flex;
  gap: 1.5rem;
  overflow-x: auto;
  padding: 1.5rem 0;
  scrollbar-width: none;
}

.recently-played-list::before,
.recently-played-list::after {
  content: '';
  min-width: 1.5rem;
  flex-shrink: 0;
}

.recently-played-list::-webkit-scrollbar {
  display: none;
}
</style>
