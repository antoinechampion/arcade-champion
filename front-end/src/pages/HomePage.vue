<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import FeaturedGame from '@/components/home/FeaturedGame.vue'
import RecentlyPlayed from '@/components/home/RecentlyPlayed.vue'
import AllGames from '@/components/home/AllGames.vue'
import { fetchRecentlyPlayed, launchGame, imageUrl } from '@/api/client'
import { usePageNavigation } from '@/composables/navigation'
import type { Game } from '@/api/types'

usePageNavigation(['featured', 'recentlyPlayed', 'allGames'])

const games = ref<Game[]>([])

onMounted(async () => {
  games.value = await fetchRecentlyPlayed()
})

const featuredGame = computed(() => games.value[0])
const recentGames = computed(() => games.value.slice(1))

function playFeatured() {
  if (featuredGame.value) {
    launchGame(featuredGame.value.platform, featuredGame.value.appId)
  }
}
</script>

<template>
  <FeaturedGame
    v-if="featuredGame"
    :title="featuredGame.title"
    :platform="featuredGame.platform"
    :release-year="featuredGame.releaseYear"
    :developer="featuredGame.developer"
    :banner-url="imageUrl(featuredGame.bannerFilename)"
    @play="playFeatured"
  />

  <RecentlyPlayed v-if="recentGames.length" :games="recentGames" />

  <AllGames />

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
