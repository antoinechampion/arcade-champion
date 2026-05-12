<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import GameCard from './GameCard.vue'
import type { Game } from '@/api/types'

const props = defineProps<{
  games: Game[]
}>()

const selectedIndex = ref(0)
const containerRef = ref<HTMLElement | null>(null)
const translateX = ref(0)

const listStyle = computed(() => ({
  transform: `translateX(${translateX.value}px)`,
}))

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
  const container = containerRef.value
  if (!container) return
  const card = container.children[selectedIndex.value] as HTMLElement | undefined
  if (!card) return
  const clipWidth = container.parentElement!.getBoundingClientRect().width
  const cardLeft = card.offsetLeft
  const cardWidth = card.offsetWidth
  const listWidth = container.scrollWidth
  const target = -(cardLeft - clipWidth / 2 + cardWidth / 2)
  const minTranslate = -(listWidth - clipWidth)
  translateX.value = Math.max(minTranslate, Math.min(0, target))
}
</script>

<template>
  <section class="recently-played-section py-6 px-12">
    <h2 class="text-lg font-bold mb-4 opacity-80 px-12">Recently Played</h2>
    <div class="recently-played-clip">
      <div ref="containerRef" class="recently-played-list" :style="listStyle">
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
    </div>
  </section>
</template>

<style scoped>
.recently-played-clip {
  overflow-x: clip;
}

.recently-played-list {
  display: flex;
  gap: 1.5rem;
  padding: 1.5rem;
  transition: transform 0.3s ease;
  width: max-content;
}
</style>
