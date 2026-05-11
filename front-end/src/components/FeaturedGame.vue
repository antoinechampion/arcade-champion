<script setup lang="ts">
import ArcadeButton from './ArcadeButton.vue'

defineProps<{
  title: string
  platform: string
  releaseYear: number
  developer: string
  bannerUrl: string
}>()
</script>

<template>
  <section class="relative h-[60vh] overflow-hidden">
    <div class="absolute inset-0">
      <img :src="bannerUrl" :alt="title" class="banner-img h-full object-cover w-[80%] ml-[30%]" />
    </div>

    <div class="relative flex flex-col items-start justify-end h-full p-24 gap-3">
      <span class="text-sm font-medium uppercase tracking-widest opacity-70 ml-0.5">{{ platform }}</span>
      <h1 class="text-5xl font-bold leading-tight">{{ title }}</h1>
      <span class="text-sm opacity-60 ml-0.5">{{ releaseYear }} - {{ developer }}</span>
      <ArcadeButton class="mt-9" label="Play">
        <template #icon>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-5">
            <path d="M8 5v14l11-7z" />
          </svg>
        </template>
      </ArcadeButton>
    </div>
  </section>
</template>

<style scoped>
.banner-img {
  animation: pan 30s ease-in-out infinite alternate;
  mask-image:
    linear-gradient(to right, transparent 0%, black 40%),
    linear-gradient(to top, transparent 0%, black 40%);
  mask-composite: intersect;
  -webkit-mask-image:
    linear-gradient(to right, transparent 0%, black 40%),
    linear-gradient(to top, transparent 0%, black 40%);
  -webkit-mask-composite: source-in;
}

@keyframes pan {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-16.67%);
  }
}
</style>
