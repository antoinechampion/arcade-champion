<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import GameCard from './GameCard.vue'
import ArcadeButton from './ArcadeButton.vue'
import { fetchAllGames } from '@/api/client'
import type { Game } from '@/api/types'

const games = ref<Game[]>([])
const searchQuery = ref('')
const selectedIndex = ref(0)
const focused = ref(false)
const gridRef = ref<HTMLElement | null>(null)

onMounted(async () => {
  games.value = await fetchAllGames()
})

function getColumnCount(): number {
  const grid = gridRef.value
  if (!grid) return 1
  return getComputedStyle(grid).gridTemplateColumns.split(' ').length
}

async function search() {
  const query = searchQuery.value.trim() || undefined
  games.value = await fetchAllGames(query)
  selectedIndex.value = 0
  await nextTick()
  gridRef.value?.focus()
}

function onKeydown(e: KeyboardEvent) {
  const len = games.value.length
  if (!len) return

  let next = selectedIndex.value
  switch (e.key) {
    case 'ArrowRight':
      next = Math.min(next + 1, len - 1)
      break
    case 'ArrowLeft':
      next = Math.max(next - 1, 0)
      break
    case 'ArrowDown':
      next = Math.min(next + getColumnCount(), len - 1)
      break
    case 'ArrowUp':
      next = Math.max(next - getColumnCount(), 0)
      break
    default:
      return
  }
  e.preventDefault()
  selectedIndex.value = next
}

function onSearchKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') search()
}
</script>

<template>
  <section class="all-games-section py-6 px-12">
    <div class="flex items-center justify-between px-12 mb-6">
      <h2 class="text-lg font-bold opacity-80">All Games</h2>
      <div class="search-bar">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search games…"
        class="search-input"
        @keydown="onSearchKeydown"
      />
      <ArcadeButton label="Search" @click="search" />
      </div>
    </div>

    <div
      ref="gridRef"
      class="games-grid pl-12"
      tabindex="0"
      @keydown="onKeydown"
      @focus="focused = true"
      @blur="focused = false"
    >
      <GameCard
        v-for="(game, index) in games"
        :key="game.id"
        :title="game.title"
        :platform="game.platform"
        :release-year="game.releaseYear"
        :developer="game.developer"
        :image-url="game.imageUrl"
        :selected="focused && index === selectedIndex"
      />
    </div>

    <p v-if="!games.length" class="px-12 opacity-50 text-sm">No games found.</p>
  </section>
</template>

<style scoped>
.search-bar {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.search-input {
  flex: 1;
  max-width: 400px;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(24px);
  color: var(--color-text);
  font-size: 1rem;
  font-family: inherit;
  outline: none;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12);
  transition: box-shadow 0.3s ease;
}

.search-input::placeholder {
  color: var(--color-text);
  opacity: 0.4;
}

.search-input:focus {
  box-shadow:
    0 0 20px var(--color-glow),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.games-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, 200px);
  gap: 1.5rem;
  padding-bottom: 2rem;
  outline: none;
}
</style>
