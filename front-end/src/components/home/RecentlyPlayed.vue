<script setup lang="ts">
import { ref, computed } from 'vue'
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
const containerRef = ref<HTMLElement | null>(null)
const translateX = ref(0)

const listStyle = computed(() => ({
  transform: `translateX(${translateX.value}px)`,
}))

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
  const container = containerRef.value
  if (!container) return
  const clip = container.parentElement
  if (!clip) return
  const clipWidth = clip.clientWidth
  const cards = container.children
  if (!cards.length) return
  const cardWidth = (cards[0] as HTMLElement).offsetWidth
  const style = getComputedStyle(container)
  const gap = parseFloat(style.columnGap) || parseFloat(style.gap) || 0
  const paddingLeft = parseFloat(style.paddingLeft) || 0
  const paddingRight = parseFloat(style.paddingRight) || 0
  const cardCenter = paddingLeft + selectedIndex.value * (cardWidth + gap) + cardWidth / 2
  const target = clipWidth / 2 - cardCenter
  const totalWidth = paddingLeft + cards.length * cardWidth + (cards.length - 1) * gap + paddingRight
  const minTranslate = -(totalWidth - clipWidth)
  translateX.value = Math.max(minTranslate, Math.min(0, target))
}
</script>

<template>
  <section ref="sectionRef" class="recently-played-section py-6 px-12">
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
          :image-url="imageUrl(game.coverFilename)"
          :selected="active && index === selectedIndex"
        />
      </div>
    </div>
  </section>
</template>

<style scoped>
.recently-played-clip {
  overflow: hidden;
  padding: 2rem 0;
  margin: -2rem 0;
}

.recently-played-list {
  display: flex;
  gap: 1.5rem;
  padding: 1.5rem;
  padding-left: 3rem;
  transition: transform 0.3s ease;
  width: max-content;
}
</style>
