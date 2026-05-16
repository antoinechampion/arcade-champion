<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import GameCard from '@/components/design-system/GameCard.vue'
import ArcadeButton from '@/components/design-system/ArcadeButton.vue'
import ArcadeTextInput from '@/components/design-system/ArcadeTextInput.vue'
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
const focusArea = ref<'search' | 'cards'>('search')
const searchFocusIdx = ref(0)
const currentSectionIdx = ref(0)
const currentCardIdx = ref(0)
const searchBarRef = ref<HTMLElement | null>(null)
const containerRef = ref<HTMLElement | null>(null)

const searchInputFocused = computed(() => active.value && focusArea.value === 'search' && searchFocusIdx.value === 0)
const searchButtonFocused = computed(() => active.value && focusArea.value === 'search' && searchFocusIdx.value === 1)

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
  return active.value && focusArea.value === 'cards' && currentSectionIdx.value === section && currentCardIdx.value === card
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

function handleSearchNav(command: NavCommand): boolean {
  switch (command) {
    case 'right':
      if (searchFocusIdx.value === 0) { searchFocusIdx.value = 1; return true }
      return true
    case 'left':
      if (searchFocusIdx.value === 1) { searchFocusIdx.value = 0; return true }
      return true
    case 'down':
      if (sections.value.length) {
        focusArea.value = 'cards'
        return true
      }
      return false
    case 'up':
      return false
    case 'confirm':
      if (searchFocusIdx.value === 1) { search(); return true }
      return true
    default:
      return false
  }
}

function handleCardsNav(command: NavCommand): boolean {
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
      focusArea.value = 'search'
      return true
    default:
      return false
  }
}

function handleNav(command: NavCommand): boolean {
  return focusArea.value === 'search' ? handleSearchNav(command) : handleCardsNav(command)
}

const { active } = useComponentNavigation('allGames', {
  onCommand: handleNav,
  onEnter(from) {
    if (from === 'up') {
      focusArea.value = 'search'
      searchFocusIdx.value = 0
    }
  },
})

function scrollToFocused() {
  if (focusArea.value === 'search') {
    searchBarRef.value?.scrollIntoView({ block: 'center', behavior: 'smooth' })
    return
  }
  const container = containerRef.value
  if (!container) return
  const sectionEl = container.children[currentSectionIdx.value] as HTMLElement | undefined
  if (!sectionEl) return
  const grid = sectionEl.querySelector('.games-grid') as HTMLElement | null
  if (!grid) return
  const card = grid.children[currentCardIdx.value] as HTMLElement | undefined
  card?.scrollIntoView({ block: 'center', behavior: 'smooth' })
}

watch([focusArea, searchFocusIdx, currentSectionIdx, currentCardIdx], scrollToFocused)

function onSearchKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') search()
}
</script>

<template>
  <section class="all-games-section py-6 px-12">
    <div ref="searchBarRef" class="flex items-center justify-between px-12 mb-6">
      <h2 class="text-lg font-bold opacity-80">All Games</h2>
      <div class="search-bar">
      <ArcadeTextInput
        v-model="searchQuery"
        placeholder="Search games…"
        :focused="searchInputFocused"
        @keydown="onSearchKeydown"
      />
      <ArcadeButton label="Search" :focused="searchButtonFocused" @click="search" />
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


.games-grid {
  display: grid;
  grid-template-columns: repeat(v-bind(COLUMNS), 1fr);
  gap: 1.5rem;
  padding-bottom: 1rem;
}

</style>
