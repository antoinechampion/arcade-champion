<script setup lang="ts">
import { ref } from 'vue'
import ArcadeButton from '@/components/design-system/ArcadeButton.vue'
import { useComponentNavigation, type NavCommand } from '@/composables/navigation'

defineProps<{
  title: string
  platform: string
  releaseYear: number
  developer: string
  bannerUrl: string
}>()

const emit = defineEmits<{
  play: []
}>()

const buttonRef = ref<InstanceType<typeof ArcadeButton> | null>(null)
const sectionRef = ref<HTMLElement | null>(null)

const { active } = useComponentNavigation('featured', {
  onCommand(command: NavCommand) {
    if (command === 'confirm') {
      emit('play')
      return true
    }
    return false
  },
  onEnter() {
    sectionRef.value?.scrollIntoView({ block: 'center', behavior: 'smooth' })
  },
})
</script>

<template>
  <section ref="sectionRef" class="relative h-[60vh] overflow-hidden">
    <div class="absolute inset-0">
      <img :src="bannerUrl" :alt="title" class="banner-img h-full w-full object-cover" />
    </div>

    <div class="relative flex flex-col items-start justify-end h-full p-24 gap-3">
      <span class="text-sm font-medium uppercase tracking-widest opacity-70 ml-0.5">{{ platform }}</span>
      <h1 class="hero-title text-5xl font-bold leading-tight">{{ title }}</h1>
      <span class="text-sm opacity-60 ml-0.5">{{ releaseYear }} - {{ developer }}</span>
      <ArcadeButton ref="buttonRef" class="play-btn mt-9" label="Play" :focused="active" @click="emit('play')">
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
  image-rendering: pixelated;
  object-position: right center;
  mask-image:
    linear-gradient(to right, transparent 0%, black 45%, black 88%, transparent 100%),
    linear-gradient(to top, transparent 0%, black 45%);
  mask-composite: intersect;
  -webkit-mask-image:
    linear-gradient(to right, transparent 0%, black 45%, black 88%, transparent 100%),
    linear-gradient(to top, transparent 0%, black 45%);
  -webkit-mask-composite: source-in;
}

.hero-title {
  text-shadow: 0 0 30px rgba(0, 0, 0, 0.6);
}

:deep(.play-btn.focused) {
  animation: play-pulse 2.4s ease-in-out infinite;
}

@keyframes pan {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-8%);
  }
}

@keyframes play-pulse {
  0%, 100% {
    box-shadow:
      0 0 30px var(--color-glow),
      0 0 60px rgba(124, 92, 224, 0.15),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
  50% {
    box-shadow:
      0 0 45px var(--color-glow-strong),
      0 0 90px rgba(124, 92, 224, 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
}
</style>
