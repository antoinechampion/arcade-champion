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
import { useRouter } from 'vue-router'
import type { Game } from '@/api/types'
import type { LaunchMode } from '@/components/home/LaunchOptions.vue'

usePageNavigation(['featured', 'recentlyPlayed', 'allGames'])

const router = useRouter()

type LaunchPhase = 'matchmaking' | 'launching'

const games = ref<Game[]>([])
const pendingGame = ref<Game | null>(null)
const phase = ref<LaunchPhase | null>(null)
let matchmakingAbort: AbortController | null = null

onMounted(async () => {
  games.value = await fetchRecentlyPlayed()
})

const featuredGame = computed(() => games.value[0])
const recentGames = computed(() => games.value.slice(1))

const overlayMessage = computed(() =>
  phase.value === 'matchmaking' ? 'Matchmaking in progress…' : 'Launching…',
)

function showLaunching() {
  phase.value = 'launching'
  setTimeout(() => router.push('/playing'), 10_000)
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

// Online matchmaking blocks until the server pairs us and the game starts, so we
// keep the matchmaking overlay up for the whole request, then switch to launching.
async function launchOnline(game: Game) {
  phase.value = 'matchmaking'
  matchmakingAbort = new AbortController()
  try {
    await launchGame(game.platform, game.appId, { mode: 'online' }, matchmakingAbort.signal)
    showLaunching()
  } catch {
    // Aborted by the user or a real failure — either way return to idle.
    phase.value = null
  } finally {
    matchmakingAbort = null
  }
}

function cancelMatchmaking() {
  matchmakingAbort?.abort()
}

function onLaunchModeSelected(mode: LaunchMode) {
  const game = pendingGame.value
  pendingGame.value = null
  if (!game) return
  if (mode === 'online') {
    launchOnline(game)
  } else {
    launchGame(game.platform, game.appId, { mode })
    showLaunching()
  }
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

  <LaunchingOverlay
    v-if="phase"
    :message="overlayMessage"
    :cancellable="phase === 'matchmaking'"
    @cancel="cancelMatchmaking"
  />

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
