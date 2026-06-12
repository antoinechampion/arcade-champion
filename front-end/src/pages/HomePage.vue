<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import FeaturedGame from '@/components/home/FeaturedGame.vue'
import RecentlyPlayed from '@/components/home/RecentlyPlayed.vue'
import AllGames from '@/components/home/AllGames.vue'
import LaunchOptions from '@/components/home/LaunchOptions.vue'
import LaunchingOverlay from '@/components/design-system/LaunchingOverlay.vue'
import HudFrame from '@/components/design-system/HudFrame.vue'
import { fetchRecentlyPlayed, launchGame, imageUrl } from '@/api/client'
import { usePageNavigation } from '@/composables/navigation'
import type { Game } from '@/api/types'
import type { LaunchMode } from '@/components/home/LaunchOptions.vue'

usePageNavigation(['featured', 'recentlyPlayed', 'allGames'])

const games = ref<Game[]>([])
const pendingGame = ref<Game | null>(null)
const launching = ref(false)

onMounted(async () => {
  games.value = await fetchRecentlyPlayed()
})

const featuredGame = computed(() => games.value[0])
const recentGames = computed(() => games.value.slice(1))

function showLaunching() {
  launching.value = true
  setTimeout(() => { launching.value = false }, 10_000)
}

function handleLaunch(game: Game) {
  if (game.platform === 'Fightcade') {
    pendingGame.value = game
  } else {
    launchGame(game.platform, game.appId)
    showLaunching()
  }
}

function playFeatured() {
  if (featuredGame.value) handleLaunch(featuredGame.value)
}

function onLaunchModeSelected(mode: LaunchMode) {
  if (pendingGame.value) {
    launchGame(pendingGame.value.platform, pendingGame.value.appId, { mode })
    showLaunching()
  }
  pendingGame.value = null
}

function onLaunchCancelled() {
  pendingGame.value = null
}
</script>

<template>
  <HudFrame />

  <FeaturedGame
    v-if="featuredGame"
    :title="featuredGame.title"
    :platform="featuredGame.platform"
    :release-year="featuredGame.releaseYear"
    :developer="featuredGame.developer"
    :banner-url="imageUrl(featuredGame.bannerFilename)"
    @play="playFeatured"
  />

  <RecentlyPlayed v-if="recentGames.length" :games="recentGames" @launch="handleLaunch" />

  <AllGames @launch="handleLaunch" />

  <LaunchOptions
    v-if="pendingGame"
    @select="onLaunchModeSelected"
    @cancel="onLaunchCancelled"
  />

  <LaunchingOverlay v-if="launching" />

  <footer class="flex justify-center py-8">
    <RouterLink to="/backoffice" class="backoffice-link">
      Back Office
    </RouterLink>
  </footer>
</template>

<style scoped>
.backoffice-link {
  color: var(--color-text);
  opacity: 0.4;
  font-size: 0.875rem;
  text-decoration: none;
  transition: opacity 0.2s ease;
}

.backoffice-link:hover,
.backoffice-link:focus-visible {
  opacity: 0.8;
}
</style>
