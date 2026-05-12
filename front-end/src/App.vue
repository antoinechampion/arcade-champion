<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import FeaturedGame from './components/FeaturedGame.vue'
import RecentlyPlayed from './components/RecentlyPlayed.vue'
import AllGames from './components/AllGames.vue'
import { fetchRecentlyPlayed } from '@/api/client'
import type { Game } from '@/api/types'

const games = ref<Game[]>([])

onMounted(async () => {
  games.value = await fetchRecentlyPlayed()
})

const featuredGame = computed(() => games.value[0])
const recentGames = computed(() => games.value.slice(1))
</script>

<template>
  <FeaturedGame
    v-if="featuredGame"
    :title="featuredGame.title"
    :platform="featuredGame.platform"
    :release-year="featuredGame.releaseYear"
    :developer="featuredGame.developer"
    :banner-url="featuredGame.bannerUrl ?? featuredGame.imageUrl"
  />

  <RecentlyPlayed :games="recentGames" />

  <AllGames />
</template>
