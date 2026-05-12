<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import GameCard from './GameCard.vue'
import ArcadeButton from './ArcadeButton.vue'
import { fetchAllGames } from '@/api/client'
import { useComponentNavigation, type NavCommand } from '@/composables/navigation'
import type { Game } from '@/api/types'

const COLUMNS = 7

interface GameSection {
  letter: string
  games: Game[]
}

const games = ref<Game[]>([])
const searchQuery = ref('')
const currentSectionIdx = ref(0)
const currentCardIdx = ref(0)
const containerRef = ref<HTMLElement | null>(null)

const sections = computed<GameSection[]>(() => {
  const grouped: GameSection[] = []
  for (const game of games.value) {
    const letter = game.title[0].toUpperCase()
    const last = grouped[grouped.length - 1]
    if (last && last.letter === letter) {
      last.games.push(game)
    } else {
      grouped.push({ letter, games: [game] })
    }
  }
  return grouped
})

function isSelected(section: number, card: number) {
  return active.value && currentSectionIdx.value === section && currentCardIdx.value === card
}

function clampCard(section: number, card: number) {
  return Math.min(card, sections.value[section].games.length - 1)
}

onMounted(async () => {
  games.value = await fetchAllGames()
})

async function search() {
  const query = searchQuery.value.trim() || undefined
  games.value = await fetchAllGames(query)
  currentSectionIdx.value = 0
  currentCardIdx.value = 0
  await nextTick()
}

function handleNav(command: NavCommand): boolean {
  if (!sections.value.length) return false

  const col = currentCardIdx.value % COLUMNS
  const row = Math.floor(currentCardIdx.value / COLUMNS)
  const rows = Math.ceil(sections.value[currentSectionIdx.value].games.length / COLUMNS)

  switch (command) {
    case 'right':
      if (currentCardIdx.value < sections.value[currentSectionIdx.value].games.length - 1) {
        currentCardIdx.value++
      } else if (currentSectionIdx.value < sections.value.length - 1) {
        currentSectionIdx.value++
        currentCardIdx.value = 0
      }
      return true
    case 'left':
      if (currentCardIdx.value > 0) {
        currentCardIdx.value--
      } else if (currentSectionIdx.value > 0) {
        currentSectionIdx.value--
        currentCardIdx.value = sections.value[currentSectionIdx.value].games.length - 1
      }
      return true
    case 'down':
      if (row < rows - 1) {
        currentCardIdx.value = clampCard(currentSectionIdx.value, (row + 1) * COLUMNS + col)
        return true
      } else if (currentSectionIdx.value < sections.value.length - 1) {
        currentSectionIdx.value++
        currentCardIdx.value = clampCard(currentSectionIdx.value, col)
        return true
      }
      return false
    case 'up':
      if (row > 0) {
        currentCardIdx.value = (row - 1) * COLUMNS + col
        return true
      } else if (currentSectionIdx.value > 0) {
        currentSectionIdx.value--
        const lastRow = Math.floor((sections.value[currentSectionIdx.value].games.length - 1) / COLUMNS)
        currentCardIdx.value = clampCard(currentSectionIdx.value, lastRow * COLUMNS + col)
        return true
      }
      return false
    default:
      return false
  }
}

const { active } = useComponentNavigation('allGames', {
  onCommand: handleNav,
  onEnter(from) {
    if (from === 'up') {
      currentSectionIdx.value = 0
      currentCardIdx.value = 0
    }
  },
})

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
      ref="containerRef"
      class="pl-12"
    >
      <div v-for="(section, sIdx) in sections" :key="section.letter" class="mb-6">
        <h3 class="text-sm font-bold opacity-50 mb-3">{{ section.letter }}</h3>
        <div class="games-grid" :data-section="section.letter">
          <GameCard
            v-for="(game, i) in section.games"
            :key="game.id"
            :title="game.title"
            :platform="game.platform"
            :release-year="game.releaseYear"
            :developer="game.developer"
            :image-url="game.imageUrl"
            :selected="isSelected(sIdx, i)"
          />
        </div>
      </div>
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
  grid-template-columns: repeat(v-bind(COLUMNS), 1fr);
  gap: 1.5rem;
  padding-bottom: 1rem;
}

</style>
