<script setup lang="ts">
import FeaturedGame from './components/FeaturedGame.vue'
import GameCard from './components/GameCard.vue'
</script>

<template>
  <FeaturedGame
    title="Street Fighter III: 3rd Strike"
    platform="Fightcade"
    :release-year="1999"
    developer="Capcom"
    banner-url="https://d1lss44hh2trtw.cloudfront.net/resize?type=webp&url=https%3A%2F%2Fshacknews-www.s3.amazonaws.com%2Fassets%2Farticle%2F2024%2F07%2F29%2F3rd-strike-header_feature.jpg&width=1032&sign=XEMEx9JTW7qcjq1YKxMlyIOL40hjUM0eBKwjqW6KlMU"
  />

  <div class="flex gap-6 p-12">
    <GameCard
      title="Street Fighter III: 3rd Strike"
      platform="Fightcade"
      :release-year="1999"
      developer="Capcom"
      image-url="https://images.igdb.com/igdb/image/upload/t_cover_big/co1wjp.jpg"
      selected
    />
    <GameCard
      title="Metal Slug 3"
      platform="MAME"
      :release-year="2000"
      developer="SNK"
      image-url="https://images.igdb.com/igdb/image/upload/t_cover_big/co1qxk.jpg"
    />
    <GameCard
      title="Guilty Gear XX Accent Core"
      platform="Fightcade"
      :release-year="2002"
      developer="Arc System Works"
      image-url="https://images.igdb.com/igdb/image/upload/t_cover_big/co1y8h.jpg"
    />
  </div>
</template>
