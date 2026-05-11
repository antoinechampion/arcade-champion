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
      image-url="https://i.pinimg.com/736x/a8/db/46/a8db46f121860572350471b0e4405c32.jpg"
      selected
    />
    <GameCard
      title="Marvel VS Capcom 2"
      platform="MAME"
      :release-year="2000"
      developer="SNK"
      image-url="https://wiki.supercombo.gg/images/thumb/2/29/MVSC2_Cover_Art.jpg/300px-MVSC2_Cover_Art.jpg"
    />
    <GameCard
      title="King of Fighters 98"
      platform="Fightcade"
      :release-year="2002"
      developer="SNK"
      image-url="https://images.igdb.com/igdb/image/upload/t_cover_big/co1y8h.jpg"
    />
    <GameCard
        title="Street Fighter 6"
        platform="Steam"
        :release-year="2025"
        developer="Capcom"
        image-url="https://upload.wikimedia.org/wikipedia/en/thumb/9/94/Street_Fighter_6_box_art.jpg/250px-Street_Fighter_6_box_art.jpg"
    />
    <GameCard
        title="Street Fighter 2: Champion Edition"
        platform="Fightcade"
        :release-year="1992"
        developer="Capcom"
        image-url="https://i.redd.it/2rgdsgr7p3pc1.jpeg"
    />
  </div>
</template>
