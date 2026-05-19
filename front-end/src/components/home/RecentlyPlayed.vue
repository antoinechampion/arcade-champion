<script setup lang="ts">
import { ref, nextTick } from 'vue'
import GameCard from '@/components/design-system/GameCard.vue'
import { useComponentNavigation, type NavCommand } from '@/composables/navigation'
import { imageUrl } from '@/api/client'
import type { Game } from '@/api/types'

const props = defineProps<{
  games: Game[]
}>()

const emit = defineEmits<{
  launch: [game: Game]
}>()

const selectedIndex = ref(0)
const sectionRef = ref<HTMLElement | null>(null)
const scrollRef = ref<HTMLElement | null>(null)
const cardRefs = ref<HTMLElement[]>([])

const { active } = useComponentNavigation('recentlyPlayed', {
  onCommand(command: NavCommand) {
    if (command === 'right') {
      if (selectedIndex.value >= props.games.length - 1) return true
      selectedIndex.value++
      scrollToSelected()
      return true
    }
    if (command === 'left') {
      if (selectedIndex.value <= 0) return true
      selectedIndex.value--
      scrollToSelected()
      return true
    }
    if (command === 'confirm') {
      const game = props.games[selectedIndex.value]
      if (game) emit('launch', game)
      return true
    }
    return false
  },
  onEnter() {
    sectionRef.value?.scrollIntoView({ block: 'center', behavior: 'smooth' })
  },
})

function scrollToSelected() {
  nextTick(() => {
    const card = cardRefs.value[selectedIndex.value]
    if (!card) return
    card.scrollIntoView({ inline: 'center', block: 'nearest', behavior: 'smooth' })
  })
}
</script>

<template>
  <section ref="sectionRef" class="recently-played-section py-6 px-12">
    <h2 class="text-lg font-bold mb-4 opacity-80 px-12">Recently Played</h2>
    <div ref="scrollRef" class="recently-played-scroll">
      <div
        v-for="(game, index) in games"
        :key="game.id"
        :ref="(el) => { if (el) cardRefs[index] = el as HTMLElement }"
      >
        <GameCard
          :title="game.title"
          :platform="game.platform"
          :release-year="game.releaseYear"
          :developer="game.developer"
          :image-url="imageUrl(game.coverFilename)"
          :selected="active && index === selectedIndex"
        />
      </div>
    </div>
  </section>
</template>

<style scoped>
.recently-played-scroll {
  display: flex;
  gap: 1.5rem;
  padding: 1.5rem 3rem;
  overflow-x: auto;
  scroll-behavior: smooth;
  scrollbar-width: none;
}

.recently-played-scroll::-webkit-scrollbar {
  display: none;
}
</style>
